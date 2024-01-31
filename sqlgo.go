package sqlgo

const (
	Dialect           sqlGoDialect = "dialect"
	MySQLDialect      sqlGoDialect = "mysql-" + Dialect
	PostgreSQLDialect sqlGoDialect = "postgresql-" + Dialect
)

type SQLGoBase interface {
	GetSQLGoParameter() SQLGoParameter
	BuildSQL() string
}

type SQLGOModuleSetter interface {
	SetSQLGoInsert(sqlGoInsert SQLGoInsert) SQLGo
	SetSQLGoUpdate(sqlGoUpdate SQLGoUpdate) SQLGo
	SetSQLGoDelete(sqlGoDelete SQLGoDelete) SQLGo
	SetSQLGoSelect(sqlGoSelect SQLGoSelect) SQLGo
	SetSQLGoValues(sqlGoValues SQLGoValues) SQLGo
	SetSQLGoFrom(sqlGoFrom SQLGoFrom) SQLGo
	SetSQLGoJoin(sqlGoJoin SQLGoJoin) SQLGo
	SetSQLGoWhere(sqlGoWhere SQLGoWhere) SQLGo
	SetSQLGoOrder(sqlGoOrder SQLGoOrder) SQLGo
	SetSQLGoOffsetLimit(sqlGoOffsetLimit SQLGoOffsetLimit) SQLGo
	SetSQLGoSchema(sqlGoSchema SQLGoSchema) SQLGo
}

type SQLGo interface {
	SQLGOModuleSetter

	SQLInsert(table sqlGoTable, columns sqlGoInsertColumnSlice, values ...sqlGoInsertValueSlice) SQLGo
	SetSQLInsert(table sqlGoTable) SQLGo
	SetSQLInsertColumn(columns ...sqlGoInsertColumn) SQLGo
	SetSQLInsertValue(values ...sqlGoInsertValue) SQLGo

	SQLUpdate(table sqlGoTable, values ...sqlGoUpdateValue) SQLGo
	SetSQLUpdate(table sqlGoTable) SQLGo
	SetSQLUpdateValue(column string, value interface{}) SQLGo
	SetSQLUpdateValueSlice(values ...sqlGoUpdateValue) SQLGo

	SQLDelete(table sqlGoTable, sqlWhere ...sqlGoWhereValue) SQLGo
	SetSQLDelete(table sqlGoTable) SQLGo

	SQLSelect(values ...sqlGoSelectValue) SQLGo
	SetSQLSelect(value interface{}, alias sqlGoAlias) SQLGo

	SQLValues(values ...sqlGoValuesValueSlice) SQLGo
	SetSQLValues(values ...sqlGoValuesValueSlice) SQLGo

	SQLFrom(table sqlGoTable, alias sqlGoAlias, columns ...sqlGoFromColumn) SQLGo
	SetSQLFrom(table sqlGoTable, alias sqlGoAlias) SQLGo
	SetSQLFromColumn(columns ...sqlGoFromColumn) SQLGo

	SQLJoin(values ...sqlGoJoinValue) SQLGo
	SetSQLJoin(joinType string, table sqlGoTable, alias sqlGoAlias, sqlWhere ...sqlGoWhereValue) SQLGo

	SQLWhere(values ...sqlGoWhereValue) SQLGo
	SetSQLWhere(whereType string, whereColumn string, operator string, value interface{}) SQLGo
	SQLWhereGroup(whereType string, values ...sqlGoWhereValue) SQLGo
	SetSQLWhereGroup(whereType string, values ...sqlGoWhereValue) SQLGo

	SQLOrder(values ...sqlGoOrderValue) SQLGo
	SetSQLOrder(value sqlGoValue, order string) SQLGo

	SQLOffsetLimit(offset int, limit int) SQLGo
	SetSQLLimit(limit int) SQLGo
	SetSQLOffset(offset int) SQLGo
	SQLPageLimit(page int, limit int) SQLGo
	SetSQLPage(page int) SQLGo

	SQLSchema(schema string) SQLGo
	SetSQLSchema(schema string) SQLGo

	SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGo
	SQLGoBase
}

