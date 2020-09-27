package launchpad

import (
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/reader"
)

func (l *LaunchpadMK2) listen() error {
	rd := reader.New(
		reader.NoLogger(),
		// write every message to the out port
		reader.Each(l.handleMidiMessage),
	)

	if !l.inputStream.IsOpen() {
		if err := l.inputStream.Open(); err != nil {
			return err
		}
	}

	if err := rd.ListenTo(l.inputStream); err != nil {
		return err
	}

	l.isListening = true
	return nil
}

func (l *LaunchpadMK2) handleMidiMessage(pos *reader.Position, msg midi.Message) {
	l.listenerMutex.RLock()
	defer l.listenerMutex.RUnlock()

	for i := range l.listener {
		//delegate to listener
		l.listener[i](pos, msg)
	}
}

func (l *LaunchpadMK2) addMidiMessageListener(listener func(pos *reader.Position, msg midi.Message)) {
	l.listenerMutex.Lock()
	defer l.listenerMutex.Unlock()

	l.listener = append(l.listener, listener)
}

func (l *LaunchpadMK2) ListenToHits() (<-chan Hit, error) {
	if !l.isListening {
		if err := l.listen(); err != nil {
			return nil, err
		}
	}

	hitChan := make(chan Hit) //unbuffered/blocking channel

	l.addMidiMessageListener(func(pos *reader.Position, msg midi.Message) {
		isHit := false
		hit := Hit{
			X:    0,
			Y:    0,
			Down: false,
		}

		if controlChange, ok := msg.(channel.ControlChange); ok {
			if controlChange.Controller() >= 104 && controlChange.Controller() <= 112 {
				isHit = true

				hit.X = int(controlChange.Controller()) - 104
				hit.Y = 8
				hit.Down = controlChange.Value() == 127
			}
		} else if noteOn, ok := msg.(channel.NoteOn); ok {
			isHit = true
			hit.Down = true

			hit.X = int(noteOn.Key())%10 - 1
			hit.Y = 8 - (int(noteOn.Key())-hit.X)/10
		} else if noteOn, ok := msg.(channel.NoteOff); ok {
			isHit = true
			hit.Down = false

			hit.X = int(noteOn.Key())%10 - 1
			hit.Y = 8 - (int(noteOn.Key())-hit.X)/10
		}

		if isHit {
			hitChan <- hit
		}
	})

	return hitChan, nil
}

func (l *LaunchpadMK2) ListenToScrollTextEndMarker() (<-chan interface{}, error) {
	if !l.isListening {
		if err := l.listen(); err != nil {
			return nil, err
		}
	}

	markerChan := make(chan interface{}) //unbuffered/blocking channel

	l.addMidiMessageListener(func(pos *reader.Position, msg midi.Message) {
		if controlChange, ok := msg.(channel.ControlChange); ok && controlChange.Value() == 3 {
			markerChan <- true
		}
	})

	return markerChan, nil
}
