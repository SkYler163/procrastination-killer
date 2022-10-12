package interfacer

import (
	"sync"

	"github.com/SkYler163/procrastination-killer/internal/model"
	"github.com/SkYler163/procrastination-killer/internal/timer"
)

// Interfacer an interfacer struct.
type Interfacer struct {
	timer              *timer.Pomodoro
	ticksChan          chan float64
	ticksReset         chan float64
	timeLeftChan       chan string
	controlSignalsChan chan model.ControlSignals
	exitChan           chan struct{}
	mu                 *sync.Mutex
	workTimeMin        int
}

// NewInterfacer creates instance of interfacer.
func NewInterfacer(
	timer *timer.Pomodoro,
	ticksChan chan float64,
	timeLeftChan chan string,
	ticksReset chan float64,
	controlSignalsChan chan model.ControlSignals,
	exitChan chan struct{},
	workTimeMin int,
	mu *sync.Mutex,
) *Interfacer {
	return &Interfacer{
		timer:              timer,
		ticksChan:          ticksChan,
		timeLeftChan:       timeLeftChan,
		ticksReset:         ticksReset,
		controlSignalsChan: controlSignalsChan,
		exitChan:           exitChan,
		workTimeMin:        workTimeMin,
		mu:                 mu,
	}
}
