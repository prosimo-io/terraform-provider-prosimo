package client

import "context"

type s3CredentialData struct {
	Data []string `json:"data"`
}

type S3Input struct {
	ID     string `json:"id"`
	Region string `json:"region"`
}

func (prosimoClient *ProsimoClient) ReadS3Credentionals(ctx context.Context, S3Input *S3Input) ([]string, error) {

	req, err := prosimoClient.api_client.NewRequest("POST", CloudS3Endpoint, S3Input)
	if err != nil {
		return nil, err
	}

	s3CredentialsData := &s3CredentialData{}
	_, err = prosimoClient.api_client.Do(ctx, req, s3CredentialsData)
	if err != nil {
		return nil, err
	}

	return s3CredentialsData.Data, nil
}
