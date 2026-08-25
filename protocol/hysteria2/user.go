package hysteria2

import (
	"github.com/sagernet/sing-box/option"
)

func (h *Inbound) UpdateUsers(users []option.Hysteria2User) error {
	names := make([]string, len(users))
	passwords := make([]string, len(users))
	for i, user := range users {
		names[i] = user.Name
		passwords[i] = user.Password
	}
	ids := h.users.reset(names)
	h.service.UpdateUsers(ids, passwords)
	return nil
}
