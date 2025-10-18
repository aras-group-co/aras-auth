package provider

import (
	"fmt"
	"sync"

	"github.com/aras-services/aras-auth/internal/domain"
)

type ProviderRegistry struct {
	providers       map[string]domain.IdentityProvider
	defaultProvider string
	mu              sync.RWMutex
}

func NewProviderRegistry() domain.ProviderRegistry {
	return &ProviderRegistry{
		providers: make(map[string]domain.IdentityProvider),
	}
}

func (r *ProviderRegistry) RegisterProvider(provider domain.IdentityProvider) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := provider.GetProviderName()
	if name == "" {
		return fmt.Errorf("provider name cannot be empty")
	}

	if _, exists := r.providers[name]; exists {
		return fmt.Errorf("provider %s already registered", name)
	}

	r.providers[name] = provider

	// Set as default if it's the first provider or if it's the local provider
	if r.defaultProvider == "" || name == "local" {
		r.defaultProvider = name
	}

	return nil
}

func (r *ProviderRegistry) GetProvider(name string) (domain.IdentityProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, exists := r.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	if !provider.IsEnabled() {
		return nil, fmt.Errorf("provider %s is disabled", name)
	}

	return provider, nil
}

func (r *ProviderRegistry) GetDefaultProvider() domain.IdentityProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.defaultProvider == "" {
		return nil
	}

	provider, exists := r.providers[r.defaultProvider]
	if !exists || !provider.IsEnabled() {
		return nil
	}

	return provider
}

func (r *ProviderRegistry) ListProviders() []domain.IdentityProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var providers []domain.IdentityProvider
	for _, provider := range r.providers {
		providers = append(providers, provider)
	}

	return providers
}

func (r *ProviderRegistry) GetEnabledProviders() []domain.IdentityProvider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var providers []domain.IdentityProvider
	for _, provider := range r.providers {
		if provider.IsEnabled() {
			providers = append(providers, provider)
		}
	}

	return providers
}


