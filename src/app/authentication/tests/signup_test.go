package tests

// import (
// 	"authentication/business"
// 	"authentication/handlers"
// 	"authentication/repository"
// 	"bytes"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"regexp"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/suite"
// 	"gorm.io/gorm"
// )

// var createUserQuery = "INSERT INTO `users` (`username`,`password`,`panCard`,`phoneNumber`,`email`) VALUES (?,?,?,?,?)"

// func getSignupRouter(gDB *gorm.DB) *gin.Engine {
// 	gin.SetMode(gin.TestMode)
// 	router := gin.New()

// 	repo := repository.NewCreateUserRepository(gDB)
// 	service := business.NewCreateUserService(repo)
// 	handler := handlers.NewCreateUserHandler(service)

// 	router.POST("/api/auth/signup", handler.HandleCreaterUser)

// 	return router
// }

// type SignupUserTestSuite struct {
// 	suite.Suite
// }

// func (suite *SignupUserTestSuite) SetupTest() {
// 	gin.SetMode(gin.TestMode)
// }

// func TestMockSignupUser(t *testing.T) {
// 	suite.Run(t, new(SignupUserTestSuite))
// }

// func (suite *SignupUserTestSuite) TestMockSignupUser201SignupSuccessfull() {
// 	t := suite.T()

// 	gDB, mock := getMockSQLDB(t)

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(createUserQuery)).
// 		WithArgs(
// 			"Arijit",
// 			sqlmock.AnyArg(),
// 			"EQZRP1234P",
// 			7568912340,
// 			"arijit@gmail.com",
// 		).
// 		WillReturnResult(sqlmock.NewResult(1, 1))
// 	mock.ExpectCommit()

// 	router := getSignupRouter(gDB)

// 	body := `{
//   "confirmPassword": "Secure@123",
//   "email": "arijit@gmail.com",
//   "panCard": "EQZRP1234P",
//   "password": "Secure@123",
//   "phoneNumber": 7568912340,
//   "username": "Arijit"
// }`

// 	request := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
// 	request.Header.Set("Content-type", "application/json")

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, request)

// 	suite.Equal(http.StatusCreated, w.Code)
// 	suite.Contains(w.Body.String(), "User created successfully")
// 	suite.NoError(mock.ExpectationsWereMet())
// }

// func (suite *SignupUserTestSuite) TestMockSignupUser409DuplicateEmail() {
// 	t := suite.T()

// 	gDB, mock := getMockSQLDB(t)

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(createUserQuery)).
// 		WithArgs(
// 			"Arijit",
// 			sqlmock.AnyArg(),
// 			"EQZRP1234P",
// 			7568912340,
// 			"arijit@gmail.com",
// 		).
// 		WillReturnError(errors.New("duplicate key value violates unique constraint : idx_users_email"))
// 	mock.ExpectRollback()

// 	router := getSignupRouter(gDB)

// 	body := `{
//   "confirmPassword": "Secure@123",
//   "email": "arijit@gmail.com",
//   "panCard": "EQZRP1234P",
//   "password": "Secure@123",
//   "phoneNumber": 7568912340,
//   "username": "Arijit"
// }`

// 	request := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
// 	request.Header.Set("Content-type", "application/json")

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, request)

// 	suite.Equal(http.StatusConflict, w.Code)
// 	suite.Contains(w.Body.String(), "email")
// 	suite.NoError(mock.ExpectationsWereMet())
// }

// func (suite *SignupUserTestSuite) TestMockSignupUser409DuplicatePanCard() {
// 	t := suite.T()

// 	gDB, mock := getMockSQLDB(t)

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(createUserQuery)).
// 		WithArgs(
// 			"Arijit",
// 			sqlmock.AnyArg(),
// 			"EQZRP1234P",
// 			7568912340,
// 			"arijit@gmail.com",
// 		).
// 		WillReturnError(errors.New("duplicate key value violates unique constraint :  idx_users_pan_card"))
// 	mock.ExpectRollback()

// 	router := getSignupRouter(gDB)

// 	body := `{
//   "confirmPassword": "Secure@123",
//   "email": "arijit@gmail.com",
//   "panCard": "EQZRP1234P",
//   "password": "Secure@123",
//   "phoneNumber": 7568912340,
//   "username": "Arijit"
// }`

