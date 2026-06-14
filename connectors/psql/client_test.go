package psql_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ClientSuite struct {
	PostgresRepositorySuite
}

func TestClientSuite(t *testing.T) {
	suite.Run(
		t,
		new(ClientSuite),
	)
}
