package client

import (
	"context"
	"fmt"
	"net/url"
)

type Edge struct {
	CloudKeyID       string             `json:"cloudKeyID,omitempty"`
	CloudRegion      string             `json:"cloudRegion,omitempty"`
	ID               string             `json:"id,omitempty"`
	NodeSizesettings *ConnectorSettings `json:"nodeSizeSettings,omitempty"`
	CloudType        string             `json:"cloudType,omitempty"`
	ClusterName      string             `json:"clusterName,omitempty"`
	PappFqdn         string             `json:"pappFqdn,omitempty"`
	RegStatus        string             `json:"regStatus,omitempty"`
	Status           string             `json:"status,omitempty"`
	Subnet           string             `json:"subnet,omitempty"`
	TeamID           string             `json:"teamID,omitempty"`
	Byoresource      *ByoResource       `json:"byoResourceDetails,omitempty"`
}

type ByoResource struct {
	VpcID string `json:"vpcId,omitempty"`
}

type EdgeList struct {
	Edges []*Edge `json:"data,omitempty"`
}

func (edge Edge) String() string {
	return fmt.Sprintf("{ClusterName:%s, PappFqdn:%s, CloudType:%s, Status:%s}",
		edge.ClusterName, edge.PappFqdn, edge.CloudType, edge.Status)
}

func (prosimoClient *ProsimoClient) CreateEdge(ctx context.Context, edge *Edge) (*ResourcePostResponseData, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", EdgeEndpoint, edge)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &ResourcePostResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) GetEdge(ctx context.Context) (*EdgeList, error) {

	req, err := prosimoClient.api_client.NewRequest("GET", EdgeEndpoint, nil)
	if err != nil {
		return nil, err
	}

	edgeList := &EdgeList{}
	_, err = prosimoClient.api_client.Do(ctx, req, edgeList)
	if err != nil {
		return nil, err
	}

	return edgeList, nil

}

func (prosimoClient *ProsimoClient) DeleteEdge(ctx context.Context, edgeId string) error {

	deleteEdgeEndpt := fmt.Sprintf("%s/%s", EdgeEndpoint, edgeId)

	req, err := prosimoClient.api_client.NewRequest("DELETE", deleteEdgeEndpt, nil)
	if err != nil {
		return err
	}

	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}

func (prosimoClient *ProsimoClient) DeployApp(ctx context.Context, edge *Edge) (*ResourcePostResponseData, error) {

	appDeployEndpt := fmt.Sprintf("%s/%s", AppDeploymentEndpoint, edge.ID)

	emptyInterface := &Edge{}
	req, err := prosimoClient.api_client.NewRequest("PUT", appDeployEndpt, emptyInterface)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &ResourcePostResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) DeleteAppDeployment(ctx context.Context, edgeId string) (*ResourcePostResponseData, error) {

	deleteAppDeploymentEndpt := fmt.Sprintf("%s/%s", AppDeploymentEndpoint, edgeId)

	req, err := prosimoClient.api_client.NewRequest("DELETE", deleteAppDeploymentEndpt, &Edge{})
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &ResourcePostResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) ForceDeleteAppDeployment(ctx context.Context, edgeId string) (*ResourcePostResponseData, error) {

	deleteAppDeploymentEndpt := fmt.Sprintf("%s/%s", AppDeploymentEndpoint, edgeId)
	u, err := url.Parse(deleteAppDeploymentEndpt)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return nil, err
	}

	// Add query parameters
	q := u.Query()
	q.Add("force", "true")
	u.RawQuery = q.Encode()

	req, err := prosimoClient.api_client.NewRequest("DELETE", u.String(), &Edge{})
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &ResourcePostResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}

func (prosimoClient *ProsimoClient) PatchSubnetRange(ctx context.Context, edgeId string, patchSubnet *Edge) error {

	patchAppDeploymentEndpt := fmt.Sprintf("%s/%s", PatchEdgeSubnetEndpoint, edgeId)

	req, err := prosimoClient.api_client.NewRequest("PATCH", patchAppDeploymentEndpt, patchSubnet)
	if err != nil {
		return err
	}

	// resourcePostResponseData := &ResourcePostResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil

}

func (prosimoClient *ProsimoClient) UpdateEdge(ctx context.Context, edgeId string, edge *Edge) (*ResourcePostResponseData, error) {

	patchAppDeploymentEndpt := fmt.Sprintf("%s/%s", PatchEdgeEndpoint, edgeId)

	req, err := prosimoClient.api_client.NewRequest("PATCH", patchAppDeploymentEndpt, edge)
	if err != nil {
		return nil, err
	}

	resourcePostResponseData := &ResourcePostResponseData{}
	_, err = prosimoClient.api_client.Do(ctx, req, resourcePostResponseData)
	if err != nil {
		return nil, err
	}

	return resourcePostResponseData, nil

}
