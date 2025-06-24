package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	playerSpeed  = 5
	bulletSpeed  = 7
	asteroidSpeed = 7
)

type Game struct {
	player      Player
	bullets     []Bullet
	asteroids   []Asteroid
	gameOver    bool
	score       int
	spawnTimer  int
}

type Player struct {
	x      float64
	y      float64
	width  float64
	height float64
}

type Bullet struct {
	x      float64
	y      float64
	active bool
}

type Asteroid struct {
	x      float64
	y      float64
	width  float64
	height float64
	active bool
}

func (g *Game) Update() error {
	if g.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.reset()
		}
		return nil
	}

	// Player movement
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.player.x > 0 {
		g.player.x -= playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && g.player.x < screenWidth-g.player.width {
		g.player.x += playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) && g.player.y > 0 {
		g.player.y -= playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && g.player.y < screenHeight-g.player.height {
		g.player.y += playerSpeed
	}

	// Shoot bullets
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.bullets = append(g.bullets, Bullet{
			x:      g.player.x + g.player.width/2 - 2,
			y:      g.player.y,
			active: true,
		})
	}

	// Update bullets
	for i := range g.bullets {
		if g.bullets[i].active {
			g.bullets[i].y -= bulletSpeed
			if g.bullets[i].y < 0 {
				g.bullets[i].active = false
			}
		}
	}

	// Spawn asteroids
	g.spawnTimer++
	if g.spawnTimer >= 60 { // Spawn every second (60 frames)
		g.spawnTimer = 0
		width := float64(rand.Intn(30) + 20)
		g.asteroids = append(g.asteroids, Asteroid{
			x:      float64(rand.Intn(screenWidth - int(width))),
			y:      -width,
			width:  width,
			height: width,
			active: true,
		})
	}

	// Update asteroids
	for i := range g.asteroids {
		if g.asteroids[i].active {
			g.asteroids[i].y += asteroidSpeed
			if g.asteroids[i].y > screenHeight {
				g.asteroids[i].active = false
				g.score++
			}
		}
	}

	// Collision detection: bullets vs asteroids
	for i := range g.bullets {
		if !g.bullets[i].active {
			continue
		}
		for j := range g.asteroids {
			if !g.asteroids[j].active {
				continue
			}
			if isColliding(g.bullets[i].x, g.bullets[i].y, 4, 10,
				g.asteroids[j].x, g.asteroids[j].y, g.asteroids[j].width, g.asteroids[j].height) {
				g.bullets[i].active = false
				g.asteroids[j].active = false
				g.score += 5
			}
		}
	}

	// Collision detection: player vs asteroids
	for i := range g.asteroids {
		if !g.asteroids[i].active {
			continue
		}
		if isColliding(g.player.x, g.player.y, g.player.width, g.player.height,
			g.asteroids[i].x, g.asteroids[i].y, g.asteroids[i].width, g.asteroids[i].height) {
			g.gameOver = true
		}
	}

	// Clean up inactive objects
	g.cleanUpObjects()

	return nil
}

func isColliding(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

func (g *Game) cleanUpObjects() {
	// Clean bullets
	var activeBullets []Bullet
	for _, b := range g.bullets {
		if b.active {
			activeBullets = append(activeBullets, b)
		}
	}
	g.bullets = activeBullets

	// Clean asteroids
	var activeAsteroids []Asteroid
	for _, a := range g.asteroids {
		if a.active {
			activeAsteroids = append(activeAsteroids, a)
		}
	}
	g.asteroids = activeAsteroids
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background
	screen.Fill(color.RGBA{0, 0, 20, 255})

	// Draw player (spaceship)
	ebitenutil.DrawRect(screen, g.player.x, g.player.y, g.player.width, g.player.height, color.RGBA{0, 255, 0, 255})
	// Draw ship's cockpit
	ebitenutil.DrawRect(screen, g.player.x+g.player.width/2-2, g.player.y-5, 4, 5, color.RGBA{255, 255, 0, 255})

	// Draw bullets
	for _, b := range g.bullets {
		if b.active {
			ebitenutil.DrawRect(screen, b.x, b.y, 4, 10, color.RGBA{255, 255, 0, 255})
		}
	}

	// Draw asteroids
	for _, a := range g.asteroids {
		if a.active {
			ebitenutil.DrawRect(screen, a.x, a.y, a.width, a.height, color.RGBA{150, 75, 0, 255})
		}
	}

	// Draw score
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", g.score), 10, 10)

	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, "GAME OVER - Press R to restart", screenWidth/2-100, screenHeight/2)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) reset() {
	g.player = Player{
		x:      screenWidth/2 - 15,
		y:      screenHeight - 40,
		width:  30,
		height: 30,
	}
	g.bullets = make([]Bullet, 0)
	g.asteroids = make([]Asteroid, 0)
	g.gameOver = false
	g.score = 0
	g.spawnTimer = 0
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game := &Game{}
	game.reset()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Space Dodger (Linux)")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
