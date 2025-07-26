package handler

import (
	"strconv"
	"time"

	"github.com/omareloui/odinls/internal/errs"
)

func parseIntIfExists(str string) (int, error) {
	if str != "" {
		num, err := strconv.Atoi(str)
		if err != nil {
			return 0, errs.ErrInvalidFloat
		}
		return num, nil
	}
	return 0, nil
}

func parseFloatIfExists(str string) (float64, error) {
	if str != "" {
		str = nonDigitRegexp.ReplaceAllLiteralString(str, "")
		num, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0, errs.ErrInvalidFloat
		}
		return num, nil
	}
	return 0, nil
}

func parseDateOnlyIfExists(str string) (time.Time, error) {
	if str != "" {
		date, err := time.Parse(time.DateOnly, str)
		if err != nil {
			return time.Time{}, errs.ErrInvalidFloat
		}
		return date, nil
	}
	return time.Time{}, nil
}

func formatDateOnlyIfNonZero(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format(time.DateOnly)
}

func formatFloatIfNonZero(f float64) string {
	if f == 0 {
		return ""
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func formatBooleanIfNonZero(b bool) string {
	if !b {
		return ""
	}
	return "on"
}
