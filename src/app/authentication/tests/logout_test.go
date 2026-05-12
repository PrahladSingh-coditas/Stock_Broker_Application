package tests

import (
	"authentication/business"
	"authentication/handlers"
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"authentication/commons/constants"
	"stock_broker_application/src/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/stretchr/testify/suite"
)

func getLogoutRouter(rclient *redis.Client) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	service := business.NewLogoutUserService(rclient)
	handler := handlers.NewLogoutUserHandler(service)

	router.POST("/api/auth/logout", func(ctx *gin.Context) {
		ctx.Set("username", "Arijit")
		ctx.Set("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Nzg1NzA2NjgsImlhdCI6MTc3ODQ4NDI2OCwic3ViIjoiQXJpaml0In0.oQV8f1D0SlRg9Gu2gFLqMvmlhzhAIFyncLv6PWpHliw")
		ctx.Set("expiry-time", 1778570668)
		ctx.Next()
	}, handler.HandleLogoutUser)
	return router
}

type LogoutTestSuite struct {
	suite.Suite
}

func (user *LogoutTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func TestMockLogout(t *testing.T) {
	suite.Run(t, new(LogoutTestSuite))
}

func (suite *LogoutTestSuite) TestMockLogoutUser200LogoutSuccessful() {
	//t := suite.T()
	ctx := context.Background()
	rClient, mock, _ := utils.GetRedisClient(ctx, false)
	token := "BLACKLISTED_TOKEN_eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Nzg1NzA2NjgsImlhdCI6MTc3ODQ4NDI2OCwic3ViIjoiQXJpaml0In0.oQV8f1D0SlRg9Gu2gFLqMvmlhzhAIFyncLv6PWpHliw"

	//mock.ExpectExists(token).SetVal(0)
	mock.ExpectSet(token, 1, 0).SetVal("OK")

	router := getLogoutRouter(rClient)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/logout", bytes.NewBufferString(`{}`))
	request.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Nzg1NzA2NjgsImlhdCI6MTc3ODQ4NDI2OCwic3ViIjoiQXJpaml0In0.oQV8f1D0SlRg9Gu2gFLqMvmlhzhAIFyncLv6PWpHliw")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "logout successfully")
	suite.NoError(mock.ExpectationsWereMet())
}

func (suite *LogoutTestSuite) TestMockLogoutUser500RedisSetError() {
	//t := suite.T()
	ctx := context.Background()
	rClient, mock, _ := utils.GetRedisClient(ctx, false)
	token := "BLACKLISTED_TOKEN_eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Nzg1NzA2NjgsImlhdCI6MTc3ODQ4NDI2OCwic3ViIjoiQXJpaml0In0.oQV8f1D0SlRg9Gu2gFLqMvmlhzhAIFyncLv6PWpHliw"

	mock.ExpectSet(token, 1, 0).SetErr(errors.New(constants.RedisSetOperationError))

	router := getLogoutRouter(rClient)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/logout", bytes.NewBufferString(`{}`))
	//request.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Nzg1NzA2NjgsImlhdCI6MTc3ODQ4NDI2OCwic3ViIjoiQXJpaml0In0.oQV8f1D0SlRg9Gu2gFLqMvmlhzhAIFyncLv6PWpHliw")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusInternalServerError, w.Code)
	suite.Contains(w.Body.String(), "error in redis set operation")
	suite.NoError(mock.ExpectationsWereMet())
}
