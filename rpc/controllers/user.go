package controllers

// @Summary Get subscriptions，当前订阅列表
// @Description Get subscriptions supporing paging
// @Router /subscriptions [get]
// @Tags Subscription
// @Param offset query int false "0"
// @Param limit query int false "10"
// @Success 200 {object} GeneralResponse{data=[]SubscriptionResponse} "code 0"
// @failure 400 {object} GeneralResponse "code 1 with msg"
// @failure 500 {object} GeneralResponse "code 2 with msg"
//func (rpc *RpcController) GetSubscriptions(c *gin.Context) {
//	ResponsePaging(c, http.StatusOK, 0, "", pagingResult, subscriptions)
//}
