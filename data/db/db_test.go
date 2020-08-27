package db_test

import (
	"time"
	"errors"
	"database/sql/driver"

	. "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
	sql "github.com/jmoiron/sqlx"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/francislennon17/studio-classes/data"
	"github.com/francislennon17/studio-classes/data/db"
)

var mockErr = errors.New("error")

var _ = Describe("Db", func() {
	Describe("CreateClass Success", func() {
		mockDb, mock, _ := sqlmock.New()

		dataSource := db.NewDataSource(sql.NewDb(mockDb, ""))

		BeforeEach(func() {
			mock.
				ExpectExec(".*").
				WithArgs("yoga", time.Date(2019, time.September, 01, 0,0,0,0, time.UTC), time.Date(2019, time.September, 10, 0,0,0,0, time.UTC), 15).
				WillReturnResult(driver.ResultNoRows)

		})
		It("Successful", func() {
			err := dataSource.CreateClass(data.Class{
				Name: "yoga",
				Capacity: 15,
				Start: time.Date(2019, time.September, 01, 0,0,0,0, time.UTC),
				End: time.Date(2019, time.September, 10, 0,0,0,0, time.UTC) ,
			})
			Expect(err).To(BeNil())
		})
	})

	Describe("CreateClass Failure", func() {
		mockDb, mock, _ := sqlmock.New()

		dataSource := db.NewDataSource(sql.NewDb(mockDb, ""))

		BeforeEach(func() {
			mock.
				ExpectExec(".*").
				WithArgs("yoga", time.Date(2019, time.September, 01, 0,0,0,0, time.UTC), time.Date(2019, time.September, 10, 0,0,0,0, time.UTC), 15).
				WillReturnError(errors.New("error"))

		})
		It("Failure", func() {
			err := dataSource.CreateClass(data.Class{
				Name: "yoga",
				Capacity: 15,
				Start: time.Date(2019, time.September, 01, 0,0,0,0, time.UTC),
				End: time.Date(2019, time.September, 10, 0,0,0,0, time.UTC) ,
			})
			Expect(err).To(Equal(mockErr))
		})
	})

	Describe("GetClasses Success", func() {
		mockDb, mock, _ := sqlmock.New()

		dataSource := db.NewDataSource(sql.NewDb(mockDb, ""))

		BeforeEach(func() {
			mock.
				ExpectQuery("SELECT class_name, start_date, end_date, capacity FROM classes WHERE '2019-09-01 00:00:00 +0000 UTC' BETWEEN start_date AND end_date || '2019-09-10 00:00:00 +0000 UTC' BETWEEN start_date AND end_date").
				WillReturnRows(sqlmock.NewRows([]string{"class_name", "start_date", "end_date", "capacity"}).
					AddRow("yoga", time.Date(2019, time.September, 01, 0,0,0,0, time.UTC), time.Date(2019, time.September, 05, 0,0,0,0, time.UTC), 10).
					AddRow("boxing", time.Date(2019, time.September, 06, 0,0,0,0, time.UTC), time.Date(2019, time.September, 15, 0,0,0,0, time.UTC), 15))
		})
		It("Success", func() {
			classes, err := dataSource.GetClasses(
				time.Date(2019, time.September, 01, 0,0,0,0, time.UTC),
				time.Date(2019, time.September, 10, 0,0,0,0, time.UTC),
			)
			Expect(err).To(BeNil())
			Expect(classes).Should(ConsistOf(
				[]data.Class{{
					Name: "yoga",
					Capacity: 10,
					Start: time.Date(2019, time.September, 01, 0,0,0,0, time.UTC),
					End: time.Date(2019, time.September, 05, 0,0,0,0, time.UTC) ,
				}, {
					Name: "boxing",
					Capacity: 15,
					Start: time.Date(2019, time.September, 06, 0,0,0,0, time.UTC),
					End: time.Date(2019, time.September, 15, 0,0,0,0, time.UTC) ,
				}},
			))
		})
	})

	Describe("GetClasses Failure", func() {
		mockDb, mock, _ := sqlmock.New()

		dataSource := db.NewDataSource(sql.NewDb(mockDb, ""))

		BeforeEach(func() {
			mock.
				ExpectQuery("SELECT class_name, start_date, end_date, capacity FROM classes WHERE '2019-09-01 00:00:00 +0000 UTC' BETWEEN start_date AND end_date || '2019-09-10 00:00:00 +0000 UTC' BETWEEN start_date AND end_date").
				WillReturnError(mockErr)

		})
		It("Failure", func() {
			_, err := dataSource.GetClasses(
				time.Date(2019, time.September, 01, 0,0,0,0, time.UTC),
				time.Date(2019, time.September, 10, 0,0,0,0, time.UTC),
			)
			Expect(err).To(Equal(mockErr))
		})
	})
})