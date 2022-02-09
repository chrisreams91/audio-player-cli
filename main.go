package main

import (
	"flag"
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
const maxHeight = 20
const spectrumWidth = 150
const spectrumOffset = 5
var decay = 0.2

var freqSpectrum = make([]float64, spectrumWidth)
var redraw chan string
var exit chan bool


type CustomSteamer struct {
	streamer beep.Streamer
}	

func (cs *CustomSteamer) Stream(samples [][2]float64) (n int, ok bool) {
	filled := 0
	for filled < len(samples) {
		n, ok := cs.streamer.Stream(samples[filled:])
		if (!ok ) {
			return 0, false
		}

		updateSpectrumValues(samples[filled:], freqSpectrum)

		filled += n
		redraw <- "redraw that bitch"
	}
	return len(samples), true
}

func (cs *CustomSteamer) Err() error {
	return nil
}

func updateSpectrumValues(samples [][2]float64, freqSpectrum []float64) {
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
			freqSpectrum[i] = math.Max(freqSpectrum[i]-decay, 0.0)
		} else {
			freqSpectrum[i] = (val + freqSpectrum[i]) / 2.0
		}
	}
}


func main() {
	var file string
    flag.StringVar(&file, "file", "", "mp3 filename")
	flag.Parse()

	openedFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("failed to open file")
	}
	streamer, format, err := mp3.Decode(openedFile)
	if err != nil {
		log.Fatalf("failed to initialize streamer")
	}
	

	sampleRate := beep.SampleRate(numSamples)
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	
	// normalizes song quality 
	resampled := beep.Resample(4, format.SampleRate, numSamples, streamer)
	exit = make(chan bool)
	speaker.Play(beep.Seq(&CustomSteamer{streamer: resampled}, beep.Silence(numSamples), beep.Callback(func() {
		exit <- true
	})))
	


	redraw = make(chan string)
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	top := buildTopBars(freqSpectrum[spectrumOffset:spectrumWidth])
	bottom := buildBottomBars(freqSpectrum[spectrumOffset:spectrumWidth])


	uiEvents := ui.PollEvents()

	for {
		select {
			case <-redraw:
				newData := preventEmpty(top.Data, freqSpectrum[spectrumOffset: spectrumWidth])
				top.Data = newData
				bottom.Data = newData

				ui.Render(top, bottom)

			case e := <-uiEvents:
				switch e.ID {
				case "q", "<C-c>":
					return
				}

			case <-exit:
					return
			}
	}
}

// some hacky shit
func preventEmpty(prevData [] float64, newData [] float64) []float64 {
	notGonnaBeEmpty := make([]float64, spectrumWidth - spectrumOffset)
	copy(notGonnaBeEmpty, newData)

	equal := true
	for i := 1; i < len(newData); i++ {
        if newData[i] != newData[0] {
            equal = false
        }
    }

	// a hack because drawing ui crashes without data points
	// sneaking this one non 0 point in off screen for graceful shutdown
	if (equal) {
		notGonnaBeEmpty[len(notGonnaBeEmpty) - 1] = 1
	}

	return notGonnaBeEmpty
}
