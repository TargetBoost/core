package models

import "strconv"

func MapToTarget(t Target) TargetService {
	return TargetService{
		UID:        t.UID,
		Title:      t.Title,
		Link:       t.Link,
		Icon:       t.Icon,
		Status:     t.Status,
		Count:      t.Count,
		Total:      t.Total,
		Cost:       t.Cost,
		Cause:      t.Cause,
		TotalPrice: strconv.Itoa(int(t.TotalPrice)),
	}
}

func MapToQueueExecutors(t Queue) QueueToService {
	return QueueToService{
		TID:    t.TID,
		UID:    t.UID,
		Cost:   t.Cost,
		Title:  t.Title,
		Status: t.Status,
	}
}
