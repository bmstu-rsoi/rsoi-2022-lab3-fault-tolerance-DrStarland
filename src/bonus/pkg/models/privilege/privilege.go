package privilege

type Privilege struct {
	ID       int    `json:"id"`
	Username string `json:"username" valid:"type(string)"`
	Status   string `json:"status" valid:"type(string)"`
	Balance  int    `json:"balance" valid:"type(int),range(0|100000000)"`
}

type PrivilegeHistory struct {
	ID            int    `json:"id"`
	PrivilegeID   int    `json:"privilegeId"`
	TicketUID     string `json:"ticketUid"`
	Date          string `json:"date"`
	BalanceDiff   int    `json:"balanceDiff"`
	OperationType string `json:"operationType"`
}

type Repository interface {
	GetPrivilegeByUsername(username string) (*Privilege, error)
	GetHistoryById(ticketUID string) ([]*PrivilegeHistory, error)
	CreateHistoryRecord(*PrivilegeHistory) error
	CreatePrivilege(*Privilege) error
	UpdatePrivilege(*Privilege) error
}
