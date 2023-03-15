package main

import (
	"context"
	"encoding/json"
	homebridge "github.com/pedrohff/autom8/golang/pkg/homebridgeclient"
)

type TemperatureAccessoryResult struct {
	Values TemperatureAccessoryValues `json:"values"`
}

type TemperatureAccessoryValues struct {
	CurrentTemperature float64 `json:"CurrentTemperature"`
}

type temperatureClient struct {
	*homebridge.Client
	accessoryId string
}

func NewTemperatureClient(client *homebridge.Client, accessoryId string) *temperatureClient {
	return &temperatureClient{
		Client:      client,
		accessoryId: accessoryId,
	}
}

func (t *temperatureClient) GetRoomTemperature(ctx context.Context) (float64, error) {
	accessoryJSON, err := t.GetAccessory(ctx, t.accessoryId)
	if err != nil {
		return 0, err
	}
	result := TemperatureAccessoryResult{}
	err = json.Unmarshal(accessoryJSON, &result)
	if err != nil {
		return 0, err
	}
	return result.Values.CurrentTemperature, nil
}
