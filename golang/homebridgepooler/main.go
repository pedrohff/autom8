package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
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
			server.logger.Info("sending temperature", zap.Float64("temperature", temperature))
			if token := server.mqttClient.Publish(server.env.MQTTTopic, 2, true, message); token.Wait() && token.Error() != nil {
				server.logger.Error("produce err: ", zap.Error(token.Error()))
			}
			time.Sleep(server.env.PoolInterval)
		}
	}()
	<-stopChan

}
