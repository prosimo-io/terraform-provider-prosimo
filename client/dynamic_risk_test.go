package client

import (
	"context"
	"fmt"
	"testing"
)

func TestGetDYRISK(t *testing.T) {

	ctx := context.Background()
	//var id string

	prosimo_client, err := NewProsimoClient("https://myeveroot1619612166376.dashboard.psonar.us/", "ElwOWlHLSna5cqgKZ-0NrHPzmUqpGcUhSJ8Kn5U0LEw=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	getDYRisk, err := prosimo_client.GetDYRisk(ctx)

	if err != nil {
		fmt.Printf("err - %s", err)
	}
	fmt.Printf("respString - %v\n", getDYRisk.DyRisk)
	for i, val := range getDYRisk.DyRisk {
		fmt.Printf("i, v %v, %v\n", i, val)
		fmt.Println(val.Id)
		//id = val.Id

	}

}

func TestPutDYRISK(t *testing.T) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient("https://myeveroot1619612166376.dashboard.psonar.us/", "ElwOWlHLSna5cqgKZ-0NrHPzmUqpGcUhSJ8Kn5U0LEw=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	thresholds := []Threshold{}
	threshold1 := Threshold{
		Name:    "alert",
		Enabled: true,
		Value:   45,
	}
	thresholds = append(thresholds, threshold1)
	threshold2 := Threshold{
		Name:    "mfa",
		Enabled: true,
		Value:   45,
	}
	thresholds = append(thresholds, threshold2)
	threshold3 := Threshold{
		Name:    "lockUser",
		Enabled: false,
		Value:   100,
	}
	_ = threshold3
	//thresholds = append(thresholds, threshold3)

	fmt.Println(thresholds)

	threshold_details := &Dynamic_Risk{
		Id:         "23ad3b7a-6837-4f39-82e8-10b9d92d0ac9",
		Thresholds: thresholds,
	}

	respString, err := prosimo_client.PutDYRisk(ctx, threshold_details)
	_ = respString

	fmt.Printf("respString2 - %v", respString.DyRisk)

	if err != nil {
		fmt.Printf("err - %s", err)
	}

}
