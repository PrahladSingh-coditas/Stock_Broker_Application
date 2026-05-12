package tests

import (
	"authentication/business"
	"authentication/handlers"
	"authentication/repository"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const queryByDetails = "SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?"
const queryForUpdateOTP = "UPDATE `users` SET `otpExpiresAt`=?,`otpSent`=? WHERE username = ?"

func forgotPasswordSQLMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, _ := sqlmock.New()

	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})

	gdb, _ := gorm.Open(dialector, &gorm.Config{})
	return gdb, mock
}

func GetForgotPasswordRouter(gdb *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	forgotPasswordRepo := repository.NewForgotPasswordRepository(gdb)
	forgotPasswordService := business.NewForgotPasswordService(forgotPasswordRepo)
	forgotPasswordHandler := handlers.NewForgotPasswordHandler(forgotPasswordService)
	r.POST("/api/auth/forgotpassword", forgotPasswordHandler.HandleForgotPassword)
	return r
}

type ForgotPasswordTestSuite struct {
	suite.Suite
}

func (s *ForgotPasswordTestSuite) TestForgotPassword_Success() {
	t := s.T()

	gdb, mock := forgotPasswordSQLMock(t)
	rows := sqlmock.NewRows([]string{"panCard", "phoneNumber", "username"}).AddRow("ABCDE1234F", 8432805566, "Sanjana")
	mock.ExpectQuery(regexp.QuoteMeta(queryByDetails)).WithArgs("Sanjana", 1).WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryForUpdateOTP)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "Sanjana").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	r := GetForgotPasswordRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/forgotpassword", bytes.NewBufferString(`{"panCard":"ABCDE1234F","phoneNumber":8432805566,"username":"Sanjana"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Forgot password OTP sent successfully, will expire in 2 minutes")
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ForgotPasswordTestSuite) TestForgotPassword_UserNotFound() {
	t := s.T()
	gdb, mock := forgotPasswordSQLMock(t)
	mock.ExpectQuery(regexp.QuoteMeta(queryByDetails)).WithArgs("Sakshi", 1).WillReturnError(gorm.ErrRecordNotFound)

	r := GetForgotPasswordRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/forgotpassword", bytes.NewBufferString(`{"panCard":"ABCDE1234F","phoneNumber":8432805566,"username":"Sakshi"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	s.Equal(http.StatusNotFound, w.Code)
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ForgotPasswordTestSuite) TestMockForgotPassword404UserNotFoundForGenerateOTP() {
	t := s.T()
	gdb, mock := forgotPasswordSQLMock(t)
	rows := sqlmock.NewRows([]string{"panCard", "phoneNumber", "username"}).AddRow("ABCDE1234F", 8432805566, "Sanjana")
	mock.ExpectQuery(regexp.QuoteMeta(queryByDetails)).WithArgs("Sanjana", 1).WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryForUpdateOTP)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "Sanjana").WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectRollback()

	r := GetForgotPasswordRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/forgotpassword", bytes.NewBufferString(`{"panCard":"ABCDE1234F","phoneNumber":8432805566,"username":"Sanjana"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	s.Equal(http.StatusNotFound, w.Code)
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ForgotPasswordTestSuite) TestForgotPassword_AuthenticationFailed() {
	t := s.T()
	gdb, mock := forgotPasswordSQLMock(t)
	mock.ExpectQuery(regexp.QuoteMeta(queryByDetails)).WithArgs("Sanjana", 1).WillReturnError(errors.New("Database Error"))

	r := GetForgotPasswordRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/forgotpassword", bytes.NewBufferString(`{"panCard":"ABCDE1234F","phoneNumber":8432805566,"username":"Sanjana"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	s.Equal(http.StatusUnauthorized, w.Code)
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ForgotPasswordTestSuite) TestMockForgotPassword401Unauthorized() {
	t := s.T()
	gdb, mock := forgotPasswordSQLMock(t)
	rows := sqlmock.NewRows([]string{"panCard", "phoneNumber", "username"}).AddRow("ABCDE1234F", 8432805566, "Sanjana")
	mock.ExpectQuery(regexp.QuoteMeta(queryByDetails)).WithArgs("Sanjana", 1).WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryForUpdateOTP)).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), "Sanjana").WillReturnError(errors.New("Database Error"))
	mock.ExpectRollback()

	r := GetForgotPasswordRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/forgotpassword", bytes.NewBufferString(`{"panCard":"ABCDE1234F","phoneNumber":8432805566,"username":"Sanjana"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	s.Equal(http.StatusUnauthorized, w.Code)
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ForgotPasswordTestSuite) TestForgotPassword_InvalidRequest() {
	t := s.T()
	gdb, mock := forgotPasswordSQLMock(t)
	r := GetForgotPasswordRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/forgotpassword", bytes.NewBufferString(`{"username":123}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	s.Equal(http.StatusBadRequest, w.Code)
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ForgotPasswordTestSuite) TestForgotPassword_ValidationError() {
	t := s.T()
	gdb, mock := forgotPasswordSQLMock(t)
	r := GetForgotPasswordRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/forgotpassword", bytes.NewBufferString(`{"username":"Sanjana"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	s.Equal(http.StatusBadRequest, w.Code)
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ForgotPasswordTestSuite) TestMockForgotPassword401InvalidCredentials() {
	t := s.T()
	gdb, mock := forgotPasswordSQLMock(t)
	rows := sqlmock.NewRows([]string{"panCard", "phoneNumber", "username"}).AddRow("ABCDE1234D", 8432805566, "Sanjana")
	mock.ExpectQuery(regexp.QuoteMeta(queryByDetails)).WithArgs("Sanjana", 1).WillReturnRows(rows)

	r := GetForgotPasswordRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/forgotpassword", bytes.NewBufferString(`{"panCard":"ABCDE1234F","phoneNumber":8432805566,"username":"Sanjana"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	s.Equal(http.StatusUnauthorized, w.Code)
	s.NoError(mock.ExpectationsWereMet())
}

func TestForgotPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(ForgotPasswordTestSuite))
}
