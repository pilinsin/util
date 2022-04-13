package util

import (
	"io/ioutil"
	"strings"
	"time"
	"unicode"
)

const Layout = "2006-1-2 15:4"

type TimeInfo struct {
	Begin, End, Loc string
}

func NewTimeInfo(begin, end, loc string) (*TimeInfo, bool) {
	L, err := time.LoadLocation(loc)
	if err != nil{return nil, err}
	bTime, err := time.ParseInLocation(Layout, begin, L)
	if err != nil{return nil, err}
	eTime, err := time.ParseInLocation(Layout, end, L)
	if err != nil{return nil, err}
	return &TimeInfo{begin.String(), end.String(), L.String()}, nil
}
func (ti TimeInfo) BeginTime() time.Time{
	L, err := time.LoadLocation(ti.Loc)
	if err != nil{return time.Time{}}
	bTime, err := time.ParseInLocation(Layout, ti.Begin, L)
	if err != nil{return time.Time{}}
	return bTime
}
func (ti TimeInfo) EndTime() time.Time{
	L, err := time.LoadLocation(ti.Loc)
	if err != nil{return time.Time{}}
	eTime, err := time.ParseInLocation(Layout, ti.End, L)
	if err != nil{return time.Time{}}
	return eTime
}
func (ti TimeInfo) WithinTime(now time.Time) bool {
	L, _ := time.LoadLocation(ti.Loc)
	bTime, _ := time.ParseInLocation(Layout, ti.Begin, L)
	eTime, _ := time.ParseInLocation(Layout, ti.End, L)
	now = now.In(L)

	return (bTime.Equal(now) || bTime.Before(now)) && eTime.After(now)
}
func (ti TimeInfo) AfterTime(now time.Time) bool {
	L, _ := time.LoadLocation(ti.Loc)
	eTime, _ := time.ParseInLocation(Layout, ti.End, L)
	now = now.In(L)

	return now.Equal(eTime) || now.After(eTime)
}

//https://stackoverflow.com/questions/40120056/get-a-list-of-valid-time-zones-in-go
func GetOsTimeZones() []string {
	var zones []string
	var zoneDirs = []string{
		// Update path according to your OS
		"/usr/share/zoneinfo/",
		"/usr/share/lib/zoneinfo/",
		"/usr/lib/locale/TZ/",
	}

	for _, zd := range zoneDirs {
		zones = walkTzDir(zd, zones)

		for idx, zone := range zones {
			zones[idx] = strings.ReplaceAll(zone, zd+"/", "")
		}
	}

	return zones
}

func walkTzDir(path string, zones []string) []string {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return zones
	}

	isAlpha := func(s string) bool {
		for _, r := range s {
			if !unicode.IsLetter(r) {
				return false
			}
		}
		return true
	}

	for _, info := range fileInfos {
		if info.Name() != strings.ToUpper(info.Name()[:1])+info.Name()[1:] {
			continue
		}

		if !isAlpha(info.Name()[:1]) {
			continue
		}

		newPath := path + "/" + info.Name()

		if info.IsDir() {
			zones = walkTzDir(newPath, zones)
		} else {
			zones = append(zones, newPath)
		}
	}

	return zones
}
