package main

import (
	"fmt"

	"github.com/0XF0D0/glut"
)

func main() {
	glut.Init()
	glut.InitDisplayMode(glut.SINGLE | glut.RGBA)
	glut.InitWindowSize(320, 240)
	glut.CreateWindow("c8 emulator")
	glut.ReshapeFunc(reshape)
	glut.DisplayFunc(display)
	glut.MainLoop()
}

func reshape(width, height int) {
	fmt.Println("reshape")
}

func display() {
	fmt.Println("display")
}
