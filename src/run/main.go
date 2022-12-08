package main

import (
	"flag"
	"io/ioutil"
	"log"
	"sync"

	"gopkg.in/yaml.v2"

	"com.github/cmyrsh/sm-reader/datamodel"
	"com.github/cmyrsh/sm-reader/mqttadapter"
	"com.github/cmyrsh/sm-reader/p1port"
)

const (
	CHANNEL_SIZE int = 100
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var usb_port, mqtt_address, mqtt_topic, mqtt_cred_file, config_path string
	var interval int
	var debug bool

	flag.BoolVar(&debug, "debug", false, "Debug Flag. If true, all the data will be printed to console")
	flag.StringVar(&usb_port, "usb", "/dev/ttyUSB0", "USB/Serial port")
	flag.StringVar(&mqtt_address, "mqtt_address", "localhost:1883", "MQTT Host and Port")
	flag.StringVar(&mqtt_topic, "mqtt_topic", "dev.sample.topic", "Topic name where Serial data needs to be sent")
	flag.IntVar(&interval, "interval", 60, "message interval in seconds. default is 60")
	flag.StringVar(&mqtt_cred_file, "mqtt_cred_file", "", "path of mqtt cred file default is blank, file should be in properties format")
	flag.StringVar(&config_path, "config", "./smartmeter-config.yml", "path of config yaml file. see sample attached")

	flag.Parse()

	config, error_config := readConfig(config_path)

	if error_config == nil {
		log.Printf("Using values from Config %s\n", config_path)
		usb_port = config.USBPort
		mqtt_address = config.MQTT.Address
		mqtt_topic = config.MQTT.Topic
		interval = int(config.Interval)
		mqtt_cred_file = config.MQTT.CredFile
		debug = config.Debug
	} else {
		log.Printf("Error reading config. Will use default values usb_port:%s mqtt_address:%s mqtt_topic:%s interval:%d mqtt_cred_file:%s",
			usb_port, mqtt_address, mqtt_topic, interval, mqtt_cred_file)
	}

	var wg sync.WaitGroup

	wg.Add(2)

	telegram_channel := make(chan datamodel.P1Telegram, CHANNEL_SIZE)

	go p1port.ReadSerial(usb_port, telegram_channel, wg, debug)

	go mqttadapter.Start(mqtt_address, "smartmeter-reader", mqtt_topic, mqtt_cred_file, telegram_channel, wg, interval)

	wg.Wait()

}

func readConfig(path string) (*conf, error) {

	log.Printf("Reading Config from %s", path)

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Unable to read %s err   #%v \n", path, err)
		return nil, err
	}

	c := &conf{}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Printf("Unmarshal: %v\n", err)
		return nil, err
	}

	return c, nil
}

type conf struct {
	Interval int64  `yaml:"send_interval"`
	MQTT     mqtt   `yaml:"mqtt"`
	USBPort  string `yaml:"usb"`
	Debug    bool   `yaml:"debug,omitempty"`
}

type mqtt struct {
	Address  string `yaml:"address"`
	Topic    string `yaml:"topic"`
	CredFile string `yaml:"cred_file,omitempty"`
}
