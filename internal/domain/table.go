package domain

type Table struct {
	TableUUid      string
	RestaurantUUID string
	tableId        string
	PositionX      float64
	PositionY      float64
	MaxUser        float64
}

type TableDTO struct {
	TableUUid      string
	RestaurantUUID string
	tableId        string
	PositionX      float64
	PositionY      float64
	MaxUser        float64
	CanBook        bool
}
