package main

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/gomidi/midi/v2/smf"
)

// playMIDI plays a MIDI file.
// It takes the file path and device index as input.
func playMIDI(filePath string, deviceIndex int) {
	// Open selected MIDI output device
	out, err := getMIDIOut(deviceIndex)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()

	// Handle Ctrl+C signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("Interrupted! Sending NoteOff...")
		time.Sleep(100 * time.Millisecond)
		sendNoteOff(out)
		os.Exit(0)
	}()

	// Read MIDI file
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Failed to read MIDI file:", err)
		return
	}

	rd := bytes.NewReader(data)
	fmt.Println("Playing MIDI file...")
	smf.ReadTracksFrom(rd).Do(func(ev smf.TrackEvent) {
		fmt.Printf("track %v @%vms %s\n", ev.TrackNo, ev.AbsMicroSeconds/1000, ev.Message)
	}).Play(out)

	fmt.Println("Playback finished")
}
