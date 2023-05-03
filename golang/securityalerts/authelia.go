package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

func NewAutheliaLineProcessor(mqttClient mqtt.Client, topic string) LineProcessor {
	return &AutheliaLogsMQTTDispatcher{
		mqttClient: mqttClient,
		topic:      topic,
	}
}

type AutheliaLogInput struct {
	Level    string    `json:"level"`
	Method   string    `json:"method"`
	Msg      string    `json:"msg"`
	Path     string    `json:"path"`
	RemoteIP string    `json:"remote_ip"`
	Time     time.Time `json:"time"`
}

type AutheliaLogsMQTTDispatcher struct {
	mqttClient mqtt.Client
	topic      string
}

func (a *AutheliaLogsMQTTDispatcher) Name() string      { return "authelia" }
func (a *AutheliaLogsMQTTDispatcher) LogSuffix() string { return "log" }
func (a *AutheliaLogsMQTTDispatcher) Process(line int, content string) error {
	log := AutheliaLogInput{}
	err := json.Unmarshal([]byte(content), &log)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("**[%s]** | **%s** | `%s` | **%s** | *%s* \n> %s", log.Level, log.Method, log.Path, log.RemoteIP, log.Time.Format(time.RFC3339), log.Msg)
	token := a.mqttClient.Publish(a.topic, 0, false, message)
	if token.Error() != nil {
		return fmt.Errorf("failed to publish to mqtt: %w", token.Error())
	}
	return nil
}
