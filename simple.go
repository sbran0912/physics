func pointInPolygon(point Vec2, polygon []Vec2) bool {
	n := len(polygon)
	for i := 0; i < n; i++ {
		edge := polygon[(i+1)%n].Sub(polygon[i])
		toPoint := point.Sub(polygon[i])
		if edge.Perp().Dot(toPoint) < 0 {
			return false
		}
	}
	return true
}

func findContactPointsSimplified(polyA, polyB []Vec2) []Vec2 {
	contacts := []Vec2{}

	// 1. Checke, welche Ecken von A in B liegen
	for _, p := range polyA {
		if pointInPolygon(p, polyB) {
			contacts = append(contacts, p)
		}
	}

	// 2. Checke, welche Ecken von B in A liegen
	for _, p := range polyB {
		if pointInPolygon(p, polyA) {
			contacts = append(contacts, p)
		}
	}

	// 3. Falls keine Ecken drin liegen (z.B. flache Berührung),
	// nutze deine alten Kantenschnitte als Backup
	if len(contacts) == 0 {
		for i := 0; i < len(polyA); i++ {
			a1, a2 := polyA[i], polyA[(i+1)%len(polyA)]
			for j := 0; j < len(polyB); j++ {
				b1, b2 := polyB[j], polyB[(j+1)%len(polyB)]
				if pt, ok := lineIntersect(a1, a2, b1, b2); ok {
					contacts = append(contacts, pt)
				}
			}
		}
	}

	return contacts
}
