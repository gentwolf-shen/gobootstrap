package gobootstrap

import (
	"runtime"

	"github.com/gentwolf-shen/gin-boost"
	"github.com/gentwolf-shen/gobatis"
	"github.com/gentwolf-shen/gobootstrap/db"
	"github.com/gentwolf-shen/gobootstrap/embed"
	"github.com/gentwolf-shen/gobootstrap/interceptor"
	"github.com/gentwolf-shen/gobootstrap/logger"
	"github.com/gentwolf-shen/gohelper-v2/config"
	"github.com/gentwolf-shen/gohelper-v2/dict"
	"github.com/gentwolf-shen/gohelper-v2/endless"
)

type Application struct {
	cfg    config.Config
	engine *gin.Engine
}

func New() *Application {
	app := &Application{}
	app.init()
	return app
}

func (a *Application) init() {
	var err error
	a.cfg, err = config.LoadDefault()
	if err != nil {
		panic("load default config error: " + err.Error())
	}

	logger.LoadDefault()

	dict.EnableEnv = true
	_ = dict.LoadDefault()

	if a.cfg.Web.IsDebug {
		a.engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		a.engine = gin.New()
	}

	a.engine.Use(a.auth())
	a.engine.Use(gin.Recovery())

	runtime.GOMAXPROCS(runtime.NumCPU())
}

func (a *Application) AllowCrossDomain(domains []string) *Application {
	if len(domains) == 1 && domains[0] == "*" {
		a.engine.Use(gin.AllowCrossDomainAll())
	} else {
		a.engine.Use(gin.AllowCrossDomain(domains))
	}
	return a
}

func (a *Application) Register(register func(app *Application)) *Application {
	register(a)
	return a
}

func (a *Application) Run() *Application {
	if err := endless.ListenAndServe(a.cfg.Web.Port, a.engine); err != nil {
		logger.Sugar.Error(err)
	}
	return a
}

func (a *Application) ShutdownHook(hook func()) {
	hook()
}

func (a *Application) UseGoBatis(embedFile embed.ItfaceEmbedFile, prefix string) *Application {
	gobatis.SetCustomLogger(logger.Sugar)
	db.UseGoBatis(a.cfg.Db, embedFile, prefix)
	return a
}

func (a *Application) UseFlyway(dbName string, embedFile embed.ItfaceEmbedFile, prefix string) *Application {
	db.UseFlyway(dbName, embedFile, prefix)
	return a
}

func (a *Application) GetWebEngine() *gin.Engine {
	return a.engine
}

func (a *Application) auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !interceptor.Valid(c) {
			c.Abort()
		}
	}
}
