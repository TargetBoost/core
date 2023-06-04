package api

import (
	"github.com/SevereCloud/vksdk/v2/api"
)

type VK struct {
	*api.VK
}

func New(token string) *VK {
	vk := api.NewVK(token)

	return &VK{
		vk,
	}
}

func (v *VK) IsMember(groupID, userID int64) (int, error) {
	member, err := v.GroupsIsMember(api.Params{
		"group_id": groupID,
		"user_id":  userID,
	})
	if err != nil {
		return 0, err
	}
	return member, nil
}

func (v *VK) GetInfo() (api.UsersGetResponse, error) {
	user, err := v.UsersGet(api.Params{})
	if err != nil {
		return nil, err
	}
	return user, nil
}
