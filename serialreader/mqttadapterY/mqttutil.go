package mqttadapterY

import (
	//"encoding/json"

	"fmt"
	"log"
	"os"
	"sync"

	"encoding/json"
	"time"

	"../datamodel"

	//import the Paho Go MQTT library
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	osname, _ = os.Hostname()
)

func buildClient(mqtt_url string, clientid string, username string, password string, topic string) MQTT.Client {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(mqtt_url)
	opts.SetClientID(fmt.Sprintf("%s@%s", clientid, osname))
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetWill(fmt.Sprintf("will:%s", topic), "GoingDown", 1, false)
	opts.SetCleanSession(true)

	opts.OnConnect = func(client MQTT.Client) {
		log.Println("Connected")
	}
	opts.OnConnectionLost = func(client MQTT.Client, err error) {
		log.Printf("Connect lost: %v", err)
	}
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())

	}
	return client
}

func MQTTy(mqtt_url string, clientid string, topic string, data_channel <-chan datamodel.P1Telegram, wg sync.WaitGroup, secs int) {

	defer wg.Done()
	/*
		var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
			log.Printf("TOPIC: %s\n", msg.Topic())
			log.Printf("MSG: %s\n", msg.Payload())
		}
		opts.SetDefaultPublishHandler(f)
	*/
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	//create and start a client using the above ClientOptions
	c1 := buildClient(mqtt_url, clientid, "", "", topic)

	defer c1.Disconnect(uint(secs))

	log.Println("Created MQTT Client")

	ticker := time.NewTicker(time.Duration(secs) * time.Second)

	for t := range ticker.C {

		//  Create slice for telegram
		var telegram_array []datamodel.P1Telegram

		select {
		case telegram := <-data_channel:
			// if message found in channel then push the message to slice
			telegram_array = append(telegram_array, telegram)
			// collect all the messages from channel
			for len(data_channel) > 0 {
				telegram_array = append(telegram_array, <-data_channel)
			}

			arr_len := len(telegram_array)

			if arr_len > 0 {
				telegram_json, err := json.Marshal(telegram_array)
				if err != nil {
					log.Fatal(err)
				}
				token := c1.Publish(topic, 1, false, telegram_json)

				//log.Println("waiting for token")
				if token.Wait() && token.Error() != nil {
					log.Panic(token.Error())
				} else {
					log.Printf("published messages: %d size: %d", arr_len, len(telegram_json))
				}

			} else {
				log.Println("array empty")
			}

		default:
			log.Println("no msg")
		}

		log.Printf("Time :%s", t)
	}

}
