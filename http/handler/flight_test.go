package handler

import (
	"bou.ke/monkey"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type GetFlightsTestSuite struct {
	suite.Suite
	sqlMock  sqlmock.Sqlmock
	e        *echo.Echo
	flight   Flight
	timeMock time.Time
}

func (suite *GetFlightsTestSuite) SetupSuite() {
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}))
	if err != nil {
		log.Fatal(err)
	}

	suite.sqlMock = sqlMock
	suite.e = echo.New()
	suite.flight = Flight{
		DB:        db,
		Validator: validator.New(),
	}

	suite.timeMock = time.Date(2020, time.January, 1, 2, 3, 0, 0, time.UTC)
}

func (suite *GetFlightsTestSuite) CallHandler(query string) (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodGet, "/flights"+query, strings.NewReader(""))
	res := httptest.NewRecorder()
	ctx := suite.e.NewContext(req, res)
	err := suite.flight.GetFlights(ctx)

	return res, err
}

func (suite *GetFlightsTestSuite) TestGetFlights_OneFlight_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	expectedMsg := `[{"id":2,"dep_city":{"id":9,"name":"Dallas"},"arr_city":{"id":6,"name":"Philadelphia"},"dep_time":"` + suite.timeMock.Format("2006-01-02T15:04:05.999999999Z07:00") + `","arr_time":"` + suite.timeMock.Format("2006-01-02T15:04:05.999999999Z07:00") + `","airplane":{"id":8,"name":"Boeing 787"},"airline":"Southwest Airlines","price":1257,"cxl_sit_id":1,"remaining_seats":67}]`

	rows := sqlmock.NewRows([]string{"id", "dep_city_id", "arr_city_id", "dep_time", "arr_time", "airplane_id", "airline", "price", "cxl_sit_id", "remaining_seats", "Airplane__id", "Airplane__name", "DepCity__id", "DepCity__name", "ArrCity__id", "ArrCity__name"}).
		AddRow(2, 9, 6, suite.timeMock, suite.timeMock, 8, "Southwest Airlines", 1257, 1, 67, 8, "Boeing 787", 9, "Dallas", 6, "Philadelphia")
	var reqStr string = "^SELECT `flights`\\.`id`,`flights`\\.`dep_city_id`,`flights`\\.`arr_city_id`,`flights`\\.`dep_time`,`flights`\\.`arr_time`,`flights`\\.`airplane_id`,`flights`\\.`airline`,`flights`\\.`price`,`flights`\\.`cxl_sit_id`,`flights`\\.`remaining_seats`,`flights`\\.`created_at`,`flights`\\.`updated_at`,`Airplane`\\.`id` AS `Airplane__id`,`Airplane`\\.`name` AS `Airplane__name`,`Airplane`\\.`created_at` AS `Airplane__created_at`,`Airplane`\\.`updated_at` AS `Airplane__updated_at`,`DepCity`\\.`id` AS `DepCity__id`,`DepCity`\\.`name` AS `DepCity__name`,`DepCity`\\.`created_at` AS `DepCity__created_at`,`DepCity`\\.`updated_at` AS `DepCity__updated_at`,`ArrCity`\\.`id` AS `ArrCity__id`,`ArrCity`\\.`name` AS `ArrCity__name`,`ArrCity`\\.`created_at` AS `ArrCity__created_at`,`ArrCity`\\.`updated_at` AS `ArrCity__updated_at` FROM `flights` LEFT JOIN `airplanes` `Airplane` ON `flights`\\.`airplane_id` \\= `Airplane`\\.`id` LEFT JOIN `cities` `DepCity` ON `flights`\\.`dep_city_id` \\= `DepCity`\\.`id` LEFT JOIN `cities` `ArrCity` ON `flights`\\.`arr_city_id` \\= `ArrCity`\\.`id` WHERE DepCity\\.name \\= \\? AND ArrCity\\.name \\= \\? AND year\\(dep_time\\) \\= \\? AND month\\(dep_time\\) \\= \\? AND day\\(dep_time\\) \\= \\?$"
	suite.sqlMock.ExpectQuery(reqStr).
		WithArgs("Dallas", "Philadelphia", 2020, 11, 24).
		WillReturnRows(rows)

	res, err := suite.CallHandler("?departure_city=Dallas&arrival_city=Philadelphia&date=2020-11-24T00:00:00Z")
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())
}

