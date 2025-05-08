package service

type apiAgeResponse struct {
	Age int `json:"age"`
}

type apiGenderResponse struct {
	Gender string `json:"gender"`
}

type apiNationalityResponse struct {
	Country []nationalizeCountry `json:"country"`
}

type nationalizeCountry struct{
	CountryID   string  `json:"country_id"`
}

type EnrichmentData struct {
    Type  string     
    Value interface{} 
}