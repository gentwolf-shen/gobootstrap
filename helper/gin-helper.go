package helper

import (
	"github.com/gentwolf-shen/gin-boost"
)

var (
	GinHelper = &ginHelper{}
)

type ginHelper struct{}

func (this *ginHelper) AllowCrossDomainAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		this.allowCrossDomain(c, "*")
	}
}

func (this *ginHelper) AllowCrossDomain(domains []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Header.Get("Origin")
		bl := false
		for _, domain := range domains {
			if domain == host {
				bl = true
				break
			}
		}

		if bl {
			this.allowCrossDomain(c, c.Request.Header.Get("Origin"))
		}
	}
}

func (this *ginHelper) allowCrossDomain(c *gin.Context, host string) {
	c.Header("Access-Control-Allow-Origin", host)
	c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, X-Requested-With, Content-Type, Authorization")
	c.Header("Access-Control-Allow-Credentials", "true")
	if c.Request.Method == "OPTIONS" {
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATH,DELETE,OPTIONS,HEAD")
		c.Header("Access-Control-Max-Age", "3600")
		c.AbortWithStatus(200)
	}
}
