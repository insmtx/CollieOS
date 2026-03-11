package config

type GithubAppConfig struct {
	AppID         int64
	PrivateKey    string
	WebhookSecret string
	BaseURL       string
}
