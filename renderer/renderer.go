package renderer

import (
	"unsafe"

	"github.com/0xF0D0/glut"
)

type Renderer struct {
	KeyboardDown chan byte
	KeyboardUp   chan byte
	screenData   [32][64][3]byte
	drawFlag     bool
}

func Initialize() *Renderer {
	nr := &Renderer{}
	nr.KeyboardDown = make(chan byte)
	nr.KeyboardUp = make(chan byte)

	glut.Init()
	glut.InitDisplayMode(glut.SINGLE | glut.RGBA)
	glut.InitWindowSize(640, 320)
	glut.CreateWindow("c8 emulator")
	glut.ReshapeFunc(nr.reshape)
	glut.DisplayFunc(nr.display)
	glut.IdleFunc(nr.display)
	glut.KeyboardFunc(nr.keyboardDown)
	glut.KeyboardUpFunc(nr.keyboardUp)

	screenData := nr.screenData
	glut.TexImage2DRGBByte(64, 32, unsafe.Pointer(&screenData[0]))
	glut.SetTexParameteri()
	glut.EnableTexture()

	return nr
}

func (r *Renderer) BindRenderInput(input <-chan []byte) {
	go func() {
		for v := range input {
			if len(v) != 32*64 {
				continue
			}
			for y := 0; y < 32; y++ {
				for x := 0; x < 64; x++ {
					var rgbv byte = 255
					if v[y*64+x] == 0 {
						rgbv = 0
					}
					r.screenData[y][x][0] = rgbv
					r.screenData[y][x][1] = rgbv
					r.screenData[y][x][2] = rgbv
				}
			}
			r.drawFlag = true
		}
	}()
}

func (r *Renderer) RunMainLoop() {
	glut.MainLoop()
}

func (r *Renderer) draw() {
	screenData := r.screenData
	glut.Clear()
	glut.TexSubImage2DRGBByte(64, 32, unsafe.Pointer(&screenData[0]))
	glut.DrawTexture(640, 320)
	glut.SwapBuffers()
}

func (r *Renderer) reshape(width, height int) {
	glut.ClearColor(0.0, 0.0, 0.5, 0.0)
	glut.SetMatrixModeProjection()
	glut.LoadIdentity()
	glut.Ortho2D(width, height)
	glut.SetMatrixModeModelView()
	glut.ViewPort(width, height)
}

func (r *Renderer) display() {
	if r.drawFlag {
		r.draw()
		r.drawFlag = false
	}
}

func (r *Renderer) keyboardDown(key byte, x, y int) {
	go func() {
		r.KeyboardDown <- key
	}()
}

func (r *Renderer) keyboardUp(key byte, x, y int) {
	go func() {
		r.KeyboardUp <- key
	}()
}
