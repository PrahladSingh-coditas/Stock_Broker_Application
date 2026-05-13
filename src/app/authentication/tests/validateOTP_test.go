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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const queryForOTP = "SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?"

func validateOTPSQLMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, _ := sqlmock.New()

	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})

	gdb, _ := gorm.Open(dialector, &gorm.Config{})
	return gdb, mock
}

func GetValidateOTPRouter(gdb *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	validateOTPRepo := repository.NewValidateUserOtpRepository(gdb)
	validateOTPService := business.NewValidateUserOtpService(validateOTPRepo)
	validateOTPHandler := handlers.NewValidateUserOtpHandler(validateOTPService)
	r.POST("/api/auth/validateotp", validateOTPHandler.HandleValidateUserOtp)
	return r
}

type ValidateOTPTestSuite struct {
	suite.Suite
}

func (s *ValidateOTPTestSuite) TestValidateOTP_Success() {
	t := s.T()

	gdb, mock := validateOTPSQLMock(t)
	rows := sqlmock.NewRows([]string{"otpSent", "otpExpiresAt", "username"}).AddRow("1234", 1979587566, "Sanjana")
	mock.ExpectQuery(regexp.QuoteMeta(queryForOTP)).WithArgs("Sanjana", 1).WillReturnRows(rows)

	r := GetValidateOTPRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/validateotp", bytes.NewBufferString(`{"otp":"1234","username":"Sanjana"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Token Generated Successfully")
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ValidateOTPTestSuite) TestValidateOTP_DatabaseError() {
	t := s.T()

	gdb, mock := validateOTPSQLMock(t)
	mock.ExpectQuery(regexp.QuoteMeta(queryForOTP)).WithArgs("Sanjana", 1).WillReturnError(errors.New("Database Error"))

	r := GetValidateOTPRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/validateotp", bytes.NewBufferString(`{"otp":"1234","username":"Sanjana"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.Contains(w.Body.String(), "Signin Failed")
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ValidateOTPTestSuite) TestValidateOTP_IncorrectOTP() {
	t := s.T()

	gdb, mock := validateOTPSQLMock(t)
	rows := sqlmock.NewRows([]string{"otpSent", "otpExpiresAt", "username"}).AddRow("1234", 1979587566, "Sanjana")
	mock.ExpectQuery(regexp.QuoteMeta(queryForOTP)).WithArgs("Sanjana", 1).WillReturnRows(rows)

	r := GetValidateOTPRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/validateotp", bytes.NewBufferString(`{"otp":"1235","username":"Sanjana"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
	s.Contains(w.Body.String(), "Incorrect OTP")
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ValidateOTPTestSuite) TestValidateOTP_ExpiredOTP() {
	t := s.T()

	gdb, mock := validateOTPSQLMock(t)
	rows := sqlmock.NewRows([]string{"otpSent", "otpExpiresAt", "username"}).AddRow("1234", 0, "Sanjana")
	mock.ExpectQuery(regexp.QuoteMeta(queryForOTP)).WithArgs("Sanjana", 1).WillReturnRows(rows)

	r := GetValidateOTPRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/validateotp", bytes.NewBufferString(`{"otp":"1234","username":"Sanjana"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
	s.Contains(w.Body.String(), "OTP Expired")
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ValidateOTPTestSuite) TestValidateOTP_InvalidPayload() {
	t := s.T()

	gdb, mock := validateOTPSQLMock(t)

	r := GetValidateOTPRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/validateotp", bytes.NewBufferString(`{"otp":1234,"username":"Sanjana"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "invalid required payload")
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ValidateOTPTestSuite) TestValidateOTP_UserNotFoundError() {
	t := s.T()

	gdb, mock := validateOTPSQLMock(t)
	mock.ExpectQuery(regexp.QuoteMeta(queryForOTP)).WithArgs("Sanjana", 1).WillReturnError(gorm.ErrRecordNotFound)

	r := GetValidateOTPRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/validateotp", bytes.NewBufferString(`{"otp":"1234","username":"Sanjana"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusNotFound, w.Code)
	s.Contains(w.Body.String(), "user not found")
	s.NoError(mock.ExpectationsWereMet())
}

func (s *ValidateOTPTestSuite) TestValidateOTP_ValidationError() {
	t := s.T()

	gdb, mock := validateOTPSQLMock(t)

	r := GetValidateOTPRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/validateotp", bytes.NewBufferString(`{"otp":"123","username":"Sanjana"}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "otp must be exactly 4 digits and only numeric")
	s.NoError(mock.ExpectationsWereMet())
}

func TestValidateOTPTestSuite(t *testing.T) {
	if err := utils.InitJWTConfig("../../../config"); err != nil {
		log.Fatalf(constants.ErrJWTConfigReadFailed, err)
	}
	suite.Run(t, new(ValidateOTPTestSuite))
}
