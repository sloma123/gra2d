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
	scene  *Scene
	win    bool
	death  bool
}

type Scene struct {
	background *ebiten.Image
	ground     *ebiten.Image
	coin       *ebiten.Image
}

type Player struct {
	Image     *ebiten.Image
	x, y      float64
	speed     float64
	velocityY float64
	gravity   float64
	isJumping bool
	groundY   float64
}

// wykonuje się co klatkę
func (g *Game) Update() error {
	speed := 2.0

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.speed = -speed
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.speed = speed
	} else {
		g.player.speed = 0
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.player.isJumping {
		g.player.velocityY = -10  // nadajemy prędkość w górę
		g.player.isJumping = true // gracz jest w powietrzu
	}

	g.player.velocityY += g.player.gravity
	g.player.y += g.player.velocityY

	g.player.x += g.player.speed

	if g.player.y >= g.player.groundY {
		g.player.y = g.player.groundY // ustaw na ziemię
		g.player.velocityY = 0        // zatrzymaj spadanie
		g.player.isJumping = false    // pozwól znowu skakać
	}

	g.win = checkCollision(g.player.x, g.player.y, 30, 30, 400, 397, 30, 30)
	if g.win {
		log.Println("Wygrałeś!")
	}
	return nil
}

// wykonuje się po Update
func (g *Game) Draw(screen *ebiten.Image) {
	// if g.player != nil && g.player.Image != nil {
	// 	screen.DrawImage(g.player.Image, nil)
	// }
	DrawScene(g.scene.background, screen, float64(screen.Bounds().Dx())/float64(g.scene.background.Bounds().Dx()), float64(screen.Bounds().Dy())/float64(g.scene.background.Bounds().Dy()), 0, 0, false)
	for i := 0; i < 18; i++ {
		x := i * 35
		DrawScene(g.scene.ground, screen, 2, 2, x, 445, true)
	}
	DrawScene(g.player.Image, screen, 2, 2, int(g.player.x), int(g.player.y), true)
	DrawScene(g.scene.coin, screen, 2, 2, 400, 397, true)

}

// wykonuje się na początku i przy zmianie rozmiaru okna
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 630, 480
}

func main() {
	// Wczytaj osadzony plik PNG
	f, err := assets.Open("assets/player.png")
	if err != nil {
		log.Fatal("Nie udało się wczytać player.png:", err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	playerImg := ebiten.NewImageFromImage(img)

	bgFile, err := assets.Open("assets/background.png")
	if err != nil {
		log.Fatal("Nie udało się wczytać background.png:", err)
	}
	bgImg, _, err := image.Decode(bgFile)
	if err != nil {
		log.Fatal(err)
	}
	backgroundImg := ebiten.NewImageFromImage(bgImg)

	grFile, err := assets.Open("assets/ground.png")
	if err != nil {
		log.Fatal("Nie udało się wczytać ground.png:", err)
	}
	grImg, _, err := image.Decode(grFile)
	if err != nil {
		log.Fatal(err)
	}
	groundImg := ebiten.NewImageFromImage(grImg)

	coFile, err := assets.Open("assets/coin.png")
	if err != nil {
		log.Fatal("Nie udało się wczytać ground.png:", err)
	}
	coImg, _, err := image.Decode(coFile)
	if err != nil {
		log.Fatal(err)
	}
	coinImg := ebiten.NewImageFromImage(coImg)

	player := &Player{Image: playerImg, x: 0, y: 397, speed: 0, velocityY: 0, gravity: 0.5, isJumping: false, groundY: 397}
	scene := &Scene{background: backgroundImg, ground: groundImg, coin: coinImg}
	game := &Game{player: player, scene: scene}

	ebiten.SetWindowSize(630, 480)
	ebiten.SetWindowTitle("Gra 2d")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func DrawScene(img *ebiten.Image, screen *ebiten.Image, scaleX float64, scaleY float64, x int, y int, translate bool) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	if translate {
		op.GeoM.Translate(float64(x), float64(y))
	}
	screen.DrawImage(img, op)

}

func checkCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	if x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2 {
		return true
	}
	return false
}
