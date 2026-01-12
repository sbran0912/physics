package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	screenWidth := int32(1200)
	screenHeight := int32(1000)
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(screenWidth, screenHeight, "Physics - raylib-go")
	defer rl.CloseWindow()

	a := CreateBox(500, 0, 150, 150)
	b := CreateBox(400, 400, 600, 50)
	a.Rotate(0.9)
	b.Rotate(0.5)
	//a.basic.velocity = rl.Vector2{0, 2}
	//b.basic.velocity = rl.Vector2{0, 0}
	b.basic.mass = math.MaxFloat32
	b.basic.inertia = math.MaxFloat32

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkBlue)

		//Gravity anwenden
		a.ApplyForce(rl.Vector2{0, 0.1}, 0)
		//a.basic.velocity = rl.Vector2ClampValue(a.basic.velocity, 0, 10)
		ok, mtv, contacts := detectCollPoly(&a, &b)
		if ok {
			fmt.Println("Kollision!!")
			forceA, angForcA, forceB, angForceB := resolveCollPoly(&a, &b, contacts, mtv)
			a.ApplyForce(forceA, angForcA)
			b.ApplyForce(forceB, angForceB)
		}
		fmt.Println("Acceleration A:", a.basic.accel, a.basic.angAccel)
		a.Update()
		b.Update()
		fmt.Println("Velocity A:", a.basic.velocity, a.basic.angVelocity)
		a.Draw(rl.Red, 3)
		b.Draw(rl.White, 2)

		rl.DrawLineEx(b.basic.location, rl.Vector2Add(b.basic.location, mtv), 3, rl.Green)
		for _, c := range contacts {
			rl.DrawCircleV(c, 5, rl.Black)
		}

		rl.EndDrawing()
	}
}

type ShapeType int

const (
	CircleShape  ShapeType = iota
	PolygonShape ShapeType = iota
)

type Shape interface {
	Update()
	Draw()
	Rotate()
	FollowMouse()
	ApplyForce()
}

type BasicShape struct {
	typ         ShapeType
	location    rl.Vector2
	velocity    rl.Vector2
	angVelocity float32
	accel       rl.Vector2
	angAccel    float32
	mass        float32
	inertia     float32
}

type Polygon struct {
	basic         BasicShape
	vertices      []rl.Vector2
	verticesCount int
}

type Circle struct {
	basic       BasicShape
	radius      int
	orientation rl.Vector2
}

func CreateBox(x float32, y float32, w float32, h float32) Polygon {
	// 1. Masse basierend auf Fläche berechnen
	mass := w * h

	// 2. Trägheitsmoment für ein Rechteck (um den Schwerpunkt)
	inertia := mass * (w*w + h*h)
	vertices := []rl.Vector2{{x, y}, {x + w, y}, {x + w, y + h}, {x, y + h}}

	box := Polygon{
		basic: BasicShape{
			typ:      PolygonShape,
			location: rl.Vector2{x + w/2, y + h/2},
			mass:     mass,
			inertia:  inertia,
		},
		vertices:      vertices,
		verticesCount: 4,
	}
	return box
}

func (poly *Polygon) ApplyForce(force rl.Vector2, angForce float32) {
	poly.basic.accel = rl.Vector2Add(poly.basic.accel, force)
	poly.basic.angAccel += angForce
}

func (poly *Polygon) Update() {
	poly.basic.velocity = rl.Vector2Add(poly.basic.velocity, poly.basic.accel)
	poly.basic.angVelocity += poly.basic.angAccel

	poly.basic.location = rl.Vector2Add(poly.basic.location, poly.basic.velocity)
	for i := range poly.verticesCount {
		poly.vertices[i] = rl.Vector2Add(poly.vertices[i], poly.basic.velocity)
	}
	poly.Rotate(poly.basic.angVelocity)
	poly.basic.accel = rl.Vector2{0, 0}
	poly.basic.angAccel = 0
}

func (poly Polygon) Draw(c rl.Color, thick float32) {
	for i := range poly.verticesCount {
		rl.DrawLineEx(poly.vertices[i], poly.vertices[(i+1)%poly.verticesCount], thick, c)
	}
	rl.DrawCircleV(poly.basic.location, 5, c)
}

func (poly *Polygon) Rotate(angle float32) {
	for i := range poly.verticesCount {
		relativePos := rl.Vector2Subtract(poly.vertices[i], poly.basic.location)
		rotatedPos := rl.Vector2Rotate(relativePos, angle)
		poly.vertices[i] = rl.Vector2Add(rotatedPos, poly.basic.location)
	}
}

/*
func (poly *Polygon) FollowMouse() {
	mouse := rl.GetMousePosition()
	delta := rl.Vector2Subtract(mouse, poly.basic.location)

	for i := range poly.verticesCount {
		poly.vertices[i] = rl.Vector2Add(poly.vertices[i], delta)
	}
	poly.basic.location = mouse
}
*/

