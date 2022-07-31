package main

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"tiny-submarine/tests"
)

func TestSubmarineTestSuite(t *testing.T) {
	suite.Run(t, new(tests.SubmarineTestSuite))
}
