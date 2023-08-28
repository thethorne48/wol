[![gowol](https://snapcraft.io/gowol/badge.svg)](https://snapcraft.io/gowol)

# GoWoL

Wake-On-Lan simple Go server


# Usage

This is a basic wrapper for creating and sending "Magic Packets" that can wake devices on a local network as long as the given device allows it. 

It works by allowing an always-on server (like a raspberry pi) to run this service, then if you have a peer-to-site VPN running (i.e. WireGuard) you can send a simple HTTP request to your server that sends a "Magic Packet" with the sent MAC address. 

```bash

```