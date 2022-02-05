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

type CustomSteamer struct {
	streamer beep.Streamer
}

func (cs *CustomSteamer) Add(streamer beep.Streamer) {
	cs.streamer = streamer
}

const spectrumWidth = 314
const maxHeight = 20
var freqSpectrum = make([]float64, 512)

func updateSpectrumValues(numberOfSamples int, samples [][2]float64, maxValue float64, freqSpectrum []float64) {

	singleChannel := make([]float64, len(samples))
	for i := 0; i < len(samples); i++ {
		singleChannel[i] = samples[i][0]
	}

	fftOutput := fft.FFTReal(singleChannel)

	for i := 0; i < spectrumWidth; i++ {
		fr := real(fftOutput[i])
		fi := imag(fftOutput[i])
		magnitude := math.Sqrt(fr*fr + fi*fi)
		val := math.Min(maxValue, math.Abs(magnitude))
		if freqSpectrum[i] > val {
			freqSpectrum[i] = math.Max(freqSpectrum[i]-8.0, 0.0)
		} else {
			freqSpectrum[i] = (val + freqSpectrum[i]) / 2.0
		}
	}
}

func (cs *CustomSteamer) Stream(samples [][2]float64) (n int, ok bool) {
	filled := 0
	// filled is broke
	// fmt.Println("filled , ", filled)
	for filled < len(samples) {
	


		// should this ok be handled ?
		n, _ := cs.streamer.Stream(samples[filled:])
		
		updateSpectrumValues(44100, samples[filled:], maxHeight, freqSpectrum)

		filled += n
		redraw <- "redraw that bitch"
	}
	return len(samples), true
}

func (cs *CustomSteamer) Err() error {
	return nil
}

var redraw chan string
var customStreamer CustomSteamer

func main() {
	redraw = make(chan string)
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	sampleRate := beep.SampleRate(44100)
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	


	var file string
    flag.StringVar(&file, "file", "", "mp3 filename")
	flag.Parse()

	f, err := os.Open("audio/" + file)
	if err != nil {
		fmt.Println(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Println(err)
	}

	resampled := beep.Resample(4, format.SampleRate, sampleRate, streamer)
	customStreamer.streamer = resampled

	speaker.Play(&customStreamer)

	bar := buildBars(freqSpectrum[0:spectrumWidth], "temp")
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
					bar.Data = freqSpectrum[0:spectrumWidth]
					ui.Render(bar)
				}
			}
	}
}
