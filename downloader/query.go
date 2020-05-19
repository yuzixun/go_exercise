package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type CallRecord struct {
	// gorm.Model

	ID             int
	OrganizationId int
	UserId         int
	RecordableId   int
	RecordableType string
	AgentId        string
	CallType       int
	DeviceType     int
	Status         int
	CallId         string
	CallingNumber  string
	CalledNumber   string
	TotalTime      int
	AgentTime      int
	UserTime       int
	Charge         float64
	FileName       string
	Key            string
	WavKey         string
	Content        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	AppType        int
}

func queryCallRecords(modelChan chan CallRecord, quitChan chan int, db *gorm.DB) {
	page, per := 0, 10000

	for cTime := time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC); cTime.Before(time.Now()); cTime = cTime.AddDate(0, 1, 0) {
		for {
			records := []CallRecord{}
			db.Offset(page*per).
				Where("created_at >= ?", cTime).Where("created_at < ?", cTime.AddDate(0, 1, 0)).
				Limit(per).Find(&records)

			if len(records) == 0 {
				break
			}

			for _, record := range records {
				modelChan <- record
			}

			page++
		}
	}

	quitChan <- 0
}
