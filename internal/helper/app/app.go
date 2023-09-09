package app

import (
	"TTPanel/internal/helper/errcode"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Ctx *gin.Context
}

type Pager struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	TotalRows int `json:"total_rows"`
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{Ctx: ctx}
}

func (r *Response) ToResponse(data interface{}) {
	hostname, _ := os.Hostname()
	if data == nil {
		data = gin.H{
			"code": 200,
			"msg":  "success",
			"host": hostname,
		}
	} else {
		data = gin.H{
			"code": 200,
			"msg":  "success",
			"data": data,
			"host": hostname,
		}
	}
	r.Ctx.JSON(http.StatusOK, data)
}
func (r *Response) ToResponseMsg(msg any) {
	hostname, _ := os.Hostname()
	data := gin.H{
		"code": 200,
		"msg":  msg,
		"host": hostname,
	}
	r.Ctx.JSON(http.StatusOK, data)
}
func (r *Response) ToResponseList(list interface{}, totalRows int, pageSize int, page int) {
	r.ToResponse(gin.H{
		"list": list,
		"pager": Pager{
			Page:      page,
			PageSize:  pageSize,
			TotalRows: totalRows,
		},
	})
}
func (r *Response) ToErrorResponse(err *errcode.Error) {
	response := gin.H{"code": err.Code(), "msg": err.Msg()}
	details := err.Details()
	if len(details) > 0 {
		response["details"] = details
	}

	r.Ctx.JSON(err.StatusCode(), response)
}
