package models

import "time"

type DayCell struct {
	AvailableReservations int
	Date                  time.Time
}
