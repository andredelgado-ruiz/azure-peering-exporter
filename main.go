package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	resourceGroupFlag  = flag.String("resource-group", "", "Azure Resource Group name")
	vnetNameFlag       = flag.String("vnet-name", "", "Azure Virtual Network name")
	intervalFlag       = flag.Int("interval", 300, "Interval in seconds between checks")
	tenantIDFlag       = flag.String("tenant-id", "", "Azure Tenant ID")
	clientIDFlag       = flag.String("client-id", "", "Azure Client ID")
	clientSecretFlag   = flag.String("client-secret", "", "Azure Client Secret")
	subscriptionIDFlag = flag.String("subscription-id", "", "Azure Subscription ID")

	peeringState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "azure_vnet_peering_state",
			Help: "State of Azure VNet peering.",
		},
		[]string{"resource_group", "vnet_name", "peering_name"},
	)
)

func init() {
	// Register the metric with Prometheus
	prometheus.MustRegister(peeringState)
}

func startHttpServer() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Starting HTTP server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func main() {
	flag.Parse() // Parse the flags from command line

	// Start the HTTP server for health checks and metrics
	go startHttpServer()

	if err := run(); err != nil {
		log.Fatalf("Error running exporter: %v", err)
	}
}

func run() error {
	client, err := createNetworkClient()
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}

	resourceGroup := *resourceGroupFlag
	vnetName := *vnetNameFlag

	if resourceGroup == "" || vnetName == "" {
		return fmt.Errorf("resource group and vnet name must be provided")
	}

	interval := time.Duration(*intervalFlag) * time.Second
	for {
		err = listPeerings(client, resourceGroup, vnetName)
		if err != nil {
			log.Printf("Error listing peerings: %v", err)
		}
		time.Sleep(interval)
	}
}

func createNetworkClient() (*armnetwork.VirtualNetworkPeeringsClient, error) {
	tenantID := *tenantIDFlag
	clientID := *clientIDFlag
	clientSecret := *clientSecretFlag
	subscriptionId := *subscriptionIDFlag

	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		return nil, fmt.Errorf("creating credentials: %w", err)
	}

	client, err := armnetwork.NewVirtualNetworkPeeringsClient(subscriptionId, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("creating network client: %w", err)
	}
	return client, nil
}

func listPeerings(client *armnetwork.VirtualNetworkPeeringsClient, resourceGroup, vnetName string) error {
	pager := client.NewListPager(resourceGroup, vnetName, nil)
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		if err != nil {
			return fmt.Errorf("retrieving next page: %w", err)
		}

		for _, peering := range resp.Value {
			state := 0
			if *peering.Properties.PeeringState == "Connected" {
				state = 1
			}
			// Update the metric with the current state of the peering
			peeringState.WithLabelValues(resourceGroup, vnetName, *peering.Name).Set(float64(state))
		}
	}
	return nil
}
