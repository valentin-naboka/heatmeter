package mbus

import (
	"encoding/xml"
	"math"
	"strconv"
)

type Calories uint16
type CubicMetre float32
type Watt uint16
type CubicMetresPerHour float32
type Celsius uint64
type Seconds uint64

type Measurement struct {
	Energy        Calories
	Volume        CubicMetre
	Power         Watt
	VolumeFlow    CubicMetresPerHour
	FlowTemp      Celsius
	ReturnTemp    Celsius
	OperatingTime Seconds
	ErrorTime     Seconds
}

func (m *Measurement) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	result := struct {
		Records []struct {
			Value string `xml:"Value"`
		} `xml:"DataRecord"`
	}{}

	if err := d.DecodeElement(&result, &start); err != nil {
		return err
	}

	{
		value, err := strconv.ParseFloat(result.Records[0].Value, 32)
		if err != nil {
			return err
		}
		m.Energy = Calories(math.Round(value))
	}

	{
		value, err := strconv.ParseFloat(result.Records[3].Value, 32)
		if err != nil {
			return err
		}
		m.Volume = CubicMetre(value)
	}

	{
		value, err := strconv.ParseFloat(result.Records[4].Value, 32)
		if err != nil {
			return err
		}
		m.Power = Watt(value)
	}

	{
		value, err := strconv.ParseFloat(result.Records[5].Value, 32)
		if err != nil {
			return err
		}
		m.VolumeFlow = CubicMetresPerHour(value)
	}

	{
		value, err := strconv.ParseFloat(result.Records[6].Value, 32)
		if err != nil {
			return err
		}
		m.FlowTemp = Celsius(math.Round(value))
	}

	{
		value, err := strconv.ParseFloat(result.Records[7].Value, 32)
		if err != nil {
			return err
		}
		m.ReturnTemp = Celsius(math.Round(value))
	}

	{
		value, err := strconv.ParseFloat(result.Records[8].Value, 32)
		if err != nil {
			return err
		}
		m.OperatingTime = Seconds(value)
	}

	{
		value, err := strconv.ParseFloat(result.Records[9].Value, 32)
		if err != nil {
			return err
		}
		m.ErrorTime = Seconds(value)
	}

	return nil
}