func (suite *GetFlightsTestSuite) TestGetFlights_MultipleFlights_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	expectedMsg := `[{"id":2,"dep_city":{"id":10,"name":"Tehran"},"arr_city":{"id":6,"name":"Philadelphia"},"dep_time":"` + suite.timeMock.Format("2006-01-02T15:04:05.999999999Z07:00") + `","arr_time":"` + suite.timeMock.Format("2006-01-02T15:04:05.999999999Z07:00") + `","airplane":{"id":8,"name":"Boeing 787"},"airline":"Southwest Airlines","price":1257,"cxl_sit_id":1,"remaining_seats":67},
					{"id":3,"dep_city":{"id":10,"name":"Tehran"},"arr_city":{"id":6,"name":"Philadelphia"},"dep_time":"` + suite.timeMock.Format("2006-01-02T15:04:05.999999999Z07:00") + `","arr_time":"` + suite.timeMock.Format("2006-01-02T15:04:05.999999999Z07:00") + `","airplane":{"id":9,"name":"Boeing 747"},"airline":"Iran Air","price":1258,"cxl_sit_id":2,"remaining_seats":68}]`

	rows := sqlmock.NewRows([]string{"id", "dep_city_id", "arr_city_id", "dep_time", "arr_time", "airplane_id", "airline", "price", "cxl_sit_id", "remaining_seats", "Airplane__id", "Airplane__name", "DepCity__id", "DepCity__name", "ArrCity__id", "ArrCity__name"}).
		AddRow(2, 10, 6, suite.timeMock, suite.timeMock, 8, "Southwest Airlines", 1257, 1, 67, 8, "Boeing 787", 10, "Tehran", 6, "Philadelphia").
		AddRow(3, 10, 6, suite.timeMock, suite.timeMock, 8, "Iran Air", 1258, 2, 68, 9, "Boeing 747", 10, "Tehran", 6, "Philadelphia")

	var reqStr string = "^SELECT `flights`\\.`id`,`flights`\\.`dep_city_id`,`flights`\\.`arr_city_id`,`flights`\\.`dep_time`,`flights`\\.`arr_time`,`flights`\\.`airplane_id`,`flights`\\.`airline`,`flights`\\.`price`,`flights`\\.`cxl_sit_id`,`flights`\\.`remaining_seats`,`flights`\\.`created_at`,`flights`\\.`updated_at`,`Airplane`\\.`id` AS `Airplane__id`,`Airplane`\\.`name` AS `Airplane__name`,`Airplane`\\.`created_at` AS `Airplane__created_at`,`Airplane`\\.`updated_at` AS `Airplane__updated_at`,`DepCity`\\.`id` AS `DepCity__id`,`DepCity`\\.`name` AS `DepCity__name`,`DepCity`\\.`created_at` AS `DepCity__created_at`,`DepCity`\\.`updated_at` AS `DepCity__updated_at`,`ArrCity`\\.`id` AS `ArrCity__id`,`ArrCity`\\.`name` AS `ArrCity__name`,`ArrCity`\\.`created_at` AS `ArrCity__created_at`,`ArrCity`\\.`updated_at` AS `ArrCity__updated_at` FROM `flights` LEFT JOIN `airplanes` `Airplane` ON `flights`\\.`airplane_id` \\= `Airplane`\\.`id` LEFT JOIN `cities` `DepCity` ON `flights`\\.`dep_city_id` \\= `DepCity`\\.`id` LEFT JOIN `cities` `ArrCity` ON `flights`\\.`arr_city_id` \\= `ArrCity`\\.`id` WHERE DepCity\\.name \\= \\? AND ArrCity\\.name \\= \\? AND year\\(dep_time\\) \\= \\? AND month\\(dep_time\\) \\= \\? AND day\\(dep_time\\) \\= \\?$"
	suite.sqlMock.ExpectQuery(reqStr).
		WithArgs("Tehran", "Philadelphia", 2020, 11, 24).
		RowsWillBeClosed().
		WillReturnRows(rows)
	err := suite.sqlMock.ExpectationsWereMet()

	res, err := suite.CallHandler("?departure_city=Tehran&arrival_city=Philadelphia&date=2020-11-24T00:00:00Z")
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())
}

