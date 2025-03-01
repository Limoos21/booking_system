package dto

type BookingDTO struct {
	BookingUuid string `json:"booking_uuid"`
	DateStart   string `json:"date_start"`
	DateEnd     string `json:"date_end"`
	Comment     string `json:"comment"`
	NumGuests   uint   `json:"num_guests"`
	Status      string `json:"status"`
	UserUuid    string `json:"-"`

	Table []string `json:"table"`
}
