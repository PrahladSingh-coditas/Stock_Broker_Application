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
	"gorm.io/gorm"
)

var queryUpdatePassword = "UPDATE `users` SET `password`=? WHERE username = ?"

func getChangePasswordRouter(gdb *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	repo := repository.NewChangePasswordRepository(gdb)
	service := business.NewChangePasswordService(repo)
	handler := handlers.NewChangePasswordHandler(service)

	router.POST("/api/auth/change-password", func(ctx *gin.Context) {
		ctx.Set("username", "Arijit")
		ctx.Next()
	}, handler.HandleChangePassword)
	return router
}

type ChanegePasswordTestSuite struct {
	suite.Suite
}

func (suite *ChanegePasswordTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func TestMockChangePassword(t *testing.T) {
	suite.Run(t, new(ChanegePasswordTestSuite))
}

func (suite *ChanegePasswordTestSuite) TestMockChangePassword200PasswordChagedSuccssfully() {
	t := suite.T()
	gdb, mock := getMockSQLDB(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryUpdatePassword)).WithArgs(sqlmock.AnyArg(), "Arijit").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	router := getChangePasswordRouter(gdb)

	reqBody := `{
		"password":"Pass@123",
		"confirmPassword":"Pass@123"
	}`
	request := httptest.NewRequest(http.MethodPost, "/api/auth/change-password", bytes.NewBufferString(reqBody))
	request.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "Password updated successfully")
	suite.NoError(mock.ExpectationsWereMet())
}

func (suite *ChanegePasswordTestSuite) TestMockChangePassword500DatabaseQueryError() {
	t := suite.T()
	gdb, mock := getMockSQLDB(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryUpdatePassword)).WithArgs(sqlmock.AnyArg(), "Arijit").WillReturnError(errors.New(constants.ErrDatabaseQueryErrorMsg))
	mock.ExpectRollback()

	router := getChangePasswordRouter(gdb)

	reqBody := `{
		"password":"Pass@123",
		"confirmPassword":"Pass@123"
	}`

	request := httptest.NewRequest(http.MethodPost, "/api/auth/change-password", bytes.NewBufferString(reqBody))
	request.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), constants.ErrDatabaseQueryErrorMsg)
	suite.NoError(mock.ExpectationsWereMet())
}

func (suite *ChanegePasswordTestSuite) TestMockChangePassword404UserNotFoundError() {
	t := suite.T()
	gdb, mock := getMockSQLDB(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryUpdatePassword)).WithArgs(sqlmock.AnyArg(), "Arijit").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	router := getChangePasswordRouter(gdb)

	reqBody := `{
		"password":"Pass@123",
		"confirmPassword":"Pass@123"
	}`

	request := httptest.NewRequest(http.MethodPost, "/api/auth/change-password", bytes.NewBufferString(reqBody))
	request.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusNotFound, w.Code)
	suite.Contains(w.Body.String(), constants.ErrUserNotFoundMsg)
	suite.NoError(mock.ExpectationsWereMet())
}

func (suite *ChanegePasswordTestSuite) TestMockChangePassword400InvalidPayload() {
	t := suite.T()
	gdb, mock := getMockSQLDB(t)

	router := getChangePasswordRouter(gdb)

	reqBody := `{
		"password":1234,
		"confirmPassword":"Pass@123"
	}`

	request := httptest.NewRequest(http.MethodPost, "/api/auth/change-password", bytes.NewBufferString(reqBody))
	request.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), constants.ErrInvalidPayload)
	suite.NoError(mock.ExpectationsWereMet())
}

func (suite *ChanegePasswordTestSuite) TestMockChangePassword400ValidationFailedError() {
	t := suite.T()
	gdb, mock := getMockSQLDB(t)

	router := getChangePasswordRouter(gdb)

	reqBody := `{
		"password":"Pass@123",
		"confirmPassword":"Pass@12"
	}`

	request := httptest.NewRequest(http.MethodPost, "/api/auth/change-password", bytes.NewBufferString(reqBody))
	request.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusBadRequest, w.Code)
	suite.Contains(w.Body.String(), "ConfirmPassword must match Password.")
	suite.NoError(mock.ExpectationsWereMet())
}
