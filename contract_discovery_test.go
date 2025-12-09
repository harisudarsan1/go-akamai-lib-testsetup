package main

import (
	"context"
	"strings"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
)

const (
	// Ion Group (has Ion Standard product)
	TestIonGroupName = "CTM LABS PRIVATE LIMITED (Kluisz)-V-5ZUL2W6"

	// Download Delivery Group (has Download Delivery product) - for reference
	TestDDGroupName = "CTM LABS PRIVATE LIMITED (Kluisz)-V-620VL0G"

	// Use Ion group for discovery tests (change this to switch between groups)
	TestGroupName = TestIonGroupName

	// Property names
	TestIonPropertyName = "test-ion-standard-property"
	TestDDPropertyName  = "test-download-delivery-property" // existing property

	// Use Ion property for tests
	TestPropertyName = TestIonPropertyName

	// Hostnames
	TestIonHostname     = "test-ion.kluisz.com"
	TestIonEdgeHostname = "test-ion.kluisz.com.edgekey.net"

	// Use Ion hostname for tests
	TestHostname     = TestIonHostname
	TestEdgeHostname = TestIonEdgeHostname

	// Custom property test
	CustomPropertyName = "propertyname"
	CustomHostname     = "example.kluisz.com"
	CustomEdgeHostname = "example.kluisz.com.edgekey.net"
)

// TestDiscoverContractByGroupName tests contract discovery by group name
func TestDiscoverContractByGroupName(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Akamai Contract Discovery Test ===\n")

	// Step 1: Authenticate
	t.Log("Step 1: Authenticating with Akamai API...")
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}
	t.Log("✅ Authentication successful")

	papiClient := papi.Client(sess)

	// Step 2: Discover contract by group name
	t.Logf("\nStep 2: Discovering Contract for Group: %q", TestGroupName)
	contractInfo, err := DiscoverContractByGroupName(ctx, papiClient, TestGroupName)
	if err != nil {
		t.Fatalf("❌ Failed to discover contract: %v", err)
	}

	// Step 3: Validate discovered information
	t.Log("\nStep 3: Validating discovered information...")
	if contractInfo.GroupName != TestGroupName {
		t.Errorf("❌ Group name mismatch: got %q, want %q", contractInfo.GroupName, TestGroupName)
	}

	if contractInfo.ContractID == "" {
		t.Error("❌ Contract ID is empty")
	}

	if contractInfo.GroupID == "" {
		t.Error("❌ Group ID is empty")
	}

	t.Log("✅ Discovered information validated")

	// Step 4: Discover Product IDs
	t.Log("\nStep 4: Discovering Product IDs...")
	productIDs, err := DiscoverProductIDs(ctx, papiClient, contractInfo.ContractID)
	if err != nil {
		t.Logf("⚠️  Warning: Failed to discover products: %v", err)
	} else {
		t.Logf("✅ Found %d product(s)", len(productIDs))
		contractInfo.ProductIDs = productIDs
	}

	// Step 5: Save to cache
	t.Log("\nStep 5: Saving configuration to cache...")
	config := &AkamaiConfig{
		ContractID:   contractInfo.ContractID,
		ContractName: contractInfo.ContractName,
		GroupID:      contractInfo.GroupID,
		GroupName:    contractInfo.GroupName,
		ProductIDs:   contractInfo.ProductIDs,
	}

	err = SaveConfig(config)
	if err != nil {
		t.Logf("⚠️  Warning: Failed to save config: %v", err)
	} else {
		t.Log("✅ Configuration saved successfully")
	}

	// Step 6: Print results
	t.Log("\n" + "=================================")
	PrintDiscoveryResult(config)
	t.Log("=================================")

	// Step 7: Test loading from cache
	t.Log("\nStep 7: Testing cache loading...")
	loadedConfig, err := LoadConfig()
	if err != nil {
		t.Logf("⚠️  Warning: Failed to load config: %v", err)
	} else if loadedConfig != nil {
		t.Log("✅ Successfully loaded configuration from cache")
		if loadedConfig.ContractID != config.ContractID {
			t.Errorf("❌ Cached Contract ID mismatch: got %q, want %q", loadedConfig.ContractID, config.ContractID)
		}
	}

	t.Log("\n✅ All tests passed!")
}

