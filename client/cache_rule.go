package client

import (
	"context"
	"errors"
	"fmt"
)

type CacheRuleListData struct {
	Data []*CacheRule `json:"data,omitempty"`
}

type CacheRuleData struct {
	Data *CacheRule `json:"data,omitempty"`
}


type CacheRule struct {
	BypassCache         bool        `json:"bypassCache,omitempty"`
	CacheControlIgnored bool        `json:"cacheControlIgnored,omitempty"`
	DefaultCache        bool        `json:"default,omitempty"`
	Editable            bool        `json:"editable,omitempty"`
	ID                  string      `json:"id,omitempty"`
	Name                string      `json:"name,omitempty"`
	ShareStaticContent  bool        `json:"shareStaticContent,omitempty"`
	TeamID              string      `json:"teamID,omitempty"`
	PathPatterns        []PATH      `json:"pathPatterns,omitempty"`
	AppDomains          []AppDomaiN `json:"appDomains,omitempty"`
	ByPassInfo          ByPassInfo  `json:"bypassInfo,omitempty"`
	IsNew               bool        `json:"isNew,omitempty"`
}

type ByPassInfo struct {
	RespHrdrs RespHdrs `json:"respHdrs,omitempty"`
}
type RespHdrs struct {
	ContentType       []string `json:"content-type,omitempty"`
	X_Jenkins_Session []string `json:"x-jenkins-session,omitempty"`
}

type AppDomaiN struct {
	AppDomainID string `json:"id,omitempty"`
	AppDomain   string `json:"domain,omitempty"`
}

type PATH struct {
	Path      string   `json:"path,omitempty"`
	ByPassURI bool     `json:"bypassURI,omitempty"`
	IsDefault bool     `json:"isDefault,omitempty"`
	Status    string   `json:"status,omitempty"`
	IsNewPath bool     `json:"isNewPath,omitempty"`
	Settings  Settings `json:"settings,omitempty"`
}

type Settings struct {
	UserIDIgnored         bool   `json:"userIDIgnored,omitempty"`
	QueryParamaterIgnored bool   `json:"queryParameterIgnored,omitempty"`
	Type                  string `json:"type,omitempty"`
	CacheControlIgnored   bool   `json:"cacheControlIgnored,omitempty"`
	CookieIgnored         bool   `json:"cookieIgnored,omitempty"`
	TTL                   Ttl    `json:"ttl,omitempty"`
}

type Ttl struct {
	Enabled  bool   `json:"enabled,omitempty"`
	Time     int    `json:"time,omitempty"`
	TimeUnit string `json:"timeUnit,omitempty"`
}

func (prosimoClient *ProsimoClient) GetCacheRule(ctx context.Context) ([]*CacheRule, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", CacheRuleEndpoint, nil)
	if err != nil {
		return nil, err
	}

	cacheRuleListData := &CacheRuleListData{}
	_, err = prosimoClient.api_client.Do(ctx, req, cacheRuleListData)
	if err != nil {
		return nil, err
	}

	return cacheRuleListData.Data, nil

}

func (prosimoClient *ProsimoClient) GetCacheRuleByName(ctx context.Context, cacheRuleName string) (*CacheRule, error) {

	cacheRuleList, err := prosimoClient.GetCacheRule(ctx)
	if err != nil {
		return nil, err
	}

	var cacheRule *CacheRule
	for _, returnedCacheRule := range cacheRuleList {
		if returnedCacheRule.Name == cacheRuleName {
			cacheRule = returnedCacheRule
			break
		}
	}

	if cacheRule == nil {
		return nil, errors.New("CacheRule doesnt exists")
	}

	return cacheRule, nil

}

func (prosimoClient *ProsimoClient) GetCacheRuleByID(ctx context.Context, cacheRuleID string) (*CacheRule, error) {

	cacheRuleList, err := prosimoClient.GetCacheRule(ctx)
	if err != nil {
		return nil, err
	}

	var cacheRule *CacheRule
	for _, returnedCacheRule := range cacheRuleList {
		if returnedCacheRule.ID == cacheRuleID {
			cacheRule = returnedCacheRule
			break
		}
	}

	if cacheRule == nil {
		return nil, errors.New("CacheRule doesnt exists")
	}

	return cacheRule, nil

}

func (prosimoClient *ProsimoClient) CreateCacheRule(ctx context.Context, NewCacheRule *CacheRule) (*CacheRuleData, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", CacheRuleEndpoint, NewCacheRule)
	if err != nil {
		return nil, err
	}

	cacheRuleListData := &CacheRuleData{}
	_, err = prosimoClient.api_client.Do(ctx, req, cacheRuleListData)
	if err != nil {
		return nil, err
	}

	return cacheRuleListData, nil

}

func (prosimoClient *ProsimoClient) UpdateCacheRule(ctx context.Context, NewCacheRule *CacheRule) (*CacheRuleData, error) {
	updateCacheRuleEndpoint := fmt.Sprintf("%s/%s", CacheRuleEndpoint, NewCacheRule.ID)

	req, err := prosimoClient.api_client.NewRequest("PUT", updateCacheRuleEndpoint, NewCacheRule)
	if err != nil {
		return nil, err
	}

	cacheRuleListData := &CacheRuleData{}
	_, err = prosimoClient.api_client.Do(ctx, req, cacheRuleListData)
	if err != nil {
		return nil, err
	}

	return cacheRuleListData, nil

}

func (prosimoClient *ProsimoClient) DeleteCacheRule(ctx context.Context, cacheruleID string) error {
	DeleteCacheRuleEndpoint := fmt.Sprintf("%s/%s", CacheRuleEndpoint, cacheruleID)

	req, err := prosimoClient.api_client.NewRequest("DELETE", DeleteCacheRuleEndpoint, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}
