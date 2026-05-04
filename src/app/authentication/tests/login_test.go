package tests

import (
	"authentication/business"
	"authentication/handlers"
	"authentication/repository"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var queryUserByUsername = "SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?"
var queryUpdateOtp = "UPDATE `users` SET `otpExpiresAt`=?,`otpSent`=? WHERE username = ?"

func getMockSQLDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })
	dialector := mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true})
	gdb, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)
	return gdb, mock
}

func getLoginRouter(gDB *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	repo := repository.NewSignInUserRepository(gDB)
	service := business.NewSignInUserService(repo)
	handler := handlers.NewSignInUserHandler(service)

	router.POST("/api/auth/signin", handler.HandleSignInUser)

	return router
}

type SignInUserTestSuite struct {
	suite.Suite
}

func (suite *SignInUserTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (suite *SignInUserTestSuite) TestMockSignInUser200SignInSuccessfull() {
	t := suite.T()

	gDB, mock := getMockSQLDB(t)

	user := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(1, "Arijit", "$2a$10$u1PjmxE.WsjEVKKis8qCVOghE2CQ0KFLuh1Kbes5BbYYeq5nZidZu")
	mock.ExpectQuery(regexp.QuoteMeta(queryUserByUsername)).WithArgs("Arijit", 1).WillReturnRows(user)

	otp := sqlmock.AnyArg()
	otpTime := sqlmock.AnyArg()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryUpdateOtp)).WithArgs(otpTime, otp, "Arijit").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := getLoginRouter(gDB)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBufferString(`{"username":"Arijit","password":"Secure@123"}`))
	request.Header.Set("Content-type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "OTP sent successfully")
}

func (suite *SignInUserTestSuite) TestMockSignInUser400InvalidPayload() {
	t := suite.T()

	gDB, _ := getMockSQLDB(t)

	router := getLoginRouter(gDB)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBufferString(`{"username":123}`))
	request.Header.Set("Content-type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "invalid required payload")
}

func (suite *SignInUserTestSuite) TestMockSignInUser404UserNotFound() {
	t := suite.T()

	gDB, mock := getMockSQLDB(t)

	//rows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(1, "Arijit", "$10$u1PjmxE.WsjEVKKis8qCVOghE2CQ0KFLuh1Kbes5BbYYeq5nZidZu")

	mock.ExpectQuery(regexp.QuoteMeta(queryUserByUsername)).WithArgs("Dharmesh", 1).WillReturnError(gorm.ErrRecordNotFound)
	router := getLoginRouter(gDB)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBufferString(`{"username":"Dharmesh","password":"Secure@123"}`))
	request.Header.Set("Content-type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Contains(w.Body.String(), "user not found")
}

func (suite *SignInUserTestSuite) TestMockSignInUser400ValidationFailed() {
	t := suite.T()

	gDB, mock := getMockSQLDB(t)

	rows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(1, "Arijit", "$10$u1PjmxE.WsjEVKKis8qCVOghE2CQ0KFLuh1Kbes5BbYYeq5nZidZu")

	mock.ExpectQuery(regexp.QuoteMeta(queryUserByUsername)).WithArgs("Arijit", 1).WillReturnRows(rows)
	router := getLoginRouter(gDB)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBufferString(`{"username":"","password":"Secure@123"}`))
	request.Header.Set("Content-type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *SignInUserTestSuite) TestMockSignInUser401WrongPassword() {
	t := suite.T()

	gDB, mock := getMockSQLDB(t)

	rows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(1, "Arijit", "$10$u1PjmxE.WsjEVKKis8qCVOghE2CQ0KFLuh1Kbes5BbYYeq5nZidZu")

	mock.ExpectQuery(regexp.QuoteMeta(queryUserByUsername)).WithArgs("Arijit", 1).WillReturnRows(rows)
	router := getLoginRouter(gDB)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBufferString(`{"username":"Arijit","password":"NotSecure@123"}`))
	request.Header.Set("Content-type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusUnauthorized, w.Code)
	suite.Contains(w.Body.String(), "password not matched")
}

func (suite *SignInUserTestSuite) TestMockSignInUser500DatabaseQueryError() {
	t := suite.T()

	gDB, mock := getMockSQLDB(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryUserByUsername)).WithArgs("Arijit", 1).WillReturnError(fmt.Errorf("database query error"))
	router := getLoginRouter(gDB)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBufferString(`{"username":"Arijit","password":"Dharmesh@123"}`))
	request.Header.Set("Content-type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "database query error")
}

func (suite *SignInUserTestSuite) TestMockSignInUser500DatabaseQueryErrorWhileStoreOtp() {
	t := suite.T()

	gDB, mock := getMockSQLDB(t)

	user := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(1, "Arijit", "$2a$10$u1PjmxE.WsjEVKKis8qCVOghE2CQ0KFLuh1Kbes5BbYYeq5nZidZu")
	mock.ExpectQuery(regexp.QuoteMeta(queryUserByUsername)).WithArgs("Arijit", 1).WillReturnRows(user)

	otp := sqlmock.AnyArg()
	otpTime := sqlmock.AnyArg()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryUpdateOtp)).WithArgs(otpTime, otp, "Arijit").WillReturnError(errors.New("database query error"))
	mock.ExpectCommit()

	router := getLoginRouter(gDB)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBufferString(`{"username":"Arijit","password":"Secure@123"}`))
	request.Header.Set("Content-type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "database query error")
}

func (suite *SignInUserTestSuite) TestMockSignInUser404UserNotFoundWhileStoreOtp() {
	t := suite.T()

	gDB, mock := getMockSQLDB(t)

	user := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(1, "Arijit", "$2a$10$u1PjmxE.WsjEVKKis8qCVOghE2CQ0KFLuh1Kbes5BbYYeq5nZidZu")
	mock.ExpectQuery(regexp.QuoteMeta(queryUserByUsername)).WithArgs("Arijit", 1).WillReturnRows(user)

	otp := sqlmock.AnyArg()
	otpTime := sqlmock.AnyArg()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryUpdateOtp)).WithArgs(otpTime, otp, "Arijit").WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectRollback()

	router := getLoginRouter(gDB)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBufferString(`{"username":"Arijit","password":"Secure@123"}`))
	request.Header.Set("Content-type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Contains(w.Body.String(), "user not found")
}

func TestMockSignInUser(t *testing.T) {
	suite.Run(t, new(SignInUserTestSuite))
}
