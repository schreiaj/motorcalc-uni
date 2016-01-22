package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

type Motor struct {
	Name        string
	FreeSpeed   float64
	FreeAmps    float64
	StallAmps   float64
	StallTorque float64
	SpecVoltage float64
}

func main() {

	// motorFile, err := os.Open("motors.csv")
	// defer motorFile.Close()
	motorsData, _ := Asset("motors.csv")
	reader := csv.NewReader(strings.NewReader(string(motorsData)))
	reader.TrimLeadingSpace = true

	rawMotorData, err := reader.ReadAll()

	var motors []Motor
	var motor Motor
	for _, row := range rawMotorData {
		motor.Name = row[0]
		motor.FreeSpeed, _ = strconv.ParseFloat(row[1], 64)
		motor.FreeAmps, err = strconv.ParseFloat(row[3], 64)
		motor.StallAmps, err = strconv.ParseFloat(row[4], 64)
		motor.StallTorque, err = strconv.ParseFloat(row[2], 64)
		motor.SpecVoltage, err = strconv.ParseFloat(row[5], 64)
		motors = append(motors, motor)
	}

	selectedMotorString := os.Args[1]
	selectedMotor := new(Motor)
	for _, m := range motors {
		if m.Name == selectedMotorString {
			selectedMotor = &m
			break
		}
	}

	// ohms := selectedMotor.SpecVoltage / selectedMotor.StallAmps
	var speed float64
	var amps float64
	var powerOut float64
	var powerIn float64
	var efficiency float64
	var torque float64
	volts := 12.0

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	switch os.Args[2] {
	case "amps":
		amps, _ = strconv.ParseFloat(os.Args[3], 64)
		torque = (amps - selectedMotor.FreeAmps) / (selectedMotor.StallAmps - selectedMotor.FreeAmps) * selectedMotor.StallTorque
		speed = selectedMotor.FreeSpeed * (1.0 - torque/selectedMotor.StallTorque)
		powerOut = (speed / 60 * 2 * math.Pi) * torque
		powerIn = volts * amps
		efficiency = 100 * (powerOut / powerIn)
	case "torque":
		torque, _ = strconv.ParseFloat(os.Args[3], 64)
		amps = selectedMotor.FreeAmps + (torque/selectedMotor.StallTorque)*(selectedMotor.StallAmps-selectedMotor.FreeAmps)
		speed = selectedMotor.FreeSpeed * (1.0 - torque/selectedMotor.StallTorque)
		powerOut = (speed / 60 * 2 * math.Pi) * torque
		powerIn = volts * amps
		efficiency = 100 * (powerOut / powerIn)

	}

	str := fmt.Sprintf("%.2f\t%.0f\t%.2f\t%.2f\t%.2f\t%.2f", torque, speed, amps, powerOut, powerIn-powerOut, efficiency)
	fmt.Fprintln(w, "torque(N*m)\trpm\tamps\toutput(W)\theat(W)\teff")
	fmt.Fprintln(w, str)
	w.Flush()
	if err != nil {
		panic(err)
	}
}
