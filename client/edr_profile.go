package client

import (
	"context"
	"fmt"
)

type EDR_Profile_List struct {
	CrowdStrike []EDR_Profile `json:"CrowdStrike,omitempty"`
}

type EDR_Profile struct {
	AuditID     string   `json:"auditID,omitempty"`
	Id          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Vendor      string   `json:"vendor,omitempty"`
	Criteria    CRITERIA `json:"criteria,omitempty"`
	CreatedTime string   `json:"createdTime,omitempty"`
	UpdatedTime string   `json:"updatedTime,omitempty"`
	Dpps        []dpps   `json:"dpps,omitempty"`
}

type dpps struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type CRITERIA struct {
	SensorActive string   `json:"sensorActive,omitempty"`
	Status       string   `json:"status,omitempty"`
	Ztascore     ZtaScore `json:"ztaScore,omitempty"`
}

type ZtaScore struct {
	From int `json:"from,omitempty"`
	To   int `json:"to,omitempty"`
}

type EDR_Profile_Res struct {
	EdrProfileRes EDR_Profile_List `json:"data,omitempty"`
}
type EDR_Profile_Update_Res struct {
	EdrProfileUpdateRes EDR_Profile `json:"data,omitempty"`
}

func (ztascore ZtaScore) String() string {
	return fmt.Sprintf("{From:%d, To :%d}",
		ztascore.From, ztascore.To)
}

func (criteria CRITERIA) String() string {
	return fmt.Sprintf("{SensorActive:%s, Status :%s}",
		criteria.SensorActive, criteria.Status)
}

func (edr_profile EDR_Profile) String() string {
	return fmt.Sprintf("{Id:%s, CreatedTime:%s, UpdatedTime:%s, Name:%s, Vendor:%s}",
		edr_profile.Id, edr_profile.CreatedTime, edr_profile.UpdatedTime, edr_profile.Name, edr_profile.Vendor)
}

func (prosimoClient *ProsimoClient) GetEDRProfile(ctx context.Context) (*EDR_Profile_Res, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", EDRProfileEndpoint, nil)
	if err != nil {
		return nil, err
	}
	edrres := &EDR_Profile_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, edrres)

	if err != nil {
		return nil, err
	}

	return edrres, nil

}

func (prosimoClient *ProsimoClient) GetEDRProfileByName(ctx context.Context, profilename string) (bool, *EDR_Profile, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", EDRProfileEndpoint, nil)
	if err != nil {
		return false, nil, err
	}
	edrres := &EDR_Profile_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, edrres)

	if err != nil {
		return false, nil, err
	}
	returnprofile := &EDR_Profile{}
	flag := false
	for _, profile := range edrres.EdrProfileRes.CrowdStrike {
		if profile.Name == profilename {
			returnprofile = &profile
			flag = true
		}
	}

	return flag, returnprofile, nil

}

func (prosimoClient *ProsimoClient) UpdateEDRProfile(ctx context.Context, edrprofileops []EDR_Profile) (*EDR_Profile_Update_Res, error) {

	updateEDRProfileEndpoint := fmt.Sprintf("%s/%s", EDRProfileEndpoint, "vendor/CrowdStrike")
	req, err := prosimoClient.api_client.NewRequest("PUT", updateEDRProfileEndpoint, edrprofileops)
	if err != nil {
		return nil, err
	}

	edrprofileputData := &EDR_Profile_Update_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, edrprofileputData)
	if err != nil {
		return nil, err
	}

	return edrprofileputData, nil

}
