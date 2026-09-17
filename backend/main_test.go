package main

import "testing"

func TestResolvePortRespectsEnvOverride(t *testing.T) {
	t.Setenv("PORT", "3000")

	if got := resolvePort(); got != "3000" {
		t.Errorf("resolvePort() = %q, want %q", got, "3000")
	}
}

func TestResolvePortFallsBackTo8080WhenUnset(t *testing.T) {
	t.Setenv("PORT", "")

	if got := resolvePort(); got != "8080" {
		t.Errorf("resolvePort() = %q, want %q", got, "8080")
	}
}
