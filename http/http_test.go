package http_test

import (
	"bytes"
	"errors"
	"strings"
	"net/http"
	"net/http/httptest"

	internalHttp "github.com/francislennon17/studio-classes/http"
	"github.com/francislennon17/studio-classes/data"
	"github.com/francislennon17/studio-classes/data/mock"
	"github.com/labstack/echo/v4"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var mockArgs mock.Args

func newTestRouter(dataSource data.Source) *echo.Echo {
	return internalHttp.NewRouter(dataSource)
}

var _ = Describe("Http", func() {
	Describe("Routes", func() {
		var router *echo.Echo
		var recorder *httptest.ResponseRecorder
		var request *http.Request
		var dataSource data.Source

		BeforeEach(func() {
			recorder = httptest.NewRecorder()

			mockArgs = mock.Args{}
			dataSource = mock.NewMockDataSource(mockArgs)
		})

		Describe("POST /classes", func() {
			JustBeforeEach(func() {
				router = newTestRouter(dataSource)
				router.ServeHTTP(recorder, request)
			})
	
			AfterEach(func() {
				mockArgs = mock.Args{}
			})

			Context("failing to bind", func() {
				BeforeEach(func() {
					body := []byte(`invalid`)
					request, _ = http.NewRequest("POST", "/classes", bytes.NewBuffer(body))
				})

				It("returns a 400", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns expected message", func() {
					Expect(strings.TrimRight(recorder.Body.String(), "\n")).To(Equal(`{"message":"failed to parse request: code=415, message=Unsupported Media Type"}`))
				})
			})

			Context("No class_name field", func() {
				BeforeEach(func() {
					request, _ = http.NewRequest("POST", "/classes", bytes.NewBuffer([]byte(mock.MockClassNoName)))
					request.Header.Set("Content-Type", "application/json")
				})

				It("returns a 400", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns expected message", func() {
					Expect(strings.TrimRight(recorder.Body.String(), "\n")).To(Equal(`{"message":"failed to validate request: class must have a non-empty name"}`))
				})
			})

			Context("Invalid capacity", func() {
				BeforeEach(func() {
					request, _ = http.NewRequest("POST", "/classes", bytes.NewBuffer([]byte(mock.MockClassInvalidCapacity)))
					request.Header.Set("Content-Type", "application/json")
				})

				It("returns a 400", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns expected message", func() {
					Expect(strings.TrimRight(recorder.Body.String(), "\n")).To(Equal(`{"message":"failed to validate request: class capacity must be a postive number greater than 0"}`))
				})
			})

			Context("Error from getClasses", func() {
				BeforeEach(func() {
					mockArgs.MockGetClassesErr = errors.New("error")
					dataSource = mock.NewMockDataSource(mockArgs)

					request, _ = http.NewRequest("POST", "/classes", bytes.NewBuffer([]byte(mock.MockClass)))
					request.Header.Set("Content-Type", "application/json")
				})

				It("returns a 500", func() {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				})
			})

			Context("class found for selected date", func() {
				BeforeEach(func() {
					mockArgs.MockClasses = []data.Class{{
						Name: "test",
					},}

					dataSource = mock.NewMockDataSource(mockArgs)

					request, _ = http.NewRequest("POST", "/classes", bytes.NewBuffer([]byte(mock.MockClass)))
					request.Header.Set("Content-Type", "application/json")
				})

				It("returns a 400", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})

				It("returns expected message", func() {
					Expect(strings.TrimRight(recorder.Body.String(), "\n")).To(Equal(`{"message":"class exists for those dates"}`))
				})
			})

			Context("class found for selected date", func() {
				BeforeEach(func() {
					mockArgs.MockStoreErr = errors.New("error")
					dataSource = mock.NewMockDataSource(mockArgs)

					request, _ = http.NewRequest("POST", "/classes", bytes.NewBuffer([]byte(mock.MockClass)))
					request.Header.Set("Content-Type", "application/json")
				})

				It("returns a 500", func() {
					Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				})
			})
		})
	})
})
