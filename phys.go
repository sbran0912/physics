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
	ClearState()
}

type BasicShape struct {
	Location    rl.Vector2
	Velocity    rl.Vector2
	AngVelocity float32
	Accel       rl.Vector2
	AngAccel    float32
	Mass        float32
	Inertia     float32
	IsGrounded  bool
}

type Polygon struct {
	Basic    BasicShape
	Vertices []rl.Vector2
}

type Circle struct {
	Basic       BasicShape
	Radius      float32
	Orientation rl.Vector2
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
		Basic: BasicShape{
			Location:   rl.Vector2{x + w/2, y + h/2},
			Mass:       mass,
			Inertia:    inertia,
			IsGrounded: false,
		},
		Vertices: vertices,
	}
	return &poly
}

func (poly *Polygon) ApplyForce(force rl.Vector2, angForce float32) {
	poly.Basic.Accel = rl.Vector2Add(poly.Basic.Accel, force)
	poly.Basic.AngAccel += angForce
}

func (poly *Polygon) Update() {
	poly.Basic.Velocity = rl.Vector2Add(poly.Basic.Velocity, poly.Basic.Accel)
	poly.Basic.AngVelocity += poly.Basic.AngAccel

	poly.Basic.Velocity = rl.Vector2Scale(poly.Basic.Velocity, 0.9995)
	poly.Basic.Location = rl.Vector2Add(poly.Basic.Location, poly.Basic.Velocity)
	for i := range poly.Vertices {
		poly.Vertices[i] = rl.Vector2Add(poly.Vertices[i], poly.Basic.Velocity)
	}
	poly.Basic.AngVelocity *= 0.9995
	poly.Rotate(poly.Basic.AngVelocity)

	poly.Basic.Accel = rl.Vector2{0, 0}
	poly.Basic.AngAccel = 0
}

func (poly *Polygon) Draw(c rl.Color, thick float32) {
	n := len(poly.Vertices)
	for i := range n {
		rl.DrawLineEx(poly.Vertices[i], poly.Vertices[(i+1)%n], thick, c)
	}
	rl.DrawCircleV(poly.Basic.Location, 5, c)
}

func (poly *Polygon) Rotate(angle float32) {
	for i := range poly.Vertices {
		relativePos := rl.Vector2Subtract(poly.Vertices[i], poly.Basic.Location)
		rotatedPos := rl.Vector2Rotate(relativePos, angle)
		poly.Vertices[i] = rl.Vector2Add(rotatedPos, poly.Basic.Location)
	}
}

func (poly *Polygon) ResetPos(delta rl.Vector2) {
	for i := range poly.Vertices {
		poly.Vertices[i] = rl.Vector2Add(poly.Vertices[i], delta)
	}
	poly.Basic.Location = rl.Vector2Add(poly.Basic.Location, delta)
}

func (poly *Polygon) ApplyGravity(gravity rl.Vector2) {
	if poly.Basic.Mass < math.MaxFloat32 {
		if !poly.Basic.IsGrounded {
			poly.ApplyForce(gravity, 0)
		} else {
			// NEU: Wenn das Objekt extrem langsam wird, zwinge es komplett zum Stillstand
			if rl.Vector2Length(poly.Basic.Velocity) < 0.5 && poly.Basic.AngVelocity < 0.1 && poly.Basic.AngVelocity > -0.1 {
				poly.Basic.Velocity = rl.Vector2{0, 0}
				poly.Basic.AngVelocity = 0
			} else {
				// Sonst dämpfe wie bisher
				dampingForce := rl.Vector2Scale(poly.Basic.Velocity, -0.5)
				dampingAngForce := poly.Basic.AngVelocity * -0.5
				poly.ApplyForce(dampingForce, dampingAngForce)
			}
		}
	}
}

func (poly *Polygon) ClearState() {
	poly.Basic.IsGrounded = false
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
		Basic: BasicShape{
			Location: rl.Vector2{x, y},
			Mass:     mass,
			Inertia:  inertia,
		},
		Radius:      r,
		Orientation: rl.Vector2{r + x, y},
	}
	return &circle
}

func (circle *Circle) ApplyForce(force rl.Vector2, angForce float32) {
	circle.Basic.Accel = rl.Vector2Add(circle.Basic.Accel, force)
	circle.Basic.AngAccel += angForce
}

