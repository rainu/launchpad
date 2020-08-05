package launchpad

import (
	"errors"
	"fmt"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"strings"
	"sync"
)

// Launchpad represents a device with an input and output MIDI stream.
type Launchpad interface {

	// ListenToHits listens the input stream for hits.
	// It will return an error if listening initialisation failed.
	ListenToHits() (<-chan Hit, error)

	// ListenToScrollTextEndMarker listens the input stream for text end marker events.
	// It will return an error if listening initialisation failed.
	ListenToScrollTextEndMarker() (<-chan interface{}, error)

	// Light lights the button at x,y with the given green and red values.
	// x and y are [0, 8], g and r are [0, 3]
	// Note that x=8 corresponds to the round scene buttons on the right side of the device,
	// and y=8 corresponds to the round buttons on the top of the device.
	Light(x, y, g, r int) error

	// Text will return a scrolling text builder whether you can build and
	// perform an text with the given color which will be scrolled on the launchpad.
	Text(g int, r int) ScrollingTextBuilder

	// TextLoop will return a scrolling text builder whether you can build and
	// perform an text with the given color which will be scrolled endless on the launchpad.
	// If you want to stop an text loop you have to build and execute an empty textLoop!
	TextLoop(g int, r int) ScrollingTextBuilder

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

type launchpad struct {
	inputStream  midi.In
	outputStream midi.Out

	isListening   bool
	listenerMutex sync.RWMutex
	listener      []func(pos *reader.Position, msg midi.Message)
}

// NewLaunchpad will create a new Launchpad instance. It will discover for connected
// Launchpads. If there is no Launchpad found, an error will returned.
func NewLaunchpad(driver midi.Driver) (Launchpad, error) {
	in, out, err := discover(driver)
	if err != nil {
		return nil, err
	}

	return &launchpad{
		inputStream:  in,
		outputStream: out,
	}, nil
}

func discover(driver midi.Driver) (midi.In, midi.Out, error) {
	ins, err := driver.Ins()
	if err != nil {
		return nil, nil, err
	}
	outs, err := driver.Outs()
	if err != nil {
		return nil, nil, err
	}

	var in midi.In
	var out midi.Out

	for i, _ := range ins {
		if strings.Contains(ins[i].String(), "Launchpad S") {
			in = ins[i]
			break
		}
	}
	for i, _ := range outs {
		if strings.Contains(outs[i].String(), "Launchpad S") {
			out = outs[i]
			break
		}
	}

	if in == nil {
		return nil, nil, fmt.Errorf("no lanchpad input stream connected")
	}
	if out == nil {
		return nil, nil, fmt.Errorf("no lanchpad output stream connected")
	}

	return in, out, nil
}

func (l *launchpad) Out() midi.Out {
	return l.outputStream
}

func (l *launchpad) In() midi.In {
	return l.inputStream
}

func (l *launchpad) Close() error {
	errMessages := make([]string, 0, 3)

	if err := l.inputStream.Close(); err != nil {
		errMessages = append(errMessages, err.Error())
	}
	if err := l.outputStream.Close(); err != nil {
		errMessages = append(errMessages, err.Error())
	}

	if len(errMessages) > 0 {
		return errors.New(strings.Join(errMessages, ";"))
	}

	return nil
}

func (l *launchpad) write(b []byte) (int, error) {
	if !l.outputStream.IsOpen() {
		if err := l.outputStream.Open(); err != nil {
			return -1, err
		}
	}

	return l.outputStream.Write(b)
}
