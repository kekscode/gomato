package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/caseymrm/menuet"
)

// TODO: Add App Icon

const (
	pomodoroIco   string = "🍅"
	shortBreakIco string = "⏸️"
	longBreakIco  string = "☕"
)

type Pomodoro struct {
	desc     string
	interval int
	icon     string
	work     bool
}

// Pomodori holds available pomodori and keeps track of their history
type Pomodori struct {
	availablePomodori map[string]Pomodoro
	elapsedPomodori   []Pomodoro
}

// NewPomodori creates a Pomodori container
func NewPomodori() *Pomodori {

	var p Pomodori

	// The maps have to be initialized or we will receive a runtime error
	// `panic: assignment to entry in nil map` later on:
	p.availablePomodori = make(map[string]Pomodoro)

	pomodoroInterval := menuet.Defaults().Integer("pomodoroInterval")
	shortBreakInterval := menuet.Defaults().Integer("shortBreakInterval")
	longBreakInterval := menuet.Defaults().Integer("longBreakInterval")

	if pomodoroInterval <= 0 {
		pomodoroInterval = 25
		menuet.Defaults().SetInteger("pomodoroInterval", pomodoroInterval)
	}

	if shortBreakInterval <= 0 {
		shortBreakInterval = 5
		menuet.Defaults().SetInteger("shortBreakInterval", shortBreakInterval)
	}

	if longBreakInterval <= 0 {
		longBreakInterval = 20
		menuet.Defaults().SetInteger("longBreakInterval", longBreakInterval)
	}

	p.availablePomodori["pomodoro"] = Pomodoro{
		desc:     "Start a new Pomodoro",
		interval: menuet.Defaults().Integer("pomodoroInterval"),
		icon:     pomodoroIco,
		work:     true,
	}
	p.availablePomodori["shortBreak"] = Pomodoro{
		desc:     "Take a short break",
		interval: menuet.Defaults().Integer("shortBreakInterval"),
		icon:     shortBreakIco,
		work:     false,
	}
	p.availablePomodori["longBreak"] = Pomodoro{
		desc:     "Take a long break",
		interval: menuet.Defaults().Integer("longBreakInterval"),
		icon:     longBreakIco,
		work:     false,
	}

	return &p
}

// StartPomodoro starts a pomodoro phase with a given interval and
// asks for the next interval after the start interval has elapsed
func StartPomodoro(startPom Pomodoro, poms *Pomodori) {
	for startPom.interval > 0 {
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: startPom.icon + strconv.Itoa(startPom.interval),
		})
		time.Sleep(1 * time.Minute)
		startPom.interval--
	}

	// Add period to history
	poms.elapsedPomodori = append(poms.elapsedPomodori, startPom)

	log.Printf("%#v", poms.elapsedPomodori)

	AlertForNextInterval(poms)
}

// AlertForNextInterval asks the user how to proceed after a pomodoro has elapsed
func AlertForNextInterval(pomodori *Pomodori) {

	history := ""
	history += fmt.Sprint("Your pomodoro history: \n\n")
	for i, p := range pomodori.elapsedPomodori {
		history += fmt.Sprintf("%d. %s (%s)\n", i+1, p.icon, p.desc)
	}

	availPomodoriKeys := []string{}
	availPomodoriDescs := []string{}

	for k, v := range pomodori.availablePomodori {
		availPomodoriKeys = append(availPomodoriKeys, k)
		availPomodoriDescs = append(availPomodoriDescs, v.desc)
	}

	response := menuet.App().Alert(menuet.Alert{
		MessageText:     fmt.Sprintf("%v elapsed", pomodori.elapsedPomodori[len(pomodori.elapsedPomodori)-1].desc),
		InformativeText: fmt.Sprintf(history),
		Buttons:         availPomodoriDescs,
	})

	switch response.Button {
	case 0:
		go StartPomodoro(pomodori.availablePomodori[availPomodoriKeys[0]], pomodori)
	case 1:
		go StartPomodoro(pomodori.availablePomodori[availPomodoriKeys[1]], pomodori)
	case 2:
		go StartPomodoro(pomodori.availablePomodori[availPomodoriKeys[2]], pomodori)
	}
}

func menuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{}

	pomodoroInterval := menuet.Defaults().String("pomodoroInterval")
	shortBreakInterval := menuet.Defaults().String("shortBreakInterval")
	longBreakInterval := menuet.Defaults().String("longBreakInterval")

	currentIntervalsDesc := fmt.Sprintf(
		`Your current intervals are set to:

%s minute for each pomodoro interval %s
%s minute for short breaks %s
%s minute for long breaks %s

Attention: After saving new values, an application restart is required to take effect.
`, pomodoroInterval, pomodoroIco, shortBreakInterval, shortBreakIco, longBreakInterval, longBreakIco)

	items = append(items, menuet.MenuItem{
		Text: "Settings...",
		Clicked: func() {
			response := menuet.App().Alert(menuet.Alert{
				MessageText:     "Set your intervals in Minutes",
				InformativeText: currentIntervalsDesc,
				Inputs:          []string{"Pomodoro", "Short break", "Long break"},
				Buttons:         []string{"Save", "Cancel"},
			})
			log.Printf("%#v", response)

			if response.Inputs[0] != "" {
				iv, err := strconv.Atoi(response.Inputs[0])
				if err != nil {
					log.Printf("Error: %v", err)
				}
				if err == nil {
					menuet.Defaults().SetInteger("pomodoroInterval", iv)
				}
			}
			if response.Inputs[1] != "" {
				iv, err := strconv.Atoi(response.Inputs[1])
				if err != nil {
					log.Printf("Error: %v", err)
				}
				if err == nil {
					menuet.Defaults().SetInteger("shortBreakInterval", iv)
				}
			}
			if response.Inputs[2] != "" {
				iv, err := strconv.Atoi(response.Inputs[2])
				if err != nil {
					log.Printf("Error: %v", err)
				}
				if err == nil {
					menuet.Defaults().SetInteger("longBreakInterval", iv)
				}
			}
		}})
	return items
}

func main() {
	app := menuet.App()
	app.Name = "Go!mato Pomodoro Timer"
	app.Label = "com.github.kekscode.gomato"
	app.Children = menuItems

	var pomodori = NewPomodori()
	go StartPomodoro(pomodori.availablePomodori["pomodoro"], pomodori)

	app.RunApplication()
}
