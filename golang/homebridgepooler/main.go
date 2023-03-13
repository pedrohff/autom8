package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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

func main() {

	envs := Envs{
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

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	mqttClient := newMQTTClient(envs)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	httpClient := http.DefaultClient
	go func() {
		for {
			ctx := context.Background()

			temperature := getTemp(ctx, httpClient, envs)
			_ = temperature
			message := fmt.Sprintf("{\"temperature\": %.2f}", temperature)
			if token := mqttClient.Publish(envs.MQTTTopic, 2, true, message); token.Wait() && token.Error() != nil {
				fmt.Println("produce err: ", token.Error().Error())
			}
			time.Sleep(envs.PoolInterval)
		}
	}()
	<-stopChan

}

func newMQTTClient(envs Envs) mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker(envs.MQTTHost).SetClientID("homebridgepooler")
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

	return mqtt.NewClient(opts)
}

func getToken(ctx context.Context, httpClient *http.Client, env Envs) (*TokenResponse, error) {
	reqBody := TokenRequest{env.HomeBridgeUser, env.HomeBridgePassword}
	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	url := env.HomeBridgeHost + "auth/login"
	fmt.Println(url)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(reqBodyJSON))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if env.Debug {
		fmt.Println(string(responseBody))
	}
	switch response.StatusCode {
	case 201:
		r := TokenResponse{}
		err := json.Unmarshal(responseBody, &r)
		if err != nil {
			return nil, err
		}
		return &r, nil

	default:
		return nil, fmt.Errorf("invalid status code %d", response.StatusCode)
	}

}

func getTemp(ctx context.Context, client *http.Client, env Envs) float64 {
	path := fmt.Sprintf("accessories/%s", env.TemperatureAccessoryID)
	request, err := http.NewRequest(http.MethodGet, env.HomeBridgeHost+path, nil)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	token, err := getToken(ctx, client, env)
	if err != nil {
		return 0
	}
	request.Header.Add("Authorization", "Bearer "+token.AccessToken)
	request.Header.Add("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer response.Body.Close()

	res := HBTemp{}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	if env.Debug {
		fmt.Println(string(responseBody))
	}
	switch response.StatusCode {
	case 200:
		err = json.Unmarshal(responseBody, &res)
		if err != nil {
			fmt.Println(err)
			return 0
		}
		return res.Values.CurrentTemperature
	default:
		fmt.Println("invalid status: ", response.StatusCode)
		return 0
	}

}

type HBTemp struct {
	Values HBTempValues `json:"values"`
}

type HBTempValues struct {
	CurrentTemperature float64 `json:"CurrentTemperature"`
}
