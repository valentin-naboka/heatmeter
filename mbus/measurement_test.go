package mbus

import (
	"bytes"
	"encoding/xml"
	"log"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/net/html/charset"
)

var inputXML string = `
<MBusData>

    <SlaveInformation>
        <Id>58630782</Id>
        <Manufacturer>DME</Manufacturer>
        <Version>65</Version>
        <ProductName></ProductName>
        <Medium>Heat: Inlet</Medium>
        <AccessNumber>214</AccessNumber>
        <Status>00</Status>
        <Signature>0000</Signature>
    </SlaveInformation>

    <DataRecord id="0">
        <Function>Instantaneous value</Function>
        <StorageNumber>0</StorageNumber>
        <Unit>Reserved</Unit>
        <Quantity>Reserved</Quantity>
        <Value>12074.000000</Value>
    </DataRecord>

    <DataRecord id="1">
        <Function>Instantaneous value</Function>
        <StorageNumber>0</StorageNumber>
        <Tariff>1</Tariff>
        <Device>0</Device>
        <Unit>Reserved</Unit>
        <Quantity>Reserved</Quantity>
        <Value>0.000000</Value>
    </DataRecord>

    <DataRecord id="2">
        <Function>Instantaneous value</Function>
        <StorageNumber>0</StorageNumber>
        <Tariff>2</Tariff>
        <Device>0</Device>
        <Unit>Reserved</Unit>
        <Quantity>Reserved</Quantity>
        <Value>0.000000</Value>
    </DataRecord>

    <DataRecord id="3">
        <Function>Instantaneous value</Function>
        <StorageNumber>0</StorageNumber>
        <Unit>m^3</Unit>
        <Quantity>Volume</Quantity>
        <Value>1442.384000</Value>
    </DataRecord>

    <DataRecord id="4">
        <Function>Instantaneous value</Function>
        <StorageNumber>0</StorageNumber>
        <Unit>W</Unit>
        <Quantity>Power</Quantity>
        <Value>443.000000</Value>
    </DataRecord>

    <DataRecord id="5">
        <Function>Instantaneous value</Function>
        <StorageNumber>0</StorageNumber>
        <Unit>m^3/h</Unit>
        <Quantity>Volume flow</Quantity>
        <Value>0.018000</Value>
    </DataRecord>

    <DataRecord id="6">
        <Function>Instantaneous value</Function>
        <StorageNumber>0</StorageNumber>
        <Unit>°C</Unit>
        <Quantity>Flow temperature</Quantity>
        <Value>76.900000</Value>
    </DataRecord>

    <DataRecord id="7">
        <Function>Instantaneous value</Function>
        <StorageNumber>0</StorageNumber>
        <Unit>°C</Unit>
        <Quantity>Return temperature</Quantity>
        <Value>55.200000</Value>
    </DataRecord>

    <DataRecord id="8">
        <Function>Instantaneous value</Function>
        <StorageNumber>0</StorageNumber>
        <Unit>s</Unit>
        <Quantity>Operating time</Quantity>
        <Value>104716800.000000</Value>
    </DataRecord>

    <DataRecord id="9">
        <Function>Instantaneous value</Function>
        <StorageNumber>0</StorageNumber>
        <Unit>s</Unit>
        <Quantity>Operating time</Quantity>
        <Value>0.000000</Value>
    </DataRecord>

</MBusData>`

func printCaller(t *testing.T, depth int) {
	function, file, line, _ := runtime.Caller(depth)
	trimName := func(n string) string {
		i := strings.LastIndex(n, "/")
		if i == -1 {
			return n
		}
		return n[i+1:]
	}
	t.Logf("%s: line %d, function: %s\n", trimName(file), line, trimName(runtime.FuncForPC(function).Name()))
}

func expectEqual(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		printCaller(t, 2)
		t.Errorf("expected: %v, actual: %v\n", expected, actual)
	}
}

func expectIntEqual(t *testing.T, expected, actual int) {
	expectEqual(t, &expected, &actual)
}

func TestUnmarshalMeasurement(t *testing.T) {
	reader := bytes.NewReader([]byte(inputXML))

	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	var m Measurement
	err := decoder.Decode(&m)

	if err != nil {
		log.Fatal(err)
	}

	expectEqual(t, Calories(12074), m.Energy)
	expectEqual(t, CubicMetre(1442.384), m.Volume)
	expectEqual(t, Watt(443), m.Power)
	expectEqual(t, CubicMetresPerHour(0.018), m.VolumeFlow)
	expectEqual(t, Celsius(77), m.FlowTemp)
	expectEqual(t, Celsius(55), m.ReturnTemp)
	expectEqual(t, Seconds(104716800), m.OperatingTime)
	expectEqual(t, Seconds(0), m.ErrorTime)
}
