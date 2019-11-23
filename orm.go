package corm

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

/**
获取一个新的DB
conn 数据库连接
*/
func GetDb(conn *sql.DB) *Db {
	db := new(Db)
	db.conn = conn
	return db
}

/**
设置数据表
table 表名
*/
func (db *Db) Tab(table string) *Db {
	db.table = table
	return db
}

/**
设置查询字段，格式：Select("id", "name", "age")
field 查询字段
*/
func (db *Db) Select(field ...string) *Db {
	db.fields = append(db.fields, field...)
	return db
}

/**
设置查询字段原生格式，格式：Select("id, name, age, IFNULL(sex=1,1,2) AS sex")
field 查询字段
*/
func (db *Db) SelectRaw(field string) *Db {
	fieldTmp := strings.Split(field, ",")
	db.fields = append(db.fields, fieldTmp...)
	return db
}

/**
查询条件，格式：Where("id", ">", 100).where("name", "=", "张三")
field 查询字段
operator 条件符号 >、<、=、<>、like、in 等
condition 条件值
*/
func (db *Db) Where(field, operator string, condition interface{}) *Db {
	db.where = append(db.where, where{
		field:     field,
		operator:  operator,
		condition: condition,
	})
	return db
}

/**
将数字字符串转换成INT
*/
func (db *Db) WhereStrToInt(field, operator string, condition string) *Db {
	//strInt, err := strconv.ParseInt(condition, 10, 64)
	if condition == "" {
		db.pushErr(errors.New("func:WhereStrToInt condition is empty"))
	}
	strInt, err := strconv.Atoi(condition)
	if err != nil {
		db.pushErr(err)
	}
	return db.Where(field, operator, strInt)
}

/**
将int64转换成字符串
*/
func (db *Db) WhereInt64ToStr(field, operator string, condition int64) *Db {
	intStr := strconv.FormatInt(condition, 10)
	return db.Where(field, operator, intStr)
}

/**
将int转换成字符串
*/
func (db *Db) WhereIntToStr(field, operator string, condition int) *Db {
	intStr := strconv.Itoa(condition)
	return db.Where(field, operator, intStr)
}

/**
查询条件原生格式，格式：Where("id > 100 and name = '张三'")
where 条件字符串
*/
func (db *Db) WhereRaw(where string) *Db {
	db.whereRaw = append(db.whereRaw, where)
	return db
}

/**
查询 In 条件，格式：WhereIn("name", "张")
where 条件字符串
*/
func (db *Db) WhereIn(field string, condition ...interface{}) *Db {
	db.where = append(db.where, where{
		field:          field,
		operator:       IN,
		conditionArray: condition,
	})
	return db
}

/**
查询 Not In 条件，格式：WhereNotIn("name", "张")
where 条件字符串
*/
func (db *Db) WhereNotIn(field string, condition ...interface{}) *Db {
	db.where = append(db.where, where{
		field:          field,
		operator:       NOT_IN,
		conditionArray: condition,
	})
	return db
}

/**
查询 like 条件，格式：WhereLike("name", "张")
where 条件字符串
*/
func (db *Db) WhereLike(field string, condition string) *Db {
	condition = "%" + condition + "%"
	db.where = append(db.where, where{
		field:     field,
		operator:  LIKE,
		condition: condition,
	})
	return db
}

/**
查询 not like 条件，格式：WhereNotLike("name", "张")
where 条件字符串
*/
func (db *Db) WhereNotLike(field string, condition string) *Db {
	condition = "%" + condition + "%"
	db.where = append(db.where, where{
		field:     field,
		operator:  NOT_LIKE,
		condition: condition,
	})
	return db
}

/**
查询 Between 条件，格式：WhereBetween("id", 100, 1000)
where 条件字符串
*/
func (db *Db) WhereBetween(field string, startCondition interface{}, endCondition interface{}) *Db {
	db.where = append(db.where, where{
		field:          field,
		operator:       BETWEEN,
		conditionArray: []interface{}{startCondition, endCondition},
	})
	return db
}

/**
查询结果过滤 Having ，格式：Having("name", "=", "张三").Having("age", ">", 18)
Having 条件字符串
*/
func (db *Db) Having(field, operator string, condition interface{}) *Db {
	db.having = append(db.having, having{
		field:     field,
		operator:  operator,
		condition: condition,
	})
	return db
}

/**
排序，格式：OrderBy("id", "desc").OrderBy("name", "asc")
field 字段
by asc或desc
*/
func (db *Db) OrderBy(field, by string) *Db {
	db.orderBy = append(db.orderBy, orderBy{
		field: field,
		by:    by,
	})
	return db
}

/**
分组，格式：GroupBy("class", "type")
field 字段
by asc或desc
*/
func (db *Db) GroupBy(field ...string) *Db {
	db.groupBy = append(db.groupBy, field...)
	return db
}

/**
查询结果数量，格式：Limit(100)
limit 数量
*/
func (db *Db) Limit(limit int) *Db {
	db.limit = limit
	return db
}

