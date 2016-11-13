package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

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
	users := populateUser()
	l := getLogsForUser(users.getUserIDFromName(user))
	now := time.Now()

	switch selector {
	case "daily":
		ret = mapping(l.getDailyCount(now, minTimestamp))
	case "weekly":
		ret = mapping(l.getWeeklyCount(now, minTimestamp))
	case "monthly":
		ret = mapping(l.getMonthlyCount(now, minTimestamp))
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
	encoder.Encode(data)
}
