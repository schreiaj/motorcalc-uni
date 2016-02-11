package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/sajari/fuzzy"
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
	reader.Comment = '#'

	rawMotorData, _ := reader.ReadAll()

	var motors []Motor
	var motorNames []string
	var motor Motor
	for _, row := range rawMotorData {
		motorNames = append(motorNames, row[0])
		motor.Name = row[0]
		motor.FreeSpeed, _ = strconv.ParseFloat(row[1], 64)
		motor.FreeAmps, _ = strconv.ParseFloat(row[3], 64)
		motor.StallAmps, _ = strconv.ParseFloat(row[4], 64)
		motor.StallTorque, _ = strconv.ParseFloat(row[2], 64)
		motor.SpecVoltage, _ = strconv.ParseFloat(row[5], 64)
		motors = append(motors, motor)
	}
	app := cli.NewApp()
	app.Name = "MotorCalc"
	app.EnableBashCompletion = true
	app.Usage = "Compute values for motors given input params"
	app.Action = func(c *cli.Context) {
		fmt.Println(c.String("motor"))
		fmt.Println(c.String("given"))
		fmt.Println(c.String("value"))
	}
	app.Commands = []cli.Command{
		{
			Name:  "amps",
			Usage: "compute stats of motor at a given current draw",
			Action: func(c *cli.Context) {
				fmt.Println(findMotor(c.String("motor"), motorNames, motors).Name)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "motor, m",
					Usage: "Which motor to run calcs for",
					Value: "cim",
				},
				cli.Float64Flag{
					Name:  "value",
					Usage: "The value being used",
				},
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "motor, m",
			Usage: "Which motor to run calcs for",
		},
		cli.Float64Flag{
			Name:  "value",
			Usage: "The value being used",
		},
	}
	app.Run(os.Args)
}
func computeTorque(amps float64, selectedMotor *Motor) float64 {
	return (amps - selectedMotor.FreeAmps) / (selectedMotor.StallAmps - selectedMotor.FreeAmps) * selectedMotor.StallTorque
}

func computeAmps(torque float64, selectedMotor *Motor) float64 {
	return selectedMotor.FreeAmps + (torque/selectedMotor.StallTorque)*(selectedMotor.StallAmps-selectedMotor.FreeAmps)
}

func computeSpeed(torque float64, selectedMotor *Motor) float64 {
	return selectedMotor.FreeSpeed * (1.0 - torque/selectedMotor.StallTorque)
}

func findMotor(motorName string, motorNames []string, motors []Motor) Motor {

	for _, motor := range motors {
		if motor.Name == motorName {
			return motor
		}
	}
	model := fuzzy.NewModel()
	model.Train(motorNames)
	motorName = model.Suggestions(motorName, false)[0]
	fmt.Println("Did you mean %s?", motorName)
	os.Exit(-1)
	return motors[0]
}

// selectedMotorString := os.Args[1]
// selectedMotor := new(Motor)
// for _, m := range motors {
// 	if m.Name == selectedMotorString {
// 		selectedMotor = &m
// 		break
// 	}
// }

// func compute()  {
// 	// ohms := selectedMotor.SpecVoltage / selectedMotor.StallAmps
// 	var speed float64
// 	var amps float64
// 	var powerOut float64
// 	var powerIn float64
// 	var efficiency float64
// 	var torque float64
// 	volts := 12.0
//
// 	w := new(tabwriter.Writer)
// 	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
//
// 	switch os.Args[2] {
// 	case "amps":
// 		amps, _ = strconv.ParseFloat(os.Args[3], 64)
// 		torque = computeTorque(amps, selectedMotor)
//
// 	case "torque":
// 		torque, _ = strconv.ParseFloat(os.Args[3], 64)
// 		amps = computeAmps(torque, selectedMotor)
// 	}
// 	speed = computeSpeed(torque, selectedMotor)
// 	powerOut = (speed / 60 * 2 * math.Pi) * torque
// 	powerIn = volts * amps
// 	efficiency = 100 * (powerOut / powerIn)
// 	str := fmt.Sprintf("%.2f\t%.0f\t%.2f\t%.2f\t%.2f\t%.2f", torque, speed, amps, powerOut, powerIn-powerOut, efficiency)
// 	fmt.Fprintln(w, "torque(N*m)\trpm\tamps\toutput(W)\theat(W)\teff")
// 	fmt.Fprintln(w, str)
// 	w.Flush()
// 	if err != nil {
// 		panic(err)
// 	}
// }
