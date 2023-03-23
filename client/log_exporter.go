package client

import (
	"context"
	"fmt"
)

type Log_Config struct {
	Id                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	IP                  string `json:"ip,omitempty"`
	TcpPort             int    `json:"tcpPort,omitempty"`
	TlsEnabled          bool   `json:"tlsEnabled,omitempty"`
	AuthenticationToken string `json:"authenticationToken,omitempty"`
	Description         string `json:"description,omitempty"`
	TeamID              string `json:"teamID,omitempty"`
	Status              string `json:"status,omitempty"`
	CreatedTime         string `json:"createdTime,omitempty"`
	UpdatedTime         string `json:"updatedTime,omitempty"`
}

type Log_Config_Res struct {
	LogConfig Log_Config `json:"data,omitempty"`
}

type Log_Config_ResList struct {
	LogConfigList []*Log_Config `json:"data,omitempty"`
}

func (log_config Log_Config) String() string {
	return fmt.Sprintf("{Id:%s, TeamID:%s, CreatedTime:%s, UpdatedTime:%s, Name:%s, IP:%s, TcpPort:%d, Status:%s, }",
		log_config.Id, log_config.Id, log_config.CreatedTime, log_config.UpdatedTime, log_config.Name, log_config.IP, log_config.TcpPort, log_config.Status)
}

func (prosimoClient *ProsimoClient) GetLogConf(ctx context.Context) (*Log_Config_ResList, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", LogExporterEndpoint, nil)
	if err != nil {
		return nil, err
	}
	logres := &Log_Config_ResList{}
	_, err = prosimoClient.api_client.Do(ctx, req, logres)

	if err != nil {
		return nil, err
	}

	return logres, nil

}

func (prosimoClient *ProsimoClient) CreateLogConf(ctx context.Context, logconfops *Log_Config) (*Log_Config_Res, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", LogExporterEndpoint, logconfops)
	if err != nil {
		return nil, err
	}

	logconfpostData := &Log_Config_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, logconfpostData)
	if err != nil {
		return nil, err
	}

	return logconfpostData, nil

}

func (prosimoClient *ProsimoClient) UpdateLogConf(ctx context.Context, logconfops *Log_Config) (*Log_Config_Res, error) {

	updateLogExporterEndpoint := fmt.Sprintf("%s/%s", LogExporterEndpoint, logconfops.Id)
	req, err := prosimoClient.api_client.NewRequest("PUT", updateLogExporterEndpoint, logconfops)
	if err != nil {
		return nil, err
	}

	logconfputData := &Log_Config_Res{}
	_, err = prosimoClient.api_client.Do(ctx, req, logconfputData)
	if err != nil {
		return nil, err
	}

	return logconfputData, nil

}

func (prosimoClient *ProsimoClient) DeleteLogConf(ctx context.Context, logID string) error {

	updateLogExporterEndpoint := fmt.Sprintf("%s/%s", LogExporterEndpoint, logID)
	req, err := prosimoClient.api_client.NewRequest("DELETE", updateLogExporterEndpoint, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}
