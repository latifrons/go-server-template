package controllers

type GeneralResponse struct {
	Code int         `json:"code" example:"0"` // Code is 0 for normal cases and positive for errors.
	Msg  string      `json:"msg"`              // Msg is "" for normal cases and message for errors.
	Data interface{} `json:"data,omitempty"`   // Optional
}

type PagingResponse struct {
	Code   int         `json:"code" example:"0"` // Code is 0 for normal cases and positive for errors.
	Msg    string      `json:"msg"`              // Msg is "" for normal cases and message for errors.
	Limit  int         `json:"limit"`            // Limit is the result count in this response
	Total  int64       `json:"total"`            // Total is the total result in database
	Offset int         `json:"offset"`           // Offset is the given params in the request starts from 0
	Data   interface{} `json:"data,omitempty"`   // Optional
}

type DebugUAResponse struct {
	BrowserName string `json:"browser_name"`
	DeviceType  string `json:"device_type"`
	OsName      string `json:"os_name"`
	OsPlatform  string `json:"os_platform"`
	DbPlatform  string `json:"db_platform"`
}