func (suite *GetFlightsTestSuite) TestGetFlights_MissingDepCityParameter_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest
	expectedMsg := "\"Bad Request\"\n"

	res, err := suite.CallHandler("?arrival_city=Philadelphia&date=2020-11-24T00:00:00Z")
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.Equal(expectedMsg, res.Body.String())
}

func (suite *GetFlightsTestSuite) TestGetFlights_MissingArrCityParameter_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest
	expectedMsg := "\"Bad Request\"\n"

	res, err := suite.CallHandler("?departure_city=Philadelphia&date=2020-11-24T00:00:00Z")
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.Equal(expectedMsg, res.Body.String())
}

func (suite *GetFlightsTestSuite) TestGetFlights_MissingDateParameter_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest
	expectedMsg := "\"Bad Request\"\n"

	res, err := suite.CallHandler("?departure_city=Dallas^arrival_city=Yazd")
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.Equal(expectedMsg, res.Body.String())
}

func (suite *GetFlightsTestSuite) TestGetFlights_Database_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusInternalServerError
	expectedMsg := "\"Internal Server Error\"\n"

	var reqStr string = "^SELECT `flights`\\.`id`,`flights`\\.`dep_city_id`,`flights`\\.`arr_city_id`,`flights`\\.`dep_time`,`flights`\\.`arr_time`,`flights`\\.`airplane_id`,`flights`\\.`airline`,`flights`\\.`price`,`flights`\\.`cxl_sit_id`,`flights`\\.`remaining_seats`,`flights`\\.`created_at`,`flights`\\.`updated_at`,`Airplane`\\.`id` AS `Airplane__id`,`Airplane`\\.`name` AS `Airplane__name`,`Airplane`\\.`created_at` AS `Airplane__created_at`,`Airplane`\\.`updated_at` AS `Airplane__updated_at`,`DepCity`\\.`id` AS `DepCity__id`,`DepCity`\\.`name` AS `DepCity__name`,`DepCity`\\.`created_at` AS `DepCity__created_at`,`DepCity`\\.`updated_at` AS `DepCity__updated_at`,`ArrCity`\\.`id` AS `ArrCity__id`,`ArrCity`\\.`name` AS `ArrCity__name`,`ArrCity`\\.`created_at` AS `ArrCity__created_at`,`ArrCity`\\.`updated_at` AS `ArrCity__updated_at` FROM `flights` LEFT JOIN `airplanes` `Airplane` ON `flights`\\.`airplane_id` \\= `Airplane`\\.`id` LEFT JOIN `cities` `DepCity` ON `flights`\\.`dep_city_id` \\= `DepCity`\\.`id` LEFT JOIN `cities` `ArrCity` ON `flights`\\.`arr_city_id` \\= `ArrCity`\\.`id` WHERE DepCity\\.name \\= \\? AND ArrCity\\.name \\= \\? AND year\\(dep_time\\) \\= \\? AND month\\(dep_time\\) \\= \\? AND day\\(dep_time\\) \\= \\?$"
	suite.sqlMock.ExpectQuery(reqStr).
		WithArgs("Tokyo", "Philadelphia", 2020, 11, 24).
		WillReturnError(errors.New("error"))

	res, err := suite.CallHandler("?departure_city=Tokyo&arrival_city=Philadelphia&date=2020-11-24T00:00:00Z")
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())
}

func TestGetFlights(t *testing.T) {
	suite.Run(t, new(GetFlightsTestSuite))
}

type GetDatesTestSuite struct {
	suite.Suite
	sqlMock   sqlmock.Sqlmock
	e         *echo.Echo
	flight    Flight
	timeMock1 time.Time
	timeMock2 time.Time
}

