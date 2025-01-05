package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Responder struct {
	C *gin.Context
}

type Response struct {
	Success  bool        `json:"success"`
	Message  string      `json:"message"`
	PageInfo *PageInfo   `json:"pageInfo,omitempty"`
	Result   interface{} `json:"result,omitempty"`
	Error    interface{} `json:"error,omitempty"`
}

type PageInfo struct {
	CurrentPage int `json:"currentPage"`
	NextPage    int `json:"nextPage"`
	PrevPage    int `json:"prevPage"`
	TotalPage   int `json:"totalPage"`
	TotalData   int `json:"totalData"`
}

func NewResponse(c *gin.Context) *Responder {
	return &Responder{C: c}
}

func (r *Responder) Success(message string, data interface{}) {
	r.C.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Result:  data,
	})
}

func (r *Responder) Created(message string, data interface{}) {
	r.C.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Result:  data,
	})
}

func (r *Responder) GetAllSuccess(message string, data interface{}, page *PageInfo) {
	r.C.JSON(http.StatusOK, Response{
		Success:  true,
		Message:  message,
		PageInfo: page,
		Result:   data,
	})
}

func (r *Responder) BadRequest(message string, err interface{}) {
	r.C.JSON(http.StatusOK, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func (r *Responder) Unauthorized(message string, err interface{}) {
	r.C.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
	r.C.Abort()
}

func (r *Responder) Forbidden(message string, err interface{}) {
	r.C.JSON(http.StatusForbidden, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
	r.C.Abort()
}

func (r *Responder) NotFound(message string, err interface{}) {
	r.C.JSON(http.StatusNotFound, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
	r.C.Abort()
}

func (r *Responder) InternalServerError(message string, err interface{}) {
	r.C.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
	r.C.Abort()
}

// func Success(c *gin.Context, statusCode int, message string, data interface{}, pageInfo *PageInfo) {
// 	c.JSON(statusCode, Response{
// 		Success:  true,
// 		Message:  message,
// 		Result:   data,
// 		PageInfo: pageInfo,
// 	})
// }

// func Error(c *gin.Context, statusCode int, message string, data interface{}) {
// 	c.JSON(statusCode, Response{
// 		Success: false,
// 		Message: message,
// 		Result:  data,
// 	})
// }

// func BadRequest(c *gin.Context, message string, data interface{}) {
// 	Error(c, http.StatusBadRequest, message, data)
// }

// func InternalServerError(c *gin.Context, message string, data interface{}) {
// 	Error(c, http.StatusInternalServerError, message, data)
// }

// func Created(c *gin.Context, message string, data interface{}, pageInfo *PageInfo) {
// 	Success(c, http.StatusCreated, message, data, pageInfo)
// }

// func Ok(c *gin.Context, message string, data interface{}, pageInfo *PageInfo) {
// 	Success(c, http.StatusOK, message, data, pageInfo)
// }

// func NotFound(c *gin.Context, message string, data interface{}) {
// 	Error(c, http.StatusNotFound, message, data)
// }

// func Unauthorized(c *gin.Context, message string, data interface{}) {
// 	Error(c, http.StatusUnauthorized, message, data)
// }
