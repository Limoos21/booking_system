package controllers

import "time"

type reservationRequest struct {
	DateStart time.Time       `json:"date_start"`
	DateEnd   time.Time       `json:"date_end"`
	Capacity  int             `json:"capacity"`
	Contacts  contactsRequest `json:"contacts"`
	Table     []string
}

type contactsRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type updateReservationRequest struct {
	Contacts contactsRequest `json:"contacts"`
	Capacity int             `json:"capacity"`
}
