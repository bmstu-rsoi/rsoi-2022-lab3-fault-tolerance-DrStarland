package ticket

import (
	"fmt"
	"log"

	// валидатор
	"github.com/asaskevich/govalidator"
)

type Ticket struct {
	ID           int    `json:"id"`
	TicketUID    string `json:"ticketUid" valid:"type(string),uuid"`
	Username     string `json:"username" valid:"type(string)"`
	FlightNumber string `json:"flightNumber" valid:"type(string)"`
	Price        int    `json:"price" valid:"type(int),range(0|100000000)"`
	Status       string `json:"status" valid:"type(string)"`
}

func (p *Ticket) Validate() error {
	_, err := govalidator.ValidateStruct(p)
	if err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			for _, fld := range allErrs.Errors() {
				data := []byte(fmt.Sprintf("field: %#v\n\n", fld))
				log.Println(data)
			}
		}
	}
	return err
}

type Repository interface {
	GetByUsername(flightNumber string) ([]*Ticket, error)
	Add(*Ticket) error
	Delete(ticketUID string) error
}
