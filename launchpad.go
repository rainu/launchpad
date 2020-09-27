package launchpad

import (
	"errors"
	"fmt"
	"gitlab.com/gomidi/midi"
	"strings"
)

// Launchpad represents a device with an input and output MIDI stream.
type Launchpad interface {

	// ListenToHits listens the input stream for hits.
	// It will return an error if listening initialisation failed.
	ListenToHits() (<-chan Hit, error)

	// ListenToScrollTextEndMarker listens the input stream for text end marker events.
	// It will return an error if listening initialisation failed.
	ListenToScrollTextEndMarker() (<-chan interface{}, error)

	// Light lights the button at x,y with the given color.
	// x and y are [0, 8] and color can be a ColorS or ColorMK2 / ColorMK2RGB (depends on connected launchpad)
	// Note that x=8 corresponds to the round scene buttons on the right side of the device,
	// and y=8 corresponds to the round buttons on the top of the device.
	Light(x, y int, color Color) error

	// Text will return a scrolling text builder whether you can build and
	// perform an text with the given color (for Launchpad MK2 only ColorMK2 will work) which will be scrolled on the launchpad.
	Text(color Color) ScrollingTextBuilder

	// TextLoop will return a scrolling text builder whether you can build and
	// perform an text with the given color (for Launchpad MK2 only ColorMK2 will work) which will be scrolled endless on the launchpad.
	// If you want to stop an text loop you have to build and execute an empty textLoop!
	TextLoop(color Color) ScrollingTextBuilder

	// Clear turns off all the LEDs on the Launchpad.
	Clear() error

	// Close will close all underlying resources such like the streams and so on.
	// It will return an error if any of the underlying resources will return an error
	// on closing.
	Close() error

	// Out will return the underlying midi output stream.
	Out() midi.Out

	// In will return the underlying midi input stream.
	In() midi.In
}

//Color can be ColorS for "Launchpad S"-Devices and ColorMK2 or ColorMK2RGB for "Launchpad MK2"-Devices
type Color interface {
	AsBytes() []byte
}

type Hit struct {
	X    int
	Y    int
	Down bool
}

// ScrollingTextBuilder is used to build and display an scrolling text on the Launchpad.
type ScrollingTextBuilder interface {
	// Add adds a text snipped with a given speed to the builder.
	// The speed can be a value from 1-7. The text must be ASCII
	// characters! Otherwise the result could be weired.
	Add(speed byte, text string) ScrollingTextBuilder

	// Perform sends the pre-built scrolling text to the launchpad.
	Perform() error
}

// NewLaunchpad will create a new Launchpad instance. It will discover for connected
// Launchpads. If there is no Launchpad found, an error will returned.
func NewLaunchpad(driver midi.Driver) (Launchpad, error) {
	backend, err := discover(driver)
	if err != nil {
		return nil, err
	}

	return backend, nil
}

func discover(driver midi.Driver) (Launchpad, error) {
	ins, err := driver.Ins()
	if err != nil {
		return nil, err
	}
	outs, err := driver.Outs()
	if err != nil {
		return nil, err
	}

	var in midi.In
	var out midi.Out

	for i, _ := range ins {
		if strings.Contains(ins[i].String(), "Launchpad") {
			in = ins[i]
			break
		}
	}
	for i, _ := range outs {
		if strings.Contains(outs[i].String(), "Launchpad") {
			out = outs[i]
			break
		}
	}

	if in == nil {
		return nil, fmt.Errorf("no lanchpad input stream connected")
	}
	if out == nil {
		return nil, fmt.Errorf("no lanchpad output stream connected")
	}

	if strings.Contains(in.String(), "Launchpad S") {
		return &LaunchpadS{
			inputStream:  in,
			outputStream: out,
		}, nil
	} else if strings.Contains(in.String(), "Launchpad MK2") {
		// Switch to the session mode.
		out.Write([]byte{0xf0, 0x00, 0x20, 0x29, 0x02, 0x18, 0x22, 0x00, 0xf7})

		return &LaunchpadMK2{
			inputStream:  in,
			outputStream: out,
		}, nil
	}

	return nil, errors.New("unsupported launchpad type")
}
