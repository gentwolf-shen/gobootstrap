package db

import (
	"regexp"
	"sort"
	"strings"

	"gobootstrap/embed"
	"gobootstrap/logger"

	"github.com/gentwolf-shen/gohelper-v2/converter"
	"github.com/gentwolf-shen/gohelper-v2/timehelper"
)

var (
	ptn            = regexp.MustCompile(`V([0-9]+)__`)
	createTableSql = map[string]string{
		"mysql": `CREATE TABLE IF NOT EXISTS flyway_schema_history_go (
					id           BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT PRIMARY KEY,
					version      INT UNSIGNED     NOT NULL DEFAULT 0,
					script       VARCHAR(200)     NOT NULL DEFAULT 0,
					installed_on CHAR(20)         NOT NULL DEFAULT '',
					success      TINYINT UNSIGNED NOT NULL DEFAULT 0
				) ENGINE = InnoDB DEFAULT CHARSET = utf8`,
	}
)

func UseFlyway(dbName string, embedFile embed.ItfaceEmbedFile, prefix string) {
	obj := &Flyway{dbName: dbName}
	obj.Run(embedFile, prefix)
}

type Flyway struct {
	versions   map[int]FlywayEntity
	dbName     string
	daoService *BaseDaoService
}

func (s *Flyway) Run(mappers embed.ItfaceEmbedFile, prefix string) {
	s.daoService = NewBaseDaoService(s.dbName, "flyway_schema_history_go")
	s.listHistory()

	files, _ := mappers.ReadDir(prefix + "/" + s.dbName)
	size := len(files)
	if size == 0 {
		return
	}

	names := make(map[int]string, size)
	arr := make([]int, size)
	i := 0

	for _, file := range files {
		rs := ptn.FindStringSubmatch(file.Name())
		version := converter.ToInt(rs[1])
		_, ok := s.versions[version]
		if ok {
			continue
		}

		arr[i] = version
		names[version] = file.Name()
		i++
	}

	arr = arr[0:i]
	sort.Ints(arr)

	for _, v := range arr {
		logger.Sugar.Infof("Flyway SQL: %s", names[v])
		filename := prefix + "/" + s.dbName + "/" + names[v]
		b, err := mappers.ReadFile(filename)
		if err != nil || b == nil {
			logger.Sugar.Errorf("read [%s] error: %v", filename, err)
			continue
		}
		s.parseSql(v, names[v], b)
	}
}

func (s *Flyway) listHistory() {
	var rows []FlywayEntity
	err := s.daoService.List(&rows, &ListParam{Field: "id,version,script,installed_on,success"})
	if err != nil {
		logger.Sugar.Warnf("Flyway error %v", err)
		s.createTable()

		rows = []FlywayEntity{}
	}

	s.versions = make(map[int]FlywayEntity)
	for _, v := range rows {
		s.versions[v.Version] = v
	}
}

func (s *Flyway) createTable() {
	logger.Sugar.Info("Flyway init table")
	// TODO: use driver name
	_, err := s.daoService.CreateTable(createTableSql["mysql"])
	if err != nil {
		panic(err)
	}
}

func (s *Flyway) parseSql(version int, filename string, b []byte) {
	var success uint8 = 1
	segments := strings.Split(strings.TrimSpace(string(b)), ";")
	for _, segment := range segments {
		if segment == "" {
			continue
		}

		if _, err := s.daoService.CreateTable(segment); err != nil {
			logger.Sugar.Errorf("Flyway exec error %s, %v\n%s", filename, err, segment)
			success = 0
		}
	}

	p := &FlywayEntity{
		Version:     version,
		Script:      filename,
		InstalledOn: timehelper.Today(),
		Success:     success,
	}
	_, err := s.daoService.Insert(p)
	if err != nil {
		logger.Sugar.Error(err)
	}
}
