# launchpad 
A package allows you to talk to your Novation Launchpad S in Go. Light buttons or read your touches.

This library is currently only working with Launchpad S (Green-Red Launchpads)

## Usage
Initialize a new Launchpad. If there are no currently connected Launchpad
device, initialization will fail with an error. You can fake a device by
creating an input and output MIDI device and name them as Launchpad.

## Drivers

For "cable" communication you need a Driver to connect with the MIDI system of your OS. 
Currently there are two multi-platform drivers available:

* package ```gitlab.com/gomidi/rtmididrv``` based on rtmidi
* package ```gitlab.com/gomidi/portmididrv``` based on portmidi

Portmidi is required to use this package.

```
$ apt-get install libportmidi-dev
# or
$ brew install portmidi
# or 
$ yay -S portmidi
```

```go
package main

import (
    "github.com/rainu/launchpad"
	driver "gitlab.com/gomidi/portmididrv"
	"log"
)
func main() {
    pad, err := launchpad.NewLaunchpad(driver.New())
    if err != nil {
        log.Fatalf("Error initializing launchpad: %v", err)
    }
    defer pad.Close()
    
    // turn off all of the lights
    pad.Clear()
}
```

### Coordinate system

The coordinate system is illustrated below.
```
+--------- arrow keys -----------+  +--- mode keys ---+
{0, 8} {1, 8} {2, 8} {3, 8} {4, 8} {5, 8} {6, 8} {7, 8} | ableton
----------------------------------------------------------------
{0, 0} {1, 0} {2, 0} {3, 0} {4, 0} {5, 0} {6, 0} {7, 0} | {8, 0} vol
----------------------------------------------------------------
{0, 1} {1, 1} {2, 1} {3, 1} {4, 1} {5, 1} {6, 1} {7, 1} | {8, 1} pan
----------------------------------------------------------------
{0, 2} {1, 2} {2, 2} {3, 2} {4, 2} {5, 2} {6, 2} {7, 2} | {8, 2} sndA
----------------------------------------------------------------
{0, 3} {1, 3} {2, 3} {3, 3} {4, 3} {5, 3} {6, 3} {7, 3} | {8, 3} sndB
----------------------------------------------------------------
{0, 4} {1, 4} {2, 4} {3, 4} {4, 4} {5, 4} {6, 4} {7, 4} | {8, 4} stop
----------------------------------------------------------------
{0, 5} {1, 5} {2, 5} {3, 5} {4, 5} {5, 5} {6, 5} {7, 5} | {8, 5} trk on
----------------------------------------------------------------
{0, 6} {1, 6} {2, 6} {3, 6} {4, 6} {5, 6} {6, 6} {7, 6} | {8, 6} solo
----------------------------------------------------------------
{0, 7} {1, 7} {2, 7} {3, 7} {4, 7} {5, 7} {6, 7} {7, 7} | {8, 7} arm
----------------------------------------------------------------
```

## Demo: Light your touchs

A simple program to light every touch:

```go
package main

import (
    "github.com/rainu/launchpad"
    driver "gitlab.com/gomidi/portmididrv"
    "log"
)
func main() {
    pad, err := launchpad.NewLaunchpad(driver.New())
    if err != nil {
        log.Fatal(err)
    }
    defer pad.Close()
    
    pad.Clear()
    
    hits, err := pad.ListenToHits()
    if err != nil {
        log.Fatal(err)
    }
    
    for {
        select {
        case hit := <-hits:
            pad.Light(hit.X, hit.Y, 3, 3)
        }
    }
}
```

## Demo: Write a scrolling text

A simple program to write text on the launchpad

```go
package main

import (
    "github.com/rainu/launchpad"
    driver "gitlab.com/gomidi/portmididrv"
    "log"
)
func main() {
    pad, err := launchpad.NewLaunchpad(driver.New())
    if err != nil {
        log.Fatal(err)
    }
    defer pad.Close()
    
    pad.Clear()
    
    // Send Text-Loop
    pad.Text(3, 0).Add(7, "Hello World!").Perform()
    
    marker, err := pad.ListenToScrollTextEndMarker()
    if err != nil {
        log.Fatal(err)
    }
    for {
        <-marker
        log.Printf("Scrolling text is ended now.")
    }
}
```
