package main

import (
	"bufio"
	"fmt"
	"github.com/itchyny/volume-go"
	"github.com/tarm/serial"
	"log"
	"strconv"
)

func main() {
	c := &serial.Config{Name: "COM3", Baud: 9600}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		output := scanner.Text()
		fmt.Println(output)

		numericVolume, err := strconv.Atoi(output)
		if err != nil {
			panic(err)
			return
		}
		err = volume.SetVolume(numericVolume)
		if err != nil {
			panic(err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	//for {
	//	reader := bufio.NewReader(s)
	//	reply, err := reader.ReadBytes('\x0a')
	//	if err != nil {
	//		panic(err)
	//	}
	//	n, err := s.Read(reply)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Printf("%q", buf[:n])
	//}

}
