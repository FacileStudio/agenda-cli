package cmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/FacileStudio/porte/devicelogin"
)

func TestDeviceOutcome(t *testing.T) {
	const server = "https://agenda.example.com/api"
	if fallback, err := deviceOutcome(devicelogin.ErrExchangeNotMounted, server); !fallback || err != nil {
		t.Fatalf("not mounted = %v, %v, want a fallback", fallback, err)
	}
	fallback, err := deviceOutcome(devicelogin.ErrAccountNotLinked, server)
	if fallback || err == nil || !strings.Contains(err.Error(), "sign in once at https://agenda.example.com in") {
		t.Fatalf("not linked = %v, %v, want the sign-in hint naming the site", fallback, err)
	}
	if fallback, err := deviceOutcome(devicelogin.ErrDenied, server); fallback || !errors.Is(err, devicelogin.ErrDenied) {
		t.Fatalf("denied = %v, %v, want the error unchanged", fallback, err)
	}
	if fallback, err := deviceOutcome(nil, server); fallback || err != nil {
		t.Fatalf("success = %v, %v, want no fallback and no error", fallback, err)
	}
}
