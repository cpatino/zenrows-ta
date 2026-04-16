package model

type DeviceProfile struct {
	ID            string            `json:"id"`
	DeviceType    string            `json:"deviceType" binding:"required,oneof=desktop mobile"`
	WindowSize    WindowSize        `json:"windowSize" binding:"required"`
	UserAgent     string            `json:"userAgent" binding:"required"`
	CountryCode   string            `json:"countryCode" binding:"required,len=2"`
	CustomHeaders map[string]string `json:"customHeaders,omitempty"`
	UserID        string            `json:"-"`
}

type WindowSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}
