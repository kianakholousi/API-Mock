package handler

import (
	"github.com/DATA-DOG/go-sqlmock"
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

type GetAirplanesTestSuite struct {
	suite.Suite
	sqlMock   sqlmock.Sqlmock
	e         *echo.Echo
	airplanes Airplane
	timeMock  time.Time
}

func (suite *GetAirplanesTestSuite) SetupSuite() {
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
	suite.airplanes = Airplane{
		DB: db,
	}
	suite.timeMock = time.Date(2020, time.January, 1, 2, 3, 0, 0, time.UTC)
}

func (suite *GetAirplanesTestSuite) CallHandler() (*httptest.ResponseRecorder, error) {
	req := httptest.NewRequest(http.MethodGet, "/airplanes", strings.NewReader(""))
	res := httptest.NewRecorder()
	c := suite.e.NewContext(req, res)
	err := suite.airplanes.Get(c)

	return res, err
}

func (suite *GetAirplanesTestSuite) TestGetAirplanes_Success() {
	require := suite.Require()
	expectedStatusCode := http.StatusOK
	expectedMsg := `[{"id":1,"name":"AirbusA320"},{"id":2,"name":"Boeing737"}]`

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(1, "AirbusA320", suite.timeMock, suite.timeMock).
		AddRow(2, "Boeing737", suite.timeMock, suite.timeMock)

	suite.sqlMock.ExpectQuery("^SELECT \\* FROM `airplanes`$").
		RowsWillBeClosed().
		WillReturnRows(rows)

	res, err := suite.CallHandler()
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())
}

func (suite *GetAirplanesTestSuite) TestGetAirplanes_Database_Failure() {
	require := suite.Require()
	expectedStatusCode := http.StatusInternalServerError
	expectedMsg := "\"Internal Server Error\"\n"

	res, err := suite.CallHandler()
	require.NoError(err)
	require.Equal(expectedStatusCode, res.Code)
	require.JSONEq(expectedMsg, res.Body.String())
}

func TestGetAirplanes(t *testing.T) {
	suite.Run(t, new(GetAirplanesTestSuite))
}
