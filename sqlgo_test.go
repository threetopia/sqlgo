package sqlgo

import (
	"testing"
)

const selectQuery string = "SELECT t.column_one AS columnOne, t.column_two AS columnTwo, t.column_three AS columnThree, t.column_no_alias FROM schema.table AS t INNER JOIN schema.join_table1 AS jt1 ON jt1.id=t.id INNER JOIN schema.join_table2 AS jt2 ON jt2.id=t.id WHERE t.column_one ILIKE ANY ($1) AND t.column_two=$2 AND t.column_three=$3 AND t.column_four IN ($4, $5) GROUP BY columnOne ORDER BY columnOne ASC, columnTwo DESC"

func TestSelectQueryPrepend(t *testing.T) {
	sql := NewSQLGo().
		SQLSchema("schema").
		SQLSelect(
			SetSQLSelect("t.column_one", "columnOne"),
			SetSQLSelect("t.column_two", "columnTwo"),
			SetSQLSelect("t.column_three", "columnThree"),
			SetSQLSelect("t.column_no_alias", ""),
		).
		SQLFrom("table", "t").
		SQLJoin(
			SetSQLJoin("INNER", "join_table1", "jt1", SetSQLJoinWhere("AND", "jt1.id", "=", "t.id")),
			SetSQLJoin("INNER", "join_table2", "jt2", SetSQLJoinWhere("AND", "jt2.id", "=", "t.id")),
		).
		SQLWhere(
			SetSQLWhere("AND", "t.column_one", "ILIKE ANY", []int{1, 2, 3}),
			SetSQLWhere("AND", "t.column_two", "=", float64(0)),
			SetSQLWhere("AND", "t.column_three", "=", 0),
			SetSQLWhere("AND", "t.column_four", "IN", []bool{true, false}),
		).
		SQLGroupBy("columnOne").
		SQLOrder(SetSQLOrder("columnOne", "ASC"), SetSQLOrder("columnTwo", "DESC"))
	if sqlStr := sql.BuildSQL(); sqlStr != selectQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", selectQuery, sqlStr)
	}
	t.Log(sql.GetSQLGoParameter().GetSQLParameter())
}

func TestSelectQueryPipeline(t *testing.T) {
	sql := NewSQLGo().
		SetSQLSchema("schema").
		SetSQLSelect("t.column_one", "columnOne").
		SetSQLSelect("t.column_two", "columnTwo").
		SetSQLSelect("t.column_three", "columnThree").
		SetSQLSelect("t.column_no_alias", "").
		SQLFrom("table", "t").
		SetSQLJoin("INNER", "join_table1", "jt1", SetSQLJoinWhere("AND", "jt1.id", "=", "t.id")).
		SetSQLJoin("INNER", "join_table2", "jt2", SetSQLJoinWhere("AND", "jt2.id", "=", "t.id")).
		SetSQLWhere("AND", "t.column_one", "ILIKE ANY", []int{1, 2, 3}).
		SetSQLWhere("AND", "t.column_two", "=", float64(0)).
		SetSQLWhere("AND", "t.column_three", "=", 0).
		SetSQLWhere("AND", "t.column_four", "IN", []bool{true, false}).
		SetSQLGroupBy("columnOne").
		SetSQLOrder("columnOne", "ASC").
		SetSQLOrder("columnTwo", "DESC")
	if sqlStr := sql.BuildSQL(); sqlStr != selectQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", selectQuery, sqlStr)
	}
	t.Log(sql.GetSQLGoParameter().GetSQLParameter())
}

