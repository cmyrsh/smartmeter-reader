module com.github/cmyrsh/sm-reader/mqttadapter

replace com.github/cmyrsh/sm-reader/datamodel => ../datamodel

go 1.19

require (
	com.github/cmyrsh/sm-reader/datamodel v0.0.0-00010101000000-000000000000
	github.com/eclipse/paho.mqtt.golang v1.4.2
	github.com/magiconair/properties v1.8.6
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
)
