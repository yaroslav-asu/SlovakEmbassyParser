package gorm

import (
	"fmt"
	"time"
)

type Proxy struct {
	Id          uint `gorm:"primaryKey"`
	Ip          string
	Port        string
	Code        string
	Country     string
	Https       bool
	IsWorking   bool
	RequestTime time.Duration
	LastWorking time.Time
}

func (p *Proxy) Url() string {
	return fmt.Sprintf("http://%s:%s", p.Ip, p.Port)
}