// TestListAllContractsAndGroups tests listing all available resources
func TestListAllContractsAndGroups(t *testing.T) {
	ctx := context.Background()

	t.Log("=== List All Contracts and Groups ===\n")

	// Authenticate
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	papiClient := papi.Client(sess)

	// List all
	infos, err := ListAllContractsAndGroups(ctx, papiClient)
	if err != nil {
		t.Fatalf("❌ Failed to list contracts and groups: %v", err)
	}

	t.Logf("\n✅ Found %d contract/group combination(s)", len(infos))

	if len(infos) == 0 {
		t.Error("❌ No contracts or groups found")
	}
}

// TestDiscoverAndCache tests the full discovery and caching workflow
func TestDiscoverAndCache(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Test Discover and Cache Workflow ===\n")

	// Authenticate
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	papiClient := papi.Client(sess)

	// Discover and cache
	config, err := DiscoverAndCache(ctx, papiClient, TestGroupName)
	if err != nil {
		t.Fatalf("❌ Failed to discover and cache: %v", err)
	}

	t.Log("✅ Discovery and caching completed successfully")

	// Validate
	if config.ContractID == "" {
		t.Error("❌ Contract ID is empty")
	}

	if config.GroupID == "" {
		t.Error("❌ Group ID is empty")
	}

	if config.GroupName != TestGroupName {
		t.Errorf("❌ Group name mismatch: got %q, want %q", config.GroupName, TestGroupName)
	}

	t.Log("\n" + "=== Cached Configuration ===")
	PrintDiscoveryResult(config)
}

