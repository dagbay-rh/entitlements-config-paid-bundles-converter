package main

import (
	"reflect"
	"testing"
)

func TestRemovePaidSkusFromSkus(t *testing.T) {
	tests := []struct {
		name     string
		input    []Bundle
		expected []Bundle
	}{
		{
			name: "removes paid SKUs from regular SKUs",
			input: []Bundle{
				{
					Name:     "test-bundle",
					Skus:     []string{"sku1", "sku2", "sku3"},
					PaidSkus: []string{"sku2"},
				},
			},
			expected: []Bundle{
				{
					Name:     "test-bundle",
					Skus:     []string{"sku1", "sku3"},
					PaidSkus: []string{"sku2"},
				},
			},
		},
		{
			name: "removes multiple paid SKUs",
			input: []Bundle{
				{
					Name:     "test-bundle",
					Skus:     []string{"sku1", "sku2", "sku3", "sku4"},
					PaidSkus: []string{"sku2", "sku4"},
				},
			},
			expected: []Bundle{
				{
					Name:     "test-bundle",
					Skus:     []string{"sku1", "sku3"},
					PaidSkus: []string{"sku2", "sku4"},
				},
			},
		},
		{
			name: "handles no overlap between SKUs and paid SKUs",
			input: []Bundle{
				{
					Name:     "test-bundle",
					Skus:     []string{"sku1", "sku2"},
					PaidSkus: []string{"sku3", "sku4"},
				},
			},
			expected: []Bundle{
				{
					Name:     "test-bundle",
					Skus:     []string{"sku1", "sku2"},
					PaidSkus: []string{"sku3", "sku4"},
				},
			},
		},
		{
			name: "handles empty paid SKUs",
			input: []Bundle{
				{
					Name:     "test-bundle",
					Skus:     []string{"sku1", "sku2"},
					PaidSkus: []string{},
				},
			},
			expected: []Bundle{
				{
					Name:     "test-bundle",
					Skus:     []string{"sku1", "sku2"},
					PaidSkus: []string{},
				},
			},
		},
		{
			name: "handles all SKUs being paid",
			input: []Bundle{
				{
					Name:     "test-bundle",
					Skus:     []string{"sku1", "sku2"},
					PaidSkus: []string{"sku1", "sku2"},
				},
			},
			expected: []Bundle{
				{
					Name:     "test-bundle",
					Skus:     []string{},
					PaidSkus: []string{"sku1", "sku2"},
				},
			},
		},
		{
			name: "preserves boolean fields",
			input: []Bundle{
				{
					Name:           "test-bundle",
					UseValidAccNum: true,
					UseValidOrgId:  true,
					UseIsInternal:  false,
					Skus:           []string{"sku1", "sku2", "sku3"},
					PaidSkus:       []string{"sku2"},
				},
			},
			expected: []Bundle{
				{
					Name:           "test-bundle",
					UseValidAccNum: true,
					UseValidOrgId:  true,
					UseIsInternal:  false,
					Skus:           []string{"sku1", "sku3"},
					PaidSkus:       []string{"sku2"},
				},
			},
		},
		{
			name: "handles multiple bundles",
			input: []Bundle{
				{
					Name:     "bundle1",
					Skus:     []string{"sku1", "sku2"},
					PaidSkus: []string{"sku1"},
				},
				{
					Name:     "bundle2",
					Skus:     []string{"sku3", "sku4"},
					PaidSkus: []string{"sku4"},
				},
			},
			expected: []Bundle{
				{
					Name:     "bundle1",
					Skus:     []string{"sku2"},
					PaidSkus: []string{"sku1"},
				},
				{
					Name:     "bundle2",
					Skus:     []string{"sku3"},
					PaidSkus: []string{"sku4"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removePaidSkusFromSkus(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("removePaidSkusFromSkus() = %v, want %v", result, tt.expected)
			}
		})
	}
}
