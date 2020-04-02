package main

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	sql "github.com/jmoiron/sqlx"
)

var (
	// EvaluatorID is your student ID
	EvaluatorID = "18307130003"
	// SubmissionDir is the relative path of the submission directory of assignment 1
	SubmissionDir = "../../../ass1/submission/"
	// User is the username used to connect to the database
	User = "root"
	// Password for the username
	Password = "Hakula"
)

// ConcurrentCompareAndInsert is similar with compareAndInsert in `main.go`, but it is concurrent and faster!
func ConcurrentCompareAndInsert(subs map[string]*Submission) {
	start := time.Now()
	defer func() {
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/ass1_result_evaluated_by_%s", User, Password, EvaluatorID))
		if err != nil {
			panic(nil)
		}
		rows, err := db.Query("SELECT COUNT(*) FROM comparison_result")
		if err != nil {
			panic(err)
		}
		rows.Next()
		var cnt int
		err = rows.Scan(&cnt)
		if err != nil {
			panic(err)
		}
		if cnt == 0 {
			panic("ConcurrentCompareAndInsert Not Implemented")
		}
		fmt.Println("ConcurrentCompareAndInsert takes ", time.Since(start))
	}()

	// Connect to the database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/ass1_result_evaluated_by_%s", User, Password, EvaluatorID))
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(64)
	db.SetMaxIdleConns(16)

	type query struct {
		sid string // the ID of submitter
		cid string // the ID of comparer
		qid int    // the ID of SQL query
		res int    // the result of comparison
	}

	// Insert results into the database
	insert := func(q query) {
		s := fmt.Sprintf("INSERT INTO comparison_result VALUES ('%s', '%s', %d, %d)", q.sid, q.cid, q.qid, q.res)
		_, err := db.Exec(s)
		if err != nil {
			fmt.Println(s)
			panic(err)
		}
	}

	// Check if the answers are correct
	var wg sync.WaitGroup
	for sid, submitter := range subs {
		for cid, comparer := range subs {
			wg.Add(1)
			go func(sid, cid string, submitter, comparer *Submission) {
				for i := 0; i < NumSQL; i++ {
					var res int
					if reflect.DeepEqual(submitter.sqlResults[i], comparer.sqlResults[i]) {
						res = 1
					} else {
						res = 0
					}
					insert(query{sid, cid, i + 1, res})
				}
				wg.Done()
			}(sid, cid, submitter, comparer)
		}
	}
	wg.Wait()
	// YOUR CODE END
}

// GetScoreSQL returns a string which contains only ONE SQL to be executed, which collects the data in table
// `comparision_result` and inserts the score of each submitter on each query into table `score`
func GetScoreSQL() string {
	var SQL string
	SQL = "SELECT 1" // ignore this line, it just makes the returned SQL a valid SQL if you haven't written yours.
	// YOUR CODE BEGIN

	// YOUR CODE END
	return SQL
}

// GetScore reads your score from table `score`
func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN

	// YOUR CODE END
}
