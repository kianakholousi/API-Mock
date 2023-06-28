package handler

import (
	"bou.ke/monkey"
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
	sqlMock sqlmock.Sqlmock
	e       *echo.Echo
	flight  Flight
}

func (suite *GetFlightsTestSuite) SetupSuite() {
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	cfg := mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}

	db, err := gorm.Open(mysql.New(cfg))
	if err != nil {
		log.Fatal(err)
	}

	timeNowMock := time.Date(2020, time.January, 1, 2, 3, 0, 0, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return timeNowMock })
	defer patch.Unpatch()

	suite.sqlMock = sqlMock
	suite.e = echo.New()
	suite.flight = Flight{
		DB:        db,
		Validator: validator.New(),
	}
}

func (suite *GetFlightsTestSuite) CallHandler(endpoint string) (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodGet, endpoint, strings.NewReader(""))
	res := httptest.NewRecorder()
	c := suite.e.NewContext(req, res)
	err := suite.flight.Get(c)

	return res, err
}

func (suite *GetFlightsTestSuite) TestGetFlights_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	expectedMsgTime := time.Now()
	expectedMsg := `[{"id":2,"dep_city":{"ID":9,"Name":"Dallas"},"arr_city":{"ID":6,"Name":"Philadelphia"},"dep_time":"` + expectedMsgTime.Format("2006-01-02T15:04:05.999999999Z07:00") + `","arr_time":"` + expectedMsgTime.Format("2006-01-02T15:04:05.999999999Z07:00") + `","airplane":{"ID":8,"Name":"Boeing 787"},"airline":"Southwest Airlines","price":1257,"cxl_sit_id":1,"remaining_seats":67}]`

	rows := sqlmock.NewRows([]string{"id", "dep_city_id", "arr_city_id", "dep_time", "arr_time", "airplane_id", "airline", "price", "cxl_sit_id", "remaining_seats", "Airplane__id", "Airplane__name", "DepCity__id", "DepCity__name", "ArrCity__id", "ArrCity__name"}).
		AddRow(2, 9, 6, expectedMsgTime, expectedMsgTime, 8, "Southwest Airlines", 1257, 1, 67, 8, "Boeing 787", 9, "Dallas", 6, "Philadelphia")
	var reqStr string = "SELECT " + "`flights`.`id`,`flights`.`dep_city_id`,`flights`.`arr_city_id`,`flights`.`dep_time`,`flights`.`arr_time`,`flights`.`airplane_id`,`flights`.`airline`,`flights`.`price`,`flights`.`cxl_sit_id`,`flights`.`left_seat`,`flights`.`created_at`,`flights`.`updated_at`,`DepCity`.`id` AS `DepCity__id`,`DepCity`.`name` AS `DepCity__name`,`DepCity`.`created_at` AS `DepCity__created_at`,`DepCity`.`updated_at` AS `DepCity__updated_at`,`ArrCity`.`id` AS `ArrCity__id`,`ArrCity`.`name` AS `ArrCity__name`,`ArrCity`.`created_at` AS `ArrCity__created_at`,`ArrCity`.`updated_at` AS `ArrCity__updated_at` FROM `flights` LEFT JOIN `cities` `DepCity` ON `flights`.`dep_city_id` = `DepCity`.`id` LEFT JOIN `cities` `ArrCity` ON `flights`.`arr_city_id` = `ArrCity`.`id` WHERE DepCity.name = 'Dallas' AND ArrCity.name = 'Philadelphia' AND year(dep_time) = 2020 AND month(dep_time) = 11 AND day(dep_time) = 24"
	log.Print(reqStr)
	suite.sqlMock.ExpectQuery("SELECT").
		RowsWillBeClosed().
		WillReturnRows(rows)

	res, err := suite.CallHandler("/flights?departure_city=Dallas&arrival_city=Philadelphia&departure_time=2020-11-24T00:00:00Z")
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())
}

func TestGetFlights(t *testing.T) {
	suite.Run(t, new(GetFlightsTestSuite))
}
