package tuic

import (
	"github.com/sagernet/sing-box/option"
	exceptions "github.com/sagernet/sing/common/exceptions"
	"github.com/gofrs/uuid/v5"
)

func (h *Inbound) UpdateUsers(users []option.TUICUser) error {
	var names []string
	var userUUIDList [][16]byte
	var userPasswordList []string
	for index, user := range users {
		if user.UUID == "" {
			return exceptions.New("missing uuid for user ", index)
		}
		userUUID, err := uuid.FromString(user.UUID)
		if err != nil {
			return exceptions.Cause(err, "invalid uuid for user ", index)
		}
		names = append(names, user.Name)
		userUUIDList = append(userUUIDList, userUUID)
		userPasswordList = append(userPasswordList, user.Password)
	}
	ids := h.users.reset(names)
	h.server.UpdateUsers(ids, userUUIDList, userPasswordList)
	return nil
}
