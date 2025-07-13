package secrets_service

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AwsSecretsService struct {
	secretsClient *secretsmanager.Client
}

func NewAwsSecretsService(ctx context.Context) (SecretsService, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &AwsSecretsService{secretsClient: secretsmanager.NewFromConfig(cfg)}, nil
}

func (s *AwsSecretsService) GetSecret(ctx context.Context, key string) (string, error) {
	secret, err := s.secretsClient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get secret: %w", err)
	}
	return *secret.SecretString, nil
}