func TestSelectQuerySetByVariable(t *testing.T) {
	sql := NewSQLGo()
	sqlSchema := NewSQLGoSchema().SetSQLSchema("schema")
	sqlSelect := NewSQLGoSelect().
		SetSQLSelect("t.column_one", "columnOne").
		SetSQLSelect("t.column_two", "columnTwo").
		SetSQLSelect("t.column_three", "columnThree").
		SetSQLSelect("t.column_no_alias", "")
	sqlFrom := NewSQLGoFrom().SQLFrom("table", "t")
	sqlJoin := NewSQLGoJoin().SetSQLJoin("INNER", "join_table1", "jt1", SetSQLJoinWhere("AND", "jt1.id", "=", "t.id")).
		SetSQLJoin("INNER", "join_table2", "jt2", SetSQLJoinWhere("AND", "jt2.id", "=", "t.id"))
	sqlWhere := NewSQLGoWhere().
		SetSQLWhere("AND", "t.column_one", "ILIKE ANY", []int{1, 2, 3}).
		SetSQLWhere("AND", "t.column_two", "=", float64(0)).
		SetSQLWhere("AND", "t.column_three", "=", 0).
		SetSQLWhere("AND", "t.column_four", "IN", []bool{true, false})
	sqlGroupBy := NewSQLGoGroupBy().SetSQLGroupBy("columnOne")
	sqlOrder := NewSQLGoOrder().SetSQLOrder("columnOne", "ASC").SetSQLOrder("columnTwo", "DESC")
	sql.SetSQLGoSchema(sqlSchema).
		SetSQLGoSelect(sqlSelect).
		SetSQLGoFrom(sqlFrom).
		SetSQLGoJoin(sqlJoin).
		SetSQLGoWhere(sqlWhere).
		SetSQLGoGroupBy(sqlGroupBy).
		SetSQLGoOrder(sqlOrder)
	if sqlStr := sql.BuildSQL(); sqlStr != selectQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", selectQuery, sqlStr)
	}
	t.Log(sql.GetSQLGoParameter().GetSQLParameter())
}

const deleteQuery string = "DELETE FROM table WHERE column_one=$1"

func TestDeleteQueryPrepend(t *testing.T) {
	sql := NewSQLGo().
		SQLDelete("table", SetSQLWhere("AND", "column_one", "=", "value_one"))
	if sqlStr := sql.BuildSQL(); sqlStr != deleteQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", deleteQuery, sqlStr)
	}
}

func TestDeleteQueryPipeline(t *testing.T) {
	sql := NewSQLGo().
		SetSQLDelete("table").
		SetSQLWhere("AND", "column_one", "=", "value_one")
	if sqlStr := sql.BuildSQL(); sqlStr != deleteQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", deleteQuery, sqlStr)
	}
}

const insertQuery string = "INSERT INTO table (col1, col2, col3) VALUES ($1, $2, $3), ($1, $2, $3), ($1, $2, $3)"

func TestInsert(t *testing.T) {
	sql := NewSQLGo().
		SetSQLInsert("table").
		SetSQLInsertColumn("col1", "col2", "col3").
		SetSQLInsertValue("val1", "val2", "val3").
		SetSQLInsertValue("val1", "val2", "val3").
		SetSQLInsertValue("val1", "val2", "val3")
	if sqlStr := sql.BuildSQL(); sqlStr != insertQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", insertQuery, sqlStr)
	}
}

const insertToTsVectorQuery string = "INSERT INTO table (col1, col2, col3) VALUES ($1, $2, $3), ($1, $2, $3), ($1, $2, to_tsvector('english', $3))"

func TestInsertToTsVector(t *testing.T) {
	sql := NewSQLGo().
		SetSQLInsert("table").
		SetSQLInsertColumn("col1", "col2", "col3").
		SetSQLInsertValue("val1", "val2", "val3").
		SetSQLInsertValue("val1", "val2", "val3").
		SetSQLInsertValue("val1", "val2", SetSQLInsertToTsVector("english", "val3"))
	if sqlStr := sql.BuildSQL(); sqlStr != insertToTsVectorQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", insertToTsVectorQuery, sqlStr)
	}
}

const updateQuery string = "UPDATE table SET col1=$1, col2=$2 WHERE col3=$3"

func TestUpdate(t *testing.T) {
	sql := NewSQLGo().
		SetSQLUpdate("table").
		SetSQLUpdateValue("col1", "val1").
		SetSQLUpdateValue("col2", "val2").
		SetSQLWhere("AND", "col3", "=", "val3")
	if sqlStr := sql.BuildSQL(); sqlStr != updateQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", updateQuery, sqlStr)
	}
}

const updateToTsVectorQuery string = "UPDATE schema.table SET col1=$1, val2=to_tsvector('english', $2) WHERE col3=$3"

func TestUpdateToTsVector(t *testing.T) {
	sql := NewSQLGo().
		SetSQLSchema("schema").
		SetSQLUpdate("table").
		SetSQLUpdateValue("col1", "val1").
		SetSQLUpdateToTsVector("val2", "english", "coalesce(title, '') || ' ' || coalesce(body, '')").
		SetSQLWhere("AND", "col3", "=", "val3")
	if sqlStr := sql.BuildSQL(); sqlStr != updateToTsVectorQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", updateQuery, sqlStr)
	}
}

