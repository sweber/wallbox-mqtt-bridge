# MQTT Bridge for Wallbox (Enhanced for EVCC)
This is a modified version of Fredrik Gustafssons **awesome** Wallbox Pulsar Plus mqtt bridge (https://github.com/jagheterfredrik/).

The changes here where made to use the mqtt bridge in conjunction with evcc (https://evcc.io).

To setup the Pulsar Plus with evcc follow the steps below:
1. Root your Wallbox and setup the bridge (see below)
2. Setup mqtt in evcc (https://docs.evcc.io/docs/reference/configuration/mqtt)
3. Extend the evcc configuration with following meter/charger/loadpoint (don't forget to replace the wallbox id placeholder)
```yaml
meters:
- name: wallboxpulsarmeter
  type: custom
  power:
    source: mqtt
    topic: wallbox_<your_wallbox_id>/charging_power/state
    timeout: 180s
  currents:
    - source: mqtt
      topic: wallbox_<your_wallbox_id>/charging_current_l1/state
    - source: mqtt
      topic: wallbox_<your_wallbox_id>/charging_current_l2/state
    - source: mqtt
      topic: wallbox_<your_wallbox_id>/charging_current_l3/state

chargers:
- name: wallboxpulsarcharger
  type: custom
  status:
    source: mqtt
    topic: wallbox_147293/control_pilot/state
  enabled: # charger enabled state (true/false or 0/1)
    source: mqtt
    topic: wallbox_147293/charging_enable/state
  enable: # set charger enabled state (true/false or 0/1)
    source: mqtt
    topic: wallbox_147293/charging_enable/set
    payload: ${enable:%d}
  maxcurrent: # set charger max current (A)
    source: mqtt
    topic: wallbox_147293/max_charging_current/set


loadpoints:
- title: wallboxpulsar
  charger: wallboxpulsarcharger
  meter: wallboxpulsarmeter
  mode: pv
```


<br><br>
**==Original Readme.md below:==**
***
# MQTT Bridge for Wallbox

This open-source project connects your Wallbox fully locally to Home Assistant, providing you with unparalleled speed and reliability.

Note: Doesn't work with firmware v6.6.x see issue https://github.com/jagheterfredrik/wallbox-mqtt-bridge/issues/63

## Features

- **Instant Sensor Data:** The Wallbox's internal state is polled every second and any updates are immediately pushed to the external MQTT broker.

- **Instant Control:** Quickly lock/unlock, pause/resume or change the max charging current, without involving the manufacturer's servers.

- **Always available:** As long as your local network is up and your Wallbox has power, you're in control! No need to rely on a third party to communicate with the device you own.

- **Home Assistant MQTT Auto Discovery:** Enjoy a hassle-free setup with Home Assistant MQTT Auto Discovery support. The integration effortlessly integrates with your existing Home Assistant environment.

<br/>
<p align="center">
   <img src="https://github.com/jagheterfredrik/wallbox-mqtt-bridge/assets/9987465/06488a5d-e6fe-4491-b11d-e7176792a7f5" height="507" />
</p>

## Getting Started

1. [Root your Wallbox](https://github.com/jagheterfredrik/wallbox-pwn)
2. Setup an MQTT Broker, if you don't already have one. Here's an example [installing it as a Home Assistant add-on](https://www.youtube.com/watch?v=dqTn-Gk4Qeo)
3. `ssh` to your Wallbox and run

```sh
curl -sSfL https://github.com/sweber/wallbox-mqtt-bridge/releases/tag/v20241029_1/install.sh > install.sh && bash install.sh
```

Note: To upgrade to new version, simply run the command from step 3 again.

## Acknowledgments

A big shoutout to [@tronikos](https://github.com/tronikos) for their valuable contributions. This project wouldn't be the same without the collaborative spirit of the open-source community.
