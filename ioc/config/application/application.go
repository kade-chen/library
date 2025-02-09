package application

import (
	"os"

	"github.com/kade-chen/library/ioc"
)

func init() {
	ioc.Config().Registry(&Application{
		AppName:      "default",
		Domain:       "localhost",
		EncryptKey:   "defualt app encrypt key",
		CipherPrefix: "@ciphered@",
	})
}

type Application struct {
	ioc.ObjectImpl

	AppName        string `json:"name" yaml:"name" toml:"name" env:"APP_NAME"`
	AppDescription string `json:"description" yaml:"description" toml:"description" env:"APP_DESCRIPTION"`
	Domain         string `json:"domain" yaml:"domain" toml:"domain" env:"APP_DOMAIN"`
	EncryptKey     string `json:"encrypt_key" yaml:"encrypt_key" toml:"encrypt_key" env:"APP_ENCRYPT_KEY"`
	CipherPrefix   string `json:"cipher_prefix" yaml:"cipher_prefix" toml:"cipher_prefix" env:"APP_CIPHER_PREFIX"`
}

func (i *Application) Init() error {
	sn := os.Getenv("OTEL_SERVICE_NAME")
	if sn == "" {
		os.Setenv("OTEL_SERVICE_NAME", i.AppName)
	}
	return nil
}

func (i *Application) Name() string {
	return AppName
}

// 优先初始化, 以供后面的组件使用
func (i *Application) Priority() int {
	return 9999
}

func (i *Application) GetAppNameWithDefault(defaultValue string) string {
	if i.AppName != "" {
		return i.AppName
	}
	return defaultValue
}
