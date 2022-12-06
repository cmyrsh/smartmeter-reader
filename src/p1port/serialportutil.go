package p1port

import (
	"log"

	"bufio"
	"encoding/json"
	"strings"
	"sync"

	"com.github/cmyrsh/sm-reader/datamodel"
	"github.com/tarm/serial"
)

const (
	BAUD_RATE int = 115200
)

func ReadSerial(usb_port string, channel_telegram chan<- datamodel.P1Telegram, wg sync.WaitGroup) {
	defer wg.Done()

	config := &serial.Config{}

	config.Parity = serial.ParityNone
	config.StopBits = serial.Stop1
	config.Size = 8
	config.Baud = BAUD_RATE
	config.Name = usb_port

	p, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	reader := bufio.NewScanner(p)

	log.Println("Opened Serial Port. Starting to read from Serial Port")

	if p == nil {
		log.Printf("%#v", p)
		log.Panic("serial port is nil")
	}
	var telegram *datamodel.P1Telegram
	for reader.Scan() {
		linex := reader.Text()

		if strings.HasPrefix(linex, "!") && nil != telegram {

			select {
			case channel_telegram <- *telegram:
				//log.Printf("sending telegram to channel %#v", telegram)
			default:
				log.Printf("buffer full. ignoring telegram %#v", telegram)
			}

		} else if strings.HasPrefix(linex, "/") {
			//log.Println("creating new telegram " + linex)
			telegram = new(datamodel.P1Telegram)
		} else if nil != telegram {

			telegram.PopulateFromLine(linex)
		}
	}

}

func ReadChannel(data_channel <-chan datamodel.P1Telegram, wg sync.WaitGroup) {
	for {
		telegram := <-data_channel
		telegram_json, err := json.Marshal(telegram)

		if err != nil {
			log.Fatal(err)
		}
		log.Println("got -- ", telegram_json)
	}
}
