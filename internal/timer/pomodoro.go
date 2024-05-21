package timer

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/SkYler163/procrastination-killer/internal/model"
	"github.com/SkYler163/procrastination-killer/internal/signaller"
)

// Pomodoro pomodoro timer struct.
type Pomodoro struct {
	signaller          *signaller.Signaller
	controlSignalsChan chan model.ControlSignals
	ticksChan          chan float64
	ticksReset         chan float64
	timeLeftChan       chan string
	exitChan           chan struct{}
	mu                 *sync.Mutex
	workPeriodMinutes  int
	shortRestMinutes   int
	longRestMinutes    int
}

// NewPomodoro creates an instance of pomodoro timer.
func NewPomodoro(
	signaller *signaller.Signaller,
	workPeriodMinutes, shortRestMinutes, longRestMinutes int,
	controlSignalsChan chan model.ControlSignals,
	exitChan chan struct{},
	ticksChan, ticksReset chan float64,
	timeLeftChan chan string,
	mu *sync.Mutex,
) *Pomodoro {
	return &Pomodoro{
		signaller:          signaller,
		workPeriodMinutes:  workPeriodMinutes,
		shortRestMinutes:   shortRestMinutes,
		longRestMinutes:    longRestMinutes,
		controlSignalsChan: controlSignalsChan,
		ticksChan:          ticksChan,
		ticksReset:         ticksReset,
		timeLeftChan:       timeLeftChan,
		exitChan:           exitChan,
		mu:                 mu,
	}
}

const (
	secondsInMinute  = 60
	roundsToLongRest = 4
)

// Run runs pomodoro timer.
func (p *Pomodoro) Run() {
	ticker := time.NewTicker(time.Second)
	roundEnd := make(chan struct{})
	ticks := p.workPeriodMinutes * secondsInMinute
	ticksNumber := make(chan int)
	ticksPassed := float64(0)
	p.ticksReset <- float64(ticks)

	go p.tickerController(roundEnd, ticksNumber, p.exitChan, ticker)

	go func() {
		for {
			select {
			case <-p.exitChan:
				return
			case tn := <-ticksNumber:
				p.mu.Lock()
				ticks = tn
				p.mu.Unlock()
			case <-ticker.C:
				p.timeLeftChan <- fmt.Sprintf("%02d:%02d", ticks/secondsInMinute, ticks%secondsInMinute)
				p.mu.Lock()

				ticks--
				ticksPassed++
				p.ticksChan <- ticksPassed

				if ticks == 0 {
					roundEnd <- struct{}{}

					ticksPassed = 0
				}

				p.mu.Unlock()
			}
		}
	}()
}

func (p *Pomodoro) tickerController(
	roundEnd chan struct{},
	ticksNumber chan int,
	exitChan chan struct{},
	ticker *time.Ticker,
) {
	signalChan := make(chan os.Signal, 1)
	isWorkPeriod := true
	roundNumber := 1

	var newTicksNumber int

	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case cs := <-p.controlSignalsChan:
			switch cs {
			case model.ControlSignalPlay:
				ticker.Reset(time.Second)
			case model.ControlSignalPause:
				ticker.Stop()
			case model.ControlSignalStop:
				ticker.Stop()
				exitChan <- struct{}{}
				exitChan <- struct{}{}

				return
			}
		case s := <-signalChan:
			switch s {
			case syscall.SIGINT, syscall.SIGTERM:
				exitChan <- struct{}{}
				exitChan <- struct{}{}

				os.Exit(1)
			}
		case <-roundEnd:
			p.signaller.Signal()

			isWorkPeriod = !isWorkPeriod

			switch {
			case isWorkPeriod:
				roundNumber++

				newTicksNumber = p.workPeriodMinutes * secondsInMinute
			case !isWorkPeriod && (roundNumber%roundsToLongRest) == 0:
				newTicksNumber = p.longRestMinutes * secondsInMinute
			default:
				newTicksNumber = p.shortRestMinutes * secondsInMinute
			}

			ticksNumber <- newTicksNumber
			p.ticksReset <- float64(newTicksNumber)
		}
	}
}
