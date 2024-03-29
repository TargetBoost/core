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

func (r *Repository) SetChatMembers(cid, countSub int64, title, userName, photoLink, bio string) {
	var q models.ChatMembersChanel
	q.CID = cid
	q.Title = title
	q.UserName = userName
	q.PhotoLink = photoLink
	q.Bio = bio
	q.CountSub = countSub

	if err := r.db.Debug().Table("chat_members_chanels").Where("c_id = ?", cid).Update("photo_link", photoLink).Update("bio", bio).Update("count_sub", countSub).Error; err != nil {
		r.db.Table("chat_members_chanels").Create(&q) // create new record from newUser
	}
}

func (r *Repository) GetStatisticTargetsOnExecutesIsTrue() []models.StatisticTargetsOnExecutesIsTrue {
	//	select cmc.c_id, u.tg, u.id, t.link, cc.c_id from users u inner join chat_members_chanels cmc on REPLACE(u.tg, '@', '') = cmc.user_name inner join queues q on u.id = q.uid inner join targets t on q.t_id = t.id right join chat_members_chanels cc on cc.user_name = replace(t.link, 'https://t.me/', '') where q.status = 3 order by u.id
	var q []models.StatisticTargetsOnExecutesIsTrue
	r.db.Table(
		"users",
	).Select(
		"chat_members_chanels.c_id as cid_users, users.tg, users.id, targets.link, cl.c_id as cid_channels, chat_members_chanels.updated_at",
	).Joins(
		"inner join chat_members_chanels on REPLACE(users.tg, '@', '') = LOWER(chat_members_chanels.user_name) inner join queues on users.id = queues.uid inner join targets on queues.t_id = targets.id right join chat_members_chanels cl on cl.user_name = replace(targets.link, 'https://t.me/', '')",
	).Where(
		"queues.status = 3",
	).Order("users.id").Find(&q)

	return q
}
