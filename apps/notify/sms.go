package notify

const (
	// DEFAULT_TENCENT_SMS_ENDPOINT todo
	DEFAULT_TENCENT_SMS_ENDPOINT = "sms.tencentcloudapi.com"
)

func NewSmsConfig() *SmsConfig {
	return &SmsConfig{
		Enabled:  false,
		Provider: PROVIDER_TENCENT,
		Tencent:  NewTencentSmsConfig(),
	}
}

// NewTencentSmsConfig todo
func NewTencentSmsConfig() *TencentSmsConfig {
	return &TencentSmsConfig{
		Endpoint: DEFAULT_TENCENT_SMS_ENDPOINT,
	}
}
