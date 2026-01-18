package lib

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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
	ResetPos()
}

type BasicShape struct {
	typ         ShapeType
	location    rl.Vector2
	Velocity    rl.Vector2
	AngVelocity float32
	accel       rl.Vector2
	angAccel    float32
	mass        float32
	inertia     float32
	IsGrounded  bool
}

type Polygon struct {
	Basic    BasicShape
	vertices []rl.Vector2
}

type Circle struct {
	basic       BasicShape
	radius      float32
	orientation rl.Vector2
}

func CreatePolygon(x float32, y float32, w float32, h float32, wall bool) Polygon {
	var mass float32
	var inertia float32

	if !wall {
		mass = w * h
		inertia = mass * (w*w + h*h) / 2
	} else {
		mass = float32(math.MaxFloat32)
		inertia = float32(math.MaxFloat32)
	}

	vertices := []rl.Vector2{{x, y}, {x + w, y}, {x + w, y + h}, {x, y + h}}

	poly := Polygon{
		Basic: BasicShape{
			typ:        PolygonShape,
			location:   rl.Vector2{x + w/2, y + h/2},
			mass:       mass,
			inertia:    inertia,
			IsGrounded: false,
		},
		vertices: vertices,
	}
	return poly
}

func (poly *Polygon) ApplyForce(force rl.Vector2, angForce float32) {
	poly.Basic.accel = rl.Vector2Add(poly.Basic.accel, force)
	poly.Basic.angAccel += angForce
}

func (poly *Polygon) Update() {
	fmt.Println("Polygon Update", "angaccel:", poly.Basic.angAccel, "accel:", poly.Basic.accel)
	poly.Basic.Velocity = rl.Vector2Add(poly.Basic.Velocity, poly.Basic.accel)
	poly.Basic.AngVelocity += poly.Basic.angAccel
	fmt.Println("Polygon Update", "angvel:", poly.Basic.AngVelocity, "vel:", poly.Basic.Velocity)

	//poly.basic.velocity = rl.Vector2Scale(poly.basic.velocity, 0.995)
	poly.Basic.location = rl.Vector2Add(poly.Basic.location, poly.Basic.Velocity)
	for i := range poly.vertices {
		poly.vertices[i] = rl.Vector2Add(poly.vertices[i], poly.Basic.Velocity)
	}
	//poly.basic.angVelocity *= 0.995
	poly.Rotate(poly.Basic.AngVelocity)

	poly.Basic.accel = rl.Vector2{0, 0}
	poly.Basic.angAccel = 0
}

func (poly *Polygon) Draw(c rl.Color, thick float32) {
	n := len(poly.vertices)
	for i := range n {
		rl.DrawLineEx(poly.vertices[i], poly.vertices[(i+1)%n], thick, c)
	}
	rl.DrawCircleV(poly.Basic.location, 5, c)
}

func (poly *Polygon) Rotate(angle float32) {
	for i := range poly.vertices {
		relativePos := rl.Vector2Subtract(poly.vertices[i], poly.Basic.location)
		rotatedPos := rl.Vector2Rotate(relativePos, angle)
		poly.vertices[i] = rl.Vector2Add(rotatedPos, poly.Basic.location)
	}
}

func (poly *Polygon) ResetPos(delta rl.Vector2) {
	for i := range poly.vertices {
		poly.vertices[i] = rl.Vector2Add(poly.vertices[i], delta)
	}
	poly.Basic.location = rl.Vector2Add(poly.Basic.location, delta)
}

func CreateCircle(x float32, y float32, r float32, wall bool) Circle {
	var mass float32
	var inertia float32

	if !wall {
		mass = r * r / 4
		inertia = r * r * r / 2
	} else {
		mass = float32(math.MaxFloat32)
		inertia = float32(math.MaxFloat32)
	}

	circle := Circle{
		basic: BasicShape{
			typ:      CircleShape,
			location: rl.Vector2{400, 200},
			mass:     mass,
			inertia:  inertia,
		},
		radius:      r,
		orientation: rl.Vector2{r + x, y},
	}
	return circle
}

