package client

import (
	"context"
	"fmt"
)

type Dynamic_Risk struct {
	Id          string      `json:"id,omitempty"`
	Thresholds  []Threshold `json:"thresholds,omitempty"`
	CreatedTime string      `json:"createdTime,omitempty"`
	UpdatedTime string      `json:"updatedTime,omitempty"`
}

type Threshold struct {
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
	Value   int    `json:"value,omitempty"`
}

type Dynamic_RiskList struct {
	DyRisk []*Dynamic_Risk `json:"data,omitempty"`
}

func (threshold Threshold) String() string {
	return fmt.Sprintf("{Name:%s, Enabled:%t, Value:%d}",
		threshold.Name, threshold.Enabled, threshold.Value)
}

func (dynamicRisk Dynamic_Risk) String() string {
	return fmt.Sprintf("{Id:%s, CreatedTime:%s, UpdatedTime:%s}",
		dynamicRisk.Id, dynamicRisk.CreatedTime, dynamicRisk.UpdatedTime)
}

func (prosimoClient *ProsimoClient) GetDYRisk(ctx context.Context) (*Dynamic_RiskList, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", DynamicRiskEndpoint, nil)
	if err != nil {
		return nil, err
	}
	dyriskList := &Dynamic_RiskList{}
	_, err = prosimoClient.api_client.Do(ctx, req, dyriskList)

	if err != nil {
		return nil, err
	}

	return dyriskList, nil

}

func (prosimoClient *ProsimoClient) PutDYRisk(ctx context.Context, dyriskops *Dynamic_Risk) (*Dynamic_RiskList, error) {

	updateDyRiskEndpt := fmt.Sprintf("%s/%s", DynamicRiskEndpoint, dyriskops.Id)
	req, err := prosimoClient.api_client.NewRequest("PUT", updateDyRiskEndpt, dyriskops)
	if err != nil {
		return nil, err
	}

	dyriskputData := &Dynamic_RiskList{}
	_, err = prosimoClient.api_client.Do(ctx, req, dyriskputData)
	if err != nil {
		return nil, err
	}

	return dyriskputData, nil

}
