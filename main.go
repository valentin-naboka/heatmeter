package main

import "C"
import (
	"fmt"
	"heatmeter/generated"
	"heatmeter/logger"
	"heatmeter/mbus"
	"heatmeter/report"
	"heatmeter/safe"
	"log"
	"os"
	"strconv"
	"time"
)

func isDayOfReport() bool {
	var startMonthVar, endMonthVar, dayOfReportVar, endMonthDayOfReportVar string
	if startMonthVar = os.Getenv("HM_START_MONTH"); startMonthVar == "" {
		log.Fatal("HM_START_MONTH variable is not set")
	}
	startMonth, err := strconv.Atoi(startMonthVar)
	if err != nil {
		log.Fatalf("Wrong start month: %s, %s", startMonthVar, err)
	}

	if endMonthVar = os.Getenv("HM_END_MONTH"); endMonthVar == "" {
		log.Fatal("HM_END_MONTH variable is not set")
	}
	endMonth, err := strconv.Atoi(endMonthVar)
	if err != nil {
		log.Fatalf("Wrong end month: %s, %s", startMonthVar, err)
	}

	if dayOfReportVar = os.Getenv("HM_DAY_OF_REPORT"); dayOfReportVar == "" {
		log.Fatal("HM_DAY_OF_REPORT variable is not set")
	}
	dayOfReport, err := strconv.Atoi(dayOfReportVar)
	if err != nil {
		log.Fatalf("Wrong day of report: %s, %s", dayOfReportVar, err)
	}

	if endMonthDayOfReportVar = os.Getenv("HM_END_MONTH_DAY_OF_REPORT"); endMonthDayOfReportVar == "" {
		log.Fatal("HM_END_MONTH_DAY_OF_REPORT variable is not set")
	}
	endMonthDayOfReport, err := strconv.Atoi(endMonthDayOfReportVar)
	if err != nil {
		log.Fatalf("Wrong end month day of report: %s, %s", endMonthDayOfReportVar, err)
	}

	currentMonth := time.Now().Month()
	currentDay := time.Now().Day()
	return (currentMonth >= time.Month(endMonth) || currentMonth < time.Month(startMonth)) && currentDay == dayOfReport ||
		currentMonth == time.Month(endMonth) && currentDay == endMonthDayOfReport
}

func main() {

	if !isDayOfReport() {
		return
	}

	data6, err := safe.Decrypt(generated.Data6)
	if err != nil {
		log.Fatal("Unable to decrypt data6: ", err)
	}

	data7, err := safe.Decrypt(generated.Data7)
	if err != nil {
		log.Fatal("Unable to decrypt data7: ", err)
	}

	data7Int, err := strconv.Atoi(data7)
	if err != nil {
		log.Fatal("data7 is not int: ", err)
	}

	logger, err := logger.NewTelegram(data6, int64(data7Int))
	if err != nil {
		log.Fatal("unable create bot: ", err)
	}
	log.SetOutput(logger)

	var device, mbusIDVar, baudrateVar string
	if device = os.Getenv("HM_DEVICE"); device == "" {
		log.Fatal("HM_DEVICE variable is not set")
	}

	if mbusIDVar = os.Getenv("HM_MBUS_ID"); mbusIDVar == "" {
		log.Fatal("HM_MBUS_ID variable is not set")
	}

	mbusID, err := strconv.Atoi(mbusIDVar)
	if err != nil {
		log.Fatalf("Wrong mbus ID: %s, %s", mbusIDVar, err)
	}

	if baudrateVar = os.Getenv("HM_BAUDRATE"); baudrateVar == "" {
		log.Fatal("HM_BAUDRATE variable is not set")
	}

	baudrate, err := strconv.Atoi(baudrateVar)
	if err != nil {
		log.Fatalf("Wrong baudrate: %s, %s", baudrateVar, err)
	}

	htmlDumpPathVar := os.Getenv("HM_HTML_DUMP_PATH")
	if htmlDumpPathVar == "" {
		htmlDumpPathVar = "/var/log/heatmeter/"
	}

	data3, err := safe.Decrypt(generated.Data3)
	if err != nil {
		log.Fatal("Unable to decrypt data3: ", err)
	}

	data4, err := safe.Decrypt(generated.Data4)
	if err != nil {
		log.Fatal("Unable to decrypt data4: ", err)
	}

	data5, err := safe.Decrypt(generated.Data5)
	if err != nil {
		log.Fatal("Unable to decrypt data5: ", err)
	}

	var reader mbus.Reader
	err = reader.Open(device, uint8(mbusID), uint16(baudrate))
	defer reader.Close()
	if err != nil {
		log.Fatal(err)
	}

	measurement, err := reader.ReadData()
	if err != nil {
		log.Fatal("Unable to get measurement: ", err)
	}

	errorHoursStr := strconv.Itoa(int(measurement.ErrorTime) / 3600)
	operatingDaysStr := strconv.Itoa(int(measurement.OperatingTime) / 3600 / 24)

	flowTempStr := strconv.Itoa(int(measurement.FlowTemp))
	returnTempStr := strconv.Itoa(int(measurement.ReturnTemp))
	powerStr := strconv.Itoa(int(measurement.Power))

	energyStr := fmt.Sprintf("%.3f", float32(measurement.Energy)/1000)
	volumeStr := fmt.Sprintf("%.3f", measurement.Volume)
	volumeFlowStr := fmt.Sprintf("%.3f", measurement.VolumeFlow/1000)

	submitter := report.NewSubmitter(data3, htmlDumpPathVar)
	ok := submitter.Execute(data4, data5,
		//energy,
		"",
		volumeStr,
		volumeFlowStr,
		powerStr,
		flowTempStr,
		returnTempStr,
		operatingDaysStr,
		errorHoursStr)

	if ok {
		log.Printf("Data has been succesfully submited.\n Energy: %s, volume: %s, volume flow: %s, power: %s, flow temperature: %s, return temperature: %s, operating days: %s, error hours: %s",
			energyStr,
			volumeStr,
			volumeFlowStr,
			powerStr,
			flowTempStr,
			returnTempStr,
			operatingDaysStr,
			errorHoursStr)
	}
}