type (
	sqlGo struct {
		sqlGoInsert      SQLGoInsert
		sqlGoUpdate      SQLGoUpdate
		sqlGoDelete      SQLGoDelete
		sqlGoSelect      SQLGoSelect
		sqlGoValues      SQLGoValues
		sqlGoFrom        SQLGoFrom
		sqlGoJoin        SQLGoJoin
		sqlGoWhere       SQLGoWhere
		sqlGoOrder       SQLGoOrder
		sqlGoOffsetLimit SQLGoOffsetLimit
		sqlGoSchema      SQLGoSchema
		sqlGoParameter   SQLGoParameter
	}

	sqlGoTable   interface{}
	sqlGoAlias   string
	sqlGoValue   interface{}
	sqlGoDialect string
)

func NewSQLGo() SQLGo {
	return &sqlGo{
		sqlGoInsert:      NewSQLGoInsert(),
		sqlGoUpdate:      NewSQLGoUpdate(),
		sqlGoDelete:      NewSQLGoDelete(),
		sqlGoSelect:      NewSQLGoSelect(),
		sqlGoValues:      NewSQLGoValues(),
		sqlGoFrom:        NewSQLGoFrom(),
		sqlGoJoin:        NewSQLGoJoin(),
		sqlGoWhere:       NewSQLGoWhere(),
		sqlGoOrder:       NewSQLGoOrder(),
		sqlGoOffsetLimit: NewSQLGoOffsetLimit(),
		sqlGoSchema:      NewSQLGoSchema(),
		sqlGoParameter:   NewSQLGoParameter(),
	}
}

func (s *sqlGo) SetSQLGoInsert(sqlGoInsert SQLGoInsert) SQLGo {
	s.sqlGoInsert = sqlGoInsert
	return s
}

func (s *sqlGo) SetSQLGoUpdate(sqlGoUpdate SQLGoUpdate) SQLGo {
	s.sqlGoUpdate = sqlGoUpdate
	return s
}

func (s *sqlGo) SetSQLGoDelete(sqlGoDelete SQLGoDelete) SQLGo {
	s.sqlGoDelete = sqlGoDelete
	return s
}

func (s *sqlGo) SetSQLGoSelect(sqlGoSelect SQLGoSelect) SQLGo {
	s.sqlGoSelect = sqlGoSelect
	return s
}

func (s *sqlGo) SetSQLGoValues(sqlGoValues SQLGoValues) SQLGo {
	s.sqlGoValues = sqlGoValues
	return s
}

func (s *sqlGo) SetSQLGoFrom(sqlGoFrom SQLGoFrom) SQLGo {
	s.sqlGoFrom = sqlGoFrom
	return s
}

func (s *sqlGo) SetSQLGoJoin(sqlGoJoin SQLGoJoin) SQLGo {
	s.sqlGoJoin = sqlGoJoin
	return s
}

func (s *sqlGo) SetSQLGoWhere(sqlGoWhere SQLGoWhere) SQLGo {
	s.sqlGoWhere = sqlGoWhere
	return s
}

func (s *sqlGo) SetSQLGoOrder(sqlGoOrder SQLGoOrder) SQLGo {
	s.sqlGoOrder = sqlGoOrder
	return s
}

func (s *sqlGo) SetSQLGoOffsetLimit(sqlGoOffsetLimit SQLGoOffsetLimit) SQLGo {
	s.sqlGoOffsetLimit = sqlGoOffsetLimit
	return s
}

func (s *sqlGo) SetSQLGoSchema(sqlGoSchema SQLGoSchema) SQLGo {
	s.sqlGoSchema = sqlGoSchema
	return s
}

func (s *sqlGo) GetSQLGoSchema() SQLGoSchema {
	return s.sqlGoSchema
}

func (s *sqlGo) SetSQLGoParameter(sqlGoParameter SQLGoParameter) SQLGo {
	s.sqlGoParameter = sqlGoParameter
	return s
}

func (s *sqlGo) GetSQLGoParameter() SQLGoParameter {
	return s.sqlGoParameter
}

