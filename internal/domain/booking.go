package domain

import (
	"errors"
	"time"
)

type Booking struct {
	BookingUuid string
	DateStart   time.Time
	DateEnd     time.Time
	Comment     string
	NumGuests   uint
	Status      string //wait/sucess/canceled
	UserUuid    string
	Contact     string
}

func (b *Booking) SetStatus(status string) error {
	switch status {
	case "wait", "success", "cancel":
		b.Status = status
		return nil
	default:
		return errors.New("invalid status")
	}
}

func (b *Booking) SetNumGuests(numGuests uint, maxNumGuests uint) bool {
	if numGuests > maxNumGuests {
		return false
	}
	b.NumGuests = numGuests
	return true
}
