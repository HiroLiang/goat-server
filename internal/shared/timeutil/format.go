package timeutil

import "time"

const (
	FormatDateTime = "2006-01-02 15:04:05"
	FormatDate     = "2006-01-02"
	FormatTime     = "15:04:05"
	FormatISO      = time.RFC3339
)

// Format formats time using layout
func Format(t time.Time, layout string) string {
	return t.Format(layout)
}

// MustParse parses a time string using layout (panic on error)
func MustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

// Parse parses string into time.Time
func Parse(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}
