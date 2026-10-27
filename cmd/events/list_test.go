package events

import (
	"testing"
)

func TestBuildQueryParams(t *testing.T) {
	params := buildQueryParams("2026-10-27", "2026-10-28", "standup")

	if params["from"] != "2026-10-27T00:00:00Z" {
		t.Errorf("expected from 2026-10-27T00:00:00Z, got %s", params["from"])
	}
	if params["to"] != "2026-10-28T23:59:59Z" {
		t.Errorf("expected to 2026-10-28T23:59:59Z, got %s", params["to"])
	}
	if params["search"] != "standup" {
		t.Errorf("expected search standup, got %s", params["search"])
	}

	rfcParams := buildQueryParams("2026-10-27T12:00:00Z", "2026-10-28T18:00:00Z", "")
	if rfcParams["from"] != "2026-10-27T12:00:00Z" {
		t.Errorf("expected from 2026-10-27T12:00:00Z, got %s", rfcParams["from"])
	}
	if rfcParams["to"] != "2026-10-28T18:00:00Z" {
		t.Errorf("expected to 2026-10-28T18:00:00Z, got %s", rfcParams["to"])
	}
}