func (poly *Polygon) ResetPos(delta rl.Vector2) {
	for i := range poly.verticesCount {
		poly.vertices[i] = rl.Vector2Add(poly.vertices[i], delta)
	}
	poly.basic.location = rl.Vector2Add(poly.basic.location, delta)
}

// Projektionsbereich eines Polygons auf einer Achse
func projectPolygon(poly Polygon, axis rl.Vector2) (min, max float32) {
	min = rl.Vector2DotProduct(poly.vertices[0], axis)
	max = min
	for _, vertex := range poly.vertices[1:] {
		proj := rl.Vector2DotProduct(vertex, axis)
		if proj < min {
			min = proj
		}
		if proj > max {
			max = proj
		}
	}
	return min, max
}

// Alle Normalenachsen eines Polygons
func getAxes(poly Polygon) []rl.Vector2 {
	axes := []rl.Vector2{}
	for i := range poly.verticesCount {
		edge := rl.Vector2Normalize(rl.Vector2Subtract(poly.vertices[(i+1)%poly.verticesCount], poly.vertices[i]))
		axes = append(axes, rl.Vector2{-edge.Y, edge.X})
	}
	return axes
}

// Prüft ob ein Punkt in einem Polygon liegt
func pointInPolygon(point rl.Vector2, poly *Polygon) bool {
	for i := range poly.verticesCount {
		edge := rl.Vector2Subtract(poly.vertices[(i+1)%poly.verticesCount], poly.vertices[i])
		toPoint := rl.Vector2Subtract(point, poly.vertices[i])
		if rl.Vector2DotProduct(edge, toPoint) < 0 {
			return false
		}
	}
	return true
}

// Findet alle Kontaktpunkte zwischen zwei Polygonen
func findContactPoints(polyA, polyB *Polygon) []rl.Vector2 {
	contacts := []rl.Vector2{}

	// 1. Checke, welche Ecken von A in B liegen
	for _, p := range polyA.vertices {
		if pointInPolygon(p, polyB) {
			contacts = append(contacts, p)
		}
	}

	// 2. Checke, welche Ecken von B in A liegen
	for _, p := range polyB.vertices {
		if pointInPolygon(p, polyA) {
			contacts = append(contacts, p)
		}
	}

	return contacts
}

// Findet die Referenzkante für die Kontaktpunktberechnung
func findReferenceEdge(poly *Polygon, normal rl.Vector2) (p1, p2 rl.Vector2) {
	bestDot := float32(-math.MaxFloat32)
	for i := range poly.verticesCount {
		a := poly.vertices[i]
		b := poly.vertices[(i+1)%poly.verticesCount]

		edge := rl.Vector2Normalize(rl.Vector2Subtract(b, a))
		edgeNormal := rl.Vector2{-edge.Y, edge.X}

		d := rl.Vector2DotProduct(edgeNormal, normal)
		if d > bestDot {
			bestDot = d
			p1 = a
			p2 = b
		}
	}
	return p1, p2
}

// Projiziert einen Punkt auf eine Kante
func projectPointOntoEdge(p, a, b rl.Vector2) rl.Vector2 {
	ab := rl.Vector2Subtract(b, a)
	t := rl.Vector2DotProduct(rl.Vector2Subtract(p, a), ab) /
		rl.Vector2DotProduct(ab, ab)

	// Clamp auf Segment
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	return rl.Vector2Add(a, rl.Vector2Scale(ab, t))
}

// Überträgt Kontaktpunkte auf die Referenzkante von Polygon A
func transferContactsToA(contacts []rl.Vector2, polyA *Polygon, normal rl.Vector2) []rl.Vector2 {
	refEdge_start, refEdge_end := findReferenceEdge(polyA, normal)

	projected := []rl.Vector2{}
	for _, c := range contacts {
		p := projectPointOntoEdge(c, refEdge_start, refEdge_end)
		projected = append(projected, p)
	}
	return projected
}

// SAT-Kollisionstest
func detectCollPoly(polyA, polyB *Polygon) (isColliding bool, mtv rl.Vector2, contacts []rl.Vector2) {
	smallestOverlap := float32(math.MaxFloat32)
	var smallestAxis rl.Vector2

	// Alle Achsen beider Polygone sammeln
	axes := append(getAxes(*polyA), getAxes(*polyB)...)

	for _, axis := range axes {
		minA, maxA := projectPolygon(*polyA, axis)
		minB, maxB := projectPolygon(*polyB, axis)

		overlap := float32(math.Min(float64(maxA), float64(maxB))) - float32(math.Max(float64(minA), float64(minB)))
		if overlap <= 0 {
			return false, rl.Vector2{}, nil
		}

		if overlap < smallestOverlap {
			smallestOverlap = overlap
			smallestAxis = axis

			// Richtung korrigieren
			dir := rl.Vector2Subtract(polyB.basic.location, polyA.basic.location)
			if rl.Vector2DotProduct(dir, smallestAxis) < 0 {
				smallestAxis = rl.Vector2Scale(smallestAxis, -1)
			}
		}
	}

	mtv = rl.Vector2Scale(smallestAxis, smallestOverlap)

	rawContacts := findContactPoints(polyA, polyB)
	// Hinweis: smallesAxis ist bereits normalisiert!
	contacts = transferContactsToA(rawContacts, polyA, rl.Vector2Scale(smallestAxis, -1))
	return true, mtv, contacts
}

