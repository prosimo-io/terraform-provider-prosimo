package client

import (
	"context"
	"fmt"
	"testing"
)

func TestGetIPREP(t *testing.T) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient("https://myeveroot1619612166376.dashboard.psonar.us/", "ElwOWlHLSna5cqgKZ-0NrHPzmUqpGcUhSJ8Kn5U0LEw=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	getIPRep, err := prosimo_client.GetIPREP(ctx)

	fmt.Printf("respString2 - %v\n", getIPRep.IpReps)

	if err != nil {
		fmt.Printf("err - %s", err)
	}

}

func TestPutIPREP(t *testing.T) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient("https://myeveroot1619612166376.dashboard.psonar.us/", "ElwOWlHLSna5cqgKZ-0NrHPzmUqpGcUhSJ8Kn5U0LEw=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}
	// blocklist := []string{}
	ipRep := &IP_Reputation{
		Enabled:   true,
		AllowList: []string{"2.1.1.1/16"},
		// BlockList: []string{},
	}
	fmt.Println(ipRep)
	respString, err := prosimo_client.PutIPREP(ctx, ipRep)

	fmt.Printf("respString2 - %v", respString.IpReps)

	if err != nil {
		fmt.Printf("err - %s", err)
	}

}
