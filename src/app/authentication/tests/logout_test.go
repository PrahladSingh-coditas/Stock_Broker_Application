package tests

import (
	"authentication/business"
	"authentication/handlers"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

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
		ctx.Set("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzgwNDgwNjEsImlhdCI6MTc3Nzk2MTY2MSwic3ViIjoiQXJpaml0In0.NdfXzmr50xBhGFQm50qPwlvvkIKbVgeY88fMHbroX-c")
		ctx.Set("expiry-time", 1778048061)
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
	token := "BLACKLISTED_TOKEN_eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzgwNDgwNjEsImlhdCI6MTc3Nzk2MTY2MSwic3ViIjoiQXJpaml0In0.NdfXzmr50xBhGFQm50qPwlvvkIKbVgeY88fMHbroX-c"

	//mock.ExpectExists(token).SetVal(0)
	mock.ExpectSet(token, 1, 0).SetVal("OK")

	router := getLogoutRouter(rClient)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/logout", bytes.NewBufferString(`{}`))
	request.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzgwNDgwNjEsImlhdCI6MTc3Nzk2MTY2MSwic3ViIjoiQXJpaml0In0.NdfXzmr50xBhGFQm50qPwlvvkIKbVgeY88fMHbroX-c")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	suite.Equal(http.StatusOK, w.Code)
	suite.Contains(w.Body.String(), "logout successfully")
	suite.NoError(mock.ExpectationsWereMet())
}
