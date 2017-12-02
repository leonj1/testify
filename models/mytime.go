package models

import "time"

type MyTime struct {
	time.Time
}

func (m *MyTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	t, err := time.Parse(time.RFC3339, s[1:len(s)-1])
	if err != nil {
		t, err = time.Parse("1970-01-01 00:00:00", s[1:len(s)-1])

	}
	m.Time = t
	return
}
