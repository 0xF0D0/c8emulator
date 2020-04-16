package main

import (
	"unsafe"

	"github.com/0XF0D0/glut"
)

var screenData [32][64][3]byte
var cnt int = 0

func main() {
	glut.Init()
	glut.InitDisplayMode(glut.SINGLE | glut.RGBA)
	glut.InitWindowSize(640, 320)
	glut.CreateWindow("c8 emulator")
	glut.ReshapeFunc(reshape)
	glut.DisplayFunc(display)
	glut.IdleFunc(display)

	glut.TexImage2DRGBByte(64, 32, unsafe.Pointer(&screenData))
	glut.SetTexParameteri()
	glut.EnableTexture()

	glut.MainLoop()

}

func draw() {
	glut.Clear()

	for i := 0; i < 32; i++ {
		for j := 0; j < 64; j++ {
			for k := 0; k < 3; k++ {
				screenData[i][j][k] = 122
			}
		}
	}

	glut.TexSubImage2DRGBByte(64, 32, unsafe.Pointer(&screenData))
	glut.DrawTexture(640, 320)
	glut.SwapBuffers()
}

func reshape(width, height int) {
	glut.ClearColor(0.0, 0.0, 0.5, 0.0)
	glut.SetMatrixModeProjection()
	glut.LoadIdentity()
	glut.Ortho2D(width, height)
	glut.SetMatrixModeModelView()
	glut.ViewPort(width, height)
}

func display() {
	if cnt < 100000 {
		cnt++
		return
	} else {
		draw()
		cnt = 0
	}
}
