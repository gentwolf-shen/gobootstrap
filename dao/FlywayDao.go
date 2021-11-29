package dao

import (
	"github.com/gentwolf-shen/gohelper-v2/gomybatis"
)

type FlywayDao struct {
	prefix string
}

func NewFlyway(target string) *FlywayDao {
	return &FlywayDao{prefix: target}
}

func (d *FlywayDao) fixName(name string) string {
	return d.prefix + ":Flyway." + name
}

func (d *FlywayDao) ListHistory(value interface{}) error {
	return gomybatis.QueryObjects(value, d.fixName("ListHistory"), map[string]interface{}{})
}

func (d *FlywayDao) CreateTable() error {
	sql := `
		CREATE TABLE IF NOT EXISTS flyway_schema_history_go
		(
			id           BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT PRIMARY KEY,
			version      INT UNSIGNED     NOT NULL DEFAULT 0,
			script       VARCHAR(200)     NOT NULL DEFAULT 0,
			installed_on CHAR(20)         NOT NULL DEFAULT '',
			success      TINYINT UNSIGNED NOT NULL DEFAULT 0
		) ENGINE = InnoDB
		  DEFAULT CHARSET = utf8`
	return d.Update(sql)
}

func (d *FlywayDao) Update(sql string) error {
	_, err := gomybatis.Update(d.fixName("Update"), map[string]interface{}{"sql": sql})
	return err
}

func (d *FlywayDao) Insert(p map[string]interface{}) bool {
	_, err := gomybatis.Insert(d.fixName("Insert"), p)
	return err != nil
}
