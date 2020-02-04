package main

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"strconv"
	"strings"
)

type Display struct {
	GeneralInfo *widgets.Paragraph
	SailorsList *widgets.List
	Overview    *widgets.Paragraph
	Log         *widgets.List
}

func NewDisplay() (*Display, error) {
	if err := ui.Init(); err != nil {
		return nil, fmt.Errorf("cannot initialize ui: %w", err)
	}

	return &Display{
		GeneralInfo: widgets.NewParagraph(),
		SailorsList: widgets.NewList(),
		Overview:    widgets.NewParagraph(),
		Log:         widgets.NewList(),
	}, nil
}

func (display *Display) Draw(world *World) {
	screenWidth, screenHeight := ui.TerminalDimensions()
	display.GeneralInfo.SetRect(0, 0, screenWidth, 4)
	display.GeneralInfo.PaddingLeft = 2
	display.GeneralInfo.PaddingRight = 2
	display.GeneralInfo.WrapText = false
	display.GeneralInfo.Title = T[GeneralInformationTitle]
	display.GeneralInfo.Text = fmt.Sprintf(T[GeneralInformationText],
		int(world.Time.Hours())/24+1,
		int(world.Time.Hours())%24,
		int(world.Time.Minutes())%60,
		world.WindForce,
		world.Ship.Speed,
		world.DistanceLeft,
		world.Ship.Hull*100,
		c(world.Ship.Hull),
		world.Ship.Sails*100,
		c(world.Ship.Sails),
		world.Ship.FloodAmount*100,
		ci(world.Ship.FloodAmount),
		world.Ship.GoodsQuality*100,
		c(world.Ship.GoodsQuality))
	ui.Render(display.GeneralInfo)

	display.SailorsList.SetRect(0, display.GeneralInfo.Rectangle.Max.Y,
		screenWidth, display.GeneralInfo.Rectangle.Max.Y+2+int(WorksCount))
	display.SailorsList.PaddingLeft = 2
	display.SailorsList.PaddingRight = 2
	display.SailorsList.Title = T[CrewAndJobsTitle]
	display.SailorsList.Rows = stringifyWorks(world)
	display.SailorsList.Border = true
	ui.Render(display.SailorsList)

	display.Overview.SetRect(0, display.SailorsList.Rectangle.Max.Y,
		screenWidth, display.SailorsList.Rectangle.Max.Y+2+5)
	display.Overview.PaddingLeft = 2
	display.Overview.PaddingRight = 2
	display.Overview.Title = T[OverviewTitle]
	display.Overview.WrapText = true
	if world.Performance == nil {
		overview := ""
		if world.Ship.Hull < 1 {
			overview += T[HullIsBrokenOverview]
		}

		if world.Ship.Sails < 1 {
			overview += T[SailIsBrokenOverview]
		}

		if world.Ship.FloodAmount > 0 {
			overview += T[ShipIsFloodedOverview]
		}

		if world.DistanceLeft == 3600 {
			overview += T[DispatchMessage]
		} else {
			if world.Ship.Speed > 0 {
				windForce := int(world.WindForce)
				if windForce < 3 {
					overview += T[ShipGoesByOarsOverview]
				} else if windForce < 8 {
					overview += T[ShipGoesBySailOverview]
				} else {
					overview += T[StormOutsideOverview]
				}
				overview += fmt.Sprintf(T[NavigationTeamOptimalSize], EffectiveNavigatorsPerWind[int(world.WindForce)])
			} else {
				overview += T[ShipIsStandingOverview]
			}
		}

		display.Overview.Text = overview
	}
	ui.Render(display.Overview)

	display.Log.SetRect(0, display.Overview.Rectangle.Max.Y, screenWidth, screenHeight)
	display.Log.PaddingLeft = 2
	display.Log.PaddingRight = 2
	display.Log.Title = T[EventsTitle]
	display.Log.WrapText = false
	start := len(world.MessageLog) - 20
	if start < 0 {
		start = 0
	}
	display.Log.Rows = world.MessageLog[start:]
	display.Log.ScrollTop()
	display.Log.ScrollBottom()
	ui.Render(display.Log)
}

func (display *Display) Release() {
	ui.Close()
}

var workNames = []LangKey{
	JobIdle,
	JobNavigating,
	JobRepairingHull,
	JobRepairingSail,
	JobPumpingOutWater,
	JobShootingCannons,
}

func stringifyWorks(world *World) []string {
	works := make([]string, WorksCount)

	for i := 0; i < int(WorksCount); i++ {
		works[i] = stringifyWork(world, Work(i))
	}

	return works
}

func stringifyWork(world *World, work Work) string {
	sb := &strings.Builder{}

	if work == WorkRest {
		sb.WriteString("    ")
	} else {
		sb.WriteString("[")
		sb.WriteString(strconv.Itoa(int(work)))
		sb.WriteString("] ")
	}
	sb.WriteString(T[workNames[work]])

	cnt := 0
	for _, sailor := range world.Ship.Crew {
		if sailor.Work == work {
			cnt++
		}
	}
	sb.WriteString(fmt.Sprintf(" (%2d):    ", cnt))

	for _, sailor := range world.Ship.Crew {
		if sailor.Work == work {
			cnt++
			sb.WriteString(stringifySailor(sailor))
			sb.WriteString(" ")
		}
	}

	return sb.String()
}

//var faces = []string{"😴", "😞", "😓", "😫", "😠", "😕", "😬", "😐", "🙂", "😀", "😎"}

func stringifySailor(sailor Sailor) string {
	//return fmt.Sprintf("[%3.0f%%](%v)%s", sailor.Stamina*100, c(sailor.Stamina), faces[int(sailor.Stamina*10)])
	return fmt.Sprintf("([%3.0f%%](%v))", sailor.Stamina*100, c(sailor.Stamina))
}

func c(v float64) string {
	v = Clamp(v)
	if v >= 1 {
		return "fg:white"
	} else if v > .75 {
		return "fg:green"
	} else if v > .5 {
		return "fg:yellow"
	} else if v > .25 {
		return "fg:red"
	} else {
		return "bg:red,fg:black"
	}
}

func ci(v float64) string {
	return c(1 - v)
}
