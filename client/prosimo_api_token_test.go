package client

import (
	"context"
	"fmt"
	"testing"
)

func TestCreateApiToken(t *testing.T) {

	ctx := context.Background()

	prosimo_client, err := NewProsimoClient("https://myevekuldeep1610576711105.dashboard.psonar.us/", "aggS7fie4O2mulcuEh2XbGH0cjJ4R1WHIhXggsnxr-I=", true)
	if err != nil {
		fmt.Printf("err - %s", err)
	}

	roles := []*Role{}
	role := Role{
		Name: "admin",
	}
	roles = append(roles, &role)
	prosimoApiTokenOpts := &ProsimoApiToken{
		Name: "test-2",
		Role: roles,
	}
	respString, err := prosimo_client.CreateAPIToken(ctx, prosimoApiTokenOpts)

	fmt.Printf("respString2 - %s", respString)

	if err != nil {
		fmt.Printf("err - %s", err)
	}

}
