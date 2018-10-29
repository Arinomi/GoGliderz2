package main

import (
	"time"
)

// contains the relevant timestamps, used in this file
type Timestamps struct {
	latest      time.Time
	oldest      time.Time
	oldestNewer time.Time
}

// this is returned to the handlers
type Ticker struct {
	Latest       time.Time     `json:"t_latest"`
	Start        time.Time     `json:"t_start"`
	Stop         time.Time     `json:"t_stop"`
	TrackIDs     []string      `json:"tracks"`
	Responsetime time.Duration `json:"responsetime"`
}

// returns a Timestamps struct
func returnTimestamps(ts time.Time) Timestamps {
	tracks := db.GetAll()

	timestamps := Timestamps{}

	if len(tracks) > 0 {
		timestamps.latest = returnLatest(tracks)
		timestamps.oldest = returnOldest(tracks)
		timestamps.oldestNewer = returnOldestNewer(ts, tracks)
	}

	return timestamps
}

// I like descriptive function names.
func returnLatest(tracks []Track) time.Time {
	var latest time.Time

	for _, element := range tracks {
		if element.Timestamp.After(latest) {
			latest = element.Timestamp
		}
	}

	return latest
}

// Still descriptive
func returnOldest(tracks []Track) time.Time {
	var oldest time.Time

	for i, element := range tracks {
		if i == 0 {
			oldest = element.Timestamp
		}

		if element.Timestamp.Before(oldest) {
			oldest = element.Timestamp
		}
	}

	return oldest
}

// Keep on being descriptive
func returnOldestNewer(input time.Time, tracks []Track) time.Time {
	ts := time.Now()
	now := ts

	for _, element := range tracks {
		if element.Timestamp.After(input) && element.Timestamp.Before(ts) {
			ts = element.Timestamp
		}
	}

	if ts == now {
		return time.Time{}
	}

	return ts
}

// Yup. Descriptive. Returns a ticker for api/ticker
func returnTicker() (Ticker, bool) {
	// return the latest timestamp
	// return "start", this is the very first timestamp of them all
	// return "stop", this is the fifth timestamp from "start"
	tracks := db.GetAllSorted("timestamp")
	var ticker Ticker
	ok := true

	if len(tracks) > 0 {
		procStart := time.Now()
		var timestamps []time.Time
		max := 5 // TODO: Make the cap a config value

		for i, element := range tracks {
			if i <= max {
				ticker.TrackIDs = append(ticker.TrackIDs, element.TrackId)
				timestamps = append(timestamps, element.Timestamp)
			}
		}

		ticker.Start = timestamps[0]
		ticker.Stop = timestamps[len(timestamps)-1]
		ticker.Latest = tracks[len(tracks)-1].Timestamp
		ticker.Responsetime = time.Since(procStart)

		return ticker, ok
	} else {
		ok = false
		return ticker, ok
	}
}

// Returns a Ticker, but made for /api/ticker/<timestamp>
func returnTickerTimestamp(input time.Time) (Ticker, bool) {
	tracks := db.GetAllSorted("timestamp")
	var ticker Ticker
	ok := true
	if len(tracks) > 0 {
		procStart := time.Now()
		timestamps := returnTimestamps(input)
		max := 5
		var tsArr []time.Time // timestamp arrays, all timestamps after the start

		ticker.Latest = timestamps.latest
		ticker.Start = timestamps.oldestNewer

		if ticker.Start.IsZero() && timestamps.oldestNewer.IsZero() {
			ok = false
			return Ticker{}, ok
		}

		for _, element := range tracks {
			if element.Timestamp.Equal(ticker.Start) {
				// Putting the first ID in there.
				// Also, doing ticker.TrackIDs[0] = element.TrackId DOES NOT WORK!
				// Because ticker.TrackIDs does not have any indices yet, thus not even a 0 index!
				// I spent TWO HOURS wondering why my program had a hissy fit over this
				// It was obvious, yet my mouse brain managed to oversee that
				ticker.TrackIDs = append(ticker.TrackIDs, element.TrackId)
			}

			if element.Timestamp.After(ticker.Start) && len(tsArr) <= max {
				tsArr = append(tsArr, element.Timestamp)
				ticker.TrackIDs = append(ticker.TrackIDs, element.TrackId)
			}
		}

		ticker.Stop = tsArr[len(tsArr)-1]
		ticker.Responsetime = time.Since(procStart)

		return ticker, ok
	} else {
		ok = false
		return ticker, ok
	}
}