func TestUpdateToTsVectorPrepend(t *testing.T) {
	sql := NewSQLGo().
		SQLSchema("schema").
		SQLUpdate("table",
			SetSQLUpdateValue("col1", "val1"),
			SetSQLUpdateToTsVector("val2", "english", "coalesce(title, '') || ' ' || coalesce(body, '')"),
		).
		SQLWhere(
			SetSQLWhere("AND", "col3", "=", "val3"),
		)
	if sqlStr := sql.BuildSQL(); sqlStr != updateToTsVectorQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", updateToTsVectorQuery, sqlStr)
	}
	t.Log(sql.GetSQLGoParameter().GetSQLParameter())
}

// func TestWhereINClause(t *testing.T) {
// 	sql := NewSQLGo().
// 		SetSQLSelect("t.col1", "col1").
// 		SetSQLSelect("t.col2", "col2").
// 		SetSQLFrom("table", "t").
// 		SetSQLWhere("AND", "asd", "IN", []string{"satu", "satu", "dua", "tiga", "empat", "satu"})
// 	fmt.Println(sql.BuildSQL(), sql.GetSQLGoParameter().GetSQLParameter())
// }

// func TestWhereAnyClause(t *testing.T) {
// 	sql := NewSQLGo().
// 		SetSQLSelect("t.col1", "col1").
// 		SetSQLSelect("t.col2", "col2").
// 		SetSQLFrom("table", "t").
// 		SetSQLWhere("AND", "asd", "ANY", []string{"satu", "satu", "dua", "tiga", "empat", "satu"})
// 	fmt.Println(sql.BuildSQL(), sql.GetSQLGoParameter().GetSQLParameter())
// }

// func TestOffsetLimit(t *testing.T) {
// 	sql := NewSQLGo().
// 		SetSQLSelect("t.column_one", "columnOne").
// 		SetSQLSelect("t.column_two", "columnTwo").
// 		SetSQLSelect("t.column_three", "columnThree").
// 		SetSQLSelect("t.column_no_alias", "").
// 		SQLFrom("table", "t").
// 		SQLPageLimit(1, 10)
// 	fmt.Println(sql.BuildSQL(), sql.GetSQLGoParameter().GetSQLParameter())
// }

// func TestALL(t *testing.T) {
// 	sql := NewSQLGo()
// 	sql.SetSQLSelect("u.id", "id")
// 	sql.SetSQLSelect("u.full_name", "full_name")
// 	sql.SetSQLSelect("u.id_card", "id_card")
// 	sql.SetSQLSelect("u.country_code", "country_code")
// 	sql.SetSQLSelect("u.search_meta", "search_meta")
// 	sql.SetSQLSelect("u.data_hash", "data_hash")
// 	sql.SetSQLSelect("u.deleted", "deleted")
// 	sql.SetSQLSelect("u.created_at", "created_at")
// 	sql.SetSQLSelect("u.updated_at", "updated_at")
// 	sql.SetSQLFrom(`"user"`, "u")
// 	sql.SetSQLJoin("INNER", "search_meta_view", "smv", SetSQLJoinWhere("AND", "smv.id", "=", "u.id"))
// 	sql.SetSQLWhere("AND", "u.id", "ANY", []string{"asd"})
// 	fmt.Println(sql.BuildSQL(), sql.GetSQLGoParameter().GetSQLParameter())
// }

const whereGroupQuery string = "SELECT t.column_one AS columnOne, t.column_two AS columnTwo, t.column_three AS columnThree, t.column_no_alias FROM schema.table AS t INNER JOIN schema.join_table1 AS jt1 ON jt1.id=t.id INNER JOIN schema.join_table2 AS jt2 ON jt2.id=t.id WHERE (t.column_two=$1 AND t.column_three=$2) AND t.column_one ILIKE ANY ($3)"

func TestWhereGroupQueryPrepend(t *testing.T) {
	sql := NewSQLGo().
		SQLSchema("schema").
		SQLSelect(
			SetSQLSelect("t.column_one", "columnOne"),
			SetSQLSelect("t.column_two", "columnTwo"),
			SetSQLSelect("t.column_three", "columnThree"),
			SetSQLSelect("t.column_no_alias", ""),
		).
		SQLFrom("table", "t").
		SQLJoin(
			SetSQLJoin("INNER", "join_table1", "jt1", SetSQLJoinWhere("AND", "jt1.id", "=", "t.id")),
			SetSQLJoin("INNER", "join_table2", "jt2", SetSQLJoinWhere("AND", "jt2.id", "=", "t.id")),
		).
		SQLWhereGroup("OR",
			SetSQLWhere("AND", "t.column_two", "=", "value_two"),
			SetSQLWhere("AND", "t.column_three", "=", 3),
		).
		SQLWhere(
			SetSQLWhere("AND", "t.column_one", "ILIKE ANY", []int{1, 2, 3}),
		)
	if sqlStr := sql.BuildSQL(); sqlStr != whereGroupQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", whereGroupQuery, sqlStr)
	}
	t.Log(sql.GetSQLGoParameter().GetSQLParameter())
}

