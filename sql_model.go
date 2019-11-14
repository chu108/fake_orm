package corm

import "database/sql"

type where struct {
	field          string
	operator       string
	condition      interface{}
	conditionArray []interface{}
}

type having struct {
	field     string
	operator  string
	condition interface{}
}

type join struct {
	table     string
	direction string
	on        string
}

type orderBy struct {
	field string
	by    string
}

type db struct {
	conn      *sql.DB
	table     string
	join      []join
	fields    []string
	where     []where
	whereRaw  []string
	orderBy   []orderBy
	groupBy   []string
	limit     int64
	having    []having
	sum       string
	count     string
	max       string
	min       string
	insert    map[string]interface{}
	update    map[string]interface{}
	chunkById int64
	err       error
}
