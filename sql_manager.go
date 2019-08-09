package godbmanager

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
)

const GO_DB_DEBUG = false

//Query class that holds the transaction query string and params
type Query struct {
	params                    []interface{}
	transactionIdOfOtherQuery int
	transactionIdOnExec       int64
}

//SqlManager class
type sqlManager struct {
	queries  []Query
	sqlQuery string
}

func (sqlManager sqlManager) getQuery(query string) {

}

//returns transactionId, rowsAffected and err if insert update failed/not happened
func (sqlManager sqlManager) Insert(query string, args ...interface{}) (int64, int64, error) {
	stmt, err := Db.Prepare(query)
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
		return 0, 0, err
	}
	res, err := stmt.Exec(args...)
	if err != nil {
		fmt.Println(err)
		return 0, 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0, 0, err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return 0, 0, err
	}
	//log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
	return lastId, rowCnt, nil
}

func (sqlManager sqlManager) Update(query string) {

}

//Performs Row Query
func (sqlManager sqlManager) QueryRow(query string, args ...interface{}) *sql.Row {
	if GO_DB_DEBUG {
		fmt.Println("query", query)
		fmt.Println("params", args)
	}

	stmt, err := Db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	return stmt.QueryRow(args...)
}

//Performs Row Query
func (sqlManager sqlManager) QueryRows(query string, args ...interface{}) (*sql.Rows, error) {
	if GO_DB_DEBUG {
		fmt.Println("query", query)
		fmt.Println("params", args)
	}

	stmt, err := Db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)

	//todo find way to write defer here as now the function who calls this method has to write defer rows.close()
	//defer func() { rows.Close() }()
	return rows, err
}

func (sqlManager *sqlManager) AddTransactionQuery(query string) {
	sqlManager.queries = make([]Query, 0)
	sqlManager.sqlQuery = query
}

//Performs bulk transaction
//@param transactionIdOfOtherQuery - The transaction id of the previous query
//@params - here pass nil where you want transaction id of previous query to be replaced with this nil param
func (sqlManager *sqlManager) AddTransactions(transactionIdOfOtherQuery int, params ...interface{}) int {
	q := Query{params: params, transactionIdOfOtherQuery: transactionIdOfOtherQuery}
	sqlManager.queries = append(sqlManager.queries, q)
	id := len(sqlManager.queries)
	if GO_DB_DEBUG {
		fmt.Println("total query added ", id)
	}
	return id - 1
}

//Perform transaction commit. if failed will rollback
func (sqlManager sqlManager) PerformTransactions() error {
	totalQueries := len(sqlManager.queries)
	if totalQueries == 0 {
		return nil
	}
	tx, err := Db.Begin()
	if err != nil {
		return err
	} else {
		fmt.Println("len of queries", len(sqlManager.queries))
		defer func() {
			// Rollback the transaction after the function returns.
			// If the transaction was already commited, this will do nothing.
			_ = tx.Rollback()
		}()
		stmt, err := tx.Prepare(sqlManager.sqlQuery)
		defer stmt.Close()
		if err != nil {
			fmt.Println("stmt", err)
			return err
		}
		for index, query := range sqlManager.queries {
			//fmt.Println(index, query.sqlQuery)
			//fmt.Println("transactionIdOfOtherQuery", query.transactionIdOfOtherQuery)
			if query.transactionIdOfOtherQuery != -1 {
				for i, a := range query.params {
					//fmt.Println(i, reflect.ValueOf(a).IsValid())
					if !reflect.ValueOf(a).IsValid() {
						query.params[i] = sqlManager.queries[query.transactionIdOfOtherQuery].transactionIdOnExec
						break
					}
				}
			}

			result, err := stmt.Exec(query.params...)
			if err != nil {
				fmt.Println("Rollback", err)
				return err
			} else {
				sqlManager.queries[index].transactionIdOnExec, err = result.LastInsertId()
				if err != nil {
					return err
				}
			}
		}
		errCommit := tx.Commit()
		if errCommit != nil {
			fmt.Println("Commit err", errCommit)
			return err
		} else {
			fmt.Println("Commit success")
			return nil
		}
	}
}

func (sqlManager sqlManager) PerformMultiTransactions(queries []string) error {
	totalQueries := len(queries)
	if totalQueries == 0 {
		return nil
	}
	tx, err := Db.Begin()
	if err != nil {
		return err
	} else {
		//fmt.Println("len of queries", len(queries))
		defer func() {
			// Rollback the transaction after the function returns.
			// If the transaction was already commited, this will do nothing.
			_ = tx.Rollback()
		}()
		for _, q := range queries {
			//fmt.Println("q", q)
			_, err := tx.Exec(q)
			if err != nil {
				fmt.Println("stmt", err)
				return err
			}
		}
		errCommit := tx.Commit()
		if errCommit != nil {
			fmt.Println("Commit err", errCommit)
			return err
		} else {
			fmt.Println("Commit success")
			return nil
		}
	}
}

//return interface to perform sql query/transation
func GetSqlHandler() SqlHandler {
	return &sqlManager{}
}

//Interface helper methods
type SqlHandler interface {
	Insert(query string, params ...interface{}) (int64, int64, error)
	Update(query string)
	QueryRow(query string, params ...interface{}) *sql.Row
	QueryRows(query string, params ...interface{}) (*sql.Rows, error)
	AddTransactionQuery(query string)
	AddTransactions(lastTransactionId int, params ...interface{}) int
	PerformTransactions() error
	PerformMultiTransactions(queriesWithParams []string) error
}
