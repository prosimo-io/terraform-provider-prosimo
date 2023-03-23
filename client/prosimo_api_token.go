package client

import (
	"context"
	"fmt"
)

type ProsimoApiToken struct {
	Name string  `json:"name,omitempty"`
	Role []*Role `json:"roles,omitempty"`
}

type Roles struct {
	Role []*Role `json:"data,omitempty"`
}

type Role struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (prosimoClient *ProsimoClient) CreateAPIToken(ctx context.Context, prosimoApiTokenOpts *ProsimoApiToken) (*Role, error) {

	// get role id for role name
	roleReq, err := prosimoClient.api_client.NewRequest("GET", RoleEndpoint, nil)
	if err != nil {
		return nil, err
	}

	prosimoRoles := &Roles{}
	_, err = prosimoClient.api_client.Do(ctx, roleReq, prosimoRoles)
	if err != nil {
		return nil, err
	}

	// create api token with role id
	prosimoAPITokenRoles := []*Role{}
	prosimoAPIToken := &ProsimoApiToken{
		Name: prosimoApiTokenOpts.Name,
		Role: prosimoAPITokenRoles,
	}
	for _, prosimoAPITokenOptsRole := range prosimoApiTokenOpts.Role {
		roleID := ""
		for _, role := range prosimoRoles.Role {
			if role.Name == prosimoAPITokenOptsRole.Name {
				roleID = role.ID
				break
			}
		}
		if roleID == "" {
			return nil, fmt.Errorf("Role Name %s doesnt match", prosimoAPITokenOptsRole.Name)
		}
		apiTokenRole := Role{}
		apiTokenRole.ID = roleID
		prosimoAPITokenRoles = append(prosimoAPITokenRoles, &apiTokenRole)
	}

	fmt.Println((prosimoAPITokenRoles[0]))
	fmt.Println((prosimoAPIToken))

	tokenReq, err := prosimoClient.api_client.NewRequest("POST", APITokenEndpoint, prosimoAPIToken)
	if err != nil {
		return nil, err
	}

	_, err = prosimoClient.api_client.Do(ctx, tokenReq, nil)
	if err != nil {
		return nil, err
	}

	return nil, nil

}
