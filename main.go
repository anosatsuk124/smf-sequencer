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
        fmt.Println("  midi-player play - <device_index>    # read SMF from stdin")
        return
    }

	switch os.Args[1] {
	case "list":
		listMIDIDevices()
    case "play":
        if len(os.Args) < 4 {
            fmt.Println("Usage:")
            fmt.Println("  midi-player play <file> <device_index>")
            fmt.Println("  midi-player play - <device_index>    # read SMF from stdin")
            return
        }
        filePath := os.Args[2]
        var deviceIndex int
        fmt.Sscanf(os.Args[3], "%d", &deviceIndex)
        if filePath == "-" {
            playMIDIReader(os.Stdin, deviceIndex)
        } else {
            playMIDI(filePath, deviceIndex)
        }
	default:
		fmt.Println("Unknown command")
	}
}
