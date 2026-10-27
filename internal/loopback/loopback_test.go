package loopback

import (
	"errors"
	"testing"
)

func TestValidateCode(t *testing.T) {
	result := waitResult{
		code:  "auth-code-123",
		state: "expected-state",
	}

	code, err := validateCode(result, "expected-state")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != "auth-code-123" {
		t.Errorf("expected code auth-code-123, got %s", code)
	}

	_, err = validateCode(result, "wrong-state")
	if err == nil {
		t.Error("expected error for mismatched state, got nil")
	}

	errResult := waitResult{err: errors.New("read failed")}
	_, err = validateCode(errResult, "expected-state")
	if err == nil {
		t.Error("expected error propagation, got nil")
	}
}

func TestRandomState(t *testing.T) {
	state1 := RandomState()
	state2 := RandomState()

	if state1 == "" || state2 == "" {
		t.Error("expected non-empty random state")
	}
	if state1 == state2 {
		t.Error("expected unique states")
	}
}
