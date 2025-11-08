package main

import (
	"embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assets embed.FS

var playerImage *ebiten.Image

// Game struktura gry
type Game struct{}

// Update – logika gry
func (g *Game) Update() error {
	return nil
}

// Draw – rysowanie na ekranie
func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(playerImage, nil)
}

// Layout – rozmiar okna
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 320, 240
}

func main() {
	// Wczytaj osadzony plik PNG
	f, err := assets.Open("assets/player.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	playerImage = ebiten.NewImageFromImage(img)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Osadzony obraz PNG w grze")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
