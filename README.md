# Ireland Electricity Tariffs

A crowd-sourced dataset of electricity tariffs in Ireland, for use in scripting
and home automation such as Home Assistant.

## Overview

This repository provides electricity tariff data for various Irish energy
providers. The data is structured to be easily consumed by home automation
systems, particularly Home Assistant.

## Available Tariffs

The following electricity providers and tariffs are currently supported:

- **Energia**
    - Smart Data - 15
- **Electric Ireland**
    - Standard
    - Smart Dual Fuel
- **Bord Gáis Energy**
    - Standard

To see the full list of available providers and tariffs, run:

```bash
go run ireland-tariffs.go -list-providers
go run ireland-tariffs.go -list-tariffs
```

## Accessing the Tariff Data

The tariff data is automatically generated via GitHub Actions and published to
the [releases page](https://github.com/conallob/ireland-electricity-tariffs/releases/).
You can access the data in several ways:

### Option 1: Direct Download

You can download the JSON files directly from the latest release:

1. Visit
   the [releases page](https://github.com/conallob/ireland-electricity-tariffs/releases/)
2. Download the desired JSON file(s):
    - Individual tariff files (e.g., `energia-smart-15.json`)
    - VAT-inclusive tariff files (e.g., `energia-smart-15.inc.vat.json`)
    - Home Assistant compatible file (`home_assistant_tariffs.json`)

### Option 2: Using curl

To fetch the latest tariff data programmatically:

```bash
# Download the Home Assistant compatible file
curl -L https://github.com/conallob/ireland-electricity-tariffs/releases/latest/download/home_assistant_tariffs.json -o home_assistant_tariffs.json

# Download a specific tariff file
curl -L https://github.com/conallob/ireland-electricity-tariffs/releases/latest/download/energia-smart-15.json -o energia-smart-15.json

# Download a VAT-inclusive tariff file
curl -L https://github.com/conallob/ireland-electricity-tariffs/releases/latest/download/energia-smart-15.inc.vat.json -o energia-smart-15.inc.vat.json
```

### Option 3: Home Assistant Integration

To use these tariffs with Home Assistant's REST sensor and utility meter:

1. Add the following to your `configuration.yaml`:

```yaml
rest:
  - resource: "https://github.com/conallob/ireland-electricity-tariffs/releases/latest/download/home_assistant_tariffs.json"
    scan_interval: 86400  # Update once per day
    sensor:
      - name: "Ireland Electricity Tariffs"
        value_template: "OK"
        json_attributes_path: "$.tariffs"
        json_attributes:
          - "energia-smart-15"
          - "electric-ireland-standard"
          - "electric-ireland-smart-dual-fuel"
          - "bord-gais-standard"
```

2. Create template sensors for specific tariffs:

```yaml
template:
  - sensor:
      - name: "Energia Smart 15 Day Rate"
        unit_of_measurement: "EUR/kWh"
        state: "{{ state_attr('sensor.ireland_electricity_tariffs', 'energia-smart-15')['rates']['day'] }}"

      - name: "Energia Smart 15 Night Rate"
        unit_of_measurement: "EUR/kWh"
        state: "{{ state_attr('sensor.ireland_electricity_tariffs', 'energia-smart-15')['rates']['night'] }}"
```

See `home_assistant_example.yaml` for a complete example including utility_meter integration and tariff automation.

## File Structure

- Individual tariff files (e.g., `energia-smart-15.json`) contain the raw tariff
  data without VAT or discounts applied
- VAT-inclusive files (e.g., `energia-smart-15.inc.vat.json`) contain the tariff
  data with VAT and discounts applied
- `home_assistant_tariffs.json` contains all tariffs organized by shortname in a
  format compatible with Home Assistant
- `home_assistant_example.yaml` contains a complete Home Assistant configuration example

## Command-Line Usage

The `ireland-tariffs.go` program supports several command-line options:

```bash
# Generate JSON for all tariffs (default)
go run ireland-tariffs.go

# List available providers
go run ireland-tariffs.go -list-providers

# List available tariffs
go run ireland-tariffs.go -list-tariffs

# Generate JSON for a specific provider
go run ireland-tariffs.go -provider "Energia" -output energia_tariffs.json

# Generate JSON for a specific tariff
go run ireland-tariffs.go -tariff "energia-smart-15" -output energia_smart_15.json

# Skip generating individual tariff files
go run ireland-tariffs.go -skip-individual

# Combine filters
go run ireland-tariffs.go -provider "Electric Ireland" -output ei_tariffs.json -skip-individual
```

### Available Flags

- `-provider` - Filter by provider name (case-insensitive)
- `-tariff` - Filter by tariff shortname (case-insensitive)
- `-output` - Output file path for Home Assistant JSON (default: `home_assistant_tariffs.json`)
- `-skip-individual` - Skip generating individual tariff JSON files
- `-list-providers` - List all available providers and exit
- `-list-tariffs` - List all available tariffs and exit

## Contributing

We welcome contributions! Here's how you can help:

### Adding a New Tariff

1. **Fork this repository**

2. **Add your tariff data to `ireland-tariffs.go`**

   Find the appropriate provider section (or create a new one) and add your tariff using the `NewTariff` helper function:

   ```go
   // For an existing provider (e.g., Energia)
   energiaTariffs := []ElectricityTariff{
       // ... existing tariffs ...
       NewTariff(
           "Energia",
           "Your Plan Name",
           "energia-plan-shortname",
           time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),  // Effective date
           0.35,    // day rate (EUR/kWh, ex VAT)
           0.37,    // peak rate (EUR/kWh, ex VAT)
           0.19,    // night rate (EUR/kWh, ex VAT)
           0.20,    // microgeneration rate (EUR/kWh)
           1.09,    // VAT rate (1.09 = 9% VAT)
           1.00,    // discount (1.00 = no discount, 0.85 = 15% discount)
       ),
   }
   ```

3. **For a new provider**, create a new section:

   ```go
   // New Provider tariffs
   newProviderTariffs := []ElectricityTariff{
       NewTariff(
           "New Provider",
           "Plan Name",
           "new-provider-plan-shortname",
           time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
           0.35, 0.37, 0.19, 0.20, 1.09, 1.00,
       ),
   }

   // Add to allTariffs
   allTariffs = append(allTariffs, newProviderTariffs...)
   ```

4. **Test your changes**

   ```bash
   # Generate and verify the output
   go run ireland-tariffs.go -list-tariffs
   go run ireland-tariffs.go -provider "Your Provider"
   ```

5. **Submit a pull request** with:
   - Clear description of the tariff being added
   - Source of the tariff information (e.g., provider website URL)
   - Effective date of the tariff

### Tariff Data Guidelines

- **Rates should be in EUR per kWh, excluding VAT**
- **VAT rate**: Use 1.09 for 9% VAT (Ireland's current rate)
- **Discount**: Use decimal format (1.00 = no discount, 0.85 = 15% discount, 0.70 = 30% discount)
- **Effective date**: Use the date when the tariff becomes active
- **Shortname**: Use lowercase with hyphens (e.g., `provider-plan-name`)
- **Always verify rates** from official provider sources before submitting

### Updating Existing Tariffs

When a provider updates their rates:

1. Update the effective date to the new date
2. Update the rate values
3. Add a comment explaining why the tariff was updated (e.g., "Updated for Q1 2025 rates")
4. Submit a pull request with a link to the provider's announcement

### Code Quality

- Run `go fmt` before submitting
- Ensure all tests pass: `go run ireland-tariffs.go -list-tariffs`
- Keep the code organized by provider
- Add inline comments for complex calculations

## License

See the [LICENSE](LICENSE) file for details.
