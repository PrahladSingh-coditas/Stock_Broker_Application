package tests

// import (
// 	"authentication/business"
// 	"authentication/handlers"
// 	"authentication/repository"
// 	"bytes"
// 	"database/sql"
// 	"net/http"
// 	"net/http/httptest"
// 	"regexp"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/gin-gonic/gin"
// 	"github.com/go-openapi/testify/v2/require"
// 	"github.com/stretchr/testify/suite"
// 	"golang.org/x/crypto/bcrypt"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// const queryByUsername = "SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?"

// func signinSQLMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
// 	t.Helper()
// 	sqlDB, mock, err := sqlmock.New()
// 	require.NoError(t, err)
// 	t.Cleanup(func() { _ = sqlDB.Close() })
// 	dialector := mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true})
// 	gdb, err := gorm.Open(dialector, &gorm.Config{})
// 	require.NoError(t, err)
// 	return gdb, mock
// }

// func GetSigninRouter(gdb *gorm.DB) *gin.Engine {
// 	gin.SetMode(gin.TestMode)
// 	r := gin.New()
// 	signinRepo := repository.NewSigninUserRepository(gdb)
// 	signinService := business.NewSigninUserService(signinRepo)
// 	signinHandler := handlers.NewSigninUserHandler(signinService)
// 	r.POST("/api/auth/signin", signinHandler.HandleSigninUser)
// 	return r
// }

// type SigninTestSuite struct {
// 	suite.Suite
// }

// func (s *SigninTestSuite) SetupTest() {
// 	gin.SetMode(gin.TestMode)
// }

// func (s *SigninTestSuite) TestSignin_Success() {
// 	t := s.T()

// 	gdb, mock := signinSQLMock(t)
// 	hashPassword, _ := bcrypt.GenerateFromPassword([]byte("Sanjana@123"), bcrypt.DefaultCost)
// 	rows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(1, "Sanjana", string(hashPassword))
// 	mock.ExpectQuery(regexp.QuoteMeta(queryByUsername)).WithArgs("Sanjana", 1).WillReturnRows(rows)
// 	r := GetSigninRouter(gdb)

// 	req := httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBufferString(`{"username":"Sanjana","password":"Sanjana@123"}`))
// 	req.Header.Set("Content-type", "application/json")
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	s.Equal(http.StatusOK, w.Code)
// 	s.Contains(w.Body.String(), "User logged in successfully")
// 	s.NoError(mock.ExpectationsWereMet())
// }

// func (s *SigninTestSuite) TestSignin_InvalidRequest() {
// 	t := s.T()
// 	gdb, mock := signinSQLMock(t)
// 	r := GetSigninRouter(gdb)

// 	req := httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBufferString(`{"username":123}`))
// 	req.Header.Set("Content-type", "application/json")
// 	w := httptest.NewRecorder()

// 	r.ServeHTTP(w, req)
// 	s.Equal(http.StatusBadRequest, w.Code)
// 	s.NoError(mock.ExpectationsWereMet())
// }

// func (s *SigninTestSuite) TestSignin_InvalidPassword() {
// 	t := s.T()
// 	gdb, mock := signinSQLMock(t)
// 	hashPassword, _ := bcrypt.GenerateFromPassword([]byte("Sanjana@123"), bcrypt.DefaultCost)
// 	rows := sqlmock.NewRows([]string{"id", "username", "password"}).AddRow(1, "Sanjana", string(hashPassword))
// 	mock.ExpectQuery(regexp.QuoteMeta(queryByUsername)).WithArgs("Sanjana", 1).WillReturnRows(rows)
// 	r := GetSigninRouter(gdb)

// 	req := httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBufferString(`{"username":"Sanjana","password":"Sanjana@12345"}`))
// 	req.Header.Set("Content-type", "application/json")
// 	w := httptest.NewRecorder()

// 	r.ServeHTTP(w, req)
// 	s.Equal(http.StatusUnauthorized, w.Code)
// 	s.NoError(mock.ExpectationsWereMet())
// }

// func (s *SigninTestSuite) TestSignin_UserNotFound() {
// 	t := s.T()
// 	gdb, mock := signinSQLMock(t)
// 	mock.ExpectQuery(regexp.QuoteMeta(queryByUsername)).WithArgs("Sakshi", 1).WillReturnError(gorm.ErrRecordNotFound)
// 	r := GetSigninRouter(gdb)

// 	req := httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBufferString(`{"username":"Sakshi","password":"Sakshi@123"}`))
// 	req.Header.Set("Content-type", "application/json")
// 	w := httptest.NewRecorder()

// 	r.ServeHTTP(w, req)
// 	s.Equal(http.StatusNotFound, w.Code)
// 	s.NoError(mock.ExpectationsWereMet())
// }

// func (s *SigninTestSuite) TestSignin_InternalServerError() {
// 	t := s.T()
// 	gdb, mock := signinSQLMock(t)
// 	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT ?")).WithArgs("Sanjana", 1).WillReturnError(sql.ErrNoRows)
// 	r := GetSigninRouter(gdb)

// 	req := httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBufferString(`{"username":"Sanjana","password":"Sanjana@123"}`))
// 	req.Header.Set("Content-type", "application/json")
// 	w := httptest.NewRecorder()

// 	r.ServeHTTP(w, req)
// 	s.Equal(http.StatusInternalServerError, w.Code)
// 	s.NoError(mock.ExpectationsWereMet())
// }

// func (s *SigninTestSuite) TestSignin_ValidationError() {
// 	t := s.T()
// 	gdb, mock := signinSQLMock(t)
// 	r := GetSigninRouter(gdb)

// 	req := httptest.NewRequest(http.MethodPost, "/api/auth/signin", bytes.NewBufferString(`{"password":"Sanjana@123"}`))
// 	req.Header.Set("Content-type", "application/json")
// 	w := httptest.NewRecorder()

// 	r.ServeHTTP(w, req)
// 	s.Equal(http.StatusBadRequest, w.Code)
// 	s.NoError(mock.ExpectationsWereMet())
// }

// func TestSigininTestSuite(t *testing.T) {
// 	suite.Run(t, new(SigninTestSuite))
// }
