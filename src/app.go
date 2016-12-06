package main

// main
func main() {

	game := NewShmupWarz(800, 600, "Spaaace Gopher")
	defer game.Destroy()
	game.Run(game)
}
