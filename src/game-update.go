/*
//  Implementation of the Update method for the Game structure
//  This method is called once at every frame (60 frames per second)
//  by ebiten, juste before calling the Draw method (game-draw.go).
//  Provided with a few utilitary methods:
//    - CheckArrival
//    - ChooseRunners
//    - HandleLaunchRun
//    - HandleResults
//    - HandleWelcomeScreen
//    - Reset
//    - UpdateAnimation
//    - UpdateRunners
*/

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"time"
)

// HandleWelcomeScreen waits for the player to push SPACE in order to
// start the game
func (g *Game) HandleWelcomeScreen() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

// HandleWaitForPlayers waits for the players to join the server
func (g *Game) HandleWaitForConnexion() bool {
    return g.stateServer == 1
}

// ChooseRunners loops over all the runners to check which sprite each
// of them selected
func (g *Game) ChooseRunners() (done bool) {
	done = true
	for i := range g.runners {
		if i == 0 {
		    currentChoose := g.runners[i].colorScheme
			done = g.runners[i].ManualChoose() && done

			if currentChoose != g.runners[i].colorScheme {
			    g.SendDeplacementMenuRunner()
            }
		} else {
			//done = g.runners[i].RandomChoose() && done
		}
	}
	return done
}

func (g *Game) HandleWaitForPlayers() bool {
    return g.stateServer == 2
}

// HandleLaunchRun countdowns to the start of a run
func (g *Game) HandleLaunchRun() bool {
	if time.Since(g.f.chrono).Milliseconds() > 1000 {
		g.launchStep++
		g.f.chrono = time.Now()
	}
	if g.launchStep >= 5 {
		g.launchStep = 0
		return true
	}
	return false
}

// UpdateRunners loops over all the runners to update each of them
func (g *Game) UpdateRunners() {
	for i := range g.runners {
		if i == 0 {
		    xpos := g.runners[i].xpos

			g.runners[i].ManualUpdate()

			if xpos != g.runners[i].xpos {
			    g.SendPosition()
			}
		} else {
		}
	}
}

// CheckArrival loops over all t he runners to check which ones are arrived
func (g *Game) CheckArrival() (finished bool) {
	finished = true
	for i := range g.runners {
		g.runners[i].CheckArrival(&g.f)
		finished = finished && g.runners[i].arrived
	}
	return finished
}

// Reset resets all the runners and the field in order to start a new run
func (g *Game) Reset() {
	for i := range g.runners {
		g.runners[i].Reset(&g.f)
	}
	g.f.Reset()
}

// UpdateAnimation loops over all the runners to update their sprite
func (g *Game) UpdateAnimation() {
	for i := range g.runners {
		g.runners[i].UpdateAnimation(g.runnerImage)
	}
}

func (g *Game) HandleWaitForResults() bool {
    return g.stateServer == 3
}

// HandleResults computes the resuls of a run and prepare them for
// being displayed
func (g *Game) HandleResults() bool {
	if time.Since(g.f.chrono).Milliseconds() > 1000 || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.resultStep++
		g.f.chrono = time.Now()
	}
	if g.resultStep >= 4 && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.resultStep = 0
		return true
	}
	return false
}

// Update is the main update function of the game. It is called by ebiten
// at each frame (60 times per second) just before calling Draw (game-draw.go)
// Depending of the current state of the game it calls the above utilitary
// function and then it may update the state of the game
func (g *Game) Update() error {
	switch g.state {
	case StateWelcomeScreen:
		done := g.HandleWelcomeScreen()
		if done {
		    go g.ConnectToServer()
			g.state++
		}
	case StateWaitForConnexion:
	    done := g.HandleWaitForConnexion()
	    if done {
	        g.state++
        }
	case StateChooseRunner:
		done := g.ChooseRunners()
		if done {
			go g.SendRunner()
			g.state++
		}
	case StateWaitForRunner:
		done := g.HandleWaitForPlayers()
		if done {
            g.UpdateAnimation()
			g.state++
		}
	case StateLaunchRun:
		done := g.HandleLaunchRun()
		if done {
			g.state++
		}
	case StateRun:
		g.UpdateRunners()
		finished := g.CheckArrival()
		g.UpdateAnimation()
		if finished {
			g.SendResults()
			g.state++
		}
	case StateWaitForResults:
	    done := g.HandleWaitForResults()
	    if done {
	        g.state++
        }
	case StateResult:
		done := g.HandleResults()
		if done {
			g.Reset()
			g.state = StateWaitForRunner
			g.RestartGame()
			g.stateServer = 1
		}
	}
	return nil
}
