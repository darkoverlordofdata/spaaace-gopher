package main

import "github.com/veandco/go-sdl2/sdl"

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
	Texture *sdl.Texture
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
	e.Sprite.Texture = this.LoadTexture("assets/images/sprite.png")
	return
}

func (this *ShmupWarz) CreateBackgroundEntity() (e *Entity) {
	e = new(Entity)
	e.EntityType.Ordinal = EntityTypeBackground
	e.Sprite.Texture = this.LoadTexture("assets/images/BackdropBlackLittleSparkBlack.png")
	return
}

func (this *ShmupWarz) CreateTextEntity() (e *Entity) {
	e = new(Entity)
	e.EntityType.Ordinal = EntityTypeText
	return
}
