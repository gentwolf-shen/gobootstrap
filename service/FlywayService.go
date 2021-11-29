package service

import (
	"github.com/gentwolf-shen/gobootstrap"
	"github.com/gentwolf-shen/gobootstrap/dao"
	"github.com/gentwolf-shen/gobootstrap/embed"
	"github.com/gentwolf-shen/gohelper-v2/convert"
	"github.com/gentwolf-shen/gohelper-v2/logger"
	"github.com/gentwolf-shen/gohelper-v2/timehelper"
	"github.com/gentwolf-shen/gohelper-v2/util"
	"regexp"
	"sort"
	"strings"
)

var (
	Flyway = &FlywayService{}
)

type FlywayService struct {
	ptn      *regexp.Regexp
	dirRoot  string
	versions map[int]gobootstrap.FlywayEntity

	daos map[string]*dao.FlywayDao
	xml  string
}

func (s *FlywayService) Init(length int) {
	s.daos = make(map[string]*dao.FlywayDao, length)
	s.ptn = regexp.MustCompile(`V([0-9]+)__`)
}

func (s *FlywayService) AddDao(name string, dao *dao.FlywayDao) {
	s.daos[name] = dao
}

func (s *FlywayService) SetXml(str string) {
	s.xml = str
}

func (s *FlywayService) GetXml() string {
	if s.xml != "" {
		return s.xml
	}
	return `<?xml version="1.0" encoding="UTF-8"?>
<mapper version="1.0">
    <update id="Update">
        ${sql}
    </update>
    <select id="ListHistory">
        SELECT * FROM flyway_schema_history_go ORDER BY id DESC
    </select>
    <insert id="Insert">
        INSERT INTO flyway_schema_history_go
        SET version = #{version}, script = #{script},installed_on = #{installedOn}, success = #{success}
    </insert>
</mapper>`
}

func (s *FlywayService) Run(mappers embed.ItfaceEmbedFile, prefix, target string) {
	s.listHistory(target)

	files, _ := mappers.ReadDir(prefix + target)
	size := len(files)
	if size == 0 {
		return
	}

	names := make(map[int]string, size)
	arr := make([]int, size)
	i := 0

	for _, file := range files {
		rs := s.ptn.FindStringSubmatch(file.Name())
		version := convert.ToInt(rs[1])
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
		logger.Infof("SQL: %s", names[v])
		filename := prefix + target + "/" + names[v]
		b, err := mappers.ReadFile(filename)
		if err != nil || b == nil {
			logger.Errorf("read [%s] error: %v", filename, err)
			continue
		}
		s.parseSql(v, target, names[v], b)
	}
}

func (s *FlywayService) listHistory(target string) {
	var rows []gobootstrap.FlywayEntity
	err := s.daos[target].ListHistory(&rows)
	if err != nil {
		logger.Info("init flyway_schema_history")
		if err = s.daos[target].CreateTable(); err != nil {
			panic(err)
		}

		_ = s.daos[target].ListHistory(&rows)
	}

	s.versions = make(map[int]gobootstrap.FlywayEntity)
	for _, v := range rows {
		s.versions[v.Version] = v
	}
}

func (s *FlywayService) parseSql(version int, target, filename string, b []byte) {
	var success uint8 = 1
	segments := strings.Split(strings.TrimSpace(string(b)), ";")
	for _, segment := range segments {
		if segment == "" {
			continue
		}

		err := s.daos[target].Update(segment)
		if err != nil {
			logger.Error(filename)
			logger.Error(err)
			success = 0
		}
	}

	p := gobootstrap.FlywayEntity{
		Version:     version,
		Script:      filename,
		InstalledOn: timehelper.Today(),
		Success:     success,
	}
	s.daos[target].Insert(util.Struct2Map(p))
}
