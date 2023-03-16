package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"net/http"
	"os"
	"pkg/homebridge"
	"pkg/log"
	"time"
)

type Server struct {
	env               Envs
	httpClient        *http.Client
	homeBridgeClient  *homebridge.Client
	temperatureClient *temperatureClient
	mqttClient        mqtt.Client
	logger            *zap.Logger
}

type Envs struct {
	HomeBridgeHost         string
	HomeBridgeUser         string
	HomeBridgePassword     string
	TemperatureAccessoryID string
	PoolInterval           time.Duration
	MQTTHost               string
	MQTTTopic              string
	Debug                  bool
}

func (s *Server) Setup() error {
	appName := "homebridgepooler"
	s.env = Envs{
		Debug:     os.Getenv("DEBUG") == "true",
		MQTTHost:  os.Getenv("MQTT_HOST"),
		MQTTTopic: os.Getenv("MQTT_TOPIC"),

		HomeBridgeHost:         os.Getenv("HOMEBRIDGE_HOST"),
		HomeBridgeUser:         os.Getenv("HOMEBRIDGE_AUTH_USER"),
		HomeBridgePassword:     os.Getenv("HOMEBRIDGE_AUTH_PASSWORD"),
		TemperatureAccessoryID: os.Getenv("HOMEBRIDGE_TEMPERATURE_ACESSORY_ID"),
		PoolInterval: func() time.Duration {
			osInterval := os.Getenv("POOL_INTERVAL")
			defaultDuration := 5000 * time.Millisecond
			if osInterval == "" {
				return defaultDuration

			}
			duration, err := time.ParseDuration(osInterval)
			if err != nil {
				return defaultDuration
			}
			return duration
		}(),
	}

	s.logger = log.NewLogger(appName, s.env.Debug)

	s.httpClient = http.DefaultClient

	s.homeBridgeClient = homebridge.NewClient(s.httpClient, &homebridge.ClientOpts{
		Host:     s.env.HomeBridgeHost,
		User:     s.env.HomeBridgeUser,
		Password: s.env.HomeBridgePassword,
	}, s.logger)

	s.temperatureClient = NewTemperatureClient(s.homeBridgeClient, s.env.TemperatureAccessoryID)
	s.starMQTTClient()
	if token := s.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (s *Server) starMQTTClient() {
	opts := mqtt.NewClientOptions().AddBroker(s.env.MQTTHost).SetClientID("homebridgepooler")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.AutoReconnect = true

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		s.logger.Warn("on connection lost", zap.Error(err))
	}
	opts.OnReconnecting = func(client mqtt.Client, options *mqtt.ClientOptions) {
		s.logger.Warn("on reconnecting")
	}
	opts.OnConnect = func(client mqtt.Client) {
		s.logger.Info("onconnect")
	}
	s.mqttClient = mqtt.NewClient(opts)
}
