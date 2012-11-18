package main

import (
	//"log"
	"math"
	//"time"

	"github.com/jteeuwen/glfw"
)

func initInput() {
	// Start listening to mouse events
	glfw.SetMouseButtonCallback(click)
	glfw.SetMousePosCallback(move)
}

var isHeld bool
var ot, op, or float64
var ox, oy int

func click(button, state int) {

	if button == glfw.Mouse1 {
		if state == 1 {
			isHeld = true
			ot, op, or = theta, phi, r
			ox, oy = glfw.MousePos()
		} else {
			isHeld = false
		}
	}
}

var isFirst bool

func move(x, y int) {
	if isHeld {
		dx, dy := float64(x-ox), float64(y-oy)
		s := 100.0
		dt, dp := dy/s, dx/s
		theta, phi = math.Mod(ot+dt, math.Pi), math.Mod(op+dp, math.Pi*2)
	}
}
