package airport

import (
	"fmt"
	"log"

	"github.com/asaskevich/govalidator"
)

type Airport struct {
	ID      int    `json:"id" valid:",optional"`
	Name    string `json:"name" valid:"type(string)"`
	City    string `json:"city" valid:"type(string)"`
	Country string `json:"country" valid:"type(string)"`
}

type Repository interface {
	GetAirportByID(airportID string) (*Airport, error)
}

func (p *Airport) Validate() error {
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
	return err
}
