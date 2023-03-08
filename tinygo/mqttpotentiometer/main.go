package main

import (
	//"tinygo.org/x/drivers/net/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	//"machine"
	"strconv"
	"time"
)

var lastVolume = 999

func main() {
	potenciometro := machine.ADC0
	potReader := machine.ADC{potenciometro}
	machine.InitADC()

	opts := mqtt.NewClientOptions()
	_ = opts
	//opts.AddBroker("192.168.15.9:1883")
	//opts.SetClientID("tinygo")

	//mqttClient := mqtt.NewClient(opts)
	//connectToken := mqttClient.Connect()
	//_ = connectToken
	//connectToken.Wait()
	//if err := connectToken.Error(); err != nil {
	//	println(err)
	//	return
	//}
	for {
		outputStr := ""
		analogicValue := potReader.Get()
		//65535
		volume := int(valueMapping(analogicValue, 0, 65034, 0, 100))

		outputStr += strconv.Itoa(volume)

		if lastVolume != volume {
			println(outputStr)
			//if !mqttClient.IsConnected() {
			//	print("broker is not connected")
			//	continue
			//}
			//publishToken := mqttClient.Publish("/controller/volume/windows", 0, false, volume)
			//if err := publishToken.Error(); err != nil {
			//	println(err)
			//}

			lastVolume = volume
		}
		time.Sleep(100 * time.Microsecond)
	}
}

func valueMapping(value, min, max, newMin, newMax uint16) float32 {
	scale := float32(value-min) / float32(max-min)
	return float32(newMin) + scale*float32(newMax-newMin)
}
