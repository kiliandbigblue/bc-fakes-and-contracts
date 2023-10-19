package inmemory

import (
	"testing"

	"github.com/kiliandbigblue/bc-fakes-and-contracts/domain/apiclient"
)

func TestInMemoryClient(t *testing.T) {
	apiclient.ClientContract{NewClient: func() apiclient.Client {
		return NewClient()
	}}.Test(t)
}
