package controllers

import (
	"github.com/atom-eight/tmt-backend/dbgorm"
	"github.com/atom-eight/tmt-backend/tools"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Get News List，新闻列表。Latest Thinking
// @Description Get News，Latest Thinking. 文章语言会根据Cookie里的Locale不同而不同
// @Router /cms/news [get]
// @Tags Latest Thinking
// @Param offset query int false "0"
// @Param limit query int false "10"
// @Success 200 {object} GeneralResponse{data=[]NewsBriefResponse} "code 0"
// @failure 400 {object} GeneralResponse "code 1 with msg"
// @failure 500 {object} GeneralResponse "code 2 with msg"
func (rpc *RpcController) GetNewsList(c *gin.Context) {
	locale := rpc.GetLocale(c)
	newsType := c.DefaultQuery("query", "latest_thinking")
	// basic project columns
	dbElems, pagingResult, err := rpc.DbOperator.LoadVisibleNewsList(tools.GetContextDefault(), newsType,
		locale, rpc.extractPagingQuery(c))

	if rpc.ResponseError(c, err) {
		return
	}

	responseList := make([]NewsBriefResponse, len(dbElems))

	for i, dbElem := range dbElems {
		p := toNewsBriefResponseItem(&dbElem)
		responseList[i] = p
	}
	ResponsePaging(c, http.StatusOK, 0, "", pagingResult, responseList)
}

// @Summary Get News，详细新闻信息。
// @Description Get News，Latest Thinking. 文章语言会根据Cookie里的Locale不同而不同
// @Router /cms/news [get]
// @Tags Latest Thinking
// @Param newsUniqueId query string true "News Unique ID，相同文章会有相同的UniqueId"
// @Param offset query int false "0"
// @Param limit query int false "10"
// @Success 200 {object} GeneralResponse{data=[]NewsResponse} "code 0"
// @failure 400 {object} GeneralResponse "code 1 with msg"
// @failure 500 {object} GeneralResponse "code 2 with msg"
func (rpc *RpcController) GetNews(c *gin.Context) {
	locale := rpc.GetLocale(c)

	uniqueId := c.Query("newsUniqueId")
	if rpc.ResponseEmptyQuery(c, uniqueId) {
		return
	}

	// basic project columns
	dbElem, err := rpc.DbOperator.LoadSingleNews(tools.GetContextDefault(), uniqueId, locale)

	if rpc.ResponseError(c, err) {
		return
	}

	p := toNewsResponseItem(&dbElem)
	Response(c, http.StatusOK, 0, "", p)
}

func toNewsBriefResponseItem(elem *dbgorm.DbNews) NewsBriefResponse {
	return NewsBriefResponse{
		UniqueId:  elem.UniqueId,
		Title:     elem.Title,
		Author:    elem.Author,
		PicUrl:    elem.PicUrl,
		Url:       elem.Url,
		Timestamp: elem.CreatedAt.Format(ShortDateFormatPattern),
	}
}

func toNewsResponseItem(elem *dbgorm.DbNews) NewsResponse {
	return NewsResponse{
		UniqueId:  elem.UniqueId,
		Title:     elem.Title,
		Author:    elem.Author,
		PicUrl:    elem.PicUrl,
		Url:       elem.Url,
		Timestamp: elem.CreatedAt.Format(ShortDateFormatPattern),
	}
}

// @Summary Post news，发布新新闻（内部接口）
// @Description Post new project
// @Router /cms/news [post]
// @Param data body NewsRequest true "News info with locale specified"
// @Tags Latest Thinking
// @Success 200 {object} GeneralResponse{data=int} "code 0 with news id created"
// @failure 404 {object} GeneralResponse "code 1 with msg"
// @failure 500 {object} GeneralResponse "code 2 with msg"
func (rpc *RpcController) PostNews(c *gin.Context) {
	var req NewsRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Response(c, http.StatusBadRequest, 1, "bad input", nil)
		return
	}

	// validation check
	id, err := rpc.DbOperator.InsertNews(tools.GetContextDefault(), dbgorm.DbNews{
		UniqueId: req.UniqueId,
		Title:    req.Title,
		Author:   req.Author,
		Url:      req.Url,
		Type:     req.Type,
		Lang:     req.Lang,
		PicUrl:   req.PicUrl,
		Visible:  true,
	})
	if rpc.ResponseError(c, err) {
		return
	}
	Response(c, http.StatusOK, 0, "", id)
	return
}
