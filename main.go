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

	a := lib.CreatePolygon(600, 100, 150, 150, false)
	b := lib.CreatePolygon(350, 400, 700, 50, true)
	a.Rotate(0.3)
	//b.Rotate(-0.3)
	//b.ApplyForce(rl.Vector2{0, 0}, -0.005)
	//a.ApplyForce(rl.Vector2{0, 2}, 0)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkBlue)

		//Gravity anwenden
		a.ApplyGravity(rl.Vector2{0, 0.3})

		//Kollision
		ok, mtv, contacts := lib.DetectCollPoly(&a, &b)
		if ok {
			//fmt.Println("Kollision!!")
			forceA, angForcA, forceB, angForceB := lib.ResolveCollPoly(&a, &b, contacts, mtv)
			a.ApplyForce(forceA, angForcA)
			b.ApplyForce(forceB, angForceB)
		}

		a.Update()
		//b.Update()
		a.Draw(rl.Red, 3)
		b.Draw(rl.White, 2)

		//rl.DrawLineEx(b.basic.location, rl.Vector2Add(b.basic.location, mtv), 3, rl.Green)
		for _, c := range contacts {
			rl.DrawCircleV(c, 5, rl.Black)
		}

		rl.EndDrawing()
	}
}