/**
查询结果数量，格式：Limit(100)
limit 数量
*/
func (db *Db) Offset(offset int) *Db {
	db.offset = offset
	return db
}

/**
左连接，格式：LeftJoin("group", "user.group_id = group.id")
table 表名
on 关联条件
*/
func (db *Db) LeftJoin(table, on string) *Db {
	db.join = append(db.join, join{
		table:     table,
		direction: LEFT_JOIN,
		on:        on,
	})
	return db
}

/**
左连接，格式：RightJoin("group", "user.group_id = group.id")
table 表名
on 关联条件
*/
func (db *Db) RightJoin(table, on string) *Db {
	db.join = append(db.join, join{
		table:     table,
		direction: RIGHT_JOIN,
		on:        on,
	})
	return db
}

/**
左连接，格式：Join("group", "user.group_id = group.id")
table 表名
on 关联条件
*/
func (db *Db) Join(table, on string) *Db {
	db.join = append(db.join, join{
		table:     table,
		direction: INNER_JOIN,
		on:        on,
	})
	return db
}

/**
查询一条数据
callable 回调函数
*/
func (db *Db) First(result ...interface{}) error {
	err := db.queryRow(db.whereToSql(), db.getWhereValue(), result...)
	if errs(err) != nil {
		return err
	}
	return nil
}

/**
查询某个字段，以字符串形式返回
*/
func (db *Db) ValueStr(field string) (string, error) {
	db.fields = nil
	db.Select(field)

	var value string
	err := db.First(&value)
	return value, err
}

/**
查询某个字段，以数字形式返回
*/
func (db *Db) ValueInt(field string) (int, error) {
	db.fields = nil
	db.Select(field)

	var value int
	err := db.First(&value)
	return value, err
}

/**
查询某个字段，以Float形式返回
*/
func (db *Db) ValueFloat(field string) (float64, error) {
	db.fields = nil
	db.Select(field)

	var value float64
	err := db.First(&value)
	return value, err
}

/**
查询多条数据
callable 回调函数
*/
func (db *Db) Get(callable func(rows *sql.Rows)) error {
	rows, err := db.query(db.whereToSql(), db.getWhereValue()...)
	if errs(err) != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		callable(rows)
	}
	return nil
}

/**
分页数据查询
page 页数
pageCount 每页记录数
callable 回调函数
*/
func (db *Db) GetPage(page, pageCount int, callable func(rows *sql.Rows)) (int64, error) {
	db.offset = (page - 1) * pageCount
	db.limit = pageCount
	//var totalCount int64
	//var err error
	//总记录数
	totalCount, err := db.Count()
	if err != nil {
		return 0, err
	}
	err = db.Get(callable)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

/**
Sum
*/
func (db *Db) Sum(sumField string) (float64, error) {
	db.sum = sumField
	var sum sql.NullFloat64
	err := db.queryRow(db.sumToSql(), db.getWhereValue(), &sum)
	if errs(err) != nil {
		return 0, err
	}
	return sum.Float64, nil
}

/**
Sum
*/
func (db *Db) Max(maxField string) (int64, error) {
	db.max = maxField
	var max sql.NullInt64
	err := db.queryRow(db.maxToSql(), db.getWhereValue(), &max)
	if errs(err) != nil {
		return 0, err
	}
	return max.Int64, nil
}

/**
Sum
*/
func (db *Db) Min(minField string) (int64, error) {
	db.min = minField
	var min sql.NullInt64
	err := db.queryRow(db.minToSql(), db.getWhereValue(), &min)
	if errs(err) != nil {
		return 0, err
	}
	return min.Int64, nil
}

/**
Count
*/
func (db *Db) Count() (int64, error) {
	var count sql.NullInt64
	err := db.queryRow(db.countToSql(), db.getWhereValue(), &count)
	if errs(err) != nil {
		return 0, err
	}
	return count.Int64, nil
}

/**
Exists 查询数据是否存在
*/
func (db *Db) Exists() (bool, error) {
	count, err := db.Count()
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

/**
打印SQL
*/
func (db *Db) PrintSql() string {
	return db.whereToSql()
}

/**
插入数据
*/
func (db *Db) Insert(insertMap map[string]interface{}) (LastInsertId int64, err error) {
	db.insert = insertMap
	insertStr, vals := db.insertToSql()

	rest, err := db.exec(insertStr, vals...)
	if err != nil {
		return 0, err
	}
	insertId, err := rest.LastInsertId()
	if err != nil {
		return 0, err
	}
	return insertId, nil
}

/**
修改数据
*/
func (db *Db) Update(updateMap map[string]interface{}) (updateNum int64, err error) {
	db.update = updateMap
	updateStr, vals := db.updateToSql()

	vals = append(vals, db.getWhereValue()...)
	rest, err := db.exec(updateStr, vals...)
	if err != nil {
		return 0, err
	}
	rows, err := rest.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

//执行事务
func (db *Db) Transaction(callable func(dbTrans *Db) error) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}

	db.tx = tx

	err = callable(db)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
