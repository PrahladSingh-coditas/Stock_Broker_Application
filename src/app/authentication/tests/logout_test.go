package tests

import (
	"authentication/business"
	"authentication/handlers"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"stock_broker_application/src/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

func GetLogoutRouter(redisClient *redis.Client) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	logoutservice := business.NewLogoutUser(redisClient)
	logouthandler := handlers.LogoutUserHandler(logoutservice)
	r.POST("/api/auth/logout", func(ctx *gin.Context) {
		ctx.Set("token", "this-is-token-string")
		ctx.Set("exp", 1798674348)
		ctx.Next()
	}, logouthandler.HandleUserLogout)
	return r
}

type LogoutTestSuite struct {
	suite.Suite
}

func (s *LogoutTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (s *LogoutTestSuite) TestLogout200LogoutSuccess() {
	mockRedisClient, mockRedisController, _ := utils.GetRedisClient(true)
	mockRedisController.ExpectSet("BLACKLISTE_TOKEN:this-is-token-string", 1, 0).SetVal("Ok")
	r := GetLogoutRouter(mockRedisClient)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Logout Successful and Token Blacklisted")
	s.NoError(mockRedisController.ExpectationsWereMet())
}

func (s *LogoutTestSuite) TestLogout500RedisOperationFailed() {
	mockRedisClient, mockRedisController, _ := utils.GetRedisClient(true)
	mockRedisController.ExpectSet("BLACKLISTE_TOKEN:this-is-token-string", 1, 0).SetErr(errors.New("Error in Performing Operation"))
	r := GetLogoutRouter(mockRedisClient)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.Contains(w.Body.String(), "authentication failed")
	s.NoError(mockRedisController.ExpectationsWereMet())
}

func (s *LogoutTestSuite) TestLogout500RedisConnectionFailed() {
	mockRedisClient, mockRedisController, _ := utils.GetRedisClient(true)
	mockRedisController.ExpectSet("BLACKLISTE_TOKEN:this-is-token-string", 1, 0).SetErr(errors.New("Redis Connection Interrupted, Failed Operation"))
	r := GetLogoutRouter(mockRedisClient)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Equal(http.StatusInternalServerError, w.Code)
	s.Contains(w.Body.String(), "authentication failed")
	s.NoError(mockRedisController.ExpectationsWereMet())
}

func TestLogoutTestSuite(t *testing.T) {
	suite.Run(t, new(LogoutTestSuite))
}
