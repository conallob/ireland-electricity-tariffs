package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type TariffPrice struct {
	Day             float64 `json:"day"`             // Price per kWh, in €
	Peak            float64 `json:"peak"`            // Price per kWh, in €
	Night           float64 `json:"night"`           // Price per kWh, in €
	MicroGeneration float64 `json:"microgeneration"` // Price per kWh, in €
	VATRate         float64 `json:"VAT"`             // Applicable VAT Rate, in decimal form
	Discount        float64 `json:"discount"`        // Applied Discount, in decimal form
}

type ElectricityTariff struct {
	Provider      string `json:"provider"`
	PlanName      string `json:"plan"`
	PlanShortName string `json:"shortname"`
	Price         TariffPrice `json:"price"`
}

func main() {

	// Map Keys are Effective Date (can be in the future)
	ieTariffs := map[time.Time]*ElectricityTariff{
		time.Date(2023, time.December, 22, 0, 0, 0, 0, time.UTC): {
			Provider:      "Energia",
			PlanName:      "Smart Data - 15",
			PlanShortName: "energia-smart-15",
			Price: TariffPrice{
				Day:             0.3731,
				Peak:            0.391,
				Night:           0.1998,
				MicroGeneration: 0.24,
				VATRate:         0.09,
				Discount:        0.15,
			},
		},		
	}

	for _, pricePlan := range ieTariffs {
		fileName := fmt.Sprintf("%s.json", pricePlan.PlanShortName)
		payload, err := json.MarshalIndent(pricePlan, "", "  ")
		if err != nil {
		    fmt.Errorf("%v", err)
		}

		e := os.WriteFile(fileName, payload, 0644)
		if e != nil {
		    fmt.Errorf("%v", e)
		}
	}

}
