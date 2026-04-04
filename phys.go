package lib

import (
	"fmt"
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
	angVelocity float32
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
	poly.basic.velocity = rl.Vector2Add(poly.basic.velocity, poly.basic.accel)
	poly.basic.angVelocity += poly.basic.angAccel

	poly.basic.velocity = rl.Vector2Scale(poly.basic.velocity, 0.9995)
	poly.basic.location = rl.Vector2Add(poly.basic.location, poly.basic.velocity)
	for i := range poly.vertices {
		poly.vertices[i] = rl.Vector2Add(poly.vertices[i], poly.basic.velocity)
	}
	poly.basic.angVelocity *= 0.9995
	poly.Rotate(poly.basic.angVelocity)

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
			// Objekt steht STABIL auf dem Boden.
			// Dämpfe die Geschwindigkeiten sanft runter (Schlaf-Zustand/Sleeping),
			// anstatt sie komplett hart auf 0 zu zwingen.
			dampingForce := rl.Vector2Scale(poly.basic.velocity, -0.5) // Reduziert Veloctiy
			dampingAngForce := poly.basic.angVelocity * -0.5           // Reduziert Rotation
			poly.ApplyForce(dampingForce, dampingAngForce)
		}
	}
}

func CreateCircle(x float32, y float32, r float32, wall bool) *Circle {
	var mass float32
	var inertia float32

	if !wall {
		mass = r * r * 2
		inertia = (r * r * r) * 100
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
	circle.basic.angVelocity += circle.basic.angAccel

	circle.basic.velocity = rl.Vector2Scale(circle.basic.velocity, 0.9995)
	circle.basic.location = rl.Vector2Add(circle.basic.location, circle.basic.velocity)
	circle.orientation = rl.Vector2Add(circle.orientation, circle.basic.velocity)

	circle.basic.angVelocity *= 0.9995
	circle.Rotate(circle.basic.angVelocity)

	circle.basic.accel = rl.Vector2{0.0, 0.0}
	circle.basic.angAccel = 0.0
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
			circle.ApplyForce(rl.Vector2Scale(circle.basic.velocity, -1), circle.basic.angVelocity*-1)
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
	for i := range n {
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

func resetCirclePolyPositionsBasedOnMass(poly *Polygon, circle *Circle, mtv rl.Vector2) {
	// Objekte um mtv zurücksetzen basierend auf ihrer Masse
	if poly.basic.mass < math.MaxFloat32 {
		poly.ResetPos(rl.Vector2Scale(mtv, -0.5))
	} else {
		circle.ResetPos(rl.Vector2Scale(mtv, 0.5))
	}

	if circle.basic.mass < math.MaxFloat32 {
		circle.ResetPos(rl.Vector2Scale(mtv, 0.5))
	} else {
		poly.ResetPos(rl.Vector2Scale(mtv, -0.5))
	}
}

func resetCirclePositionsBasedOnMass(circleA *Circle, circleB *Circle, mtv rl.Vector2) {
	// Objekte um mtv zurücksetzen basierend auf ihrer Masse
	if circleA.basic.mass < math.MaxFloat32 {
		circleA.ResetPos(rl.Vector2Scale(mtv, -0.5))
	} else {
		circleB.ResetPos(rl.Vector2Scale(mtv, 0.5))
	}

	if circleB.basic.mass < math.MaxFloat32 {
		circleB.ResetPos(rl.Vector2Scale(mtv, 0.5))
	} else {
		circleA.ResetPos(rl.Vector2Scale(mtv, -0.5))
	}
}

func calculatePolyRelativeVectors(polyA, polyB *Polygon, collisionPoint rl.Vector2) (rAP_perp, rBP_perp, VgesamtA, velocity_AB rl.Vector2) {

	// Linie von A.location zu Kollisionspunkt
	rAP := rl.Vector2Subtract(collisionPoint, polyA.basic.location)
	// Linie von B.location zu Kollisionspunkt
	rBP := rl.Vector2Subtract(collisionPoint, polyB.basic.location)

	rAP_perp = rl.Vector2{-rAP.Y, rAP.X}
	rBP_perp = rl.Vector2{-rBP.Y, rBP.X}

	VtanA := rl.Vector2Scale(rAP_perp, polyA.basic.angVelocity)
	VtanB := rl.Vector2Scale(rBP_perp, polyB.basic.angVelocity)
	VgesamtA = rl.Vector2Add(polyA.basic.velocity, VtanA)
	VgesamtB := rl.Vector2Add(polyB.basic.velocity, VtanB)
	velocity_AB = rl.Vector2Subtract(VgesamtA, VgesamtB)

	return rAP_perp, rBP_perp, VgesamtA, velocity_AB
}

func calculateCirclePolyRelativeVectors(poly *Polygon, circle *Circle, collisionPoint rl.Vector2) (rAP_perp, rBP_perp, VgesamtA, velocity_AB rl.Vector2) {

	// Linie von A.location zu Kollisionspunkt
	rAP := rl.Vector2Subtract(collisionPoint, poly.basic.location)
	// Linie von B.location zu Kollisionspunkt
	rBP := rl.Vector2Subtract(collisionPoint, circle.basic.location)

	rAP_perp = rl.Vector2{-rAP.Y, rAP.X}
	rBP_perp = rl.Vector2{-rBP.Y, rBP.X}

	VtanA := rl.Vector2Scale(rAP_perp, poly.basic.angVelocity)
	VtanB := rl.Vector2Scale(rBP_perp, circle.basic.angVelocity)
	VgesamtA = rl.Vector2Add(poly.basic.velocity, VtanA)
	VgesamtB := rl.Vector2Add(circle.basic.velocity, VtanB)
	velocity_AB = rl.Vector2Subtract(VgesamtA, VgesamtB)

	return rAP_perp, rBP_perp, VgesamtA, velocity_AB
}

func calculateCircleRelativeVectors(circleA *Circle, circleB *Circle, mtv rl.Vector2) (rAP_perp, rBP_perp, VgesamtA, velocity_AB rl.Vector2) {

	// mtv wird auf Radius des jeweiligen Kreises skaliert
	// Kollisionspunkt ist bei Kreisen nicht relevant, nur der mtv = Normale
	rA := rl.Vector2Scale(rl.Vector2Normalize(mtv), -circleA.radius)
	rB := rl.Vector2Scale(rl.Vector2Normalize(mtv), circleB.radius)

	rAP_perp = rl.Vector2{-rA.Y, rA.X}
	rBP_perp = rl.Vector2{-rB.Y, rB.X}

	VtanA := rl.Vector2Scale(rAP_perp, circleA.basic.angVelocity)
	VtanB := rl.Vector2Scale(rBP_perp, circleB.basic.angVelocity)
	VgesamtA = rl.Vector2Add(circleA.basic.velocity, VtanA)
	VgesamtB := rl.Vector2Add(circleB.basic.velocity, VtanB)
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

func isPolyGrounded(contacts []rl.Vector2, VgesamtA, velocityAB, mtv rl.Vector2) (isGrounded bool) {
	// Prüfen, ob die Bedingungen für "Grounded" erfüllt sind
	// Zwei Kontakte, geringe Geschwindigkeiten und möglichst horizontale Lage
	// es reicht wenn VgesamtA und velocity_AB geprüft werden, dann ist VgesamtB zwangsläufig auch klein
	isGrounded = len(contacts) > 1 && rl.Vector2Length(velocityAB) < 1.5 && rl.Vector2Length(VgesamtA) < 1.5 && isHorizontal(mtv, rl.Vector2{0, 1})
	return isGrounded
}

func findPolyTop(polyA, polyB *Polygon, mtv rl.Vector2) *Polygon {
	// mtv geht immer von PolyA zu PolyB
	if rl.Vector2DotProduct(mtv, rl.Vector2{0, 1}) > 0 {
		//mtv und gravitiy zeigen in die selbe Richtung, dann PolyA oben
		return polyA
	}
	return polyB
}

func checkContactStability(poly *Polygon, contacts []rl.Vector2, mtv rl.Vector2) (stable bool) {
	var ratio float32
	normal := rl.Vector2Scale(mtv, -1)
	edgeStart, edgeEnd := findReferenceEdge(poly.vertices, normal)
	edgeLength := rl.Vector2Distance(edgeStart, edgeEnd)
	contactDistance := rl.Vector2Distance(contacts[0], contacts[1])

	if edgeLength > 0 {
		ratio = contactDistance / edgeLength
	}
	stable = (ratio > 0.3)
	return stable
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

func calculateCirclePolyImpulse(
	poly *Polygon,
	circle *Circle,
	mtv, rAP_perp, velocity_AB rl.Vector2,
	e float32) float32 {

	jDenominator := rl.Vector2DotProduct(
		rl.Vector2Scale(velocity_AB, -(1+e)),
		mtv,
	)

	jDivLinear := rl.Vector2DotProduct(
		mtv,
		rl.Vector2Scale(mtv, (1/poly.basic.mass+1/circle.basic.mass)),
	)
	//nur für box
	jDivAngular := float32(math.Pow(float64(rl.Vector2DotProduct(rAP_perp, mtv)), 2)) / poly.basic.inertia

	return jDenominator / (jDivLinear + jDivAngular)
}

func calculateCircleImpulse(
	circleA *Circle,
	circleB *Circle,
	mtv, velocity_AB rl.Vector2,
	e float32) float32 {

	jDenominator := rl.Vector2DotProduct(
		rl.Vector2Scale(velocity_AB, -(1+e)),
		mtv,
	)

	jDivLinear := rl.Vector2DotProduct(
		mtv,
		rl.Vector2Scale(mtv, (1/circleA.basic.mass+1/circleB.basic.mass)),
	)

	return jDenominator / jDivLinear
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
	var e float32 = 0.4

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

func calculateCirclePolyCollisionForces(
	poly *Polygon,
	circle *Circle,
	mtv, rAP_perp, rBP_perp, velocity_AB rl.Vector2) (forceA rl.Vector2, angForceA float32, forceB rl.Vector2, angForceB float32) {

	// Impulskonstante festlegen
	var e float32 = 0.3

	// Reibungsvektor t berechnen
	t := calculateFrictionVector(mtv, velocity_AB)
	var friction float32 = -0.15 // Reibungskoeffizient

	// Impuls berechnen
	j := calculateCirclePolyImpulse(poly, circle, mtv, rAP_perp, velocity_AB, e)

	//kalkuliere Kraft für poly
	forceA = rl.Vector2Add(rl.Vector2Scale(mtv, (j/poly.basic.mass)), rl.Vector2Scale(t, (friction*-j/poly.basic.mass)))
	angForceA = rl.Vector2DotProduct(rAP_perp, rl.Vector2Add(rl.Vector2Scale(mtv, j/poly.basic.inertia), rl.Vector2Scale(t, friction*-j/poly.basic.inertia)))

	//kalkuliere Kraft für Circle
	forceB = rl.Vector2Add(rl.Vector2Scale(mtv, (-j/circle.basic.mass)), rl.Vector2Scale(t, (friction*j/circle.basic.mass)))
	angForceB = rl.Vector2DotProduct(rBP_perp, rl.Vector2Scale(t, friction*j/circle.basic.inertia))

	return forceA, angForceA, forceB, angForceB
}

func calculateCircleCollisionForces(
	circleA *Circle,
	circleB *Circle,
	mtv, rAP_perp, rBP_perp, velocity_AB rl.Vector2) (forceA rl.Vector2, angForceA float32, forceB rl.Vector2, angForceB float32) {

	// Impulskonstante festlegen
	var e float32 = 0.3

	// Reibungsvektor t berechnen
	t := calculateFrictionVector(mtv, velocity_AB)
	var friction float32 = -0.15 // Reibungskoeffizient

	// Impuls berechnen
	j := calculateCircleImpulse(circleA, circleB, mtv, velocity_AB, e)

	/*
		 	force = vec2_add(vec2_scale(normal, (0.8*j/ballA->mass)), vec2_scale(t, (0.2*-j/ballA->mass)));
			force_ang = vec2_dot(rA_perp, vec2_scale(t, 0.1*-j/ballA->inertia));

	*/
	//kalkuliere Kraft für CircleA
	forceA = rl.Vector2Add(rl.Vector2Scale(mtv, (j/circleA.basic.mass)), rl.Vector2Scale(t, (friction*j/circleA.basic.mass)))
	angForceA = rl.Vector2DotProduct(rAP_perp, rl.Vector2Scale(t, friction*j/circleA.basic.inertia))

	//kalkuliere Kraft für CircleB
	forceB = rl.Vector2Add(rl.Vector2Scale(mtv, (-j/circleB.basic.mass)), rl.Vector2Scale(t, (friction*-j/circleB.basic.mass)))
	angForceB = rl.Vector2DotProduct(rBP_perp, rl.Vector2Scale(t, friction*-j/circleB.basic.inertia))

	return forceA, angForceA, forceB, angForceB
}

func isCenterOfMassSupported(location rl.Vector2, contacts []rl.Vector2, normal rl.Vector2) bool {
	// Bei weniger als 2 Kontakten kann ein Objekt (meistens) nicht balancieren
	if len(contacts) < 2 {
		return false
	}

	// Tangente zum Boden bilden (senkrecht zur Kollisionsnormale)
	tangent := rl.Vector2{-normal.Y, normal.X}

	// Schwerpunkt auf diese Achse projizieren
	comProj := rl.Vector2DotProduct(location, tangent)

	minProj := float32(math.MaxFloat32)
	maxProj := float32(-math.MaxFloat32)

	// Kontaktpunkte auf die Achse projizieren
	for _, cp := range contacts {
		proj := rl.Vector2DotProduct(cp, tangent)
		if proj < minProj {
			minProj = proj
		}
		if proj > maxProj {
			maxProj = proj
		}
	}

	// Liegt der Schwerpunkt über der Fläche zwischen den äußeren Kontaktpunkten?
	// Wir geben eine minimale Toleranz (Epsilon), damit Objekte an Kanten etwas früher kippen
	epsilon := float32(0.5)
	return comProj >= (minProj-epsilon) && comProj <= (maxProj+epsilon)
}

// liegt das Polygon beinahe horizontal (mtv / gravity > 0.9)?
func isHorizontal(mtv rl.Vector2, gravity rl.Vector2) bool {
	//return rl.FloatEquals(rl.Vector2DotProduct(mtv, gravity), rl.Vector2Length(mtv)*rl.Vector2Length(gravity))
	result := rl.Vector2DotProduct(mtv, gravity) / (rl.Vector2Length(mtv) * rl.Vector2Length(gravity))

	return result > 0.9 || result < -0.9
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
func ResolveCollisionPoly(
	polyA *Polygon,
	polyB *Polygon,
	contacts []rl.Vector2,
	mtv rl.Vector2) (forceA rl.Vector2, angForceA float32, forceB rl.Vector2, angForceB float32) {

	// Positionen der Polygone zurücksetzen
	resetPolyPositionsBasedOnMass(polyA, polyB, mtv)

	// Danach muss mtv normalisiert werden
	mtv = rl.Vector2Normalize(mtv)

	// Für den Grounded-Check brauchen wir weiterhin einen zentralen Referenzwert
	// (Wir lassen die Logik für isGrounded intakt, auch wenn wir es vorerst ignorieren)
	centerCollisionPoint := contacts[0]
	if len(contacts) > 1 {
		centerCollisionPoint = rl.Vector2Scale(rl.Vector2Add(contacts[0], contacts[1]), 0.5)
	}

	// Werte nur für den Grounded-Check berechnen
	_, _, vGesamtA_center, velocity_AB_center := calculatePolyRelativeVectors(polyA, polyB, centerCollisionPoint)

	// Prüfe ob Polygon auf Grund gelaufen ist
	polyTop := findPolyTop(polyA, polyB, mtv)
	if isPolyGrounded(contacts, vGesamtA_center, velocity_AB_center, mtv) {
		// NEU: Prüfen ob der Schwerpunkt über den Kontaktpunkten liegt
		if isCenterOfMassSupported(polyTop.basic.location, contacts, mtv) {
			polyTop.basic.isGrounded = true
		} else {
			// Schwerpunkt hängt in der Luft -> Es muss kippen!
			polyTop.basic.isGrounded = false
		}
	} else {
		polyTop.basic.isGrounded = false
	}

	// --- NEU: Impulse für jeden Kontaktpunkt einzeln berechnen ---
	var sumForceA, sumForceB rl.Vector2
	var sumAngForceA, sumAngForceB float32
	var appliedContacts float32 = 0

	// Berechne den Impuls für jeden Kontaktpunkt einzeln
	for _, cp := range contacts {
		rAP_perp_cp, rBP_perp_cp, _, velocity_AB_cp := calculatePolyRelativeVectors(polyA, polyB, cp)

		// Prüfe für jeden Punkt individuell, ob er sich auf Kollisionskurs befindet
		if rl.Vector2DotProduct(velocity_AB_cp, rl.Vector2Scale(mtv, -1)) < 0 {
			fA, aFA, fB, aFB := calculatePolyCollisionForces(polyA, polyB, mtv, rAP_perp_cp, rBP_perp_cp, velocity_AB_cp)
			sumForceA = rl.Vector2Add(sumForceA, fA)
			sumAngForceA += aFA
			sumForceB = rl.Vector2Add(sumForceB, fB)
			sumAngForceB += aFB
			appliedContacts++
		}
	}

	// Durchschnitt der angewendeten Impulse bilden
	if appliedContacts > 0 {
		forceA = rl.Vector2Scale(sumForceA, 1.0/appliedContacts)
		angForceA = sumAngForceA / appliedContacts
		forceB = rl.Vector2Scale(sumForceB, 1.0/appliedContacts)
		angForceB = sumAngForceB / appliedContacts
	} else {
		forceA = rl.Vector2{0, 0}
		angForceA = 0.0
		forceB = rl.Vector2{0, 0}
		angForceB = 0.0
	}

	return forceA, angForceA, forceB, angForceB
}

func DetectCollisionCirclePoly(poly *Polygon, circle *Circle) (isColliding bool, mtv rl.Vector2, contacts []rl.Vector2) {
	n := len(poly.vertices)
	for i := range n {
		edge := rl.Vector2Subtract(poly.vertices[(i+1)%n], poly.vertices[i])
		verticeToCircle := rl.Vector2Subtract(circle.basic.location, poly.vertices[i])
		t := rl.Vector2DotProduct(verticeToCircle, edge) / rl.Vector2DotProduct(edge, edge)
		if t > 0 && t <= 1 {
			edge_proj := rl.Vector2Scale(edge, t)
			lineToContact := rl.Vector2Add(rl.Vector2Scale(verticeToCircle, -1), edge_proj)
			distSquare := float64(rl.Vector2DotProduct(lineToContact, lineToContact))

			if distSquare <= float64(circle.radius*circle.radius) {
				//Distanz Kreismittelpunkt kleiner als Kreisradius -> Kollision
				overlap := float32(math.Sqrt(distSquare)) - circle.radius
				mtv = rl.Vector2Scale(rl.Vector2Normalize(lineToContact), overlap)
				contact := rl.Vector2Add(poly.vertices[i], edge_proj)
				contacts = append(contacts, contact)
				return true, mtv, contacts
			}
		}
	}
	return false, rl.Vector2{}, []rl.Vector2{}
}

func ResolveCollisionCirclePoly(
	poly *Polygon,
	circle *Circle,
	contacts []rl.Vector2,
	mtv rl.Vector2) (forceA rl.Vector2, angForceA float32, forceB rl.Vector2, angForceB float32) {

	//Positionen der Polygone zurücksetzen
	resetCirclePolyPositionsBasedOnMass(poly, circle, mtv)

	// danach muss mtv normalisiert werden
	mtv = rl.Vector2Normalize(mtv)

	//es gibt nur einen Kontaktpunkt
	collisionpoint := contacts[0]

	//finde relative Vektoren aus der Kollision
	rAP_perp, rBP_perp, _, velocity_AB := calculateCirclePolyRelativeVectors(poly, circle, collisionpoint)

	// wenn negativ, dann auf Kollisionskurs
	if rl.Vector2DotProduct(velocity_AB, rl.Vector2Scale(mtv, -1)) < 0 {
		forceA, angForceA, forceB, angForceB = calculateCirclePolyCollisionForces(poly, circle, mtv, rAP_perp, rBP_perp, velocity_AB)
		fmt.Println(forceA, angForceA, forceB, angForceB)
	} else {
		forceA = rl.Vector2{0, 0}
		angForceA = 0.0
		forceB = rl.Vector2{0, 0}
		angForceB = 0.0
	}
	return forceA, angForceA, forceB, angForceB
}

func DetectCollisionCircle(circleA *Circle, circleB *Circle) (isColliding bool, mtv rl.Vector2, contacts []rl.Vector2) {
	//Distanz ermitteln
	collisionLine := rl.Vector2Subtract(circleA.basic.location, circleB.basic.location)
	dist := rl.Vector2Length(collisionLine)
	radiusTotal := circleA.radius + circleB.radius
	if dist < radiusTotal {
		overlap := dist - radiusTotal
		mtv = rl.Vector2Scale(rl.Vector2Normalize(collisionLine), overlap)
		return true, mtv, nil
	}
	return false, rl.Vector2{}, nil
}

func ResolveCollisionCircle(
	circleA *Circle,
	circleB *Circle,
	mtv rl.Vector2) (forceA rl.Vector2, angForceA float32, forceB rl.Vector2, angForceB float32) {

	//Positionen der Polygone zurücksetzen
	resetCirclePositionsBasedOnMass(circleA, circleB, mtv)

	// danach muss mtv normalisiert werden
	mtv = rl.Vector2Normalize(mtv)

	//finde relative Vektoren aus der Kollision
	rAP_perp, rBP_perp, _, velocity_AB := calculateCircleRelativeVectors(circleA, circleB, mtv)

	// wenn negativ, dann auf Kollisionskurs
	if rl.Vector2DotProduct(velocity_AB, rl.Vector2Scale(mtv, -1)) < 0 {
		forceA, angForceA, forceB, angForceB = calculateCircleCollisionForces(circleA, circleB, mtv, rAP_perp, rBP_perp, velocity_AB)
	} else {
		forceA = rl.Vector2{0, 0}
		angForceA = 0.0
		forceB = rl.Vector2{0, 0}
		angForceB = 0.0
	}
	return forceA, angForceA, forceB, angForceB
}

// Dispatcher: ruft die konkrete Detect-Funktion je Typkombination auf.
func DetectCollision(a, b Shape) (isColliding bool, mtv rl.Vector2, contacts []rl.Vector2) {

	switch shapeA := a.(type) {
	case *Polygon:
		switch shapeB := b.(type) {
		case *Polygon:
			return DetectCollisionPoly(shapeA, shapeB)
		case *Circle:
			return DetectCollisionCirclePoly(shapeA, shapeB)
		default:
			return false, rl.Vector2{}, nil
		}

	case *Circle:
		switch shapeB := b.(type) {
		case *Polygon:
			return DetectCollisionCirclePoly(shapeB, shapeA)
		case *Circle:
			// TODO: Implementiere DetectCollCircleCircle(A, B)
			return DetectCollisionCircle(shapeA, shapeB)
		default:
			return false, rl.Vector2{}, nil
		}
	default:
		return false, rl.Vector2{}, nil
	}
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
			return ResolveCollisionCirclePoly(shapeA, shapeB, contacts, mtv)
		default:
			return rl.Vector2{}, 0, rl.Vector2{}, 0
		}

	case *Circle:
		switch shapeB := b.(type) {
		case *Polygon:
			return ResolveCollisionCirclePoly(shapeB, shapeA, contacts, mtv)
		case *Circle:
			return ResolveCollisionCircle(shapeA, shapeB, mtv)
		default:
			return rl.Vector2{}, 0, rl.Vector2{}, 0
		}

	default:
		return rl.Vector2{}, 0, rl.Vector2{}, 0
	}
}

// applyForce funkion analog detectCollison
func ApplyForce(
	a Shape,
	b Shape,
	forceA rl.Vector2,
	angForceA float32,
	forceB rl.Vector2,
	angForceB float32,
) {

	switch shapeA := a.(type) {
	case *Polygon:
		switch shapeB := b.(type) {
		case *Polygon:
			shapeA.ApplyForce(forceA, angForceA)
			shapeB.ApplyForce(forceB, angForceB)
		case *Circle:
			shapeA.ApplyForce(forceA, angForceA)
			shapeB.ApplyForce(forceB, angForceB)
		}

	case *Circle:
		switch shapeB := b.(type) {
		case *Polygon:
			shapeA.ApplyForce(forceB, angForceB)
			shapeB.ApplyForce(forceA, angForceA)

		case *Circle:
			shapeA.ApplyForce(forceA, angForceA)
			shapeB.ApplyForce(forceB, angForceB)
		}
	}
}
