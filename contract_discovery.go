package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
)

// ContractInfo holds discovered contract and group information
type ContractInfo struct {
	ContractID   string   `json:"contractId"`
	ContractName string   `json:"contractName"`
	GroupID      string   `json:"groupId"`
	GroupName    string   `json:"groupName"`
	ProductIDs   []string `json:"productIds,omitempty"`
}

// AkamaiConfig holds cached configuration
type AkamaiConfig struct {
	ContractID   string   `json:"contractId"`
	ContractName string   `json:"contractName"`
	GroupID      string   `json:"groupId"`
	GroupName    string   `json:"groupName"`
	ProductIDs   []string `json:"productIds,omitempty"`
}

const configFileName = ".akamai-config.json"

// DiscoverContractByGroupName finds contract ID by exact group name match
func DiscoverContractByGroupName(ctx context.Context, papiClient papi.PAPI, groupName string) (*ContractInfo, error) {
	fmt.Printf("🔍 Searching for Group: %q\n", groupName)

	// Get all groups
	groupsResp, err := papiClient.GetGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	// Search for exact match
	for _, group := range groupsResp.Groups.Items {
		if group.GroupName == groupName {
			if len(group.ContractIDs) == 0 {
				return nil, fmt.Errorf("group %q has no associated contracts", groupName)
			}

			info := &ContractInfo{
				ContractID: group.ContractIDs[0],
				GroupID:    group.GroupID,
				GroupName:  group.GroupName,
			}

			// Get contract name
			contractsResp, err := papiClient.GetContracts(ctx)
			if err == nil {
				for _, contract := range contractsResp.Contracts.Items {
					if contract.ContractID == info.ContractID {
						info.ContractName = contract.ContractTypeName
						break
					}
				}
			}

			fmt.Printf("✅ Found matching group!\n")
			return info, nil
		}
	}

	return nil, fmt.Errorf("group %q not found", groupName)
}

// DiscoverContractByGroupID finds contract ID by group ID
func DiscoverContractByGroupID(ctx context.Context, papiClient papi.PAPI, groupID string) (*ContractInfo, error) {
	fmt.Printf("🔍 Searching for Group ID: %s\n", groupID)

	groupsResp, err := papiClient.GetGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	for _, group := range groupsResp.Groups.Items {
		if group.GroupID == groupID {
			if len(group.ContractIDs) == 0 {
				return nil, fmt.Errorf("group %s has no associated contracts", groupID)
			}

			info := &ContractInfo{
				ContractID: group.ContractIDs[0],
				GroupID:    group.GroupID,
				GroupName:  group.GroupName,
			}

			// Get contract name
			contractsResp, err := papiClient.GetContracts(ctx)
			if err == nil {
				for _, contract := range contractsResp.Contracts.Items {
					if contract.ContractID == info.ContractID {
						info.ContractName = contract.ContractTypeName
						break
					}
				}
			}

			fmt.Printf("✅ Found matching group!\n")
			return info, nil
		}
	}

	return nil, fmt.Errorf("group ID %s not found", groupID)
}

