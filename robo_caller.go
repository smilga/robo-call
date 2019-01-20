package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// RoboCaller is main application controller
type RoboCaller struct {
	player            player
	phone             phone
	phoneBook         phoneBook
	statusCh          chan status
	currentContact    *contact
	processedContacts []*contact
}

func (r *RoboCaller) nextNuber() string {
	if r.currentContact != nil {
		r.processedContacts = append(r.processedContacts, r.currentContact)
	}

	next := r.phoneBook.nextNumber()
	r.currentContact = &contact{
		number: next,
	}
	return next
}

func (r *RoboCaller) callNext() {
	err := r.phone.call(r.nextNuber())
	if err != nil {
		log.Fatal(err)
	}
}

func (r *RoboCaller) exit() {
	close(r.statusCh)
}

func (r *RoboCaller) start() {
	go r.updateStatus()

	// make specifications for every case
	for callStatus := range r.statusCh {
		switch {
		case callStatus == RINGING:
			// incoming call, drop!
			r.phone.drop()
		case callStatus == DIALING:
			// just wait and rest to dial
		case callStatus == IDLE && r.phoneBook.hasNumbersLeft():
			fmt.Println("call next")
			// IDLE can call next
			r.callNext()
		case callStatus == IDLE && r.currentContact.inProcess:
			// call has been droped
			fmt.Println("call droped")
			r.player.stop()
			r.currentContact.inProcess = false
		case callStatus == ACTIVE && r.currentContact.inProcess && !r.player.isPlaying():
			fmt.Println("drop call")
			// connection active but record have been finished playing
			r.player.stop()
			r.currentContact.inProcess = false
			r.phone.drop()
		case callStatus == ACTIVE && !r.currentContact.inProcess && !r.player.isPlaying():
			fmt.Println("play sound")
			// not playing and current contact havent heard record
			r.player.play()
			r.currentContact.inProcess = true
		}
	}

	fmt.Println("Application finished, exiting")
}

func (r *RoboCaller) updateStatus() {
	for {
		status := r.phone.status()
		fmt.Println("Call status: ", status)
		r.statusCh <- status
		time.Sleep(time.Millisecond * 1000)
	}
}

// NewRoboCaller returns new RoboCaller
func newRoboCaller(audioFile *os.File, numbers []string) *RoboCaller {
	return &RoboCaller{
		player:    newPlayer(audioFile),
		phone:     newTelephone(),
		phoneBook: newPhoneBook(numbers),
		statusCh:  make(chan status),
	}
}
