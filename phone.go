package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type phone interface {
	call(string) error
	drop()
	status() status
}

type status = string

// Call statuses
const (
	IDLE         status = "IDLE"
	DIALING      status = "DIALING"
	ACTIVE       status = "ACTIVE"
	CONNECTING   status = "CONNECTING"
	DISCONNECTED status = "DISCONNECTED"
	RINGING      status = "RINGING"
)

var (
	callCmd    = []string{"shell", "am", "start", "-a", "android.intent.action.CALL", "-d"}
	endCallCmd = []string{"shell", "input", "keyevent", "KEYCODE_ENDCALL"}
	dumpCmd    = []string{"shell", "dumpsys", "telecom", "|", "head", "|", "sed", "-n", "3p", "|", "awk", "'{print $2}'"}
)

type telephone struct {
	log chan string
}

func (t *telephone) call(number string) error {
	cmd := append(callCmd, fmt.Sprintf("tel:+371%s", number))
	_, err := exec.Command("adb", cmd...).Output()
	if err != nil {
		return err
	}
	return nil
}

func (t *telephone) drop() {
	_, err := exec.Command("adb", endCallCmd...).Output()
	if err != nil {
		log.Fatal("Error droping phone: ", err)
	}
}

func (t *telephone) status() status {
	bb, err := exec.Command("adb", dumpCmd...).Output()
	if err != nil {
		log.Fatal("Error getting status: ", err)
	}
	statusStr := strings.TrimRight(strings.TrimSpace(string(bb)), ",")
	if statusStr == "" {
		statusStr = "IDLE"
	}
	return status(statusStr)
}

func newTelephone() phone {
	return &telephone{}
}
