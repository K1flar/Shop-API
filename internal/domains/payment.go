package domains

import "time"

type Payment struct {
	ID      uint32
	OrderID uint32
	Payment float64
	Date    time.Time
}
