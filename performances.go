package main

import (
	"math/rand"
	"time"
)

type PerfCoralReef struct {
	DamageMin float64
	DamageMax float64
}

func (p *PerfCoralReef) Init(w *World) {
	w.Log(T[PerfCoralReefDescr])
}

func (p *PerfCoralReef) Process(w *World) (finished bool) {
	w.DamageHull(Lerp(p.DamageMin, p.DamageMax, rand.Float64()))
	return true
}

type PerfBottle struct {
	DamageSailsMin float64
	DamageSailsMax float64
	DamageHullMin  float64
	DamageHullMax  float64

	time time.Duration
	c    int
}

func (p *PerfBottle) Init(w *World) {
	if rand.ExpFloat64() < .3 {
		w.Log(T[FoundBottleEmpty])
		p.c = -1
	} else {
		w.Log(T[FoundBottleDjinn])
		p.c = rand.Intn(5)
		w.Log(T[LangKey(int(DjinnCase0)+p.c)])
		p.time = w.Time
	}
}

func (p *PerfBottle) Process(w *World) (finished bool) {
	switch p.c {
	case 0:
		w.DamageSails(Lerp(p.DamageSailsMin, p.DamageSailsMax, rand.Float64()))
	case 1:
		w.DamageHull(Lerp(p.DamageHullMin, p.DamageHullMax, rand.Float64()))
	case 2:
		w.Ship.Hull = 1
		w.Ship.Sails = 1
		w.Ship.GoodsQuality = 1
	case 3:
		w.Ship.Speed *= 2
		return w.Time-p.time > time.Hour*24
	case 4:
		return w.Time-p.time > time.Hour*24*2
	}
	return true
}

type PerfPirates struct {
	DamagePerShootMin float64
	DamagePerShootMax float64
	Durability        int
	hull              int
}

func (p *PerfPirates) Init(w *World) {
	w.Log(T[PiratesAttacking])
	p.hull = p.Durability
}

func (p *PerfPirates) Process(w *World) (finished bool) {
	if rand.Float64() < .3 {
		w.Log(T[PiratesShot])
		r := rand.Float64()
		if r <= .05 && len(w.Ship.Crew) > 0 {
			w.Log(T[PiratesKilledSailor])
			w.Ship.Crew = w.Ship.Crew[:len(w.Ship.Crew)-1]
		} else if r < .6 {
			w.DamageHull(Lerp(p.DamagePerShootMin, p.DamagePerShootMax, rand.Float64()))
		} else {
			w.DamageSails(Lerp(p.DamagePerShootMin, p.DamagePerShootMax, rand.Float64()))
		}
	}

	if rand.Float64() < w.GetCrewEffectiveness(WorkShoot, 20) {
		w.Log(T[YouShot])
		p.hull -= 1
	}

	if p.hull <= 0 {
		w.Log(T[PiratesDead])
		reward := rand.Float64()
		if reward < .3 {
			w.Log(T[FoundPiratesGoods])
			w.GoodsQualityImproved(.1)
		} else if reward < .6 {
			w.Log(T[FoundPiratesSurvivor])
			w.Ship.Crew = append(w.Ship.Crew, Sailor{
				Id: SailorId(len(w.Ship.Crew)),
			})
		}
		return true
	}
	return false
}

type PerfParty struct {
	DeadProbability float64
}

func (p PerfParty) Init(w *World) {
	w.Log(T[Party])
}

func (p PerfParty) Process(w *World) (finished bool) {
	for i, s := range w.Ship.Crew {
		s.Stamina = 0
		w.Ship.Crew[i] = s
	}

	if rand.Float64() < p.DeadProbability && len(w.Ship.Crew) > 0 {
		w.Log(T[DrunkedDead])
		w.Ship.Crew = w.Ship.Crew[:len(w.Ship.Crew)-1]
	}
	return true
}

type PerfTrader struct {
	QualityDeltaMin float64
	QualityDeltaMax float64
}

func (p *PerfTrader) Init(w *World) {
	w.Log(T[Trader])
}

func (p *PerfTrader) Process(w *World) (finished bool) {
	r := Lerp(p.QualityDeltaMin, p.QualityDeltaMax, rand.Float64())
	if r > 0 {
		w.Log(T[TraderPositive])
		w.GoodsQualityImproved(r)
	} else {
		w.Log(T[TraderNegative])
		w.GoodsQualityReduced(-r)
	}
	return true
}

type PerfBlow struct {
	MinDamage float64
	MaxDamage float64
}

func (p *PerfBlow) Init(w *World) {
	w.Log(T[Blow])
}

func (p *PerfBlow) Process(w *World) (finished bool) {
	w.DamageHull(Lerp(p.MinDamage, p.MaxDamage, rand.Float64()))
	if len(w.Ship.Crew) > 0 {
		w.Ship.Crew = w.Ship.Crew[:len(w.Ship.Crew)-1]
	}
	return true
}
