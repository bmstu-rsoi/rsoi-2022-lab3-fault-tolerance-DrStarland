package flight

import (
	"fmt"
	"log"

	// валидатор
	"github.com/asaskevich/govalidator"
)

type Flight struct {
	ID            int    `json:"id" valid:",optional"`
	FlightNumber  string `json:"flightNumber" valid:"type(string)"`
	Date          string `json:"date" valid:"type(string)"`
	FromAirportId int    `json:"fromAirportId" valid:"type(int)"`
	ToAirportId   int    `json:"toAirportId" valid:"type(int)"`
	Price         int    `json:"price" valid:"type(int)"`
}
type Repository interface {
	GetAllFlights() ([]*Flight, error)
	GetFlightByNumber(flightNumber string) (*Flight, error)
}

func (p *Flight) Validate() error {
	_, err := govalidator.ValidateStruct(p)
	if err != nil {
		if allErrs, ok := err.(govalidator.Errors); ok {
			for _, fld := range allErrs.Errors() {
				data := []byte(fmt.Sprintf("field: %#v\n\n", fld))
				log.Println(data)
				//w.Write(data)
			}
		}
	}
	return err // mya?
}
