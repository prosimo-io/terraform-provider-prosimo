package client

import (
	"context"
	"fmt"
	"os"
	"time"
)

var backoffSchedule = []time.Duration{
	1 * time.Second,
	3 * time.Second,
	10 * time.Second,
}

type DevicePosture_Settings struct {
	Enabled bool `json:"enabled,omitempty"`
}

type DevicePosture_Profile struct {
	AuditID     string        `json:"auditID,omitempty"`
	TaskID      string        `json:"taskID,omitempty"`
	Id          string        `json:"id,omitempty"`
	Name        string        `json:"name,omitempty"`
	RiskLevel   string        `json:"riskLevel,omitempty"`
	Enabled     bool          `json:"enabled,omitempty"`
	Criteria    CRITERIADP    `json:"criteria,omitempty"`
	CreatedTime string        `json:"createdTime,omitempty"`
	UpdatedTime string        `json:"updatedTime,omitempty"`
	EdrProfiles []EDR_Profile `json:"edrProfiles,omitempty"`
}

type CRITERIADP struct {
	DiskEncryptionStatus string       `json:"diskEncryptionStatus,omitempty"`
	DomainOfInterest     []string     `json:"domainOfInterest,omitempty"`
	FirewallStatus       string       `json:"firewallStatus,omitempty"`
	Os                   string       `json:"os,omitempty"`
	RunningProcess       []string     `json:"runningProcess,omitempty"`
	OsOperator           string       `json:"osOperator,omitempty"`
	OsVersions           []Os_Version `json:"osVersions,omitempty"`
}

type Os_Version struct {
	Build string `json:"build,omitempty"`
	Patch string `json:"patch,omitempty"`
}

type DevicePosture_Profile_Res struct {
	DPProfileRes []DevicePosture_Profile `json:"data,omitempty"`
}
type DevicePosture_Profile_Update_Res struct {
	DPProfileUpdateRes DevicePosture_Profile `json:"data,omitempty"`
}

type TaskStatusPost struct {
	IdList []string `json:"idList,omitempty"`
}

type Task struct {
	TaskID string `json:"taskID,omitempty"`
	Status string `json:"status,omitempty"`
	// Name   string `json:"name,omitempty"`
}

type Item struct {
	TaskItemID string `json:"taskID,omitempty"`
	Status     string `json:"status,omitempty"`
	Name       string `json:"name,omitempty"`
	Details    string `json:"details,omitempty"`
}

// type Context struct {
// 	Type string `json:"type,omitempty"`
// }

type Task_Status struct {
	TaskID   string `json:"taskID,omitempty"`
	ItemList []Item `json:"items,omitempty"`
	Status   string `json:"status,omitempty"`
	// ContextType Context `json:"context,omitempty"`
}

type Task_Status_res struct {
	Records []*Task_Status `json:"data,omitempty"`
}

//	type Task_Status_res_data struct {
//		TaskStatusRes *Task_Status_res `json:"data,omitempty"`
//	}
type DP_Profile_Read_Res struct {
	High   []DevicePosture_Profile `json:"high,omitempty"`
	Medium []DevicePosture_Profile `json:"medium,omitempty"`
	Low    []DevicePosture_Profile `json:"low,omitempty"`
}

func (osversion Os_Version) String() string {
	return fmt.Sprintf("{Build:%s, Patch :%s}",
		osversion.Build, osversion.Patch)
}

func (criteria CRITERIADP) String() string {
	return fmt.Sprintf("{DiskEncryptionStatus:%s, FirewallStatus:%s, Os:%s, OsOperator:%s}",
		criteria.DiskEncryptionStatus, criteria.FirewallStatus, criteria.Os, criteria.OsOperator)
}

func (dp_profile DevicePosture_Profile) String() string {
	return fmt.Sprintf("{Id:%s, CreatedTime:%s, UpdatedTime:%s, Name:%s, RiskLevel:%s, Enabled:%t}",
		dp_profile.Id, dp_profile.CreatedTime, dp_profile.UpdatedTime, dp_profile.Name, dp_profile.RiskLevel, dp_profile.Enabled)
}

func (prosimoClient *ProsimoClient) GetDPProfile(ctx context.Context) (*DP_Profile_Read_Res, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", DPProfileEndpoint, nil)
	if err != nil {
		return nil, err
	}
	dpres := &DP_Profile_Read_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, dpres)

	if err != nil {
		return nil, err
	}

	return dpres, nil

}

