package gorm

type Proxy struct {
	Id          uint `gorm:"primaryKey"`
	Url         string
	RequestTime int
	Country     string
	Https       bool
	Working     bool
}
