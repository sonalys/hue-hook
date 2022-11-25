# HUE-HOOK

hue-hook is a cli tool that can be used to receive plex webhook and switch lights on / off whenever
you are playing something.

# Usage

You can either input your bridge ip / user yourself or let the cli register itself and generate the `config.yaml` file.

You can run `hue-hook -lights` to get all light ids registered on your bridge.

You need to run the program once, let the plex webhook call it with the player's UUID. copy it and then put it in the config file with the lamps you want to vinculate.

## Example config.yaml

```yaml
user: QH5KQLBOCHv50BqKwMoNoZD8tHjzWCEnJ2cfRF0v
bridgeHost: 192.168.1.111
port: 3304
players:
  b3kcaxv4szkm6tejy85qyibs: [1,3,6]
```

# Building

To build it simply run `make build`. for Raspberry pi and ARMv7 architecture you can use `make build_arm`.