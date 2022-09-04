package constant

import (
	"fmt"
)

type Country string

const (
	ID Country = "ID" //Indonesia
	MY Country = "MY" //Malaysia
	PH Country = "PH" //Philipines
	SG Country = "SG" //Vietnam
	TW Country = "TW" //Taiwan
	TH Country = "TH" //Thailand
	VN Country = "VN" //Vietnam
	BR Country = "BR" //Brazil
	MX Country = "MX" //Mexico
	CO Country = "CO" //Colombia
	CL Country = "CL" //Chile
	ES Country = "ES" //Spain
	FR Country = "FR" //France
	AR Country = "AR" //Argentina
	PL Country = "PL" //Poland
)

var (
	countryTimezone = map[Country]string{
		ID: "Asia/Jakarta",
		SG: "Asia/Singapore",
		TH: "Asia/Bangkok",
		MY: "Asia/Kuala_Lumpur",
		PH: "Asia/Manila",
		VN: "Asia/Ho_Chi_Minh",
		TW: "Asia/Taipei",
		BR: "America/Sao_Paulo",
		MX: "America/Mexico_City",
		CO: "America/Bogota",
		CL: "America/Santiago",
		AR: "America/Argentina/Buenos_Aires",
		PL: "Europe/Warsaw",
		ES: "Europe/Madrid",
		FR: "Europe/Paris",
	}

	countryDST = []Country{MX, CL, ES, FR, PL}

	countryDomain = map[Country]string{
		ID: ".co.id",
		SG: ".sg",
		MY: ".com.my",
		VN: ".vn",
		PH: ".ph",
		TH: ".co.th",
		TW: ".tw",
		BR: ".com.br",
		MX: ".com.mx",
		CO: ".com.co",
		AR: ".com.ar",
		PL: ".pl",
		CL: ".cl",
		ES: ".es",
		FR: ".fr",
	}

	countryCoinMultiplier = map[Country]float64{
		"ID": 1,
		"SG": 1,
		"MY": 1,
		"VN": 0.01,
		"PH": 100,
		"TH": 100,
		"TW": 100,
		"BR": 1,
		"MX": 1,
		"CO": 1,
		"ES": 1,
		"CL": 1,
		"AR": 1,
		"PL": 1,
	}
)

func (c Country) Validate() error {
	switch c {
	case ID, MY, PH, SG, TW, TH, VN, BR, MX, CO, CL, ES, FR, AR, PL:
		return nil
	}
	return fmt.Errorf("invalid_country %s", c)
}

func (c Country) TimeZone() (t string) {
	t = countryTimezone[c]
	if len(t) == 0 {
		t = "UTC"
	}

	return
}

func (c Country) IsDST() bool {
	for _, id := range countryDST {
		if id == c {
			return true
		}
	}

	return false
}

func (c Country) Domain() string {
	return countryDomain[c]
}

func (c Country) CoinMultiplier() (m float64) {
	m = countryCoinMultiplier[c]
	if m == 0 {
		//undefined, set as 1
		m = 1
	}

	return
}