func (s *sqlGo) SQLInsert(table sqlGoTable, columns sqlGoInsertColumnSlice, values ...sqlGoInsertValueSlice) SQLGo {
	s.sqlGoInsert.SQLInsert(table, columns, values...)
	return s
}

func (s *sqlGo) SetSQLInsert(table sqlGoTable) SQLGo {
	s.sqlGoInsert.SetSQLInsert(table)
	return s
}

func (s *sqlGo) SetSQLInsertColumn(columns ...sqlGoInsertColumn) SQLGo {
	s.sqlGoInsert.SetSQLInsertColumn(columns...)
	return s
}

func (s *sqlGo) SetSQLInsertValue(values ...sqlGoInsertValue) SQLGo {
	s.sqlGoInsert.SetSQLInsertValue(values...)
	return s
}

func (s *sqlGo) SQLUpdate(table sqlGoTable, values ...sqlGoUpdateValue) SQLGo {
	s.sqlGoUpdate.SQLUpdate(table, values...)
	return s
}

func (s *sqlGo) SetSQLUpdate(table sqlGoTable) SQLGo {
	s.sqlGoUpdate.SetSQLUpdate(table)
	return s
}

func (s *sqlGo) SetSQLUpdateValue(column string, value interface{}) SQLGo {
	s.sqlGoUpdate.SetSQLUpdateValue(column, value)
	return s
}

func (s *sqlGo) SetSQLUpdateValueSlice(values ...sqlGoUpdateValue) SQLGo {
	s.sqlGoUpdate.SetSQLUpdateValueSlice(values...)
	return s
}

func (s *sqlGo) SQLDelete(table sqlGoTable, sqlWhere ...sqlGoWhereValue) SQLGo {
	s.sqlGoDelete.SQLDelete(table)
	s.sqlGoWhere.SQLWhere(sqlWhere...)
	return s
}

func (s *sqlGo) SetSQLDelete(table sqlGoTable) SQLGo {
	s.sqlGoDelete.SQLDelete(table)
	return s
}

func (s *sqlGo) SQLSelect(values ...sqlGoSelectValue) SQLGo {
	s.sqlGoSelect.SQLSelect(values...)
	return s
}

func (s *sqlGo) SetSQLSelect(value interface{}, alias sqlGoAlias) SQLGo {
	s.sqlGoSelect.SetSQLSelect(value, alias)
	return s
}

func (s *sqlGo) SQLValues(values ...sqlGoValuesValueSlice) SQLGo {
	s.sqlGoValues.SQLValues(values...)
	return s
}

func (s *sqlGo) SetSQLValues(values ...sqlGoValuesValueSlice) SQLGo {
	s.sqlGoValues.SetSQLValues(values...)
	return s
}

func (s *sqlGo) SQLFrom(table sqlGoTable, alias sqlGoAlias, columns ...sqlGoFromColumn) SQLGo {
	s.sqlGoFrom.SQLFrom(table, alias, columns...)
	return s
}

func (s *sqlGo) SetSQLFrom(table sqlGoTable, alias sqlGoAlias) SQLGo {
	s.sqlGoFrom.SetSQLFrom(table, alias)
	return s
}

func (s *sqlGo) SetSQLFromColumn(columns ...sqlGoFromColumn) SQLGo {
	s.sqlGoFrom.SetSQLFromColumn(columns...)
	return s
}

func (s *sqlGo) SQLJoin(values ...sqlGoJoinValue) SQLGo {
	s.sqlGoJoin.SQLJoin(values...)
	return s
}

func (s *sqlGo) SetSQLJoin(joinType string, table sqlGoTable, alias sqlGoAlias, sqlWhere ...sqlGoWhereValue) SQLGo {
	s.sqlGoJoin.SetSQLJoin(joinType, table, alias, sqlWhere...)
	return s
}

func (s *sqlGo) SQLWhere(values ...sqlGoWhereValue) SQLGo {
	s.sqlGoWhere.SQLWhere(values...)
	return s
}

func (s *sqlGo) SetSQLWhere(whereType string, whereColumn string, operator string, value interface{}) SQLGo {
	s.sqlGoWhere.SetSQLWhere(whereType, whereColumn, operator, value)
	return s
}

