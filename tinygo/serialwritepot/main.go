package main

import (
	"machine"
	"strconv"
	"time"
)

var lastVolume = 999

func main() {
	potenciometro := machine.ADC0
	potReader := machine.ADC{potenciometro}
	machine.InitADC()

	for {
		outputStr := ""
		analogicValue := potReader.Get()
		//65535
		volume := int(valueMapping(analogicValue, 0, 65034, 0, 100))

		outputStr += strconv.Itoa(volume)

		if lastVolume != volume {
			println(outputStr)

			lastVolume = volume
		}
		time.Sleep(100 * time.Microsecond)
	}
}

func valueMapping(value, min, max, newMin, newMax uint16) float32 {
	scale := float32(value-min) / float32(max-min)
	return float32(newMin) + scale*float32(newMax-newMin)
}
