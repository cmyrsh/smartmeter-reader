# smartmeter-reader
This project houses code used to read standard smartmeter in Netherlands.

# Introduction
If you stay in the Netherlands, then your energy provider has propbably installed a smartmeter at your house. This smartmeter has a P1 port, using which you can extract energy readings (Electricity and Gas) every 10 seconds. The readings can be used further for your purpose. 

This repo has 

The diagram shows high level picture of how the code is deployed
<!-- 
![alt text](https://github.com/cmyrsh/smartmeter-reader/blob/master/smartmeter_reader.jpg "Diagram") -->

```mermaid
 flowchart LR
    SmartMeter(P1 Port of Smart Meter)
    Pi(RaspBerry Pi USB Port)
    MQTT(MQTT Server)
    SmartMeter --> Pi
    Pi -- Json --> MQTT
```

## Serial Reader
This module is a go program. It reads the P1 Telegram from smartmeter and creates a JSON message. After creating the message, it will send it to MQTT server.

### Building

Serial Reader needs to be compiled for target platform. Following code builds for ARM devices (example: Raspberry Pi / C.H.I.P etc)
```{r, engine='bash', count_lines}
./build_arm.sh
```
To compile on Linux we use 
```{r, engine='bash', count_lines}
./build_linux.sh
```

### Running
To know all options, run 
```{r, engine='bash', count_lines}
smartmeter_reader_arm -h
```
#### Options
```{r, engine='bash', count_lines}
Usage of ./smartmeter_reader_arm:
  -interval int
        message interval in seconds. default is 60 (default 60)
  -mqtt_address string
        MQTT Host and Port (default "localhost:1883")
  -mqtt_cred_file string
        path of mqtt cred file default is blank, file should be in properties format
  -mqtt_topic string
        Topic name where Serial data needs to be sent (default "dev.sample.topic")
  -usb string
        USB/Serial port (default "/dev/ttyUSB0")
```


