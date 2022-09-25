package models

func MapToTarget(t Target) TargetService {
	return TargetService{
		UID:    t.UID,
		Title:  t.Title,
		Link:   t.Link,
		Icon:   t.Icon,
		Status: t.Status,
		Count:  t.Count,
		Cost:   t.Cost,
	}
}
