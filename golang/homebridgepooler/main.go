package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	server := Server{}
	err := server.Setup()
	if err != nil {
		panic(err)
		return
	}

	go func() {
		for {
			ctx := context.Background()

			temperature, err := server.temperatureClient.GetRoomTemperature(ctx)
			if err != nil {
				fmt.Println(err)
				time.Sleep(time.Minute)
				continue
			}
			_ = temperature
			message := fmt.Sprintf("{\"temperature\": %.2f}", temperature)
			if token := server.mqttClient.Publish(server.env.MQTTTopic, 2, true, message); token.Wait() && token.Error() != nil {
				fmt.Println("produce err: ", token.Error().Error())
			}
			time.Sleep(server.env.PoolInterval)
		}
	}()
	<-stopChan

}
