package testingutils

import (
	"testing"
)

func AssertEquals[Type comparable](t *testing.T, got, want Type) {
	t.Helper()

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertNotEquals[Type comparable](t *testing.T, got, want Type) {
	t.Helper()

	if got == want {
		t.Errorf("got %v want %v", got, want)
	}
}
