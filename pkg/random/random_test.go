package random

import (
	"math"
	"testing"
)

const (
	testBytesLength = 32
)

func TestGenerateBytes(t *testing.T) {
	b1, err := GenerateBytes(testBytesLength)
	if err != nil {
		t.Fatalf("Errored: %v", err)
	}

	b2, err := GenerateBytes(testBytesLength)
	if err != nil {
		t.Fatalf("Errored: %v", err)
	}

	if len(b1) != testBytesLength {
		t.Fatalf("Expected length to be %d, got %d", testBytesLength, len(b1))
	}

	if len(b1) != len(b2) {
		t.Fatal("Expected b1 & b2 length to be equal")
	}

	if areByteSlicesEqual(b1, b2) {
		t.Fatal("Expected b1 & b2 to be different")
	}
}

func TestGenerateString(t *testing.T) {
	s1, err := GenerateString(testBytesLength)
	if err != nil {
		t.Fatalf("Errored: %v", err)
	}

	s2, err := GenerateString(testBytesLength)
	if err != nil {
		t.Fatalf("Errored: %v", err)
	}

	b64Length := int(math.Ceil(testBytesLength/3)*4) + 4
	if len(s1) != b64Length {
		t.Fatalf("Expected length to be %d, got %d", b64Length, len(s1))
	}

	if len(s1) != len(s2) {
		t.Fatal("Expected s1 & s2 length to be equal")
	}

	if s1 == s2 {
		t.Fatal("Expected s1 & b2 to be different")
	}
}

func areByteSlicesEqual(a, b []byte) bool {
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
