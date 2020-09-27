package launchpad

import (
	"errors"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"strings"
	"sync"
)

type LaunchpadMK2 struct {
	inputStream  midi.In
	outputStream midi.Out

	isListening   bool
	listenerMutex sync.RWMutex
	listener      []func(pos *reader.Position, msg midi.Message)
}

func (l *LaunchpadMK2) Out() midi.Out {
	return l.outputStream
}

func (l *LaunchpadMK2) In() midi.In {
	return l.inputStream
}

func (l *LaunchpadMK2) Close() error {
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

func (l *LaunchpadMK2) write(b []byte) (int, error) {
	if !l.outputStream.IsOpen() {
		if err := l.outputStream.Open(); err != nil {
			return -1, err
		}
	}

	return l.outputStream.Write(b)
}
