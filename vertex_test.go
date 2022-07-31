package main

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"vertex/tests"
)

func TestVertexTestSuite(t *testing.T) {
	suite.Run(t, new(tests.VertexTestSuite))
}
