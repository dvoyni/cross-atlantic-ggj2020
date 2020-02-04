package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

func NewWorld(seed int64) *World {
	rand.Seed(seed)
	world := &World{
		Ship: &Ship{
			Crew:         []Sailor{},
			Hull:         1,
			Sails:        1,
			GoodsQuality: 1,
		},
		DistanceLeft: 3600,
		Time:         9 * time.Hour,
	}

	world.Ship.Crew = make([]Sailor, 20)
	for i := 0; i < len(world.Ship.Crew); i++ {
		world.Ship.Crew[i] = Sailor{
			Id:      SailorId(i),
			Stamina: 1,
			Work:    WorkRest,
		}
	}

	world.Events = []*Event{
		{
			ProbabilityPerDay:  1 / 5.,
			Interval:           1,
			CurrentProbability: -1,
			RequiresMovement:   true,
			Performance: &PerfCoralReef{
				DamageMin: 0.05,
				DamageMax: 0.25,
			},
		},
		{
			ProbabilityPerDay: 1 / 2.,
			Interval:          1,
			RequiresMovement:  true,
			Performance: &PerfBottle{
				DamageSailsMin: 0.25,
				DamageSailsMax: 0.5,
				DamageHullMin:  0.1,
				DamageHullMax:  0.25,
			},
		},
		{
			ProbabilityPerDay: 1 / 2.,
			Interval:          2,
			RequiresMovement:  false,
			Performance: &PerfPirates{
				DamagePerShootMin: 0,
				DamagePerShootMax: .1,
				Durability:        5,
			},
		},
		{
			ProbabilityPerDay: 1 / 3.,
			Interval:          3,
			Performance: &PerfParty{
				DeadProbability: .2,
			},
		},
		{
			ProbabilityPerDay: 1 / 2.,
			Interval:          2,
			Performance: &PerfTrader{
				QualityDeltaMin: -.1,
				QualityDeltaMax: .1,
			},
		},
		{
			ProbabilityPerDay:  1 / 5.,
			Interval:           5,
			CurrentProbability: 0,
			Performance: &PerfBlow{
				MinDamage: 0.3,
				MaxDamage: 0.4,
			},
		},
	}

	world.Log(T[JourneyStarted])

	return world
}

func (world *World) Update() {
	if world.Finished {
		return
	}

	dt := time.Hour // 2
	world.DeltaTime = dt
	perHour := float64(dt) / float64(time.Hour)
	perDay := perHour / 24.

	world.Time += dt

	if world.WindTimeB < world.Time {
		world.WindTimeA = world.Time
		world.WindTimeB = world.Time + time.Duration(rand.Intn(4)+2)*time.Hour
		world.WindForceA = world.WindForceB
		world.WindForceB = rand.Float64() * (MaxWindForce + 1)
	}

	windT := float64(world.Time-world.WindTimeA) / float64(world.WindTimeB-world.WindTimeA)
	world.WindForce = Lerp(world.WindForceA, world.WindForceB, windT)

	for i, s := range world.Ship.Crew {
		s.Stamina = math.Min(1, math.Max(0, s.Stamina+StaminaDeltaPerHour[s.Work]*perHour))
		world.Ship.Crew[i] = s
		if s.Stamina == 0 {
			world.UnassignSailor(s.Id)
		}
	}

	world.Ship.Hull = Clamp(world.Ship.Hull + RepairingHullPerHour*world.GetCrewEffectiveness(WorkRepairHull, 10))
	if world.Ship.Hull == 0 {
		world.Log(T[ShipIsFullyBroken])
		world.Log(T[JourneyEnded])
		world.Finished = true
	}

	world.Ship.Sails = Clamp(world.Ship.Sails + RepairingSailsPerHour*world.GetCrewEffectiveness(WorkRepairSail, 10))

	world.Ship.FloodAmount = Clamp(world.Ship.FloodAmount + Lerp(MaxFloodingPerHour, 0, world.Ship.Hull)*perHour)
	world.Ship.FloodAmount = Clamp(world.Ship.FloodAmount - PumpingOutPerHour*world.GetCrewEffectiveness(WorkPumpOutWater, 10)*perHour)
	if world.Ship.FloodAmount == 1 {
		world.Log(T[ShipIsFullyFlooded])
		world.Log(T[JourneyEnded])
		world.Finished = true
	}

	world.Ship.GoodsQuality = Clamp(world.Ship.GoodsQuality - GoodQualityReducePerDay*perDay)
	world.Ship.GoodsQuality = Clamp(world.Ship.GoodsQuality - Lerp(0, MaxGoodQualityFloodReducePerHour, world.Ship.FloodAmount)*perHour)

	windForce := int(world.WindForce)
	world.Ship.Speed = MaxSpeedPerWind[windForce] * world.GetCrewEffectiveness(WorkNavigation, EffectiveNavigatorsPerWind[windForce]) * world.Ship.Sails

	if !world.Finished {
		if world.Performance == nil {
			for _, e := range world.Events {
				passed := true
				if e.RequiresMovement {
					passed = passed && world.Ship.Speed > 0
				}
				if passed && e.IsHappened(dt) {
					world.Performance = e.Performance
					world.Performance.Init(world)
					break
				}
			}
		}

		if world.Performance != nil {
			if world.Performance.Process(world) {
				world.Performance = nil
			}
		}

		world.DistanceLeft -= world.Ship.Speed * perHour
		if world.DistanceLeft <= 0 {
			world.DistanceLeft = 0
			score := int(1000*world.Ship.GoodsQuality + 100*world.Ship.Hull + 100*world.Ship.Sails)
			world.Log(fmt.Sprintf(T[ReachedDestination], score))
			world.Log(T[JourneyEnded])
			world.Finished = true
		}
	} else {
		world.Ship.Speed = 0
	}
}

