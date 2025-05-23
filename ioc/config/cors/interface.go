package cors

type CORS struct {
	Enabled        bool     `toml:"enabled" json:"enabled" yaml:"enabled"  env:"ENABLED"`
	AllowedHeaders []string `json:"cors_allowed_headers" yaml:"cors_allowed_headers" toml:"cors_allowed_headers" env:"ALLOWED_HEADERS" envSeparator:","`
	AllowedDomains []string `json:"cors_allowed_domains" yaml:"cors_allowed_domains" toml:"cors_allowed_domains" env:"ALLOWED_DOMAINS" envSeparator:","`
	AllowedMethods []string `json:"cors_allowed_methods" yaml:"cors_allowed_methods" toml:"cors_allowed_methods" env:"ALLOWED_METHODS" envSeparator:","`
	ExposeHeaders  []string `json:"cors_expose_headers" yaml:"cors_expose_headers" toml:"cors_expose_headers" env:"EXPOSE_HEADERS" envSeparator:","`
	AllowCookies   bool     `toml:"cors_allow_cookies" json:"cors_allow_cookies" yaml:"cors_allow_cookies"  env:"ALLOW_COOKIES"`
	// 单位秒, 默认12小时
	MaxAge int `toml:"max_age" json:"max_age" yaml:"max_age"  env:"HTTP_CORS_MAX_AGE"`
}

func Default() *CORS {
	return &CORS{
		Enabled:        true,
		AllowedHeaders: []string{".*"},
		AllowedDomains: []string{".*"},
		AllowedMethods: []string{"HEAD", "OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"},
		MaxAge:         12 * 60 * 60,
	}
}
