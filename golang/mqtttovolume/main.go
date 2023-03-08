package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/itchyny/volume-go"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func readMessage(client mqtt.Client, message mqtt.Message) {
	output := string(message.Payload())
	println(output)
	numericVolume, err := strconv.Atoi(output)
	if err != nil {
		panic(err)
		return
	}
	err = volume.SetVolume(numericVolume)
	if err != nil {
		panic(err)
		return
	}
}

func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	opts := mqtt.NewClientOptions().AddBroker("192.168.15.9:1883").SetClientID("gotrivial")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.AutoReconnect = true
	topic := "mytest"

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Println("on connection lost ", err)
	}
	opts.OnReconnecting = func(client mqtt.Client, options *mqtt.ClientOptions) {
		fmt.Println("on reconnecting")
	}
	opts.OnConnect = func(client mqtt.Client) {
		fmt.Println("onconnect")
		if token := client.Subscribe(topic, 0, readMessage); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	<-stopChan
}
