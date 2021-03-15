package main

import "C"
import (
	"fmt"
	"heatmeter/mbus"
	"log"
)

func main() {
	var reader mbus.Reader
	err := reader.Open("/dev/cu.usbserial-1410", 13, 2400)
	defer reader.Close()
	if err != nil {
		log.Fatal(err)
	}

	measurement, _ := reader.ReadData()
	fmt.Printf("%v\n", measurement)
}
