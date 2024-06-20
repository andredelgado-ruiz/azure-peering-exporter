# Azure Peering Exporter for Prometheus

The Azure Peering Exporter is a Prometheus exporter designed to fetch and expose peering state metrics from Azure Virtual Networks. It supports dynamic configuration to monitor resources across multiple environments.

## Metrics Supported

- `azure_peering_state`: Displays the current state of each peering connection within a specified virtual network.

Azure Virtual Network documentation: [Azure Virtual Network Documentation](https://docs.microsoft.com/en-us/azure/virtual-network/)

## TODO

- [x] Implement basic peering state metrics.
- [ ] Add more detailed metrics regarding peering statuses.
- [ ] Implement error handling and retries for API requests.
- [ ] Introduce security enhancements for sensitive data handling.

## Configuration

The exporter accepts configuration through command-line arguments, allowing for easy deployment and configuration changes:

```sh
Usage:
  azure_exporter

Application Options:
  --tenant-id string
        Azure Tenant ID. (required)
  --client-id string
        Azure Client ID. (required)
  --client-secret string
        Azure Client Secret. (required)
  --subscription-id string
        Azure Subscription ID. (required)
  --resource-group string
        Azure Resource Group name. (required)
  --vnet-name string
        Azure Virtual Network name. (required)
  --interval int
        Time in seconds between each metric fetch cycle. (default -> 300)
```


## Usage

### Option A) Running Natively

Ensure you have Go installed, then build and run the exporter:

```sh
user@host: go build -o azure_exporter
user@host: ./azure_exporter --tenant-id xxxx --client-id xxxx --client-secret xxxx --subscription-id xxxx --resource-group xxxx --vnet-name xxxx --interval 300
```

### Option B) Docker

Build the Docker image and run it with the necessary arguments:

```sh
docker build -t azure_peering_exporter .
docker run --rm -it -p 8080:8080 azure_peering_exporter --tenant-id xxxx --client-id xxxx --client-secret xxxx --subscription-id xxxx --resource-group xxxx --vnet-name xxxx --interval 300
```

### Metrics

Examples of the metrics exposed by this exporter:

```sh
# HELP azure_peering_state Current state of the Azure VNet peering.
# TYPE azure_peering_state gauge
azure_peering_state{peering_name="example-peering", resource_group="example-rg", vnet_name="example-vnet"} 1
```

### Contribute

Contributions are welcome! Feel free to open an issue or pull request if you have suggestions or improvements.
