package main

import (
	"time"
)

var statistic AppStatistic

func init() {
	statistic.Status = "ok"
	statistic.Time.Start = time.Now()
}

type AppStatistic struct {
	// mux    sync.Mutex
	Status string `json:"status,omitempty" yaml:"status,omitempty"`
	Time   struct {
		Start time.Time `json:"start,omitempty" yaml:"start,omitempty"`
		Lives struct {
			Duration time.Duration `json:"duration,omitempty" yaml:"duration,omitempty"`
			Sec      float64       `json:"seconds,omitempty" yaml:"seconds,omitempty"`
			Min      float64       `json:"minutes,omitempty" yaml:"minutes,omitempty"`
			Hour     float64       `json:"hours,omitempty" yaml:"hours,omitempty"`
		} `json:"lives,omitempty" yaml:"lives,omitempty"`
	} `json:"time,omitempty" yaml:"time,omitempty"`
}

func (as *AppStatistic) Upgrade() {
	as.Time.Lives.Duration = time.Since(as.Time.Start)
	as.Time.Lives.Sec = as.Time.Lives.Duration.Seconds()
	as.Time.Lives.Min = as.Time.Lives.Duration.Minutes()
	as.Time.Lives.Hour = as.Time.Lives.Duration.Hours()
}
