package bridge

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jagheterfredrik/wallbox-mqtt-bridge/app/wallbox"
)

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	panic("Connection to MQTT lost")
}

func RunBridge(configPath string) {
	c := LoadConfig(configPath)
	w := wallbox.New()
	w.RefreshData()

	serialNumber := w.SerialNumber()
	entityConfig := getEntities(w)
	if c.Settings.DebugSensors {
		for k, v := range getDebugEntities(w) {
			entityConfig[k] = v
		}
	}

	if c.Settings.PowerBoostEnabled {
		for k, v := range getPowerBoostEntities(w, c) {
			entityConfig[k] = v
		}
	}

	topicPrefix := "wallbox_" + serialNumber
	availabilityTopic := topicPrefix + "/availability"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", c.MQTT.Host, c.MQTT.Port))
	opts.SetUsername(c.MQTT.Username)
	opts.SetPassword(c.MQTT.Password)
	opts.SetWill(availabilityTopic, "offline", 1, true)
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for key, val := range entityConfig {
		component := val.Component
		uid := serialNumber + "_" + key
		config := map[string]interface{}{
			"~":                  topicPrefix + "/" + key,
			"availability_topic": availabilityTopic,
			"state_topic":        "~/state",
			"unique_id":          uid,
			"device": map[string]string{
				"identifiers": serialNumber,
				"name":        c.Settings.DeviceName,
			},
		}
		if val.Setter != nil {
			config["command_topic"] = "~/set"
		}
		for k, v := range val.Config {
			config[k] = v
		}
		jsonPayload, _ := json.Marshal(config)
		token := client.Publish("homeassistant/"+component+"/"+uid+"/config", 1, true, jsonPayload)
		token.Wait()
	}

	token := client.Publish(availabilityTopic, 1, true, "online")
	token.Wait()

	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		field := strings.Split(msg.Topic(), "/")[1]
		payload := string(msg.Payload())
		setter := entityConfig[field].Setter
		fmt.Println("Setting", field, payload)
		setter(payload)
	}

	topic := topicPrefix + "/+/set"
	client.Subscribe(topic, 1, messageHandler)

	ticker := time.NewTicker(time.Duration(c.Settings.PollingIntervalSeconds) * time.Second)
	defer ticker.Stop()

	published := make(map[string]interface{})
	lastUpdateIntervalUpdatedTopics := make(map[string]time.Time)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			w.RefreshData()
			for key, val := range entityConfig {
				payload := val.Getter()
				bytePayload := []byte(fmt.Sprint(payload))
				isIntervalUpdatedTopic := slices.Index(c.Settings.IntervalUpdatedTopics, key) > -1
				lastUpdateTimeForTopic := lastUpdateIntervalUpdatedTopics[key]
				isUpdateForTopicNecessary := isUpdateForTopicNecessary(lastUpdateTimeForTopic, c)
				if c.Settings.VerboseOutput {
					fmt.Printf("Check if publishing of '%s' is necessary. (isIntervalUpdatedTopic: '%t', isUpdateForTopicNecessary: '%t', lastUpdateTimeForTopic: '%s')\n",
						key, isIntervalUpdatedTopic, isUpdateForTopicNecessary, lastUpdateTimeForTopic.String())
				}
				if isIntervalUpdatedTopic && isUpdateForTopicNecessary || published[key] != payload {
					if val.RateLimit != nil && !val.RateLimit.Allow(strToFloat(payload)) {
						if c.Settings.VerboseOutput {
							fmt.Printf("Rate limit exceed for key '%s' and payload '%s'\n", key, payload)
						}
						continue
					}
					fmt.Println("Publishing: ", key, payload)
					token := client.Publish(topicPrefix+"/"+key+"/state", 1, true, bytePayload)
					token.Wait()
					published[key] = payload
					lastUpdateIntervalUpdatedTopics[key] = time.Now()
				}
			}
		case <-interrupt:
			fmt.Println("Interrupted. Exiting...")
			token := client.Publish(availabilityTopic, 1, true, "offline")
			token.Wait()
			client.Disconnect(250)
			return
		}
	}
}

func isUpdateForTopicNecessary(lastUpdateTimeForTopic time.Time, config *WallboxConfig) bool {
	if lastUpdateTimeForTopic.IsZero() {
		return true
	}
	lastUpdatedBeforeForTopic := time.Since(lastUpdateTimeForTopic)
	return int(lastUpdatedBeforeForTopic.Seconds()) > (*config).Settings.IntervalUpdatedTopicsSeconds
}
