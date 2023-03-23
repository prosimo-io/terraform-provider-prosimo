package client

import (
	"context"
	"fmt"
	"testing"
)

func TestUpdateAllowList(t *testing.T) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient("https://myeveroot1619612166376.dashboard.psonar.us/", "ElwOWlHLSna5cqgKZ-0NrHPzmUqpGcUhSJ8Kn5U0LEw=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}
	deleteusers := []Users{}
	addusers := []Users{}
	adduser := Users{
		Email:  "abc@demo.com",
		Reason: "def",
	}
	userStatus, err := prosimo_client.CheckUserDetails(ctx, adduser.Email)
	fmt.Println(userStatus)
	if userStatus.UserAvailabilityStatus.Present == false {
		addusers = append(addusers, adduser)
	}
	postallowlist := &PostAllowlist{
		AddUsers:    addusers,
		DeleteUsers: deleteusers,
	}
	respString, err := prosimo_client.UpdateAllowList(ctx, postallowlist)
	_ = respString
	//fmt.Printf("respString2 - %v", respString.postallowlist)

	if err != nil {
		fmt.Printf("err - %s", err)
	}

}

func TestSearchAllowList(t *testing.T) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient("https://myeveroot1619612166376.dashboard.psonar.us/", "ElwOWlHLSna5cqgKZ-0NrHPzmUqpGcUhSJ8Kn5U0LEw=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	//postallowlist := &GetAllowlist{}
	searchallowlist, err := prosimo_client.SearchAllowList(ctx)

	if err != nil {
		fmt.Printf("err - %s", err)
	}

	fmt.Printf("respString2 %v\n", searchallowlist.GetAllowlist)

}

func TestDeleteAllowList(t *testing.T) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient("https://myeveroot1619612166376.dashboard.psonar.us/", "ElwOWlHLSna5cqgKZ-0NrHPzmUqpGcUhSJ8Kn5U0LEw=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	//postallowlist := &GetAllowlist{}
	res := prosimo_client.DeleteAllUsers(ctx)

	fmt.Printf("respString2 %v\n", res)

}
