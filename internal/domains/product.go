package domains

type Product struct {
	ID          uint32
	Name        string
	Description string
	Price       float64
	Quantity    uint32
	CategoryID  uint32
	ImagePath   string
}
