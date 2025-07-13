package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"
)

// TariffPrice represents the pricing structure for an electricity tariff
type TariffPrice struct {
	Day             float64 `json:"day"`             // Price per kWh, in €
	Peak            float64 `json:"peak"`            // Price per kWh, in €
	Night           float64 `json:"night"`           // Price per kWh, in €
	MicroGeneration float64 `json:"microgeneration"` // Price per kWh, in €
	VATRate         float64 `json:"VAT"`             // Applicable VAT Rate, in decimal form
	Discount        float64 `json:"discount"`        // Applied Discount, in decimal form
}

// ElectricityTariff represents a single electricity tariff plan
type ElectricityTariff struct {
	Provider      string      `json:"provider"`
	PlanName      string      `json:"plan"`
	PlanShortName string      `json:"shortname"`
	Price         TariffPrice `json:"price"`
	EffectiveDate time.Time   `json:"effective_date"` // Date when this tariff becomes effective
}

// HomeAssistantOutput represents the structure expected by Home Assistant REST integration
type HomeAssistantOutput struct {
	Providers map[string]map[string]interface{} `json:"providers"` // Provider -> Plan -> Tariff data
}

// CurrencyRounding rounds a currency amount to the specified precision
func CurrencyRounding(amount float64, precision int) float64 {
	p := math.Pow(10, float64(precision))
	return math.Round(amount*p) / p
}

// ApplyVATAndDiscount applies VAT and discount to a price
func ApplyVATAndDiscount(price float64, vatRate float64, discount float64) float64 {
	return CurrencyRounding(price*vatRate*discount, 4)
}

func main() {
	// Define tariffs for different providers
	tariffs := []ElectricityTariff{
		{
			Provider:      "Energia",
			PlanName:      "Smart Data - 15",
			PlanShortName: "energia-smart-15",
			EffectiveDate: time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
			Price: TariffPrice{
				Day:             0.3451,
				Peak:            0.3617,
				Night:           0.1848,
				MicroGeneration: 0.20,
				VATRate:         1.09, // +9% VAT
				Discount:        0.85, // -15% discount
			},
		},
		{
			Provider:      "Electric Ireland",
			PlanName:      "Standard",
			PlanShortName: "electric-ireland-standard",
			EffectiveDate: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			Price: TariffPrice{
				Day:             0.3521,
				Peak:            0.3721,
				Night:           0.1921,
				MicroGeneration: 0.19,
				VATRate:         1.09, // +9% VAT
				Discount:        1.00, // No discount
			},
		},
		{
			Provider:      "Bord Gáis Energy",
			PlanName:      "Standard",
			PlanShortName: "bord-gais-standard",
			EffectiveDate: time.Date(2025, time.March, 1, 0, 0, 0, 0, time.UTC),
			Price: TariffPrice{
				Day:             0.3489,
				Peak:            0.3689,
				Night:           0.1889,
				MicroGeneration: 0.18,
				VATRate:         1.09, // +9% VAT
				Discount:        0.90, // -10% discount
			},
		},
	}

	// Group tariffs by provider
	providerMap := make(map[string][]ElectricityTariff)
	for _, tariff := range tariffs {
		providerMap[tariff.Provider] = append(providerMap[tariff.Provider], tariff)
	}

	// Create individual tariff files (for backward compatibility)
	for _, tariff := range tariffs {
		// Ex VAT, without discounts
		fileName := fmt.Sprintf("%s.json", tariff.PlanShortName)
		payload, err := json.MarshalIndent(tariff, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			continue
		}

		e := os.WriteFile(fileName, payload, 0644)
		if e != nil {
			fmt.Printf("Error writing file: %v\n", e)
			continue
		}

		// Create a copy for the VAT-inclusive version
		tariffWithVAT := tariff
		tariffWithVAT.Price.Day = ApplyVATAndDiscount(tariff.Price.Day, tariff.Price.VATRate, tariff.Price.Discount)
		tariffWithVAT.Price.Peak = ApplyVATAndDiscount(tariff.Price.Peak, tariff.Price.VATRate, tariff.Price.Discount)
		tariffWithVAT.Price.Night = ApplyVATAndDiscount(tariff.Price.Night, tariff.Price.VATRate, tariff.Price.Discount)

		fileName = fmt.Sprintf("%s.inc.vat.json", tariff.PlanShortName)
		payload, err = json.MarshalIndent(tariffWithVAT, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			continue
		}

		e = os.WriteFile(fileName, payload, 0644)
		if e != nil {
			fmt.Printf("Error writing file: %v\n", e)
			continue
		}
	}

	// Create Home Assistant compatible output
	haOutput := HomeAssistantOutput{
		Providers: make(map[string]map[string]interface{}),
	}

	// Populate the Home Assistant output structure
	for provider, providerTariffs := range providerMap {
		haOutput.Providers[provider] = make(map[string]interface{})

		for _, tariff := range providerTariffs {
			// Create a map for this tariff plan
			planMap := map[string]interface{}{
				"provider":       tariff.Provider,
				"plan":           tariff.PlanName,
				"shortname":      tariff.PlanShortName,
				"effective_date": tariff.EffectiveDate.Format("2006-01-02"),
				"price": map[string]interface{}{
					"day":             tariff.Price.Day,
					"peak":            tariff.Price.Peak,
					"night":           tariff.Price.Night,
					"microgeneration": tariff.Price.MicroGeneration,
					"vat_rate":        tariff.Price.VATRate,
					"discount":        tariff.Price.Discount,
				},
				"price_inc_vat": map[string]interface{}{
					"day":             ApplyVATAndDiscount(tariff.Price.Day, tariff.Price.VATRate, tariff.Price.Discount),
					"peak":            ApplyVATAndDiscount(tariff.Price.Peak, tariff.Price.VATRate, tariff.Price.Discount),
					"night":           ApplyVATAndDiscount(tariff.Price.Night, tariff.Price.VATRate, tariff.Price.Discount),
					"microgeneration": tariff.Price.MicroGeneration,
				},
			}

			// Add to the provider's plans
			haOutput.Providers[provider][tariff.PlanShortName] = planMap
		}
	}

	// Write the Home Assistant compatible file
	haPayload, err := json.MarshalIndent(haOutput, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling Home Assistant JSON: %v\n", err)
		return
	}

	e := os.WriteFile("home_assistant_tariffs.json", haPayload, 0644)
	if e != nil {
		fmt.Printf("Error writing Home Assistant file: %v\n", e)
		return
	}

	fmt.Println("Successfully generated tariff files and Home Assistant compatible output.")
}