// 	request := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
// 	request.Header.Set("Content-type", "application/json")

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, request)

// 	suite.Equal(http.StatusConflict, w.Code)
// 	suite.Contains(w.Body.String(), "panCard")
// 	suite.NoError(mock.ExpectationsWereMet())
// }

// func (suite *SignupUserTestSuite) TestMockSignupUser409DuplicateUsername() {
// 	t := suite.T()

// 	gDB, mock := getMockSQLDB(t)

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(createUserQuery)).
// 		WithArgs(
// 			"Arijit",
// 			sqlmock.AnyArg(),
// 			"EQZRP1234P",
// 			7568912340,
// 			"arijit@gmail.com",
// 		).
// 		WillReturnError(errors.New("duplicate key value violates unique constraint"))
// 	mock.ExpectRollback()

// 	router := getSignupRouter(gDB)

// 	body := `{
//   "confirmPassword": "Secure@123",
//   "email": "arijit@gmail.com",
//   "panCard": "EQZRP1234P",
//   "password": "Secure@123",
//   "phoneNumber": 7568912340,
//   "username": "Arijit"
// }`

// 	request := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
// 	request.Header.Set("Content-type", "application/json")

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, request)

// 	suite.Equal(http.StatusConflict, w.Code)
// 	suite.Contains(w.Body.String(), "username")
// 	suite.NoError(mock.ExpectationsWereMet())
// }

// func (suite *SignupUserTestSuite) TestMockSignupUser500InternalServerError() {
// 	t := suite.T()

// 	gDB, mock := getMockSQLDB(t)

// 	mock.ExpectBegin()
// 	mock.ExpectExec(regexp.QuoteMeta(createUserQuery)).
// 		WithArgs(
// 			"Arijit",
// 			sqlmock.AnyArg(),
// 			"EQZRP1234P",
// 			7568912340,
// 			"arijit@gmail.com",
// 		).
// 		WillReturnError(errors.New("internal server error"))
// 	mock.ExpectRollback()

// 	router := getSignupRouter(gDB)

// 	body := `{
//   "confirmPassword": "Secure@123",
//   "email": "arijit@gmail.com",
//   "panCard": "EQZRP1234P",
//   "password": "Secure@123",
//   "phoneNumber": 7568912340,
//   "username": "Arijit"
// }`

// 	request := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
// 	request.Header.Set("Content-type", "application/json")

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, request)

// 	suite.Equal(http.StatusInternalServerError, w.Code)
// 	suite.Contains(w.Body.String(), "failed to create user")
// 	suite.NoError(mock.ExpectationsWereMet())
// }

// func (suite *SignupUserTestSuite) TestMockSignupUser400InvalidPayload() {
// 	t := suite.T()

// 	gDB, mock := getMockSQLDB(t)

// 	router := getSignupRouter(gDB)

// 	body := `{
//   "confirmPassword": "Secure@123",
//   "email": "arijit@gmail.com",
//   "panCard": "EQZRP1234P",
//   "password": "Secure@123",
//   "phoneNumber": 7568912340,
//   "username": 123
// }`

// 	request := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
// 	request.Header.Set("Content-type", "application/json")

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, request)

// 	suite.Equal(http.StatusBadRequest, w.Code)
// 	suite.Contains(w.Body.String(), "invalid required payload")
// 	suite.NoError(mock.ExpectationsWereMet())
// }

// func (suite *SignupUserTestSuite) TestMockSignupUser400ValidationError() {
// 	t := suite.T()

// 	gDB, mock := getMockSQLDB(t)

// 	router := getSignupRouter(gDB)

// 	body := `{
//   "confirmPassword": "Secure@123",
//   "email": "arijit@gmail.com",
//   "panCard": "EQZRP1234P",
//   "password": "Secure@12",
//   "phoneNumber": 7568912340,
//   "username": "Arijit"
// }`

// 	request := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
// 	request.Header.Set("Content-type", "application/json")

// 	w := httptest.NewRecorder()
// 	router.ServeHTTP(w, request)

// 	suite.Equal(http.StatusBadRequest, w.Code)
// 	suite.Contains(w.Body.String(), "ConfirmPassword must match Password.")
// 	suite.NoError(mock.ExpectationsWereMet())
// }