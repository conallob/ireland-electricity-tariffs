package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type TariffPrice struct {
	Day      	float64 `json:"day"`      		// Price per kWh, in €
	Peak     	float64 `json:"peak"`     		// Price per kWh, in €
	Night    	float64 `json:"night"`    		// Price per kWh, in €
	MicroGeneration	float64 `json:"microgeneration"`    	// Price per kWh, in €
	VATRate  	float64 `json:"VAT"`      		// Applicable VAT Rate, in decimal form
	Discount 	float64 `json:"discount"` 		// Applied Discount, in decimal form
}

type ElectricityTariff struct {
	Provider         string    `json:"provider"`
	PlanName         string    `json:"plan"`
	PlanShortName    string `json:"shortname"`
	Price            TariffPrice
}

func main() {

        // Map Keys are Effective Date (can be in the future)
	ieTariffs := map[time.Time]*ElectricityTariff{
		time.Date(2023, time.June, 13, 0, 0, 0, 0, time.UTC): {
			Provider: "Energia",
			PlanName: "Smart Data - 15",
			PlanShortName: "energia-smart-15",
			Price: TariffPrice{
				Day:		0.4576,
				Peak:     	0.4794,
				Night:		0.2450,
				MicroGeneration: 0.18,
				VATRate:  	0.09,
				Discount: 	0.15,
			},
		},
		time.Date(2023, time.September, 15, 0, 0, 0, 0, time.UTC): {
			Provider: "Energia",
			PlanName: "Smart Data - 15",
			PlanShortName: "energia-smart-15",
			Price: TariffPrice{
				Day:		0.4576,
				Peak:     	0.4794,
				Night:		0.2450,
				MicroGeneration: 0.24,
				VATRate:  	0.09,
				Discount: 	0.15,
			},
		},
			time.Date(2023, time.October, 2, 0, 0, 0, 0, time.UTC): {
			Provider: "Energia",
			PlanName: "Smart Data - 15",
			PlanShortName: "energia-smart-15",
			Price: TariffPrice{
				Day:	  		0.3731,
				Peak:	  		0.213,
				Night:	  		0.245,
				MicroGeneration:	0.24,
				VATRate:  		0.09,
				Discount: 		0.15,
			},
		},

	}

	fmt.Println(json.MarshalIndent(ieTariffs, "", "  "))

}
