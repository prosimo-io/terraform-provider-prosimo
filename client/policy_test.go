package client

import (
	"fmt"
	"testing"
)

func TestCreatePolicy(t *testing.T) {

	// ctx := context.Background()

	prosimo_client, err := NewProsimoClient("https://myeveroot1641293042702.dashboard.psonar.us/", "XXIPZXcmMb3jXF91zhAFvVsyNUXP4PIfy6wl9_Ave2k=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	//fmt.Println(policy1)
	res := prosimo_client.ReadJson()
	fmt.Println(res)
}
