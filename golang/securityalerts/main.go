package main

import (
	"go.uber.org/zap"
	"os"
	"os/signal"
	"pkg/log"
	"syscall"
)

func main() {

	// 1 list files in directory after x timestamp
	logger := log.NewLogger("securityalerts", false)
	pointer, err := NewPointer(logger)
	if err != nil {
		logger.Fatal("failed to open pointer file", zap.Error(err))
		return
	}

	stopWorker := pointer.StartWorker()
	defer stopWorker()

	mqttHost := "mosquitto:1883"

	if envNewHost := os.Getenv("MQTT_HOST"); envNewHost != "" {
		mqttHost = envNewHost
	}

	mqttClient := startMQTTClient(mqttHost, logger)
	if token := mqttClient.Connect(); token.Error() != nil {
		logger.Fatal("failed to connect to mqtt broker", zap.Error(err))
		return
	}

	topic := "/home/servers/glados/notifications"
	if envNewTopic := os.Getenv("MQTT_TOPIC"); envNewTopic != "" {
		topic = envNewTopic
	}
	autheliaProcessor := NewAutheliaLineProcessor(mqttClient, topic)

	reader := NewFileReader(autheliaProcessor, pointer, logger)

	directory := "."
	if envNewDir := os.Getenv("TAIL_DIRECTORY"); envNewDir != "" {
		directory = envNewDir
	}

	go func() {
		err := reader.TailDir(directory)
		if err != nil {
			logger.Fatal("failed to read directory", zap.Error(err), zap.String("directory", directory))
			return
		}
	}()
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan
}
