# smartmeter-reader
This project houses code used to read standard smartmeter in Netherlands.

# Introduction
If you stay in Netherlands, then your energy provider has propbably installed new smartmeter at your house. This smartmeter has a P1 port, using which you can extract energy readings (Electricity and Gas) every 10 seconds. The readings can be used further for your purpose. This repo has code to send the data to Graphite database via MQTT server.

The diagram shows high level picture of how the code is deployed

![alt text](https://github.com/cmyrsh/smartmeter-reader/blob/master/smartmeter_reader.jpg "Diagram")

The repo contains 2 modules, Serial Reader and Graphite Feeder

## Serial Reader
This module is a go program. It reads the P1 Telegram from smartmeter and creates a JSON message. After creating the message, it will send it to MQTT server.

### Building

Serial Reader needs to be compiled for target platform. Following code builds for ARM devices (example: Raspberry Pi / C.H.I.P etc)
`env GOARM=7;GOOS=linux;GOARCH=arm go build -v -o smartmeter_reader_arm`

## Graphite Feeder
This is a node in node-red. This node reads P1 JSON message from MQTT topic and extract Smartmeter data. Later it sends data to Graphite instance.



