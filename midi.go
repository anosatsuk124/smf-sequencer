package main

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

// listMIDIDevices lists available MIDI output devices.
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

// getMIDIOut retrieves a MIDI output device by index.
// It returns the MIDI output device and an error, if any.
func getMIDIOut(deviceIndex int) (drivers.Out, error) {
	drv, err := rtmididrv.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MIDI driver: %w", err)
	}

	outs, err := drv.Outs()
	if err != nil {
		drv.Close()
		return nil, fmt.Errorf("failed to get MIDI outputs: %w", err)
	}

	if deviceIndex < 0 || deviceIndex >= len(outs) {
		drv.Close()
		return nil, fmt.Errorf("invalid device index")
	}

	out := outs[deviceIndex]
	if err := out.Open(); err != nil {
		drv.Close()
		return nil, fmt.Errorf("failed to open MIDI output device: %w", err)
	}

	// Close the driver when the function returns
	// This ensures that the driver is closed even if there is an error
	return out, nil
}
