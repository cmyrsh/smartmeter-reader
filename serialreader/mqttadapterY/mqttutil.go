package mqttadapterY

import (
	//"encoding/json"

	"fmt"
	"log"
	"os"
	"sync"

	"encoding/json"
	"time"

	"github.com/magiconair/properties"

	"../datamodel"

	//import the Paho Go MQTT library
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	osname, _ = os.Hostname()
)

const LOG_INTERVAL = 3600 * 1000

func buildClient(mqtt_url string, clientid string, username string, password string, topic string) MQTT.Client {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(mqtt_url)
	opts.SetClientID(fmt.Sprintf("%s@%s", clientid, osname))
	if len(username) > 0 {
		opts.SetUsername(username)
		if len(password) > 0 {
			opts.SetPassword(password)
		} else {
			log.Fatal("Password not provided for username")
		}
	}

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

func getcredentials(cred_file string) (user, pass string) {

	if len(cred_file) > 0 {
		p := properties.MustLoadFile(cred_file, properties.UTF8)

		return p.MustGetString("user"), p.MustGetString("password")
	}

	return "", ""
}

func MQTTy(mqtt_url string, clientid string, topic string, cred_file string, data_channel <-chan datamodel.P1Telegram, wg sync.WaitGroup, secs int) {

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

	usr, psswd := getcredentials(cred_file)

	c1 := buildClient(mqtt_url, clientid, usr, psswd, topic)

	defer c1.Disconnect(uint(secs))

	log.Println(fmt.Sprintf("Connected to MQTT Server %s with ClientId %s", mqtt_url, clientid))

	ticker := time.NewTicker(time.Duration(secs) * time.Second)

	log_threshhold := time.Now()

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
				} else if t.UnixMilli()-LOG_INTERVAL > log_threshhold.UnixMilli() {
					log.Printf("published messages: %d size: %d", arr_len, len(telegram_json))
					log_threshhold = t
				}

			} else {
				log.Println("array empty")
			}

		default:
			log.Println("no msg")
		}

	}

}
