package mock

import (
	"time"

	"github.com/francislennon17/studio-classes/data"
)

type dataSource struct {
	mockGetClassesErr error
	mockGetBookingsErr error
	mockStoreErr error
	mockBookings []data.Booking
	mockClasses []data.Class
}

type Args struct {
	MockGetClassesErr error
	MockGetBookingsErr error
	MockStoreErr error
	MockBookings []data.Booking
	MockClasses []data.Class
}

func NewMockDataSource(args Args) data.Source {
	return &dataSource {
		mockGetClassesErr: args.MockGetClassesErr,
		mockGetBookingsErr: args.MockGetBookingsErr,
		mockStoreErr: args.MockStoreErr,
		mockBookings: args.MockBookings,
		mockClasses: args.MockClasses,
	}
}

func (ds *dataSource) CreateClass(data.Class) error {
	return ds.mockStoreErr
}

func (ds *dataSource) GetClasses(time.Time, time.Time) ([]data.Class, error) {
	return ds.mockClasses, ds.mockGetClassesErr
}

func (ds *dataSource) CreateBooking(data.Booking) error {
	return ds.mockStoreErr
}

func (ds *dataSource) GetBookings(time.Time) ([]data.Booking, error) {
	return ds.mockBookings, ds.mockGetBookingsErr
}

var MockClassNoName = `{
	"start_date": "2020-10-01",
	"end_date": "2020-10-10",
	"capacity": 2
}`

var MockClassInvalidCapacity = `{
    "class_name": "test1",
    "start_date": "2020-10-01",
    "end_date": "2020-10-10",
    "capacity": -1
}`

var MockClass = `{
    "class_name": "test1",
    "start_date": "2020-10-01",
    "end_date": "2020-10-10",
    "capacity": 2
}`