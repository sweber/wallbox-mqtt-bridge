package bridge

import (
	"gopkg.in/ini.v1"
)

type WallboxConfig struct {
	MQTT struct {
		Host     string `ini:"host"`
		Port     int    `ini:"port"`
		Username string `ini:"username"`
		Password string `ini:"password"`
	} `ini:"mqtt"`

	Settings struct {
		PollingIntervalSeconds       int      `ini:"polling_interval_seconds"`
		DeviceName                   string   `ini:"device_name"`
		DebugSensors                 bool     `ini:"debug_sensors"`
		PowerBoostEnabled            bool     `ini:"power_boost_enabled"`
		IntervalUpdatedTopics        []string `ini:"interval_updated_topics"`
		IntervalUpdatedTopicsSeconds int      `ini:"interval_updated_topics_seconds"`
		VerboseOutput                bool     `ini:"verbose_output"`
	} `ini:"settings"`
}

func (w *WallboxConfig) SaveTo(path string) {
	cfg := ini.Empty()
	cfg.ReflectFrom(w)
	cfg.SaveTo(path)
}

func LoadConfig(path string) *WallboxConfig {
	cfg, _ := ini.Load(path)

	var config WallboxConfig
	if err := cfg.MapTo(&config); err != nil {
		return nil
	}

	return &config
}
