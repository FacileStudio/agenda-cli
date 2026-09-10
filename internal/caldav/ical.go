package caldav

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/emersion/go-ical"
)

// Event is a calendar event in a format-neutral representation.
type Event struct {
	UID         string
	Summary     string
	Description string
	Location    string
	Start       time.Time
	End         time.Time
	AllDay      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RenderICS serializes an event into an iCalendar (RFC 5545) payload.
func RenderICS(event *Event) (string, error) {
	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropVersion, "2.0")
	cal.Props.SetText(ical.PropProductID, "-//Facile Studio//Agenda CLI//EN")

	comp := ical.NewComponent(ical.CompEvent)
	comp.Props.SetText(ical.PropUID, event.UID)
	comp.Props.SetText(ical.PropSummary, event.Summary)

	if event.Description != "" {
		comp.Props.SetText(ical.PropDescription, event.Description)
	}
	if event.Location != "" {
		comp.Props.SetText(ical.PropLocation, event.Location)
	}

	if event.AllDay {
		comp.Props.SetDate(ical.PropDateTimeStart, event.Start)
		comp.Props.SetDate(ical.PropDateTimeEnd, event.End)
	} else {
		comp.Props.SetDateTime(ical.PropDateTimeStart, event.Start)
		comp.Props.SetDateTime(ical.PropDateTimeEnd, event.End)
	}

	comp.Props.SetDateTime(ical.PropDateTimeStamp, time.Now())
	comp.Props.SetDateTime(ical.PropCreated, event.CreatedAt)
	comp.Props.SetDateTime(ical.PropLastModified, event.UpdatedAt)

	cal.Children = append(cal.Children, comp)

	var buf bytes.Buffer
	enc := ical.NewEncoder(&buf)
	if err := enc.Encode(cal); err != nil {
		return "", fmt.Errorf("marshal ical: %w", err)
	}

	return buf.String(), nil
}

// ParseICS decodes an iCalendar payload into events.
func ParseICS(data string) ([]*Event, error) {
	dec := ical.NewDecoder(strings.NewReader(data))
	cal, err := dec.Decode()
	if err != nil {
		return nil, fmt.Errorf("parse ical: %w", err)
	}

	var events []*Event
	for _, child := range cal.Children {
		if child.Name != ical.CompEvent {
			continue
		}
		events = append(events, &Event{
			UID:         textProp(child, ical.PropUID),
			Summary:     textProp(child, ical.PropSummary),
			Description: textProp(child, ical.PropDescription),
			Location:    textProp(child, ical.PropLocation),
		})
	}

	return events, nil
}

func textProp(comp *ical.Component, name string) string {
	if p := comp.Props.Get(name); p != nil {
		return p.Value
	}
	return ""
}
