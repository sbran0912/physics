package lib

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Shape interface {
	Update()
	Draw(c rl.Color, thick float32)
	Rotate(angle float32)
	ApplyForce(force rl.Vector2, angForce float32)
	ResetPos(delta rl.Vector2)
	ApplyGravity(gravity rl.Vector2)
}

type BasicShape struct {
	location    rl.Vector2
	velocity    rl.Vector2
	AngVelocity float32
	accel       rl.Vector2
	angAccel    float32
	mass        float32
	inertia     float32
	isGrounded  bool
}

type Polygon struct {
	basic    BasicShape
	vertices []rl.Vector2
}

type Circle struct {
	basic       BasicShape
	radius      float32
	orientation rl.Vector2
}

func CreatePolygon(x float32, y float32, w float32, h float32, wall bool) *Polygon {
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
		basic: BasicShape{
			location:   rl.Vector2{x + w/2, y + h/2},
			mass:       mass,
			inertia:    inertia,
			isGrounded: false,
		},
		vertices: vertices,
	}
	return &poly
}

func (poly *Polygon) ApplyForce(force rl.Vector2, angForce float32) {
	poly.basic.accel = rl.Vector2Add(poly.basic.accel, force)
	poly.basic.angAccel += angForce
}

func (poly *Polygon) Update() {
	//fmt.Println("Polygon Update", "angaccel:", poly.Basic.angAccel, "accel:", poly.Basic.accel)
	poly.basic.velocity = rl.Vector2Add(poly.basic.velocity, poly.basic.accel)
	poly.basic.AngVelocity += poly.basic.angAccel
	//fmt.Println("Polygon Update", "angvel:", poly.Basic.AngVelocity, "vel:", poly.Basic.velocity)

	//poly.basic.velocity = rl.Vector2Scale(poly.basic.velocity, 0.995)
	poly.basic.location = rl.Vector2Add(poly.basic.location, poly.basic.velocity)
	for i := range poly.vertices {
		poly.vertices[i] = rl.Vector2Add(poly.vertices[i], poly.basic.velocity)
	}
	//poly.basic.angVelocity *= 0.995
	poly.Rotate(poly.basic.AngVelocity)

	poly.basic.accel = rl.Vector2{0, 0}
	poly.basic.angAccel = 0
}

func (poly *Polygon) Draw(c rl.Color, thick float32) {
	n := len(poly.vertices)
	for i := range n {
		rl.DrawLineEx(poly.vertices[i], poly.vertices[(i+1)%n], thick, c)
	}
	rl.DrawCircleV(poly.basic.location, 5, c)
}

func (poly *Polygon) Rotate(angle float32) {
	for i := range poly.vertices {
		relativePos := rl.Vector2Subtract(poly.vertices[i], poly.basic.location)
		rotatedPos := rl.Vector2Rotate(relativePos, angle)
		poly.vertices[i] = rl.Vector2Add(rotatedPos, poly.basic.location)
	}
}

func (poly *Polygon) ResetPos(delta rl.Vector2) {
	for i := range poly.vertices {
		poly.vertices[i] = rl.Vector2Add(poly.vertices[i], delta)
	}
	poly.basic.location = rl.Vector2Add(poly.basic.location, delta)
}

func (poly *Polygon) ApplyGravity(gravity rl.Vector2) {
	if poly.basic.mass < math.MaxFloat32 {
		if !poly.basic.isGrounded {
			poly.ApplyForce(gravity, 0)
		} else {
			//Rückstoß der Kollision wird neutralisiert
			poly.ApplyForce(rl.Vector2Scale(poly.basic.velocity, -1), poly.basic.AngVelocity*-1)
		}
	}
}

func CreateCircle(x float32, y float32, r float32, wall bool) *Circle {
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
			location: rl.Vector2{x, y},
			mass:     mass,
			inertia:  inertia,
		},
		radius:      r,
		orientation: rl.Vector2{r + x, y},
	}
	return &circle
}

func (circle *Circle) ApplyForce(force rl.Vector2, angForce float32) {
	circle.basic.accel = rl.Vector2Add(circle.basic.accel, force)
	circle.basic.angAccel += angForce
}

