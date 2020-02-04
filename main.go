package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"strings"
	"time"
)

func main() {
	display, err := NewDisplay()
	if err != nil {
		return
	}

	defer display.Release()
	uiEvents := ui.PollEvents()

	menu(uiEvents)
	game(display, uiEvents)
}

func menu(uiEvents <-chan ui.Event) {
	screenWidth, screenHeight := ui.TerminalDimensions()
	image := widgets.NewParagraph()
	image.SetRect(0, 0, screenWidth, screenHeight-3)
	image.Title = T[CrossAtlantic]
	image.Text =
		"\n" +
			"     |>\n" +
			"     |\\\n" +
			"  ___|_\\__\n" +
			"  \\      /\n" +
			"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\n"
	ui.Render(image)

	image = widgets.NewParagraph()
	image.SetRect(0, screenHeight-3, screenWidth, screenHeight)
	image.Text = T[PressAKey]
	ui.Render(image)
	for {
		e := <-uiEvents
		switch e.ID {
		case "r":
			T = langRu
			return
		case "e":
			T = langEn
			return
		}
	}
}

func game(display *Display, uiEvents <-chan ui.Event) {
	world := NewWorld(time.Now().UnixNano())

	for {
		select {
		case e := <-uiEvents:
			work := WorksCount
			assign := false
			mostStamined := false

			switch e.ID {
			case "<C-c>":
				return
			case "1":
				work = WorkNavigation
				assign = true
			case "2":
				work = WorkRepairHull
				assign = true
			case "3":
				work = WorkRepairSail
				assign = true
			case "4":
				work = WorkPumpOutWater
				assign = true
			case "5":
				work = WorkShoot
				assign = true
			case "!":
				work = WorkNavigation
				//mostStamined = true
			case "@":
				work = WorkRepairHull
				//mostStamined = true
			case "#":
				work = WorkRepairSail
				//mostStamined = true
			case "$":
				work = WorkPumpOutWater
				//mostStamined = true
			case "%":
				work = WorkShoot
			//mostStamined = true
			case "<F1>":
				lines := strings.Split(T[ControlsHelp], "\n")
				for _, l := range lines {
					world.Log(l)
				}
			case "<F10>":
				world = NewWorld(time.Now().UnixNano())

				/*case "<C-1>":
					work = WorkNavigation
				case "<C-2>":
					work = WorkRepairHull
				case "<C-3>":
					work = WorkRepairSail
				case "<C-4>":
					work = WorkPumpOutWater
				case "<C-5>":
					work = WorkShoot*/
			}

			if work != WorksCount && !world.Finished {
				if assign {
					sailorId := world.FindSailor(WorkRest, true)
					if sailorId != -1 {
						world.AssignSailor(sailorId, work)
					}
				} else {
					sailorId := world.FindSailor(work, mostStamined)
					if sailorId != -1 {
						world.UnassignSailor(sailorId)
					}
				}
			}

		case <-time.After(1 * time.Second):
			world.Update()
		}

		display.Draw(world)
	}
}
