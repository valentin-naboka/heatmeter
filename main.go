package main

import "C"
import (
	"fmt"
	"heatmeter/mbus"
	"log"
	"os"
	"strconv"
)

func main() {
	var device, mbusIDVar, baudrateVar string
	if device = os.Getenv("HM_DEVICE"); device == "" {
		log.Fatal("HM_DEVICE variable is not set")
	}

	if mbusIDVar = os.Getenv("HM_MBUS_ID"); mbusIDVar == "" {
		log.Fatal("HM_MBUS_ID variable is not set")
	}

	mbusID, err := strconv.Atoi(mbusIDVar)
	if err != nil{
		log.Fatalf("Wrong mbus ID: %s, %s", mbusIDVar, err)
	}

	if baudrateVar=os.Getenv("HM_BAUDRATE"); baudrateVar == "" {
		log.Fatal("HM_BAUDRATE variable is not set")
	}

	baudrate, err := strconv.Atoi(baudrateVar)
	if err != nil{
		log.Fatalf("Wrong baudrate: %s, %s", baudrateVar, err)
	}

	var reader mbus.Reader
	err = reader.Open(device, uint8(mbusID), uint16(baudrate))
	defer reader.Close()
	if err != nil {
		log.Fatal(err)
	}

	measurement, _ := reader.ReadData()
	fmt.Printf("%v\n", *measurement)
}
