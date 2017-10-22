package providers

import "time"

type Tokener interface {
	GetToken() string
	GetUsernameUniq() (string, error)
	GetProviderName() string
	GetRefreshToken() string
	GetExpires() *time.Time
}

type Providerer interface {
	Name() string
	GetCallbackUrl(string) (string, error)
	GetToken(string, string) (Tokener, error)
	QueryAuthCode() string
	RefreshToken(string) (Tokener, error)
}

type Providers map[string]Providerer

func (p *Providers) Add(provider Providerer) {
	(*p)[provider.Name()] = provider
}
func (p *Providers) Get(name string) (Providerer, bool) {
	provider, ok := (*p)[name]
	return provider, ok
}
