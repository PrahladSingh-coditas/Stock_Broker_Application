package tests

import (
	"authentication/business"
	"authentication/handlers"
	"authentication/repository"
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"stock_broker_application/src/constants"
	"stock_broker_application/src/utils"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

var queryGetUserByUsername = "SELECT * FROM `users` WHERE username = ?  ORDER BY `users`.`id` LIMIT ?"

func getValidateOtpRouter(gDB *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	repo := repository.NewValidateUserOtpRepository(gDB)
	service := business.NewValidateUserOtpService(repo)
	handler := handlers.NewValidateUserOtpHandler(service)

	router.POST("/api/auth/validate-otp", handler.HandleValidateUserOtp)

	return router
}

type ValidateOtpTestSuite struct {
	suite.Suite
}

func (suite *ValidateOtpTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func TestMockValidateOtp(t *testing.T) {
	if err := utils.InitJWTConfig("../../../config"); err != nil {
		log.Fatalf(constants.ErrJWTConfigReadFailed, err)
	}
	suite.Run(t, new(ValidateOtpTestSuite))
}

func (suite *ValidateOtpTestSuite) TestMockValidateOtp200ValidatedSuccessfully() {
	t := suite.T()

	gDB, mock := getMockSQLDB(t)

	//change this time to future time always to pass this test case
	rows := sqlmock.NewRows([]string{"id", "username", "otpSent", "otpExpiresAt"}).AddRow(1, "Arijit", 1234, 1797961756)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetUserByUsername)).WithArgs("Arijit", 1).WillReturnRows(rows)

	router := getValidateOtpRouter(gDB)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/validate-otp", bytes.NewBufferString(`{"username":"Arijit", "otp":"1234"}`))
	request.Header.Set("Content-type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "OTP validated successfully")
	suite.NoError(mock.ExpectationsWereMet())
}

func (suite *ValidateOtpTestSuite) TestMockValidateOtp500InternalServerError() {
	t := suite.T()

	gDB, mock := getMockSQLDB(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetUserByUsername)).WithArgs("Arijit", 1).WillReturnError(errors.New("internal server error"))

	router := getValidateOtpRouter(gDB)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/validate-otp", bytes.NewBufferString(`{"username":"Arijit", "otp":"1234"}`))
	request.Header.Set("Content-type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "failed to sign in user")
	suite.NoError(mock.ExpectationsWereMet())
}

func (suite *ValidateOtpTestSuite) TestMockValidateOtp404UserNotFoundError() {
	t := suite.T()

	gDB, mock := getMockSQLDB(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetUserByUsername)).WithArgs("Arijit", 1).WillReturnError(gorm.ErrRecordNotFound)

	router := getValidateOtpRouter(gDB)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/validate-otp", bytes.NewBufferString(`{"username":"Arijit", "otp":"1234"}`))
	request.Header.Set("Content-type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Contains(w.Body.String(), "user not found")
	suite.NoError(mock.ExpectationsWereMet())
}

func (suite *ValidateOtpTestSuite) TestMockValidateOtp401IncorrectOtpError() {
	t := suite.T()

	gDB, mock := getMockSQLDB(t)

	//change this time to future time always to pass this test case
	rows := sqlmock.NewRows([]string{"id", "username", "otpSent", "otpExpiresAt"}).AddRow(1, "Arijit", 4321, 1797961756)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetUserByUsername)).WithArgs("Arijit", 1).WillReturnRows(rows)

	router := getValidateOtpRouter(gDB)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/validate-otp", bytes.NewBufferString(`{"username":"Arijit", "otp":"1234"}`))
	request.Header.Set("Content-type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Contains(w.Body.String(), "entered OTP is not correct")
	suite.NoError(mock.ExpectationsWereMet())
}

func (suite *ValidateOtpTestSuite) TestMockValidateOtp401ExpiredOtpError() {
	t := suite.T()

	gDB, mock := getMockSQLDB(t)

	//keep this time always same to pass this test case
	rows := sqlmock.NewRows([]string{"id", "username", "otpSent", "otpExpiresAt"}).AddRow(1, "Arijit", 1234, 0)

	mock.ExpectQuery(regexp.QuoteMeta(queryGetUserByUsername)).WithArgs("Arijit", 1).WillReturnRows(rows)

	router := getValidateOtpRouter(gDB)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/validate-otp", bytes.NewBufferString(`{"username":"Arijit", "otp":"1234"}`))
	request.Header.Set("Content-type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Contains(w.Body.String(), "OTP expired")
	suite.NoError(mock.ExpectationsWereMet())
}

func (suite *ValidateOtpTestSuite) TestMockValidateOtp400InvalidPayloadError() {
	t := suite.T()

	gDB, mock := getMockSQLDB(t)

	router := getValidateOtpRouter(gDB)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/validate-otp", bytes.NewBufferString(`{"username":123, "otp":"1234"}`))
	request.Header.Set("Content-type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "invalid required payload")
	suite.NoError(mock.ExpectationsWereMet())
}

func (suite *ValidateOtpTestSuite) TestMockValidateOtp400ValidationError() {
	t := suite.T()

	gDB, mock := getMockSQLDB(t)

	router := getValidateOtpRouter(gDB)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/validate-otp", bytes.NewBufferString(`{"username":"123", "otp":"1234"}`))
	request.Header.Set("Content-type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.NoError(mock.ExpectationsWereMet())
}