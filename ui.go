package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/h8gi/canvas"
	"golang.org/x/image/colornames"
	"math"
)

type UI struct {
	Canvas    *canvas.Canvas
	FrameTime float64
	State     *State
	ctx       *canvas.Context
}

func NewUI() *UI {
	c := canvas.NewCanvas(&canvas.CanvasConfig{
		Width:     1000,
		Height:    1000,
		FrameRate: 60,
		Title:     "Beads Simulation",
	})

	c.Setup(func(ctx *canvas.Context) {
		ctx.SetColor(colornames.White)
		ctx.Clear()
	})

	state := new(State)
	state.SetIntegration("config.json")
	state.StartIntegrationLoop()

	ui := new(UI)
	ui.FrameTime = 0.0
	ui.State = state
	n := 0.0

	c.Draw(func(ctx *canvas.Context) {
		n += float64(state.Integration.Params.PlaySpeed)
		ui.ctx = ctx
		ctx.Push()
		timeStepSize := state.Integration.Params.TimeStepSize
		lastTimeIndex := len(state.Integration.TimeSlices) - 1
		lastTime := float64(lastTimeIndex) * timeStepSize

		if state.PLaying {
			playRate := 1.0
			if state.LastFrameTime != 0 {
				timeDiff := (n - state.LastFrameTime) / 1000
				newTime := state.AnimationTime + playRate*timeDiff
				if !(newTime > lastTime) {
					ui.setAnimationTime(newTime)
				} else if !(state.AnimationTime > lastTime) {
					ui.setAnimationTime(lastTime)
				}
			}
			state.LastFrameTime = n
		}

		timeIndex := math.Floor(state.AnimationTime / timeStepSize)
		if int(timeIndex) > lastTimeIndex {
			ui.drawIntegration(lastTimeIndex)
		} else {
			ui.drawIntegration(int(timeIndex))
		}
	})

	ui.Canvas = c

	return ui
}

func (receiver *UI) setAnimationTime(time float64) {
	receiver.State.AnimationTime = time
}

func (receiver *UI) drawIntegration(timeIndex int) {
	params := receiver.State.Integration.Params
	beads := receiver.State.Integration.TimeSlices[timeIndex]
	walls := receiver.State.Integration.BeakerWalls
	ctx := receiver.ctx
	zoom := receiver.State.Integration.Params.Zoom
	xo := receiver.State.Integration.Params.XOffset
	yo := receiver.State.Integration.Params.YOffset

	receiver.HandleMouse()

	ctx.SetColor(colornames.White)
	ctx.Clear()

	ctx.SetColor(colornames.Black)
	ctx.DrawRectangle(0, 0, 1000, 0*zoom+yo)
	ctx.Fill()

	ctx.SetLineWidth(params.BeakerThickness * zoom)
	ctx.SetLineJoin(gg.LineJoinRound)
	ctx.SetLineCap(gg.LineCapRound)
	ctx.MoveTo(walls[0]*zoom+xo, walls[1]*zoom+yo)
	for i := 0; i < len(walls); i += 2 {
		if i > 0 {
			ctx.LineTo(walls[i]*zoom+xo, walls[i+1]*zoom+yo)
		}
	}
	ctx.Stroke()

	ctx.MoveTo(beads[0]*zoom+xo, beads[1]*zoom+yo)
	ctx.SetLineWidth(0.005 * zoom)
	ctx.SetColor(colornames.Blue)
	for i := 0; i < len(beads); i += 2 {
		if i > 0 {
			ctx.LineTo(beads[i]*zoom+xo, beads[i+1]*zoom+yo)
		}
	}
	ctx.Stroke()

	ctx.SetColor(colornames.Red)
	for i := 0; i < len(beads); i += 2 {
		ctx.DrawArc(beads[i]*zoom+xo, beads[i+1]*zoom+yo, 0.005*zoom, 0, math.Pi*2)
		ctx.Fill()
	}

	maxHeight := 0.0
	for i := 0; i < len(beads); i += 4 {
		if beads[i+1] > maxHeight {
			maxHeight = beads[i+1]
		}
	}

	ctx.SetLineWidth(0.004 * zoom)
	if maxHeight > receiver.State.MaxHeight {
		receiver.State.MaxHeight = maxHeight
	}
	ctx.DrawLine(0, receiver.State.MaxHeight*zoom+yo, 1000*zoom+xo, receiver.State.MaxHeight*zoom+yo)
	ctx.RotateAbout(math.Pi, 500, 500)
	ctx.ScaleAbout(-1, 1, 500, 500)
	ctx.DrawStringAnchored(
		fmt.Sprintf(
			"Max Height: %.6fm", receiver.State.MaxHeight-receiver.State.Integration.Params.InitialHeight,
		),
		100,
		1000-(receiver.State.MaxHeight*zoom+yo)-10,
		0.5,
		0.5,
	)
	ctx.Stroke()
	ctx.ScaleAbout(-1, 1, 500, 500)
	ctx.RotateAbout(math.Pi, 500, 500)
}

func (receiver *UI) HandleMouse() {
	ctx := receiver.ctx
	if !ctx.IsMouseDragged {
		return
	}
	dx := ctx.Mouse.X - ctx.PMouse.X
	dy := ctx.Mouse.Y - ctx.PMouse.Y
	receiver.State.Integration.Params.XOffset += dx
	receiver.State.Integration.Params.YOffset += dy
}
