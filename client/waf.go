package client

import (
	"context"
	"errors"
	"fmt"
)

type WafListData struct {
	Data []*Waf `json:"data,omitempty"`
}

type WafData struct {
	Data *Waf `json:"data,omitempty"`
}

type Waf struct {
	DefaultWaf bool            `json:"default,omitempty"`
	ID         string          `json:"id,omitempty"`
	Name       string          `json:"name,omitempty"`
	Mode       string          `json:"mode,omitempty"`
	Threshold  int             `json:"threshold,omitempty"`
	TeamID     string          `json:"teamID,omitempty"`
	WafRuleSet *WafRuleSet     `json:"rulesets,omitempty"`
	AppDomains []*WafAppDomain `json:"appDomains,omitempty"`
}

type WafRuleSet struct {
	Basic *WafRuleGroups `json:"basic,omitempty"`
	OWASP *WafRuleGroups `json:"owasp-crs-v32,omitempty"`
}

type WafRuleGroups struct {
	Name       string   `json:"name,omitempty"`
	Rulegroups []string `json:"rulegroups,omitempty"`
}

type WafAppDomain struct {
	AppID  string `json:"appID,omitempty"`
	ID     string `json:"id,omitempty"`
	Domain string `json:"domain,omitempty"`
}

func (waf Waf) String() string {
	return fmt.Sprintf("{Name:%s, Mode:%s, Threshold:%d, WafRuleSet:%v, ID:%s}",
		waf.Name, waf.Mode, waf.Threshold, waf.WafRuleSet, waf.ID)
}

func (wafRuleSet WafRuleSet) String() string {
	return fmt.Sprintf("{Basic:%v, OWASP:%v}",
		wafRuleSet.Basic, wafRuleSet.OWASP)
}

func (wafRuleGroups WafRuleGroups) String() string {
	return fmt.Sprintf("{Name:%s, Rulegroups:%s}",
		wafRuleGroups.Name, wafRuleGroups.Rulegroups)
}

func (prosimoClient *ProsimoClient) GetWaf(ctx context.Context) ([]*Waf, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", WafEndpoint, nil)
	if err != nil {
		return nil, err
	}

	wafListData := &WafListData{}
	_, err = prosimoClient.api_client.Do(ctx, req, wafListData)
	if err != nil {
		return nil, err
	}

	return wafListData.Data, nil

}

func (prosimoClient *ProsimoClient) GetWafByName(ctx context.Context, wafName string) (*Waf, error) {

	wafList, err := prosimoClient.GetWaf(ctx)
	if err != nil {
		return nil, err
	}

	var waf *Waf
	for _, returnedWaf := range wafList {
		if returnedWaf.Name == wafName {
			waf = returnedWaf
			break
		}
	}

	if waf == nil {
		return nil, errors.New("Waf doesnt exists")
	}

	return waf, nil

}

func (prosimoClient *ProsimoClient) GetWafByID(ctx context.Context, wafID string) (*Waf, error) {

	wafList, err := prosimoClient.GetWaf(ctx)
	if err != nil {
		return nil, err
	}

	var waf *Waf
	for _, returnedWaf := range wafList {
		if returnedWaf.ID == wafID {
			waf = returnedWaf
			break
		}
	}

	if waf == nil {
		return nil, errors.New("Waf doesnt exists")
	}

	return waf, nil

}

func (prosimoClient *ProsimoClient) CreateWaf(ctx context.Context, waf *Waf) (*Waf, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", WafEndpoint, waf)
	if err != nil {
		return nil, err
	}

	wafData := &WafData{}
	_, err = prosimoClient.api_client.Do(ctx, req, wafData)
	if err != nil {
		return nil, err
	}

	return wafData.Data, nil

}

func (prosimoClient *ProsimoClient) UpdateWaf(ctx context.Context, waf *Waf) (*Waf, error) {

	updateWafEndpt := fmt.Sprintf("%s/%s", WafEndpoint, waf.ID)

	req, err := prosimoClient.api_client.NewRequest("PUT", updateWafEndpt, waf)
	if err != nil {
		return nil, err
	}

	wafData := &WafData{}
	_, err = prosimoClient.api_client.Do(ctx, req, wafData)
	if err != nil {
		return nil, err
	}

	return wafData.Data, nil

}

func (prosimoClient *ProsimoClient) DeleteWaf(ctx context.Context, wafID string) error {

	deleteWafEndpt := fmt.Sprintf("%s/%s", WafEndpoint, wafID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", deleteWafEndpt, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}

func (prosimoClient *ProsimoClient) GetWafRuleSet(ctx context.Context) (*WafRuleSet, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", WafRuleSetEndpoint, nil)
	if err != nil {
		return nil, err
	}

	wafData := &WafData{}
	_, err = prosimoClient.api_client.Do(ctx, req, wafData)
	if err != nil {
		return nil, err
	}

	return wafData.Data.WafRuleSet, nil

}

type WafAppDomainIds struct {
	AddDomainIDs    []string `json:"addDomainIDs,omitempty"`
	DeleteDomainIDs []string `json:"deleteDomainIDs,omitempty"`
}

func (wafAppDomainIds WafAppDomainIds) String() string {
	return fmt.Sprintf("{AddDomainIDs:%s, DeleteDomainIDs:%s}",
		wafAppDomainIds.AddDomainIDs, wafAppDomainIds.DeleteDomainIDs)
}

func (prosimoClient *ProsimoClient) UpdateWafAppDomains(ctx context.Context, wafAppDomainIds *WafAppDomainIds, wafID string) (*Waf, error) {

	updateWafAppDomainEndpt := fmt.Sprintf(WafAppDomainEndpoint, wafID)

	req, err := prosimoClient.api_client.NewRequest("PUT", updateWafAppDomainEndpt, wafAppDomainIds)
	if err != nil {
		return nil, err
	}

	wafData := &WafData{}
	_, err = prosimoClient.api_client.Do(ctx, req, wafData)
	if err != nil {
		return nil, err
	}

	return wafData.Data, nil

}
