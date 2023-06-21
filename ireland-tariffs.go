package main

import (
	"fmt"
	"encoding/json"
	"time"
)

type TariffPrice struct {
	Day float64       // Price per kWh, in €
	Peak float64      // Price per kWh, in €
	Night float64     // Price per kWh, in €
  VATRate float64   // Applicable VAT Rate, in decimal form
  Discount float64  // Applicable Discount Rate, in decimal form
}


type ElectricityTariff struct {
  Provider string
  PlanName string
	Date time.Time			// Effective Date (can be in the future)
  Price TariffPrice
}


func main() {

	ieTariffs := map[string]*ElectricityTariff {
		"energia-smart-15": {
			Provider: "Energia",
			PlanName: "Smart Data - 15",
			Date: time.Date(2023, time.June, 13, 0, 0, 0, 0, time.UTC),
			Price: TariffPrice{
				Day: 0.4576,
				Peak: 0.4794,
				Night: 0.2450,
				VATRate: 0.09,
				Discount:	0.15,
			},
		},
	}

  fmt.Println(json.MarshalIndent(ieTariffs, "", "  "))

}