func (prosimoClient *ProsimoClient) GetDPProfileBYRiskLevel(ctx context.Context, risk_level string) (*DevicePosture_Profile_Res, error) {
	updateDPProfileEndpoint := fmt.Sprintf("%s/%s/%s", DPProfileEndpoint, "risk-level", risk_level)
	req, err := prosimoClient.api_client.NewRequest("GET", updateDPProfileEndpoint, nil)
	if err != nil {
		return nil, err
	}
	dpres := &DevicePosture_Profile_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, dpres)

	if err != nil {
		return nil, err
	}

	return dpres, nil

}

func (prosimoClient *ProsimoClient) GetDPSettings(ctx context.Context) (*DevicePosture_Settings, error) {

	updateDPProfileEndpoint := fmt.Sprintf("%s/%s", DPProfileEndpoint, "settings")
	req, err := prosimoClient.api_client.NewRequest("GET", updateDPProfileEndpoint, nil)
	if err != nil {
		return nil, err
	}

	dpsettingData := &DevicePosture_Settings{}
	_, err = prosimoClient.api_client.Do(ctx, req, dpsettingData)
	if err != nil {
		return nil, err
	}

	return dpsettingData, nil

}

func (prosimoClient *ProsimoClient) UpdateDPSettings(ctx context.Context, dpprofilesetting DevicePosture_Settings) (*DevicePosture_Profile_Update_Res, error) {

	updateDPProfileEndpoint := fmt.Sprintf("%s/%s", DPProfileEndpoint, "settings")
	req, err := prosimoClient.api_client.NewRequest("PUT", updateDPProfileEndpoint, dpprofilesetting)
	if err != nil {
		return nil, err
	}

	dpprofileputData := &DevicePosture_Profile_Update_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, dpprofileputData)
	if err != nil {
		return nil, err
	}

	return dpprofileputData, nil

}

func (prosimoClient *ProsimoClient) UpdateDPProfileHigh(ctx context.Context, dpprofileops []DevicePosture_Profile) (*DevicePosture_Profile_Update_Res, error) {

	updateDPProfileEndpoint := fmt.Sprintf("%s/%s", DPProfileEndpoint, "risk-level/high")
	req, err := prosimoClient.api_client.NewRequest("PUT", updateDPProfileEndpoint, dpprofileops)
	if err != nil {
		return nil, err
	}

	dpprofileputData := &DevicePosture_Profile_Update_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, dpprofileputData)
	if err != nil {
		return nil, err
	}

	return dpprofileputData, nil

}

func (prosimoClient *ProsimoClient) UpdateDPProfileMedium(ctx context.Context, dpprofileops []DevicePosture_Profile) (*DevicePosture_Profile_Update_Res, error) {

	updateDPProfileEndpoint := fmt.Sprintf("%s/%s", DPProfileEndpoint, "risk-level/medium")
	req, err := prosimoClient.api_client.NewRequest("PUT", updateDPProfileEndpoint, dpprofileops)
	if err != nil {
		return nil, err
	}

	dpprofileputData := &DevicePosture_Profile_Update_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, dpprofileputData)
	if err != nil {
		return nil, err
	}

	return dpprofileputData, nil

}

func (prosimoClient *ProsimoClient) UpdateDPProfilelow(ctx context.Context, dpprofileops []DevicePosture_Profile) (*DevicePosture_Profile_Update_Res, error) {

	updateDPProfileEndpoint := fmt.Sprintf("%s/%s", DPProfileEndpoint, "risk-level/low")
	req, err := prosimoClient.api_client.NewRequest("PUT", updateDPProfileEndpoint, dpprofileops)
	if err != nil {
		return nil, err
	}

	dpprofileputData := &DevicePosture_Profile_Update_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, dpprofileputData)
	if err != nil {
		return nil, err
	}

	return dpprofileputData, nil

}

func (prosimoClient *ProsimoClient) GetTaskStatus(ctx context.Context, taskID string) (*Task_Status, error) {
	var err error
	var taskidList []string
	taskidList = append(taskidList, taskID)
	taskInput := TaskStatusPost{
		IdList: taskidList,
	}

	// UpdatedTaskEndpoint := fmt.Sprintf("%s/%s", TaskEndpoint, taskID)
	req, err := prosimoClient.api_client.NewRequest("POST", TaskEndpointSearch, taskInput)
	if err != nil {
		return nil, err
	}
	taskStatusData := &Task_Status_res{}
	for _, backoff := range backoffSchedule {
		_, err = prosimoClient.api_client.Do(ctx, req, taskStatusData)
		if err == nil {
			break
		}

		fmt.Fprintf(os.Stderr, "Request error: %+v\n", err)
		fmt.Fprintf(os.Stderr, "Retrying in %v\n", backoff)
		time.Sleep(backoff)
	}
	if err != nil {
		return nil, err
	}

	var taskDetails *Task_Status
	for _, taskStatusRes := range taskStatusData.Records {
		taskDetails = taskStatusRes
	}
	return taskDetails, nil
}
