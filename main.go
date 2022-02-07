package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	ui "github.com/gizak/termui/v3"
	"github.com/mjibson/go-dsp/fft"
)


const numSamples = 48000
// tinker with these numbers
const peakFalloff = .3
const maxHeight = 20
const spectrumWidth = 100
const spectrumOffset = 0

var freqSpectrum = make([]float64, spectrumWidth)
var redraw chan string

type CustomSteamer struct {
	streamer beep.Streamer
}	

func (cs *CustomSteamer) Stream(samples [][2]float64) (n int, ok bool) {
	filled := 0
	for filled < len(samples) {
		n, ok := cs.streamer.Stream(samples[filled:])
		if (!ok ) {
			return len(samples), false
		}

		updateSpectrumValues(numSamples, samples[filled:], freqSpectrum)

		filled += n
		redraw <- "redraw that bitch"
	}
	return len(samples), true
}

func (cs *CustomSteamer) Err() error {
	return nil
}

func updateSpectrumValues(numberOfSamples int, samples [][2]float64, freqSpectrum []float64) {

	singleChannel := make([]float64, len(samples))
	for i := 0; i < len(samples); i++ {
		singleChannel[i] = samples[i][0]
	}

	fftOutput := fft.FFTReal(singleChannel)

	for i := 0; i < spectrumWidth; i++ {
		fr := real(fftOutput[i])
		fi := imag(fftOutput[i])
		magnitude := math.Sqrt(fr*fr + fi*fi)
		val := math.Min(maxHeight, math.Abs(magnitude))
		if freqSpectrum[i] > val {
			freqSpectrum[i] = math.Max(freqSpectrum[i]-peakFalloff, 0.0)
		} else {
			freqSpectrum[i] = (val + freqSpectrum[i]) / 2.0
		}
	}
}


func main() {
	redraw = make(chan string)
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	
	defer ui.Close()

	sampleRate := beep.SampleRate(numSamples)
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	

	var file string
    flag.StringVar(&file, "file", "", "mp3 filename")
	flag.Parse()

	f, err := os.Open("audio/" + file + ".mp3")
	if err != nil {
		fmt.Println(err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Println(err)
	}
	
	// normalizes song quality 
	resampled := beep.Resample(4, format.SampleRate, numSamples, streamer)
	speaker.Play(&CustomSteamer{streamer: resampled})

	top := buildTopBars(freqSpectrum[spectrumOffset:spectrumWidth], f.Name())
	bottom := buildBottomBars(freqSpectrum[spectrumOffset:spectrumWidth])

	uiEvents := ui.PollEvents()

	for {
		select {
			case e := <-uiEvents:
				switch e.ID {
				case "q", "<C-c>":
					return
				}
			case <-redraw:
				if (freqSpectrum[0] != 0) {
					newData := decay(top.Data, freqSpectrum[spectrumOffset: spectrumWidth])
					top.Data = newData
					bottom.Data = newData

					// bottom bar has to be first?
					// ui.Render(bottom, top)
					ui.Render(bottom)

					// ui.Render(top)

				}
			}
	}
}

// keeps bars from flashing as much
func decay(prevData [] float64, newData [] float64) []float64 {
	decayedData := make([]float64, spectrumWidth - spectrumOffset)
	for i, num := range newData {
		if newData[i] == 0 && prevData[i] > 0 {
			decayedData[i] = prevData[i]
		} else {
			decayedData[i] = num
		}
    }
	return decayedData
}
