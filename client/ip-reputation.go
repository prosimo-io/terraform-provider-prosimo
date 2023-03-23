package client

import (
	"context"
	"fmt"
)

type IP_Reputation struct {
	Enabled   bool     `json:"enabled,omitempty"`
	AllowList []string `json:"allowlist,omitempty"`
	BlockList []string `json:"blocklist,omitempty"`
}

type IP_ReputationData struct {
	IpReps *IP_Reputation `json:"data,omitempty"`
}

func (ipr IP_Reputation) String() string {
	return fmt.Sprintf("{Enabled:%t, Allowlist:%s}",
		ipr.Enabled, ipr.AllowList)
}

func (prosimoClient *ProsimoClient) GetIPREP(ctx context.Context) (*IP_ReputationData, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", IPRepEndpoint, nil)
	if err != nil {
		return nil, err
	}
	iprepList := &IP_ReputationData{}
	_, err = prosimoClient.api_client.Do(ctx, req, iprepList)

	if err != nil {
		return nil, err
	}

	return iprepList, nil

}

func (prosimoClient *ProsimoClient) PutIPREP(ctx context.Context, iprops *IP_Reputation) (*IP_ReputationData, error) {
	req, err := prosimoClient.api_client.NewRequest("PUT", IPRepEndpoint, iprops)
	if err != nil {
		return nil, err
	}

	prepListData := &IP_ReputationData{}
	_, err = prosimoClient.api_client.Do(ctx, req, prepListData)
	if err != nil {
		return nil, err
	}

	return prepListData, nil

}
