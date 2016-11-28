package main

import (
	"errors"
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

func getAllUsers() (userList, error) {
	var ret userList

	err := db.Select(&ret, "select * from users")
	if err != nil {
		return nil, err
	}

	return ret, err
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

func getUserByName(name string) (user, error) {
	ret := []user{}
	err := db.Select(&ret, "select * from users where name=$1", name)
	if err != nil {
		return user{}, err
	}
	if len(ret) != 1 {
		return user{}, errors.New("unexpected number of results returned")
	}

	return ret[0], nil
}

func toggleUserActive(u user) error {
	isActive := !u.Active
	_, err := db.Exec("update users set active=$1 where id=$2", isActive, u.UserID)
	return err
}

func deleteUser(u user) error {
	_, err := db.Exec("delete from users where id=$1", u.UserID)
	return err
}

func updateUser(u user) error {
	_, err := db.Exec("update users set name=$1, mail=$2 where id=$3", u.Name, u.Mail, u.UserID)
	return err
}

func insertUser(u user) error {
	_, err := db.Exec("insert into users (name, mail, active) values ($1,$2,$3)",
		u.Name, u.Mail, u.Active)

	return err
}

func getActiveUsers() (userList, error) {
	u := userList{}
	err := db.Select(&u, "select * from users where active=true")
	if err != nil {
		return userList{}, err
	}

	return u, nil
}
