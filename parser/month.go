package parser

import (
	"main/internal/utils/funcs"
	"main/models"
	"strconv"
)

func (p Parser) moveToMonth(city models.City, month int, year int) (string, error) {
	res, err := p.Get(funcs.Linkefy("calendar.do?month=", strconv.Itoa(month-p.month+(year-p.year)*12), "&consularPost=", city.Id))
	p.month = month
	p.year = year
	return res, err
}

func (p Parser) GetMonth(city models.City, month int, year int) (string, error) {
	return p.moveToMonth(city, month, year)
}
