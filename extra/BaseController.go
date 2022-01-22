package extra

import (
	"github.com/gentwolf-shen/gin-boost"
	"github.com/gentwolf-shen/gobootstrap/validator"
	"github.com/gentwolf-shen/gohelper-v2/converter"
	"github.com/gentwolf-shen/gohelper-v2/dict"
)

var (
	SucceedCode    = 200
	SucceedMessage = "succeed"
)

type BaseController struct {
}

// 绑定参数(uri, query, form post, json post)
func (ctl *BaseController) BindRequest(c *gin.Context, p interface{}, cb func(rs *ResponseMessage)) {
	ctl.bind(c, func() error {
		if p == nil {
			return nil
		}

		return c.BindRequest(p)
	}, func(rs *ResponseMessage) {
		cb(rs)
	})
}

// 绑定参数 (form post, json post)
func (ctl *BaseController) BindBody(c *gin.Context, p interface{}, cb func(rs *ResponseMessage)) {
	ctl.bind(c, func() error {
		if p == nil {
			return nil
		}

		return c.ShouldBind(p)
	}, func(rs *ResponseMessage) {
		cb(rs)
	})
}

// 绑定参数 (query)
func (ctl *BaseController) BindQuery(c *gin.Context, p interface{}, cb func(rs *ResponseMessage)) {
	ctl.bind(c, func() error {
		if p == nil {
			return nil
		}

		return c.ShouldBindQuery(p)
	}, func(rs *ResponseMessage) {
		cb(rs)
	})
}

// 绑定参数 (uri)
func (ctl *BaseController) BindUri(c *gin.Context, p interface{}, cb func(rs *ResponseMessage)) {
	ctl.bind(c, func() error {
		if p == nil {
			return nil
		}

		return c.ShouldBindUri(p)
	}, func(rs *ResponseMessage) {
		cb(rs)
	})
}

func (ctl *BaseController) bind(c *gin.Context, bindTarget func() error, cb func(rs *ResponseMessage)) {
	rs := &ResponseMessage{Code: SucceedCode}

	if err := bindTarget(); err != nil {
		rs.Message = validator.Translate(err)
		ctl.ShowCustomError(c, rs)
		return
	}

	cb(rs)

	if rs.Message != nil {
		ctl.ShowCustomError(c, rs)
		return
	}

	if rs.Code > 0 {
		rs.Message = dict.Get(converter.ToStr(rs.Code))
		ctl.ShowCodeError(c, rs)
		return
	}

	rs.Message = SucceedMessage
	c.JSON(200, rs)
}

func (ctl *BaseController) ShowCodeError(c *gin.Context, rs *ResponseMessage) {
	rs.Message = dict.Get(converter.ToStr(rs.Code))
	c.JSON(rs.Code/10000, rs)
}

func (ctl *BaseController) ShowCustomError(c *gin.Context, rs *ResponseMessage) {
	rs.Code = 4000000
	c.JSON(400, rs)
}

func (ctl *BaseController) ShowSucceed(c *gin.Context, rs *ResponseMessage) {
	c.JSON(200, rs)
}
