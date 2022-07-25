package rpc

// swaggo
// Run swag init -g controller.go
// Copy swagger.json to doc/swagger.json
// Can be rendered by SwaggerUI or Redoc
// Notation Doc: https://github.com/swaggo/swag

// @title Some API
// @version 1.0
// @description Some API
// @BasePath /v1
import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/latifrons/lbserver/berror"
	"github.com/latifrons/lbserver/debug"
	"github.com/latifrons/lbserver/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type RpcController struct {
	Flags        *debug.Flags `container:"type"`
	AllowOrigins []string
}

func (rpc *RpcController) Response(c *gin.Context, status int, code int, msg string, data interface{}) {
	if rpc.Flags.ResponseLog {
		logrus.WithField("data", data).WithField("msg", msg).WithField("code", code).WithField("status", status).Info("resp")
	}
	c.JSON(status, GeneralResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
func (rpc *RpcController) ResponseOK(c *gin.Context, data interface{}) {
	if rpc.Flags.ResponseLog {
		logrus.WithField("data", data).Info("resp ok")
	}
	c.JSON(http.StatusOK, GeneralResponse{
		Code: 0,
		Msg:  "",
		Data: data,
	})
}

func (rpc *RpcController) ResponsePaging(c *gin.Context, pagingResult model.PagingResult, data interface{}, list interface{}) {
	if pagingResult.Limit != 0 {
		c.JSON(http.StatusOK, PagingResponse{
			GeneralResponse: GeneralResponse{
				Code: 0,
				Msg:  "",
				Data: data,
			},
			List:  list,
			Size:  pagingResult.Limit,
			Total: pagingResult.Total,
			Page:  pagingResult.Offset/pagingResult.Limit + 1,
		})
	} else {
		c.JSON(http.StatusOK, PagingResponse{
			GeneralResponse: GeneralResponse{
				Code: 0,
				Msg:  "",
				Data: data,
			},
			List:  list,
			Size:  pagingResult.Limit,
			Total: pagingResult.Total,
			Page:  1,
		})
	}

}

func (rpc *RpcController) ResponseBadRequest(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	logrus.WithError(err).Debug("bad request")
	if rpc.Flags.ReturnDetailError {
		rpc.Response(c, http.StatusBadRequest, berror.ErrBadRequest, err.Error(), nil)
	} else {
		rpc.Response(c, http.StatusBadRequest, berror.ErrBadRequest, "Bad Request. Check your input", nil)
	}
	return true
}

func (rpc *RpcController) ResponseInternalServerError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	logrus.WithError(err).Error("internal error")
	if rpc.Flags.ReturnDetailError {
		rpc.Response(c, http.StatusInternalServerError, berror.ErrInternal, err.Error(), nil)
	} else {
		rpc.Response(c, http.StatusInternalServerError, berror.ErrInternal, "Internal server error", nil)
	}
	return true
}

func (rpc *RpcController) ResponseError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	switch err.(type) {
	case *berror.BError:
		berr := err.(*berror.BError)

		var msg = "fail"
		if rpc.Flags.ReturnDetailError {
			msg = berr.Msg
		}

		// response http code according to DTM
		switch berr.ErrorCategory {
		case berror.CategoryBusinessFail:
			rpc.Response(c, http.StatusOK, berr.Code, msg, nil)
		case berror.CategoryBusinessTemporary:
			rpc.Response(c, http.StatusOK, berr.Code, msg, nil)
		case berror.CategorySystemTemporary:
			rpc.Response(c, http.StatusOK, berr.Code, msg, nil)
		}
	default:
		fmt.Printf("%s", err.Error())
		if v, ok := err.(berror.StackTracer); ok {
			fmt.Println(v.StackTrace())
		} else {
			fmt.Printf("%s", err.Error())
		}

		rpc.Response(c, http.StatusInternalServerError, berror.ErrInternal, err.Error(), nil)
	}
	return true
}

func (rpc *RpcController) ResponseEmptyQuery(c *gin.Context, value string) bool {
	if value == "" {
		logrus.Debug("param missing")
		rpc.Response(c, http.StatusBadRequest, berror.ErrBadRequest, "param missing", nil)
		return true
	}
	return false
}

func (rpc *RpcController) ToStringArray(query string) (arr []string, err error) {
	return strings.Split(query, "$"), nil
}

func (rpc *RpcController) extractOrderQuery(c *gin.Context) (p model.OrderParams, err error) {
	dir := strings.ToUpper(c.DefaultQuery("order_direction", ""))
	if dir != model.OrderDirectionASC && dir != model.OrderDirectionDESC && dir != "" {
		err = berror.NewFail(berror.ErrBadRequest, "bad order direction: "+dir)
		return
	}
	p = model.OrderParams{
		OrderBy:        c.DefaultQuery("order_by", ""),
		OrderDirection: dir,
	}
	return
}

func (rpc *RpcController) mustNotDuplicate(c *gin.Context, key string, nonce uint) bool {
	// TODO duplicate check
	if nonce == 0 {
		rpc.ResponseBadRequest(c, errors.New("duplicate request or empty nonce"))
		return true
	}
	return false
}

func (rpc *RpcController) parseBoolValue(v string) (vout bool, err error) {
	v = strings.ToLower(v)
	v = strings.TrimSpace(v)
	switch v {
	case "0":
		vout = false
	case "false":
		vout = false
	case "1":
		vout = true
	case "true":
		vout = true
	default:
		err = errors.New("bad bool value")
	}
	return
}
