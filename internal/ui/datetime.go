package ui

import "time"

var dutchWeekdays = [...]string{
	"zondag",
	"maandag",
	"dinsdag",
	"woensdag",
	"donderdag",
	"vrijdag",
	"zaterdag",
}

var dutchMonths = [...]string{
	"januari",
	"februari",
	"maart",
	"april",
	"mei",
	"juni",
	"juli",
	"augustus",
	"september",
	"oktober",
	"november",
	"december",
}

// formatTime formats the given time as "HH:MM" or "HH MM" (with a blinking colon).
func formatTime(t time.Time) string {
	// Blink the colon every second.
	//
	// Both ":" and " " occupy one character in a mono font,
	// so the clock never changes width.
	separator := ":"
	if t.Second()%2 == 1 {
		separator = " "
	}

	return t.Format("15") + separator + t.Format("04")
}

// formatDate formats t as a Dutch date, e.g. "maandag 2 maart 2026".
func formatDate(t time.Time) string {
	return dutchWeekdays[t.Weekday()] + " " +
		t.Format("2") + " " +
		dutchMonths[t.Month()-1] + " " +
		t.Format("2006")
}
