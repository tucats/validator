package validator

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ParseDuration parses a duration string and returns the corresponding time.Duration.
// It handles extensions that allow specifying days in a duration, which are
// converted to hours before using the builtin time.ParseDuration function.
func ParseDuration(durationString string) (time.Duration, error) {
	var (
		days  int
		hours int
		mins  int
		secs  int
		ms    int
		err   error
	)

	// Does the string even contain the "d" character? If not, just parse the duration string as-is.
	if !strings.Contains(durationString, "d") {
		return time.ParseDuration(durationString)
	}

	// Scan the string character by character, converting the numeric parts to integers.
	days, hours, mins, secs, ms, err = parseDurationWithDays(durationString)
	if err != nil {
		return 0, err
	}

	// Now reconstruct a Go native duration string using the parsed values. First,
	// merge days into the hours value.
	hours += days * 24
	result := fmt.Sprintf("%dh%dm%ds%dms", hours, mins, secs, ms)

	// Convert the interval string to a time.Duration value
	duration, err := time.ParseDuration(result)

	return duration, err
}

// parseDurationWithDays scans the duration string character by character, converting the numeric parts to integers.
// Unlike the builtin function in the time package, this parse function allows the duration string to specify days
// with a suffix of "d".
func parseDurationWithDays(durationString string) (days int, hours int, mins int, secs int, ms int, err error) {
	chars := ""
	mSeen := false

	for _, ch := range durationString {
		value := 0
		if chars != "" {
			value, err = strconv.Atoi(chars)
			if err != nil {
				return days, hours, mins, secs, ms, ErrInvalidInteger.Context(chars)
			}
		}

		switch ch {
		case 'd':
			days = value
			mSeen = false
			chars = ""

		case 'h':
			hours = value
			mSeen = false
			chars = ""

		case 'm':
			mSeen = true

		case 's':
			if mSeen {
				ms = value
				chars = ""
			} else if chars != "" {
				secs = value
				chars = ""
			}

			mSeen = false

		default:
			if mSeen {
				if chars != "" {
					mins = value

					chars = ""
				}

				mSeen = false
			}

			if !unicode.IsSpace(ch) {
				chars += string(ch)
			}
		}
	}

	// If at the end of the string, we saw an "m" but nothing that followed it,
	// we still may have a minute value to add.
	if mSeen {
		if chars != "" {
			mins, err = strconv.Atoi(chars)
			if err != nil {
				return days, hours, mins, secs, ms, ErrInvalidInteger.Context(chars)
			}
		}

		chars = ""
	}

	// if anything left in the chars buffer, then unparsable content.
	if strings.TrimSpace(chars) != "" {
		return days, hours, mins, secs, ms, ErrInvalidDuration.Context(chars)
	}

	// If we've gotten this far without returning an error, the duration string
	// was valid.
	return days, hours, mins, secs, ms, nil
}
