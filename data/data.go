package data

import (
	"time"
	"encoding/json"
)

type Source interface {
	CreateClass(Class) error
	GetClasses(time.Time, time.Time) ([]Class, error)
	CreateBooking(Booking) error
	GetBookings(time.Time) ([]Booking, error)
}

type Class struct {
	Name string `db:"class_name"`
	Start time.Time `db:"start_date"`
	End time.Time `db:"end_date"`
	Capacity int `db:"capacity"`
}

//custom unmarshal to format dates from request YYYY-MM-DD
func (class *Class) UnmarshalJSON(b []byte) error {
	temp := struct{
		Name string `json:"class_name"`
		Start string `json:"start_date"`
		End string `json:"end_date"`
		Capacity int `json:"capacity"`
	}{}

	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}

	format := "2006-01-02"
	start, err := time.Parse(format , temp.Start)
	if err != nil {
		return err
	}
	end, err := time.Parse(format , temp.End)
	if err != nil {
		return err
	}
	class.Name = temp.Name
	class.Capacity = temp.Capacity
	class.Start = start
	class.End = end
	return nil
}

type Booking struct {
	Name string `db:"name"`
	Date time.Time `db:"date"`
}

//custom unmarshal to format dates from request YYYY-MM-DD
func (booking *Booking) UnmarshalJSON(b []byte) error {
	temp := struct{
		Name string `json:"name"`
		Date string `json:"date"`
	}{}

	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}

	format := "2006-01-02"
	date, err := time.Parse(format , temp.Date)
	if err != nil {
		return err
	}
	booking.Name = temp.Name
	booking.Date = date
	return nil
}