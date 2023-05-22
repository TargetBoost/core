package models

type TargetService struct {
	// основные данные
	UID    uint    `json:"uid"`    // кто создал задачу
	Title  string  `json:"title"`  // заголовок
	Link   string  `json:"link"`   // ссылка на задание
	Icon   string  `json:"icon"`   // иконка задания
	Status int64   `json:"status"` // открыта/закрыта
	Count  int64   `json:"count"`  // количетсво заданий
	Total  float64 `json:"total"`
	Cost   float64 `json:"cost"` // цена одного задания

	// гео данные
	Country string `json:"country"`  // список стран исполнителей
	City    string `json:"city"`     // список городов исполнителей
	OldFrom int64  `json:"old_from"` // возраст исполнителя от
	OldTo   int64  `json:"old_to"`   // возраст исполнителя до
	Gender  string `json:"gender"`   // половой признак

	Type  string `json:"type"`
	Cause string `json:"cause"`
}

type TaskToService struct {
	UID uint `json:"uid"`
	TID uint `json:"tid"`
}
