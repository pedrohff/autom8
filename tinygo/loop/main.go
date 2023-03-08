package main

import (
	"machine"
	"strconv"
	"time"

	tservo "tinygo.org/x/drivers/servo"
)

var isPressed = false

var lastPot = 0

var intervalMS = 600

func main() {
	// machine.I2C0.Configure(machine.I2CConfig{Frequency: machine.TWI_FREQ_100KHZ})
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	button := machine.PD2
	button.Configure(machine.PinConfig{Mode: machine.PinInput})

	potenciometro := machine.ADC0
	potReader := machine.ADC{potenciometro}
	machine.InitADC()

	// servo := Init(machine.PD6)
	// go servo.ServoRoutine()

	newServo, err := tservo.New(machine.Timer1, machine.D9)
	if err != nil {
		for {
			println("could not configure servo")
			println(err.Error())
			time.Sleep(time.Second)
		}
		return
	}

	for {
		// loopState := button.Get()
		// if isPressed != loopState {
		// 	if isPressed {
		// 		println("pressed")
		// 		led.Set(!led.Get())
		// 	}
		// 	isPressed = loopState
		// }
		intervalDuration := time.Millisecond * time.Duration(intervalMS)

		analogicValue := potReader.Get()
		currPercentage := getPercentageFromAnalogicInput(analogicValue)
		if false {
			// if lastPot != int(analogicValue) {

			println("currPercentage " + strconv.Itoa(int(currPercentage)) + "   (" + strconv.Itoa(int(analogicValue)) + ")")
			// fmt.Printf()
			// fmt.Printf("currPercentage %.2g (%d)\n", currPercentage, analogicValue)
			lastPot = int(analogicValue)
		}
		// angle := valueMapping(analogicValue, 0, 65535, 500, 2500)
		// // println("angle: " + strconv.Itoa(int(angle)))
		// newServo.SetMicroseconds(int16(angle))
		if intervalMS < 900 {

			newServo.SetMicroseconds(1100)
			time.Sleep(intervalDuration)
			newServo.SetMicroseconds(1800)
			time.Sleep(intervalDuration)
		}
		maxInterval := 1000
		minInterval := 10
		intervalMS = int(valueMapping(analogicValue, 0, 65535, uint16(minInterval), uint16(maxInterval)))

	}
}

func getPercentageFromAnalogicInput(v uint16) float64 {
	return float64(100.0/65535) * float64(v)
}

type Device struct {
	pin      machine.Pin
	on       bool
	angle    uint8
	pulseMin uint16
	pulseMax uint16
}

func Init(pin machine.Pin) Device {
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	return Device{pin: pin, on: false, pulseMin: 600, pulseMax: 2400}
}

func (d *Device) ServoRoutine() {
	for {
		if d.on {
			// println("set angle ---  " + strconv.Itoa(int(d.angle)))
			pulse := valueMapping(uint16(d.angle), 0, 180, d.pulseMin, d.pulseMax)
			d.pin.High()
			time.Sleep(time.Microsecond * time.Duration(pulse))
			d.pin.Low()
			time.Sleep(time.Microsecond * time.Duration(20000.0-pulse))
		}
	}
}

func (d *Device) Angle(angle uint8) {
	if angle < 0 {
		angle = 0
	} else if angle > 180 {
		angle = 180
	}
	d.on = true
	d.angle = angle
}

func (d *Device) PulseRange(min, max uint16) {
	d.pulseMin = min
	d.pulseMax = max
}

func (d *Device) Deinit() {
	d.on = false
}

func valueMapping(value, min, max, newMin, newMax uint16) float32 {
	scale := float32(value-min) / float32(max-min)
	return float32(newMin) + scale*float32(newMax-newMin)
}
