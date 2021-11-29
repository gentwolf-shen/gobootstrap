package controller

import (
	"github.com/gentwolf-shen/gin-boost"
	"github.com/gentwolf-shen/gobootstrap/entity"
	"github.com/gentwolf-shen/gobootstrap/validator"
	"github.com/gentwolf-shen/gohelper-v2/convert"
	"github.com/gentwolf-shen/gohelper-v2/dict"
)

type BaseController struct{}

func (ctl *BaseController) BindRequest(c *gin.Context, p interface{}, cb func(rs *entity.ResponseMessage)) {
	rs := &entity.ResponseMessage{}

	if err := c.BindRequest(p); err != nil {
		rs.Message = validator.Translate(err)
		ctl.ShowCustomError(c, rs)
		return
	}

	cb(rs)

	if rs.Code > 0 {
		rs.Message = dict.Get(convert.ToStr(rs.Code))
		ctl.ShowCodeError(c, rs)
		return
	}

	if rs.Message != nil {
		ctl.ShowCustomError(c, rs)
		return
	}

	if rs.Data == nil {
		rs.Data = "success"
	}

	c.JSON(200, rs)
}

func (ctl *BaseController) ShowCodeError(c *gin.Context, rs *entity.ResponseMessage) {
	rs.Message = dict.Get(convert.ToStr(rs.Code))
	c.JSON(rs.Code/10000, rs)
}

func (ctl *BaseController) ShowCustomError(c *gin.Context, rs *entity.ResponseMessage) {
	rs.Code = 4000000
	c.JSON(400, rs)
}
