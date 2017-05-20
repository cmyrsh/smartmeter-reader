package mqttadapter

import (
	//"encoding/json"
	"log"
	"os"
	"sync"

	"../datamodel"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/mqtt"
	"time"
	"encoding/json"
	//import the Paho Go MQTT library
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	osname, _ = os.Hostname()
)

func MQTTx(mqtt_url string, clientid string, topic string, data_channel <-chan datamodel.P1Telegram, wg sync.WaitGroup, secs int) {

	defer wg.Done()

	mqttAdaptor := mqtt.NewAdaptor(mqtt_url, clientid+"@"+osname)

	work := func() {

		gobot.Every(time.Duration(secs) * time.Second, func() {

			//log.Println("checking for messages...")
			var telegram_array []datamodel.P1Telegram

			select {
			case telegram := <-data_channel:

				telegram_array = append(telegram_array, telegram)
				//log.Printf("datagram array -- %d", len(telegram_array))

				for len(data_channel) > 0 {
					telegram_array = append(telegram_array, <-data_channel)
					//log.Printf("datagram array -- %d", len(telegram_array))
				}

				arr_len := len(telegram_array)

				if arr_len > 0 {
					telegram_json, err := json.Marshal(telegram_array)
					if err != nil {
						log.Fatal(err)
					}

					//published := mqttAdaptor.Publish(topic, telegram_json)

					log.Printf("published: %b messages: %d size: %d", true, arr_len, len(telegram_json))
				} else {
					log.Println("array empty")
				}


			default:
				log.Println("no msg")
			}


		})

	}
	robot := gobot.NewRobot(clientid,[]gobot.Connection{mqttAdaptor},work,)

	robot.Start()

}

type mqttD struct {
	c_o MQTT.ClientOptions
	c MQTT.Client

	/*
		// Name returns the label for the Adaptor
	Name() string
	// SetName sets the label for the Adaptor
	SetName(n string)
	// Connect initiates the Adaptor
	Connect() error
	// Finalize terminates the Adaptor
	Finalize() error
	 */
}

func (d mqttD) Name() (name string) {
	return d.c_o.ClientID
}

func (d mqttD) SetName(name string) {
	d.c_o.SetClientID(name)
}

func (d mqttD) Connect() error {
	if token := d.c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	} else {
		return nil
	}
}


func (d mqttD) Finalize() error {
	d.c.Disconnect(10)
	return nil
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
	opts := MQTT.NewClientOptions().AddBroker(mqtt_url)
	opts.SetClientID(clientid)
	opts.SetWill("equipments/meterreader", "GoingDown", 1, false)
	opts.SetCleanSession(true)


	//create and start a client using the above ClientOptions
	c1 := MQTT.NewClient(opts)


	mymqttadp :=mqttD{*opts,c1}

	work := func() {

		gobot.Every(time.Duration(secs) * time.Second, func() {

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
					token := mymqttadp.c.Publish(topic, 1, false, telegram_json)

					//log.Println("waiting for token")
					if(token.Wait() && token.Error() != nil) {
						log.Panic(token.Error())
					} else {
						log.Printf("published: %b messages: %d size: %d", true, arr_len, len(telegram_json))
					}


				} else {
					log.Println("array empty")
				}


			default:
				log.Println("no msg")
			}


		})

	}
	robot := gobot.NewRobot(clientid,[]gobot.Connection{mymqttadp},work,)

	robot.Start()

}
