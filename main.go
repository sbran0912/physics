package main

import (
	"physics/lib"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1400)
	screenHeight := int32(700)
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(screenWidth, screenHeight, "Physics - raylib-go")
	defer rl.CloseWindow()

	// Seed für Zufallszahlen

	shapes := []lib.Shape{}

	// Boden (unbeweglich)
	shapes = append(shapes, lib.CreatePolygon(350, 500, 1000, 50, true))

	// Mauer aus 4 gestapelten Blöcken mit zufälliger Breite und Höhe
	baseX := float32(800)
	baseY := float32(500)
	currentY := baseY

	for i := 0; i < 4; i++ {
		// Breite zwischen 60 und 100
		blockW := 60 + lib.RandomFloat(0, 40)
		// Höhe zwischen 35 und 70
		blockH := 35 + lib.RandomFloat(0, 35)
		// Zentriert
		blockX := baseX - blockW/2

		shapes = append(shapes, lib.CreatePolygon(blockX, currentY-blockH, blockW, blockH, false))
		currentY -= blockH // Für nächsten Block
	}

	// Zwei Kreise
	circle1 := lib.CreateCircle(400, 400, 42, false)
	circle2 := lib.CreateCircle(300, 300, 35, false)
	shapes = append(shapes, circle1, circle2)

	// Einmaliger Impuls nach rechts
	circle1.ApplyForce(rl.Vector2{10, 0}, 0)
	circle2.ApplyForce(rl.Vector2{12, 0}, 0)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkBlue)

		// Gravität auf alle Objekte anwenden
		for i := range shapes {
			shapes[i].ApplyGravity(rl.Vector2{0, 0.4})
			shapes[i].ClearState()
		}

		// Kollisionen prüfen und auflösen
		for i := 0; i < len(shapes)-1; i++ {
			for j := i + 1; j < len(shapes); j++ {
				ok, mtv, contacts := lib.DetectCollision(shapes[i], shapes[j])
				if ok {
					forceA, angForceA, forceB, angForceB := lib.ResolveCollision(shapes[i], shapes[j], contacts, mtv)
					shapes[i].ApplyForce(forceA, angForceA)
					shapes[j].ApplyForce(forceB, angForceB)
				}
				// Kontaktpunkte zeichnen entfernt (wie in TypeScript)
			}
		}

		// Alle Objekte aktualisieren und zeichnen
		for i := range shapes {
			shapes[i].Update()
			shapes[i].Draw(rl.Beige, 3)
		}
		rl.EndDrawing()
	}
}