func resolveCollPoly(
	polyA *Polygon,
	polyB *Polygon,
	contacts []rl.Vector2,
	mtv rl.Vector2) (forceA rl.Vector2, angForceA float32, forceB rl.Vector2, angForceB float32) {

	//Objekte um mtv zurücksetzen
	if polyA.basic.mass < math.MaxFloat32 {
		polyA.ResetPos(rl.Vector2Scale(mtv, -0.5))
	} else {
		polyB.ResetPos(rl.Vector2Scale(mtv, 0.5))
	}
	if polyB.basic.mass < math.MaxFloat32 {
		polyB.ResetPos(rl.Vector2Scale(mtv, 0.5))
	} else {
		polyA.ResetPos(rl.Vector2Scale(mtv, -0.5))
	}

	// danach muss mtv normalisiert werden
	mtv = rl.Vector2Normalize(mtv)

	//finde bei zwei Kontaktpunkten den Kollisionspunkt in der Mitte
	collisionpoint := contacts[0]
	if len(contacts) > 1 {
		collisionpoint = rl.Vector2Scale(rl.Vector2Add(collisionpoint, contacts[1]), 0.5)
	}

	// Linie von A.location zu Kollisionspunkt
	rAP := rl.Vector2Subtract(collisionpoint, polyA.basic.location)
	// Linie von B.location zu Kollisionspunkt
	rBP := rl.Vector2Subtract(collisionpoint, polyB.basic.location)

	rAP_perp := rl.Vector2{-rAP.Y, rAP.X}
	rBP_perp := rl.Vector2{-rBP.Y, rBP.X}
	VtanA := rl.Vector2Scale(rAP_perp, polyA.basic.angVelocity)
	VtanB := rl.Vector2Scale(rBP_perp, polyB.basic.angVelocity)
	VgesamtA := rl.Vector2Add(polyA.basic.velocity, VtanA)
	VgesamtB := rl.Vector2Add(polyB.basic.velocity, VtanB)
	velocity_AB := rl.Vector2Subtract(VgesamtA, VgesamtB)
	fmt.Println("VtanA:", VtanA)
	fmt.Println("Gesamt velocity AB mit Tangentialgeschwindigkeit:", velocity_AB)

	if rl.Vector2DotProduct(velocity_AB, rl.Vector2Scale(mtv, -1)) < 0 { // wenn negativ, dann auf Kollisionskurs
		e := float32(0.4) //inelastischer Stoß
		j_denominator := rl.Vector2DotProduct(rl.Vector2Scale(velocity_AB, -(1+e)), mtv)
		j_divLinear := rl.Vector2DotProduct(mtv, rl.Vector2Scale(mtv, (1/polyA.basic.mass+1/polyB.basic.mass)))
		j_divAngular := float32(math.Pow(float64(rl.Vector2DotProduct(rAP_perp, mtv)), 2))/polyA.basic.inertia + float32(math.Pow(float64(rl.Vector2DotProduct(rBP_perp, mtv)), 2))/polyB.basic.inertia
		j := j_denominator / (j_divLinear + j_divAngular)
		// Grundlage für Friction berechnen (t)
		t := rl.Vector2{-mtv.Y, mtv.X}
		t_scalarprodukt := rl.Vector2DotProduct(velocity_AB, t)
		t = rl.Vector2Normalize(rl.Vector2Scale(t, t_scalarprodukt))

		friction := float32(-0.45)
		//kalkuliere Kraft polyA
		forceA = rl.Vector2Add(rl.Vector2Scale(mtv, (j/polyA.basic.mass)), rl.Vector2Scale(t, (friction*-j/polyA.basic.mass)))
		angForceA = rl.Vector2DotProduct(rAP_perp, rl.Vector2Add(rl.Vector2Scale(mtv, j/polyA.basic.inertia), rl.Vector2Scale(t, friction*-j/polyA.basic.inertia)))
		//kalkuliere Kraft polyB
		forceB = rl.Vector2Add(rl.Vector2Scale(mtv, (-j/polyB.basic.mass)), rl.Vector2Scale(t, (friction*j/polyB.basic.mass)))
		angForceB = rl.Vector2DotProduct(rAP_perp, rl.Vector2Add(rl.Vector2Scale(mtv, j/polyB.basic.inertia), rl.Vector2Scale(t, friction*-j/polyB.basic.inertia)))
	} else {
		forceA = rl.Vector2{0, 0}
		angForceA = 0.0
		forceB = rl.Vector2{0, 0}
		angForceB = 0.0
	}
	return forceA, angForceA, forceB, angForceB
}
