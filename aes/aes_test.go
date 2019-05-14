package main

import "testing"

func TestDumb(t *testing.T) {
	if dumb() != 1 {
		t.Errorf("Something really strange happened. Expected %d, got %d", 1, dumb())
	}
}
