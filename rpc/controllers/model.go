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

type ProjectPagingResponse struct {
	PagingResponse
	HasAny bool `json:"has_any"` // has any project in database regardless of status
}

type NewsRequest struct {
	UniqueId string `json:"unique_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Author   string `json:"author" binding:"required"`
	PicUrl   string `json:"pic_url" binding:"required"`
	//Html     string `json:"html" binding:"required"`
	Url  string `json:"url"`
	Type string `json:"type"`
	Lang string `json:"lang" binding:"required"`
}

type NewsBriefResponse struct {
	UniqueId  string `json:"unique_id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	PicUrl    string `json:"pic_url"`
	Url       string `json:"url"`
	Timestamp string `json:"timestamp"`
}

type NewsResponse struct {
	UniqueId string `json:"unique_id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	//Html      string `json:"html"`
	PicUrl    string `json:"pic_url"`
	Url       string `json:"url"`
	Timestamp string `json:"timestamp"`
}

type ProjectResponse struct {
	//Id                int      `json:"id"`
	UniqueId                   string   `json:"unique_id"`
	Title                      string   `json:"title"`
	PicUrl                     string   `json:"pic_url"`
	LogoPicUrl                 string   `json:"logo_pic_url"`
	TargetReturn               string   `json:"target_return"`
	TargetSize                 string   `json:"target_size"`
	Lockup                     string   `json:"lockup"`
	InvestmentHorizon          string   `json:"investment_horizon"`
	StrategyType               string   `json:"strategy_type"`
	Geographies                string   `json:"geographies"`
	Industry                   string   `json:"industry"`
	Website                    string   `json:"website"`
	StartAtTimestamp           int64    `json:"start_at_timestamp"`
	EndAtTimestamp             int64    `json:"end_at_timestamp"`
	TradingPlatform            string   `json:"trading_platform"`
	TradingPlatformDescription string   `json:"trading_platform_description"`
	TradingStatusDescription   string   `json:"trading_status_description"`
	Status                     string   `json:"status"`
	Tags                       []string `json:"tags"`
}

type ArticleResponse struct {
	ArticlePicUrl         string `json:"article_pic_url"`
	Html                  string `json:"html"`
	Author                string `json:"author"`
	AuthorTitle           string `json:"author_title"`
	AvatarUrl             string `json:"avatar_url"`
	AuthorDescriptionHtml string `json:"author_description_html"`
}

type ProjectDetailResponse struct {
	ProjectResponse
	ArticleResponse
}

type ProjectRequest struct {
	Title             string   `json:"title" binding:"required"`
	UniqueId          string   `json:"unique_id" binding:"required"`
	PicUrl            string   `json:"pic_url"`
	LogoPicUrl        string   `json:"logo_pic_url" binding:"required"`
	TargetReturn      string   `json:"target_return"`
	TargetSize        string   `json:"target_size"`
	Lockup            string   `json:"lockup"`
	InvestmentHorizon string   `json:"investment_horizon"`
	StrategyType      string   `json:"strategy_type"`
	Geographies       string   `json:"geographies"`
	Industry          string   `json:"industry"`
	Website           string   `json:"website"`
	StartAtTimestamp  int64    `json:"start_at_timestamp" binding:"required"`
	EndAtTimestamp    int64    `json:"end_at_timestamp" binding:"required"`
	TradingPlatform   string   `json:"trading_platform"`
	Status            string   `json:"status" binding:"required"`
	Lang              string   `json:"lang" binding:"required"`
	Tags              []string `json:"tags"`
}

type ArticleRequest struct {
	PicUrl                string `json:"pic_url" binding:"required"`
	Html                  string `json:"html" binding:"required"`
	Author                string `json:"author" binding:"required"`
	AuthorTitle           string `json:"author_title" binding:"required"`
	AvatarUrl             string `json:"avatar_url" binding:"required"`
	AuthorDescriptionHtml string `json:"author_description_html" binding:"required"`
	Lang                  string `json:"lang" binding:"required"`
	Platform              string `json:"platform" binding:"required"`
	ProjectUniqueId       string `json:"project_unique_id" binding:"required"`
}

type TagsRequest struct {
	ProjectUniqueId string   `json:"unique_id" binding:"required"`
	Tags            []string `json:"tags" binding:"required"`
}

type ContactResponse struct {
	Id        uint   `json:"id"`
	Name      string `json:"name" example:"Ryan Gan"`
	Email     string `json:"email" example:"latifrons88@gmail.com"`
	Subject   string `json:"subject" example:"I want to try"`                   // Subject to record
	Message   string `json:"message" example:"I want to try this product.\n\n"` // Message to record
	Timestamp int64  `json:"timestamp"`
	Ip        string `json:"ip"`
}

type ContactRequest struct {
	Name    string `json:"name" example:"Ryan Gan" binding:"required"`
	Email   string `json:"email" example:"latifrons88@gmail.com" binding:"required"`
	Subject string `json:"subject" example:"I want to try"`                   // Subject to record
	Message string `json:"message" example:"I want to try this product.\n\n"` // Message to record
}

type SubscribeRequest struct {
	Email string `json:"email" example:"latifrons88@gmail.com" binding:"required"`
}

type SubscriptionResponse struct {
	Id        uint   `json:"id"`
	Email     string `json:"email" example:"latifrons88@gmail.com"`
	Timestamp int64  `json:"timestamp"`
	Ip        string `json:"ip"`
}

type DebugUAResponse struct {
	BrowserName string `json:"browser_name"`
	DeviceType  string `json:"device_type"`
	OsName      string `json:"os_name"`
	OsPlatform  string `json:"os_platform"`
	DbPlatform  string `json:"db_platform"`
}
