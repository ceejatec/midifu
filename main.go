package main

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"

	// replace with e.g. "gitlab.com/gomidi/rtmididrv" for real midi connections
	"gitlab.com/gomidi/midi/writer"
	driver "gitlab.com/gomidi/rtmididrv"
)

// This example reads from the first input and and writes to the first output port
func main() {
	drv, err := driver.New()
	must(err)

	// make sure to close all open ports at the end
	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	fmt.Println(ins)
	outs, err := drv.Outs()
	must(err)
	fmt.Println(outs)
	in, out := ins[1], outs[1]

	must(in.Open())
	must(out.Open())

	defer in.Close()
	defer out.Close()

	// the writer we are writing to
	wr := writer.New(out)

	var velcount uint8 = 0

	// to disable logging, pass mid.NoLogger() as option
	rd := reader.New(
		reader.NoLogger(),
		// write every message to the out port
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			fmt.Printf("got %s\n", msg)
		}),
		reader.NoteOn(func(pos *reader.Position, channel uint8, key uint8, velocity uint8) {
			fmt.Printf("channel %v, key %v, vel %v\n", channel, key, velcount)
			writer.NoteOn(wr, key, velcount)
			velcount++
		}),
	)

	// listen for MIDI
	err = rd.ListenTo(in)
	must(err)

	err = writer.NoteOn(wr, 60, 100)
	must(err)

	time.Sleep(30000 * time.Millisecond)
	err = writer.NoteOff(wr, 60)

	must(err)
	// Output: got channel.NoteOn channel 0 key 60 velocity 100
	// got channel.NoteOff channel 0 key 60
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