func TestWhereGroupPipeline(t *testing.T) {
	sql := NewSQLGo().
		SetSQLSchema("schema").
		SetSQLSelect("t.column_one", "columnOne").
		SetSQLSelect("t.column_two", "columnTwo").
		SetSQLSelect("t.column_three", "columnThree").
		SetSQLSelect("t.column_no_alias", "").
		SQLFrom("table", "t").
		SetSQLJoin("INNER", "join_table1", "jt1", SetSQLJoinWhere("AND", "jt1.id", "=", "t.id")).
		SetSQLJoin("INNER", "join_table2", "jt2", SetSQLJoinWhere("AND", "jt2.id", "=", "t.id")).
		SetSQLWhereGroup("OR",
			SetSQLWhere("AND", "t.column_two", "=", "value_two"),
			SetSQLWhere("AND", "t.column_three", "=", 3),
		).
		SetSQLWhere("AND", "t.column_one", "ILIKE ANY", []int{1, 2, 3})
	if sqlStr := sql.BuildSQL(); sqlStr != whereGroupQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", whereGroupQuery, sqlStr)
	} else {
		t.Log(sqlStr)
	}
	t.Log(sql.GetSQLGoParameter().GetSQLParameter())
}

const whereBetween string = "SELECT t.column_one AS columnOne, t.column_two AS columnTwo, t.column_three AS columnThree, t.column_no_alias FROM schema.table AS t INNER JOIN schema.join_table1 AS jt1 ON jt1.id=t.id INNER JOIN schema.join_table2 AS jt2 ON jt2.id=t.id WHERE t.column_one ILIKE ANY ($1) OR (t.column_three=$2 AND (t.column_two BETWEEN $3 AND $4))"

func TestBetween(t *testing.T) {
	sql := NewSQLGo().
		SQLSchema("schema").
		SQLSelect(
			SetSQLSelect("t.column_one", "columnOne"),
			SetSQLSelect("t.column_two", "columnTwo"),
			SetSQLSelect("t.column_three", "columnThree"),
			SetSQLSelect("t.column_no_alias", ""),
		).
		SQLFrom("table", "t").
		SQLJoin(
			SetSQLJoin("INNER", "join_table1", "jt1", SetSQLJoinWhere("AND", "jt1.id", "=", "t.id")),
			SetSQLJoin("INNER", "join_table2", "jt2", SetSQLJoinWhere("AND", "jt2.id", "=", "t.id")),
		).
		SQLWhere(
			SetSQLWhere("AND", "t.column_one", "ILIKE ANY", []int{1, 2, 3}),
		).
		SQLWhereGroup("OR",
			SetSQLWhere("AND", "t.column_three", "=", 3),
			SetSQLWhereBetween("AND", "t.column_two", 100000, 2000000),
		)
	if sqlStr := sql.BuildSQL(); sqlStr != whereBetween {
		t.Errorf("result must be (%s) BuildSQL give (%s)", whereBetween, sqlStr)
	}
	t.Log(sql.GetSQLGoParameter().GetSQLParameter())
}

const whereTsQuery string = "SELECT ts_rank(tsv, to_tsquery('english', $1)) AS rank"

func TestToTsQuery(t *testing.T) {
	sql := NewSQLGo().
		SQLSchema("schema").
		SQLSelect(
			SetSQLSelectTsRank("tsv", "english", "open & source & software", "rank"),
		).
		SQLFrom("table", "t").
		SQLWhere(SetSQLWhereToTsQuery("AND", "tsv", "english", "open & source & software"))
	if sqlStr := sql.BuildSQL(); sqlStr != whereTsQuery {
		t.Errorf("result must be (%s) BuildSQL give (%s)", whereTsQuery, sqlStr)
	}
	t.Log(sql.GetSQLGoParameter().GetSQLParameter())
}
