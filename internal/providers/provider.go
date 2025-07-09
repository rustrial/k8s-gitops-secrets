package providers

import (
	"context"
	"encoding/json"
	"fmt"

	secretsv1beta1 "github.com/rustrial/k8s-gitops-secrets/api/secrets/v1beta1"
)

var providerFactories []ProviderFactory = make([]ProviderFactory, 0)

// RegisterProviderFactory registers ProviderFactory instance.
func RegisterProviderFactory(factory ProviderFactory) {
	providerFactories = append(providerFactories, factory)
}

type Audience interface {
	Audience() map[string]string
}

// ProviderFactory implements logic to decrypt and envelop encrypted secret value.
type ProviderFactory interface {
	// Return DecryptionProvider or nil if it does not support the provided Provider spec.
	NewProvider(ctx context.Context, provider *secretsv1beta1.Provider) (DecryptionProvider, error)
}

// DecryptionProvider implements logic to decrypt and envelop encrypted secret value.
type DecryptionProvider interface {
	// Extract audience information from the object.
	Audience(ctx context.Context, sealedSecret *secretsv1beta1.SealedSecret, envelope *secretsv1beta1.Envelope) Audience
	// Decrypt envelope encrypted secret value.
	Decrypt(ctx context.Context, envelope *secretsv1beta1.Envelope, audience Audience) ([]byte, error)
}

// GetProvider returns provider
func GetProvider(ctx context.Context, provider *secretsv1beta1.Provider) (DecryptionProvider, error) {
	if provider != nil {
		for _, factory := range providerFactories {
			provider, err := factory.NewProvider(ctx, provider)
			if err != nil {
				return nil, fmt.Errorf("Error while check provider: %w", err)
			}
			if provider != nil {
				return provider, nil
			}
		}
		txt, _ := json.Marshal(provider)
		return nil, fmt.Errorf("Invalid (unsupported) provider specification: %s", txt)
	}
	return nil, fmt.Errorf("Empty (nil) provider")
}
