package main

import (
	"crypto/md5"
	"encoding/hex"
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
	Active bool   `db:"active"`
	Mail   string `db:"mail"`
	Image  string
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

	dbString := os.Getenv("DB")
	if dbString == "" {
		fmt.Printf("DB is not set")
		os.Exit(1)
	}

	db, err = sqlx.Open("postgres", dbString)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("server started")
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./static/fonts"))))
	mux.HandleFunc("/", logcoffee)
	mux.HandleFunc("/stats", stats)
	mux.HandleFunc("/graph", graph)
	mux.HandleFunc("/json", jsonAjax)
	mux.HandleFunc("/dailyChart", dayGraph)
	mux.HandleFunc("/jsonDaily", jsonDaily)
	mux.HandleFunc("/weeklyChart", weekGraph)
	mux.HandleFunc("/jsonWeekly", jsonWeekly)
	mux.HandleFunc("/admin", admin)
	mux.HandleFunc("/editUser", editUser)
	mux.HandleFunc("/createUser", addUser)

	http.ListenAndServe(":8080", mux)

}

func renderPage(store bool, users []user) map[string]interface{} {
	p := make(map[string]interface{}, 0)
	var data []user
	for _, v := range users {
		data = append(data, v)
	}
	p["Users"] = data
	if _, ok := p["Error"]; !ok {
		p["Store"] = store
	}
	return p
}

func calculateUserAverages(u user, minTimestamp int64) (dailyAvg, weeklyAvg, monthlyAvg float64) {
	var (
		//number of coffies
		countDaily   int64
		countWeekly  int64
		countMonthly int64

		//number of days/weeks/month considered
		cntDays   int64
		cntWeeks  int64
		cntMonths int64
	)

	now := time.Now()
	subDay := 60 * 60 * 24
	subWeek := subDay * 7

	toDaily := now.Truncate(truncDay).Unix()
	fromDaily := toDaily - int64(subDay)

	l := getLogsForUser(u.UserID)
	for fromDaily >= minTimestamp {
		cntDays++
		countDaily += l.calculateSumForUser(fromDaily, toDaily)
		toDaily = fromDaily
		fromDaily = fromDaily - int64(subDay)
	}

	toWeekly := now.Truncate(truncWeek).Unix()
	fromWeekly := toWeekly - int64(subWeek)

	for fromWeekly >= minTimestamp {
		cntWeeks++

		countWeekly += l.calculateSumForUser(fromWeekly, toWeekly)

		toWeekly = fromWeekly
		fromWeekly = fromWeekly - int64(subWeek)
	}

	toTmp := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	toMonthly := toTmp.Add(-time.Second).Unix()
	fromTmp := time.Date(toTmp.Year(), toTmp.Month()-1, 1, 0, 0, 0, 0, time.UTC)
	fromMonthly := fromTmp.Unix()
	for fromMonthly >= minTimestamp {
		cntMonths++

		countMonthly += l.calculateSumForUser(fromMonthly, toMonthly)

		toTmp = fromTmp.Add(-time.Second)
		fromTmp = time.Date(toTmp.Year(), toTmp.Month()-1, 1, 0, 0, 0, 0, time.UTC)
		toMonthly = toTmp.Unix()
		fromMonthly = fromTmp.Unix()
	}

	dailyAvg = checkNaN(float64(countDaily) / float64(cntDays))
	weeklyAvg = checkNaN(float64(countWeekly) / float64(cntWeeks))
	monthlyAvg = checkNaN(float64(countMonthly) / float64(cntMonths))

	return dailyAvg, weeklyAvg, monthlyAvg
}

func checkNaN(in float64) float64 {
	if math.IsNaN(in) {
		return 0.0
	}
	return in
}

func (users userList) emptyMap() map[string]float64 {
	ret := make(map[string]float64, 0)
	for _, v := range users {
		ret[v.Name] = 0.0
	}
	return ret
}

func (logs userlogs) getDailyCount(now time.Time, minTimestamp int64) map[time.Time]int64 {
	ret := make(map[time.Time]int64, 0)

	var subDay int64 = 60 * 60 * 24

	from := now.Truncate(truncDay).Unix()
	to := from + subDay
	for from >= minTimestamp {
		ret[time.Unix(from, 0)] = logs.calculateSumForUser(from, to)

		to = from
		from = from - subDay
	}

	return ret
}

func (logs userlogs) getWeeklyCount(now time.Time, minTimestamp int64) map[time.Time]int64 {
	ret := make(map[time.Time]int64, 0)

	var subWeek int64 = 60 * 60 * 24 * 7

	from := now.Truncate(truncWeek).Unix()
	to := from + subWeek
	for to >= minTimestamp {
		//fix
		ret[time.Unix(from, 0)] = logs.calculateSumForUser(from, to)

		to = from
		from = from - subWeek
	}

	return ret
}

func (logs userlogs) getMonthlyCount(now time.Time, minTimestamp int64) map[time.Time]int64 {
	ret := make(map[time.Time]int64, 0)

	fromTmp := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	from := fromTmp.Unix()
	toTmp := fromTmp.AddDate(0, 1, -1)
	to := toTmp.Unix()
	for to >= minTimestamp {
		//fix
		ret[time.Unix(from, 0)] = logs.calculateSumForUser(from, to)

		toTmp = fromTmp.AddDate(0, 0, -1)
		fromTmp = toTmp.AddDate(0, -1, 1)
		to = toTmp.Unix()
		from = fromTmp.Unix()
	}

	return ret
}

func (users userList) getUserIDFromName(name string) int {
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

func hasLocalImage(name string) (bool, string) {
	wd, err := os.Getwd()
	if err != nil {
		return false, ""
	}
	path := fmt.Sprintf("%s/static/%s.png", wd, name)
	fh, err := os.Open(path)
	if err != nil {
		return false, ""
	}
	fh.Close()

	ret := fmt.Sprintf("/static/%s.png", name)

	return true, ret
}

func hasGravatarImage(mail string) (bool, string) {
	if mail == "" {
		return false, ""
	}
	hash := md5.Sum([]byte(mail))
	url := fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=404&s=60", hex.EncodeToString(hash[:]))

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, ""
	}

	return true, url
}
