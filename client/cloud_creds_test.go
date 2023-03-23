package client

import (
	"context"
	"fmt"
	"testing"
)

func TestCreateCloudCreds(t *testing.T) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient("https://myeveroot1619612166376.dashboard.psonar.us/", "ElwOWlHLSna5cqgKZ-0NrHPzmUqpGcUhSJ8Kn5U0LEw=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	filepath := "/Users/sibaprasadtripathy/Downloads/prosimo-test-bf8bbf15b37c_copy.json"

	cloudCreds := &CloudCreds{
		CloudType: GCPCloudType,
		Nickname:  "demo2",
		KeyType:   GCPKeyType,
	}

	fmt.Println(cloudCreds)
	cloudCredsData, err := prosimo_client.UploadGcpCloudCreds(ctx, cloudCreds, filepath)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("cloudCreds created id - %s", cloudCredsData.CloudCreds.ID)

}
