package models

import (
	"go.uber.org/zap"
	"main/internal/utils/funcs"
	"strings"
	"time"
)

type Date struct {
	date time.Time
}

func NewDate(year int, month int, day int, hour int, minute int) Date {
	return Date{
		date: time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC),
	}
}
func NewDateYMD(year int, month int, day int) Date {
	return NewDate(year, month, day, 0, 0)
}

func NewDateYM(year int, month int) Date {
	return NewDate(year, month, 0, 0, 0)
}

func NewBlankDate() Date {
	return NewDate(0, 0, 0, 0, 0)
}

func (d *Date) ChangeYear(year int) {
	d.Change(d.Minutes(), d.Hour(), d.Day(), d.Month(), year)
}

func (d *Date) ChangeMonth(month int) {
	d.Change(d.Minutes(), d.Hour(), d.Day(), month, d.Year())
}

func (d *Date) ChangeDay(day int) {
	d.Change(d.Minutes(), d.Hour(), day, d.Month(), d.Year())
}

func (d *Date) ChangeHour(hour int) {
	d.Change(d.Minutes(), hour, d.Day(), d.Month(), d.Year())
}

func (d *Date) ChangeMinutes(minutes int) {
	d.Change(minutes, d.Hour(), d.Day(), d.Month(), d.Year())
}

func (d *Date) Change(minutes int, hour int, day int, month int, year int) {
	d.date = time.Date(year, time.Month(month), day, hour, minutes, 0, 0, time.UTC)
}

func (d *Date) Minutes() int {
	return d.date.Hour()
}
func (d *Date) Hour() int {
	return d.date.Hour()
}

func (d *Date) Day() int {
	return d.date.Day()
}

func (d *Date) Month() int {
	return int(d.date.Month())
}

func (d *Date) Year() int {
	return d.date.Year()
}

func ParseDateFromString(dateString string) Date {
	dateElements := strings.Split(dateString, ".")
	if len(dateElements) != 3 {
		zap.L().Error("Got unexpect string to parse month cell date: " + dateString)
		return NewBlankDate()
	}
	intDate := funcs.StringsToIntArray(dateElements)
	day, month, year := intDate[0], intDate[1], intDate[2]
	return NewDateYMD(year, month, day)
}
