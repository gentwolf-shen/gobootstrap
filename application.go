package gobootstrap

import (
	"database/sql"
	"embed"
	"github.com/gentwolf-shen/gin-boost"
	"github.com/gentwolf-shen/gobootstrap/helper"
	"github.com/gentwolf-shen/gobootstrap/interceptor"
	"github.com/gentwolf-shen/gohelper-v2/config"
	"github.com/gentwolf-shen/gohelper-v2/dict"
	"github.com/gentwolf-shen/gohelper-v2/endless"
	"github.com/gentwolf-shen/gohelper-v2/gomybatis"
	"github.com/gentwolf-shen/gohelper-v2/logger"
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

	logger.LoadDefault()

	if this.cfg.Web.IsDebug {
		this.engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		this.engine = gin.New()
	}

	this.engine.Use(helper.GinHelper.AllowCrossDomainAll())
	this.engine.Use(this.auth())
	this.engine.Use(gin.Recovery())

	runtime.GOMAXPROCS(runtime.NumCPU())
}

func (this *Application) Register(register func(app *Application)) *Application {
	register(this)
	return this
}

func (this *Application) Run() *Application {
	if err := endless.ListenAndServe(this.cfg.Web.Port, this.engine); err != nil {
		logger.Error(err)
	}
	return this
}

func (this *Application) ShutdownHook(hook func()) {
	this.closeDb()
	hook()

	this.closeDb()
}

func (this *Application) SetDbMapper(mappers embed.FS, prefix string) *Application {
	this.dbConnections = make(map[string]*sql.DB, len(this.cfg.Db))
	var err error

	for name, c := range this.cfg.Db {
		this.dbConnections[name], err = sql.Open(c.Type, c.Dsn)
		if err != nil {
			logger.Errorf("init database error %s %v", name, err)
			continue
		}

		this.dbConnections[name].SetMaxIdleConns(c.MaxIdleConnections)
		this.dbConnections[name].SetMaxOpenConns(c.MaxOpenConnections)

		dirs, err1 := mappers.ReadDir(prefix)
		if err1 != nil {
			logger.Errorf("read mapper error %s %v", name, err1)
			continue
		}

		for _, dir := range dirs {
			files, _ := mappers.ReadDir(prefix + "/" + dir.Name())
			for _, n := range files {
				info, _ := n.Info()
				b, _ := mappers.ReadFile(prefix + "/" + dir.Name() + "/" + info.Name())
				gomybatis.SetMapper(this.dbConnections[name], info.Name(), string(b))
			}
		}
	}

	return this
}

func (this *Application) closeDb() {
	for _, db := range this.dbConnections {
		_ = db.Close()
	}
}

func (this *Application) GetDb(name string) *sql.DB {
	return this.dbConnections[name]
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
