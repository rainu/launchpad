package launchpad

import (
	"gitlab.com/gomidi/midi"
)

// ScrollingTextBuilder is used to build and display an scrolling text on the Launchpad.
type ScrollingTextBuilder interface {
	// Add adds a text snipped with a given speed to the builder.
	// The speed can be a value from 1-7. The text must be ASCII
	// characters! Otherwise the result could be weired.
	Add(speed byte, text string) ScrollingTextBuilder

	// Perform sends the pre-built scrolling text to the launchpad.
	Perform() error
}

type scrollingTextBuilder struct {
	Seq          []byte
	outputStream midi.Out
}

func (l *launchpad) Text(g int, r int) ScrollingTextBuilder {
	return l.text(g, r, false)
}

func (l *launchpad) TextLoop(g int, r int) ScrollingTextBuilder {
	return l.text(g, r, true)
}

func (l *launchpad) text(g int, r int, loop bool) ScrollingTextBuilder {
	color := 16*g + r + 8 + 4
	if loop {
		color += 64
	}

	return &scrollingTextBuilder{
		Seq:          []byte{0xF0, 0x00, 0x20, 0x29, 0x09, byte(color)},
		outputStream: l.outputStream,
	}
}

func (s *scrollingTextBuilder) Add(speed byte, text string) ScrollingTextBuilder {
	if speed > 7 {
		speed = 7
	} else if speed < 1 {
		speed = 1
	}

	s.Seq = append(s.Seq, speed)
	s.Seq = append(s.Seq, []byte(text)...)

	return s
}

func (s *scrollingTextBuilder) Perform() error {
	s.Seq = append(s.Seq, 0xF7)

	// the syntax of the scrolling text message:
	// F0 00 20 29 09 <colour> <text inclusive speed ...> F7
	_, err := s.outputStream.Write(s.Seq)
	return err
}
