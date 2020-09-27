package main

import (
	// replace with e.g. "gitlab.com/gomidi/rtmididrv" for real midi connections
	driver "gitlab.com/gomidi/midi/testdrv"
	"time"

	"github.com/rainu/launchpad"
	"log"
)

func main() {
	pad, err := launchpad.NewLaunchpad(driver.New("fake"))
	if err != nil {
		log.Fatalf("error while openning connection to launchpad: %v", err)
	}
	defer pad.Close()

	pad.Clear()

	// Set <0,8> to yellow (via RGB mode).
	pad.Light(0, 8, launchpad.ColorMK2RGB{16, 16, 0})

	// Set <1,8> to yellow (via RGB mode).
	pad.Light(1, 8, launchpad.ColorMK2{15})

	// Show the whole color palette
	c := 0
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			pad.Light(x, y, launchpad.ColorMK2{c})
			c++
		}
	}
	time.Sleep(5 * time.Second)
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			pad.Light(x, y, launchpad.ColorMK2{c})
			c++
		}
	}

	pad.Text(launchpad.ColorMK2{9}).
		Add(7, "Hello ").
		Add(1, "World!").
		Perform()

	hits, err := pad.ListenToHits()
	if err != nil {
		panic(err)
	}

	marker, err := pad.ListenToScrollTextEndMarker()
	if err != nil {
		panic(err)
	}

	for {
		select {
		case hit := <-hits:
			log.Printf("Button pressed at <x=%d, y=%d>", hit.X, hit.Y)
			// Turn to green.
			if hit.Down {
				pad.Light(hit.X, hit.Y, launchpad.ColorMK2{13})
			} else {
				pad.Light(hit.X, hit.Y, launchpad.ColorMK2RGB{0, 16, 16})
			}
		case <-marker:
			log.Print("Text ends")
		}

	}
}
