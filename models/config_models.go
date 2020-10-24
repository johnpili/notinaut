package models

// Config ...
type Config struct {
	HTTP struct {
		Port       int    `yaml:"port"`
		IsTLS      bool   `yaml:"is_tls"`
		ServerCert string `yaml:"server_cert"`
		ServerKey  string `yaml:"server_key"`
	} `yaml:"http"`

	System struct {
		SerialName string `yaml:"serial_name"`
		SerialBaud int    `yaml:"serial_baud"`
		CookieName string `yaml:"cookie_name"`
		CookieKey  string `yaml:"cookie_key"`
	} `yaml:"system"`

	Extraction struct {
		HeaderKey   string `yaml:"header_key"`
		DebugHeader bool   `yaml:"debug_header"`
	} `yaml:"extraction"`

	IPCacheControl struct {
		ExpirySec int64 `yaml:"expiry_sec"`
		PurgeSec  int64 `yaml:"purge_sec"`
	} `yaml:"ip_cache_control"`
}

// IPInfo ...
type IPInfo struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user-agent"`
}
