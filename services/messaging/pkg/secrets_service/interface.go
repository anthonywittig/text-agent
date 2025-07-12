package secrets_service

import "context"

type SecretsService interface {
	GetSecret(ctx context.Context, key string) (string, error)
}
