package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	screenWidth := int32(1200)
	screenHeight := int32(800)
	rl.SetConfigFlags(rl.FlagVsyncHint) // rl.FlagFullscreenMode |
	rl.InitWindow(screenWidth, screenHeight, "Physics - raylib-go")
	defer rl.CloseWindow()

	a := CreateBox(200, 100, 150, 150)
	b := CreateBox(250, 240, 150, 150)
	//b.rotate(-0.1)

	//fmt.Println(detectCollBox(a, b))
	isColliding, mtv, contacts := SATCollision(a.vertices[:], b.vertices[:])
	fmt.Println(isColliding, mtv, contacts)

	for !rl.WindowShouldClose() {

		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkBlue)
		//a.rotate(0.01)
		a.Draw(rl.Red, 3)
		b.Draw(rl.White, 2)
		rl.DrawLineEx(b.location, rl.Vector2Add(b.location, mtv), 3, rl.Green)
		for _, c := range contacts {
			rl.DrawCircleV(c, 5, rl.Black)
		}
		rl.EndDrawing()
	}
}

type Box struct {
	vertices [4]rl.Vector2
	location rl.Vector2
}

func CreateBox(x float32, y float32, w float32, h float32) Box {
	box := Box{
		location: rl.Vector2{x + w/2, y + h/2},
		vertices: [4]rl.Vector2{{x, y}, {x + w, y}, {x + w, y + h}, {x, y + h}},
	}

	return box
}

func (box Box) Draw(c rl.Color, thick float32) {
	countVertices := len(box.vertices)
	for i := range countVertices {
		rl.DrawLineEx(box.vertices[i], box.vertices[(i+1)%countVertices], thick, c)
	}
	rl.DrawCircleV(box.location, 5, c)
}

func (box *Box) rotate(angle float32) {
	countVertices := len(box.vertices)

	for i := range countVertices {
		relativePos := rl.Vector2Subtract(box.vertices[i], box.location)
		rotatedPos := rl.Vector2Rotate(relativePos, angle)
		box.vertices[i] = rl.Vector2Add(rotatedPos, box.location)
	}
}
func intersectLine(start_a, end_a, start_b, end_b rl.Vector2) (distance float32, point rl.Vector2) {
	// line a und line b sind die beiden Geraden, welche auf Überschneidung getestet werden
	// s ist der Faktor für die Linie von start_a zum Überscheindungspunkt
	// u ist der Faktor die Linie von start_b zum Überschneidungspunkt
	// Überscheindungspunkt = Vector start_a zuzüglich Vector line_a mit s skaliert oder
	// Vector start_b zuzüglich Vector line_b mit u skaliert

	line_a := rl.Vector2Subtract(end_a, start_a)
	line_b := rl.Vector2Subtract(end_b, start_b)
	cross1 := rl.Vector2CrossProduct(line_a, line_b)
	cross2 := rl.Vector2CrossProduct(line_b, line_a)
	if math.Abs(float64(cross1)) > 0.01 { //Float kann man nicht direkt auf 0 testen!!!
		s := rl.Vector2CrossProduct(rl.Vector2Subtract(start_b, start_a), line_b) / cross1
		u := rl.Vector2CrossProduct(rl.Vector2Subtract(start_a, start_b), line_a) / cross2
		if s > 0.0001 && s < 1 && u > 0.0001 && u < 1 {
			return s, rl.Vector2Add(start_a, rl.Vector2Scale(line_a, s))
		}
	}
	return 0, rl.Vector2{0, 0}
}

func detectCollBox(boxA Box, boxB Box) (isCollision bool, collisionPoint rl.Vector2, normal rl.Vector2) {
	// Geprüft wird, ob eine Ecke von boxA in die Kante von boxB schneidet
	// Zusätzlich muss die Linie von Mittelpunkt boxA und Mittelpunkt boxB durch Kante von boxB gehen
	// i ist Index von Ecke und j ist Index von Kante
	// d = Diagonale von A.Mittelpunkt zu A.vertices(i)
	// e = Kante von B(j) zu B(j+1)%countVerticesB
	// z = Linie von A.Mittelpunkt zu B.Mittelpunkt
	// _perp = Perpendicularvektor
	// scalar_d Faktor von d für den Schnittpunkt d/e
	// scalar_z Faktor von z für den Schnittpunkt z/e
	// mtv = minimal translation vector (überlappender Teil von d zur Kante e)
	countVerticesA := len(boxA.vertices)
	countVerticesB := len(boxB.vertices)
	for i := range countVerticesA {
		for j := range countVerticesB {
			// Prüfung auf intersection von Diagonale d zu Kante e
			dist_d, _ := intersectLine(boxA.location, boxA.vertices[i], boxB.vertices[j], boxB.vertices[(j+1)%countVerticesB])
			if dist_d > 0.01 {
				// Prüfung auf intersection Linie z zu Kante e
				dist_z, _ := intersectLine(boxA.location, boxB.location, boxB.vertices[j], boxB.vertices[(j+1)%countVerticesB])
				if dist_z > 0.01 {
					// Collision findet statt
					// Objekte zurücksetzen und normal_e berechnen. Kollisionspunkt ist Ecke i von BoxA
					e := rl.Vector2Subtract(boxB.vertices[(j+1)%countVerticesB], boxB.vertices[j])
					e_perp := rl.Vector2{-(e.Y), e.X}
					d := rl.Vector2Subtract(boxA.vertices[i], boxA.location)
					d = rl.Vector2Scale(d, 1-dist_d)
					e_perp = rl.Vector2Normalize(e_perp)
					distance := rl.Vector2DotProduct(e_perp, d)
					e_perp = rl.Vector2Scale(e_perp, -distance) //mtv
					//shapeResetPos(boxA, vec2_scale(e_perp, 0.5f));
					//shapeResetPos(boxB, vec2_scale(e_perp, -0.5f));
					e_perp = rl.Vector2Normalize(e_perp) // normal_e
					return true, boxA.vertices[i], e_perp
				}
			}
		}
	}
	return false, rl.Vector2{0, 0}, rl.Vector2{0, 0}
}

