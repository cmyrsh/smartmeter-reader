package main

import (
	"flag"
	"sync"

	"./mqttadapterY"

	"./datamodel"
	"./serialport"
)

const (
	CHANNEL_SIZE int = 100
)

func main() {

	var usb_port, mqtt_address, mqtt_topic, mqtt_cred_file string
	var interval int

	flag.StringVar(&usb_port, "usb", "/dev/ttyUSB0", "USB/Serial port")
	flag.StringVar(&mqtt_address, "mqtt_address", "localhost:1883", "MQTT Host and Port")
	flag.StringVar(&mqtt_topic, "mqtt_topic", "dev.sample.topic", "Topic name where Serial data needs to be sent")
	flag.IntVar(&interval, "interval", 60, "message interval in seconds. default is 60")
	flag.StringVar(&mqtt_cred_file, "mqtt_cred_file", "", "path of mqtt cred file default is blank, file should be in properties format")

	flag.Parse()

	var wg sync.WaitGroup

	wg.Add(2)

	telegram_channel := make(chan datamodel.P1Telegram, CHANNEL_SIZE)

	go serialport.ReadSerial(usb_port, telegram_channel, wg)

	go mqttadapterY.MQTTy(mqtt_address, "smartmeter-reader", mqtt_topic, mqtt_cred_file, telegram_channel, wg, interval)

	wg.Wait()

}
