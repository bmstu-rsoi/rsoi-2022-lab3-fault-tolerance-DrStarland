package flights

type Flight struct {
	ID            int    `json:"id"`
	FlightNumber  string `json:"flightNumber"`
	Date          string `json:"date"`
	FromAirportId int    `json:"fromAirportId"`
	ToAirportId   int    `json:"toAirportId"`
	Price         int    `json:"price"`
}

type FlightInfo struct {
	FlightNumber string `json:"flightNumber"`
	Date         string `json:"date"`
	FromAirport  string `json:"fromAirport"`
	ToAirport    string `json:"toAirport"`
	Price        int    `json:"price"`
}

type FlightsLimited struct {
	Page          int           `json:"page"`
	PageSize      int           `json:"pageSize"`
	TotalElements int           `json:"totalElements"`
	Items         *[]FlightInfo `json:"items"`
}
