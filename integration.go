package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type Spawning struct {
	x             float64
	y             float64
	i             int
	di            int
	CoilLinkWidth int
}

type Config struct {
	BeadMass        float64 `json:"beadMass"`
	LinkLength      float64 `json:"linkLength"`
	InitialHeight   float64 `json:"initialHeight"`
	TimeStepSize    float64 `json:"timeStepSize"`
	SubSteps        int     `json:"subSteps"`
	LinkStiffness   float64 `json:"linkStiffness"`
	Gravity         float64 `json:"gravity"`
	BeakerWidth     float64 `json:"beakerWidth"`
	BeakerHeight    float64 `json:"beakerHeight"`
	BeakerThickness float64 `json:"beakerThickness"`
	BeakerStiffness float64 `json:"beakerStiffness"`
	TotalBeads      int     `json:"totalBeads"`
	XOffset         float64 `json:"XOffset"`
	YOffset         float64 `json:"YOffset"`
	Zoom            float64 `json:"zoom"`
	PlaySpeed       int     `json:"playSpeed"`
}

func GetBeakerWalls(config Config) []float64 {
	w := config.BeakerWidth
	h := config.BeakerHeight
	t := config.BeakerThickness
	y := config.InitialHeight
	return []float64{
		-(w + t) / 2, y + h - t/2,
		-(w + t) / 2, y - t/2,
		(w + t) / 2, y - t/2,
		(w + t) / 2, y + h - t/2,
	}
}

func CreateConfig(file string) Config {
	config := Config{
		BeadMass:        5,
		LinkLength:      0.01,
		InitialHeight:   1,
		TimeStepSize:    0.001,
		SubSteps:        10,
		LinkStiffness:   1e7,
		Gravity:         9.8,
		BeakerWidth:     0.08,
		BeakerHeight:    0.02,
		BeakerThickness: 0.05,
		BeakerStiffness: 1e7,
		TotalBeads:      300,
		XOffset:         400,
		YOffset:         50,
		Zoom:            600,
		PlaySpeed:       2,
	}
	marshal, err := json.Marshal(config)
	if err != nil {
		return Config{}
	}
	configFIle, err := os.Create(file)
	if err != nil {
		_, _ = fmt.Println(err.Error())
		return Config{}
	}
	var out bytes.Buffer
	err = json.Indent(&out, marshal, "", "    ")
	if err != nil {
		return Config{}
	}
	_, err = out.WriteTo(configFIle)
	if err != nil {
		return Config{}
	}
	return config
}

func LoadConfig(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		_, _ = fmt.Println(err.Error())
		return CreateConfig(file)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return Config{}
	}
	return config
}

type Integration struct {
	Params        Config
	BeakerWalls   []float64
	TimeSlices    [][]float64
	Spawning      Spawning
	PreviousBeads []float64
}

func NewIntegration(file string) *Integration {
	config := LoadConfig(file)
	integration := new(Integration)
	integration.Params = config
	integration.BeakerWalls = GetBeakerWalls(config)
	InitialBeads := make([]float64, 2*config.TotalBeads)
	CoilLinkWidth := int(math.Floor(config.BeakerWidth / config.LinkLength))
	BeakerLinkHeight := int(math.Ceil(config.BeakerHeight / config.LinkLength))
	BeakerLinkThickness := int(math.Ceil(config.BeakerThickness / config.LinkLength))

	j := 0
	y := config.InitialHeight - config.LinkLength
	x := config.BeakerWidth/2 + float64(BeakerLinkThickness)*config.LinkLength
	InitialBeads[j] = x
	InitialBeads[j+1] = y
	for k := 0; k < BeakerLinkHeight+1; k++ {
		y += config.LinkLength
		InitialBeads[j] = x
		InitialBeads[j+1] = y
		j += 2
	}
	for k := 0; k < BeakerLinkThickness; k++ {
		x -= config.LinkLength
		InitialBeads[j] = x
		InitialBeads[j+1] = y
		j += 2
	}
	for k := 0; k < BeakerLinkHeight; k++ {
		y -= config.LinkLength
		InitialBeads[j] = x
		InitialBeads[j+1] = y
		j += 2
	}
	i := 0
	di := -1
	for ; j < 2*config.TotalBeads; j += 2 {
		x += float64(di) * config.LinkLength
		i += di
		if i <= -CoilLinkWidth || i >= 0 {
			di = -di
		}
		InitialBeads[j] = x
		InitialBeads[j+1] = y
	}

	integration.Spawning = Spawning{x, y, i, di, CoilLinkWidth}
	integration.TimeSlices = [][]float64{InitialBeads}
	return integration
}

