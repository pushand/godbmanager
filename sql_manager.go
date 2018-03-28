package godbmanager

import (
	"database/sql"
	"fmt"
	"reflect"
)

//Query class that holds the transaction query string and params
type Query struct {
	sqlQuery                  string
	params                    []interface{}
	transactionIdOfOtherQuery int
	transactionIdOnExec       int64
}

//SqlManager class
type sqlManager struct {
	queries []Query
}

func (sqlManager sqlManager) getQuery(query string) {

}

func (sqlManager sqlManager) Insert(query string) {

}

func (sqlManager sqlManager) Update(query string) {

}

//Performs Row Query
func (sqlManager sqlManager) QueryRow(query string, args ...interface{}) *sql.Row {
	fmt.Println("query", query)
	fmt.Println("params", args)
	return Db.QueryRow(query, args...)
}

//Performs Row Query
func (sqlManager sqlManager) QueryRows(query string, args ...interface{}) (*sql.Rows, error) {
	fmt.Println("query", query)
	fmt.Println("params", args)
	rows, err := Db.Query(query, args...)
	//todo find way to write defer here as now the function who calls this method has to write defer rows.close()
	//defer func() { rows.Close() }()
	return rows, err
}

//Performs bulk transaction
//@param transactionIdOfOtherQuery - The transaction id of the previous query
//@params - here pass nil where you want transaction id of previous query to be replaced with this nil param
func (sqlManager *sqlManager) AddTransaction(query string, transactionIdOfOtherQuery int, params ...interface{}) int {

	q := Query{sqlQuery: query, params: params, transactionIdOfOtherQuery: transactionIdOfOtherQuery}

	sqlManager.queries = append(sqlManager.queries, q)
	id := len(sqlManager.queries)
	fmt.Println("total query added ", id)
	return id - 1
}

//Perform transaction commit. if failed will rollback
func (sqlManager sqlManager) PerformTransactions() bool {
	tx, err := Db.Begin()
	if err != nil {
		return false
	} else {
		fmt.Println("len of queries", len(sqlManager.queries))
		for index, query := range sqlManager.queries {
			fmt.Println(index, query.sqlQuery)
			fmt.Println("transactionIdOfOtherQuery", query.transactionIdOfOtherQuery)
			if query.transactionIdOfOtherQuery != -1 {
				for i, a := range query.params {
					fmt.Println(i, reflect.ValueOf(a).IsValid())
					if !reflect.ValueOf(a).IsValid() {
						query.params[i] = sqlManager.queries[query.transactionIdOfOtherQuery].transactionIdOnExec
						break
					}
				}
			}

			result, err := tx.Exec(query.sqlQuery, query.params...)
			if err != nil {
				fmt.Println("Rollback", err)
				tx.Rollback()
				break
			} else {
				sqlManager.queries[index].transactionIdOnExec, err = result.LastInsertId()
				if err != nil {
					break
				}
			}
		}
		err := tx.Commit()
		if err != nil {
			fmt.Println("Commit err", err)
			tx.Rollback()
			return false
		} else {
			fmt.Println("Commit success")
			return true
		}
	}
}

//return interface to perform sql query/transation
func GetSqlHandler() SqlHandler {
	return &sqlManager{}
}

//Interface helper methods
type SqlHandler interface {
	Insert(query string)
	Update(query string)
	QueryRow(query string, params ...interface{}) *sql.Row
	QueryRows(query string, params ...interface{}) (*sql.Rows, error)
	AddTransaction(query string, lastTransactionId int, params ...interface{}) int
	PerformTransactions() bool
}
