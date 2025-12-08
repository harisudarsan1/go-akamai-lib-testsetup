package main

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
)

// Example functions demonstrating how to use the property helpers

// ExampleCreateIonStandardProperty demonstrates creating an Ion Standard property
func ExampleCreateIonStandardProperty(ctx context.Context, papiClient papi.PAPI) error {
	helper := NewPropertyHelper(ctx, papiClient)

	config := PropertyConfig{
		PropertyName:      "example-ion-property",
		ContractID:        ContractID,
		GroupID:           GroupID,
		Domain:            "www.example.com",
		OriginHostname:    "origin.example.com",
		CPCode:            123456,
		EnableCompression: true,
		EnableHTTP2:       true,
		EnableIPv6:        true,
		CacheTTL:          86400, // 1 day
	}

	prop, err := helper.CreateIonStandardProperty(config)
	if err != nil {
		return fmt.Errorf("failed to create Ion Standard property: %w", err)
	}

	fmt.Printf("Created Ion Standard property: %s (ID: %s)\n", prop.PropertyName, prop.PropertyID)
	return nil
}

// ExampleCreateDownloadDeliveryProperty demonstrates creating a Download Delivery property
func ExampleCreateDownloadDeliveryProperty(ctx context.Context, papiClient papi.PAPI) error {
	helper := NewPropertyHelper(ctx, papiClient)

	config := PropertyConfig{
		PropertyName:      "example-download-property",
		ContractID:        ContractID,
		GroupID:           GroupID,
		Domain:            "downloads.example.com",
		OriginHostname:    "origin-downloads.example.com",
		CPCode:            123457,
		EnableCompression: false, // Usually disabled for downloads
		EnableHTTP2:       true,
		CacheTTL:          604800, // 7 days
	}

	prop, err := helper.CreateDownloadDeliveryProperty(config)
	if err != nil {
		return fmt.Errorf("failed to create Download Delivery property: %w", err)
	}

	fmt.Printf("Created Download Delivery property: %s (ID: %s)\n", prop.PropertyName, prop.PropertyID)
	return nil
}

// ExampleCreateDynamicSiteDeliveryProperty demonstrates creating a Dynamic Site Delivery property
func ExampleCreateDynamicSiteDeliveryProperty(ctx context.Context, papiClient papi.PAPI) error {
	helper := NewPropertyHelper(ctx, papiClient)

	config := PropertyConfig{
		PropertyName:      "example-dynamic-property",
		ContractID:        ContractID,
		GroupID:           GroupID,
		Domain:            "app.example.com",
		OriginHostname:    "origin-app.example.com",
		CPCode:            123458,
		EnableCompression: true,
		EnableHTTP2:       true,
		CacheTTL:          0, // No caching for dynamic content
	}

	prop, err := helper.CreateDynamicSiteDeliveryProperty(config)
	if err != nil {
		return fmt.Errorf("failed to create Dynamic Site Delivery property: %w", err)
	}

	fmt.Printf("Created Dynamic Site Delivery property: %s (ID: %s)\n", prop.PropertyName, prop.PropertyID)
	return nil
}

// ExampleCreateMediaDeliveryProperty demonstrates creating an Adaptive Media Delivery property
func ExampleCreateMediaDeliveryProperty(ctx context.Context, papiClient papi.PAPI) error {
	helper := NewPropertyHelper(ctx, papiClient)

	config := PropertyConfig{
		PropertyName:      "example-media-property",
		ContractID:        ContractID,
		GroupID:           GroupID,
		Domain:            "media.example.com",
		OriginHostname:    "origin-media.example.com",
		CPCode:            123459,
		EnableCompression: false, // Don't compress media
		EnableHTTP2:       true,
		CacheTTL:          2592000, // 30 days
	}

	prop, err := helper.CreateMediaDeliveryProperty(config)
	if err != nil {
		return fmt.Errorf("failed to create Media Delivery property: %w", err)
	}

	fmt.Printf("Created Media Delivery property: %s (ID: %s)\n", prop.PropertyName, prop.PropertyID)
	return nil
}

// ExampleCreatePropertyWithCustomBehaviors demonstrates adding custom behaviors
func ExampleCreatePropertyWithCustomBehaviors(ctx context.Context, papiClient papi.PAPI) error {
	helper := NewPropertyHelper(ctx, papiClient)

	// Define custom behaviors
	customBehaviors := []papi.RuleBehavior{
		{
			Name: "gzipResponse",
			Options: papi.RuleOptionsMap{
				"behavior": "ALWAYS",
			},
		},
		{
			Name: "modifyOutgoingResponseHeader",
			Options: papi.RuleOptionsMap{
				"action":                "ADD",
				"standardAddHeaderName": "OTHER",
				"customHeaderName":      "X-Custom-Header",
				"headerValue":           "CustomValue",
			},
		},
	}

	config := PropertyConfig{
		PropertyName:      "example-custom-property",
		ContractID:        ContractID,
		GroupID:           GroupID,
		Domain:            "custom.example.com",
		OriginHostname:    "origin-custom.example.com",
		CPCode:            123460,
		EnableCompression: true,
		EnableHTTP2:       true,
		CacheTTL:          3600,
		CustomBehaviors:   customBehaviors,
	}

	prop, err := helper.CreateIonStandardProperty(config)
	if err != nil {
		return fmt.Errorf("failed to create property with custom behaviors: %w", err)
	}

	fmt.Printf("Created property with custom behaviors: %s (ID: %s)\n", prop.PropertyName, prop.PropertyID)
	return nil
}

