package tests

import (
	"github.com/lestrrat-go/jwx/jwk"
	"os"
	"strings"
	"vertex/config"
)

func (s *VertexKeysTestSuite) TestLocalKeyLoad() {
	if os.Getenv("TEST_KEY_FILE") == "" && strings.HasPrefix(config.Config.Strings("keys")[0], "http") {
		s.T().Skip("Local key file path is not specified, please, set it up via config or use TEST_KEY_FILE env")
	} else {
		var keyPath string
		if os.Getenv("TEST_KEY_FILE") != "" {
			keyPath = os.Getenv("TEST_KEY_FILE")
		} else {
			keyPath = config.Config.Strings("keys")[0]
		}
		var key *jwk.Key
		key, err := config.LoadKey(keyPath)
		s.Assert().Nil(err)
		s.Assert().NotEmpty(key)

	}
}

func (s *VertexKeysTestSuite) TestKeyLoadRemote() {
	if os.Getenv("TEST_KEY_REMOTE_PATH") == "" && !strings.HasPrefix(config.Config.Strings("keys")[0], "http") {
		s.T().Skip("Remote key path is not specified, please, set it up via config or use TEST_KEY_REMOTE_PATH env")
	} else {
		var keyPath string
		if os.Getenv("TEST_KEY_REMOTE_PATH") != "" {
			keyPath = os.Getenv("TEST_KEY_REMOTE_PATH")
		} else {
			keyPath = config.Config.Strings("keys")[0]
		}
		var key *jwk.Key
		key, err := config.LoadKey(keyPath)
		s.Assert().Nil(err)
		s.Assert().NotEmpty(key)
	}
}

func (s *VertexKeysTestSuite) TestKeyLoadRemoteIncorrect() {
	keyPath := "https://someurlthatdoesnt.exists/public.jwk"
	key, err := config.LoadKey(keyPath)
	s.Assert().Nil(key)
	s.Assert().Error(err)
}
