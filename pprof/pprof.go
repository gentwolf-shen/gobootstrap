package pprof

import (
	"net/http/pprof"

	"github.com/gentwolf-shen/gin-boost"
)

func Register(engine *gin.Engine, prefix string) {
	if prefix == "" {
		prefix = "/debug/pprof"
	}

	p := engine.Group(prefix)
	{
		p.GET("/", gin.WrapF(pprof.Index))
		p.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		p.GET("/profile", gin.WrapF(pprof.Profile))
		p.POST("/symbol", gin.WrapF(pprof.Symbol))
		p.GET("/symbol", gin.WrapF(pprof.Symbol))
		p.GET("/trace", gin.WrapF(pprof.Trace))
		p.GET("/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))
		p.GET("/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
		p.GET("/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
		p.GET("/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
		p.GET("/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
		p.GET("/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
	}
}