func (circle *Circle) ApplyForce(force rl.Vector2, angForce float32) {
	circle.basic.accel = rl.Vector2Add(circle.basic.accel, force)
	circle.basic.angAccel += angForce
}

func (circle *Circle) Update() {
	circle.basic.Velocity = rl.Vector2Add(circle.basic.Velocity, circle.basic.accel)
	circle.basic.accel = rl.Vector2{0.0, 0.0}
	circle.basic.AngVelocity += circle.basic.angAccel
	circle.basic.angAccel = 0.0

	circle.basic.location = rl.Vector2Add(circle.basic.location, circle.basic.Velocity)
	circle.orientation = rl.Vector2Add(circle.orientation, circle.basic.Velocity)
	circle.Rotate(circle.basic.AngVelocity)
}

func (circle *Circle) draw(thick float32, c rl.Color) {
	rl.DrawRing(circle.basic.location, circle.radius-3, circle.radius, 0, 360, 1, c)
	rl.DrawCircleV(circle.basic.location, 3, c)
	rl.DrawLineEx(circle.basic.location, circle.orientation, thick, c)
}

func (circle *Circle) Rotate(angle float32) {
	relativePos := rl.Vector2Subtract(circle.orientation, circle.basic.location)
	rotatedPos := rl.Vector2Rotate(relativePos, angle)
	circle.orientation = rl.Vector2Add(rotatedPos, circle.basic.location)
}

func (circle *Circle) ResetPos(delta rl.Vector2) {

	circle.basic.location = rl.Vector2Add(circle.basic.location, delta)
	circle.orientation = rl.Vector2Add(circle.orientation, delta)

}

