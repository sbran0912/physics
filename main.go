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
	p := []lib.Shape{}
	p = append(p, lib.CreatePolygon(490, 120, 90, 70, false))
	p = append(p, lib.CreatePolygon(350, 200, 150, 100, false))
	p = append(p, lib.CreatePolygon(400, 350, 150, 150, false))
	p = append(p, lib.CreatePolygon(350, 520, 700, 50, true))
	p[0].Rotate(0.1)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkBlue)

		// slice p iterieren
		for i := range p {
			//Gravity anwenden
			p[i].ApplyGravity(rl.Vector2{0, 0.2})
		}

		//Kollision für alle elemente des slice prüfen
		for i := 0; i < len(p)-1; i++ {
			for j := i + 1; j < len(p); j++ {
				ok, mtv, contacts := lib.DetectCollision(p[i], p[j])
				//fmt.Println(ok, mtv, contacts)
				if ok {
					//fmt.Println("Kollision!!")
					forceA, angForcA, forceB, angForceB := lib.ResolveCollision(p[i], p[j], contacts, mtv)
					p[i].ApplyForce(forceA, angForcA)
					p[j].ApplyForce(forceB, angForceB)
				}
				for _, contact := range contacts {
					rl.DrawCircleV(contact, 5, rl.Black)
				}

			}
		}

		for i := range p {
			p[i].Update()
			if i == 3 {
				p[i].Draw(rl.Red, 3)
			} else {
				p[i].Draw(rl.Beige, 3)
			}
		}

		rl.EndDrawing()
	}
}
