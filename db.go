package main

import (
	"fmt"
	"strconv"
	"time"
)

func getUsersToday(id string) int {
	now := time.Now()
	from := now.Truncate(truncDay)

	var cnt []int
	err := db.Select(&cnt, "select sum(cnt) from log where ts > $1 and userid = $2", from.Unix(), id)
	if err == nil {
		return cnt[0]
	}
	panic(err)
}

//TODO kill the magic number
func getMinTimestamp() int64 {
	lowestMin := int64(1451602800)

	var ts int64

	err := db.Get(&ts, "select min(ts) as mints from log")
	if err != nil {
		panic(err)
	}

	if ts < lowestMin {
		return lowestMin
	}
	return ts
}

func getLogsForUser(id int) userlogs {
	var ret userlogs

	db.Select(&ret, "select * from log where userid=$1 order by ts desc", id)

	return ret
}

func getAllLogs() userlogs {
	var ret userlogs

	db.Select(&ret, "select * from log order by ts asc")

	return ret
}

func storeLog(id string) error {
	user, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	l := log{
		UserID:    user,
		Number:    1,
		Timestamp: time.Now().Unix(),
	}

	tx, err := db.Begin()

	if err != nil {
		fmt.Println(err)
		return err
	}

	query := "insert into log (userid,cnt,ts) values ($1,$2,$3)"

	_, err = db.Exec(query, l.UserID, l.Number, l.Timestamp)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return err
	}

	tx.Commit()
	return nil
}