func (circle *Circle) Update() {
	circle.basic.velocity = rl.Vector2Add(circle.basic.velocity, circle.basic.accel)
	circle.basic.accel = rl.Vector2{0.0, 0.0}
	circle.basic.AngVelocity += circle.basic.angAccel
	circle.basic.angAccel = 0.0

	circle.basic.location = rl.Vector2Add(circle.basic.location, circle.basic.velocity)
	circle.orientation = rl.Vector2Add(circle.orientation, circle.basic.velocity)
	circle.Rotate(circle.basic.AngVelocity)
}

func (circle *Circle) Draw(c rl.Color, thick float32) {
	rl.DrawRing(circle.basic.location, circle.radius-thick, circle.radius, 0, 360, 1, c)
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

func (circle *Circle) ApplyGravity(gravity rl.Vector2) {
	if circle.basic.mass < math.MaxFloat32 {
		if !circle.basic.isGrounded {
			circle.ApplyForce(gravity, 0)
		} else {
			//Rückstoß der Kollision wird neutralisiert
			circle.ApplyForce(rl.Vector2Scale(circle.basic.velocity, -1), circle.basic.AngVelocity*-1)
		}
	}

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

func resetPolyPositionsBasedOnMass(polyA, polyB *Polygon, mtv rl.Vector2) {
	// Objekte um mtv zurücksetzen basierend auf ihrer Masse
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
}

func calculatePolyRelativeVectors(polyA, polyB *Polygon, collisionPoint rl.Vector2) (rAP_perp, rBP_perp, VgesamtA, velocity_AB rl.Vector2) {

	// Linie von A.location zu Kollisionspunkt
	rAP := rl.Vector2Subtract(collisionPoint, polyA.basic.location)
	// Linie von B.location zu Kollisionspunkt
	rBP := rl.Vector2Subtract(collisionPoint, polyB.basic.location)

	rAP_perp = rl.Vector2{-rAP.Y, rAP.X}
	rBP_perp = rl.Vector2{-rBP.Y, rBP.X}

	VtanA := rl.Vector2Scale(rAP_perp, polyA.basic.AngVelocity)
	VtanB := rl.Vector2Scale(rBP_perp, polyB.basic.AngVelocity)
	VgesamtA = rl.Vector2Add(polyA.basic.velocity, VtanA)
	VgesamtB := rl.Vector2Add(polyB.basic.velocity, VtanB)
	velocity_AB = rl.Vector2Subtract(VgesamtA, VgesamtB)

	return rAP_perp, rBP_perp, VgesamtA, velocity_AB
}

func updatePolyGroundedMass(polyA, polyB *Polygon, contacts []rl.Vector2, velocityAB rl.Vector2) {
	// liegen Polygone auf Grund?
	if len(contacts) > 1 && rl.Vector2Length(velocityAB) < 2.0 {
		// Prüfe ob PolyB eine Wand ist (unendliche Masse/Trägheit)
		if polyB.basic.mass == math.MaxFloat32 && polyB.basic.inertia == math.MaxFloat32 {
			polyA.basic.isGrounded = true
		} else {
			polyA.basic.isGrounded = false
		}

		// Prüfe ob PolyA eine Wand ist
		if polyA.basic.mass == math.MaxFloat32 && polyA.basic.inertia == math.MaxFloat32 {
			polyB.basic.isGrounded = true
		} else {
			polyB.basic.isGrounded = false
		}
	} else {
		polyA.basic.isGrounded = false
		polyB.basic.isGrounded = false
	}
}

func updatePolyGrounded(polyA, polyB *Polygon, contacts []rl.Vector2, VgesamtA, velocityAB rl.Vector2) {
	// liegen Polygone auf Grund?
	if len(contacts) > 1 && rl.Vector2Length(velocityAB) < 1.0 && rl.Vector2Length(VgesamtA) < 1.0 {
		// beide Objekte ruhen
		polyA.basic.isGrounded = true
		polyB.basic.isGrounded = true
	} else {
		polyA.basic.isGrounded = false
		polyB.basic.isGrounded = false
	}
}

func calculatePolyImpulse(
	polyA, polyB *Polygon,
	mtv, rAP_perp, rBP_perp, velocity_AB rl.Vector2,
	e float32) float32 {

	jDenominator := rl.Vector2DotProduct(
		rl.Vector2Scale(velocity_AB, -(1+e)),
		mtv,
	)

	jDivLinear := rl.Vector2DotProduct(
		mtv,
		rl.Vector2Scale(mtv, (1/polyA.basic.mass+1/polyB.basic.mass)),
	)

	jDivAngular := float32(math.Pow(float64(rl.Vector2DotProduct(rAP_perp, mtv)), 2))/polyA.basic.inertia +
		float32(math.Pow(float64(rl.Vector2DotProduct(rBP_perp, mtv)), 2))/polyB.basic.inertia

	return jDenominator / (jDivLinear + jDivAngular)
}

func calculateFrictionVector(mtv, velocity_AB rl.Vector2) rl.Vector2 {
	// Grundlage für Friction berechnen (t)
	t := rl.Vector2{-mtv.Y, mtv.X}
	tScalarprodukt := rl.Vector2DotProduct(velocity_AB, t)
	t = rl.Vector2Normalize(rl.Vector2Scale(t, tScalarprodukt))
	return t
}
func calculatePolyCollisionForces(
	polyA, polyB *Polygon,
	mtv, rAP_perp, rBP_perp, velocity_AB rl.Vector2) (forceA rl.Vector2, angForceA float32, forceB rl.Vector2, angForceB float32) {

	// Impulskonstante festlegen
	var e float32 = 0.5

	// Reibungsvektor t berechnen
	t := calculateFrictionVector(mtv, velocity_AB)
	var friction float32 = -0.15 // Reibungskoeffizient

	// Impuls berechnen
	j := calculatePolyImpulse(polyA, polyB, mtv, rAP_perp, rBP_perp, velocity_AB, e)

	//kalkuliere Kraft für polyA
	forceA = rl.Vector2Add(rl.Vector2Scale(mtv, (j/polyA.basic.mass)), rl.Vector2Scale(t, (friction*-j/polyA.basic.mass)))
	angForceA = rl.Vector2DotProduct(rAP_perp, rl.Vector2Add(rl.Vector2Scale(mtv, j/polyA.basic.inertia), rl.Vector2Scale(t, friction*-j/polyA.basic.inertia)))

	//kalkuliere Kraft für polyB
	forceB = rl.Vector2Add(rl.Vector2Scale(mtv, (-j/polyB.basic.mass)), rl.Vector2Scale(t, (friction*j/polyB.basic.mass)))
	angForceB = rl.Vector2DotProduct(rBP_perp, rl.Vector2Add(rl.Vector2Scale(mtv, -j/polyB.basic.inertia), rl.Vector2Scale(t, friction*j/polyB.basic.inertia)))

	return forceA, angForceA, forceB, angForceB
}

// SAT-Kollisionstest
func DetectCollisionPoly(polyA, polyB *Polygon) (isColliding bool, mtv rl.Vector2, contacts []rl.Vector2) {
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
			dir := rl.Vector2Subtract(polyB.basic.location, polyA.basic.location)
			if rl.Vector2DotProduct(dir, smallestAxis) < 0 {
				smallestAxis = rl.Vector2Scale(smallestAxis, -1)
			}
		}
	}

	mtv = rl.Vector2Scale(smallestAxis, smallestOverlap)

	rawContacts := findContactPoints(polyA.vertices, polyB.vertices)
	// Hinweis: smallestAxis ist bereits normalisiert!
	contacts = transferContactsToA(rawContacts, polyA.vertices, rl.Vector2Scale(smallestAxis, -1))
	if len(contacts) > 0 {
		return true, mtv, contacts
	} else {
		return false, rl.Vector2{}, nil
	}
}

