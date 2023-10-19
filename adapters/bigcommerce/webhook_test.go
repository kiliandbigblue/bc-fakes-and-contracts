package bigcommerce

import (
	"testing"

	"github.com/kiliandbigblue/bc-fakes-and-contracts/domain/apiclient"
)

func TestBigcommerceWebhook(t *testing.T) {
	apiclient.ClientContract{NewClient: func() apiclient.Client {
		return NewClient(WithStore("mmvgd5qy8s", "ladbzc57v7b2l76lgpvhur2bnilu2ne"))
	}}.Test(t)
}
