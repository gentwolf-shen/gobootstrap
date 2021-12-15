package gobootstrap

import (
	"database/sql"
	"github.com/gentwolf-shen/gin-boost"
	"github.com/gentwolf-shen/gobatis"
	"github.com/gentwolf-shen/gobootstrap/db"
	"github.com/gentwolf-shen/gobootstrap/embed"
	"github.com/gentwolf-shen/gobootstrap/interceptor"
	"github.com/gentwolf-shen/gobootstrap/logger"
	"github.com/gentwolf-shen/gohelper-v2/config"
	"github.com/gentwolf-shen/gohelper-v2/dict"
	"github.com/gentwolf-shen/gohelper-v2/endless"
	"runtime"
)

type Application struct {
	cfg           config.Config
	engine        *gin.Engine
	dbConnections map[string]*sql.DB
}

func New() *Application {
	app := &Application{}
	app.init()
	return app
}

func (this *Application) init() {
	var err error
	this.cfg, err = config.LoadDefault()
	if err != nil {
		panic("load default config error: " + err.Error())
	}

	dict.EnableEnv = true
	_ = dict.LoadDefault()

	if this.cfg.Web.IsDebug {
		this.engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		this.engine = gin.New()
	}

	this.engine.Use(this.auth())
	this.engine.Use(gin.Recovery())

	runtime.GOMAXPROCS(runtime.NumCPU())
}

func (this *Application) AllowCrossDomain(domains []string) *Application {
	if len(domains) == 1 && domains[0] == "*" {
		this.engine.Use(gin.AllowCrossDomainAll())
	} else {
		this.engine.Use(gin.AllowCrossDomain(domains))
	}
	return this
}

func (this *Application) Register(register func(app *Application)) *Application {
	register(this)
	return this
}

func (this *Application) Run() *Application {
	if err := endless.ListenAndServe(this.cfg.Web.Port, this.engine); err != nil {
		logger.Sugar.Error(err)
	}
	return this
}

func (this *Application) ShutdownHook(hook func()) {
	hook()
}

func (this *Application) UseGoBatis(embedFile embed.ItfaceEmbedFile, prefix string) *Application {
	gobatis.SetCustomLogger(logger.Sugar)
	db.UseGoBatis(this.cfg.Db, embedFile, prefix)
	return this
}

func (this *Application) UseFlyway(dbName string, embedFile embed.ItfaceEmbedFile, prefix string) *Application {
	db.UseFlyway(dbName, embedFile, prefix)
	return this
}

func (this *Application) GetWebEngine() *gin.Engine {
	return this.engine
}

func (this *Application) auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !interceptor.Valid(c) {
			c.Writer.WriteHeader(401)
			c.Abort()
		}
	}
}
