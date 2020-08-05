package launchpad

func (l *launchpad) Light(x, y, g, r uint8) error {
	note := x + 16*y
	velocity := 16*g + r + 8 + 4
	if y >= 8 {
		note = x + 104

		_, err := l.write([]byte{0xB0, note, velocity})
		return err
	}

	_, err := l.write([]byte{0x90, note, velocity})
	return err
}

func (l *launchpad) Clear() error {
	_, err := l.write([]byte{0xB0, 0x00, 0x00})
	return err
}
