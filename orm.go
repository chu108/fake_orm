package corm

import (
	"database/sql"
	"strings"
)

/**
获取一个新的DB
conn 数据库连接
*/
func GetDb(conn *sql.DB) *db {
	db := new(db)
	db.conn = conn
	return db
}

/**
设置数据表
table 表名
*/
func (db *db) Tab(table string) *db {
	db.table = table
	return db
}

/**
设置查询字段，格式：Select("id", "name", "age")
field 查询字段
*/
func (db *db) Select(field ...string) *db {
	db.fields = append(db.fields, field...)
	return db
}

/**
设置查询字段原生格式，格式：Select("id, name, age, IFNULL(sex=1,1,2) AS sex")
field 查询字段
*/
func (db *db) SelectRaw(field string) *db {
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
func (db *db) Where(field, operator string, condition interface{}) *db {
	db.where = append(db.where, where{
		field:     field,
		operator:  operator,
		condition: condition,
	})
	return db
}

/**
查询条件原生格式，格式：Where("id > 100 and name = '张三'")
where 条件字符串
*/
func (db *db) WhereRaw(where string) *db {
	db.whereRaw = append(db.whereRaw, where)
	return db
}

/**
查询 In 条件，格式：WhereIn("name", "张")
where 条件字符串
*/
func (db *db) WhereIn(field string, condition ...interface{}) *db {
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
func (db *db) WhereNotIn(field string, condition ...interface{}) *db {
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
func (db *db) WhereLike(field string, condition interface{}) *db {
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
func (db *db) WhereNotLike(field string, condition interface{}) *db {
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
func (db *db) WhereBetween(field string, startCondition interface{}, endCondition interface{}) *db {
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
func (db *db) Having(field, operator string, condition interface{}) *db {
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
func (db *db) OrderBy(field, by string) *db {
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
func (db *db) GroupBy(field ...string) *db {
	db.groupBy = append(db.groupBy, field...)
	return db
}

/**
查询结果数量，格式：Limit(100)
limit 数量
*/
func (db *db) Limit(limit int64) *db {
	db.limit = limit
	return db
}

/**
左连接，格式：LeftJoin("group", "user.group_id = group.id")
table 表名
on 关联条件
*/
func (db *db) LeftJoin(table, on string) *db {
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
func (db *db) RightJoin(table, on string) *db {
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
func (db *db) Join(table, on string) *db {
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
func (db *db) First(result ...interface{}) error {
	where := make([]interface{}, 0, 5)
	if len(db.where) > 0 {
		for _, w := range db.where {
			where = append(where, w.condition)
		}
	}
	err := db.conn.QueryRow(db.whereToSql(), where...).Scan(result...)
	if errs(err) != nil {
		return err
	}
	return nil
}

/**
查询多条数据
callable 回调函数
*/
func (db *db) Get(callable func(rows *sql.Rows)) error {
	where := make([]interface{}, 0, 5)
	if len(db.where) > 0 {
		for _, w := range db.where {
			where = append(where, w.condition)
		}
	}
	rows, err := db.conn.Query(db.whereToSql(), where...)
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
Sum
*/
func (db *db) Sum(sumField string) (float64, error) {
	db.sum = sumField
	var sum sql.NullFloat64
	where := make([]interface{}, 0, 5)
	if len(db.where) > 0 {
		for _, w := range db.where {
			where = append(where, w.condition)
		}
	}
	err := db.conn.QueryRow(db.sumToSql(), where...).Scan(&sum)
	if errs(err) != nil {
		return 0, err
	}
	return sum.Float64, nil
}

/**
Sum
*/
func (db *db) Max(maxField string) (int64, error) {
	db.max = maxField
	var max sql.NullInt64
	where := make([]interface{}, 0, 5)
	if len(db.where) > 0 {
		for _, w := range db.where {
			where = append(where, w.condition)
		}
	}
	err := db.conn.QueryRow(db.maxToSql(), where...).Scan(&max)
	if errs(err) != nil {
		return 0, err
	}
	return max.Int64, nil
}

/**
Sum
*/
func (db *db) Min(minField string) (int64, error) {
	db.min = minField
	var min sql.NullInt64
	where := make([]interface{}, 0, 5)
	if len(db.where) > 0 {
		for _, w := range db.where {
			where = append(where, w.condition)
		}
	}
	err := db.conn.QueryRow(db.minToSql(), where...).Scan(&min)
	if errs(err) != nil {
		return 0, err
	}
	return min.Int64, nil
}

/**
Count
*/
func (db *db) Count() (int64, error) {
	var count sql.NullInt64
	where := make([]interface{}, 0, 5)
	if len(db.where) > 0 {
		for _, w := range db.where {
			where = append(where, w.condition)
		}
	}
	err := db.conn.QueryRow(db.countToSql(), where...).Scan(&count)
	if errs(err) != nil {
		return 0, err
	}
	return count.Int64, nil
}

/**
Exists 查询数据是否存在
*/
func (db *db) Exists() (bool, error) {
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
func (db *db) PrintSql() string {
	return db.whereToSql()
}

/**
插入数据
*/
func (db *db) Insert(insertMap map[string]interface{}) (LastInsertId int64, err error) {
	db.insert = insertMap
	insertStr, vals := db.insertToSql()

	stmt, err := db.conn.Prepare(insertStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	rest, err := stmt.Exec(vals...)
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
func (db *db) Update(updateMap map[string]interface{}) (updateNum int64, err error) {
	db.update = updateMap
	updateStr, vals := db.updateToSql()

	where := make([]interface{}, 0, 5)
	if len(db.where) > 0 {
		for _, w := range db.where {
			if len(w.conditionArray) > 0 {
				where = append(where, w.conditionArray...)
			}
			if w.condition != nil {
				where = append(where, w.condition)
			}
		}
	}
	vals = append(vals, where...)

	stmt, err := db.conn.Prepare(updateStr)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	rest, err := stmt.Exec(vals...)
	if err != nil {
		return 0, err
	}
	rows, err := rest.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

func errs(err error) error {
	if err != nil {
		//数据库记录不存在，报此错，可忽略
		if err == sql.ErrNoRows {
			return nil
		}
		//字段值为NULL时，报此错，可忽略，默认为类型的零值
		if strings.Index(err.Error(), "sql: Scan error on column index") != -1 {
			return nil
		}
	}
	return err
}
