package controllers

import (
	"github.com/atom-eight/tmt-backend/dbgorm"
	"github.com/atom-eight/tmt-backend/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//
// @Summary Get projects by filter，获取项目列表
// @Description Get projects by filter，默认返回全部，除非用status做指定（locale无关，给status_open或status_closed）
// @Router /cms/projects [get]
// @Tags STO Project
// @Param status query string false "status_open"
// @Param offset query int false "0"
// @Param limit query int false "10"
// @Success 200 {object} GeneralResponse{data=[]ProjectResponse} "code 0"
// @failure 400 {object} GeneralResponse "code 1 with msg"
// @failure 500 {object} GeneralResponse "code 2 with msg"
func (rpc *RpcController) GetProjects(c *gin.Context) {
	locale := rpc.GetLocale(c)

	filterGeographies, err := rpc.ToStringArray(c.DefaultQuery("geographies", ""))
	filterIndustry, err := rpc.ToStringArray(c.DefaultQuery("industry", ""))

	var projects []ProjectResponse
	// basic project columns
	dbProjects, pagingResult, err := rpc.DbOperator.LoadProjects(tools.GetContextDefault(),
		dbgorm.ProjectFilter{
			Lang:        rpc.GetLocale(c),
			Status:      c.DefaultQuery("status", ""),
			Geographies: filterGeographies,
			Industries:  filterIndustry,
		},
		rpc.extractPagingQuery(c))

	if rpc.ResponseError(c, err) {
		return
	}
	projects = make([]ProjectResponse, len(dbProjects))
	for i, dbProject := range dbProjects {
		p := toProjectResponse(dbProject)
		rpc.translateProject(locale, &dbProject, &p)

		// tags
		tags, err := rpc.DbOperator.LoadTags(tools.GetContextDefault(), dbProject.ID)
		if rpc.ResponseError(c, err) {
			return
		}
		p.Tags = toTags(tags)
		projects[i] = p
	}

	hasAny, err := rpc.DbOperator.HasAnyProject(tools.GetContextDefault(),
		dbgorm.ProjectFilter{
			Lang: rpc.GetLocale(c),
		})
	if rpc.ResponseError(c, err) {
		return
	}

	ResponseProjectPaging(c, http.StatusOK, 0, "", pagingResult, projects, hasAny)
}

//
// @Summary Get one project detail with article，项目详情页API
// @Description Get project detail with article
// @Router /cms/project [get]
// @Param projectUniqueId query string true "Project Unique ID"
// @Tags STO Project
// @Success 200 {object} GeneralResponse{data=ProjectDetailResponse} "code 0"
// @failure 404 {object} GeneralResponse "code 1 with msg"
// @failure 500 {object} GeneralResponse "code 2 with msg"
func (rpc *RpcController) GetProjectDetail(c *gin.Context) {

	locale := rpc.GetLocale(c)
	platform := rpc.GetPlatform(c)

	uniqueId := c.Query("projectUniqueId")
	if rpc.ResponseEmptyQuery(c, uniqueId) {
		return
	}

	// basic project columns
	dbProject, err := rpc.DbOperator.LoadProject(tools.GetContextDefault(), uniqueId, locale)
	if rpc.ResponseError(c, err) {
		return
	}

	var projectDetail ProjectDetailResponse
	projectDetail.ProjectResponse = toProjectResponse(dbProject)
	rpc.translateProject(locale, &dbProject, &projectDetail.ProjectResponse)

	// tags
	tags, err := rpc.DbOperator.LoadTags(tools.GetContextDefault(), dbProject.ID)
	if rpc.ResponseError(c, err) {
		return
	}
	projectDetail.Tags = toTags(tags)

	// load article
	dbArticle, err := rpc.DbOperator.LoadLatestArticle(tools.GetContextDefault(), uniqueId, locale, platform)
	if err != nil {
		Response(c, http.StatusNotFound, 1, "article not found", nil)
		return
	}
	projectDetail.ArticleResponse = toArticleResponse(dbArticle)

	Response(c, http.StatusOK, 0, "", projectDetail)
}

// @Summary Post new project，发布新Project（内部接口）
// @Description Post new project
// @Router /cms/project [post]
// @Param data body ProjectRequest true "Project info with locale specified"
// @Tags STO Project
// @Success 200 {object} GeneralResponse{data=int} "code 0"
// @failure 404 {object} GeneralResponse "code 1 with msg"
// @failure 500 {object} GeneralResponse "code 2 with msg"
func (rpc *RpcController) PostProject(c *gin.Context) {
	var req ProjectRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Response(c, http.StatusBadRequest, 1, "bad input", nil)
		return
	}

	rpc.DbOperator.Begin()

	// validation check
	id, err := rpc.DbOperator.InsertProject(tools.GetContextDefault(), dbgorm.DbProject{
		Title:             req.Title,
		UniqueId:          req.UniqueId,
		LogoPicUrl:        req.LogoPicUrl,
		PicUrl:            req.PicUrl,
		TargetReturn:      req.TargetReturn,
		TargetSize:        req.TargetSize,
		Lockup:            req.Lockup,
		InvestmentHorizon: req.InvestmentHorizon,
		StrategyType:      req.StrategyType,
		Geographies:       req.Geographies,
		Industry:          req.Industry,
		Website:           req.Website,
		StartAt:           time.Unix(req.StartAtTimestamp, 0),
		EndAt:             time.Unix(req.EndAtTimestamp, 0),
		TradingPlatform:   req.TradingPlatform,
		Status:            req.Status,
		Lang:              req.Lang,
		Weight:            0,
	})
	if rpc.ResponseError(c, err) {
		rpc.DbOperator.Rollback()
		return
	}
	// create tags
	if req.Tags != nil {
		dbtags := []dbgorm.DbTag{}
		for _, v := range req.Tags {
			dbtags = append(dbtags, dbgorm.DbTag{
				Tag:       v,
				ProjectId: id,
			})
		}
		err = rpc.DbOperator.InsertTags(tools.GetContextDefault(), dbtags)
		if rpc.ResponseError(c, err) {
			rpc.DbOperator.Rollback()
			return
		}
	}
	rpc.DbOperator.Commit()

	Response(c, http.StatusOK, 0, "", id)
	return
}

