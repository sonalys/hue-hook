# HUE-HOOK

hue-hook is a cli tool that can be used to receive plex webhooks and turn groups / lights on / off whenever
you are playing something.

# Usage

You can either input your bridge ip / user yourself or let the cli register itself and generate the `config.yaml` file.

You can run `hue-hook -lights` to get all light ids registered on your bridge.

You need to run the program once, let the plex webhook call it with the players UUID you want to use, then put then in config file with the lamps you want to vinculate.

# Building

To build it simply run `make build`. for Raspberry pi and ARMv7 architecture you can use `make build_arm`.