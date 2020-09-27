package launchpad

import (
	"gitlab.com/gomidi/midi"
)

type scrollingTextBuilderS struct {
	Seq          []byte
	outputStream midi.Out
}

func (l *LaunchpadS) Text(color Color) ScrollingTextBuilder {
	return l.text(color, false)
}

func (l *LaunchpadS) TextLoop(color Color) ScrollingTextBuilder {
	return l.text(color, true)
}

func (l *LaunchpadS) text(color Color, loop bool) ScrollingTextBuilder {
	c := color.AsBytes()[0]
	if loop {
		c += 64
	}

	return &scrollingTextBuilderS{
		Seq:          []byte{0xF0, 0x00, 0x20, 0x29, 0x09, c},
		outputStream: l.outputStream,
	}
}

func (s *scrollingTextBuilderS) Add(speed byte, text string) ScrollingTextBuilder {
	if speed > 7 {
		speed = 7
	} else if speed < 1 {
		speed = 1
	}

	s.Seq = append(s.Seq, speed)
	s.Seq = append(s.Seq, []byte(text)...)

	return s
}

func (s *scrollingTextBuilderS) Perform() error {
	s.Seq = append(s.Seq, 0xF7)

	// the syntax of the scrolling text message:
	// F0 00 20 29 09 <colour> <text inclusive speed ...> F7
	_, err := s.outputStream.Write(s.Seq)
	return err
}
