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

	a := lib.CreatePolygon(500, 0, 150, 150, false)
	b := lib.CreatePolygon(400, 400, 600, 50, true)
	a.Rotate(0.9)
	//b.Rotate(0.2)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkBlue)

		//Gravity anwenden
		a.ApplyForce(rl.Vector2{0, 0.1}, 0)
		//a.basic.velocity = rl.Vector2ClampValue(a.basic.velocity, 0, 10)
		ok, mtv, contacts := lib.DetectCollPoly(&a, &b)
		if ok {
			//fmt.Println("Kollision!!")
			forceA, angForcA, forceB, angForceB := lib.ResolveCollPoly(&a, &b, contacts, mtv)
			a.ApplyForce(forceA, angForcA)
			b.ApplyForce(forceB, angForceB)
		}
		//fmt.Println("Acceleration A:", a.basic.accel, a.basic.angAccel)
		a.Update()
		b.Update()
		//fmt.Println("Velocity A:", a.basic.velocity, a.basic.angVelocity)
		a.Draw(rl.Red, 3)
		b.Draw(rl.White, 2)

		//rl.DrawLineEx(b.basic.location, rl.Vector2Add(b.basic.location, mtv), 3, rl.Green)
		for _, c := range contacts {
			rl.DrawCircleV(c, 5, rl.Black)
		}

		rl.EndDrawing()
	}
}
