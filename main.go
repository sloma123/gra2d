package main

import (
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var assets embed.FS

const (
	screenW = 630
	screenH = 480

	playerW = 30
	playerH = 30

	coinX = 400.0
	coinY = 397.0

	obstacleX = 250.0
	obstacleY = 396.0

	groundY = 397.0
	gravity = 0.4

	stepSpeed   = 3.0
	jumpImpulse = -10.0
)

type Game struct {
	player *Player
	scene  *Scene
	win    bool
	lose   bool
	agent  *Agent
	mode   string // "manual" lub "rl"
}

type Scene struct {
	background *ebiten.Image
	ground     *ebiten.Image
	coin       *ebiten.Image
	obstacle   *ebiten.Image
}

type Player struct {
	Image     *ebiten.Image
	x, y      float64
	speed     float64
	velocityY float64
	isJumping bool
}

type Agent struct {
	Q        map[int][4]float64 // Q[s][a]
	epsilon  float64
	alpha    float64
	gamma    float64
	xBins    int
	xBinSize float64
}

func limit(x, lo, hi int) int {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func distance(a, b float64) float64 {
	d := a - b
	if d < 0 {
		return -d
	}
	return d
}

func checkCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

func DrawScene(img *ebiten.Image, screen *ebiten.Image, scaleX float64, scaleY float64, x int, y int, translate bool) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	if translate {
		op.GeoM.Translate(float64(x), float64(y))
	}
	screen.DrawImage(img, op)
}

func (g *Game) Update() error {
	if g.win || g.lose {
		return nil
	}

	// 🔹 TRYB RĘCZNY
	if g.mode == "manual" {
		speed := 3.0

		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			g.player.speed = -speed
		} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			g.player.speed = speed
		} else {
			g.player.speed = 0
		}

		if ebiten.IsKeyPressed(ebiten.KeySpace) && !g.player.isJumping {
			g.player.velocityY = jumpImpulse
			g.player.isJumping = true
		}
	} else {
		// 🔹 TRYB RL
		s := g.agent.stateKey(g.player.x, g.player.isJumping)
		action, _ := g.agent.bestAction(s)
		switch action {
		case 0:
			g.player.speed = -stepSpeed
		case 1:
			g.player.speed = stepSpeed
		case 2:
			if !g.player.isJumping {
				g.player.velocityY = jumpImpulse
				g.player.isJumping = true
			}
			g.player.speed = 0
		case 3:
			g.player.speed = 0
		}
	}

	// FIZYKA
	g.player.velocityY += gravity
	g.player.y += g.player.velocityY
	g.player.x += g.player.speed

	if g.player.y >= groundY {
		g.player.y = groundY
		g.player.velocityY = 0
		g.player.isJumping = false
	}
	if g.player.x < 0 {
		g.player.x = 0
	}
	if g.player.x > float64(screenW-playerW) {
		g.player.x = float64(screenW - playerW)
	}

	// KOLIZJE
	if checkCollision(g.player.x, g.player.y, playerW, playerH, coinX, coinY, 30, 30) {
		g.win = true
		log.Println("Wygrana")
	}
	if checkCollision(g.player.x, g.player.y, playerW, playerH, obstacleX, obstacleY, 30, 30) {
		g.lose = true
		log.Println("Przegrana")
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	DrawScene(g.scene.background, screen, float64(screen.Bounds().Dx())/float64(g.scene.background.Bounds().Dx()), float64(screen.Bounds().Dy())/float64(g.scene.background.Bounds().Dy()), 0, 0, false)
	for i := 0; i < 18; i++ {
		x := i * 35
		DrawScene(g.scene.ground, screen, 2, 2, x, 445, true)
	}
	DrawScene(g.player.Image, screen, 2, 2, int(g.player.x), int(g.player.y), true)
	DrawScene(g.scene.coin, screen, 2, 2, int(coinX), int(coinY), true)
	DrawScene(g.scene.obstacle, screen, 2, 2, int(obstacleX), int(obstacleY), true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenW, screenH
}

func main() {

	var modeChoice int
	fmt.Print("Wybierz tryb (1 = manualny, 2 = RL): ")
	fmt.Scan(&modeChoice)

	mode := "manual"
	if modeChoice == 2 {
		mode = "rl"
	}

	agent := &Agent{
		Q:        make(map[int][4]float64),
		epsilon:  1.0,
		alpha:    0.3,
		gamma:    0.96,
		xBins:    64,
		xBinSize: float64(screenW) / 64.0,
	}

	if mode == "rl" {
		log.Println("Trening agenta...")
		agent.Train(10000, 240)
		log.Println("Trening zakończony.")
	}

	playerImg := mustLoadPNG("assets/player.png")
	backgroundImg := mustLoadPNG("assets/background.png")
	groundImg := mustLoadPNG("assets/ground.png")
	coinImg := mustLoadPNG("assets/coin.png")
	obstacleImg := mustLoadPNG("assets/obstacle.png")

	player := &Player{Image: playerImg, x: 0, y: groundY}
	scene := &Scene{background: backgroundImg, ground: groundImg, coin: coinImg, obstacle: obstacleImg}
	game := &Game{player: player, scene: scene, agent: agent, mode: mode}

	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("RL Coin Grabber (Manual / RL)")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func mustLoadPNG(path string) *ebiten.Image {
	f, err := assets.Open(path)
	if err != nil {
		log.Fatalf("Nie udało się wczytać %s: %v", path, err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}
