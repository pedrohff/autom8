package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"time"
)

func startMQTTClient(mqttHost string, logger *zap.Logger) mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker(mqttHost).SetClientID("securityalerts")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.AutoReconnect = true

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		logger.Warn("on connection lost", zap.Error(err))
	}
	opts.OnReconnecting = func(client mqtt.Client, options *mqtt.ClientOptions) {
		logger.Warn("on reconnecting")
	}
	opts.OnConnect = func(client mqtt.Client) {
		logger.Info("onconnect")
	}
	return mqtt.NewClient(opts)
}
