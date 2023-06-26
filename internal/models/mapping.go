package models

import (
	"strconv"
)

func MapToTarget(t Target) TargetService {
	return TargetService{
		NameCompany:        t.NameCompany,
		DescriptionCompany: t.DescriptionCompany,
		Type:               t.Type,
		Link:               t.Link,
		Limit:              t.Limit,
		//TypeAd: struct {
		//	gorm.Model
		//	Value string `json:"value"`
		//	Label string `json:"label"`
		//	Color string `json:"color"`
		//}{
		//	Model: gorm.Model{
		//		ID:        t.TypeAd.ID,
		//		CreatedAt: time.Time{},
		//		UpdatedAt: time.Time{},
		//		DeletedAt: gorm.DeletedAt{
		//			Time:  time.Time{},
		//			Valid: false,
		//		},
		//	},
		//	Value: t.TypeAd.Value,
		//	Label: t.TypeAd.Label,
		//	Color: t.TypeAd.Color,
		//},
	}
}

func MapToTargetAdmin(t TargetToAdmin) TargetService {
	return TargetService{
		NameCompany:        t.NameCompany,
		DescriptionCompany: t.DescriptionCompany,
		Type:               t.Type,
		Link:               t.Link,
		Limit:              t.Limit,
		//TypeAd: struct {
		//	gorm.Model
		//	Value string `json:"value"`
		//	Label string `json:"label"`
		//	Color string `json:"color"`
		//}{
		//	Model: gorm.Model{
		//		ID:        t.TypeAd.ID,
		//		CreatedAt: time.Time{},
		//		UpdatedAt: time.Time{},
		//		DeletedAt: gorm.DeletedAt{
		//			Time:  time.Time{},
		//			Valid: false,
		//		},
		//	},
		//	Value: t.TypeAd.Value,
		//	Label: t.TypeAd.Label,
		//	Color: t.TypeAd.Color,
		//},
	}
}

func MapToQueueExecutors(t QueueToExecutors) QueueToService {
	return QueueToService{
		ID:        t.ID,
		TID:       t.TID,
		UID:       t.UID,
		Cost:      t.Cost - 1,
		Title:     t.Title,
		Status:    t.Status,
		Icon:      t.Icon,
		Link:      t.Link,
		PhotoLink: t.PhotoLink,
		Bio:       t.Bio,
		CountSub:  t.CountSub,
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