// Dispatcher: ruft die konkrete Detect-Funktion je Typkombination auf.
func DetectCollision(a, b Shape) (isColliding bool, mtv rl.Vector2, contacts []rl.Vector2) {

	switch shapeA := a.(type) {
	case *Polygon:
		switch shapeB := b.(type) {
		case *Polygon:
			return DetectCollisionPoly(shapeA, shapeB)
		case *Circle:
			// TODO: Implementiere DetectCollPolyCircle(A, B)
			// return DetectCollPolyCircle(A, B)
			return false, rl.Vector2{}, nil
		default:
			return false, rl.Vector2{}, nil
		}
	/*
		case *Circle:
			switch shapeB := b.(type) {
			case *Polygon:
				// Für jetzt: TODO
				return false, rl.Vector2{}, nil
			case *Circle:
				// TODO: Implementiere DetectCollCircleCircle(A, B)
				return false, rl.Vector2{}, nil
			default:
				return false, rl.Vector2{}, nil
			}

	*/
	default:
		return false, rl.Vector2{}, nil
	}
}

func ResolveCollisionPoly(
	polyA *Polygon,
	polyB *Polygon,
	contacts []rl.Vector2,
	mtv rl.Vector2) (forceA rl.Vector2, angForceA float32, forceB rl.Vector2, angForceB float32) {

	//Positionen der Polygone zurücksetzen
	resetPolyPositionsBasedOnMass(polyA, polyB, mtv)

	// danach muss mtv normalisiert werden
	mtv = rl.Vector2Normalize(mtv)

	//finde bei zwei Kontaktpunkten den Kollisionspunkt in der Mitte
	collisionpoint := contacts[0]
	if len(contacts) > 1 {
		collisionpoint = rl.Vector2Scale(rl.Vector2Add(contacts[0], contacts[1]), 0.5)
	}

	//finde relative Vektoren aus der Kollision
	rAP_perp, rBP_perp, vGesamtA, velocity_AB := calculatePolyRelativeVectors(polyA, polyB, collisionpoint)

	//liegen Polygone auf Grund, wenn objekte auf der ebene liegen?
	if checkParallel(mtv, rl.Vector2{0, 1}) {
		updatePolyGrounded(polyA, polyB, contacts, vGesamtA, velocity_AB)
	}

	// wenn negativ, dann auf Kollisionskurs
	if rl.Vector2DotProduct(velocity_AB, rl.Vector2Scale(mtv, -1)) < 0 {
		forceA, angForceA, forceB, angForceB = calculatePolyCollisionForces(polyA, polyB, mtv, rAP_perp, rBP_perp, velocity_AB)
	} else {
		forceA = rl.Vector2{0, 0}
		angForceA = 0.0
		forceB = rl.Vector2{0, 0}
		angForceB = 0.0
	}
	return forceA, angForceA, forceB, angForceB
}

