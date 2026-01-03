// Hilfsstruktur für eine Kante
type Edge struct {
	Max    Vec2 // Der am weitesten entfernte Punkt in Normalenrichtung
	V1, V2 Vec2 // Die beiden Endpunkte der Kante
}

// Findet die "bestpassende" Kante in eine bestimmte Richtung
func getBestEdge(poly []Vec2, normal Vec2) Edge {
	maxDist := -math.MaxFloat64
	idx := 0
	for i, v := range poly {
		projection := v.Dot(normal)
		if projection > maxDist {
			maxDist = projection
			idx = i
		}
	}

	v := poly[idx]
	vPrev := poly[(idx-1+len(poly))%len(poly)]
	vNext := poly[(idx+1)%len(poly)]

	l := v.Sub(vNext)
	r := v.Sub(vPrev)

	// Wähle die Nachbarkante, die orthogonaler zur Normalen steht
	if r.Dot(normal) <= l.Dot(normal) {
		return Edge{v, vPrev, v}
	}
	return Edge{v, v, vNext}
}

// Clipping-Funktion (Sutherland-Hodgman Prinzip)
func clip(v1, v2 Vec2, normal Vec2, offset float64) []Vec2 {
	points := []Vec2{}
	d1 := normal.Dot(v1) - offset
	d2 := normal.Dot(v2) - offset

	if d1 >= 0 {
		points = append(points, v1)
	}
	if d2 >= 0 {
		points = append(points, v2)
	}

	if d1*d2 < 0 {
		e := v2.Sub(v1)
		t := d1 / (d1 - d2)
		points = append(points, v1.Add(e.Mul(t)))
	}
	return points
}

func findContactPointsImproved(polyA, polyB []Vec2, mtv Vec2) []Vec2 {
	normal := mtv.Mul(1.0 / math.Sqrt(mtv.Dot(mtv))) // Normalisieren

	// 1. Kanten finden
	edgeA := getBestEdge(polyA, normal)
	edgeB := getBestEdge(polyB, normal.Mul(-1))

	// Bestimme Referenz- und Inzidenzkante (die flachere ist Referenz)
	// Zur Vereinfachung nehmen wir hier oft Edge A als Referenz
	ref, inc := edgeA, edgeB

	// 2. Clip gegen die Seiten der Referenzkante
	refDir := ref.V2.Sub(ref.V1)
	refDir = refDir.Mul(1.0 / math.Sqrt(refDir.Dot(refDir)))

	o1 := refDir.Dot(ref.V1)
	clipped := clip(inc.V1, inc.V2, refDir, o1)
	if len(clipped) < 2 {
		return clipped
	}

	o2 := -refDir.Dot(ref.V2)
	clipped = clip(clipped[0], clipped[1], refDir.Mul(-1), o2)
	if len(clipped) < 2 {
		return clipped
	}

	// 3. Nur Punkte behalten, die hinter der Referenzkante liegen
	refNorm := refDir.Perp()
	if refNorm.Dot(mtv) < 0 {
		refNorm = refNorm.Mul(-1)
	}

	finalPoints := []Vec2{}
	maxDepth := refNorm.Dot(ref.Max)
	for _, p := range clipped {
		if refNorm.Dot(p) <= maxDepth {
			finalPoints = append(finalPoints, p)
		}
	}
	return finalPoints
}
