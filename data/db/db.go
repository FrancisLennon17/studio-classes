package db

import (
	"fmt"
	"time"
	"context"

	sql "github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	sq "github.com/Masterminds/squirrel"
    _ "github.com/go-sql-driver/mysql"
	"github.com/francislennon17/studio-classes/data"
)

type dbSource struct {
	db *sql.DB
}

const dbTimeout = 5 * time.Second

func NewDataSource(db *sql.DB) data.Source {
	return dbSource{
		db: db,
	}
}

func(source dbSource) CreateClass(class data.Class) error {
	sql, args, err := getInsertClassQuery(class).ToSql()
	if err != nil {
		return err
	}

	log.Debugf("executing query: %v\nwith args: %v", sql, args)
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, execErr := source.db.ExecContext(ctx, sql, args...)

	return execErr
}

func getInsertClassQuery(c data.Class) sq.InsertBuilder {
	return sq.
	Insert("classes").
	Columns("class_name", "start_date", "end_date", "capacity").
	Values(c.Name, c.Start, c.End, c.Capacity)
}


func(source dbSource) GetClasses(startDate time.Time, endDate time.Time) ([]data.Class, error) {
	classes := make([]data.Class, 0)
	//squirrel doesn't support "between"
	sql := fmt.Sprintf(`SELECT class_name, start_date, end_date, capacity FROM classes WHERE '%v' BETWEEN start_date AND end_date || '%v' BETWEEN start_date AND end_date`, startDate, endDate)

	log.Debugf("executing query: %v", sql)
	selectErr := source.db.Select(&classes, sql)
	return classes, selectErr
}

func(source dbSource) CreateBooking(booking data.Booking) error {
	sql, args, err := getInsertBookingQuery(booking).ToSql()
	if err != nil {
		return err
	}

	log.Debugf("executing query: %v\nwith args: %v", sql, args)
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, execErr := source.db.ExecContext(ctx, sql, args...)

	return execErr
}

func getInsertBookingQuery(b data.Booking) sq.InsertBuilder {
	return sq.
	Insert("bookings").
	Columns("name", "date").
	Values(b.Name, b.Date)
}

func(source dbSource) GetBookings(date time.Time) ([]data.Booking, error) {
	bookings := make([]data.Booking, 0)

	sql, args, err := getFetchBookingQuery(date).ToSql()
	if err != nil {
		return bookings, err
	}

	log.Debugf("executing query: %v", sql)
	selectErr := source.db.Select(&bookings, sql, args...)
	return bookings, selectErr
}

func getFetchBookingQuery(date time.Time) sq.SelectBuilder {
	return sq.
	Select("name, date").
	From("bookings").
	Where(sq.Eq{"date": date})
}