package main

import (
	"embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// dyrektywa kompilatora do osadzenia plików z assets
//
//go:embed assets/*
var assets embed.FS

// FS - file system

type Game struct {
	player *Player
}

type Player struct {
	Image *ebiten.Image
}

// wykonuje się co klatkę
func (g *Game) Update() error {
	return nil
}

// wykonuje się co klatkę
func (g *Game) Draw(screen *ebiten.Image) {
	if g.player != nil && g.player.Image != nil {
		screen.DrawImage(g.player.Image, nil)
	}
}

// wykonuje się na początku i przy zmianie rozmiaru okna
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
	playerImg := ebiten.NewImageFromImage(img)

	player := &Player{Image: playerImg}

	game := &Game{player: player}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Gra 2d")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
