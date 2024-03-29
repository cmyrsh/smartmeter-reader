module com.github/cmyrsh/sm-reader/run

replace com.github/cmyrsh/sm-reader/datamodel => ../datamodel

replace com.github/cmyrsh/sm-reader/p1port => ../p1port

replace com.github/cmyrsh/sm-reader/mqttadapter => ../mqttadapter

go 1.19

require (
	com.github/cmyrsh/sm-reader/datamodel v0.0.0-00010101000000-000000000000
	com.github/cmyrsh/sm-reader/mqttadapter v0.0.0-00010101000000-000000000000
	com.github/cmyrsh/sm-reader/p1port v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/eclipse/paho.mqtt.golang v1.4.2 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07 // indirect
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
)