// @Summary Post new article，发布Project的详细介绍文章
// @Description Post new article under the project，（不绑定project）
// @Router /cms/article [post]
// @Param data body ArticleRequest true "Article info with locale specified"
// @Tags STO Project
// @Success 200 {object} GeneralResponse{data=int} "code 0"
// @failure 404 {object} GeneralResponse "code 1 with msg"
// @failure 500 {object} GeneralResponse "code 2 with msg"
func (rpc *RpcController) PostArticle(c *gin.Context) {
	var req ArticleRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Response(c, http.StatusBadRequest, 1, "bad input", nil)
		return
	}

	// validation check
	id, err := rpc.DbOperator.InsertArticle(tools.GetContextDefault(), dbgorm.DbArticle{
		ProjectUniqueId:       req.ProjectUniqueId,
		PicUrl:                req.PicUrl,
		Html:                  req.Html,
		Author:                req.Author,
		AuthorTitle:           req.AuthorTitle,
		AvatarUrl:             req.AvatarUrl,
		AuthorDescriptionHtml: req.AuthorDescriptionHtml,
		Lang:                  req.Lang,
		Platform:              req.Platform,
	})
	if rpc.ResponseError(c, err) {
		return
	}
	Response(c, http.StatusOK, 0, "", id)
	return
}

func toProjectResponse(v dbgorm.DbProject) (p ProjectResponse) {
	p = ProjectResponse{
		//Id:                v.Id,
		UniqueId:                   v.UniqueId,
		Title:                      v.Title,
		PicUrl:                     v.PicUrl,
		LogoPicUrl:                 v.LogoPicUrl,
		TargetReturn:               v.TargetReturn,
		TargetSize:                 v.TargetSize,
		Lockup:                     v.Lockup,
		InvestmentHorizon:          v.InvestmentHorizon,
		StrategyType:               v.StrategyType,
		Geographies:                v.Geographies,
		Industry:                   v.Industry,
		Website:                    v.Website,
		StartAtTimestamp:           v.StartAt.Unix(),
		EndAtTimestamp:             v.EndAt.Unix(),
		TradingPlatform:            v.TradingPlatform,
		TradingPlatformDescription: "",
		TradingStatusDescription:   "",
		Status:                     v.Status,
		Tags:                       nil,
	}
	return p
}

func toArticleResponse(article dbgorm.DbArticle) (a ArticleResponse) {
	a = ArticleResponse{
		ArticlePicUrl:         article.PicUrl,
		Html:                  article.Html,
		Author:                article.Author,
		AvatarUrl:             article.AvatarUrl,
		AuthorDescriptionHtml: article.AuthorDescriptionHtml,
		AuthorTitle:           article.AuthorTitle,
	}
	return
}

func toTags(tags []dbgorm.DbTag) (vs []string) {
	vs = make([]string, len(tags))
	for i, v := range tags {
		vs[i] = v.Tag
	}
	return
}
