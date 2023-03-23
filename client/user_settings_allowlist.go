package client

import (
	"bytes"
	"context"
	"fmt"
)

type GetAllowlistData struct {
	GetAllowlist *GetAllowlist `json:"data,omitempty"`
}

type PostAllowlistData struct {
	PostAllowlist *PostAllowlist `json:"data,omitempty"`
}

type UserAvailabilityStatusData struct {
	UserAvailabilityStatus *UserAvailabilityStatus `json:"data,omitempty"`
}

type PostAllowlist struct {
	AddUsers    []Users `json:"addUsers,omitempty"`
	DeleteUsers []Users `json:"deleteUsers,omitempty"`
}

type Users struct {
	Email       string `json:"email,omitempty"`
	Reason      string `json:"reason,omitempty"`
	CreatedTime string `json:"createdTime,omitempty"`
}

type GetAllowlist struct {
	TotalCount int     `json:"totalCount,omitempty"`
	Records    []Users `json:"records,omitempty"`
}

type UserAvailabilityStatus struct {
	Present bool `json:"present,omitempty"`
}

func (gl GetAllowlist) String() string {
	return fmt.Sprintf("{TotalCount:%d}", gl.TotalCount)
}

func (users Users) String() string {
	return fmt.Sprintf("{Email:%s, reason:%s", users.Email, users.Reason)
}

func (prosimoClient *ProsimoClient) UpdateAllowList(ctx context.Context, allowlistOps *PostAllowlist) (*PostAllowlistData, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", UserAllowlistEndpoint, allowlistOps)
	if err != nil {
		return nil, err
	}

	gv := &PostAllowlistData{}
	_, err = prosimoClient.api_client.Do(ctx, req, gv)
	if err != nil {
		return nil, err
	}

	return gv, nil

}

func (prosimoClient *ProsimoClient) SearchAllowList(ctx context.Context) (*GetAllowlistData, error) {
	updateUserAllowlistEndpoint := fmt.Sprintf("%s/%s", UserAllowlistEndpoint, "search")

	var jsonStr = []byte(`{"page":{"size":25,"start":0}}`)
	req, err := prosimoClient.api_client.NewRequest("POST", updateUserAllowlistEndpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	getAllowlistData := &GetAllowlistData{}
	_, err = prosimoClient.api_client.Do(ctx, req, getAllowlistData)
	if err != nil {
		return nil, err
	}

	return getAllowlistData, nil

}

func (prosimoClient *ProsimoClient) DeleteAllUsers(ctx context.Context) error {

	req, err := prosimoClient.api_client.NewRequest("DELETE", UserAllowlistEndpoint, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}

func (prosimoClient *ProsimoClient) CheckUserDetails(ctx context.Context, email string) (*UserAvailabilityStatusData, error) {

	updateUserAllowlistEndpoint := fmt.Sprintf("%s/%s", UserAllowlistEndpoint, "check")
	req, err := prosimoClient.api_client.NewRequest("POST", updateUserAllowlistEndpoint, email)
	if err != nil {
		return nil, err
	}

	getUserStatus := &UserAvailabilityStatusData{}
	_, err = prosimoClient.api_client.Do(ctx, req, getUserStatus)
	if err != nil {
		return nil, err
	}

	return getUserStatus, nil

}