// Dispatcher: ruft die konkrete Resolve-Funktion je Typkombination auf.
func ResolveCollision(
	a Shape,
	b Shape,
	contacts []rl.Vector2,
	mtv rl.Vector2,
) (rl.Vector2, float32, rl.Vector2, float32) {

	switch shapeA := a.(type) {

	case *Polygon:
		switch shapeB := b.(type) {

		case *Polygon:
			return ResolveCollisionPoly(shapeA, shapeB, contacts, mtv)

		case *Circle:
			// TODO
			// return ResolveCollPolyCircle(shapeA, shapeB, contacts, mtv)
			return rl.Vector2{}, 0, rl.Vector2{}, 0

		default:
			return rl.Vector2{}, 0, rl.Vector2{}, 0
		}
	/*
		case *Circle:
			switch shapeB := b.(type) {

			case *Polygon:
				// Achtung: MTV-Richtung evtl. invertieren!
				// return ResolveCollPolyCircle(shapeB, shapeA, contacts, -mtv)
				return rl.Vector2{}, 0, rl.Vector2{}, 0

			case *Circle:
				// TODO
				// return ResolveCollCircle(shapeA, shapeB, contacts, mtv)
				return rl.Vector2{}, 0, rl.Vector2{}, 0

			default:
				return rl.Vector2{}, 0, rl.Vector2{}, 0
			}
	*/

	default:
		return rl.Vector2{}, 0, rl.Vector2{}, 0
	}
}

// zeigen mtv und gravity in die gleiche Richtung?
func checkParallel(mtv rl.Vector2, gravity rl.Vector2) bool {
	return rl.FloatEquals(rl.Vector2DotProduct(mtv, gravity), rl.Vector2Length(mtv)*rl.Vector2Length(gravity))
}
