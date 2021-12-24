package db

import (
	"github.com/gentwolf-shen/gobootstrap/logger"
	"github.com/gentwolf-shen/gobootstrap/util"
)

var (
	mapperName = "BaseMapper"
)

// 数据操作，方便直接使用struct或map[string]interface{}访问数据库
type BaseDaoService struct {
	dbName    string
	tableName string
}

func NewBaseDaoService(dbName, tableName string) *BaseDaoService {
	return &BaseDaoService{
		dbName:    dbName,
		tableName: tableName,
	}
}

func (s *BaseDaoService) getSelector(name string) string {
	return mapperName + "." + name
}

// 分页查询
// value: 记录struct
// p: 查询及分页参数
func (s *BaseDaoService) Page(value interface{}, p *ListParam) *PaginationEntity {
	rs := &PaginationEntity{
		Page:      p.Page,
		Size:      p.Size,
		Count:     0,
		TotalPage: 0,
		Items:     []interface{}{},
	}

	var err error
	rs.Count, err = s.Count(p.Where)
	if err != nil {
		logger.Sugar.Error(err)
	}

	if rs.Count == 0 {
		return rs
	}

	rs.TotalPage = util.Ceil(rs.Size, rs.Count)
	if p.Page > rs.TotalPage {
		return rs
	}

	if err = s.List(value, p); err != nil {
		logger.Sugar.Error(err)
		return rs
	}

	rs.Items = value

	return rs
}

// 简单查询，返回数据不分页
// value: 记录struct，查询字段为value的tag为db
// p: 查询及分页参数
func (s *BaseDaoService) List(value interface{}, p *ListParam) error {
	if p.Field == "" {
		p.Field = util.QueryDbTagField(value)
	}

	inputValue := map[string]interface{}{
		"table":       s.tableName,
		"field":       p.Field,
		"whereValues": p.Where,
		"page":        p.Page,
		"size":        p.Size,
		"orderBy":     p.OrderBy,
	}
	if p.Page > 0 && p.Size > 0 {
		inputValue["offset"] = util.ToOffset(p.Page, p.Size)
		inputValue["size"] = p.Size
	}

	return GetGoBatis(s.dbName).QueryObjects(value, s.getSelector("List"), inputValue)
}

// 查询记录数
func (s *BaseDaoService) Count(p map[string]interface{}) (int64, error) {
	inputValue := map[string]interface{}{
		"table":       s.tableName,
		"whereValues": p,
	}
	var c int64
	err := GetGoBatis(s.dbName).QueryObject(&c, s.getSelector("Count"), inputValue)
	return c, err
}

func (s *BaseDaoService) QueryById(value interface{}, field string, id int64) error {
	return s.Query(value, field, map[string]interface{}{"id": id})
}

// 查询单条记录
// value: 记录struct
// field: 查询的字段
// p: 查询的条件，如果是struct，则转换为map[string]interface{}
func (s *BaseDaoService) Query(value interface{}, field string, p interface{}) error {
	if field == "" {
		field = util.QueryDbTagField(value)
	}
	inputValue := map[string]interface{}{
		"table":       s.tableName,
		"field":       field,
		"whereValues": util.ToMap(p),
	}
	return GetGoBatis(s.dbName).QueryObject(value, s.getSelector("Query"), inputValue)
}

// 插入记录
// p: map[string]interface{}，如是是struct，则查询tag为db,或且标记有insert的key
func (s *BaseDaoService) Insert(p interface{}) (int64, error) {
	keys, values := util.ToArray(util.QueryDbTagMap(p, "insert"))
	inputValue := map[string]interface{}{
		"table":  s.tableName,
		"keys":   keys,
		"values": values,
	}
	return GetGoBatis(s.dbName).Insert(s.getSelector("Insert"), inputValue)
}

func (s *BaseDaoService) UpdateById(p interface{}, id int64) (int64, error) {
	return s.Update(p, map[string]interface{}{"id": id})
}

// 更新记录
// p: map[string]interface{}，如是是struct，则查询tag为db,或且标记有update的key
// argsWhere: 更新的条件
func (s *BaseDaoService) Update(p interface{}, argsWhere map[string]interface{}) (int64, error) {
	inputValue := map[string]interface{}{
		"table":        s.tableName,
		"updateValues": util.QueryDbTagMap(p, "update"),
		"whereValues":  argsWhere,
	}
	return GetGoBatis(s.dbName).Update(s.getSelector("Update"), inputValue)
}

func (s *BaseDaoService) DeleteById(id int64) (int64, error) {
	return s.Delete(map[string]interface{}{"id": id})
}

// 删除记录
func (s *BaseDaoService) Delete(p map[string]interface{}) (int64, error) {
	inputValue := map[string]interface{}{
		"table":       s.tableName,
		"whereValues": p,
	}
	return GetGoBatis(s.dbName).Update(s.getSelector("Delete"), inputValue)
}

// 创建表
// sql: SQL的DDL
func (s *BaseDaoService) CreateTable(sql string) (int64, error) {
	return GetGoBatis(s.dbName).Update(s.getSelector("CreateTable"), map[string]interface{}{"sql": sql})
}
