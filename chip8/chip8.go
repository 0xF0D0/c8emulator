//Package chip8 inplements Chip8 emulator.
package chip8

import (
	"fmt"
	"log"
	"math/rand"
	"os"
)

var fontset = [80]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //0
	0x20, 0x60, 0x20, 0x20, 0x70, //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
	0x90, 0x90, 0xF0, 0x10, 0x10, //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
	0xF0, 0x10, 0x20, 0x40, 0x40, //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80, //F
}

// Chip8 emulator.
type Chip8 struct {
	gfxChannel    chan []byte
	opcode        uint16
	memory        [4096]byte
	v             [16]byte
	indexRegister uint16
	pc            uint16
	gfx           [64 * 32]byte
	soundTimer    byte
	delayTimer    byte
	stack         [16]uint16
	sp            uint16
	key           [16]byte
}

//Initialize chip8 emulator
func Initialize() *Chip8 {
	c8 := &Chip8{}
	c8.pc = 0x200

	for i := 0; i < 80; i++ {
		c8.memory[i] = fontset[i]
	}
	c8.gfxChannel = make(chan []byte)

	return c8
}

//EmulateCycle emulates one instruction.
func (c *Chip8) EmulateCycle() {
	c.opcode = uint16(c.memory[c.pc]<<8) | uint16(c.memory[c.pc+1])
	switch c.opcode & 0xF000 {
	case 0x0000:
		switch c.opcode & 0x000F {
		case 0x0000: // 0x00E0: Clear screen
			for i := 0; i < 2048; i++ {
				c.gfx[i] = 0x0
			}
			c.gfxChannel <- c.gfx[:]
			c.pc += 2
		case 0x000E: // 0x00EE: Return from subroutine
			c.sp--
			c.pc = c.stack[c.sp]
			c.pc += 2
		default:
			log.Fatalf("Unknown opcode: 0x%X", c.opcode)
		}
	case 0x1000: // 0x1NNN: Jump to address NNN
		c.pc = c.opcode & 0x0FFF
	case 0x2000: // 0x2NNN: Call subroutine at NNN
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = c.opcode & 0x0FFF
	case 0x3000: // 0x3XNN: Skip next instruction if VX equals NN
		if c.v[(c.opcode&0x0F00)>>8] == byte(c.opcode&0x00FF) {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0x4000: // 0x4XNN: Skip next instruction if VX not equals NN
		if c.v[(c.opcode&0x0F00)>>8] != byte(c.opcode&0x00FF) {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0x5000: // 0x5XY0: Skip next instruction if VX equals VY
		if c.v[(c.opcode&0x0F00)>>8] == c.v[(c.opcode&0x00F0)>>4] {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0x6000: // 0x6XNN: Set VX to NN
		c.v[(c.opcode&0x0F00)>>8] = byte(c.opcode & 0x00FF)
		c.pc += 2
	case 0x7000: // 0x7XNN: Add NN to VX
		c.v[(c.opcode&0x0F00)>>8] += byte(c.opcode & 0x00FF)
		c.pc += 2
	case 0x8000:
		switch c.opcode & 0x000F {
		case 0x0000: // 0x8XY0: Set VX to value of VY
			c.v[(c.opcode&0x0F00)>>8] = c.v[(c.opcode&0x00F0)>>4]
			c.pc += 2
		case 0x0001: // 0x8XY1: Set VX to "VX OR VY"
			c.v[(c.opcode&0x0F00)>>8] |= c.v[(c.opcode&0x00F0)>>4]
			c.pc += 2
		case 0x0002: // 0x8XY2: Set VX to "VX AND VY"
			c.v[(c.opcode&0x0F00)>>8] &= c.v[(c.opcode&0x00F0)>>4]
			c.pc += 2
		case 0x0003: // 0x8XY3: Set VX to "VX XOR VY"
			c.v[(c.opcode&0x0F00)>>8] ^= c.v[(c.opcode&0x00F0)>>4]
			c.pc += 2
		case 0x0004: // 0x8XY4: Add VY to VX. VF is set when there is carry
			if c.v[(c.opcode&0x00F0)>>4] > (0xFF - c.v[(c.opcode&0x0F00)>>8]) {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[(c.opcode&0x0F00)>>8] += c.v[(c.opcode&0x00F0)>>4]
			c.pc += 2
		case 0x0005: // 0x8XY5: Substract VY from VX. VF is 0 when underflow, otherwise, set.
			if c.v[(c.opcode&0x00F0)>>4] > c.v[(c.opcode&0x0F00)>>8] {
				c.v[0xF] = 0
			} else {
				c.v[0xF] = 1
			}
			c.v[(c.opcode&0x0F00)>>8] -= c.v[(c.opcode&0x00F0)>>4]
			c.pc += 2
		case 0x0006: // 0x8XY6: Shift right VX by one. VF is set if lsb was 1
			c.v[0xF] = c.v[(c.opcode&0x0F00)>>8] & 0x1
			c.v[(c.opcode&0x0F00)>>8] >>= 1
			c.pc += 2
		case 0x0007: // 0x8XY7: Set VX to VY minus VX. VF is 0 when underflow, otherwise, set.
			if c.v[(c.opcode&0x0F00)>>8] > c.v[(c.opcode&0x00F0)>>4] {
				c.v[0xF] = 0
			} else {
				c.v[0xF] = 1
			}
			c.v[(c.opcode&0x0F00)>>8] = c.v[(c.opcode&0x00F0)>>4] - c.v[(c.opcode&0x0F00)>>8]
			c.pc += 2
		case 0x000E: // 0x8XYE: Shift left VX by one. VF is set if msb was 1
			c.v[0xF] = c.v[(c.opcode&0x0F00)>>8] >> 7
			c.v[(c.opcode&0x0F00)>>8] <<= 1
			c.pc += 2
		default:
			log.Fatalf("Unknown opcode: 0x%X", c.opcode)
		}
	case 0x9000: // 0x9XY0: Skip next instruction if VX is not equal to VY
		if c.v[(c.opcode&0x0F00)>>8] != c.v[(c.opcode&0x00F0)>>4] {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0xA000: // 0xANNN: Set I to address NNN
		c.indexRegister = c.opcode & 0x0FFF
		c.pc += 2
	case 0xB000: // 0xBNNN: Jump to address "NNN plus V0"
		c.pc = (c.opcode & 0x0FFF) + uint16(c.v[0])
	case 0xC000: // 0xCNNN: Set VX to random number AND NN
		c.v[(c.opcode&0x0F00)>>8] = byte(rand.Intn(0xFF)) & byte(c.opcode&0x00FF)
		c.pc += 2
	case 0xD000:
		// 0xDXYN: Draw sprite at (VX, VY) with width of 8 and height of N.
		// Each row of 8 pixels are read from I register.
		// VF is set if any pixel has been flipped
		x := uint16(c.v[(c.opcode&0x0F00)>>8])
		y := uint16(c.v[(c.opcode&0x00F0)>>4])
		height := c.opcode & 0x000F
		c.v[0xF] = 0

		for yline := uint16(0); yline < height; yline++ {
			pixel := c.memory[c.indexRegister+yline]
			for xline := uint16(0); xline < 8; xline++ {
				if (pixel & (0x80 >> xline)) != 0 {
					if c.gfx[x+xline+(y+yline)*64] == 1 {
						c.v[0xF] = 1
					}
					c.gfx[x+xline+(y+yline)*64] ^= 1
				}
			}
		}

		c.gfxChannel <- c.gfx[:]
		c.pc += 2
	case 0xE000:
		switch c.opcode & 0x00FF {
		case 0x009E: // 0xEX9E: Skip next instruction if key[VX] is pressed
			if c.key[c.v[(c.opcode&0x0F00)>>8]] != 0 {
				c.pc += 4
			} else {
				c.pc += 2
			}
		case 0x00A1: // 0xEXA1: Skip next instruction if key[VX] is not pressed
			if c.key[c.v[(c.opcode&0x0F00)>>8]] == 0 {
				c.pc += 4
			} else {
				c.pc += 2
			}
		default:
			log.Fatalf("Unknown opcode: 0x%X", c.opcode)
		}
	case 0xF000:
		switch c.opcode & 0x00FF {
		case 0x0007: // 0xFX07: set VX to value of delaytimer
			c.v[(c.opcode&0x0F00)>>8] = c.delayTimer
			c.pc += 2
		case 0x000A: // 0xFX0A: keypress is awaited, then store in VX
			keyPress := false
			for i := 0; i < 16; i++ {
				if c.key[i] != 0 {
					c.v[(c.opcode&0x0F00)>>8] = byte(i)
					keyPress = true
				}
			}

			// If there was no keypress, skip this cycle
			if !keyPress {
				return
			}
			c.pc += 2
		case 0x0015: // 0xFX15: Set delaytimer to VX
			c.delayTimer = c.v[(c.opcode&0x0F00)>>8]
			c.pc += 2
		case 0x0018: // 0xFX18: Set soundtimer to VX
			c.soundTimer = c.v[(c.opcode&0x0F00)>>8]
			c.pc += 2
		case 0x001E: // 0xFX1E: Add VX to I
			// VF is set when there is range overflow(0xFFF)
			if c.indexRegister+uint16(c.v[(c.opcode&0x0F00)>>8]) > 0xFFF {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.indexRegister += uint16(c.v[(c.opcode&0x0F00)>>8])
			c.pc += 2
		case 0x0029: // 0xFX29: Set I to the location of charater in VX
			c.indexRegister = uint16(c.v[(c.opcode&0x0F00)>>8]) * 0x5
		case 0x0033: // 0xFX33: Store Binary-coded decimal representation of VX at I
			c.memory[c.indexRegister] = c.v[(c.opcode&0x0F00)>>8] / 100
			c.memory[c.indexRegister+1] = (c.v[(c.opcode&0x0F00)>>8] / 10) % 10
			c.memory[c.indexRegister+2] = (c.v[(c.opcode&0x0F00)>>8] % 100) % 10
			c.pc += 2
		case 0x0055: // 0xFX55: Store V0 to VX in memory starting at address I
			for i := uint16(0); i < (c.opcode*0x0F00)>>8; i++ {
				c.memory[c.indexRegister+i] = c.v[i]
			}

			// Then I is set to I + X + 1
			c.indexRegister += (c.opcode&0x0F00)>>8 + 1
			c.pc += 2
		case 0x0065: // 0xFX65: Load V0 to VX from memory starting at address I
			for i := uint16(0); i < (c.opcode*0x0F00)>>8; i++ {
				c.v[i] = c.memory[c.indexRegister+i]
			}

			// Then I is set to I + X + 1
			c.indexRegister += (c.opcode&0x0F00)>>8 + 1
			c.pc += 2
		default:
			log.Fatalf("Unknown opcode: 0x%X", c.opcode)
		}
	default:
		log.Fatalf("Unknown opcode: 0x%X", c.opcode)
	}

	if c.delayTimer > 0 {
		c.delayTimer--
	}
	if c.soundTimer > 0 {
		if c.soundTimer == 1 {
			fmt.Println("Beeeep!!!")
		}
		c.soundTimer--
	}
}

//LoadGame loads game in given path
func (c *Chip8) LoadGame(dir string) {
	f, err := os.Open(dir)
	if err != nil {
		log.Fatal(err)
	}
	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fSize := fi.Size()
	if fSize >= (4096 - 512) {
		log.Fatal("Error: Rom is too big for memory")
	}
	buffer := make([]byte, fSize)
	_, err = f.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < int(fSize); i++ {
		c.memory[i+512] = buffer[i]
	}

	f.Close()
}

//GfxChannel returns gfx output.
func (c *Chip8) GfxChannel() <-chan []byte {
	return c.gfxChannel
}
