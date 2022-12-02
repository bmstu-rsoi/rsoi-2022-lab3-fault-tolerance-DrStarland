package user

import (
	"gateway/pkg/models/privilege"
	"gateway/pkg/models/tickets"
)

type UserInfo struct {
	Privilege   *privilege.PrivilegeShortInfo `json:"privilege"`
	TicketsInfo *[]tickets.TicketInfo         `json:"tickets"`
}

// type UserInfoCircuitBreaker struct {
// 	Privilege   string                `json:"privilege"`
// 	TicketsInfo *[]tickets.TicketInfo `json:"tickets"`
// }
