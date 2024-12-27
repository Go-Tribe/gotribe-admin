package util

import (
	"gotribe-admin/pkg/api/known"
	"time"
)

func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(known.TimeFormat)
}
