package main

import (
	"collision/lib"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1200)
	screenHeight := int32(1000)
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(screenWidth, screenHeight, "Physics - raylib-go")
	defer rl.CloseWindow()

	//slice für lib.Shape erstellen
	shapes := []lib.Shape{}
	shapes = append(shapes, lib.CreatePolygon(900, 80, 100, 70, false))
	shapes = append(shapes, lib.CreatePolygon(650, 197, 300, 100, false))
	shapes = append(shapes, lib.CreatePolygon(350, 500, 1000, 50, true))
	shapes = append(shapes, lib.CreatePolygon(800, 301, 150, 150, false))

	shapes[2].Rotate(-0.08)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkBlue)

		// slice p iterieren
		for i := range shapes {
			//Gravity anwenden
			shapes[i].ApplyGravity(rl.Vector2{0, 0.4})
		}

		//Kollision für alle elemente des slice prüfen
		for i := 0; i < len(shapes)-1; i++ {
			for j := i + 1; j < len(shapes); j++ {
				ok, mtv, contacts := lib.DetectCollision(shapes[i], shapes[j])
				//fmt.Println(ok, mtv, contacts)
				if ok {
					//fmt.Println("Kollision!!")
					forceA, angForcA, forceB, angForceB := lib.ResolveCollision(shapes[i], shapes[j], contacts, mtv)
					shapes[i].ApplyForce(forceA, angForcA)
					shapes[j].ApplyForce(forceB, angForceB)
				}
				//zeichne Kontaktpunkte
				for _, contact := range contacts {
					rl.DrawCircleV(contact, 5, rl.Black)
				}

			}
		}
		for i := range shapes {
			shapes[i].Update()
			shapes[i].Draw(rl.Beige, 3)
		}
		rl.EndDrawing()
	}
}
