package caldav

import (
	"fmt"
	"strings"
	"time"

	"github.com/emersion/go-ical"
)

// Attendee is an invited person with an RFC 5545 role and participation status.
type Attendee struct {
	Email    string
	Name     string
	Role     string
	Response string
}

// Reminder is a VALARM that fires Minutes before the start.
type Reminder struct {
	Minutes int
	Method  string
}

// Extras holds the optional VEVENT parts beyond the core fields.
type Extras struct {
	Color     string
	Attendees []Attendee
	Reminders []Reminder
	ExDates   []time.Time
}

func (x *Extras) apply(comp *ical.Component, allDay bool) {
	if x == nil {
		return
	}
	if x.Color != "" {
		comp.Props.SetText(ical.PropColor, x.Color)
	}
	for _, a := range x.Attendees {
		comp.Props.Add(attendeeProp(a))
	}
	for _, t := range x.ExDates {
		p := ical.NewProp(ical.PropExceptionDates)
		p.SetDateTime(t.UTC())
		if allDay {
			p.SetDate(t.UTC())
		}
		comp.Props.Add(p)
	}
	for _, r := range x.Reminders {
		comp.Children = append(comp.Children, alarm(r))
	}
}

func attendeeProp(a Attendee) *ical.Prop {
	p := ical.NewProp(ical.PropAttendee)
	p.Value = "mailto:" + a.Email
	if a.Name != "" {
		p.Params.Set(ical.ParamCommonName, strings.ReplaceAll(a.Name, `"`, ""))
	}
	p.Params.Set(ical.ParamRole, "REQ-PARTICIPANT")
	if a.Role == "optional" {
		p.Params.Set(ical.ParamRole, "OPT-PARTICIPANT")
	}
	if a.Response != "" {
		p.Params.Set(ical.ParamParticipationStatus, strings.ToUpper(a.Response))
	}
	return p
}

func alarm(r Reminder) *ical.Component {
	comp := ical.NewComponent(ical.CompAlarm)
	method := strings.ToUpper(r.Method)
	if method == "" {
		method = "DISPLAY"
	}
	action := ical.NewProp(ical.PropAction)
	action.Value = method
	comp.Props.Add(action)
	trigger := ical.NewProp(ical.PropTrigger)
	trigger.Value = fmt.Sprintf("-PT%dM", r.Minutes)
	comp.Props.Add(trigger)
	comp.Props.SetText(ical.PropDescription, "Reminder")
	return comp
}