// Projektionsbereich eines Polygons auf einer Achse
func projectPolygon(vertices []rl.Vector2, axis rl.Vector2) (min, max float32) {
	if len(vertices) == 0 {
		return 0, 0
	}
	min = rl.Vector2DotProduct(vertices[0], axis)
	max = min
	for _, vertex := range vertices[1:] {
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
func getAxes(vertices []rl.Vector2) []rl.Vector2 {
	axes := []rl.Vector2{}
	n := len(vertices)
	if n == 0 {
		return axes
	}
	for i := 0; i < n; i++ {
		a := vertices[i]
		b := vertices[(i+1)%n]
		edge := rl.Vector2Normalize(rl.Vector2Subtract(b, a))
		axes = append(axes, rl.Vector2{-edge.Y, edge.X})
	}
	return axes
}

// Prüft ob ein Punkt in einem Polygon liegt (convex)
func pointInPolygon(point rl.Vector2, vertices []rl.Vector2) bool {
	n := len(vertices)
	if n == 0 {
		return false
	}
	for i := 0; i < n; i++ {
		a := vertices[i]
		b := vertices[(i+1)%n]
		edge := rl.Vector2Subtract(b, a)
		toPoint := rl.Vector2Subtract(point, a)
		if rl.Vector2DotProduct(edge, toPoint) < 0 {
			return false
		}
	}
	return true
}

// Findet alle Kontaktpunkte zwischen zwei Polygonen (Ecken, die im anderen liegen)
func findContactPoints(verticesA, verticesB []rl.Vector2) []rl.Vector2 {
	contacts := []rl.Vector2{}

	// 1. Checke, welche Ecken von A in B liegen
	for _, p := range verticesA {
		if pointInPolygon(p, verticesB) {
			contacts = append(contacts, p)
		}
	}

	// 2. Checke, welche Ecken von B in A liegen
	for _, p := range verticesB {
		if pointInPolygon(p, verticesA) {
			contacts = append(contacts, p)
		}
	}

	return contacts
}

// Findet die Referenzkante für die Kontaktpunktberechnung
func findReferenceEdge(vertices []rl.Vector2, normal rl.Vector2) (p1, p2 rl.Vector2) {
	bestDot := float32(-math.MaxFloat32)
	n := len(vertices)
	if n == 0 {
		return rl.Vector2{}, rl.Vector2{}
	}
	for i := 0; i < n; i++ {
		a := vertices[i]
		b := vertices[(i+1)%n]

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
func transferContactsToA(contacts []rl.Vector2, verticesA []rl.Vector2, normal rl.Vector2) []rl.Vector2 {
	refEdge_start, refEdge_end := findReferenceEdge(verticesA, normal)

	projected := []rl.Vector2{}
	for _, c := range contacts {
		p := projectPointOntoEdge(c, refEdge_start, refEdge_end)
		projected = append(projected, p)
	}
	return projected
}

// SAT-Kollisionstest
func DetectCollPoly(polyA, polyB *Polygon) (isColliding bool, mtv rl.Vector2, contacts []rl.Vector2) {
	smallestOverlap := float32(math.MaxFloat32)
	var smallestAxis rl.Vector2

	// Alle Achsen beider Polygone sammeln
	axes := append(getAxes(polyA.vertices), getAxes(polyB.vertices)...)

	for _, axis := range axes {
		minA, maxA := projectPolygon(polyA.vertices, axis)
		minB, maxB := projectPolygon(polyB.vertices, axis)

		overlap := float32(math.Min(float64(maxA), float64(maxB))) - float32(math.Max(float64(minA), float64(minB)))
		if overlap <= 0 {
			return false, rl.Vector2{}, nil
		}

		if overlap < smallestOverlap {
			smallestOverlap = overlap
			smallestAxis = axis

			// Richtung korrigieren
			dir := rl.Vector2Subtract(polyB.Basic.location, polyA.Basic.location)
			if rl.Vector2DotProduct(dir, smallestAxis) < 0 {
				smallestAxis = rl.Vector2Scale(smallestAxis, -1)
			}
		}
	}

	mtv = rl.Vector2Scale(smallestAxis, smallestOverlap)

	rawContacts := findContactPoints(polyA.vertices, polyB.vertices)
	// Hinweis: smallestAxis ist bereits normalisiert!
	contacts = transferContactsToA(rawContacts, polyA.vertices, rl.Vector2Scale(smallestAxis, -1))
	return true, mtv, contacts
}

func ResolveCollPoly(
	polyA *Polygon,
	polyB *Polygon,
	contacts []rl.Vector2,
	mtv rl.Vector2) (forceA rl.Vector2, angForceA float32, forceB rl.Vector2, angForceB float32) {

	//Objekte um mtv zurücksetzen
	if polyA.Basic.mass < math.MaxFloat32 {
		polyA.ResetPos(rl.Vector2Scale(mtv, -0.5))
		fmt.Println("PolyA position zurückgesetzt")
	} else {
		polyB.ResetPos(rl.Vector2Scale(mtv, 0.5))
	}
	if polyB.Basic.mass < math.MaxFloat32 {
		polyB.ResetPos(rl.Vector2Scale(mtv, 0.5))
	} else {
		polyA.ResetPos(rl.Vector2Scale(mtv, -0.5))
		fmt.Println("PolyA position zurückgesetzt")
	}

	// danach muss mtv normalisiert werden
	mtv = rl.Vector2Normalize(mtv)

	//finde bei zwei Kontaktpunkten den Kollisionspunkt in der Mitte
	collisionpoint := contacts[0]
	if len(contacts) > 1 {
		collisionpoint = rl.Vector2Scale(rl.Vector2Add(contacts[0], contacts[1]), 0.5)
	}

	// Linie von A.location zu Kollisionspunkt
	rAP := rl.Vector2Subtract(collisionpoint, polyA.Basic.location)
	// Linie von B.location zu Kollisionspunkt
	rBP := rl.Vector2Subtract(collisionpoint, polyB.Basic.location)

	rAP_perp := rl.Vector2{-rAP.Y, rAP.X}
	rBP_perp := rl.Vector2{-rBP.Y, rBP.X}
	VtanA := rl.Vector2Scale(rAP_perp, polyA.Basic.AngVelocity)
	VtanB := rl.Vector2Scale(rBP_perp, polyB.Basic.AngVelocity)
	VgesamtA := rl.Vector2Add(polyA.Basic.Velocity, VtanA)
	VgesamtB := rl.Vector2Add(polyB.Basic.Velocity, VtanB)
	velocity_AB := rl.Vector2Subtract(VgesamtA, VgesamtB)

	//liegen Polygone auf Grund?
	if polyB.Basic.mass == math.MaxFloat32 && polyB.Basic.inertia == math.MaxFloat32 && len(contacts) > 1 && rl.Vector2Length(velocity_AB) < 2.0 {
		//polyA.Basic.AngVelocity = polyB.Basic.AngVelocity
		//polyA.Basic.Velocity = polyB.Basic.Velocity
		//fmt.Println("PolyA velocity zurückgesetzt")
		polyA.Basic.IsGrounded = true
	} else {
		polyA.Basic.IsGrounded = false
	}

	fmt.Println("Anzahl Kontakte:", len(contacts))
	fmt.Println("velocity AB:", velocity_AB, rl.Vector2Length(velocity_AB))
	fmt.Println("isGrounded:", polyA.Basic.IsGrounded)

	// wenn negativ, dann auf Kollisionskurs
	if rl.Vector2DotProduct(velocity_AB, rl.Vector2Scale(mtv, -1)) < 0 {
		// Impulskonstante e in Abhängigkeit der Geschwindigkeit velocity_AB ab Stärke 3 (von 0 bis 0.7)
		e := rl.Clamp(rl.Vector2Length(velocity_AB)/5, 0, 0.5)
		//e := float32(0.8)
		j_denominator := rl.Vector2DotProduct(rl.Vector2Scale(velocity_AB, -(1+e)), mtv)
		j_divLinear := rl.Vector2DotProduct(mtv, rl.Vector2Scale(mtv, (1/polyA.Basic.mass+1/polyB.Basic.mass)))
		j_divAngular := float32(math.Pow(float64(rl.Vector2DotProduct(rAP_perp, mtv)), 2))/polyA.Basic.inertia + float32(math.Pow(float64(rl.Vector2DotProduct(rBP_perp, mtv)), 2))/polyB.Basic.inertia
		j := j_denominator / (j_divLinear + j_divAngular)

		// Grundlage für Friction berechnen (t)
		t := rl.Vector2{-mtv.Y, mtv.X}
		t_scalarprodukt := rl.Vector2DotProduct(velocity_AB, t)
		t = rl.Vector2Normalize(rl.Vector2Scale(t, t_scalarprodukt))
		friction := float32(-0.15)

		//kalkuliere Kraft für polyA
		forceA = rl.Vector2Add(rl.Vector2Scale(mtv, (j/polyA.Basic.mass)), rl.Vector2Scale(t, (friction*-j/polyA.Basic.mass)))
		angForceA = rl.Vector2DotProduct(rAP_perp, rl.Vector2Add(rl.Vector2Scale(mtv, j/polyA.Basic.inertia), rl.Vector2Scale(t, friction*-j/polyA.Basic.inertia)))

		//kalkuliere Kraft für polyB
		forceB = rl.Vector2Add(rl.Vector2Scale(mtv, (-j/polyB.Basic.mass)), rl.Vector2Scale(t, (friction*j/polyB.Basic.mass)))
		angForceB = rl.Vector2DotProduct(rBP_perp, rl.Vector2Add(rl.Vector2Scale(mtv, -j/polyB.Basic.inertia), rl.Vector2Scale(t, friction*j/polyB.Basic.inertia)))

	} else {
		forceA = rl.Vector2{0, 0}
		angForceA = 0.0
		forceB = rl.Vector2{0, 0}
		angForceB = 0.0
	}
	fmt.Println("forceA", forceA, "angforceA", angForceA, "forceB", forceB, "angForceB", angForceB)
	fmt.Println("-----")
	return forceA, angForceA, forceB, angForceB
}