func (circle *Circle) Update() {
	circle.Basic.Velocity = rl.Vector2Add(circle.Basic.Velocity, circle.Basic.Accel)
	circle.Basic.AngVelocity += circle.Basic.AngAccel

	circle.Basic.Velocity = rl.Vector2Scale(circle.Basic.Velocity, 0.9995)
	circle.Basic.Location = rl.Vector2Add(circle.Basic.Location, circle.Basic.Velocity)
	circle.Orientation = rl.Vector2Add(circle.Orientation, circle.Basic.Velocity)

	circle.Basic.AngVelocity *= 0.9995
	circle.Rotate(circle.Basic.AngVelocity)

	circle.Basic.Accel = rl.Vector2{0.0, 0.0}
	circle.Basic.AngAccel = 0.0
}

func (circle *Circle) Draw(c rl.Color, thick float32) {
	rl.DrawRing(circle.Basic.Location, circle.Radius-thick, circle.Radius, 0, 360, 1, c)
	rl.DrawCircleV(circle.Basic.Location, 3, c)
	rl.DrawLineEx(circle.Basic.Location, circle.Orientation, thick, c)
}

func (circle *Circle) Rotate(angle float32) {
	relativePos := rl.Vector2Subtract(circle.Orientation, circle.Basic.Location)
	rotatedPos := rl.Vector2Rotate(relativePos, angle)
	circle.Orientation = rl.Vector2Add(rotatedPos, circle.Basic.Location)
}

func (circle *Circle) ResetPos(delta rl.Vector2) {

	circle.Basic.Location = rl.Vector2Add(circle.Basic.Location, delta)
	circle.Orientation = rl.Vector2Add(circle.Orientation, delta)

}

func (circle *Circle) ApplyGravity(gravity rl.Vector2) {
	if circle.Basic.Mass < math.MaxFloat32 {
		if !circle.Basic.IsGrounded {
			circle.ApplyForce(gravity, 0)
		} else {
			// Objekt steht STABIL auf dem Boden.
			// Dämpfe die Geschwindigkeiten sanft runter (Schlaf-Zustand/Sleeping),
			// anstatt sie komplett hart auf 0 zu zwingen.
			dampingForce := rl.Vector2Scale(circle.Basic.Velocity, -0.5) // Reduziert Veloctiy
			dampingAngForce := circle.Basic.AngVelocity * -0.5           // Reduziert Rotation
			circle.ApplyForce(dampingForce, dampingAngForce)
		}
	}
}

