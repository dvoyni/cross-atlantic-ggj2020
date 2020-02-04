package main

import (
	"math/rand"
	"time"
)

type World struct {
	Ship *Ship

	WindForce  float64
	WindForceA float64
	WindForceB float64
	WindTimeA  time.Duration
	WindTimeB  time.Duration

	DistanceLeft float64
	Time         time.Duration

	Events      []*Event
	Performance Performance
	MessageLog  []string

	Finished  bool
	DeltaTime time.Duration
}

type Ship struct {
	Crew         []Sailor
	Speed        float64
	Hull         float64
	Sails        float64
	FloodAmount  float64
	GoodsQuality float64
}

type SailorId int

type Sailor struct {
	Id      SailorId
	Stamina float64
	Work    Work
}

type Work int

const (
	WorkRest Work = iota
	WorkNavigation
	WorkRepairHull
	WorkRepairSail
	WorkPumpOutWater
	WorkShoot
	WorksCount
)

type Performance interface {
	Init(w *World)
	Process(w *World) (finished bool)
}

type Event struct {
	ProbabilityPerDay  float64
	Interval           float64
	CurrentProbability float64
	RequiresMovement   bool
	Performance        Performance
}

func (e *Event) IsHappened(dt time.Duration) bool {
	d := float64(dt) / float64(time.Hour) / 24.
	if e.CurrentProbability < 0 {
		e.CurrentProbability += d
	}
	if e.CurrentProbability >= 0 {
		if rand.Float64() < e.ProbabilityPerDay*d {
			e.CurrentProbability = -e.Interval
			return true
		}
	}
	return false
}

func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func Clamp(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}
