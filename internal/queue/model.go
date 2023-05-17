package queue

type Task struct {
	ID     uint   `json:"id"`
	UID    int64  `json:"uid"`
	Cost   int64  `json:"cost"`
	Title  string `json:"title"`
	Status int64  `json:"status"`
}