func (receiver *Integration) Integrate() {
	params := receiver.Params
	s := receiver.Spawning
	curIdx := len(receiver.TimeSlices) - 1
	beads := receiver.TimeSlices[curIdx]
	newBeads := make([]float64, len(beads))
	var prevBeads []float64
	if curIdx > 0 {
		prevBeads = receiver.PreviousBeads
	}
	curAccel := make([]float64, len(beads))
	beadsOnFloor := 0

	dt := params.TimeStepSize / float64(params.SubSteps)

	for t := 0; t < params.SubSteps; t++ {
		receiver.Acceleration(beads, curAccel)
		beadsOnFloor = 0

		for i := 0; i < len(beads); i += 2 {
			if len(prevBeads) != 0 {
				newBeads[i] = 2*beads[i] - prevBeads[i] + curAccel[i]*(math.Pow(dt, 2))
				newBeads[i+1] = 2*beads[i+1] - prevBeads[i+1] + curAccel[i+1]*(math.Pow(dt, 2))
			} else {
				newBeads[i] = beads[i] + curAccel[i]*(math.Pow(dt, 2))/2
				newBeads[i+1] = beads[i+1] + curAccel[i+1]*(math.Pow(dt, 2))/2
			}

			if newBeads[i+1] <= 0 || beads[i+1] <= 0 {
				newBeads[i+1] = 0
				beadsOnFloor++
			}
		}

		for i := 0; beadsOnFloor > 20 && i < 100; i++ {
			if t == 0 {
				if newBeads[i] <= 0 {
					beadsOnFloor--
				}
				newBeads = newBeads[2:]
				beads = beads[2:]
				x := newBeads[len(newBeads)-2]
				y := newBeads[len(newBeads)-1]
				if s.di > 0 && x+params.LinkLength > params.BeakerWidth/2 {
					s.di = -s.di
				}
				if s.di < 0 && x-params.LinkLength < -params.BeakerWidth/2 {
					s.di = -s.di
				}
				x += float64(s.di) * params.LinkLength
				newBeads = append(newBeads, x, y)
				beads = append(beads, x, y)
			}
		}

		if len(prevBeads) == 0 {
			prevBeads = make([]float64, len(beads))
		}
		prevBeads, beads, newBeads = beads, newBeads, prevBeads
	}

	receiver.TimeSlices = append(receiver.TimeSlices, beads)
	receiver.PreviousBeads = prevBeads
}

func (receiver *Integration) Acceleration(beads []float64, accel []float64) {
	params := receiver.Params
	effLinkStiffness := params.LinkStiffness / params.BeadMass
	effBeakerStiffness := params.BeakerStiffness / params.BeadMass
	walls := receiver.BeakerWalls
	var ax, ay, dist, displacementX, displacementY, relStress, wallLength, lineProj, minDist2, wallDist, absAccel float64
	wpv := make([]float64, 2)
	wv := make([]float64, 2)
	pointDistances := make([]float64, len(walls)/2)
	perpProjs := make([]float64, len(walls)-2)
	wallDistances := make([]float64, len(pointDistances)-1)

	for i := 0; i < len(beads); i += 2 {
		accel[i] = 0
		accel[i+1] = 0

		accel[i+1] -= params.Gravity

		if i > 0 {
			displacementX = beads[i] - beads[i-2]
			displacementY = beads[i+1] - beads[i-1]
			dist = math.Sqrt(math.Pow(displacementX, 2) + math.Pow(displacementY, 2))
			relStress = 1 - params.LinkLength/dist
			ax = effLinkStiffness * displacementX * relStress
			ay = effLinkStiffness * displacementY * relStress
			accel[i] -= ax
			accel[i+1] -= ay
			accel[i-2] += ax
			accel[i-1] += ay
		}

		for j := 0; j < len(walls); j += 2 {
			pointDistances[j/2] = math.Pow(beads[i]-walls[j], 2) + math.Pow(beads[i+1]-walls[j+1], 2)
		}

		for j := 0; j < len(walls)-2; j += 2 {
			wallLength = math.Sqrt(math.Pow(walls[j+2]-walls[j], 2) + math.Pow(walls[j+3]-walls[j+1], 2))
			wv[0] = walls[j+2] - walls[j]
			wv[1] = walls[j+3] - walls[j+1]
			wpv[0] = beads[i] - walls[j]
			wpv[1] = beads[i+1] - walls[j+1]
			lineProj = ((wpv[0])*(wv[0]) + (wpv[1])*(wv[1])) / wallLength
			if lineProj >= 0 && lineProj <= wallLength {
				perpProjs[j] = wpv[0] - lineProj*wv[0]/wallLength
				perpProjs[j+1] = wpv[1] - lineProj*wv[1]/wallLength
				wallDistances[j/2] = math.Pow(perpProjs[j], 2) + math.Pow(perpProjs[j+1], 2)
			} else {
				wallDistances[j/2] = math.Inf(1)
			}
		}

		minDist2 = math.Min(MinFromSlice(wallDistances), MinFromSlice(pointDistances))
		if minDist2 < math.Pow(params.BeakerThickness/2, 2) {
			if MinFromSlice(wallDistances) < MinFromSlice(pointDistances) {
				for j := 0; j < len(wallDistances); j++ {
					if wallDistances[j] == minDist2 {
						wallDist = math.Sqrt(wallDistances[j])
						absAccel = effBeakerStiffness * (params.BeakerThickness/2 - wallDist)
						accel[i] += absAccel * perpProjs[2*j] / wallDist
						accel[i+1] += absAccel * perpProjs[2*j+1] / wallDist
						break
					}
				}
			} else {
				for j := 0; j < len(pointDistances); j++ {
					if pointDistances[j] == minDist2 {
						wallDist = math.Sqrt(pointDistances[j])
						absAccel = effBeakerStiffness * (params.BeakerThickness/2 - wallDist)
						accel[i] += absAccel * (beads[i] - walls[2*j]) / wallDist
						accel[i+1] += absAccel * (beads[i+1] - walls[2*j+1]) / wallDist
						break
					}
				}
			}
		}
	}
}
