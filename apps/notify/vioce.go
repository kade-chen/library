package notify

func NewVoiceConfig() *VoiceConfig {
	return &VoiceConfig{
		Enabled:  false,
		Provider: PROVIDER_TENCENT,
		Tencent:  DefaultTencentVoiceConfig(),
	}
}

func DefaultTencentVoiceConfig() *TencentVoiceConfig {
	return &TencentVoiceConfig{
		Endpoint:   "vms.tencentcloudapi.com",
		Region:     "ap-guangzhou",
		SignMethod: "TC3-HMAC-SHA256",
		ReqMethod:  "POST",
	}
}
