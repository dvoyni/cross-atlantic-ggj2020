package main

type LangKey int

const (
	GeneralInformationTitle LangKey = iota
	GeneralInformationText
	CrewAndJobsTitle
	JobIdle
	JobNavigating
	JobRepairingHull
	JobRepairingSail
	JobPumpingOutWater
	JobShootingCannons
	OverviewTitle
	DispatchMessage
	ShipIsStandingOverview
	ShipGoesByOarsOverview
	ShipGoesBySailOverview
	StormOutsideOverview
	NavigationTeamOptimalSize
	PerfCoralReefDescr
	EventsTitle
	HullIsDamaged
	SailsAreDamaged
	LogMessage
	JourneyStarted
	ShipIsFullyBroken
	ShipIsFullyFlooded
	ReachedDestination
	JourneyEnded
	HullIsBrokenOverview
	SailIsBrokenOverview
	ShipIsFloodedOverview
	ControlsHelp
	FoundBottleEmpty
	FoundBottleDjinn
	DjinnCase0
	DjinnCase1
	DjinnCase2
	DjinnCase3
	DjinnCase4
	PiratesAttacking
	PiratesShot
	PiratesDead
	PiratesKilledSailor
	FoundPiratesGoods
	FoundPiratesSurvivor
	YouShot
	Party
	DrunkedDead
	GoodsQualityImproved
	GoodsQualityReduced
	Trader
	TraderPositive
	TraderNegative
	Blow
	CrossAtlantic
	PressAKey
)

type LangMap map[LangKey]string

var T = langEn
