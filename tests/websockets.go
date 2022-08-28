package tests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"vertex/communication"
)

func (s *VertexTestSuite) TestWebsocketDisconnection() {
	testServer := httptest.NewServer(s.Gin)
	defer testServer.Close()
	conn, _, err := s.Dialer.Dial(
		fmt.Sprintf("%s/ws", strings.Replace(testServer.URL, "http", "ws", 1)),
		http.Header{
			"Authorization": {"Bearer " + s.SignedToken},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	err = conn.WriteJSON(communication.UserRequest{
		Service: "vertex",
		Method:  "disconnect",
		Data:    gin.H{},
	})
	if err != nil {
		log.Fatal(err)
	}
	err = conn.ReadJSON(nil)
	// If connection is closed, err on read will be occurred
	assert.ErrorContains(s.T(), err, "close")
}
