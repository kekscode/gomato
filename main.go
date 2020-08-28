package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/caseymrm/menuet"
)

type pomodoro struct {
	desc     string
	interval int
	icon     string
	work     bool
}

// Pomodori holds available pomodori and keeps track of their history
type Pomodori struct {
	availablePomodori map[string]pomodoro
	elapsedPomodori   []pomodoro
}

// NewPomodori creates a Pomodori container
func NewPomodori() *Pomodori {

	var p Pomodori

	// The maps have to be initialized or we will receive a runtime error
	// `panic: assignment to entry in nil map` later on:
	p.availablePomodori = make(map[string]pomodoro)

	p.availablePomodori["pomodoro"] = pomodoro{
		desc:     "Start a new Pomodoro",
		interval: 25,
		icon:     "🍅",
		work:     true,
	}
	p.availablePomodori["shortBreak"] = pomodoro{
		desc:     "Take a short break",
		interval: 5,
		icon:     "⏸️",
		work:     false,
	}
	p.availablePomodori["longBreak"] = pomodoro{
		desc:     "Take a long break",
		interval: 20,
		icon:     "☕",
		work:     false,
	}

	return &p
}

// StartPomodoro starts a pomodoro phase with a given interval and
// asks for the next interval after the start interval has elapsed
func StartPomodoro(startPom pomodoro, poms *Pomodori) {
	for startPom.interval >= 0 {
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: startPom.icon + strconv.Itoa(startPom.interval),
		})
		time.Sleep(1 * time.Minute)
		startPom.interval--
	}

	// Add period to history
	poms.elapsedPomodori = append(poms.elapsedPomodori, startPom)

	fmt.Printf("%#v", poms.elapsedPomodori)

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

func main() {
	var pomodori = NewPomodori()
	go StartPomodoro(pomodori.availablePomodori["pomodoro"], pomodori)
	app := menuet.App()
	app.Name = "Go!mato Pomodoro Timer"
	app.Label = "com.github.kekscode.gomato"
	app.RunApplication()
}
