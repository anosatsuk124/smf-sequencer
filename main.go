package main

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"

	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

var activeOut drivers.Out

func listMIDIDevices() {
	drv, err := rtmididrv.New()
	if err != nil {
		fmt.Println("Failed to initialize MIDI driver:", err)
		return
	}
	defer drv.Close()

	outs, err := drv.Outs()
	if err != nil {
		fmt.Println("Failed to get MIDI outputs:", err)
		return
	}

	fmt.Println("Available MIDI Output Devices:")
	for i, out := range outs {
		fmt.Printf("[%d] %s\n", i, out.String())
	}
}

func sendNoteOff(out drivers.Out) {
	fmt.Println("Sending NoteOff messages...")
	for ch := range 16 {
		for note := range 128 {
			out.Send(midi.NoteOff(uint8(ch), uint8(note)))
		}
	}
}

func playMIDI(filePath string, deviceIndex int) {
	// Load MIDI driver
	drv, err := rtmididrv.New()
	if err != nil {
		fmt.Println("Failed to initialize MIDI driver:", err)
		return
	}
	defer drv.Close()

	// Get MIDI output devices
	outs, err := drv.Outs()
	if err != nil || len(outs) == 0 {
		fmt.Println("No MIDI output devices found")
		return
	}

	if deviceIndex < 0 || deviceIndex >= len(outs) {
		fmt.Println("Invalid device index")
		return
	}

	// Open selected MIDI output device
	out := outs[deviceIndex]
	if err := out.Open(); err != nil {
		fmt.Println("Failed to open MIDI output device:", err)
		return
	}
	defer out.Close()
	activeOut = out

	// Handle Ctrl+C signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("Interrupted! Sending NoteOff...")
		time.Sleep(100 * time.Millisecond)
		sendNoteOff(out)
		drv.Close()
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  midi-player list")
		fmt.Println("  midi-player play <file> <device_index>")
		return
	}

	switch os.Args[1] {
	case "list":
		listMIDIDevices()
	case "play":
		if len(os.Args) < 4 {
			fmt.Println("Usage: midi-player play <file> <device_index>")
			return
		}
		filePath := os.Args[2]
		var deviceIndex int
		fmt.Sscanf(os.Args[3], "%d", &deviceIndex)
		playMIDI(filePath, deviceIndex)
	default:
		fmt.Println("Unknown command")
	}
}
