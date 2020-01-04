package dto

import (
	"time"

	"github.com/jackc/pgtype"
)

func wrapDoubleQuotes(s string) string {
	return "\"" + s + "\""
}

// UUID is created to fulfill json.Marshaler
type UUID struct {
	pgtype.UUID
}

// MarshalJSON fulfills json.Marshaler
func (u UUID) MarshalJSON() (b []byte, err error) {
	s := ""
	_ = u.AssignTo(&s)
	b = []byte(wrapDoubleQuotes(s))
	return
}

type Timestamptz struct {
	pgtype.Timestamptz
}

func (t Timestamptz) MarshalJSON() (b []byte, err error) {
	b = []byte(wrapDoubleQuotes(t.Time.Format(time.RFC3339)))
	return
}
