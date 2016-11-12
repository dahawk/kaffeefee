package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
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
type userList []user

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

//TODO should be testable
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
	populate := populateUser()
	avgsDay := populate.emptyMap()
	for k, v := range countsDaily {
		avgsDay[k] = (float64(v) / float64(cntDaily[k]))
		if math.IsNaN(avgsDay[k]) {
			avgsDay[k] = 0
		}
	}
	ret["DailyAvgs"] = avgsDay

	avgsWeek := populate.emptyMap()
	for k, v := range countsWeekly {

		avgsWeek[k] = (float64(v) / float64(cntWeekly[k]))
		if math.IsNaN(avgsWeek[k]) {
			avgsWeek[k] = 0
		}
	}
	ret["WeeklyAvgs"] = avgsWeek

	avgsMonth := populate.emptyMap()
	for k, v := range countsMonthly {

		avgsMonth[k] = (float64(v) / float64(cntMonthly[k]))
		if math.IsNaN(avgsMonth[k]) {
			avgsMonth[k] = 0
		}
	}
	ret["MonthlyAvgs"] = avgsMonth

	return ret
}

func (users userList) emptyMap() map[string]float64 {
	ret := make(map[string]float64, 0)
	for _, v := range users {
		ret[v.Name] = 0.0
	}
	return ret
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

//TODO pull out populateUser
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

//TODO possibly split into get user and get count
func populateUser() userList {
	var u userList
	from, _ := getPeriod(time.Now())

	err := db.Select(&u, "select * from users where active=true")
	if err != nil {
		panic(err)
	}

	for i, user := range u {
		var cnt []int
		err = db.Select(&cnt, "select sum(cnt) from log where ts > $1 and userid = $2", from.Unix(), user.UserID)
		if err == nil {
			u[i].Today = cnt[0]
		} else {
			fmt.Println(err)
		}
	}

	return u
}
