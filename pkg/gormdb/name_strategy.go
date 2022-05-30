/*
@Date: 2022/1/10 16:55
@Author: yvanz
@File : name_strategy
*/

package gormdb

import (
	"gorm.io/gorm/schema"
)

// MyNamingStrategy 只改了ColumnName，直接返回结构体的字段名(用于大驼峰标准)
type MyNamingStrategy struct {
	ns schema.NamingStrategy
}

func (mns MyNamingStrategy) TableName(str string) string {
	return mns.ns.TableName(str)
}

func (mns MyNamingStrategy) SchemaName(table string) string {
	return mns.ns.SchemaName(table)
}

// ColumnName convert string to column name
func (mns MyNamingStrategy) ColumnName(table, column string) string {
	_ = table
	return column
}

// JoinTableName convert string to join table name
func (mns MyNamingStrategy) JoinTableName(str string) string {
	return mns.ns.JoinTableName(str)
}

// RelationshipFKName generate fk name for relation
func (mns MyNamingStrategy) RelationshipFKName(rel schema.Relationship) string {
	return mns.ns.RelationshipFKName(rel)
}

// CheckerName generate checker name
func (mns MyNamingStrategy) CheckerName(table, column string) string {
	return mns.ns.CheckerName(table, column)
}

// IndexName generate index name
func (mns MyNamingStrategy) IndexName(table, column string) string {
	return mns.ns.IndexName(table, column)
}
