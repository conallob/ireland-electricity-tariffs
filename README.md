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
- **Bord Gáis Energy**
    - Standard

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

To use these tariffs with Home Assistant's REST sensor:

1. Add the following to your `configuration.yaml`:

```yaml
sensor:
  - platform: rest
    name: ireland_electricity_tariffs
    resource: https://github.com/conallob/ireland-electricity-tariffs/releases/latest/download/home_assistant_tariffs.json
    value_template: "{{ value_json }}"
    scan_interval: 86400  # Update once per day
```

2. Create template sensors for specific tariffs:

```yaml
template:
  - sensor:
      - name: "Energia Smart Data Day Rate"
        unit_of_measurement: "€/kWh"
        state: "{{ state_attr('sensor.ireland_electricity_tariffs', 'providers')['Energia']['energia-smart-15']['price_inc_vat']['day'] }}"

      - name: "Energia Smart Data Night Rate"
        unit_of_measurement: "€/kWh"
        state: "{{ state_attr('sensor.ireland_electricity_tariffs', 'providers')['Energia']['energia-smart-15']['price_inc_vat']['night'] }}"
```

## File Structure

- Individual tariff files (e.g., `energia-smart-15.json`) contain the raw tariff
  data without VAT or discounts applied
- VAT-inclusive files (e.g., `energia-smart-15.inc.vat.json`) contain the tariff
  data with VAT and discounts applied
- `home_assistant_tariffs.json` contains all tariffs organized by provider in a
  format compatible with Home Assistant

## Contributing

To contribute new tariff data:

1. Fork this repository
2. Add your tariff data to `ireland-tariffs.go`
3. Submit a pull request

## License

See the [LICENSE](LICENSE) file for details.