func (s *sqlGo) SQLWhereGroup(whereType string, values ...sqlGoWhereValue) SQLGo {
	s.sqlGoWhere.SQLWhereGroup(whereType, values...)
	return s
}

func (s *sqlGo) SetSQLWhereGroup(whereType string, values ...sqlGoWhereValue) SQLGo {
	s.sqlGoWhere.SetSQLWhereGroup(whereType, values...)
	return s
}

func (s *sqlGo) SQLOrder(values ...sqlGoOrderValue) SQLGo {
	s.sqlGoOrder.SQLOrder(values...)
	return s
}

func (s *sqlGo) SetSQLOrder(value sqlGoValue, order string) SQLGo {
	s.sqlGoOrder.SetSQLOrder(value, order)
	return s
}

func (s *sqlGo) SQLOffsetLimit(offset int, limit int) SQLGo {
	s.sqlGoOffsetLimit.SQLOffsetLimit(offset, limit)
	return s
}

func (s *sqlGo) SetSQLLimit(limit int) SQLGo {
	s.sqlGoOffsetLimit.SetSQLLimit(limit)
	return s
}

func (s *sqlGo) SetSQLOffset(offset int) SQLGo {
	s.sqlGoOffsetLimit.SetSQLOffset(offset)
	return s
}

func (s *sqlGo) SQLPageLimit(page int, limit int) SQLGo {
	s.sqlGoOffsetLimit.SQLPageLimit(page, limit)
	return s
}

func (s *sqlGo) SetSQLPage(page int) SQLGo {
	s.sqlGoOffsetLimit.SetSQLPage(page)
	return s
}

func (s *sqlGo) SQLSchema(schema string) SQLGo {
	s.sqlGoSchema.SQLSchema(schema)
	return s
}

func (s *sqlGo) SetSQLSchema(schema string) SQLGo {
	s.sqlGoSchema.SetSQLSchema(schema)
	return s
}

func (s *sqlGo) BuildSQL() string {
	sql := ""
	s.sqlGoInsert.
		SetSQLGoParameter(s.GetSQLGoParameter()).
		SetSQLGoSchema(s.GetSQLGoSchema())
	sql = combineSQL(sql, s.sqlGoInsert.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoInsert.GetSQLGoParameter())

	s.sqlGoUpdate.
		SetSQLGoParameter(s.GetSQLGoParameter()).
		SetSQLGoSchema(s.GetSQLGoSchema())
	sql = combineSQL(sql, s.sqlGoUpdate.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoUpdate.GetSQLGoParameter())

	s.sqlGoDelete.
		SetSQLGoParameter(s.GetSQLGoParameter()).
		SetSQLGoSchema(s.GetSQLGoSchema())
	sql = combineSQL(sql, s.sqlGoDelete.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoDelete.GetSQLGoParameter())

	s.sqlGoSelect.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = combineSQL(sql, s.sqlGoSelect.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoSelect.GetSQLGoParameter())

	s.sqlGoFrom.
		SetSQLGoParameter(s.GetSQLGoParameter()).
		SetSQLGoSchema(s.GetSQLGoSchema())
	sql = combineSQL(sql, s.sqlGoFrom.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoFrom.GetSQLGoParameter())

	s.sqlGoJoin.
		SetSQLGoParameter(s.GetSQLGoParameter()).
		SetSQLGoSchema(s.GetSQLGoSchema())
	sql = combineSQL(sql, s.sqlGoJoin.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoJoin.GetSQLGoParameter())

	s.sqlGoWhere.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = combineSQL(sql, s.sqlGoWhere.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoWhere.GetSQLGoParameter())

	s.sqlGoValues.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = combineSQL(sql, s.sqlGoValues.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoValues.GetSQLGoParameter())

	s.sqlGoOrder.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = combineSQL(sql, s.sqlGoOrder.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoOrder.GetSQLGoParameter())

	s.sqlGoOffsetLimit.SetSQLGoParameter(s.GetSQLGoParameter())
	sql = combineSQL(sql, s.sqlGoOffsetLimit.BuildSQL())
	s.SetSQLGoParameter(s.sqlGoOffsetLimit.GetSQLGoParameter())

	return sql
}
