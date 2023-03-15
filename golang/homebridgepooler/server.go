package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	homebridge "github.com/pedrohff/autom8/golang/pkg/homebridgeclient"
	"net/http"
	"os"
	"time"
)

type Server struct {
	env               Envs
	httpClient        *http.Client
	homeBridgeClient  *homebridge.Client
	temperatureClient *temperatureClient
	mqttClient        mqtt.Client
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
	s.env = Envs{
		Debug:     os.Getenv("DEBUG") != "false",
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

	s.httpClient = http.DefaultClient

	s.homeBridgeClient = homebridge.NewClient(s.httpClient, &homebridge.ClientOpts{
		Host:     s.env.HomeBridgeHost,
		User:     s.env.HomeBridgeUser,
		Password: s.env.HomeBridgePassword,
		Debug:    s.env.Debug,
	})

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
		fmt.Println("on connection lost ", err)
	}
	opts.OnReconnecting = func(client mqtt.Client, options *mqtt.ClientOptions) {
		fmt.Println("on reconnecting")
	}
	opts.OnConnect = func(client mqtt.Client) {
		fmt.Println("onconnect")
	}
	s.mqttClient = mqtt.NewClient(opts)
}
