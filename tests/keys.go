package tests

import (
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"os"
	"vertex/config"
)

func (s *VertexKeysTestSuite) TestLocalKeyLoad() {
	if config.Config.Bool("keys.is_remote") && os.Getenv("TEST_KEY_FILE") == "" {
		s.T().Skip("Local key file path is not specified, please, set it up via config or use TEST_KEY_FILE env")
	} else {
		var keyPath string
		if os.Getenv("TEST_KEY_FILE") != "" {
			keyPath = os.Getenv("TEST_KEY_FILE")
		} else {
			keyPath = config.Config.String("keys.public")
		}
		var key *jwk.Key
		key, err := config.LoadKey(keyPath, false, jwa.RS256)
		s.Assert().Nil(err)
		s.Assert().NotEmpty(key)

	}
}

func (s *VertexKeysTestSuite) TestKeyLoadRemote() {
	if !config.Config.Bool("is_remote") && os.Getenv("TEST_KEY_REMOTE_PATH") == "" {
		s.T().Skip("Remote key path is not specified, please, set it up via config or use TEST_KEY_REMOTE_PATH env")
	} else {
		var keyPath string
		if os.Getenv("TEST_KEY_REMOTE_PATH") != "" {
			keyPath = os.Getenv("TEST_KEY_REMOTE_PATH")
		} else {
			keyPath = config.Config.String("keys.public")
		}
		var key *jwk.Key
		key, err := config.LoadKey(keyPath, true, jwa.RS256)
		s.Assert().Nil(err)
		s.Assert().NotEmpty(key)
	}
}

func (s *VertexKeysTestSuite) TestKeyLoadRemoteIncorrect() {
	keyPath := "https://someurlthatdoesnt.exists/public.pem"
	key, err := config.LoadKey(keyPath, true, jwa.RS256)
	s.Assert().Nil(key)
	s.Assert().Error(err)
}
