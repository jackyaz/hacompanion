# Desktop Companion for Home Assistant

This was forked from [tobias-kuendig/hacompanion](https://github.com/tobias-kuendig/hacompanion)

This is an unofficial Desktop Companion App for [Home Assistant](https://www.home-assistant.io/) written in Go.

The companion is running as a background process and sends local hardware information to your Home Assistant instance.
Additionally, you can send notifications from Home Assistant to your Computer and display them using `notify-send`.

Currently, **Linux** is the only supported operating system.


## Supported sensors

* CPU temperature
* CPU usage
* Load average
* Memory usage
* Uptime
* Power stats
* Online check
* Audio volume
* Webcam process count
* Custom scripts

## Installation

You can build this from source by running the below commands:

```shell
git clone https://github.com/jackyaz/hacompanion.git
cd hacompanion
go build
sudo cp hacompanion /usr/local/bin/hacompanion
sudo chown root:root /usr/local/bin/hacompanion
```

You can now start the companion with the `hacompanion` command. But before doing so, you have to set up
the configuration:

## Configuration and Setup

1. Make sure you have the [Mobile App integration](https://www.home-assistant.io/integrations/mobile_app/) enabled in Home Assistant (it is on by default).
1. Download a copy of the [configuration file](hacompanion.toml). Save it to `~/.config/hacompanion.toml`.
    ```shell
    sudo mkdir -p /etc/hacompanion
    sudo wget -O /etc/hacompanion/hacompanion.toml https://raw.githubusercontent.com/jackyaz/hacompanion/main/hacompanion.toml
    ```
1. In Home Assistant, generate a token by
   visting [your profile page](https://www.home-assistant.io/docs/authentication/#your-account-profile), then click on `Generate Token` at
   the end of the page.
1. Update your /etc/hacompanion/hacompanion.toml` file's `[homeassistant]` section with the generated `token`.
1. Set the display name of your device (`device_name`) and the URL of your Home Assistant instance (`host`).
1. To receive notifications on a specific IP address you may need to change the 
`push_url` and `listen` settings under `[notifications]` to point respectively 
to your local IP address and the listen port. Without any value hacompanion will 
use your default NIC and listen on port `8080`.
1. Configure all sensors in the configuration file as you see fit.

## Run the companion on system boot

If your system is using Systemd, you can use the following unit file to run the companion on system boot:

```shell
sudo ${EDITOR:-nano} "/etc/systemd/system/hacompanion.service"
```

```ini
[Unit]
Description=Home Assistant Desktop Companion
Documentation=https://github.com/jackyaz/hacompanion
# Uncomment the line below if you are using NetworkManager to ensure hacompanion
# only starts after your host is online
# After=NetworkManager-wait-online.service

[Service]
ExecStart=/usr/local/bin/hacompanion -config=/etc/hacompanion/hacompanion.toml
Restart=on-failure
RestartSec=5
Type=simple

[Install]
WantedBy=default.target
```

Start the companion by running:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now hacompanion
# check status with:
# sudo systemctl status hacompanion
# and logs with:
# sudo journalctl -xlf -u hacompanion
```

You should now see your new sensors under `Settings -> Integrations -> Mobile App -> Your Device`.

## Custom scripts

You can add any number of custom scripts in your configuration file.

The companion will call these scripts and send the output to Home Assistant. It does not matter what language the script is written in, as
long as it can be executed from the command line.

The output of your script has to be as follows:

```
my_state_value
icon:mdi:home-assistant
custom_attribute_1:value 1
custom_attribute_2:value 2
```

The above would be translated to the following json payload:

```json
{
  "icon": "mdi:home-assistant",
  "state": "my_state_value",
  "attributes": {
    "custom_attribute_1": "value 1",
    "custom_attribute_2": "value 2"
  }
}
```

The state (first line) is required.
If `icon` is not set then the icon defined in the config file will be used.
Attributes are optional.

### Example script

The following bash script reports the current time to Home Assistant.

It can be registered like this:

```toml
[script.custom_time]
path = "/path/to/script.sh"
name = "The current time"
icon = "mdi:clock-outline"
type = "sensor"
```

The script content:

```bash
#!/bin/bash
date "+%H:%M"             # First line, state of the sensor
echo formatted:$(date)    # Custom "formatted" Attribute
echo unix:$(date "+%s")   # Custom "unix" Attribute
```

The output:

```text
16:34
formatted:Sa 15 Mai 2021 16:34:40 CEST
unix:1621089280
```

## Receiving notifications

The companion can receive notifications from Home Assistant and display them using `notify-send`. To test the integration, start the companion
and execute the following service in Home Assistant:

```yaml
service: notify.mobile_app_your_device # change this!
data:
  title: "Message Title"
  message: "Message Body"
  data:
    expire: 4000 # display for 4 seconds
    urgency: normal
```

## Automation ideas

Feel free to share your automation ideas [in the Discussions section](https://github.com/jackyaz/hacompanion/discussions) of this
repo.