// Projektionsbereich auf einer Achse
func projectPolygon(polygon []rl.Vector2, axis rl.Vector2) (min, max float32) {
	min = rl.Vector2DotProduct(polygon[0], axis)
	max = min
	for _, v := range polygon[1:] {
		proj := rl.Vector2DotProduct(v, axis)
		if proj < min {
			min = proj
		}
		if proj > max {
			max = proj
		}
	}
	return
}

// Alle Normalenachsen eines Polygons
func getAxes(polygon []rl.Vector2) []rl.Vector2 {
	axes := []rl.Vector2{}
	n := len(polygon)
	for i := 0; i < n; i++ {
		edge := rl.Vector2Normalize(rl.Vector2Subtract(polygon[(i+1)%n], polygon[i]))
		axes = append(axes, rl.Vector2{-edge.Y, edge.X})
	}
	return axes
}

// Mittelpunkt eines Polygons
func polygonCenter(polygon []rl.Vector2) rl.Vector2 {
	var center rl.Vector2
	for _, v := range polygon {
		center = rl.Vector2Add(center, v)
	}
	n := float32(len(polygon))
	return rl.Vector2Scale(center, 1/n)
}

// Schnittpunkt zweier Linien
func lineIntersect(start_a, end_a, start_b, end_b rl.Vector2) (distance float32, point rl.Vector2, isColliding bool) {
	// line a und line b sind die beiden Geraden, welche auf Überschneidung getestet werden
	// s ist der Faktor für die Linie von start_a zum Überscheindungspunkt
	// u ist der Faktor die Linie von start_b zum Überschneidungspunkt
	// Überscheindungspunkt = Vector start_a zuzüglich Vector line_a mit s skaliert oder
	// Vector start_b zuzüglich Vector line_b mit u skaliert

	line_a := rl.Vector2Subtract(end_a, start_a)
	line_b := rl.Vector2Subtract(end_b, start_b)
	cross1 := rl.Vector2CrossProduct(line_a, line_b)
	cross2 := rl.Vector2CrossProduct(line_b, line_a)
	if math.Abs(float64(cross1)) > 0.0001 {
		s := rl.Vector2CrossProduct(rl.Vector2Subtract(start_b, start_a), line_b) / cross1
		u := rl.Vector2CrossProduct(rl.Vector2Subtract(start_a, start_b), line_a) / cross2
		if s > 0.0001 && s < 1 && u > 0.0001 && u < 1 {
			return s, rl.Vector2Add(start_a, rl.Vector2Scale(line_a, s)), true
		}
	}
	return 0, rl.Vector2{0, 0}, false
}

// Kontaktpunkte aller Kanten
func findContactPoints(polyA, polyB []rl.Vector2) []rl.Vector2 {
	points := []rl.Vector2{}
	for i := range polyA {
		a1 := polyA[i]
		a2 := polyA[(i+1)%len(polyA)]
		for j := range polyB {
			b1 := polyB[j]
			b2 := polyB[(j+1)%len(polyB)]
			_, pt, ok := lineIntersect(a1, a2, b1, b2)
			if ok {
				points = append(points, pt)
			}
		}
	}
	return points
}

// SAT-Kollisionstest
func SATCollision(polyA, polyB []rl.Vector2) (bool, rl.Vector2, []rl.Vector2) {
	smallestOverlap := float32(math.MaxFloat32)
	var smallestAxis rl.Vector2
	axes := append(getAxes(polyA), getAxes(polyB)...)

	for _, axis := range axes {
		minA, maxA := projectPolygon(polyA, axis)
		minB, maxB := projectPolygon(polyB, axis)

		overlap := float32(math.Min(float64(maxA), float64(maxB))) - float32(math.Max(float64(minA), float64(minB)))
		if overlap <= 0 {
			return false, rl.Vector2{}, nil
		}

		if overlap < smallestOverlap {
			smallestOverlap = overlap
			smallestAxis = axis
			// Richtung korrigieren
			dir := rl.Vector2Subtract(polygonCenter(polyB), polygonCenter(polyA))
			if rl.Vector2DotProduct(dir, axis) < 0 {
				smallestAxis = rl.Vector2Scale(axis, -1)
			}
		}
	}

	//mtv := rl.Vector2Scale(rl.Vector2Normalize(smallestAxis), smallestOverlap)
	mtv := rl.Vector2Scale(smallestAxis, smallestOverlap)
	contacts := findContactPoints(polyA, polyB)

	return true, mtv, contacts
}
