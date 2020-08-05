package main

import (
	// replace with e.g. "gitlab.com/gomidi/rtmididrv" for real midi connections
	driver "gitlab.com/gomidi/midi/testdrv"

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

	// Set <0,0> to yellow.
	pad.Light(0, 0, 2, 2)

	pad.Text(3, 0).
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
				pad.Light(hit.X, hit.Y, 0, 3)
			} else {
				pad.Light(hit.X, hit.Y, 3, 0)
			}
		case <-marker:
			log.Print("Text ends")
		}

	}
}
