package launchpad

//ColorMK2 represents a color from the "Launchpad MK2"s color palette
type ColorMK2 struct {
	//Code represents the color code of color palette.
	//It can be from 0 to 127
	Code int
}

func (c ColorMK2) AsBytes() []byte {
	return []byte{
		byte(c.Code),
	}
}

//ColorMK2RGB represents a rgb-color for the "Launchpad MK2"
type ColorMK2RGB struct {
	//Red part of color. It can be from 0 to 63!
	Red int

	//Green part of color. It can be from 0 to 63!
	Green int

	//Blue part of color. It can be from 0 to 63!
	Blue int
}

func (c ColorMK2RGB) AsBytes() []byte {
	return []byte{
		byte(c.Red),
		byte(c.Green),
		byte(c.Blue),
	}
}

func (l *LaunchpadMK2) Light(x, y int, color Color) error {
	note := (8-y)*10 + x + 1
	if y >= 8 {
		note = x + 104
	}

	velocity := color.AsBytes()
	if len(velocity) == 1 {
		//simple color (color palette)
		code := 0x90
		if y >= 8 {
			code = 0xB0
		}

		_, err := l.write([]byte{byte(code), byte(note), velocity[0]})
		return err
	}

	_, err := l.write([]byte{0xF0, 0x00, 0x20, 0x29, 0x02, 0x18, 0x0B, byte(note), velocity[0], velocity[1], velocity[2], 0xF7})
	return err
}

func (l *LaunchpadMK2) Clear() error {
	_, err := l.write([]byte{0xF0, 0x00, 0x20, 0x29, 0x02, 0x18, 0x0E, 0x00, 0xF7})
	return err
}
