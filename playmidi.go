package main

import (
    "fmt"
    "io"
    "os"
    "os/signal"
    "syscall"
    "time"

    "gitlab.com/gomidi/midi/v2/smf"
)

// playMIDI plays a MIDI file from a file path.
// It opens the file and delegates to playMIDIReader.
func playMIDI(filePath string, deviceIndex int) {
    f, err := os.Open(filePath)
    if err != nil {
        fmt.Println("Failed to open MIDI file:", err)
        return
    }
    defer f.Close()
    playMIDIReader(f, deviceIndex)
}

// playMIDIReader plays a MIDI stream from the provided reader.
func playMIDIReader(r io.Reader, deviceIndex int) {
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

    fmt.Println("Playing MIDI...")
    smf.ReadTracksFrom(r).Do(func(ev smf.TrackEvent) {
        fmt.Printf("track %v @%vms %s\n", ev.TrackNo, ev.AbsMicroSeconds/1000, ev.Message)
    }).Play(out)

    fmt.Println("Playback finished")
}
