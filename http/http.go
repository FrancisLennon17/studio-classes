package http

import (
	"net/http"
	"fmt"
	"errors"
	"time"
	
	log "github.com/sirupsen/logrus"
	"github.com/labstack/echo/v4"
	"github.com/francislennon17/studio-classes/data"
)

func ListenAndServe(dataSource data.Source) {
	e := NewRouter(dataSource)
	e.Logger.Fatal(e.Start(":1323"))
}

func NewRouter(dataSource data.Source) *echo.Echo {
	e := echo.New()
	addRoutes(e, dataSource)
	return e
}

func addRoutes(e *echo.Echo, dataSource data.Source) {
	e.POST("/classes", createClass(dataSource))
	e.POST("/bookings", createBooking(dataSource))
}

func createClass(dataSource data.Source) echo.HandlerFunc {
	return func(c echo.Context) error {
		var class data.Class
		err := c.Bind(&class)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to parse request: %v", err))
		}

		log.Debugf("recieved request: %+v", class)
		log.Debug("validating class")
		if err := validateClass(class); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to validate request: %v", err))
		}

		log.Debug("fetching stored classes")
		classes, err := dataSource.GetClasses(class.Start, class.End)
		if err != nil {
			log.Errorf("failed to fetch stored classes: %v", err)
			return c.NoContent(http.StatusInternalServerError) 
		}

		if len(classes) > 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "class exists for those dates")
		}

		log.Debug("creating classes")
		createErr := dataSource.CreateClass(class)
		if createErr != nil {
			log.Errorf("failed to create class: %v", createErr)
			return c.NoContent(http.StatusInternalServerError)
		}

		log.Debug("returning 201 created")
		return c.NoContent(http.StatusCreated)
	}
}

func validateClass(class data.Class) error {
	switch {
	case class.End.Before(class.Start):
		return errors.New("start date must be before end date")
	case class.Start.Before(time.Now()):
		return errors.New("start date must not be in the past")
	case class.Capacity <= 0:
		return errors.New("class capacity must be a postive number greater than 0")
	case class.Name == "":
		return errors.New("class must have a non-empty name")
	default:
		return nil
	}
}


func getBookings(dataSource data.Source) echo.HandlerFunc {
	return func(c echo.Context) error {
		class := data.Class{
			Name: "abc",
		}
		err := dataSource.CreateClass(class)
		return c.String(http.StatusOK, fmt.Sprintf("YEOO: %v", err))
	}
}

func createBooking(dataSource data.Source) echo.HandlerFunc {
	return func(c echo.Context) error {
		var booking data.Booking
		err := c.Bind(&booking)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to parse request: %v", err))
		}

		log.Debugf("recieved request: %+v", booking)
		log.Debug("fetching current bookings")

		bookings, err := dataSource.GetBookings(booking.Date)
		if err != nil {
			log.Errorf("failed to fetch stored bookings: %v", err)
			return c.NoContent(http.StatusInternalServerError)
		}

		log.Debug("fetching class for selected date")
		classes, err := dataSource.GetClasses(booking.Date, booking.Date)
		if err != nil {
			log.Errorf("failed to fetch stored classes: %v", err)
			return c.NoContent(http.StatusInternalServerError) 
		} else if len(classes) > 1 {
			log.Errorf("Multiple classes found for selected date")
			return c.NoContent(http.StatusInternalServerError)
		}

		log.Debug("validating bookings")
		validateErr := validateBooking(bookings, classes)
		if validateErr != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to validate booking: %v", validateErr)) 
		}

		log.Debug("creating booking")
		createErr := dataSource.CreateBooking(booking)
		if createErr != nil {
			log.Errorf("failed to create booking: %v", createErr)
			return c.NoContent(http.StatusInternalServerError)
		}

		log.Debug("returning 201 created")
		return c.NoContent(http.StatusCreated)
	}
}

func validateBooking(bookings []data.Booking, classes []data.Class) error {
	if len(classes) == 0 {
		return errors.New("No classes available to book on this date")
	} else if len(bookings) >= classes[0].Capacity {
		return errors.New("reached capacity of class on this date")
	}
	return nil
}