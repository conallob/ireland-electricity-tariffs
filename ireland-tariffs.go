package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
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

// HomeAssistantTariff represents a tariff in the format expected by Home Assistant
type HomeAssistantTariff struct {
	Provider      string             `json:"provider"`
	Plan          string             `json:"plan"`
	EffectiveDate string             `json:"effective_date"`
	Rates         map[string]float64 `json:"rates"`
	RatesExVAT    map[string]float64 `json:"rates_ex_vat"`
	VATRate       float64            `json:"vat_rate"`
	Discount      float64            `json:"discount"`
}

// HomeAssistantOutput represents the structure expected by Home Assistant REST integration
type HomeAssistantOutput struct {
	Tariffs map[string]HomeAssistantTariff `json:"tariffs"` // Tariff shortname -> Tariff data
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

// ProviderTariffs organizes tariffs by provider for easier contribution management
type ProviderTariffs struct {
	Provider string
	Tariffs  []ElectricityTariff
}

// NewTariff creates a new tariff with common parameters
func NewTariff(provider, planName, planShortName string, effectiveDate time.Time, day, peak, night, microGen, vatRate, discount float64) ElectricityTariff {
	return ElectricityTariff{
		Provider:      provider,
		PlanName:      planName,
		PlanShortName: planShortName,
		EffectiveDate: effectiveDate,
		Price: TariffPrice{
			Day:             day,
			Peak:            peak,
			Night:           night,
			MicroGeneration: microGen,
			VATRate:         vatRate,
			Discount:        discount,
		},
	}
}

// filterTariffs filters tariffs based on provider and tariff shortname
func filterTariffs(tariffs []ElectricityTariff, provider, tariffName string) []ElectricityTariff {
	if provider == "" && tariffName == "" {
		return tariffs
	}

	filtered := []ElectricityTariff{}
	for _, tariff := range tariffs {
		matchProvider := provider == "" || strings.EqualFold(tariff.Provider, provider)
		matchTariff := tariffName == "" || strings.EqualFold(tariff.PlanShortName, tariffName)

		if matchProvider && matchTariff {
			filtered = append(filtered, tariff)
		}
	}
	return filtered
}

func main() {
	// Command-line flags for filtering output
	providerFilter := flag.String("provider", "", "Filter by provider name (case-insensitive)")
	tariffFilter := flag.String("tariff", "", "Filter by tariff shortname (case-insensitive)")
	outputFile := flag.String("output", "home_assistant_tariffs.json", "Output file path for Home Assistant JSON")
	skipIndividual := flag.Bool("skip-individual", false, "Skip generating individual tariff JSON files")
	listProviders := flag.Bool("list-providers", false, "List all available providers and exit")
	listTariffs := flag.Bool("list-tariffs", false, "List all available tariffs and exit")

	flag.Parse()

	// ========================================
	// TARIFF DEFINITIONS
	// ========================================
	// To contribute a new provider or tariff, add entries below using the NewTariff helper function
	// Parameters: provider, planName, planShortName, effectiveDate, day, peak, night, microGeneration, vatRate, discount

	// Energia tariffs
	energiaTariffs := []ElectricityTariff{
		NewTariff(
			"Energia",
			"Smart Data - 15",
			"energia-smart-15",
			time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
			0.3451, // day rate
			0.3617, // peak rate
			0.1848, // night rate
			0.20,   // microgeneration rate
			1.09,   // VAT rate (1.09 = 9% VAT)
			0.85,   // discount (0.85 = 15% discount)
		),
	}

	// Electric Ireland tariffs
	electricIrelandTariffs := []ElectricityTariff{
		NewTariff(
			"Electric Ireland",
			"Standard",
			"electric-ireland-standard",
			time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			0.3521, // day rate
			0.3721, // peak rate
			0.1921, // night rate
			0.19,   // microgeneration rate
			1.09,   // VAT rate
			1.00,   // no discount
		),
		NewTariff(
			"Electric Ireland",
			"Smart Dual Fuel",
			"electric-ireland-smart-dual-fuel",
			time.Date(2025, time.August, 10, 0, 0, 0, 0, time.UTC),
			0.23779, // day rate
			0.25368, // peak rate
			0.12495, // night rate
			0.195,   // microgeneration rate
			1.09,    // VAT rate
			0.7,     // 30% discount
		),
	}

	// Bord Gáis Energy tariffs
	bordGaisTariffs := []ElectricityTariff{
		NewTariff(
			"Bord Gáis Energy",
			"Standard",
			"bord-gais-standard",
			time.Date(2025, time.March, 1, 0, 0, 0, 0, time.UTC),
			0.3489, // day rate
			0.3689, // peak rate
			0.1889, // night rate
			0.18,   // microgeneration rate
			1.09,   // VAT rate
			0.90,   // 10% discount
		),
	}

	// Combine all provider tariffs
	// To add a new provider, create a new slice above and append it here
	allTariffs := []ElectricityTariff{}
	allTariffs = append(allTariffs, energiaTariffs...)
	allTariffs = append(allTariffs, electricIrelandTariffs...)
	allTariffs = append(allTariffs, bordGaisTariffs...)

	// Handle list flags
	if *listProviders {
		providers := make(map[string]bool)
		for _, tariff := range allTariffs {
			providers[tariff.Provider] = true
		}
		fmt.Println("Available providers:")
		for provider := range providers {
			fmt.Printf("  - %s\n", provider)
		}
		return
	}

	if *listTariffs {
		fmt.Println("Available tariffs:")
		for _, tariff := range allTariffs {
			fmt.Printf("  - %s (%s - %s)\n", tariff.PlanShortName, tariff.Provider, tariff.PlanName)
		}
		return
	}

	// Apply filters
	tariffs := filterTariffs(allTariffs, *providerFilter, *tariffFilter)

	if len(tariffs) == 0 {
		fmt.Println("No tariffs match the specified filters.")
		if *providerFilter != "" {
			fmt.Printf("  Provider filter: %s\n", *providerFilter)
		}
		if *tariffFilter != "" {
			fmt.Printf("  Tariff filter: %s\n", *tariffFilter)
		}
		fmt.Println("\nUse -list-providers or -list-tariffs to see available options.")
		os.Exit(1)
	}

	fmt.Printf("Generating output for %d tariff(s)...\n", len(tariffs))

	// Create individual tariff files (for backward compatibility)
	if !*skipIndividual {
		fmt.Println("Generating individual tariff files...")
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
		fmt.Printf("Generated individual files for %d tariff(s)\n", len(tariffs))
	}

	// Create Home Assistant compatible output
	fmt.Printf("Generating Home Assistant JSON to: %s\n", *outputFile)
	haOutput := HomeAssistantOutput{
		Tariffs: make(map[string]HomeAssistantTariff),
	}

	// Populate the Home Assistant output structure
	for _, tariff := range tariffs {
		haTariff := HomeAssistantTariff{
			Provider:      tariff.Provider,
			Plan:          tariff.PlanName,
			EffectiveDate: tariff.EffectiveDate.Format("2006-01-02"),
			VATRate:       tariff.Price.VATRate,
			Discount:      tariff.Price.Discount,
			Rates: map[string]float64{
				"day":             ApplyVATAndDiscount(tariff.Price.Day, tariff.Price.VATRate, tariff.Price.Discount),
				"peak":            ApplyVATAndDiscount(tariff.Price.Peak, tariff.Price.VATRate, tariff.Price.Discount),
				"night":           ApplyVATAndDiscount(tariff.Price.Night, tariff.Price.VATRate, tariff.Price.Discount),
				"microgeneration": tariff.Price.MicroGeneration,
			},
			RatesExVAT: map[string]float64{
				"day":             tariff.Price.Day,
				"peak":            tariff.Price.Peak,
				"night":           tariff.Price.Night,
				"microgeneration": tariff.Price.MicroGeneration,
			},
		}

		haOutput.Tariffs[tariff.PlanShortName] = haTariff
	}

	// Write the Home Assistant compatible file
	haPayload, err := json.MarshalIndent(haOutput, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling Home Assistant JSON: %v\n", err)
		return
	}

	e := os.WriteFile(*outputFile, haPayload, 0644)
	if e != nil {
		fmt.Printf("Error writing Home Assistant file: %v\n", e)
		return
	}

	fmt.Printf("\n✓ Successfully generated Home Assistant output: %s\n", *outputFile)
	if !*skipIndividual {
		fmt.Printf("✓ Generated %d individual tariff file(s)\n", len(tariffs)*2)
	}
}