func (world *World) UnassignSailor(id SailorId) {
	for i, s := range world.Ship.Crew {
		if s.Id == id {
			s.Work = WorkRest
			world.Ship.Crew[i] = s
			break
		}
	}
}

func (world *World) AssignSailor(id SailorId, work Work) {
	for i, s := range world.Ship.Crew {
		if s.Id == id && s.Stamina > 0 {
			s.Work = work
			world.Ship.Crew[i] = s
			break
		}
	}
}

func (world *World) FindSailor(work Work, mostStamined bool) SailorId {
	minId := SailorId(-1)
	maxId := SailorId(1)
	minStamina := 2.
	maxStamina := -1.

	for _, s := range world.Ship.Crew {
		if s.Work == work {
			if s.Stamina > maxStamina {
				maxId = s.Id
				maxStamina = s.Stamina
			}
			if s.Stamina < minStamina {
				minId = s.Id
				minStamina = s.Stamina
			}
		}
	}

	if mostStamined {
		return maxId
	} else {
		return minId
	}
}

func (world *World) GetCrewEffectiveness(work Work, requiredSailors int) float64 {
	var workers []float64
	for _, sailor := range world.Ship.Crew {
		if sailor.Work == work {
			workers = append(workers, math.Min(1, .5+sailor.Stamina))
		}
	}

	sort.Slice(workers, func(i, j int) bool {
		return workers[i] > workers[j]
	})

	e := 0.

	for i := 0; i < requiredSailors && i < len(workers); i++ {
		e += workers[i]
	}
	e /= float64(requiredSailors)

	for i := requiredSailors; i < len(workers); i++ {
		e += workers[i] * .1 / float64(len(workers)-requiredSailors)
	}

	return e
}

func (world *World) Log(s string) {
	world.MessageLog = append(world.MessageLog,
		fmt.Sprintf(T[LogMessage],
			int(world.Time.Hours())/24+1,
			int(world.Time.Hours())%24,
			int(world.Time.Minutes())%60,
			s))
}

func (world *World) DamageHull(amount float64) {
	world.Ship.Hull = Clamp(world.Ship.Hull - amount)
	world.Log(fmt.Sprintf(T[HullIsDamaged], amount*100))
}

func (world *World) DamageSails(amount float64) {
	world.Ship.Sails = Clamp(world.Ship.Sails - amount)
	world.Log(fmt.Sprintf(T[SailsAreDamaged], amount*100))
}

func (world *World) GoodsQualityImproved(amount float64) {
	world.Ship.GoodsQuality = Clamp(world.Ship.GoodsQuality + amount)
	world.Log(fmt.Sprintf(T[GoodsQualityImproved], amount*100))
}

func (world *World) GoodsQualityReduced(amount float64) {
	world.Ship.GoodsQuality = Clamp(world.Ship.GoodsQuality - amount)
	world.Log(fmt.Sprintf(T[GoodsQualityReduced], amount*100))
}
