package tests

import (
	"authentication/business"
	"authentication/commons/constants"
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

const insertUserQuery = "INSERT INTO `users` (`username`,`password`,`panCard`,`phoneNumber`,`email`) VALUES (?,?,?,?,?)"

func signupSQLMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, _ := sqlmock.New()

	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})

	gdb, _ := gorm.Open(dialector, &gorm.Config{})
	return gdb, mock
}

func GetSignupRouter(gdb *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	signupRepo := repository.NewCreateUserRepository(gdb)
	signupService := business.NewCreateUserService(signupRepo)
	signupHandler := handlers.NewCreateUserHandler(signupService)
	r.POST("/api/auth/signup", signupHandler.HandleCreaterUser)
	return r
}

type SignupTestSuite struct {
	suite.Suite
}

func (s *SignupTestSuite) TestSignup_Success() {
	gdb, mock := signupSQLMock(s.T())

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(insertUserQuery)).
		WithArgs(
			"Dinesh",
			sqlmock.AnyArg(),
			"AABCD1234F",
			9876543210,
			"dinesh@gmail.com",
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	r := GetSignupRouter(gdb)

	body := `{
		"username":"Dinesh",
		"password":"Dinesh@123",
		"confirmPassword":"Dinesh@123",
		"panCard":"AABCD1234F",
		"phoneNumber":9876543210,
		"email":"dinesh@gmail.com"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusCreated, w.Code)
	s.Contains(w.Body.String(), "User created successfully")
	s.NoError(mock.ExpectationsWereMet())
}

func (s *SignupTestSuite) TestSignup_InvalidRequest() {
	gdb, mock := signupSQLMock(s.T())
	r := GetSignupRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup",
		bytes.NewBufferString(`{"username":123}`))

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	s.NoError(mock.ExpectationsWereMet())
}

func (s *SignupTestSuite) TestSignup_ValidationError() {
	gdb, mock := signupSQLMock(s.T())
	r := GetSignupRouter(gdb)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup",
		bytes.NewBufferString(`{"username":"Dinesh"}`))

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	s.NoError(mock.ExpectationsWereMet())
}

func (s *SignupTestSuite) TestSignup_DuplicateUsername() {
	gdb, mock := signupSQLMock(s.T())

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(insertUserQuery)).
		WillReturnError(errors.New(constants.ErrUniqueConstraintViolation))

	mock.ExpectRollback()

	r := GetSignupRouter(gdb)

	body := `{
			"username":"Dinesh",
			"password":"Dinesh@123",
			"confirmPassword":"Dinesh@123",
			"panCard":"AABCD1234F",
			"phoneNumber":9876543210,
			"email":"dinesh@gmail.com"
		}`

	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusConflict, w.Code)
	s.NoError(mock.ExpectationsWereMet())
}

func (s *SignupTestSuite) TestSignup_DuplicatePanCard() {
	gdb, mock := signupSQLMock(s.T())

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(insertUserQuery)).
		WillReturnError(errors.New("duplicate key value violates unique constraint : idx_users_pan_card"))

	mock.ExpectRollback()

	r := GetSignupRouter(gdb)

	body := `{
			"username":"Sanjana",
			"password":"Sanjana@123",
			"confirmPassword":"Sanjana@123",
			"panCard":"AABCD1234F",
			"phoneNumber":9906543210,
			"email":"sanjana@gmail.com"
		}`

	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusConflict, w.Code)
	s.NoError(mock.ExpectationsWereMet())
}

func (s *SignupTestSuite) TestSignup_DuplicateEmail() {
	gdb, mock := signupSQLMock(s.T())

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(insertUserQuery)).
		WillReturnError(errors.New("duplicate key value violates unique constraint : idx_users_email"))

	mock.ExpectRollback()

	r := GetSignupRouter(gdb)

	body := `{
			"username":"Sanjana",
			"password":"Sanjana@123",
			"confirmPassword":"Sanjana@123",
			"panCard":"ABCDE1234F",
			"phoneNumber":9906543210,
			"email":"dinesh@gmail.com"
		}`

	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusConflict, w.Code)
	s.NoError(mock.ExpectationsWereMet())
}

func (s *SignupTestSuite) TestSignup_InternalError() {
	gdb, mock := signupSQLMock(s.T())

	mock.ExpectBegin()

	mock.ExpectExec(regexp.QuoteMeta(insertUserQuery)).
		WillReturnError(errors.New("Database Query Error"))

	mock.ExpectRollback()

	r := GetSignupRouter(gdb)

	body := `{
			"username":"Dinesh",
			"password":"Dinesh@123",
			"confirmPassword":"Dinesh@123",
			"panCard":"AABCD1234F",
			"phoneNumber":9876543210,
			"email":"dinesh@gmail.com"
		}`

	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.Contains(w.Body.String(), "failed to create user")
	s.NoError(mock.ExpectationsWereMet())
}

func TestSignupTestSuite(t *testing.T) {
	suite.Run(t, new(SignupTestSuite))
}
