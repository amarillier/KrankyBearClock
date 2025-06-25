package main

import (
	"log"
	"os"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/generators"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
)

func playMp3(name string) {
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
		playBeep("up")
		return
	}

	if debug == 1 {
		log.Println("playing: ", name)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
		playBeep("up")
		return
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}

func playMid(name string) {
	// not using for now - soundfont is ridiculously big
	return
	/*
		var sampleRate beep.SampleRate = 44100

		err := speaker.Init(sampleRate, sampleRate.N(time.Second/30))
		if err != nil {
			log.Fatal(err)
		}

		// Load a soundfont.
		soundFontFile, err := os.Open("Florestan-Basic-GM-GS-by-Nando-Florestan(Public-Domain).sf2")
		if err != nil {
			log.Fatal(err)
		}
		soundFont, err := midi.NewSoundFont(soundFontFile)
		if err != nil {
			log.Fatal(err)
		}

		// Load a midi track
		midiFile, err := os.Open(name)
		if err != nil {
			log.Fatal(err)
		}
		s, format, err := midi.Decode(midiFile, soundFont, sampleRate)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Song duration: %v\n", format.SampleRate.D(s.Len()))
		speaker.PlayAndWait(s)
	*/
}

func playWav(name string) {
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
		playBeep("up")
		return
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
		playBeep("up")
		return
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}

func playBeep(style string) {
	// accept updown, up, down, ding
	sr := beep.SampleRate(48000)
	speaker.Init(sr, 4800)

	ch := make(chan struct{})
	buzzer1, _ := generators.SawtoothTone(sr, float64(750))
	buzzer2, _ := generators.SawtoothTone(sr, float64(850))
	buzzer3, _ := generators.SawtoothTone(sr, float64(950))
	buzzer4, _ := generators.SawtoothTone(sr, float64(1050))
	buzzer5, _ := generators.SawtoothTone(sr, float64(1150))
	// Play 1/n second of each tone
	t := sr.N(time.Second / 10)
	f := sr.N(time.Second / 5)
	switch style {
	case "updown":
		buzz := []beep.Streamer{
			beep.Take(t, buzzer1),
			beep.Take(t, buzzer2),
			beep.Take(t, buzzer3),
			beep.Take(t, buzzer4),
			beep.Take(t, buzzer5),
			beep.Take(t, buzzer4),
			beep.Take(t, buzzer3),
			beep.Take(t, buzzer2),
			beep.Take(f, buzzer1),
			beep.Callback(func() {
				ch <- struct{}{}
			}),
		}
		speaker.Play(beep.Seq(buzz...))
		<-ch
	case "up":
		buzz := []beep.Streamer{
			beep.Take(t, buzzer1),
			beep.Take(t, buzzer2),
			beep.Take(t, buzzer3),
			beep.Take(t, buzzer4),
			beep.Take(t, buzzer5),
			beep.Callback(func() {
				ch <- struct{}{}
			}),
		}
		speaker.Play(beep.Seq(buzz...))
		<-ch
	case "down":
		buzz := []beep.Streamer{
			beep.Take(t, buzzer5),
			beep.Take(t, buzzer4),
			beep.Take(t, buzzer3),
			beep.Take(t, buzzer2),
			beep.Take(f, buzzer1),
			beep.Callback(func() {
				ch <- struct{}{}
			}),
		}
		speaker.Play(beep.Seq(buzz...))
		<-ch
	case "ding":
		t = sr.N(time.Second / 4)
		buzzer1, _ := generators.SawtoothTone(sr, float64(350))
		buzz := []beep.Streamer{
			beep.Take(t, buzzer1),
			beep.Callback(func() {
				ch <- struct{}{}
			}),
		}
		speaker.Play(beep.Seq(buzz...))
		<-ch
	}
}

// "Now this is not the end. It is not even the beginning of the end. But it is, perhaps, the end of the beginning." Winston Churchill, November 10, 1942
