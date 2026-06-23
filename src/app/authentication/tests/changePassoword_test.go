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

const queryForUsername = "SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?"
const queryForUpdatePassword = "UPDATE `users` SET `password`=? WHERE username = ?"

func changePasswordSQLMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, _ := sqlmock.New()

	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})

	gdb, _ := gorm.Open(dialector, &gorm.Config{})
	return gdb, mock
}

func GetChangePasswordRouter(gdb *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	changePasswordRepo := repository.NewChangePasswordRepository(gdb)
	changePasswordService := business.NewChangePasswordService(changePasswordRepo)
	changePasswordHandler := handlers.NewChangePasswordHandler(changePasswordService)
	r.POST("/api/auth/changepassword", func(ctx *gin.Context) {
		ctx.Set("username", "Sanjana")
		ctx.Next()
	}, changePasswordHandler.HandleChangePassword)
	return r
}

type ChangePasswordTestSuite struct {
	suite.Suite
}

func (s *ChangePasswordTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *ChangePasswordTestSuite) TestMockChangePassword200Success() {
	t := s.T()

	gdb, mock := changePasswordSQLMock(t)
	rows := sqlmock.NewRows([]string{"username", "password"}).AddRow("Sanjana", "Sanjana@123")
	mock.ExpectQuery(regexp.QuoteMeta(queryForUsername)).WithArgs("Sanjana", 1).WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryForUpdatePassword)).WithArgs(sqlmock.AnyArg(), "Sanjana").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	r := GetChangePasswordRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/changepassword", bytes.NewBufferString(`{"confirmPassword":"SanjanaS@123","newpassword":"SanjanaS@123"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Password Changed Successfully")
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ChangePasswordTestSuite) TestMockChangePassword404UserNotFound() {
	t := s.T()

	gdb, mock := changePasswordSQLMock(t)
	mock.ExpectQuery(regexp.QuoteMeta(queryForUsername)).WithArgs("Sanjana", 1).WillReturnError(gorm.ErrRecordNotFound)

	r := GetChangePasswordRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/changepassword", bytes.NewBufferString(`{"confirmPassword":"SanjanaS@123","newpassword":"SanjanaS@123"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), "user not found")
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ChangePasswordTestSuite) TestMockChangePassword404UserNotFoundForUpdate() {
	t := s.T()

	gdb, mock := changePasswordSQLMock(t)
	rows := sqlmock.NewRows([]string{"username", "password"}).AddRow("Sanjana", "Sanjana@123")
	mock.ExpectQuery(regexp.QuoteMeta(queryForUsername)).WithArgs("Sanjana", 1).WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryForUpdatePassword)).WithArgs(sqlmock.AnyArg(), "Sanjana").WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectRollback()

	r := GetChangePasswordRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/changepassword", bytes.NewBufferString(`{"confirmPassword":"SanjanaS@123","newpassword":"SanjanaS@123"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), "user not found")
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ChangePasswordTestSuite) TestMockChangePassword500DatabaseErrorForGetUsername() {
	t := s.T()

	gdb, mock := changePasswordSQLMock(t)
	mock.ExpectQuery(regexp.QuoteMeta(queryForUsername)).WithArgs("Sanjana", 1).WillReturnError(errors.New("Databse Error"))

	r := GetChangePasswordRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/changepassword", bytes.NewBufferString(`{"confirmPassword":"SanjanaS@123","newpassword":"SanjanaS@123"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.Contains(w.Body.String(), "authentication failed")
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ChangePasswordTestSuite) TestMockChangePassword401ChangePasswordFailed() {
	t := s.T()

	gdb, mock := changePasswordSQLMock(t)
	rows := sqlmock.NewRows([]string{"username", "password"}).AddRow("Sanjana", "Sanjana@123")
	mock.ExpectQuery(regexp.QuoteMeta(queryForUsername)).WithArgs("Sanjana", 1).WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryForUpdatePassword)).WithArgs(sqlmock.AnyArg(), "Sanjana").WillReturnError(errors.New("failed to change password: %s"))
	mock.ExpectRollback()

	r := GetChangePasswordRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/changepassword", bytes.NewBufferString(`{"confirmPassword":"SanjanaS@123","newpassword":"SanjanaS@123"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
	s.Contains(w.Body.String(), "failed to change password: %s")
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ChangePasswordTestSuite) TestMockChangePassword401SameAsOldPasswordError() {
	t := s.T()

	gdb, mock := changePasswordSQLMock(t)
	rows := sqlmock.NewRows([]string{"username", "password"}).AddRow("Sanjana", "Sanjana@123")
	mock.ExpectQuery(regexp.QuoteMeta(queryForUsername)).WithArgs("Sanjana", 1).WillReturnRows(rows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryForUpdatePassword)).WithArgs(sqlmock.AnyArg(), "Sanjana").WillReturnError(errors.New("new password cannot be the same as the old password"))
	mock.ExpectRollback()

	r := GetChangePasswordRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/changepassword", bytes.NewBufferString(`{"confirmPassword":"Sanjana@123","newpassword":"Sanjana@123"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
	s.Contains(w.Body.String(), "new password cannot be the same as the old password")
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ChangePasswordTestSuite) TestChangePassword400ValidationError() {
	t := s.T()

	gdb, mock := changePasswordSQLMock(t)
	r := GetChangePasswordRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/changepassword", bytes.NewBufferString(`{"confirmPassword":"SanjanaS@12","newpassword":"SanjanaS@123"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "ConfirmPassword must match Password.")
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ChangePasswordTestSuite) TestChangePassword400InvalidRequest() {
	t := s.T()

	gdb, mock := changePasswordSQLMock(t)
	r := GetChangePasswordRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/changepassword", bytes.NewBufferString(`{"confirmPassword":123,"newpassword":"SanjanaS@123"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "invalid required payload")
	s.NoError(mock.ExpectationsWereMet())
}

func TestChangePasswordTestSuite(t *testing.T) {
	suite.Run(t, new(ChangePasswordTestSuite))
}
