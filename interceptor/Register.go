package interceptor

import (
	"github.com/gentwolf-shen/gin-boost"
)

var (
	items []RegistryInfo
)

/**

URL 认证检测

interceptor.Registry("portal", func(c *gin.Context) bool {
	// 检测是否有权限访问, pattern为正则表达式
	// true / false
	return true
}).AddPathPatterns("(.*):/api/").
	ExcludePathPatterns("(.*):/api/passport")

*/
func Registry(name string, fun func(c *gin.Context) bool) *Interceptor {
	if len(items) == 0 {
		items = make([]RegistryInfo, 0)
	}

	item := RegistryInfo{name, NewInterceptor(), fun}
	items = append(items, item)

	return item.Target
}

func Valid(c *gin.Context) bool {
	bl := true
	path := []byte(c.Request.Method + ":" + c.Request.URL.Path)

	for _, item := range items {
		if item.Target.IsMustAuthorize(path) {
			bl = item.Fun(c)
			break
		}
	}

	return bl
}

type RegistryInfo struct {
	Name   string
	Target *Interceptor
	Fun    func(c *gin.Context) bool
}
