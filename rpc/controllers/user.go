package controllers

import (
	"github.com/atom-eight/tmt-backend/dbgorm"
	"github.com/atom-eight/tmt-backend/tools"
	"github.com/gin-gonic/gin"
	"net/http"
)

//func (rpc *RpcController) test(c *gin.Context) {
//	d, err := rpc.DbOperator.LoadContacts(tools.GetContextDefault(), db.PagingParams{
//		Offset: 0,
//		Limit:  100,
//	})
//	if err != nil {
//		logrus.WithError(err).Error("xxx")
//	}
//	fmt.Println(len(d))
//}
//

// @Summary Create Contact，联系我们
// @Description User submit a new contact info
// @Router /contact [post]
// @Tags Contacts
// @Param data body ContactRequest true "body"
// @Success 200 {object} GeneralResponse{data=int} "code 0 with generated id "
// @failure 400 {object} GeneralResponse "code 1 with msg"
// @failure 500 {object} GeneralResponse "code 2 with msg"
func (rpc *RpcController) PostContact(c *gin.Context) {
	var req ContactRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Response(c, http.StatusBadRequest, 1, "bad input", nil)
		return
	}

	// validation check
	id, err := rpc.DbOperator.InsertContact(tools.GetContextDefault(), dbgorm.DbContact{
		Name:    req.Name,
		Email:   req.Email,
		Subject: req.Subject,
		Message: req.Message,
		Ip:      c.ClientIP(),
	})
	if rpc.ResponseError(c, err) {
		return
	}
	Response(c, http.StatusOK, 0, "", id)
	return
}

//
// @Summary Create Subscription，创建订阅，用户提交Email
// @Description User submit a new subscription info
// @Router /subscription [post]
// @Tags Subscription
// @Param data body SubscribeRequest true "body"
// @Success 200 {object} GeneralResponse{data=int} "code 0 with generated id "
// @failure 400 {object} GeneralResponse "code 1 with msg"
// @failure 500 {object} GeneralResponse "code 2 with msg"
func (rpc *RpcController) PostSubscription(c *gin.Context) {
	var req SubscribeRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		Response(c, http.StatusBadRequest, 1, "bad input", nil)
		return
	}

	// validation check
	id, err := rpc.DbOperator.InsertSubscription(tools.GetContextDefault(), dbgorm.DbSubscription{
		Email: req.Email,
		Ip:    c.ClientIP(),
	})
	if err != nil {

		return
	}
	Response(c, http.StatusOK, 0, "", id)
	return
}

// @Summary Get contacts，获取当前所有提交了的联系人列表
// @Description Get contacts supporing paging
// @Router /contacts [get]
// @Tags Contacts
// @Param offset query int false "0"
// @Param limit query int false "10"
// @Success 200 {object} GeneralResponse{data=[]ContactResponse} "code 0"
// @failure 400 {object} GeneralResponse "code 1 with msg"
// @failure 500 {object} GeneralResponse "code 2 with msg"
func (rpc *RpcController) GetContacts(c *gin.Context) {
	// basic project columns
	dbContactResponse, pagingResponse, err := rpc.DbOperator.LoadContacts(tools.GetContextDefault(), rpc.extractPagingQuery(c))

	if rpc.ResponseError(c, err) {
		return
	}

	contacts := make([]ContactResponse, len(dbContactResponse))
	for i, dbContact := range dbContactResponse {
		p := toContactResponse(dbContact)
		contacts[i] = p
	}
	ResponsePaging(c, http.StatusOK, 0, "", pagingResponse, contacts)
}

func toContactResponse(c dbgorm.DbContact) ContactResponse {
	return ContactResponse{
		Id:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		Subject:   c.Subject,
		Message:   c.Message,
		Timestamp: c.CreatedAt.Unix(),
		Ip:        c.Ip,
	}

}

// @Summary Get subscriptions，当前订阅列表
// @Description Get subscriptions supporing paging
// @Router /subscriptions [get]
// @Tags Subscription
// @Param offset query int false "0"
// @Param limit query int false "10"
// @Success 200 {object} GeneralResponse{data=[]SubscriptionResponse} "code 0"
// @failure 400 {object} GeneralResponse "code 1 with msg"
// @failure 500 {object} GeneralResponse "code 2 with msg"
func (rpc *RpcController) GetSubscriptions(c *gin.Context) {
	// basic project columns
	dbSubscriptions, pagingResult, err := rpc.DbOperator.LoadSubscriptions(tools.GetContextDefault(), rpc.extractPagingQuery(c))

	if rpc.ResponseError(c, err) {
		return
	}

	subscriptions := make([]SubscriptionResponse, len(dbSubscriptions))
	for i, dbSubscription := range dbSubscriptions {
		p := toSubscription(dbSubscription)
		subscriptions[i] = p
	}
	ResponsePaging(c, http.StatusOK, 0, "", pagingResult, subscriptions)
}

func toSubscription(s dbgorm.DbSubscription) SubscriptionResponse {
	return SubscriptionResponse{
		Id:        s.ID,
		Email:     s.Email,
		Timestamp: s.CreatedAt.Unix(),
		Ip:        s.Ip,
	}
}
