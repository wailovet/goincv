package goincv

import (
	"math/rand"
	"sort"
	"time"
)

type Rander struct {
	Elements          []string
	Weights           []int
	TotalWeight       int
	calibratedWeights bool
	precision         int
	calibrateValue    int
}

func (wrc *Rander) AddElement(element string, weight int) {
	weight *= wrc.calibrateValue
	i := sort.Search(len(wrc.Weights), func(i int) bool { return wrc.Weights[i] > weight })
	wrc.Weights = append(wrc.Weights, 0)
	wrc.Elements = append(wrc.Elements, "")
	copy(wrc.Weights[i+1:], wrc.Weights[i:])
	copy(wrc.Elements[i+1:], wrc.Elements[i:])
	wrc.Weights[i] = weight
	wrc.Elements[i] = element
	wrc.TotalWeight += weight
}

func (wrc *Rander) AddElements(elements map[string]int) {
	for element, weight := range elements {
		wrc.AddElement(element, weight)
	}
}

func (wrc *Rander) GetRandomChoice() (string, int) {
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(wrc.TotalWeight)
	if !wrc.calibratedWeights {
		wrc.calibrateWeights()
	}
	for key, weight := range wrc.Weights {
		value -= weight
		if value <= 0 {
			return wrc.Elements[key], key
		}
	}
	return "", -1
}

func (wrc *Rander) calibrateWeights() {
	if wrc.TotalWeight/wrc.precision < 1 {
		wrc.calibrateValue = wrc.precision / wrc.TotalWeight
		wrc.TotalWeight = 0
		for key := range wrc.Weights {
			wrc.Weights[key] *= wrc.calibrateValue
			wrc.TotalWeight += wrc.Weights[key]
		}
		wrc.calibratedWeights = true
	}
}

func NewRander(arguments ...int) *Rander {
	var precision = 1000
	if len(arguments) > 0 {
		precision = arguments[0]
	}
	return &Rander{
		precision:         precision,
		calibratedWeights: false,
		calibrateValue:    1,
	}
}
