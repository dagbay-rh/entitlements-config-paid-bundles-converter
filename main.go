package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Bundle struct {
	Name           	string   	`yaml:"name"`
	UseValidAccNum 	bool     	`yaml:"use_valid_acc_num,omitempty"`
	UseValidOrgId  	bool     	`yaml:"use_valid_org_id,omitempty"`
	UseIsInternal  	bool     	`yaml:"use_is_internal,omitempty"`
	Skus           	[]string 	`yaml:"skus,omitempty"`
	PaidSkus		[]string 	`yaml:"paid_skus,omitempty"`
}

type ConfigMap struct {
	ApiVersion string                 `yaml:"apiVersion"`
	Kind       string                 `yaml:"kind"`
	Metadata   map[string]interface{} `yaml:"metadata"`
	Data       map[string]string      `yaml:"data"`
}

type Template struct {
	Kind       string      `yaml:"kind"`
	ApiVersion string      `yaml:"apiVersion"`
	Objects    []ConfigMap `yaml:"objects"`
}

var originalTemplate Template

// writeBundlesToConfigMap writes bundles to a ConfigMap template file using the original template structure
func writeBundlesToConfigMap(bundles []Bundle, outputFilepath string) error {
	// Marshal the bundles to YAML
	bundlesYaml, err := yaml.Marshal(bundles)
	if err != nil {
		return fmt.Errorf("failed to marshal bundles to YAML: %w", err)
	}

	// Update the ConfigMap data with the new bundles
	originalTemplate.Objects[0].Data["bundles.yml"] = string(bundlesYaml)

	// Marshal the entire template back to YAML
	outputData, err := yaml.Marshal(&originalTemplate)
	if err != nil {
		return fmt.Errorf("failed to marshal template to YAML: %w", err)
	}

	// Write to the output file
	err = os.WriteFile(outputFilepath, outputData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write output file %s: %w", outputFilepath, err)
	}

	return nil
}

// removePaidSkusFromSkus creates a new bundle slice with paid SKUs removed from the Skus list
func removePaidSkusFromSkus(bundles []Bundle) []Bundle {
	newBundles := make([]Bundle, len(bundles))

	for i, bundle := range bundles {
		// Create a map of paid SKUs for quick lookup
		paidSkusMap := make(map[string]bool)
		for _, paidSku := range bundle.PaidSkus {
			paidSkusMap[paidSku] = true
		}

		// Filter out SKUs that are in PaidSkus
		filteredSkus := make([]string, 0, len(bundle.Skus))
		for _, sku := range bundle.Skus {
			if !paidSkusMap[sku] {
				filteredSkus = append(filteredSkus, sku)
			}
		}

		// Create new bundle with filtered SKUs
		newBundles[i] = Bundle{
			Name:           bundle.Name,
			UseValidAccNum: bundle.UseValidAccNum,
			UseValidOrgId:  bundle.UseValidOrgId,
			UseIsInternal:  bundle.UseIsInternal,
			Skus:           filteredSkus,
			PaidSkus:       bundle.PaidSkus,
		}
	}

	return newBundles
}

// readBundlesFromConfigMap reads bundles from a ConfigMap template file
func readBundlesFromConfigMap(filepath string) ([]Bundle, error) {
	// Read the ConfigMap template file
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filepath, err)
	}

	// Parse as OpenShift Template
	err = yaml.Unmarshal(data, &originalTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Template YAML: %w", err)
	}

	// Validate template structure
	if originalTemplate.Kind != "Template" {
		return nil, fmt.Errorf("expected Kind to be 'Template', got '%s'", originalTemplate.Kind)
	}

	if len(originalTemplate.Objects) == 0 {
		return nil, fmt.Errorf("template has no objects")
	}

	// Extract the ConfigMap
	configMap := originalTemplate.Objects[0]
	if configMap.Kind != "ConfigMap" {
		return nil, fmt.Errorf("expected first object to be 'ConfigMap', got '%s'", configMap.Kind)
	}

	// Extract bundles.yml from the ConfigMap data
	bundlesYaml, exists := configMap.Data["bundles.yml"]
	if !exists {
		return nil, fmt.Errorf("ConfigMap does not contain 'bundles.yml' key in data")
	}

	// Parse the bundles YAML content
	var bundles []Bundle
	err = yaml.Unmarshal([]byte(bundlesYaml), &bundles)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bundles YAML: %w", err)
	}

	return bundles, nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <filepath>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Convert entitlements config bundles to paid/not paid bundles.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -v, --verbose    Enable verbose output (prints loaded bundles)\n")
		fmt.Fprintf(os.Stderr, "  -h, --help       Show this help message\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  <filepath>       Path to the bundles ConfigMap template file\n")
	}

	// Define flags
	verbose := flag.Bool("v", false, "verbose output")
	flag.BoolVar(verbose, "verbose", false, "verbose output")

	// Parse flags
	flag.Parse()

	// Check if filepath argument is provided
	if flag.NArg() < 1 {
		log.Fatal("Usage: program [-v|--verbose] <filepath>")
	}

	bundlesFilepath := flag.Arg(0)

	// Read bundles from the ConfigMap template file
	bundles, err := readBundlesFromConfigMap(bundlesFilepath)
	if err != nil {
		log.Fatal(err)
	}

	if *verbose {
		fmt.Printf("Successfully loaded %d bundles from ConfigMap:\n", len(bundles))
		for i, bundle := range bundles {
			fmt.Printf("  [%d] %s (SKUs: %v, Paid SKUs: %v, UseValidAccNum: %t, UseValidOrgId: %t, UseIsInternal: %t)\n",
			 i,
			 bundle.Name, bundle.Skus, bundle.PaidSkus, bundle.UseValidAccNum, bundle.UseValidOrgId, bundle.UseIsInternal)
		}
	}
	
	// Remove paid SKUs from the regular SKUs list
	filteredBundles := removePaidSkusFromSkus(bundles)

	// Write the filtered bundles to a new ConfigMap file
	outputFilepath := bundlesFilepath + ".new"
	err = writeBundlesToConfigMap(filteredBundles, outputFilepath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully wrote filtered bundles to: %s\n", outputFilepath)
}