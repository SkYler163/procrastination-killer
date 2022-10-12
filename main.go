package main

import (
	"embed"
	"log"
	"sync"

	"github.com/pkg/errors"

	"github.com/SkYler163/procrastination-killer/internal/interfacer"
	"github.com/SkYler163/procrastination-killer/internal/model"
	"github.com/SkYler163/procrastination-killer/internal/signaller"
	"github.com/SkYler163/procrastination-killer/internal/timer"
)

//go:embed static
var f embed.FS

func main() {
	var (
		workPeriodMinutes      = 25
		shortRestPeriodMinutes = 5
		longRestPeriodMinutes  = 15
	)

	controlSignalChan := make(chan model.ControlSignals)
	exitChan := make(chan struct{})
	ticksChan, ticksResetChan := make(chan float64), make(chan float64)
	timeLeftChan := make(chan string)
	timerMutex := sync.Mutex{}

	signalFile, err := f.Open("static/cuckoo-clock.mp3")
	if err != nil {
		log.Println(errors.Wrap(err, "failed to open signalFile"))
	}

	s, err := signaller.NewSignaller(signalFile)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to run signaller"))

		return
	}

	pomodoro := timer.NewPomodoro(
		s,
		workPeriodMinutes, shortRestPeriodMinutes, longRestPeriodMinutes,
		controlSignalChan, exitChan,
		ticksChan, ticksResetChan, timeLeftChan, &timerMutex,
	)

	interfaceLocker := sync.Mutex{}

	render, err := interfacer.
		NewInterfacer(pomodoro, ticksChan, timeLeftChan, ticksResetChan, controlSignalChan,
			exitChan, workPeriodMinutes, &interfaceLocker).
		Render(f)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to render interface"))

		return
	}

	render.ShowAndRun()
}
