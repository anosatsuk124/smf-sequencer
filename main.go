/*
Package main provides a command-line MIDI player.
It supports listing available MIDI devices and playing MIDI files.
*/
package main

import (
	"fmt"
	"os"
)

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
