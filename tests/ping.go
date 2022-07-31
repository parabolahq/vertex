package tests

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
)

func (s *VertexTestSuite) TestPing() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	s.Gin.ServeHTTP(w, req)

	assert.Equal(s.T(), w.Code, 200)
	asJson := gin.H{}
	json.Unmarshal(w.Body.Bytes(), &asJson)
	assert.Equal(s.T(), asJson["ok"], true)
}
