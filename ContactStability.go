func AnalyzeContactStabilityBoth(polyA, polyB *Polygon, contacts []rl.Vector2, mtv rl.Vector2) (
    ratioA float32, 
    ratioB float32, 
    hasStableContact bool) {
    
    if len(contacts) < 2 {
        return 0, 0, false
    }
    
    // Normale der Kollision
    normal := rl.Vector2Scale(mtv, -1)
    
    // Für Polygon A analysieren
    edgeStartA, edgeEndA := findReferenceEdge(polyA.vertices, normal)
    edgeLengthA := rl.Vector2Distance(edgeStartA, edgeEndA)
    contactDistance := rl.Vector2Distance(contacts[0], contacts[1])
    
    if edgeLengthA > 0 {
        ratioA = contactDistance / edgeLengthA
    }
    
    // Für Polygon B analysieren (inverse Normale)
    edgeStartB, edgeEndB := findReferenceEdge(polyB.vertices, rl.Vector2Scale(normal, -1))
    edgeLengthB := rl.Vector2Distance(edgeStartB, edgeEndB)
    
    if edgeLengthB > 0 {
        ratioB = contactDistance / edgeLengthB
    }
    
    // Hat stabilen Kontakt wenn mindestens eines ein gutes Verhältnis hat
    hasStableContact = (ratioA > 0.3) || (ratioB > 0.3)
    
    return ratioA, ratioB, hasStableContact
}