// DiscoverProductIDs discovers available product IDs for a contract
func DiscoverProductIDs(ctx context.Context, papiClient papi.PAPI, contractID string) ([]string, error) {
	fmt.Printf("🔍 Discovering Product IDs for Contract: %s\n", contractID)

	productsResp, err := papiClient.GetProducts(ctx, papi.GetProductsRequest{
		ContractID: contractID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	if len(productsResp.Products.Items) == 0 {
		return nil, fmt.Errorf("no products found for contract %s", contractID)
	}

	productIDs := make([]string, 0, len(productsResp.Products.Items))
	fmt.Println("\n📦 Available Products:")
	for i, product := range productsResp.Products.Items {
		fmt.Printf("   %d. %s (%s)\n", i+1, product.ProductName, product.ProductID)
		productIDs = append(productIDs, product.ProductID)
	}

	return productIDs, nil
}

// ListAllContractsAndGroups lists all available contracts and groups
func ListAllContractsAndGroups(ctx context.Context, papiClient papi.PAPI) ([]ContractInfo, error) {
	fmt.Println("📋 Listing all Contracts and Groups...")

	// Get contracts
	contractsResp, err := papiClient.GetContracts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get contracts: %w", err)
	}

	// Get groups
	groupsResp, err := papiClient.GetGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	var infos []ContractInfo

	fmt.Println("\n=== Available Contracts ===")
	for i, contract := range contractsResp.Contracts.Items {
		fmt.Printf("%d. %s (%s)\n", i+1, contract.ContractID, contract.ContractTypeName)
	}

	fmt.Println("\n=== Available Groups ===")
	for i, group := range groupsResp.Groups.Items {
		fmt.Printf("%d. %s (%s)\n", i+1, group.GroupName, group.GroupID)
		if len(group.ContractIDs) > 0 {
			fmt.Printf("   Contract(s): %v\n", group.ContractIDs)

			info := ContractInfo{
				ContractID: group.ContractIDs[0],
				GroupID:    group.GroupID,
				GroupName:  group.GroupName,
			}

			// Find contract name
			for _, contract := range contractsResp.Contracts.Items {
				if contract.ContractID == info.ContractID {
					info.ContractName = contract.ContractTypeName
					break
				}
			}

			infos = append(infos, info)
		}
	}

	return infos, nil
}

// SaveConfig saves discovered configuration to a JSON file
func SaveConfig(config *AkamaiConfig) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, configFileName)

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(configPath, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("\n💾 Configuration saved to: %s\n", configPath)
	return nil
}

// LoadConfig loads configuration from JSON file
func LoadConfig() (*AkamaiConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, configFileName)

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Config doesn't exist, not an error
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AkamaiConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	fmt.Printf("📖 Loaded cached configuration from: %s\n", configPath)
	return &config, nil
}

// DiscoverAndCache discovers contract/group information and caches it
func DiscoverAndCache(ctx context.Context, papiClient papi.PAPI, groupName string) (*AkamaiConfig, error) {
	// Try to load cached config first
	cachedConfig, err := LoadConfig()
	if err != nil {
		fmt.Printf("⚠️  Failed to load cached config: %v\n", err)
	} else if cachedConfig != nil && cachedConfig.GroupName == groupName {
		fmt.Println("✅ Using cached configuration")
		return cachedConfig, nil
	}

	// Discover contract info
	contractInfo, err := DiscoverContractByGroupName(ctx, papiClient, groupName)
	if err != nil {
		return nil, err
	}

	// Discover product IDs
	productIDs, err := DiscoverProductIDs(ctx, papiClient, contractInfo.ContractID)
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to discover products: %v\n", err)
		productIDs = []string{} // Continue without products
	}

	// Create config
	config := &AkamaiConfig{
		ContractID:   contractInfo.ContractID,
		ContractName: contractInfo.ContractName,
		GroupID:      contractInfo.GroupID,
		GroupName:    contractInfo.GroupName,
		ProductIDs:   productIDs,
	}

	// Save to cache
	if err := SaveConfig(config); err != nil {
		fmt.Printf("⚠️  Warning: Failed to save config: %v\n", err)
		// Continue even if save fails
	}

	return config, nil
}

// PrintDiscoveryResult prints the discovered information in a readable format
func PrintDiscoveryResult(config *AkamaiConfig) {
	fmt.Println("\n" + "=== Discovered Information ===")
	fmt.Printf("Group Name:    %s\n", config.GroupName)
	fmt.Printf("Group ID:      %s\n", config.GroupID)
	fmt.Printf("Contract ID:   %s\n", config.ContractID)
	fmt.Printf("Contract Name: %s\n", config.ContractName)

	if len(config.ProductIDs) > 0 {
		fmt.Printf("\nAvailable Product IDs (%d):\n", len(config.ProductIDs))
		for i, productID := range config.ProductIDs {
			fmt.Printf("  %d. %s\n", i+1, productID)
		}
	}

	fmt.Println("\n" + "=== Usage in main.go ===")
	fmt.Printf("const ContractID = %q\n", config.ContractID)
	fmt.Printf("const GroupID    = %q\n", config.GroupID)
	if len(config.ProductIDs) > 0 {
		fmt.Printf("const ProductID  = %q  // or choose from available products\n", config.ProductIDs[0])
	}
}