func (circle *Circle) ClearState() {
	circle.Basic.IsGrounded = false
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
	for i := range n {
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
	if polyA.Basic.Mass < math.MaxFloat32 {
		polyA.ResetPos(rl.Vector2Scale(mtv, -0.5))
	} else {
		polyB.ResetPos(rl.Vector2Scale(mtv, 0.5))
	}

	if polyB.Basic.Mass < math.MaxFloat32 {
		polyB.ResetPos(rl.Vector2Scale(mtv, 0.5))
	} else {
		polyA.ResetPos(rl.Vector2Scale(mtv, -0.5))
	}
}

func resetCirclePolyPositionsBasedOnMass(poly *Polygon, circle *Circle, mtv rl.Vector2) {
	// Objekte um mtv zurücksetzen basierend auf ihrer Masse
	if poly.Basic.Mass < math.MaxFloat32 {
		poly.ResetPos(rl.Vector2Scale(mtv, -0.5))
	} else {
		circle.ResetPos(rl.Vector2Scale(mtv, 0.5))
	}

	if circle.Basic.Mass < math.MaxFloat32 {
		circle.ResetPos(rl.Vector2Scale(mtv, 0.5))
	} else {
		poly.ResetPos(rl.Vector2Scale(mtv, -0.5))
	}
}

func resetCirclePositionsBasedOnMass(circleA *Circle, circleB *Circle, mtv rl.Vector2) {
	// Objekte um mtv zurücksetzen basierend auf ihrer Masse
	if circleA.Basic.Mass < math.MaxFloat32 {
		circleA.ResetPos(rl.Vector2Scale(mtv, -0.5))
	} else {
		circleB.ResetPos(rl.Vector2Scale(mtv, 0.5))
	}

	if circleB.Basic.Mass < math.MaxFloat32 {
		circleB.ResetPos(rl.Vector2Scale(mtv, 0.5))
	} else {
		circleA.ResetPos(rl.Vector2Scale(mtv, -0.5))
	}
}

func calculatePolyRelativeVectors(polyA, polyB *Polygon, collisionPoint rl.Vector2) (rAP_perp, rBP_perp, VgesamtA, velocity_AB rl.Vector2) {

	// Linie von A.location zu Kollisionspunkt
	rAP := rl.Vector2Subtract(collisionPoint, polyA.Basic.Location)
	// Linie von B.location zu Kollisionspunkt
	rBP := rl.Vector2Subtract(collisionPoint, polyB.Basic.Location)

	rAP_perp = rl.Vector2{-rAP.Y, rAP.X}
	rBP_perp = rl.Vector2{-rBP.Y, rBP.X}

	VtanA := rl.Vector2Scale(rAP_perp, polyA.Basic.AngVelocity)
	VtanB := rl.Vector2Scale(rBP_perp, polyB.Basic.AngVelocity)
	VgesamtA = rl.Vector2Add(polyA.Basic.Velocity, VtanA)
	VgesamtB := rl.Vector2Add(polyB.Basic.Velocity, VtanB)
	velocity_AB = rl.Vector2Subtract(VgesamtA, VgesamtB)

	return rAP_perp, rBP_perp, VgesamtA, velocity_AB
}

func calculateCirclePolyRelativeVectors(poly *Polygon, circle *Circle, collisionPoint rl.Vector2) (rAP_perp, rBP_perp, VgesamtA, velocity_AB rl.Vector2) {

	// Linie von A.location zu Kollisionspunkt
	rAP := rl.Vector2Subtract(collisionPoint, poly.Basic.Location)
	// Linie von B.location zu Kollisionspunkt
	rBP := rl.Vector2Subtract(collisionPoint, circle.Basic.Location)

	rAP_perp = rl.Vector2{-rAP.Y, rAP.X}
	rBP_perp = rl.Vector2{-rBP.Y, rBP.X}

	VtanA := rl.Vector2Scale(rAP_perp, poly.Basic.AngVelocity)
	VtanB := rl.Vector2Scale(rBP_perp, circle.Basic.AngVelocity)
	VgesamtA = rl.Vector2Add(poly.Basic.Velocity, VtanA)
	VgesamtB := rl.Vector2Add(circle.Basic.Velocity, VtanB)
	velocity_AB = rl.Vector2Subtract(VgesamtA, VgesamtB)

	return rAP_perp, rBP_perp, VgesamtA, velocity_AB
}

func calculateCircleRelativeVectors(circleA *Circle, circleB *Circle, mtv rl.Vector2) (rAP_perp, rBP_perp, VgesamtA, velocity_AB rl.Vector2) {

	// mtv wird auf Radius des jeweiligen Kreises skaliert
	// Kollisionspunkt ist bei Kreisen nicht relevant, nur der mtv = Normale
	rA := rl.Vector2Scale(rl.Vector2Normalize(mtv), -circleA.Radius)
	rB := rl.Vector2Scale(rl.Vector2Normalize(mtv), circleB.Radius)

	rAP_perp = rl.Vector2{-rA.Y, rA.X}
	rBP_perp = rl.Vector2{-rB.Y, rB.X}

	VtanA := rl.Vector2Scale(rAP_perp, circleA.Basic.AngVelocity)
	VtanB := rl.Vector2Scale(rBP_perp, circleB.Basic.AngVelocity)
	VgesamtA = rl.Vector2Add(circleA.Basic.Velocity, VtanA)
	VgesamtB := rl.Vector2Add(circleB.Basic.Velocity, VtanB)
	velocity_AB = rl.Vector2Subtract(VgesamtA, VgesamtB)

	return rAP_perp, rBP_perp, VgesamtA, velocity_AB
}

func updatePolyGroundedMass(polyA, polyB *Polygon, contacts []rl.Vector2, velocityAB rl.Vector2) {
	// liegen Polygone auf Grund?
	if len(contacts) > 1 && rl.Vector2Length(velocityAB) < 2.0 {
		// Prüfe ob PolyB eine Wand ist (unendliche Masse/Trägheit)
		if polyB.Basic.Mass == math.MaxFloat32 && polyB.Basic.Inertia == math.MaxFloat32 {
			polyA.Basic.IsGrounded = true
		} else {
			polyA.Basic.IsGrounded = false
		}

		// Prüfe ob PolyA eine Wand ist
		if polyA.Basic.Mass == math.MaxFloat32 && polyA.Basic.Inertia == math.MaxFloat32 {
			polyB.Basic.IsGrounded = true
		} else {
			polyB.Basic.IsGrounded = false
		}
	} else {
		polyA.Basic.IsGrounded = false
		polyB.Basic.IsGrounded = false
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
	edgeStart, edgeEnd := findReferenceEdge(poly.Vertices, normal)
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
		rl.Vector2Scale(mtv, (1/polyA.Basic.Mass+1/polyB.Basic.Mass)),
	)

	jDivAngular := float32(math.Pow(float64(rl.Vector2DotProduct(rAP_perp, mtv)), 2))/polyA.Basic.Inertia +
		float32(math.Pow(float64(rl.Vector2DotProduct(rBP_perp, mtv)), 2))/polyB.Basic.Inertia

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
		rl.Vector2Scale(mtv, (1/poly.Basic.Mass+1/circle.Basic.Mass)),
	)
	//nur für box
	jDivAngular := float32(math.Pow(float64(rl.Vector2DotProduct(rAP_perp, mtv)), 2)) / poly.Basic.Inertia

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
		rl.Vector2Scale(mtv, (1/circleA.Basic.Mass+1/circleB.Basic.Mass)),
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
	forceA = rl.Vector2Add(rl.Vector2Scale(mtv, (j/polyA.Basic.Mass)), rl.Vector2Scale(t, (friction*-j/polyA.Basic.Mass)))
	angForceA = rl.Vector2DotProduct(rAP_perp, rl.Vector2Add(rl.Vector2Scale(mtv, j/polyA.Basic.Inertia), rl.Vector2Scale(t, friction*-j/polyA.Basic.Inertia)))

	//kalkuliere Kraft für polyB
	forceB = rl.Vector2Add(rl.Vector2Scale(mtv, (-j/polyB.Basic.Mass)), rl.Vector2Scale(t, (friction*j/polyB.Basic.Mass)))
	angForceB = rl.Vector2DotProduct(rBP_perp, rl.Vector2Add(rl.Vector2Scale(mtv, -j/polyB.Basic.Inertia), rl.Vector2Scale(t, friction*j/polyB.Basic.Inertia)))

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
	forceA = rl.Vector2Add(rl.Vector2Scale(mtv, (j/poly.Basic.Mass)), rl.Vector2Scale(t, (friction*-j/poly.Basic.Mass)))
	angForceA = rl.Vector2DotProduct(rAP_perp, rl.Vector2Add(rl.Vector2Scale(mtv, j/poly.Basic.Inertia), rl.Vector2Scale(t, friction*-j/poly.Basic.Inertia)))

	//kalkuliere Kraft für Circle
	forceB = rl.Vector2Add(rl.Vector2Scale(mtv, (-j/circle.Basic.Mass)), rl.Vector2Scale(t, (friction*j/circle.Basic.Mass)))
	angForceB = rl.Vector2DotProduct(rBP_perp, rl.Vector2Scale(t, friction*j/circle.Basic.Inertia))

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
	forceA = rl.Vector2Add(rl.Vector2Scale(mtv, (j/circleA.Basic.Mass)), rl.Vector2Scale(t, (friction*j/circleA.Basic.Mass)))
	angForceA = rl.Vector2DotProduct(rAP_perp, rl.Vector2Scale(t, friction*j/circleA.Basic.Inertia))

	//kalkuliere Kraft für CircleB
	forceB = rl.Vector2Add(rl.Vector2Scale(mtv, (-j/circleB.Basic.Mass)), rl.Vector2Scale(t, (friction*-j/circleB.Basic.Mass)))
	angForceB = rl.Vector2DotProduct(rBP_perp, rl.Vector2Scale(t, friction*-j/circleB.Basic.Inertia))

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

func isCircleGrounded(circleAngVelocity float32, velocity_AB, mtv rl.Vector2) bool {
	// Überprüfe auf geringe lineare Geschwindigkeit an der Kontaktstelle
	if rl.Vector2Length(velocity_AB) > 0.5 { // Schwellenwert kann angepasst werden
		return false
	}

	// Überprüfe auf geringe Winkelgeschwindigkeit (für rollende Objekte)
	if math.Abs(float64(circleAngVelocity)) > 0.1 { // Schwellenwert kann angepasst werden
		return false
	}

	// Überprüfe, ob die Kollisionsnormale (mtv) generell entgegengesetzt zur Schwerkraft (0,1) zeigt.
	// mtv zeigt von Objekt A zu Objekt B. Wenn A auf B aufliegt, sollte mtv tendenziell nach oben zeigen.
	gravity := rl.Vector2{0, 1}

	dot := rl.Vector2DotProduct(mtv, gravity)
	lenMtv := rl.Vector2Length(mtv)
	lenGravity := rl.Vector2Length(gravity)

	if lenMtv == 0 || lenGravity == 0 {
		return false
	}

	// Ein Wert unter -0.9 bedeutet, dass mtv stark antiparallel zur Schwerkraft ist (also nach oben zeigt).
	// Dies deutet darauf hin, dass Objekt A auf Objekt B aufliegt.
	return dot/(lenMtv*lenGravity) < -0.9
}

// SAT-Kollisionstest
func DetectCollisionPoly(polyA, polyB *Polygon) (isColliding bool, mtv rl.Vector2, contacts []rl.Vector2) {
	smallestOverlap := float32(math.MaxFloat32)
	var smallestAxis rl.Vector2

	// Alle Achsen beider Polygone sammeln
	axes := append(getAxes(polyA.Vertices), getAxes(polyB.Vertices)...)

	for _, axis := range axes {
		minA, maxA := projectPolygon(polyA.Vertices, axis)
		minB, maxB := projectPolygon(polyB.Vertices, axis)

		overlap := float32(math.Min(float64(maxA), float64(maxB))) - float32(math.Max(float64(minA), float64(minB)))
		if overlap <= 0 {
			return false, rl.Vector2{}, nil
		}

		if overlap < smallestOverlap {
			smallestOverlap = overlap
			smallestAxis = axis

			// Richtung korrigieren
			dir := rl.Vector2Subtract(polyB.Basic.Location, polyA.Basic.Location)
			if rl.Vector2DotProduct(dir, smallestAxis) < 0 {
				smallestAxis = rl.Vector2Scale(smallestAxis, -1)
			}
		}
	}

	mtv = rl.Vector2Scale(smallestAxis, smallestOverlap)

	rawContacts := findContactPoints(polyA.Vertices, polyB.Vertices)
	// Hinweis: smallestAxis ist bereits normalisiert!
	contacts = transferContactsToA(rawContacts, polyA.Vertices, rl.Vector2Scale(smallestAxis, -1))
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

	// NEU: Statischer Slop (Erlaube 0.5 Pixel Überlappung)
	slop := float32(0.5)
	mtvLength := rl.Vector2Length(mtv)

	if mtvLength > slop {
		// Schiebe die Objekte nur um den Betrag raus, der den Slop überschreitet
		mtvCorrection := rl.Vector2Scale(rl.Vector2Normalize(mtv), mtvLength-slop)
		resetPolyPositionsBasedOnMass(polyA, polyB, mtvCorrection)
	}

	// Danach muss mtv normalisiert werden
	mtv = rl.Vector2Normalize(mtv)
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
		if isCenterOfMassSupported(polyTop.Basic.Location, contacts, mtv) {
			polyTop.Basic.IsGrounded = true
		} else {
			// Schwerpunkt hängt in der Luft -> Es muss kippen!
			polyTop.Basic.IsGrounded = false
		}
	} else {
		polyTop.Basic.IsGrounded = false
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
	n := len(poly.Vertices)
	closestDistSq := float32(math.MaxFloat32)
	var closestPoint rl.Vector2
	inside := true

	for i := range n {
		a := poly.Vertices[i]
		b := poly.Vertices[(i+1)%n]
		edge := rl.Vector2Subtract(b, a)
		toCircle := rl.Vector2Subtract(circle.Basic.Location, a)

		// Check if center is outside this edge (assuming CW and inward normal)
		normal := rl.Vector2{-edge.Y, edge.X}
		if rl.Vector2DotProduct(normal, toCircle) < 0 {
			inside = false
		}

		// Find closest point on segment
		edgeLenSq := rl.Vector2DotProduct(edge, edge)
		t := rl.Vector2DotProduct(toCircle, edge) / edgeLenSq
		var currentClosest rl.Vector2
		if t <= 0 {
			currentClosest = a
		} else if t >= 1 {
			currentClosest = b
		} else {
			currentClosest = rl.Vector2Add(a, rl.Vector2Scale(edge, t))
		}

		diff := rl.Vector2Subtract(circle.Basic.Location, currentClosest)
		distSq := rl.Vector2DotProduct(diff, diff)

		if distSq < closestDistSq {
			closestDistSq = distSq
			closestPoint = currentClosest
		}
	}

	if inside || closestDistSq <= circle.Radius*circle.Radius {
		// Collision!
		var axis rl.Vector2
		var overlap float32

		if inside {
			// Center is inside. Find closest EDGE to push out.
			minOverlap := float32(math.MaxFloat32)
			for i := range n {
				a := poly.Vertices[i]
				b := poly.Vertices[(i+1)%n]
				edge := rl.Vector2Subtract(b, a)
				normal := rl.Vector2Normalize(rl.Vector2{-edge.Y, edge.X})
				toCircle := rl.Vector2Subtract(circle.Basic.Location, a)

				// dist is inward (positive)
				dist := rl.Vector2DotProduct(normal, toCircle)
				overlapVal := dist + circle.Radius
				if overlapVal < minOverlap {
					minOverlap = overlapVal
					axis = rl.Vector2Scale(normal, -1) // Push OUT
				}
			}
			overlap = minOverlap
		} else {
			dist := float32(math.Sqrt(float64(closestDistSq)))
			overlap = circle.Radius - dist
			if dist > 1e-6 {
				axis = rl.Vector2Scale(rl.Vector2Subtract(circle.Basic.Location, closestPoint), 1.0/dist)
			} else {
				// Circle center is exactly on a vertex/edge
				axis = rl.Vector2{0, -1} // Fallback
			}
		}

		mtv = rl.Vector2Scale(axis, overlap)
		contacts = append(contacts, closestPoint)
		return true, mtv, contacts
	}

	return false, rl.Vector2{}, nil
}

func ResolveCollisionCirclePoly(
	poly *Polygon,
	circle *Circle,
	contacts []rl.Vector2,
	mtv rl.Vector2) (forceA rl.Vector2, angForceA float32, forceB rl.Vector2, angForceB float32) {

	// NEU: Statischer Slop (Erlaube 0.5 Pixel Überlappung)
	slop := float32(0.5)
	mtvLength := rl.Vector2Length(mtv)

	if mtvLength > slop {
		// Schiebe die Objekte nur um den Betrag raus, der den Slop überschreitet
		mtvCorrection := rl.Vector2Scale(rl.Vector2Normalize(mtv), mtvLength-slop)
		resetCirclePolyPositionsBasedOnMass(poly, circle, mtvCorrection)
	}

	// danach muss mtv normalisiert werden
	mtv = rl.Vector2Normalize(mtv)

	//es gibt nur einen Kontaktpunkt
	collisionpoint := contacts[0]

	//finde relative Vektoren aus der Kollision
	rAP_perp, rBP_perp, _, velocity_AB := calculateCirclePolyRelativeVectors(poly, circle, collisionpoint)

	// Check for grounding
	// mtv points from poly to circle.
	// Is circle grounded on poly?
	circle.Basic.IsGrounded = isCircleGrounded(circle.Basic.AngVelocity, rl.Vector2Scale(velocity_AB, -1), rl.Vector2Scale(mtv, -1))

	// wenn negativ, dann auf Kollisionskurs
	if rl.Vector2DotProduct(velocity_AB, rl.Vector2Scale(mtv, -1)) < 0 {
		forceA, angForceA, forceB, angForceB = calculateCirclePolyCollisionForces(poly, circle, mtv, rAP_perp, rBP_perp, velocity_AB)
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
	collisionLine := rl.Vector2Subtract(circleA.Basic.Location, circleB.Basic.Location)
	dist := rl.Vector2Length(collisionLine)
	radiusTotal := circleA.Radius + circleB.Radius
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

	// NEU: Statischer Slop (Erlaube 0.5 Pixel Überlappung)
	slop := float32(0.5)
	mtvLength := rl.Vector2Length(mtv)

	if mtvLength > slop {
		// Schiebe die Objekte nur um den Betrag raus, der den Slop überschreitet
		mtvCorrection := rl.Vector2Scale(rl.Vector2Normalize(mtv), mtvLength-slop)
		resetCirclePositionsBasedOnMass(circleA, circleB, mtvCorrection)
	}

	// danach muss mtv normalisiert werden
	mtv = rl.Vector2Normalize(mtv)

	//finde relative Vektoren aus der Kollision
	rAP_perp, rBP_perp, _, velocity_AB := calculateCircleRelativeVectors(circleA, circleB, mtv)

	// Check for grounding
	circleA.Basic.IsGrounded = isCircleGrounded(circleA.Basic.AngVelocity, velocity_AB, mtv)
	circleB.Basic.IsGrounded = isCircleGrounded(circleB.Basic.AngVelocity, rl.Vector2Scale(velocity_AB, -1), rl.Vector2Scale(mtv, -1))

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
