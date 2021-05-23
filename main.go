package main

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	"github.com/SkYler163/procrastination-killer/internal/interfacer"
	"github.com/SkYler163/procrastination-killer/internal/model"
	"github.com/SkYler163/procrastination-killer/internal/signaller"
	"github.com/SkYler163/procrastination-killer/internal/timer"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(errors.Wrap(err, "error loading .env file"))

		return
	}

	var (
		workPeriodMinutes      = 25
		shortRestPeriodMinutes = 5
		longRestPeriodMinutes  = 15
	)

	if wpm := os.Getenv("WORK_PERIOD_MINUTES"); wpm != "" {
		v, err := strconv.ParseInt(wpm, 10, 32)
		if err != nil {
			log.Println(errors.Wrap(err, "failed parse work period minutes env"))

			return
		}

		workPeriodMinutes = int(v)
	}

	if srpm := os.Getenv("SHORT_REST_PERIOD_MINUTES"); srpm != "" {
		v, err := strconv.ParseInt(srpm, 10, 32)
		if err != nil {
			log.Println(errors.Wrap(err, "failed parse short rest period minutes env"))

			return
		}

		shortRestPeriodMinutes = int(v)
	}

	if lrpm := os.Getenv("LONG_REST_PERIOD_MINUTES"); lrpm != "" {
		v, err := strconv.ParseInt(lrpm, 10, 32)
		if err != nil {
			log.Println(errors.Wrap(err, "failed parse long rest period minutes env"))

			return
		}

		longRestPeriodMinutes = int(v)
	}

	controlSignalChan := make(chan model.ControlSignals)
	exitChan := make(chan struct{})
	ticksChan, ticksResetChan := make(chan float64), make(chan float64)
	timeLeftChan := make(chan string)
	timerMutex := sync.Mutex{}

	s, err := signaller.NewSignaller("static/cuckoo-clock.mp3")
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
		Render()
	if err != nil {
		log.Println(errors.Wrap(err, "failed to render interface"))

		return
	}

	render.ShowAndRun()
}
