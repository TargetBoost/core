package queue

type Task struct {
	TID    uint   `json:"tid"`
	UID    int64  `json:"uid"`
	Cost   int64  `json:"cost"`
	Title  string `json:"title"`
	Status int64  `json:"status"`
}
