package launchpad

//ColorS represents a color for the "Launchpad S"
type ColorS struct {
	//Red part of color. It can be from 0 to 3!
	Red int

	//Green part of color. It can be from 0 to 3!
	Green int
}

func (c ColorS) AsBytes() []byte {
	return []byte{byte(16*c.Green + c.Red + 8 + 4)}
}

func (l *LaunchpadS) Light(x, y int, color Color) error {
	note := x + 16*y
	if y >= 8 {
		note = x + 104

		_, err := l.write([]byte{0xB0, byte(note), color.AsBytes()[0]})
		return err
	}

	_, err := l.write([]byte{0x90, byte(note), color.AsBytes()[0]})
	return err
}

func (l *LaunchpadS) Clear() error {
	_, err := l.write([]byte{0xB0, 0x00, 0x00})
	return err
}
