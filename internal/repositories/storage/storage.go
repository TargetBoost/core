package storage

import (
	"core/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewStorageRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetFileByKey(key string) *models.FileStorage {
	var fileStorage models.FileStorage

	r.db.Table("file_storages").Where("key = ?", key).Find(&fileStorage)

	if fileStorage.Key == "" {
		return nil
	}

	return &fileStorage
}

func (r *Repository) SetChatMembers(cid int64, title, userName string) {
	var q models.ChatMembersChanel
	q.CID = cid
	q.Title = title
	q.UserName = userName
	r.db.Table("chat_members_chanels").Create(&q)
}

func (r *Repository) GetStatisticTargetsOnExecutesIsTrue() {
	//	select cmc.c_id, u.tg, u.id, t.link, cc.c_id from users u inner join chat_members_chanels cmc on REPLACE(u.tg, '@', '') = cmc.user_name inner join queues q on u.id = q.uid inner join targets t on q.t_id = t.id right join chat_members_chanels cc on cc.user_name = replace(t.link, 'https://t.me/', '') where q.status = 3 order by u.id

	r.db.Table(
		"users u",
	).Select(
		"cmc.c_id as cid_users, u.tg, u.id, t.link, cc.c_id as cid_channels",
	).Joins(
		"inner join chat_members_chanels cmc on REPLACE(u.tg, '@', '') = cmc.user_name inner join queues q on u.id = q.uid inner join targets t on q.t_id = t.id right join chat_members_chanels cc on cc.user_name = replace(t.link, 'https://t.me/', '')",
	).Where(
		"q.status = 3",
	).Order("u.id")
}
