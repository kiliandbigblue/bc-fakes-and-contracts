package inmemory

import (
	"testing"

	"github.com/bigbluedisco/bc-fakes-and-contracts/domain/apiclient"
)

func TestInmemoryAPI1(t *testing.T) {
	apiclient.ClientContract{NewClient: func() apiclient.Client {
		return NewClient()
	}}
}
