package privilege

type PrivilegeHistory struct {
	ID            int    `json:"id"`
	PrivilegeID   int    `json:"privilegeId"`
	TicketUID     string `json:"ticketUid"`
	Date          string `json:"date"`
	BalanceDiff   int    `json:"balanceDiff"`
	OperationType string `json:"operationType"`
}

type PrivilegeShortInfo struct {
	Status  string `json:"status"`
	Balance int    `json:"balance"`
}

type PrivilegeInfo struct {
	Balance int                 `json:"balance"`
	Status  string              `json:"status"`
	History *[]PrivilegeHistory `json:"history"`
}

type Privilege struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Status   string `json:"status"`
	Balance  int    `json:"balance"`
}
