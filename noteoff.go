package main

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

// sendNoteOff sends NoteOff messages to all channels and notes.
func sendNoteOff(out drivers.Out) {
	fmt.Println("Sending NoteOff messages...")
	// Iterate over all MIDI channels (0-15)
	for ch := range 16 {
		// Iterate over all MIDI notes (0-127)
		for note := range 128 {
			out.Send(midi.NoteOff(uint8(ch), uint8(note)))
		}
	}
}
