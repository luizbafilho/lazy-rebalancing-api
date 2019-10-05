package main

import (
	"math"
	"sort"
)

type Asset struct {
	Name                    string  `json:"name,omitempty"`
	Value                   float64 `json:"value,omitempty"`
	ActualAllocation        float64 `json:"actual_allocation,omitempty"`
	TargetAllocationPercent float64 `json:"target_allocation_percent,omitempty"`
	TargetValue             float64 `json:"target_value,omitempty"`
	Deviation               float64 `json:"deviation,omitempty"`
	Delta                   float64 `json:"delta,omitempty"`
}

type Portfolio []Asset

func (p Portfolio) Len() int           { return len(p) }
func (p Portfolio) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Portfolio) Less(i, j int) bool { return p[i].Deviation < p[j].Deviation }

func lazyRebalance(amountToContribute float64, assets Portfolio) Portfolio {
	portfolioTotal := float64(0)
	for _, a := range assets {
		portfolioTotal += a.Value
	}
	total := portfolioTotal + amountToContribute

	for i := range assets {
		a := &assets[i]
		targetValue := total * a.TargetAllocationPercent

		deviation := (a.Value / targetValue) - 1.0

		if portfolioTotal <= 0.0 {
			a.ActualAllocation = 0.0
		} else {
			a.ActualAllocation = a.Value / portfolioTotal
		}

		a.TargetValue = targetValue
		a.Deviation = deviation
	}

	sort.Sort(assets)

	_h := 0.0
	amountLeftToContribute := amountToContribute
	_k := 0.0

	lastKnownIndex := -1

	for i, a := range assets {
		if math.Abs(amountLeftToContribute) <= 0.0 {
			break
		}

		lastKnownIndex = i

		_k = a.Deviation

		targetValue := a.TargetValue
		_h = _h + targetValue

		nextLeastDeviation := 0.0

		if !(i >= len(assets)-1) {
			nextLeastDeviation = assets[i+1].Deviation
		}

		_t := _h * (nextLeastDeviation - _k)

		if math.Abs(_t) <= math.Abs(amountLeftToContribute) {
			amountLeftToContribute = amountLeftToContribute - _t
			_k = nextLeastDeviation
		} else {
			_k = _k + (amountLeftToContribute / _h)

			break
		}
	}

	// indexToStop := lastKnownIndex
	indexToStop := 0
	if lastKnownIndex > -1 {
		indexToStop = lastKnownIndex + 1
	}

	for i := range assets {
		a := &assets[i]
		if i >= indexToStop {
			break
		}

		targeValue := a.TargetValue
		deviation := a.Deviation
		delta := targeValue * (_k - deviation)

		a.Delta = delta
	}

	return assets
}
