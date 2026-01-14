package lib

import (
	"errors"
	"math/rand/v2" // Verwendung des neueren math/rand/v2 Pakets
)

func RandomInt(start, end int) int {
	return rand.IntN(end-start+1) + start
}

func RandomFloat(start float32, end float32) float32 {
	scale := rand.Float32() // [0, 1.0]
	return start + scale*(end-start)
}

// Node verwendet einen Typ-Parameter T
type Node[T any] struct {
	value T
	next  *Node[T]
}

// Stack ist nun ebenfalls generisch
type Stack[T any] struct {
	top  *Node[T]
	size int
}

// Push legt ein Element vom Typ T oben ab
func (s *Stack[T]) Push(val T) {
	newNode := &Node[T]{
		value: val,
		next:  s.top,
	}
	s.top = newNode
	s.size++
}

// Pop gibt T und einen Fehler zurück
func (s *Stack[T]) Pop() (T, error) {
	if s.IsEmpty() {
		var zero T // Erzeugt den Nullwert für den Typ T
		return zero, errors.New("stack ist leer")
	}

	val := s.top.value
	s.top = s.top.next
	s.size--
	return val, nil
}

// IsEmpty ist eine hilfreiche Utility-Funktion
func (s *Stack[T]) IsEmpty() bool {
	return s.top == nil
}

/*Beispiel für Stack
    stringStack := Stack[string]{}
	stringStack.Push("Welt")
	stringStack.Push("Hallo")

	sVal, _ := stringStack.Pop()
	fmt.Println(sVal) // "Hallo"

*/
