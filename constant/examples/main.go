package main

import (
	"log"

	"github.com/muhammad-fakhri/go-libs/constant"
)

func main() {
	// valid country
	validCountry := constant.Country("ID")
	log.Printf("country %s, timezone %s, dst %v", validCountry, validCountry.TimeZone(), validCountry.IsDST())

	// valid country with DST
	validCountryDST := constant.Country("CL")
	log.Printf("country %s, timezone %s, dst %v", validCountryDST, validCountryDST.TimeZone(), validCountryDST.IsDST())

	// invalid country
	invalidCountry := constant.Country("IX")
	log.Printf("country %s, timezone %s, dst %v", invalidCountry, invalidCountry.TimeZone(), invalidCountry.IsDST())
}
