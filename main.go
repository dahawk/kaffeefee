package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type log struct {
	ID        int
	UserID    int   `db:"userid"`
	Number    int   `db:"cnt"`
	Timestamp int64 `db:"ts"`
}

type user struct {
	Name   string `db:"name"`
	UserID int    `db:"id"`
	Today  int
}

type userlogs []log

type dayData struct {
	Hour string
	Cnt  int
}

type weekData struct {
	Day string
	Cnt int
}

var row = make([]interface{}, 0)
var db *sqlx.DB

const (
	truncDay  = time.Hour * 24
	truncWeek = truncDay * 7
)

func main() {
	var err error
	defaultDB := "postgres://<user>:<pwd>@<host>/<db>?sslmode=disable"

	dbString := os.Getenv("DB")
	if dbString == "" {
		dbString = defaultDB
	}

	db, err = sqlx.Open("postgres", dbString)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("server started")
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.HandleFunc("/", logcoffee)
	mux.HandleFunc("/stats", stats)
	mux.HandleFunc("/graph", graph)
	mux.HandleFunc("/json", jsonAjax)
	mux.HandleFunc("/dailyChart", dayGraph)
	mux.HandleFunc("/jsonDaily", jsonDaily)
	mux.HandleFunc("/weeklyChart", weekGraph)
	mux.HandleFunc("/jsonWeekly", jsonWeekly)

	http.ListenAndServe(":8080", mux)

}