func (suite *GetDatesTestSuite) SetupSuite() {
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}))
	if err != nil {
		log.Fatal(err)
	}

	suite.sqlMock = sqlMock
	suite.e = echo.New()
	suite.flight = Flight{
		DB:        db,
		Validator: validator.New(),
	}
	suite.timeMock1 = time.Date(2020, time.January, 1, 2, 3, 0, 0, time.UTC)
	suite.timeMock2 = time.Date(2021, time.January, 1, 2, 3, 0, 0, time.UTC)

}

func (suite *GetDatesTestSuite) CallHandler() (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodGet, "/flights/dates", strings.NewReader(""))
	res := httptest.NewRecorder()
	ctx := suite.e.NewContext(req, res)
	err := suite.flight.GetDates(ctx)

	return res, err
}

func (suite *GetDatesTestSuite) TestGetDates_OneDate_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	expectedMsg := `{"Dates":["2020-01-01"]}`
	rows := sqlmock.NewRows([]string{"dep_time"}).
		AddRow(suite.timeMock1)
	var reqStr = "^SELECT DISTINCT DATE\\(dep_time\\) FROM `flights`$"
	suite.sqlMock.ExpectQuery(reqStr).WillReturnRows(rows)

	res, err := suite.CallHandler()
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())

}

func (suite *GetDatesTestSuite) TestGetDates_MultipleDates_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	expectedMsg := `{"Dates":["2020-01-01", "2021-01-01"]}`
	rows := sqlmock.NewRows([]string{"dep_time"}).
		AddRow(suite.timeMock1).
		AddRow(suite.timeMock2)
	var reqStr = "^SELECT DISTINCT DATE\\(dep_time\\) FROM `flights`$"
	suite.sqlMock.ExpectQuery(reqStr).WillReturnRows(rows)

	res, err := suite.CallHandler()
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())

}

func (suite *GetDatesTestSuite) TestGetDates_Database_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusInternalServerError
	expectedMsg := "\"Internal Server Error\"\n"

	var reqStr = "^SELECT DISTINCT DATE\\(dep_time\\) FROM `flights`$"
	suite.sqlMock.ExpectQuery(reqStr).
		WillReturnError(errors.New("error"))

	res, err := suite.CallHandler()
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())
}

func TestGetDates(t *testing.T) {
	suite.Run(t, new(GetDatesTestSuite))
}

type ReserveTestSuite struct {
	suite.Suite
	sqlMock  sqlmock.Sqlmock
	e        *echo.Echo
	flight   Flight
	timeMock time.Time
}

func (suite *ReserveTestSuite) SetupSuite() {
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}))
	if err != nil {
		log.Fatal(err)
	}

	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	})

	suite.sqlMock = sqlMock
	suite.e = echo.New()
	suite.flight = Flight{
		DB:        db,
		Validator: validator.New(),
	}

	suite.timeMock = time.Date(2020, time.January, 1, 2, 3, 0, 0, time.UTC)
}

func (suite *ReserveTestSuite) TearDownSuite() {
	monkey.UnpatchAll()
}

func (suite *ReserveTestSuite) CallHandler(body string) (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodPost, "/flights/reserve", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	ctx := suite.e.NewContext(req, res)
	err := suite.flight.Reserve(ctx)

	return res, err
}

func (suite *ReserveTestSuite) TestReserve_HaveMoreSeats_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusAccepted
	expectedMsg := "\"Accepted\"\n"
	rows := sqlmock.NewRows([]string{"id", "dep_city_id", "arr_city_id", "dep_time", "arr_time", "airplane_id", "airline", "price", "cxl_sit_id", "remaining_seats", "Airplane__id", "Airplane__name", "DepCity__id", "DepCity__name", "ArrCity__id", "ArrCity__name"}).
		AddRow(1, 9, 6, suite.timeMock, suite.timeMock, 8, "Southwest Airlines", 1257, 1, 67, 8, "Boeing 787", 9, "Dallas", 6, "Philadelphia")

	reqStr := "^SELECT \\* FROM `flights` WHERE id \\= \\? ORDER BY `flights`\\.`id` LIMIT 1$"
	suite.sqlMock.ExpectQuery(reqStr).
		WithArgs(1).
		WillReturnRows(rows)

	reqStr = "^UPDATE `flights` SET `remaining_seats`\\=\\?,`updated_at`\\=\\? WHERE id = \\?"
	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectExec(reqStr).
		WithArgs(57, time.Now(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.sqlMock.ExpectCommit()

	reqBody := "{\"flight_id\":1, \"count\":10}"
	res, err := suite.CallHandler(reqBody)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())
}

