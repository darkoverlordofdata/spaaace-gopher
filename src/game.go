package main

//import "C"

import (
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
	mix "github.com/veandco/go-sdl2/sdl_mixer"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
)

const (
	StateRun = iota
	StateFlap
	StateDead
)

// Game
// base game class modeled after Microsoft.Xna.Game
type Game struct {
	Title    string
	Width    int
	Height   int
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Delta    float64
	Mode     int
	running  bool
	err      error
}

// IGame
// provide implementation for Game subclass
type IGame interface {
	Initialize()
	Start()
	Update(delta float64)
	Draw(delta float64)
	OnEvent(evt sdl.Event)
}

// Running
// check if loop is running
func (this *Game) Running() bool {
	return this.running
}

// Start the main loop
func (this *Game) Start() {
	this.running = true
}

// Quit Quits main loop
func (this *Game) Quit() {
	this.running = false
}

// Initialize the game
func (this *Game) Initialize() {
	runtime.LockOSThread()
	this.err = sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO)
	if this.err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "Init: %s\n", this.err)
		return
	}

	//mix.INIT_OGG
	this.err = mix.Init(this.Mode)
	if this.err != nil {
		return
	}
	this.err = mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, mix.DEFAULT_CHANNELS, 3072)
	if this.err != nil {
		return
	}
	this.err = ttf.Init()
	if this.err != nil {
		return
	}
	this.Window, this.err = sdl.CreateWindow(this.Title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, this.Width, this.Height, sdl.WINDOW_SHOWN)
	if this.err != nil {
		return
	}

	this.Renderer, this.err = sdl.CreateRenderer(this.Window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if this.err != nil {
		return
	}

}

// Destroy the game
func (this *Game) Destroy() {
	this.Renderer.Destroy()
	this.Window.Destroy()
	sdl.Quit()

}

// Run the game loop
func (this *Game) Run(subclass IGame) {
	var lastTime float64
	var curTime float64
	var event sdl.Event

	subclass.Initialize()
	subclass.Start()
	lastTime = float64(time.Now().UnixNano()) / 1000000.0
	for this.Running() {
		curTime = float64(time.Now().UnixNano()) / 1000000.0
		this.Delta = curTime - lastTime
		lastTime = curTime

		event = sdl.PollEvent()
		if event != nil {
			subclass.OnEvent(event)
		}

		subclass.Update(this.Delta)
		subclass.Draw(this.Delta)
	}

}

// LoadTexture from path
func (this *Game) LoadTexture(path string) (texture *sdl.Texture) {
	var err error

	texture, err = img.LoadTexture(this.Renderer, path)
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "Load Texture: %s\n", err)
	}
	return
}

// LoadFont from path
func (this *Game) LoadFont(path string, size int) (font *ttf.Font) {
	var err error

	font, err = ttf.OpenFont(path, size)
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "Load Font: %s\n", err)
	}
	return
}

// LoadMusic from path
func (this *Game) LoadMusic(path string) (music *mix.Music) {
	var err error

	music, err = mix.LoadMUS(path)
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "Load Music: %s\n", err)
	}
	return
}

// LoadSound from path
func (this *Game) LoadSound(path string) (sound *mix.Chunk) {
	var err error

	sound, err = mix.LoadWAV(path)
	if err != nil {
		sdl.LogError(sdl.LOG_CATEGORY_APPLICATION, "Load Sound: %s\n", err)
	}
	return
}
