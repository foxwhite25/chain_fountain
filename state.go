package main

import "time"

type State struct {
	Integration   Integration
	Transform     bool
	Integrating   bool
	PLaying       bool
	LastFrameTime float64
	AnimationTime float64
	MaxHeight     float64
}

func (receiver *State) StopIntegrationLoop() {
	receiver.Integrating = false
}

func (receiver *State) SetIntegration(file string) {
	receiver.StopIntegrationLoop()
	receiver.Integration = *NewIntegration(file)
	receiver.SetAnimationTime(0)
}

func (receiver *State) SetAnimationTime(time float64) {
	receiver.AnimationTime = time
}

func (receiver *State) Pause() {
	receiver.PLaying = false
	receiver.LastFrameTime = -1
}

func (receiver *State) StartIntegrationLoop() {
	receiver.Integrating = true
	receiver.PLaying = true
	go func() {
		for receiver.Integrating {
			bufferRoom := 1.0
			lastTimeIndex := len(receiver.Integration.TimeSlices) - 1
			lastTime := float64(lastTimeIndex) * receiver.Integration.Params.TimeStepSize

			if lastTime < (receiver.AnimationTime + bufferRoom) {
				receiver.Integration.Integrate()
				//newTime := lastTime + receiver.Integration.Params.TimeStepSize
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}