// ExampleAddHostnameAndActivate demonstrates adding a hostname and activating
func ExampleAddHostnameAndActivate(ctx context.Context, papiClient papi.PAPI, prop *papi.Property) error {
	helper := NewPropertyHelper(ctx, papiClient)

	config := PropertyConfig{
		PropertyName:   prop.PropertyName,
		ContractID:     ContractID,
		GroupID:        GroupID,
		Domain:         "new.example.com",
		OriginHostname: "origin.example.com",
	}

	// Add hostname
	edgeHostname := "new.example.com.edgekey.net"
	err := helper.AddHostnameToProperty(prop, config, edgeHostname)
	if err != nil {
		return fmt.Errorf("failed to add hostname: %w", err)
	}

	// Activate to staging
	notifyEmails := []string{"ops@example.com"}
	activationResp, err := helper.ActivateProperty(prop, config, papi.ActivationNetworkStaging, notifyEmails)
	if err != nil {
		return fmt.Errorf("failed to activate property: %w", err)
	}

	fmt.Printf("Property activated to staging: ActivationID=%s\n", activationResp.ActivationID)
	return nil
}

// ExampleCompleteOnboarding demonstrates a complete onboarding flow with property helpers
func ExampleCompleteOnboarding(ctx context.Context, papiClient papi.PAPI) error {
	helper := NewPropertyHelper(ctx, papiClient)

	// Step 1: Create Ion Standard property
	fmt.Println("=== Step 1: Creating Ion Standard Property ===")
	config := PropertyConfig{
		PropertyName:      "complete-onboarding-example",
		ContractID:        ContractID,
		GroupID:           GroupID,
		Domain:            "onboard.example.com",
		OriginHostname:    "origin-onboard.example.com",
		CPCode:            123461,
		EnableCompression: true,
		EnableHTTP2:       true,
		EnableIPv6:        true,
		CacheTTL:          86400,
	}

	prop, err := helper.CreateIonStandardProperty(config)
	if err != nil {
		return fmt.Errorf("step 1 failed: %w", err)
	}

	// Step 2: Add hostname
	fmt.Println("\n=== Step 2: Adding Hostname ===")
	edgeHostname := "onboard.example.com.edgekey.net"
	err = helper.AddHostnameToProperty(prop, config, edgeHostname)
	if err != nil {
		return fmt.Errorf("step 2 failed: %w", err)
	}

	// Step 3: Activate to staging
	fmt.Println("\n=== Step 3: Activating to Staging ===")
	notifyEmails := []string{"ops@example.com", "devops@example.com"}
	activationResp, err := helper.ActivateProperty(prop, config, papi.ActivationNetworkStaging, notifyEmails)
	if err != nil {
		return fmt.Errorf("step 3 failed: %w", err)
	}

	fmt.Println("\n=== Onboarding Complete ===")
	fmt.Printf("Property: %s (ID: %s)\n", prop.PropertyName, prop.PropertyID)
	fmt.Printf("Activation: %s\n", activationResp.ActivationID)
	fmt.Println("Next steps:")
	fmt.Println("1. Wait for staging activation to complete")
	fmt.Println("2. Test on staging network")
	fmt.Println("3. Activate to production when ready")

	return nil
}

// ExampleMultiplePropertyTypes demonstrates creating multiple property types
func ExampleMultiplePropertyTypes(ctx context.Context, papiClient papi.PAPI) error {
	helper := NewPropertyHelper(ctx, papiClient)

	properties := []struct {
		name     string
		propType PropertyType
		domain   string
		origin   string
		cpcode   int
	}{
		{"website", PropertyTypeIonStandard, "www.example.com", "origin-www.example.com", 100001},
		{"downloads", PropertyTypeDownloadDelivery, "downloads.example.com", "origin-downloads.example.com", 100002},
		{"api", PropertyTypeDynamicSiteDelivery, "api.example.com", "origin-api.example.com", 100003},
		{"video", PropertyTypeMediaDelivery, "video.example.com", "origin-video.example.com", 100004},
	}

	for _, p := range properties {
		fmt.Printf("\n=== Creating %s (%s) ===\n", p.name, p.propType)

		config := PropertyConfig{
			PropertyName:      p.name,
			ContractID:        ContractID,
			GroupID:           GroupID,
			Domain:            p.domain,
			OriginHostname:    p.origin,
			CPCode:            p.cpcode,
			EnableCompression: p.propType != PropertyTypeMediaDelivery && p.propType != PropertyTypeDownloadDelivery,
			EnableHTTP2:       true,
			CacheTTL:          86400,
		}

		var prop *papi.Property
		var err error

		switch p.propType {
		case PropertyTypeIonStandard:
			prop, err = helper.CreateIonStandardProperty(config)
		case PropertyTypeDownloadDelivery:
			prop, err = helper.CreateDownloadDeliveryProperty(config)
		case PropertyTypeDynamicSiteDelivery:
			prop, err = helper.CreateDynamicSiteDeliveryProperty(config)
		case PropertyTypeMediaDelivery:
			prop, err = helper.CreateMediaDeliveryProperty(config)
		}

		if err != nil {
			return fmt.Errorf("failed to create %s property: %w", p.name, err)
		}

		fmt.Printf("✓ Created: %s (ID: %s)\n", prop.PropertyName, prop.PropertyID)
	}

	fmt.Println("\n=== All Properties Created Successfully ===")
	return nil
}
