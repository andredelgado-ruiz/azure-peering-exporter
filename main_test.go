package main

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func series(t *testing.T, name string) map[string]float64 {
	t.Helper()
	families, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]float64{}
	for _, f := range families {
		if f.GetName() != name {
			continue
		}
		for _, m := range f.GetMetric() {
			key := ""
			for _, l := range m.GetLabel() {
				key += l.GetName() + "=" + l.GetValue() + ","
			}
			out[key] = m.GetGauge().GetValue()
		}
	}
	return out
}

func TestPublishDropsPreviousSyncStatus(t *testing.T) {
	publish([]peeringSample{{"rg", "vnet", "peer", 1, "LocalNotInSync", 0}})
	publish([]peeringSample{{"rg", "vnet", "peer", 1, "FullyInSync", 1}})

	got := series(t, "azure_vnet_peering_sync_level")
	want := "peering_name=peer,resource_group=rg,sync_status=FullyInSync,vnet_name=vnet,"
	if len(got) != 1 || got[want] != 1 {
		t.Fatalf("expected only %q = 1, got %v", want, got)
	}
}

func TestPublishDropsDeletedPeering(t *testing.T) {
	publish([]peeringSample{
		{"rg", "vnet", "kept", 1, "FullyInSync", 1},
		{"rg", "vnet", "deleted", 1, "FullyInSync", 1},
	})
	publish([]peeringSample{{"rg", "vnet", "kept", 1, "FullyInSync", 1}})

	if got := series(t, "azure_vnet_peering_state"); len(got) != 1 {
		t.Fatalf("expected 1 azure_vnet_peering_state series, got %v", got)
	}
}
