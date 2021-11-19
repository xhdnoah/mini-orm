package session

import (
	"mini-orm/clause"
	"reflect"
)

// Insert 将已存在对象的各字段值铺平
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	// 多次调用 Set 构造各个子句
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordValues...)
	// 调用 Build 构造最终 SQL 语句
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Find 根据已铺平字段的值构造出对象 session.Find(&[]User{})
func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()                                   // 获取切片元素类型
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable() // 根据类型映射表结构

	// 根据表结构构造 SELECT 语句，查询所有符合条件记录 rows
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	// 遍历每一行记录，利用反射创建 destType 实例, 将 dest 所有字段铺平构造切片 values
	for rows.Next() {
		dest := reflect.New(destType).Elem() // User{}
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface()) // [&User.Age, &User.Name]
		}
		// 将该行每一列的值依次赋值给 values 中的字段
		if err := rows.Scan(values...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
