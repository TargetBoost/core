package models

import "strconv"

func MapToTarget(t Target) TargetService {
	return TargetService{
		ID:         t.ID,
		UID:        t.UID,
		Title:      t.Title,
		Link:       t.Link,
		Icon:       t.Icon,
		Status:     t.Status,
		Count:      strconv.Itoa(int(t.Count)),
		Total:      strconv.Itoa(int(t.Total)),
		Cost:       t.Cost,
		Cause:      t.Cause,
		TotalPrice: strconv.Itoa(int(t.TotalPrice)),
		CMFileID:   t.CMFileID,
		Bio:        t.Bio,
	}
}

func MapToTargetAdmin(t TargetToAdmin) TargetService {
	return TargetService{
		ID:         t.ID,
		UID:        t.UID,
		Title:      t.Title,
		Link:       t.Link,
		Icon:       t.Icon,
		Status:     t.Status,
		Count:      strconv.Itoa(int(t.Count)),
		Total:      strconv.Itoa(int(t.Total)),
		Cost:       t.Cost,
		Cause:      t.Cause,
		TotalPrice: strconv.Itoa(int(t.TotalPrice)),
		Login:      t.Login,
	}
}

func MapToQueueExecutors(t QueueToExecutors) QueueToService {
	return QueueToService{
		ID:     t.ID,
		TID:    t.TID,
		UID:    t.UID,
		Cost:   t.Cost - 1,
		Title:  t.Title,
		Status: t.Status,
		Icon:   t.Icon,
		Link:   t.Link,
	}
}

func MapToTasksUser(user TaskCash) TaskCashToService {
	return TaskCashToService{
		ID:            user.ID,
		UID:           user.UID,
		TransactionID: user.TransactionID,
		Status:        user.Status,
		Number:        user.Number,
		Total:         strconv.Itoa(int(user.Total)),
	}
}
