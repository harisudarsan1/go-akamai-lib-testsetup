package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
)

// ProductFetcher helps with discovering and validating product IDs
type ProductFetcher struct {
	client papi.PAPI
	ctx    context.Context
}

// NewProductFetcher creates a new ProductFetcher instance
func NewProductFetcher(ctx context.Context, client papi.PAPI) *ProductFetcher {
	return &ProductFetcher{
		client: client,
		ctx:    ctx,
	}
}

// GetAllProducts retrieves all available products for a contract
func (pf *ProductFetcher) GetAllProducts(contractID string) (*papi.GetProductsResponse, error) {
	products, err := pf.client.GetProducts(pf.ctx, papi.GetProductsRequest{
		ContractID: contractID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	return products, nil
}

// ListProducts prints all available products in a formatted table
func (pf *ProductFetcher) ListProducts(contractID string) error {
	products, err := pf.GetAllProducts(contractID)
	if err != nil {
		return err
	}

	fmt.Printf("\nAvailable Products for Contract: %s\n", contractID)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%-50s %s\n", "Product Name", "Product ID")
	fmt.Println(strings.Repeat("-", 80))

	for _, product := range products.Products.Items {
		fmt.Printf("%-50s %s\n", product.ProductName, product.ProductID)
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Total Products: %d\n\n", len(products.Products.Items))

	return nil
}

// FindProductByName searches for a product by name (case-insensitive, partial match)
func (pf *ProductFetcher) FindProductByName(contractID, productName string) (string, error) {
	products, err := pf.GetAllProducts(contractID)
	if err != nil {
		return "", err
	}

	productNameLower := strings.ToLower(productName)

	// Try exact match first (case-insensitive)
	for _, product := range products.Products.Items {
		if strings.ToLower(product.ProductName) == productNameLower {
			return product.ProductID, nil
		}
	}

	// Try partial match
	var matches []papi.ProductItem
	for _, product := range products.Products.Items {
		if strings.Contains(strings.ToLower(product.ProductName), productNameLower) {
			matches = append(matches, product)
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no product found matching '%s'", productName)
	}

	if len(matches) == 1 {
		return matches[0].ProductID, nil
	}

	// Multiple matches found
	fmt.Printf("Multiple products found matching '%s':\n", productName)
	for i, match := range matches {
		fmt.Printf("  %d. %s (%s)\n", i+1, match.ProductName, match.ProductID)
	}
	return "", fmt.Errorf("multiple products match '%s', please be more specific", productName)
}

// FindProductByID verifies that a product ID exists in the contract
func (pf *ProductFetcher) FindProductByID(contractID, productID string) (string, error) {
	products, err := pf.GetAllProducts(contractID)
	if err != nil {
		return "", err
	}

	for _, product := range products.Products.Items {
		if product.ProductID == productID {
			return product.ProductName, nil
		}
	}

	return "", fmt.Errorf("product ID '%s' not found in contract", productID)
}

// ValidateProductID checks if a product ID is valid for the contract
func (pf *ProductFetcher) ValidateProductID(contractID, productID string) (bool, error) {
	_, err := pf.FindProductByID(contractID, productID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetRecommendedProductID suggests a product ID based on common names
func (pf *ProductFetcher) GetRecommendedProductID(contractID string, propertyType PropertyType) (string, error) {
	// Map property types to common search terms
	searchTerms := map[PropertyType][]string{
		PropertyTypeIonStandard:         {"ion standard", "ion", "web performance"},
		PropertyTypeDownloadDelivery:    {"download delivery", "download", "file download"},
		PropertyTypeDynamicSiteDelivery: {"dynamic site", "site accel", "dsa", "dynamic"},
		PropertyTypeMediaDelivery:       {"adaptive media", "media delivery", "video", "streaming"},
		PropertyTypeObjectDelivery:      {"object delivery", "object"},
		PropertyTypeAPIAcceleration:     {"api acceleration", "api"},
	}

	products, err := pf.GetAllProducts(contractID)
	if err != nil {
		return "", err
	}

	terms, exists := searchTerms[propertyType]
	if !exists {
		return "", fmt.Errorf("unknown property type: %s", propertyType)
	}

	// Try each search term
	for _, term := range terms {
		productID, err := pf.FindProductByName(contractID, term)
		if err == nil {
			return productID, nil
		}
	}

	// If nothing found, list available products
	fmt.Printf("\nNo matching product found for type: %s\n", propertyType)
	fmt.Println("Available products:")
	for _, product := range products.Products.Items {
		fmt.Printf("  - %s (%s)\n", product.ProductName, product.ProductID)
	}

	return "", fmt.Errorf("no suitable product found for type: %s", propertyType)
}

// ProductMapper helps map property types to actual product IDs for a contract
type ProductMapper struct {
	fetcher    *ProductFetcher
	contractID string
	cache      map[PropertyType]string
}

// NewProductMapper creates a new ProductMapper
func NewProductMapper(ctx context.Context, client papi.PAPI, contractID string) *ProductMapper {
	return &ProductMapper{
		fetcher:    NewProductFetcher(ctx, client),
		contractID: contractID,
		cache:      make(map[PropertyType]string),
	}
}

// GetProductID returns the product ID for a property type, using cache if available
func (pm *ProductMapper) GetProductID(propertyType PropertyType) (string, error) {
	// Check cache first
	if productID, exists := pm.cache[propertyType]; exists {
		return productID, nil
	}

	// Fetch from API
	productID, err := pm.fetcher.GetRecommendedProductID(pm.contractID, propertyType)
	if err != nil {
		return "", err
	}

	// Cache the result
	pm.cache[propertyType] = productID
	return productID, nil
}

// SetProductID manually sets a product ID for a property type
func (pm *ProductMapper) SetProductID(propertyType PropertyType, productID string) error {
	// Validate the product ID exists
	_, err := pm.fetcher.FindProductByID(pm.contractID, productID)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	pm.cache[propertyType] = productID
	return nil
}

// LoadFromConfig loads product mappings from a configuration
func (pm *ProductMapper) LoadFromConfig(config map[PropertyType]string) error {
	for propType, productID := range config {
		if err := pm.SetProductID(propType, productID); err != nil {
			return fmt.Errorf("failed to set product ID for %s: %w", propType, err)
		}
	}
	return nil
}

// GetAllMappings returns all cached property type to product ID mappings
func (pm *ProductMapper) GetAllMappings() map[PropertyType]string {
	result := make(map[PropertyType]string)
	for k, v := range pm.cache {
		result[k] = v
	}
	return result
}

// PrintMappings displays all current mappings
func (pm *ProductMapper) PrintMappings() {
	fmt.Println("\nProperty Type to Product ID Mappings:")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("%-30s %s\n", "Property Type", "Product ID")
	fmt.Println(strings.Repeat("-", 70))

	for propType, productID := range pm.cache {
		fmt.Printf("%-30s %s\n", propType, productID)
	}

	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("Total Mappings: %d\n\n", len(pm.cache))
}

// Helper function to discover and print all products for a contract
func DiscoverProducts(ctx context.Context, papiClient papi.PAPI, contractID string) error {
	fetcher := NewProductFetcher(ctx, papiClient)
	return fetcher.ListProducts(contractID)
}

// Helper function to find a specific product
func FindProduct(ctx context.Context, papiClient papi.PAPI, contractID, productName string) (string, error) {
	fetcher := NewProductFetcher(ctx, papiClient)
	return fetcher.FindProductByName(contractID, productName)
}

// ExampleUsage demonstrates how to use the product utilities
func ExampleProductUtilsUsage(ctx context.Context, papiClient papi.PAPI, contractID string) {
	// Example 1: List all products
	fmt.Println("=== Example 1: List All Products ===")
	fetcher := NewProductFetcher(ctx, papiClient)
	fetcher.ListProducts(contractID)

	// Example 2: Find a specific product
	fmt.Println("\n=== Example 2: Find Product by Name ===")
	productID, err := fetcher.FindProductByName(contractID, "ion")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Found Ion product: %s\n", productID)
	}

	// Example 3: Validate a product ID
	fmt.Println("\n=== Example 3: Validate Product ID ===")
	valid, err := fetcher.ValidateProductID(contractID, "prd_Ion")
	if err != nil {
		fmt.Printf("Product ID invalid: %v\n", err)
	} else {
		fmt.Printf("Product ID valid: %v\n", valid)
	}

	// Example 4: Use ProductMapper for property creation
	fmt.Println("\n=== Example 4: Use ProductMapper ===")
	mapper := NewProductMapper(ctx, papiClient, contractID)

	// Get recommended product ID for Ion Standard
	productID, err = mapper.GetProductID(PropertyTypeIonStandard)
	if err != nil {
		fmt.Printf("Error getting product ID: %v\n", err)
		return
	}

	fmt.Printf("Using Product ID for Ion Standard: %s\n", productID)

	// Show all mappings
	mapper.PrintMappings()

	// Example 5: Manually set product IDs
	fmt.Println("\n=== Example 5: Manual Product ID Mapping ===")
	customMappings := map[PropertyType]string{
		PropertyTypeIonStandard:      "prd_SPM",               // Your custom product ID
		PropertyTypeDownloadDelivery: "prd_Download_Delivery", // Your custom product ID
	}

	err = mapper.LoadFromConfig(customMappings)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	mapper.PrintMappings()
}
