package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	EntityTypeBackground = iota
	EntityTypeGopher
	EntityTypeText
)

var entityType = map[int]string{
	EntityTypeBackground: "Background",
	EntityTypeGopher:     "Gopher",
	EntityTypeText:       "Text",
}

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
	Texture *sdl.Texture
	Clip    []*sdl.Rect
}
type TextComponent struct {
	AbstractComponent
	Value   string
	Rect    sdl.Rect
}

var nextEntityId = 0

func GetEntityId() int {
	nextEntityId++
	return nextEntityId
}

type Entity struct {
	Id         int
	Active     bool
	EntityType EntityTypeComponent
	State      StateComponent
	Position   PositionComponent
	Color      ColorComponent
	Sprite     SpriteComponent
	Text       TextComponent
}

func (this *ShmupWarz) CreateGopherEntity(state int, clip []*sdl.Rect) (e *Entity) {
	e = new(Entity)
	e.Id = GetEntityId()
	e.EntityType.Ordinal = EntityTypeGopher
	e.Sprite.Texture = this.LoadTexture("assets/images/sprite.png")
	e.Sprite.Clip = clip
	e.State.Value = state
	return e
}

func (this *ShmupWarz) CreateBackgroundEntity() (e *Entity) {
	e = new(Entity)
	e.Id = GetEntityId()
	e.Active = true
	e.EntityType.Ordinal = EntityTypeBackground
	e.Sprite.Texture = this.LoadTexture("assets/images/BackdropBlackLittleSparkBlack.png")
	return e
}

func (this *ShmupWarz) CreateTextEntity(state int, text string) (e *Entity) {
	e = new(Entity)
	e.Id = GetEntityId()
	e.EntityType.Ordinal = EntityTypeText
	s, err := this.Font.RenderUTF8_Blended(text, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		return
	}
	defer s.Free()
	t, _ := this.Renderer.CreateTextureFromSurface(s)
	_, _, tW, tH, _ := t.Query()
	e.Sprite.Texture = t
	e.Text.Value = text
	e.Text.Rect.W = tW
	e.Text.Rect.H = tH
	e.State.Value = state
	return e
}

type System interface {
	Initialize()
	Update(delta float64)
}

type RenderSystem struct {
	Game  *ShmupWarz
	Group []int
}

func NewRenderSystem(game *ShmupWarz) (this *RenderSystem) {
	this = &RenderSystem{}
	this.Game = game
	return
}

func (this *RenderSystem) Initialize() {
	this.Group = append(this.Group, 0)
}
func (this *RenderSystem) Update(delta float64) {

	for i := 0; i < len(this.Group); i++ {
		entity := this.Game.Entities[this.Group[i]]
		if entity.Active {
			this.Game.Renderer.Copy(entity.Sprite.Texture, nil, nil)
		}

	}

}
