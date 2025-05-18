package configs

type CORSConfig struct {
	AllowOrigin      string `json:"allow_origin"`
	AllowHeaders     string `json:"allow_headers"`
	AllowMethods     string `json:"allow_methods"`
	MaxAge           string `json:"max_age"`
	AllowCredentials string `json:"allow_credentials"`
}

func (c *CORSConfig) LoadAndValidate() error {
	if c.AllowOrigin == "" {
		c.AllowOrigin = "*"
	}
	if c.AllowHeaders == "" {
		c.AllowHeaders = "*"
	}
	if c.AllowMethods == "" {
		c.AllowMethods = "GET, POST, PUT, DELETE, OPTIONS"
	}
	if c.MaxAge == "" {
		c.MaxAge = "86400" // 24 hours
	}
	if c.AllowCredentials == "" {
		c.AllowCredentials = "true"
	}
	return nil
}
