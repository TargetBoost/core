package models

func MapToTarget(t Target) TargetService {
	return TargetService{
		UID:    t.UID,
		Title:  t.Title,
		Link:   t.Link,
		Icon:   t.Icon,
		Status: t.Status,
		Count:  t.Count,
		Total:  t.Total,
		Cost:   t.Cost,
		Cause:  t.Cause,
	}
}

func MapToTargetExecutors(t TargetToExecutors) TargetServiceToExecutors {
	return TargetServiceToExecutors{
		UID: t.UID,
		TID: t.TID,
	}
}
