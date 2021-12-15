package db

import (
	_ "embed"
	"github.com/gentwolf-shen/gobatis"
	"github.com/gentwolf-shen/gobootstrap/embed"
	"github.com/gentwolf-shen/gobootstrap/logger"
	"github.com/gentwolf-shen/gohelper-v2/config"
	"strings"
)

var (
	targets map[string]*gobatis.GoBatis

	//go:embed BaseMapper.xml
	baseMapper []byte
)

func GetGoBatis(dbName string) *gobatis.GoBatis {
	return targets[dbName]
}

func UseGoBatis(configs map[string]config.DbConfig, embedFile embed.ItfaceEmbedFile, prefix string) {
	targets = make(map[string]*gobatis.GoBatis, len(configs))

	prefix = strings.TrimRight(prefix, "/")
	logger.Sugar.Debug(prefix)
	for dbName, cfg := range configs {
		c := gobatis.DbConfig{
			Driver:             cfg.Type,
			Dsn:                cfg.Dsn,
			MaxOpenConnections: cfg.MaxOpenConnections,
			MaxIdleConnections: cfg.MaxIdleConnections,
			MaxLifeTime:        cfg.MaxLifeTime,
			MaxIdleTime:        cfg.MaxIdleTime,
		}
		targets[dbName] = gobatis.NewGoBatis(c)
		logger.Sugar.Debug("GoBatis db ", dbName)

		files, _ := embedFile.ReadDir(prefix + "/" + dbName)
		for _, n := range files {
			name := n.Name()
			logger.Sugar.Debug("GoBatis XML file ", name)
			b, err := embedFile.ReadFile(prefix + "/" + dbName + "/" + name)
			if err != nil {
				logger.Sugar.Error(err)
				continue
			}

			loadXml(dbName, strings.TrimSuffix(name, ".xml"), b)
		}

		loadXml(dbName, "BaseMapper", baseMapper)
	}
}

func loadXml(dbName, xmlName string, bytes []byte) {
	logger.Sugar.Debugf("GoBatis set xml[%s] -> db[%s]", xmlName, dbName)
	if err := targets[dbName].LoadFromBytes(xmlName, bytes); err != nil {
		logger.Sugar.Error(err)
	}
}
