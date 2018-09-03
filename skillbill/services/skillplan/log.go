package main

import (
	"os"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	Info 	string = "Info"
	Error 	string = "Error"
	Warning	string = "Warning"
	Fatal	string = "Fatal"
	Panic	string = "Panic"
)

const (
	filepath	string = "/var/log/skillplan.log"
)

func setUpLogging(fileName string){

	// Create the log file if doesn't exist. And append to it if it already exists.
    f, err := os.OpenFile(fileName, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
    Formatter := new(log.TextFormatter)
    // You can change the Timestamp format. But you have to use the same date and time.
    // "2006-02-02 15:04:06" Works. If you change any digit, it won't work
	// ie "Mon Jan 2 15:04:05 MST 2006" is the reference time. You can't change it
	
	currentTime := time.Now()
    Formatter.TimestampFormat = currentTime.Format("2006-01-02 15:04:06")
    Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)

	if err != nil {
        // Cannot open log file. Logging to stderr
        fmt.Println(err)
    }else{
        log.SetOutput(f)
    }
}