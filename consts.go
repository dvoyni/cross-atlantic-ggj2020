package main

var StaminaDeltaPerHour = []float64{
	2. / 13.,
	-1. / 20.,
	-1. / 10.,
	-1. / 16.,
	-1. / 8.,
	-1. / 18.,
}

const MaxWindForce = 9.

var MaxSpeedPerWind = []float64{
	10, //0
	11, //1
	12, //2
	15, //3
	17, //4
	19, //5
	20, //6
	15, //7
	10, //8
	8,  //9
}

var EffectiveNavigatorsPerWind = []int{
	20, //0
	20, //1
	20, //2
	5,  //3
	5,  //4
	6,  //5
	6,  //6
	7,  //7
	20, //8
	20, //9
}

const GoodQualityReducePerDay = .02

const MaxGoodQualityFloodReducePerHour = .01
const MaxFloodingPerHour = .07
const PumpingOutPerHour = .1
const RepairingHullPerHour = .07
const RepairingSailsPerHour = .05
