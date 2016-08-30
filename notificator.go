package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/op/go-logging"
)

var (
	title     string
	text      string
	sleeptime int
	logger    *logging.Logger
	err       error
	notifier  Notifications
)

type Notifications struct {
	Text       string
	Title      string
	Sleeptime  int
	ExpireTime int
	SoundFile  string
}

func (n Notifications) Notify() {
	logger.Info("Start Notification")
	for true == true {
		command := exec.Command("notify-send", "-t", strconv.Itoa(n.ExpireTime*1000), n.Title, n.Text)
		command.Run()
		command = exec.Command("aplay", n.SoundFile)
		command.Start()
		logger.Info("Notification " + n.Title + " send")
		time.Sleep(time.Duration(n.Sleeptime) * time.Second)
	}
}

// Get config from configfile
func init() {
	// Initialize logger
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	format := logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{module}.%{shortfunc}() ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}")
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
	logger = logging.MustGetLogger("Notificator")
	logging.SetLevel(logging.DEBUG, "Notificator")

	logger.Info("Initialized logging")
}

func main() {

	file, _ := os.Open("/mnt/backup/Go/src/github.com/adiran/KDENotificator/notifications.json")

	dec := json.NewDecoder(file)

	// read open bracket
	_, err := dec.Token()
	if err != nil {
		logger.Fatal(err)
	}

	// while the array contains values
	for dec.More() {

		// decode an array value (Message)
		err := dec.Decode(&notifier)
		if err != nil {
			logger.Fatal(err)
		}

		go notifier.Notify()
		logger.Info(notifier.Text)
		logger.Info(notifier.Title)
		logger.Info(notifier.Sleeptime)
	}

	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		logger.Fatal(err)
	}

	// keep thread running
	c := make(chan int)
	pla := <-c
	logger.Info(pla)
}
