build:
	go-bindata -o motors.go motors.csv
	go build 
debug:
	go-bindata -o motors.go --debug motors.csv