func logcoffee(w http.ResponseWriter, r *http.Request) {
	users := populateUser()
	var occurred error
	from, to := getPeriod(time.Now())

	t, err := template.ParseFiles("tpl/log.tpl")
	store := false
	if err != nil {
		fmt.Println(err)
		return
	}

	if id, ok := r.URL.Query()["id"]; r.Method == "GET" && ok {
		occurred = storeLog(id[0])
	}

	p, err := renderPage(store, users)
	p["From"] = from.Format("02.01.2006")
	p["To"] = to.Format("02.01.2006")
	if occurred != nil {
		p["Error"] = occurred
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		fmt.Println(err)
		return
	}
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

func renderPage(store bool, users []user) (map[string]interface{}, error) {
	p := make(map[string]interface{}, 0)
	var data []user
	for _, v := range users {
		data = append(data, v)
	}
	p["Users"] = data
	if _, ok := p["Error"]; !ok {
		p["Store"] = store
	}
	return p, nil
}

func stats(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("tpl/stats.tpl")

	if err != nil {
		fmt.Println(err)
		return
	}

	dayAggr := make(map[string]int, 0)
	weekAggr := make(map[string]int, 0)
	monthAggr := make(map[string]int, 0)
	total := make(map[string]int, 0)
	day := time.Now().Truncate(truncDay)
	week := time.Now().Truncate(truncWeek)
	now := time.Now()
	month := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	users := populateUser()

	for _, v := range users {
		dayAggr[v.Name] = 0
		weekAggr[v.Name] = 0
		total[v.Name] = 0
		var logs []log
		er := db.Select(&logs, "select * from log where userid=$1", v.UserID)
		if er != nil {
			fmt.Println(er)
			return
		}

		for _, l := range logs {

			total[v.Name] += l.Number
			if day.Unix() <= l.Timestamp {
				dayAggr[v.Name] += l.Number
			}
			if week.Unix() <= l.Timestamp {
				weekAggr[v.Name] += l.Number
			}

			if month.Unix() <= l.Timestamp {
				monthAggr[v.Name] += l.Number
			}
		}
	}
	p := make(map[string]interface{}, 0)
	avgs := calculateAverages(users)
	p["Weekly"] = weekAggr
	p["Daily"] = dayAggr
	p["Monthly"] = monthAggr
	p["Total"] = total
	p["DailyAvgs"] = avgs["DailyAvgs"]
	p["WeeklyAvgs"] = avgs["WeeklyAvgs"]
	p["MonthlyAvgs"] = avgs["MonthlyAvgs"]
	errExec := t.Execute(w, p)
	if errExec != nil {
		fmt.Println(errExec)
	}
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

func calculateAverages(users []user) map[string]map[string]float64 {
	minTimestamp := getMinTimestamp()
	ret := make(map[string]map[string]float64, 0)

	now := time.Now()
	subDay := 60 * 60 * 24
	subWeek := subDay * 7

	countsDaily := make(map[string]int64, 0)
	countsWeekly := make(map[string]int64, 0)
	countsMonthly := make(map[string]int64, 0)
	cntDaily := make(map[string]int64, 0)
	cntWeekly := make(map[string]int64, 0)
	cntMonthly := make(map[string]int64, 0)
	for _, u := range users {

		countsDaily[u.Name] = 0
		countsWeekly[u.Name] = 0
		countsMonthly[u.Name] = 0

		cntDaily[u.Name] = 0
		cntWeekly[u.Name] = 0
		cntMonthly[u.Name] = 0
	}
	for _, u := range users {
		toDaily := now.Truncate(truncDay).Unix()
		fromDaily := toDaily - int64(subDay)
		toWeekly := now.Truncate(truncWeek).Unix()
		fromWeekly := toWeekly - int64(subWeek)

		l := getLogsForUser(u.UserID)
		for fromDaily >= minTimestamp {
			cntDaily[u.Name]++
			countsDaily[u.Name] += l.calculateSumForUser(fromDaily, toDaily)
			toDaily = fromDaily
			fromDaily = fromDaily - int64(subDay)
		}

		toWeekly = now.Truncate(truncWeek).Unix()
		fromWeekly = toWeekly - int64(subWeek)

		for fromWeekly >= minTimestamp {
			cntWeekly[u.Name]++

			countsWeekly[u.Name] += l.calculateSumForUser(fromWeekly, toWeekly)

			toWeekly = fromWeekly
			fromWeekly = fromWeekly - int64(subWeek)
		}

		toTmp := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		toMonthly := toTmp.Add(-time.Second).Unix()
		fromTmp := time.Date(toTmp.Year(), toTmp.Month()-1, 1, 0, 0, 0, 0, time.UTC)
		fromMonthly := fromTmp.Unix()
		for fromMonthly >= minTimestamp {
			cntMonthly[u.Name]++

			countsMonthly[u.Name] += l.calculateSumForUser(fromMonthly, toMonthly)

			toTmp = fromTmp.Add(-time.Second)
			fromTmp = time.Date(toTmp.Year(), toTmp.Month()-1, 1, 0, 0, 0, 0, time.UTC)
			toMonthly = toTmp.Unix()
			fromMonthly = fromTmp.Unix()
		}
	}
	avgsDay := emptyMap()
	for k, v := range countsDaily {
		avgsDay[k] = (float64(v) / float64(cntDaily[k]))
		if math.IsNaN(avgsDay[k]) {
			avgsDay[k] = 0
		}
	}
	ret["DailyAvgs"] = avgsDay

	avgsWeek := emptyMap()
	for k, v := range countsWeekly {

		avgsWeek[k] = (float64(v) / float64(cntWeekly[k]))
		if math.IsNaN(avgsWeek[k]) {
			avgsWeek[k] = 0
		}
	}
	ret["WeeklyAvgs"] = avgsWeek

	avgsMonth := emptyMap()
	for k, v := range countsMonthly {

		avgsMonth[k] = (float64(v) / float64(cntMonthly[k]))
		if math.IsNaN(avgsMonth[k]) {
			avgsMonth[k] = 0
		}
	}
	ret["MonthlyAvgs"] = avgsMonth

	return ret
}

func emptyMap() map[string]float64 {
	ret := make(map[string]float64, 0)
	users := populateUser()
	for _, v := range users {
		ret[v.Name] = 0.0
	}
	return ret
}

func graph(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tpl/graph.tpl")
	if err != nil {
		fmt.Println(err)
		return
	}

	uname := r.URL.Query().Get("user")
	p := make(map[string]interface{}, 0)
	p["UserId"] = uname

	t.Execute(w, p)
}

func jsonAjax(w http.ResponseWriter, r *http.Request) {
	minTimestamp := getMinTimestamp()
	encoder := json.NewEncoder(w)

	selector := r.URL.Query().Get("interval")
	user := r.URL.Query().Get("user")

	var ret []map[string]string

	switch selector {
	case "daily":
		ret = mapping(getDailyCount(user, minTimestamp))
	case "weekly":
		ret = mapping(getWeeklyCount(user, minTimestamp))
	case "monthly":
		ret = mapping(getMonthlyCount(user, minTimestamp))
	}

	encoder.Encode(ret)
}

func dayGraph(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tpl/distributionDay.tpl")
	if err != nil {
		fmt.Println(err)
		return
	}

	errExec := t.Execute(w, nil)
	if errExec != nil {
		fmt.Println(errExec)
	}
}

func jsonDaily(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	logs := getAllLogs()

	data := make([]dayData, 24)
	for i := 0; i < 24; i++ {
		var hour string
		hour = fmt.Sprintf("%d:00", i)
		if i < 10 {
			hour = fmt.Sprintf("0%d:00", i)
		}

		data[i].Hour = hour
		data[i].Cnt = 0
	}

	currentDay := 0
	days := 0

	for _, l := range logs {
		t := time.Unix(l.Timestamp, 0)
		data[t.Hour()].Cnt++
		if t.Day() != currentDay {
			currentDay = t.Day()
			days++
		}
	}

	/*for i := 0; i < 24; i++ {
		if days > 0 {
			data[i].Cnt = data[i].Cnt / days
		}
	}*/

	encoder.Encode(data)
}

func weekGraph(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tpl/distributionWeek.tpl")
	if err != nil {
		fmt.Println(err)
		return
	}

	errExec := t.Execute(w, nil)
	if errExec != nil {
		fmt.Println(errExec)
	}
}

func jsonWeekly(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	logs := getAllLogs()

	data := make([]weekData, 7)
	weeks := 0
	currentWeek := 0
	data[0].Day = "Sunday"
	data[1].Day = "Monday"
	data[2].Day = "Tuesday"
	data[3].Day = "Wednesday"
	data[4].Day = "Thursday"
	data[5].Day = "Friday"
	data[6].Day = "Saturday"

	for i := 0; i < 7; i++ {
		data[i].Cnt = 0
	}

	for _, l := range logs {
		t := time.Unix(l.Timestamp, 0)
		data[t.Weekday()].Cnt++
		_, w := t.ISOWeek()
		if w != currentWeek {
			currentWeek = w
			weeks++
		}
	}
	/*for i := 0; i < 7; i++ {
		if weeks > 0 {
			data[i].Cnt = data[i].Cnt / weeks
		}
	}*/

	encoder.Encode(data)
}

func getDailyCount(userID string, minTimestamp int64) map[time.Time]int64 {
	ret := make(map[time.Time]int64, 0)
	l := getLogsForUser(getUserIDFromName(userID))

	now := time.Now()
	var subDay int64 = 60 * 60 * 24

	from := now.Truncate(truncDay).Unix()
	to := from + subDay
	for from >= minTimestamp {
		ret[time.Unix(from, 0)] = l.calculateSumForUser(from, to)

		to = from
		from = from - subDay
	}

	return ret
}

func getWeeklyCount(userID string, minTimestamp int64) map[time.Time]int64 {
	ret := make(map[time.Time]int64, 0)
	l := getLogsForUser(getUserIDFromName(userID))

	now := time.Now()
	var subWeek int64 = 60 * 60 * 24 * 7

	from := now.Truncate(truncWeek).Unix()
	to := from + subWeek
	for to >= minTimestamp {
		//fix
		ret[time.Unix(from, 0)] = l.calculateSumForUser(from, to)

		to = from
		from = from - subWeek
	}

	return ret
}

func getMonthlyCount(userID string, minTimestamp int64) map[time.Time]int64 {
	ret := make(map[time.Time]int64, 0)
	l := getLogsForUser(getUserIDFromName(userID))

	now := time.Now()

	fromTmp := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	from := fromTmp.Unix()
	toTmp := fromTmp.AddDate(0, 1, -1)
	to := toTmp.Unix()
	for to >= minTimestamp {
		//fix
		ret[time.Unix(from, 0)] = l.calculateSumForUser(from, to)

		toTmp = fromTmp.AddDate(0, 0, -1)
		fromTmp = toTmp.AddDate(0, -1, 1)
		to = toTmp.Unix()
		from = fromTmp.Unix()
	}

	return ret
}

func getUserIDFromName(name string) int {
	users := populateUser()
	for _, v := range users {
		if v.Name == name {
			return v.UserID
		}
	}
	return -1
}

func (logs userlogs) calculateSumForUser(from, to int64) int64 {
	var sum int64
	for _, l := range logs {

		if l.Timestamp >= from && l.Timestamp <= to {
			sum++
		}
	}
	return sum
}

func mapping(input map[time.Time]int64) []map[string]string {
	var ret []map[string]string

	for k, v := range input {
		entry := make(map[string]string)

		entry["datum"] = k.Format("2006-01-02")
		entry["count"] = fmt.Sprintf("%d", v)

		ret = append(ret, entry)
	}

	return ret
}

func getPeriod(now time.Time) (time.Time, time.Time) {
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	for ; from.Weekday() != time.Monday; from = from.AddDate(0, 0, -1) {
	}
	for ; to.Weekday() != time.Sunday; to = to.AddDate(0, 0, 1) {
	}

	return from, to
}

func populateUser() []user {
	var users []user
	from, _ := getPeriod(time.Now())

	err := db.Select(&users, "select * from users where active=true")
	if err != nil {
		panic(err)
	}

	for i, u := range users {
		var cnt []int
		err = db.Select(&cnt, "select sum(cnt) from log where ts > $1 and userid = $2", from.Unix(), u.UserID)
		if err == nil {
			users[i].Today = cnt[0]
		} else {
			fmt.Println(err)
		}
	}

	return users
}

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