func (suite *ReserveTestSuite) TestReserve_HaveOnlyEnoughSeats_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusAccepted
	expectedMsg := "\"Accepted\"\n"
	rows := sqlmock.NewRows([]string{"id", "dep_city_id", "arr_city_id", "dep_time", "arr_time", "airplane_id", "airline", "price", "cxl_sit_id", "remaining_seats", "Airplane__id", "Airplane__name", "DepCity__id", "DepCity__name", "ArrCity__id", "ArrCity__name"}).
		AddRow(1, 9, 6, suite.timeMock, suite.timeMock, 8, "Southwest Airlines", 1257, 1, 67, 8, "Boeing 787", 9, "Dallas", 6, "Philadelphia")

	reqStr := "^SELECT \\* FROM `flights` WHERE id \\= \\? ORDER BY `flights`\\.`id` LIMIT 1$"
	suite.sqlMock.ExpectQuery(reqStr).
		WithArgs(1).
		WillReturnRows(rows)

	reqStr = "^DELETE FROM `flights` WHERE id = \\?"
	suite.sqlMock.ExpectBegin()
	suite.sqlMock.ExpectExec(reqStr).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.sqlMock.ExpectCommit()

	reqBody := "{\"flight_id\":1, \"count\":67}"
	res, err := suite.CallHandler(reqBody)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())
}

func (suite *ReserveTestSuite) TestReserve_NotEnoughSeats_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest
	expectedMsg := "\"Not Enough Seats\"\n"
	rows := sqlmock.NewRows([]string{"id", "dep_city_id", "arr_city_id", "dep_time", "arr_time", "airplane_id", "airline", "price", "cxl_sit_id", "remaining_seats", "Airplane__id", "Airplane__name", "DepCity__id", "DepCity__name", "ArrCity__id", "ArrCity__name"}).
		AddRow(1, 9, 6, suite.timeMock, suite.timeMock, 8, "Southwest Airlines", 1257, 1, 67, 8, "Boeing 787", 9, "Dallas", 6, "Philadelphia")

	reqStr := "^SELECT \\* FROM `flights` WHERE id \\= \\? ORDER BY `flights`\\.`id` LIMIT 1$"
	suite.sqlMock.ExpectQuery(reqStr).
		WithArgs(1).
		WillReturnRows(rows)

	reqBody := "{\"flight_id\":1, \"count\":77}"
	res, err := suite.CallHandler(reqBody)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())
}

func (suite *ReserveTestSuite) TestReserve_MissingFlightId_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest
	expectedMsg := "\"Bad Request\"\n"

	reqBody := `{"count":10}`
	res, err := suite.CallHandler(reqBody)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())
}

func (suite *ReserveTestSuite) TestReserve_MissingCount_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusBadRequest
	expectedMsg := "\"Bad Request\"\n"

	reqBody := `{"flight_id":10}`
	res, err := suite.CallHandler(reqBody)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())
}

func (suite *ReserveTestSuite) TestReserve_Database_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusInternalServerError
	expectedMsg := "\"Internal Server Error\"\n"

	reqStr := "^SELECT \\* FROM `flights` WHERE id \\= \\? ORDER BY `flights`\\.`id` LIMIT 1$"
	suite.sqlMock.ExpectQuery(reqStr).
		WillReturnError(errors.New("error"))

	reqBody := "{\"flight_id\":1, \"count\":10}"
	res, err := suite.CallHandler(reqBody)
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())
}

func TestReserve(t *testing.T) {
	suite.Run(t, new(ReserveTestSuite))
}