// TestCreateIonPropertyIfNotExists tests creating a property with Ion product
func TestCreateIonPropertyIfNotExists(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Test Create Ion Property If Not Exists ===\n")

	// Step 1: Authenticate
	t.Log("Step 1: Authenticating with Akamai API...")
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}
	t.Log("✅ Authentication successful")

	papiClient := papi.Client(sess)

	// Step 2: Discover Ion group configuration (force fresh discovery for Ion group)
	t.Log("\nStep 2: Discovering Ion group configuration...")
	t.Logf("   Target Group: %s", TestIonGroupName)

	config, err := DiscoverContractByGroupName(ctx, papiClient, TestIonGroupName)
	if err != nil {
		t.Fatalf("❌ Failed to discover Ion group: %v", err)
	}

	t.Logf("✅ Using Contract: %s, Group: %s", config.ContractID, config.GroupID)

	// Discover products for this contract
	productIDs, err := DiscoverProductIDs(ctx, papiClient, config.ContractID)
	if err != nil {
		t.Logf("⚠️  Warning: Failed to discover products: %v", err)
		productIDs = []string{}
	} else {
		config.ProductIDs = productIDs
	}

	// Step 3: Check if Ion product is available
	t.Log("\nStep 3: Checking for Ion product...")
	var ionProductID string
	for _, productID := range config.ProductIDs {
		// Look for Ion product variants:
		// - prd_Ion, prd_SPM: Standard Ion product IDs
		// - prd_Fresca: Internal Akamai name for Ion Standard
		// - Any product containing "ion" in the name
		productLower := strings.ToLower(productID)
		if productID == "prd_Ion" || productID == "prd_SPM" || productID == "prd_Fresca" ||
			strings.Contains(productLower, "ion") {
			ionProductID = productID
			t.Logf("   Found Ion-compatible product: %s", productID)
			break
		}
	}

	if ionProductID == "" {
		t.Fatalf("❌ Ion product not found in group: %s\n   Available products: %v\n   This group doesn't have Ion Standard or Fresca product!",
			TestIonGroupName, config.ProductIDs)
	}

	t.Logf("✅ Using Ion product: %s", ionProductID)

	// Step 4: Check if property already exists
	t.Log("\nStep 4: Checking if property exists...")
	properties, err := papiClient.GetProperties(ctx, papi.GetPropertiesRequest{
		ContractID: config.ContractID,
		GroupID:    config.GroupID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to list properties: %v", err)
	}

	var existingProperty *papi.Property
	for _, prop := range properties.Properties.Items {
		if prop.PropertyName == TestPropertyName {
			existingProperty = prop
			break
		}
	}

	if existingProperty != nil {
		t.Logf("✅ Property already exists: %s (ID: %s, Version: %d)",
			existingProperty.PropertyName,
			existingProperty.PropertyID,
			existingProperty.LatestVersion)
		t.Log("   Skipping creation (property already exists)")
		return
	}

	// Step 5: Create new property
	t.Logf("\nStep 5: Creating new property: %s", TestPropertyName)
	createResp, err := papiClient.CreateProperty(ctx, papi.CreatePropertyRequest{
		ContractID: config.ContractID,
		GroupID:    config.GroupID,
		Property: papi.PropertyCreate{
			ProductID:    ionProductID,
			PropertyName: TestPropertyName,
			RuleFormat:   "latest",
		},
	})
	if err != nil {
		t.Fatalf("❌ Failed to create property: %v", err)
	}

	t.Logf("✅ Property created successfully!")
	t.Logf("   Property ID: %s", createResp.PropertyID)

	// Step 6: Verify property was created
	t.Log("\nStep 6: Verifying property creation...")
	propResp, err := papiClient.GetProperty(ctx, papi.GetPropertyRequest{
		ContractID: config.ContractID,
		GroupID:    config.GroupID,
		PropertyID: createResp.PropertyID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to verify property: %v", err)
	}

	if propResp.Property == nil {
		t.Fatal("❌ Property verification failed: property is nil")
	}

	t.Log("✅ Property verified successfully!")
	t.Logf("   Property Name: %s", propResp.Property.PropertyName)
	t.Logf("   Property ID: %s", propResp.Property.PropertyID)
	t.Logf("   Latest Version: %d", propResp.Property.LatestVersion)
	t.Logf("   Product: %s", ionProductID)
	t.Logf("   Contract ID: %s", config.ContractID)
	t.Logf("   Group ID: %s", config.GroupID)

	// Step 7: Display summary
	t.Log("\n" + "=== Property Creation Summary ===")
	t.Log("✅ Property created and verified")
	t.Log("   You can now:")
	t.Log("   - Add hostnames to the property")
	t.Log("   - Configure edge hostnames")
	t.Log("   - Set up origin server")
	t.Log("   - Activate to staging/production")
	t.Log("=================================")

	t.Log("\n✅ Test completed successfully!")
}

// TestAddHostnameToProperty tests adding a hostname to an existing property
func TestAddHostnameToProperty(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Test Add Hostname to Property ===\n")

	// Step 1: Authenticate
	t.Log("Step 1: Authenticating with Akamai API...")
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}
	t.Log("✅ Authentication successful")

	papiClient := papi.Client(sess)

	// Step 2: Load configuration
	t.Log("\nStep 2: Loading configuration...")
	config, err := LoadConfig()
	if err != nil || config == nil {
		t.Log("⚠️  No cached config, discovering...")
		config, err = DiscoverAndCache(ctx, papiClient, TestGroupName)
		if err != nil {
			t.Fatalf("❌ Failed to discover configuration: %v", err)
		}
	}
	t.Logf("✅ Using Contract: %s, Group: %s", config.ContractID, config.GroupID)

	// Step 3: Get the test property
	t.Log("\nStep 3: Getting test property...")
	properties, err := papiClient.GetProperties(ctx, papi.GetPropertiesRequest{
		ContractID: config.ContractID,
		GroupID:    config.GroupID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to list properties: %v", err)
	}

	var testProperty *papi.Property
	for _, prop := range properties.Properties.Items {
		if prop.PropertyName == TestPropertyName {
			testProperty = prop
			break
		}
	}

	if testProperty == nil {
		t.Skip("⚠️  Test property not found. Run TestCreateIonPropertyIfNotExists first.")
	}

	t.Logf("✅ Found property: %s (ID: %s, Version: %d)",
		testProperty.PropertyName,
		testProperty.PropertyID,
		testProperty.LatestVersion)

	// Step 4: Get current hostnames
	t.Log("\nStep 4: Checking current hostnames...")
	hostnamesResp, err := papiClient.GetPropertyVersionHostnames(ctx, papi.GetPropertyVersionHostnamesRequest{
		PropertyID:      testProperty.PropertyID,
		PropertyVersion: testProperty.LatestVersion,
		ContractID:      config.ContractID,
		GroupID:         config.GroupID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to get hostnames: %v", err)
	}

	t.Logf("   Current hostname count: %d", len(hostnamesResp.Hostnames.Items))

	// Check if test hostname already exists
	hostnameExists := false
	for _, h := range hostnamesResp.Hostnames.Items {
		if h.CnameFrom == TestHostname {
			hostnameExists = true
			t.Logf("   ✅ Hostname already exists: %s → %s", h.CnameFrom, h.CnameTo)
			break
		}
	}

	if hostnameExists {
		t.Log("   Skipping addition (hostname already exists)")
		t.Log("\n✅ Test completed successfully (hostname already configured)!")
		return
	}

	// Step 5: Add new hostname
	t.Logf("\nStep 5: Adding hostname: %s → %s", TestHostname, TestEdgeHostname)

	newHostnames := append(hostnamesResp.Hostnames.Items, papi.Hostname{
		CnameType:            papi.HostnameCnameTypeEdgeHostname,
		CnameFrom:            TestHostname,
		CnameTo:              TestEdgeHostname,
		CertProvisioningType: "DEFAULT", // Using default/shared certificate
	})

	_, err = papiClient.UpdatePropertyVersionHostnames(ctx, papi.UpdatePropertyVersionHostnamesRequest{
		PropertyID:      testProperty.PropertyID,
		PropertyVersion: testProperty.LatestVersion,
		ContractID:      config.ContractID,
		GroupID:         config.GroupID,
		Hostnames:       newHostnames,
	})
	if err != nil {
		t.Fatalf("❌ Failed to add hostname: %v", err)
	}

	t.Log("✅ Hostname added successfully!")

	// Step 6: Verify hostname was added
	t.Log("\nStep 6: Verifying hostname addition...")
	verifyResp, err := papiClient.GetPropertyVersionHostnames(ctx, papi.GetPropertyVersionHostnamesRequest{
		PropertyID:      testProperty.PropertyID,
		PropertyVersion: testProperty.LatestVersion,
		ContractID:      config.ContractID,
		GroupID:         config.GroupID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to verify hostname: %v", err)
	}

	// Check if our hostname is in the list
	hostnameFound := false
	for _, h := range verifyResp.Hostnames.Items {
		if h.CnameFrom == TestHostname {
			hostnameFound = true
			t.Log("✅ Hostname verified successfully!")
			t.Logf("   CnameFrom: %s", h.CnameFrom)
			t.Logf("   CnameTo: %s", h.CnameTo)
			t.Logf("   CnameType: %s", h.CnameType)
			t.Logf("   CertProvisioningType: %s", h.CertProvisioningType)
			break
		}
	}

	if !hostnameFound {
		t.Fatal("❌ Hostname verification failed: hostname not found in property")
	}

	t.Logf("   Total hostnames now: %d", len(verifyResp.Hostnames.Items))

	// Step 7: Display summary
	t.Log("\n" + "=== Hostname Configuration Summary ===")
	t.Log("✅ Hostname added and verified")
	t.Logf("   Property: %s (ID: %s)", testProperty.PropertyName, testProperty.PropertyID)
	t.Logf("   Hostname: %s", TestHostname)
	t.Logf("   Edge Hostname: %s", TestEdgeHostname)
	t.Log("\n   Next steps:")
	t.Log("   - The edge hostname (*.edgekey.net) must be created via API")
	t.Log("   - Configure DNS CNAME: test.kluisz.com → test.kluisz.com.edgekey.net")
	t.Log("   - Activate property to staging/production")
	t.Log("=================================")

	t.Log("\n✅ Test completed successfully!")
}

// TestAddHostnameToCustomProperty tests adding a hostname to the custom property "propertyname"
func TestAddHostnameToCustomProperty(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Test Add Hostname to Custom Property ===\n")

	// Step 1: Authenticate
	t.Log("Step 1: Authenticating with Akamai API...")
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}
	t.Log("✅ Authentication successful")

	papiClient := papi.Client(sess)

	// Step 2: Discover Ion group configuration
	t.Log("\nStep 2: Discovering Ion group configuration...")
	t.Logf("   Target Group: %s", TestIonGroupName)

	config, err := DiscoverContractByGroupName(ctx, papiClient, TestIonGroupName)
	if err != nil {
		t.Fatalf("❌ Failed to discover Ion group: %v", err)
	}

	t.Logf("✅ Using Contract: %s, Group: %s", config.ContractID, config.GroupID)

	// Step 3: Verify property exists
	t.Logf("\nStep 3: Verifying property '%s' exists...", CustomPropertyName)
	properties, err := papiClient.GetProperties(ctx, papi.GetPropertiesRequest{
		ContractID: config.ContractID,
		GroupID:    config.GroupID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to list properties: %v", err)
	}

	var targetProperty *papi.Property
	for _, prop := range properties.Properties.Items {
		if prop.PropertyName == CustomPropertyName {
			targetProperty = prop
			break
		}
	}

	if targetProperty == nil {
		t.Fatalf("❌ Property '%s' not found in Ion group (grp_303793)\n"+
			"   Please create the property first or verify the property name.\n"+
			"   Available properties: %d", CustomPropertyName, len(properties.Properties.Items))
	}

	t.Logf("✅ Found property: %s (ID: %s, Version: %d)",
		targetProperty.PropertyName,
		targetProperty.PropertyID,
		targetProperty.LatestVersion)

	// Step 4: Get current hostnames
	t.Log("\nStep 4: Checking current hostnames...")
	hostnamesResp, err := papiClient.GetPropertyVersionHostnames(ctx, papi.GetPropertyVersionHostnamesRequest{
		PropertyID:      targetProperty.PropertyID,
		PropertyVersion: targetProperty.LatestVersion,
		ContractID:      config.ContractID,
		GroupID:         config.GroupID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to get hostnames: %v", err)
	}

	t.Logf("   Current hostname count: %d", len(hostnamesResp.Hostnames.Items))
	for _, h := range hostnamesResp.Hostnames.Items {
		t.Logf("   - %s → %s", h.CnameFrom, h.CnameTo)
	}

	// Check if custom hostname already exists
	hostnameExists := false
	for _, h := range hostnamesResp.Hostnames.Items {
		if h.CnameFrom == CustomHostname {
			hostnameExists = true
			t.Logf("   ✅ Hostname already exists: %s → %s", h.CnameFrom, h.CnameTo)
			break
		}
	}

	if hostnameExists {
		t.Log("   Skipping addition (hostname already exists)")
		t.Log("\n✅ Test completed successfully (hostname already configured)!")
		return
	}

	// Step 5: Add new hostname
	t.Logf("\nStep 5: Adding hostname: %s → %s", CustomHostname, CustomEdgeHostname)

	newHostnames := append(hostnamesResp.Hostnames.Items, papi.Hostname{
		CnameType:            papi.HostnameCnameTypeEdgeHostname,
		CnameFrom:            CustomHostname,
		CnameTo:              CustomEdgeHostname,
		CertProvisioningType: "DEFAULT", // Using default/shared certificate
	})

	_, err = papiClient.UpdatePropertyVersionHostnames(ctx, papi.UpdatePropertyVersionHostnamesRequest{
		PropertyID:      targetProperty.PropertyID,
		PropertyVersion: targetProperty.LatestVersion,
		ContractID:      config.ContractID,
		GroupID:         config.GroupID,
		Hostnames:       newHostnames,
	})
	if err != nil {
		t.Fatalf("❌ Failed to add hostname: %v", err)
	}

	t.Log("✅ Hostname added successfully!")

	// Step 6: Verify hostname was added
	t.Log("\nStep 6: Verifying hostname addition...")
	verifyResp, err := papiClient.GetPropertyVersionHostnames(ctx, papi.GetPropertyVersionHostnamesRequest{
		PropertyID:      targetProperty.PropertyID,
		PropertyVersion: targetProperty.LatestVersion,
		ContractID:      config.ContractID,
		GroupID:         config.GroupID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to verify hostname: %v", err)
	}

	// Find the newly added hostname
	hostnameFound := false
	for _, h := range verifyResp.Hostnames.Items {
		if h.CnameFrom == CustomHostname {
			hostnameFound = true
			t.Log("✅ Hostname verified successfully!")
			t.Logf("   CnameFrom: %s", h.CnameFrom)
			t.Logf("   CnameTo: %s", h.CnameTo)
			t.Logf("   CnameType: %s", h.CnameType)
			t.Logf("   CertProvisioningType: %s", h.CertProvisioningType)
			break
		}
	}

	if !hostnameFound {
		t.Fatal("❌ Hostname verification failed: hostname not found after addition")
	}

	t.Logf("   Total hostnames now: %d", len(verifyResp.Hostnames.Items))

	// Step 7: Display summary
	t.Log("\n" + "=== Hostname Configuration Summary ===")
	t.Log("✅ Hostname added and verified")
	t.Logf("   Property: %s (ID: %s)", targetProperty.PropertyName, targetProperty.PropertyID)
	t.Logf("   Hostname: %s", CustomHostname)
	t.Logf("   Edge Hostname: %s", CustomEdgeHostname)
	t.Log("\n   Next steps:")
	t.Log("   - The edge hostname (*.edgekey.net) must be created via API")
	t.Logf("   - Configure DNS CNAME: %s → %s", CustomHostname, CustomEdgeHostname)
	t.Log("   - Activate property to staging/production")
	t.Log("=================================")

	t.Log("\n✅ Test completed successfully!")
}
