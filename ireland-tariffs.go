package main

import (
	"encoding/json"
	"fmt"
	"math"
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
	Provider      string      `json:"provider"`
	PlanName      string      `json:"plan"`
	PlanShortName string      `json:"shortname"`
	Price         TariffPrice `json:"price"`
}

func CurrencyRounding(amount float64, precision int) float64 {
	p := math.Pow(10, float64(precision))
	return (amount * p) / p
}

func main() {

	// Map Keys are Effective Date (can be in the future)
	ieTariffs := map[time.Time]*ElectricityTariff{
		time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC): {
			Provider:      "Energia",
			PlanName:      "Smart Data - 15",
			PlanShortName: "energia-smart-15",
			Price: TariffPrice{
				Day:             0.3451,
				Peak:            0.3617,
				Night:           0.1848,
				MicroGeneration: 0.24,
				VATRate:         1.09, // +9% VAT
				Discount:        0.85, // -15% discount
			},
		},
	}

	for _, pricePlan := range ieTariffs {
		// Ex VAT, without discounts
		fileName := fmt.Sprintf("%s.json", pricePlan.PlanShortName)
		payload, err := json.MarshalIndent(pricePlan, "", "  ")
		if err != nil {
			fmt.Errorf("%v", err)
		}

		e := os.WriteFile(fileName, payload, 0644)
		if e != nil {
			fmt.Errorf("%v", e)
		}
		// Inc VAT, with discounts
		pricePlan.Price.Day *= pricePlan.Price.VATRate * pricePlan.Price.Discount
		pricePlan.Price.Day = CurrencyRounding(pricePlan.Price.Day, 4)
		pricePlan.Price.Peak *= pricePlan.Price.VATRate * pricePlan.Price.Discount
		pricePlan.Price.Peak = CurrencyRounding(pricePlan.Price.Peak, 4)
		pricePlan.Price.Night *= pricePlan.Price.VATRate * pricePlan.Price.Discount
		pricePlan.Price.Night = CurrencyRounding(pricePlan.Price.Night, 4)
		fileName = fmt.Sprintf("%s.inc.vat.json", pricePlan.PlanShortName)
		payload, err = json.MarshalIndent(pricePlan, "", "  ")
		if err != nil {
			fmt.Errorf("%v", err)
		}

		e = os.WriteFile(fileName, payload, 0644)
		if e != nil {
			fmt.Errorf("%v", e)
		}
	}

}
