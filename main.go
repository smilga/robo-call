package main

import (
	"flag"
	"log"
	"os/exec"
)

const (
	updateTimout = 1000
)

func main() {
	_, err := exec.LookPath("adb")
	if err != nil {
		log.Fatal("ADB not installed")
	}
	flag.Parse()

	
	numbers := readCSVFile(*flag.String("file", "phones.csv", "CSV with telephone nrs"))
	audioFile := readFile(*flag.String("audio", "sound.wav", "wav audio to playback when connected"))

	robo := newRoboCaller(audioFile, numbers)
	robo.start()
}
