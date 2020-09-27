package launchpad

import (
	"gitlab.com/gomidi/midi"
)

type scrollingTextBuilderMK2 struct {
	Seq          []byte
	outputStream midi.Out
}

func (l *LaunchpadMK2) Text(color Color) ScrollingTextBuilder {
	return l.text(color, false)
}

func (l *LaunchpadMK2) TextLoop(color Color) ScrollingTextBuilder {
	return l.text(color, true)
}

func (l *LaunchpadMK2) text(color Color, loop bool) ScrollingTextBuilder {
	colorCode := color.AsBytes()

	bLoop := byte(0x00)
	if loop {
		bLoop = 0x01
	}

	return &scrollingTextBuilderMK2{
		Seq:          []byte{0xF0, 0x00, 0x20, 0x29, 0x02, 0x18, 0x14, colorCode[0], bLoop},
		outputStream: l.outputStream,
	}
}

func (s *scrollingTextBuilderMK2) Add(speed byte, text string) ScrollingTextBuilder {
	if speed > 7 {
		speed = 7
	} else if speed < 1 {
		speed = 1
	}

	s.Seq = append(s.Seq, speed)
	s.Seq = append(s.Seq, []byte(text)...)

	return s
}

func (s *scrollingTextBuilderMK2) Perform() error {
	s.Seq = append(s.Seq, 0xF7)

	// the syntax of the scrolling text message:
	// F0 00 20 29 02 18 14h <Colour> <Loop> <Text> F7
	_, err := s.outputStream.Write(s.Seq)
	return err
}
