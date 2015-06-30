package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	register("/okundumu", okundumu)
	register("/iftar", iftar)
	register("/sahur", sahur)
	fill()
}

const (
	timeFormat     = "_2 Jan 2006"
	timeFormatLong = "_2 Jan 2006 15:04"
)

var no = []string{
	"okunmadi",
	"hayir",
	"Hayir",
	"Hayir.",
	"hayir.",
	"hayir okunmadi",
	"no",
	"no.",
	"NO.",
	"NO!",
	"NO",
	"nope",
	"nicht",
}

type timePair struct {
	iftar time.Time
	sahur time.Time
}

var callTime = map[string]timePair{}

var _callTime = map[string]struct {
	iftar, sahur string
}{
	"18 Jun 2015": {"20:48", "03:22"},
	"19 Jun 2015": {"20:48", "03:22"},
	"20 Jun 2015": {"20:48", "03:22"},
	"21 Jun 2015": {"20:48", "03:22"},
	"22 Jun 2015": {"20:49", "03:23"},
	"23 Jun 2015": {"20:49", "03:23"},
	"24 Jun 2015": {"20:49", "03:23"},
	"25 Jun 2015": {"20:49", "03:23"},
	"26 Jun 2015": {"20:49", "03:24"},
	"27 Jun 2015": {"20:49", "03:24"},
	"28 Jun 2015": {"20:49", "03:25"},
	"29 Jun 2015": {"20:49", "03:26"},
	"30 Jun 2015": {"20:49", "03:26"},
	"1 Jul 2015":  {"20:49", "03:27"},
	"2 Jul 2015":  {"20:49", "03:28"},
	"3 Jul 2015":  {"20:49", "03:28"},
	"4 Jul 2015":  {"20:49", "03:29"},
	"5 Jul 2015":  {"20:48", "03:30"},
	"6 Jul 2015":  {"20:48", "03:31"},
	"7 Jul 2015":  {"20:48", "03:32"},
	"8 Jul 2015":  {"20:48", "03:33"},
	"9 Jul 2015":  {"20:47", "03:35"},
	"10 Jul 2015": {"20:47", "03:35"},
	"11 Jul 2015": {"20:46", "03:37"},
	"12 Jul 2015": {"20:46", "03:38"},
	"13 Jul 2015": {"20:46", "03:39"},
	"14 Jul 2015": {"20:45", "03:40"},
	"15 Jul 2015": {"20:44", "03:41"},
	"16 Jul 2015": {"20:44", "03:43"},
}

func fill() {
	loc, _ := time.LoadLocation("Europe/Istanbul")

	for k, v := range _callTime {
		iftar, _ := time.ParseInLocation(timeFormatLong, fmt.Sprintf("%v %v", k, v.iftar), loc)
		sahur, _ := time.ParseInLocation(timeFormatLong, fmt.Sprintf("%v %v", k, v.sahur), loc)

		callTime[k] = timePair{iftar, sahur}
	}
}

func okundumu(args ...string) string {
	loc, _ := time.LoadLocation("Europe/Istanbul")

	now := time.Now().In(loc)
	nowstr := now.Format(timeFormat)

	timepair, ok := callTime[nowstr]
	if !ok {
		return "galiba oruc bitti"
	}

	if now.After(timepair.iftar) {
		return "okundu"
	}

	if now.Before(timepair.sahur) {
		return "sahur daha okunmadi"
	}

	if now.After(timepair.sahur) && now.Hour() < 6 {
		return "sahur okundu da daha iftara cok var"
	}

	// after sahur and before iftar, hence NO
	return no[rand.Intn(len(no))]
}

func iftar(args ...string) string {
	loc, _ := time.LoadLocation("Europe/Istanbul")

	now := time.Now().In(loc)
	nowstr := now.Format(timeFormat)

	timepair, ok := callTime[nowstr]
	if !ok {
		return "galiba oruc bitti"
	}

	return timepair.iftar.Format("15:04")
}

func sahur(args ...string) string {
	loc, _ := time.LoadLocation("Europe/Istanbul")

	now := time.Now().In(loc)
	nowstr := now.Format(timeFormat)

	timepair, ok := callTime[nowstr]
	if !ok {
		return "galiba oruc bitti"
	}

	return timepair.sahur.Format("15:04")
}
