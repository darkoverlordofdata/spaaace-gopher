package main

//import "C"

import (
	"github.com/veandco/go-sdl2/sdl"
	mix "github.com/veandco/go-sdl2/sdl_mixer"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

const SIZE = 128

const (
	EntityTypeBackground = iota
	EntityTypeGopher     = iota
	EntityTypeText       = iota
)

type AbstractComponent struct{}
type EntityTypeComponent struct {
	AbstractComponent
	Ordinal int
}
type StateComponent struct {
	AbstractComponent
	Value int
}
type PositionComponent struct {
	AbstractComponent
	X int
	Y int
}
type ColorComponent struct {
	AbstractComponent
	R byte
	G byte
	B byte
	A byte
}
type SpriteComponent struct {
	AbstractComponent
	Sprite *sdl.Texture
}
type TextComponent struct {
	AbstractComponent
	Value   *string
	Rect    sdl.Rect
	Texture *sdl.Texture
}

type Entity struct {
	EntityType EntityTypeComponent
	State      StateComponent
	Position   PositionComponent
	Color      ColorComponent
	Sprite     SpriteComponent
	Text       TextComponent
}

func (this *ShmupWarz) CreateGopherEntity() (e *Entity) {
	e = new(Entity)
	e.EntityType.Ordinal = EntityTypeGopher
	e.Sprite = this.LoadTexture("assets/images/sprite.png")
	return
}

func (this *ShmupWarz) CreateBackgroundEntity() (e *Entity) {
	e = new(Entity)
	e.EntityType.Ordinal = EntityTypeBackground
	e.Sprite = this.LoadTexture("assets/images/BackdropBlackLittleSparkBlack.png")
	return
}

func (this *ShmupWarz) CreateTextEntity() (e *Entity) {
	e = new(Entity)
	e.EntityType.Ordinal = EntityTypeText
	return
}

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
	State      int
	Sprite     *sdl.Texture
	Background *sdl.Texture
	Font       *ttf.Font
	Music      *mix.Music
	Sound      *mix.Chunk
	StateText  map[int]*Text
	Sprites    []*sdl.Texture
	rects      []*sdl.Rect
	clips      []*sdl.Rect
	frame      int
	alpha      uint8
	text       *Text
	r          byte
	g          byte
	b          byte
	a          byte
	//Entities   []interface{}
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
	// GopherEntity{},
	// TextEntity{}}

	// this.Entities = []interface{}{
	// 	this.CreateGopherEntity(),
	// 	this.CreateTextEntity()}

	// Sprite rects
	for x := 0; x < 6; x++ {
		rect := &sdl.Rect{X: int32(SIZE * x), Y: 0, W: SIZE, H: SIZE}
		this.rects = append(this.rects, rect)
	}

	// Load resources
	this.Sprite = this.LoadTexture("assets/images/sprite.png")
	this.Background = this.LoadTexture("assets/images/BackdropBlackLittleSparkBlack.png")
	this.Font = this.LoadFont("assets/fonts/skranji.regular.ttf", 24)
	this.Music = this.LoadMusic("assets/music/frantic-gameplay.ogg")
	this.Sound = this.LoadSound("assets/sounds/click.wav")

	this.StateText = map[int]*Text{}
	// pre-render the text
	for k, v := range stateText {
		s, e := this.Font.RenderUTF8_Blended(v, sdl.Color{R: 255, G: 255, B: 255, A: 0})
		if e != nil {
			continue
		}
		defer s.Free()
		t, _ := this.Renderer.CreateTextureFromSurface(s)
		_, _, tW, tH, _ := t.Query()
		this.StateText[k] = &Text{tW, tH, t}
	}

	// Play music
	this.Music.Play(-1)

	this.alpha = 255
	this.text = this.StateText[StateRun]
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
				this.text = this.StateText[StateFlap]
				this.State = StateFlap
			} else if this.State == StateFlap {
				this.text = this.StateText[StateDead]
				this.State = StateDead
			} else if this.State == StateDead {
				this.text = this.StateText[StateRun]
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
		this.clips = this.rects[0:2]

	case StateFlap:
		this.r = 251
		this.g = 231
		this.b = 240
		this.a = 255
		this.clips = this.rects[2:4]

	case StateDead:
		this.r = 255
		this.g = 250
		this.b = 205
		this.a = 255
		this.clips = this.rects[4:6]
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
	w, h := this.Window.GetSize()
	x, y := int32(w/2), int32(h/2)
	clip := this.clips[this.frame/2]

	this.Renderer.Clear()
	this.Renderer.Copy(this.Background, nil, nil)
	this.Sprite.SetColorMod(this.r, this.g, this.b)
	this.Renderer.Copy(this.Sprite, clip, &sdl.Rect{X: x - (SIZE / 2), Y: y - (SIZE / 2), W: SIZE, H: SIZE})
	this.text.Texture.SetAlphaMod(this.alpha)
	this.Renderer.Copy(this.text.Texture, nil, &sdl.Rect{X: x - (this.text.Width / 2), Y: y - SIZE*1.5, W: this.text.Width, H: this.text.Height})
	this.Renderer.Present()

}

// Destroy Destroys SDL and releases the memory
// overrides the base class Destroy
func (this *ShmupWarz) Destroy() {
	for _, v := range this.StateText {
		v.Texture.Destroy()
	}

	this.Sprite.Destroy()
	this.Font.Close()
	this.Music.Free()
	this.Sound.Free()

	ttf.Quit()
	mix.CloseAudio()
	mix.Quit()

	this.Game.Destroy()

}
