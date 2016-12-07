package main

import (
	"github.com/veandco/go-sdl2/sdl"
	mix "github.com/veandco/go-sdl2/sdl_mixer"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

const SIZE = 128

// Text State text structure
type TextEntity struct {
	Value   *string
	Rect    sdl.Rect
	Texture *sdl.Texture
}

// stateText States text
var stateText = map[int]string{
	StateRun:  "RUN",
	StateFlap: "FLAP",
	StateDead: "DEAD",
}

// Text State text structure
type Text struct {
	Width   int32
	Height  int32
	Texture *sdl.Texture
}

// ShmupWarz SDL game structure
type ShmupWarz struct {
	Game
	Font     *ttf.Font
	Music    *mix.Music
	Sound    *mix.Chunk
	State    int
	frame    int
	alpha    uint8
	r        byte
	g        byte
	b        byte
	a        byte
	Entities []*Entity
	Systems  []System
}

// NewShmupWarz Returns new shmupwarz
func NewShmupWarz(width int, height int, title string) (this *ShmupWarz) {
	this = new(ShmupWarz)
	this.Width = width
	this.Height = height
	this.Title = title
	this.Mode = mix.INIT_OGG
	return
}

// Initialize the game
// overrides the base class Initialize
// called by the game engine prior to start
// use to initialize SDL submodules
func (this *ShmupWarz) Initialize() {
	this.Game.Initialize()

	return
}

// Start
// overrides the base class start
// called by the game engine prior to the game loop
// use to load resources
func (this *ShmupWarz) Start() {

	// Load resources
	this.Font = this.LoadFont("assets/fonts/skranji.regular.ttf", 24)
	this.Music = this.LoadMusic("assets/music/frantic-gameplay.ogg")
	this.Sound = this.LoadSound("assets/sounds/click.wav")

	var rects []*sdl.Rect

	// Sprite rects
	for x := 0; x < 6; x++ {
		rect := &sdl.Rect{X: int32(SIZE * x), Y: 0, W: SIZE, H: SIZE}
		rects = append(rects, rect)
	}

	this.Entities = []*Entity{
		this.CreateBackgroundEntity(),
		this.CreateGopherEntity(0, rects[0:2]),
		this.CreateGopherEntity(1, rects[2:4]),
		this.CreateGopherEntity(2, rects[4:6]),
		this.CreateTextEntity(0, stateText[0]),
		this.CreateTextEntity(1, stateText[1]),
		this.CreateTextEntity(2, stateText[2])}

	this.Systems = append(this.Systems, NewRenderSystem(this))
	// Play music
	this.Music.Play(-1)

	this.alpha = 255
	this.Game.Start()
}

// OnEvent
func (this *ShmupWarz) OnEvent(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.QuitEvent:
		this.Game.Quit()

	case *sdl.MouseButtonEvent:
		this.Sound.Play(2, 0)
		if t.Type == sdl.MOUSEBUTTONDOWN && t.Button == sdl.BUTTON_LEFT {
			this.alpha = 255
			if this.State == StateRun {
				this.State = StateFlap
			} else if this.State == StateFlap {
				this.State = StateDead
			} else if this.State == StateDead {
				this.State = StateRun
			}
		}

	case *sdl.KeyDownEvent:
		if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE || t.Keysym.Scancode == sdl.SCANCODE_AC_BACK {
			this.Game.Quit()
		}
	}

}

// Update
// Implenents the abstract method Update
// game logic, physics, etc goes here
func (this *ShmupWarz) Update(delta float64) {

	switch this.State {
	case StateRun:
		this.r = 168
		this.g = 235
		this.b = 254
		this.a = 255

	case StateFlap:
		this.r = 251
		this.g = 231
		this.b = 240
		this.a = 255

	case StateDead:
		this.r = 255
		this.g = 250
		this.b = 205
		this.a = 255
	}

	this.frame++
	if this.frame/2 >= 2 {
		this.frame = 0
	}

	this.alpha -= 10
	if this.alpha <= 10 {
		this.alpha = 255
	}
}

// Draw
// Implenents the abstract method Draw
// do all the rendering
func (this *ShmupWarz) Draw(delta float64) {
	this.Renderer.Clear()
	this.SystemDraw()
	this.Renderer.Present()

}

func (this *ShmupWarz) SystemDraw() {

	w, h := this.Window.GetSize()
	x, y := int32(w/2), int32(h/2)

	for i := 0; i < len(this.Entities); i++ {
		en := this.Entities[i]
		switch en.EntityType.Ordinal {
		case EntityTypeBackground:
			this.Renderer.Copy(en.Sprite.Texture, nil, nil)

		case EntityTypeGopher:
			if this.State == en.State.Value {
				en.Sprite.Texture.SetColorMod(this.r, this.g, this.b)
				clip := en.Sprite.Clip[this.frame/2]
				this.Renderer.Copy(en.Sprite.Texture, clip, &sdl.Rect{X: x - (SIZE / 2), Y: y - (SIZE / 2), W: SIZE, H: SIZE})

			}
		case EntityTypeText:
			if this.State == en.State.Value {
				en.Sprite.Texture.SetAlphaMod(this.alpha)
				this.Renderer.Copy(en.Sprite.Texture, nil, &sdl.Rect{X: x - (en.Text.Rect.W / 2), Y: y - SIZE*1.5, W: en.Text.Rect.W, H: en.Text.Rect.H})

			}
		}
	}
}

// Destroy Destroys SDL and releases the memory
// overrides the base class Destroy
func (this *ShmupWarz) Destroy() {
	this.Font.Close()
	this.Music.Free()
	this.Sound.Free()

	ttf.Quit()
	mix.CloseAudio()
	mix.Quit()

	this.Game.Destroy()

}
