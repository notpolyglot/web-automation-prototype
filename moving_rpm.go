package main

import (
	"math"
	"time"
)

// @author Robin Verlangen
// Moving average implementation for Go

type RollingRPM struct {
	Window           int
	values           []float64
	valPos           int
	slotsFilled      bool
	lastInsert       float64
	timeOfLastInsert time.Time
}

func (ma *RollingRPM) Avg() float64 {
	var sum = float64(0)
	var c = ma.Window - 1

	// Are all slots filled? If not, ignore unused
	if !ma.slotsFilled {
		c = ma.valPos - 1
		if c < 0 {
			// Empty register
			return 0
		}
	}

	// Sum values
	var ic = 0
	for i := 0; i <= c; i++ {
		sum += ma.values[i]
		ic++
	}

	// Finalize average and return
	avg := sum / float64(ic)
	return avg
}

func (ma *RollingRPM) Max() float64 {
	var max = float64(0)
	var c = ma.Window - 1

	// Are all slots filled? If not, ignore unused
	if !ma.slotsFilled {
		c = ma.valPos - 1
		if c < 0 {
			// Empty register
			return 0
		}
	}

	// Calculate max value
	for i := 0; i <= c; i++ {
		if ma.values[i] > max || i == 0 && ma.values[0] != 0 {
			max = ma.values[i]
		}
	}

	return max
}

func (ma *RollingRPM) Min() float64 {
	var min = float64(0)
	var c = ma.Window - 1

	// Are all slots filled? If not, ignore unused
	if !ma.slotsFilled {
		c = ma.valPos - 1
		if c < 0 {
			// Empty register
			return 0
		}
	}

	// Calculate min value
	for i := 0; i <= c; i++ {
		if ma.values[i] < min || i == 0 && ma.values[0] != 0 {
			min = ma.values[i]
		}
	}

	return min
}

func (ma *RollingRPM) Values() []float64 {
	// return all values
	return ma.values
}

func (ma *RollingRPM) SlotsFilled() bool {
	// return all values
	return ma.slotsFilled
}

func (ma *RollingRPM) Add(val float64) {
	if math.IsNaN(val) {
		panic("Value to add is NaN.")
	}

	//find how many responses since last check
	timeSinceLast := time.Now().Sub(ma.timeOfLastInsert)
	amountCheckedSinceLast := val - ma.lastInsert

	//calculate and apply the factor to the value
	fac := int64(time.Minute / timeSinceLast)
	rpm := amountCheckedSinceLast * float64(fac)
	rounded := math.Round(rpm*100) / 100

	ma.timeOfLastInsert = time.Now()
	ma.lastInsert = val

	// Put into values array
	ma.values[ma.valPos] = rounded

	// Increment value position
	ma.valPos = (ma.valPos + 1) % ma.Window

	// Did we just go back to 0, effectively meaning we filled all registers?
	if !ma.slotsFilled && ma.valPos == 0 {
		ma.slotsFilled = true
	}
}

func NewMovingRPM(window int) *RollingRPM {
	return &RollingRPM{
		Window:      window,
		values:      make([]float64, window),
		valPos:      0,
		slotsFilled: false,
	}
}

func (p *WorkerPool) sampleRPM() {
	total := 0
	for k, v := range p.counters {
		_, isDefaultCounter := defaultCounters[k]
		if k == "error" {
			continue
		}
		if isDefaultCounter {
			total += v
		}
	}
	p.rpm.Add(float64(total))

}
