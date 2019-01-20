package main

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type player interface {
	play()
	isPlaying() bool
	stop()
}

type beepPlayer struct {
	buf    *beep.Buffer
	stream beep.StreamSeeker
}

func (p *beepPlayer) play() {
	speaker.Clear()
	p.stream = p.buf.Streamer(0, p.buf.Len())
	speaker.Play(p.stream)
}

func (p *beepPlayer) isPlaying() bool {
	return (p.stream.Position() != p.buf.Len()) &&
		p.stream.Position() != 0
}

func (p *beepPlayer) stop() {
	speaker.Clear()
	p.stream = p.buf.Streamer(0, p.buf.Len())
}

func newPlayer(file *os.File) player {
	streamer, format, err := wav.Decode(file)
	if err != nil {
		log.Fatal("Error decoding audio file: ", err)
	}

	buf := beep.NewBuffer(format)
	buf.Append(streamer)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	return &beepPlayer{
		buf:    buf,
		stream: buf.Streamer(0, buf.Len()),
	}
}
