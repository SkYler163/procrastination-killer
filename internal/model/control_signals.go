package model

// ControlSignals control signal type.
type ControlSignals int8

// ControlsSignals constants.
const (
	ControlSignalPlay  ControlSignals = 0
	ControlSignalPause ControlSignals = 1
	ControlSignalStop  ControlSignals = 2
)
