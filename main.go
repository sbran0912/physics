package main

import (
	"physics/lib"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1400)
	screenHeight := int32(1000)
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(screenWidth, screenHeight, "Physics - raylib-go")
	defer rl.CloseWindow()

	shapes := []lib.Shape{}

	// Boden (unbeweglich)
	shapes = append(shapes, lib.CreatePolygon(350, 700, 1000, 50, true))

	// Mauer aus gestapelten Blöcken (4 Blöcke übereinander)
	blockWidth := float32(80)
	blockHeight := float32(50)
	startX := float32(800) - blockWidth/2 // zentriert bei x=800
	baseY := float32(700) - blockHeight   // Bodenoberkante = 700, erster Block sitzt darauf

	for i := range 4 {
		y := baseY - float32(i)*blockHeight
		block := lib.CreatePolygon(startX, y, blockWidth, blockHeight, false)
		shapes = append(shapes, block)
	}

	// Kreis, der auf die Mauer geschossen wird
	circle := lib.CreateCircle(400, 600, 40, false)
	shapes = append(shapes, circle)

	// Einmaliger kräftiger Impuls nach rechts
	circle.ApplyForce(rl.Vector2{20, 0}, 0)

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
				// Kontaktpunkte zeichnen (optional)
				for _, contact := range contacts {
					rl.DrawCircleV(contact, 5, rl.Black)
				}
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
