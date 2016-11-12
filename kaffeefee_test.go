package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGetPeriod(t *testing.T) {
	now := time.Now()

	from, to := getPeriod(now)

	if from.Weekday() != time.Monday {
		t.Error("expected from: Monday, got ", from.Weekday())
	}
	if from.Hour() != 0 || from.Minute() != 0 || from.Second() != 0 || from.Nanosecond() != 0 {
		t.Error("expected from: 00:00:00.0, got ",
			fmt.Sprintf("%d:%d:%d.%d", from.Hour(), from.Minute(), from.Second(), from.Nanosecond()))
	}

	if to.Weekday() != time.Sunday {
		t.Error("expected from: Sunday, got ", from.Weekday())
	}
	if to.Hour() != 23 || to.Minute() != 59 || to.Second() != 59 || to.Nanosecond() != 0 {
		t.Error("expected from: 00:00:00.0, got ",
			fmt.Sprintf("%d:%d:%d.%d", to.Hour(), to.Minute(), to.Second(), to.Nanosecond()))
	}
}

func TestMapping(t *testing.T) {
	input := map[time.Time]int64{}
	layout := "2006-01-02 15:04:05"

	inputDate, err := time.Parse(layout, "2016-11-12 15:00:00")
	if err != nil {
		t.Error(err)
	}
	input[inputDate] = 10
	inputDate, err = time.Parse(layout, "2016-11-10 12:00:00")
	if err != nil {
		t.Error(err)
	}
	input[inputDate] = 11
	inputDate, err = time.Parse(layout, "2016-11-11 15:00:00")
	if err != nil {
		t.Error(err)
	}
	input[inputDate] = 12

	output := mapping(input)

	if len(output) != 3 {
		t.Error("expected length 3, was ", len(output))
	}

	element := output[0]
	if d, ok := element["datum"]; !ok || d != "2016-11-12" {
		t.Error("error accessing date")
	}
	if c, ok := element["count"]; !ok || c != "10" {
		t.Error("error accessing count")
	}
}

func TestCalculateSumForUser(t *testing.T) {
	logs := userlogs{
		log{ID: 1, UserID: 1, Number: 1, Timestamp: 10},
		log{ID: 2, UserID: 1, Number: 1, Timestamp: 12},
		log{ID: 3, UserID: 1, Number: 1, Timestamp: 14},
		log{ID: 4, UserID: 1, Number: 1, Timestamp: 16},
		log{ID: 5, UserID: 1, Number: 1, Timestamp: 18},
		log{ID: 1, UserID: 1, Number: 1, Timestamp: 20},
		log{ID: 2, UserID: 1, Number: 1, Timestamp: 22},
		log{ID: 3, UserID: 1, Number: 1, Timestamp: 24},
		log{ID: 4, UserID: 1, Number: 1, Timestamp: 26},
		log{ID: 5, UserID: 1, Number: 1, Timestamp: 28},
	}

	sum := logs.calculateSumForUser(10, 28)
	if sum != 10 {
		t.Error("expected 10, got ", sum)
	}
	sum = logs.calculateSumForUser(20, 28)
	if sum != 5 {
		t.Error("expected 5, got ", sum)
	}
}

func TestEmptyUser(t *testing.T) {
	users := userList{
		user{UserID: 1, Name: "user 1", Today: 10},
		user{UserID: 2, Name: "user 2", Today: 10},
		user{UserID: 3, Name: "user 3", Today: 10},
	}

	userMap := users.emptyMap()
	if len(userMap) != 3 {
		t.Error("expected 3 elements, got ", len(userMap))
	}
	if cnt, ok := userMap["user 1"]; !ok || cnt != 0.0 {
		t.Error("error accessing userMap")
	}
}
