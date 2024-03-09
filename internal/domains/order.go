package domains

import "time"

type Order struct {
	ID              uint32
	UserID          uint32
	Date            time.Time
	Price           float64
	Status          string
	PickupMethod    string
	DeliveryAddress string
}
