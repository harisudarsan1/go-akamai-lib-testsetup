package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
)

// PropertyType represents different Akamai product types
type PropertyType string

const (
	// PropertyTypeIonStandard represents Ion Standard (Web Performance)
	PropertyTypeIonStandard PropertyType = "Ion_Standard"

	// PropertyTypeDownloadDelivery represents Download Delivery
	PropertyTypeDownloadDelivery PropertyType = "Download_Delivery"

	// PropertyTypeDynamicSiteDelivery represents Dynamic Site Delivery (formerly DSA)
	PropertyTypeDynamicSiteDelivery PropertyType = "Dynamic_Site_Delivery"

	// PropertyTypeMediaDelivery represents Adaptive Media Delivery
	PropertyTypeMediaDelivery PropertyType = "Adaptive_Media_Delivery"

	// PropertyTypeObjectDelivery represents Object Delivery (formerly Web Performance + File Download)
	PropertyTypeObjectDelivery PropertyType = "Object_Delivery"

	// PropertyTypeAPIAcceleration represents API Acceleration
	PropertyTypeAPIAcceleration PropertyType = "API_Acceleration"
)

// PropertyConfig contains configuration for creating a property
type PropertyConfig struct {
	PropertyName      string
	ContractID        string
	GroupID           string
	ProductID         string
	Domain            string
	OriginHostname    string
	CPCode            int
	EdgeHostnameID    string
	CertEnrollmentID  int
	EnableCompression bool
	EnableHTTP2       bool
	EnableIPv6        bool
	CacheTTL          int // in seconds
	CustomBehaviors   []papi.RuleBehavior
}

// PropertyHelper provides methods for creating different property types
type PropertyHelper struct {
	client papi.PAPI
	ctx    context.Context
}

// NewPropertyHelper creates a new PropertyHelper instance
func NewPropertyHelper(ctx context.Context, client papi.PAPI) *PropertyHelper {
	return &PropertyHelper{
		client: client,
		ctx:    ctx,
	}
}

// CreateProperty creates a property with the given configuration
func (ph *PropertyHelper) CreateProperty(config PropertyConfig) (*papi.Property, error) {
	fmt.Printf(">> Creating property: %s (Product: %s)\n", config.PropertyName, config.ProductID)

	// Check if property already exists
	properties, err := ph.client.GetProperties(ph.ctx, papi.GetPropertiesRequest{
		ContractID: config.ContractID,
		GroupID:    config.GroupID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list properties: %w", err)
	}

	// Search for existing property
	for _, prop := range properties.Properties.Items {
		if prop.PropertyName == config.PropertyName {
			fmt.Printf(">> Found existing property: %s (ID: %s)\n", prop.PropertyName, prop.PropertyID)
			return prop, nil
		}
	}

	// Create new property
	createResp, err := ph.client.CreateProperty(ph.ctx, papi.CreatePropertyRequest{
		ContractID: config.ContractID,
		GroupID:    config.GroupID,
		Property: papi.PropertyCreate{
			ProductID:    config.ProductID,
			PropertyName: config.PropertyName,
			RuleFormat:   "latest",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create property: %w", err)
	}

	// Fetch the created property
	propResp, err := ph.client.GetProperty(ph.ctx, papi.GetPropertyRequest{
		ContractID: config.ContractID,
		GroupID:    config.GroupID,
		PropertyID: createResp.PropertyID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created property: %w", err)
	}

	fmt.Printf(">> Property created: %s (ID: %s, Version: %d)\n",
		config.PropertyName, propResp.Property.PropertyID, propResp.Property.LatestVersion)

	return propResp.Property, nil
}

// CreateIonStandardProperty creates an Ion Standard property optimized for web performance
func (ph *PropertyHelper) CreateIonStandardProperty(config PropertyConfig) (*papi.Property, error) {
	config.ProductID = "prd_Ion"

	prop, err := ph.CreateProperty(config)
	if err != nil {
		return nil, err
	}

	// Configure Ion Standard optimizations
	err = ph.configureIonStandard(prop, config)
	if err != nil {
		return nil, fmt.Errorf("failed to configure Ion Standard: %w", err)
	}

	return prop, nil
}

// CreateDownloadDeliveryProperty creates a Download Delivery property optimized for large file downloads
func (ph *PropertyHelper) CreateDownloadDeliveryProperty(config PropertyConfig) (*papi.Property, error) {
	config.ProductID = "prd_Download_Delivery"

	prop, err := ph.CreateProperty(config)
	if err != nil {
		return nil, err
	}

	// Configure Download Delivery optimizations
	err = ph.configureDownloadDelivery(prop, config)
	if err != nil {
		return nil, fmt.Errorf("failed to configure Download Delivery: %w", err)
	}

	return prop, nil
}

// CreateDynamicSiteDeliveryProperty creates a Dynamic Site Delivery property
func (ph *PropertyHelper) CreateDynamicSiteDeliveryProperty(config PropertyConfig) (*papi.Property, error) {
	config.ProductID = "prd_Site_Accel"

	prop, err := ph.CreateProperty(config)
	if err != nil {
		return nil, err
	}

	// Configure Dynamic Site Delivery optimizations
	err = ph.configureDynamicSiteDelivery(prop, config)
	if err != nil {
		return nil, fmt.Errorf("failed to configure Dynamic Site Delivery: %w", err)
	}

	return prop, nil
}

// CreateMediaDeliveryProperty creates an Adaptive Media Delivery property
func (ph *PropertyHelper) CreateMediaDeliveryProperty(config PropertyConfig) (*papi.Property, error) {
	config.ProductID = "prd_Adaptive_Media_Delivery"

	prop, err := ph.CreateProperty(config)
	if err != nil {
		return nil, err
	}

	// Configure Media Delivery optimizations
	err = ph.configureMediaDelivery(prop, config)
	if err != nil {
		return nil, fmt.Errorf("failed to configure Media Delivery: %w", err)
	}

	return prop, nil
}

// configureIonStandard configures Ion Standard specific behaviors
func (ph *PropertyHelper) configureIonStandard(prop *papi.Property, config PropertyConfig) error {
	fmt.Println(">> Configuring Ion Standard behaviors")

	ruleTree, err := ph.getRuleTree(prop, config)
	if err != nil {
		return err
	}

	// Ion Standard behaviors
	behaviors := []papi.RuleBehavior{
		// Origin configuration
		{
			Name: "origin",
			Options: papi.RuleOptionsMap{
				"originType":         "CUSTOMER",
				"hostname":           config.OriginHostname,
				"forwardHostHeader":  "REQUEST_HOST_HEADER",
				"cacheKeyHostname":   "REQUEST_HOST_HEADER",
				"compress":           config.EnableCompression,
				"enableTrueClientIp": false,
				"verificationMode":   "PLATFORM_SETTINGS",
				"httpPort":           80,
				"httpsPort":          443,
			},
		},
		// CP Code
		{
			Name: "cpCode",
			Options: papi.RuleOptionsMap{
				"value": map[string]interface{}{
					"id": config.CPCode,
				},
			},
		},
		// Caching
		{
			Name: "caching",
			Options: papi.RuleOptionsMap{
				"behavior": "MAX_AGE",
				"ttl":      fmt.Sprintf("%ds", config.CacheTTL),
			},
		},
		// HTTP/2
		{
			Name: "http2",
			Options: papi.RuleOptionsMap{
				"enabled": config.EnableHTTP2,
			},
		},
		// SureRoute (Ion specific)
		{
			Name: "sureRoute",
			Options: papi.RuleOptionsMap{
				"enabled":         true,
				"type":            "PERFORMANCE",
				"testObjectUrl":   "/akamai/sureroute-test-object.html",
				"forceSslForward": false,
				"raceStatTtl":     "30m",
				"toHostStatus":    "INCOMING_HH",
			},
		},
		// Prefetch (Ion specific)
		{
			Name: "prefetch",
			Options: papi.RuleOptionsMap{
				"enabled": true,
			},
		},
		// Advanced caching features
		{
			Name: "cacheError",
			Options: papi.RuleOptionsMap{
				"enabled": true,
				"ttl":     "10s",
			},
		},
	}

	// Merge with existing behaviors
	ruleTree.Rules.Behaviors = mergeBehaviors(ruleTree.Rules.Behaviors, behaviors)

	// Add custom behaviors if provided
	if len(config.CustomBehaviors) > 0 {
		ruleTree.Rules.Behaviors = append(ruleTree.Rules.Behaviors, config.CustomBehaviors...)
	}

	return ph.updateRuleTree(prop, config, ruleTree)
}

// configureDownloadDelivery configures Download Delivery specific behaviors
func (ph *PropertyHelper) configureDownloadDelivery(prop *papi.Property, config PropertyConfig) error {
	fmt.Println(">> Configuring Download Delivery behaviors")

	ruleTree, err := ph.getRuleTree(prop, config)
	if err != nil {
		return err
	}

	// Download Delivery behaviors
	behaviors := []papi.RuleBehavior{
		// Origin configuration
		{
			Name: "origin",
			Options: papi.RuleOptionsMap{
				"originType":         "CUSTOMER",
				"hostname":           config.OriginHostname,
				"forwardHostHeader":  "REQUEST_HOST_HEADER",
				"cacheKeyHostname":   "REQUEST_HOST_HEADER",
				"compress":           false, // Usually false for downloads
				"enableTrueClientIp": false,
				"verificationMode":   "PLATFORM_SETTINGS",
				"httpPort":           80,
				"httpsPort":          443,
			},
		},
		// CP Code
		{
			Name: "cpCode",
			Options: papi.RuleOptionsMap{
				"value": map[string]interface{}{
					"id": config.CPCode,
				},
			},
		},
		// Aggressive caching for downloads
		{
			Name: "caching",
			Options: papi.RuleOptionsMap{
				"behavior": "MAX_AGE",
				"ttl":      "7d", // Long cache for downloads
			},
		},
		// Large file optimization
		{
			Name: "largeFileOptimization",
			Options: papi.RuleOptionsMap{
				"enabled":                    true,
				"enablePartialObjectCaching": true,
				"minimumSize":                "100MB",
				"maximumSize":                "323GB",
			},
		},
		// Prefetching for sequential access
		{
			Name: "prefetchable",
			Options: papi.RuleOptionsMap{
				"enabled": true,
			},
		},
		// Allow POST for downloads
		{
			Name: "allowPost",
			Options: papi.RuleOptionsMap{
				"enabled":                   true,
				"allowWithoutContentLength": false,
			},
		},
		// HTTP/2 for better performance
		{
			Name: "http2",
			Options: papi.RuleOptionsMap{
				"enabled": true,
			},
		},
	}

	// Merge with existing behaviors
	ruleTree.Rules.Behaviors = mergeBehaviors(ruleTree.Rules.Behaviors, behaviors)

	// Add custom behaviors if provided
	if len(config.CustomBehaviors) > 0 {
		ruleTree.Rules.Behaviors = append(ruleTree.Rules.Behaviors, config.CustomBehaviors...)
	}

	return ph.updateRuleTree(prop, config, ruleTree)
}

// configureDynamicSiteDelivery configures Dynamic Site Delivery behaviors
func (ph *PropertyHelper) configureDynamicSiteDelivery(prop *papi.Property, config PropertyConfig) error {
	fmt.Println(">> Configuring Dynamic Site Delivery behaviors")

	ruleTree, err := ph.getRuleTree(prop, config)
	if err != nil {
		return err
	}

	// Dynamic Site Delivery behaviors
	behaviors := []papi.RuleBehavior{
		// Origin configuration
		{
			Name: "origin",
			Options: papi.RuleOptionsMap{
				"originType":         "CUSTOMER",
				"hostname":           config.OriginHostname,
				"forwardHostHeader":  "REQUEST_HOST_HEADER",
				"cacheKeyHostname":   "REQUEST_HOST_HEADER",
				"compress":           config.EnableCompression,
				"enableTrueClientIp": true, // Important for dynamic sites
				"verificationMode":   "PLATFORM_SETTINGS",
				"httpPort":           80,
				"httpsPort":          443,
			},
		},
		// CP Code
		{
			Name: "cpCode",
			Options: papi.RuleOptionsMap{
				"value": map[string]interface{}{
					"id": config.CPCode,
				},
			},
		},
		// Conservative caching for dynamic content
		{
			Name: "caching",
			Options: papi.RuleOptionsMap{
				"behavior": "NO_STORE",
			},
		},
		// HTTP/2
		{
			Name: "http2",
			Options: papi.RuleOptionsMap{
				"enabled": config.EnableHTTP2,
			},
		},
		// Real User Monitoring
		{
			Name: "realUserMonitoring",
			Options: papi.RuleOptionsMap{
				"enabled": true,
			},
		},
		// Allow all HTTP methods
		{
			Name: "allHttpInCacheHierarchy",
			Options: papi.RuleOptionsMap{
				"enabled": true,
			},
		},
	}

	// Merge with existing behaviors
	ruleTree.Rules.Behaviors = mergeBehaviors(ruleTree.Rules.Behaviors, behaviors)

	// Add custom behaviors if provided
	if len(config.CustomBehaviors) > 0 {
		ruleTree.Rules.Behaviors = append(ruleTree.Rules.Behaviors, config.CustomBehaviors...)
	}

	return ph.updateRuleTree(prop, config, ruleTree)
}

// configureMediaDelivery configures Adaptive Media Delivery behaviors
func (ph *PropertyHelper) configureMediaDelivery(prop *papi.Property, config PropertyConfig) error {
	fmt.Println(">> Configuring Adaptive Media Delivery behaviors")

	ruleTree, err := ph.getRuleTree(prop, config)
	if err != nil {
		return err
	}

	// Media Delivery behaviors
	behaviors := []papi.RuleBehavior{
		// Origin configuration
		{
			Name: "origin",
			Options: papi.RuleOptionsMap{
				"originType":         "CUSTOMER",
				"hostname":           config.OriginHostname,
				"forwardHostHeader":  "REQUEST_HOST_HEADER",
				"cacheKeyHostname":   "REQUEST_HOST_HEADER",
				"compress":           false, // Don't compress media
				"enableTrueClientIp": false,
				"verificationMode":   "PLATFORM_SETTINGS",
				"httpPort":           80,
				"httpsPort":          443,
			},
		},
		// CP Code
		{
			Name: "cpCode",
			Options: papi.RuleOptionsMap{
				"value": map[string]interface{}{
					"id": config.CPCode,
				},
			},
		},
		// Long cache for media
		{
			Name: "caching",
			Options: papi.RuleOptionsMap{
				"behavior": "MAX_AGE",
				"ttl":      "30d",
			},
		},
		// Adaptive media delivery
		{
			Name: "adaptiveMediaDelivery",
			Options: papi.RuleOptionsMap{
				"enabled": true,
			},
		},
		// Segment delivery
		{
			Name: "segmentedContentProtection",
			Options: papi.RuleOptionsMap{
				"enabled": false, // Enable if using token auth
			},
		},
		// HTTP/2 for better streaming
		{
			Name: "http2",
			Options: papi.RuleOptionsMap{
				"enabled": true,
			},
		},
	}

	// Merge with existing behaviors
	ruleTree.Rules.Behaviors = mergeBehaviors(ruleTree.Rules.Behaviors, behaviors)

	// Add custom behaviors if provided
	if len(config.CustomBehaviors) > 0 {
		ruleTree.Rules.Behaviors = append(ruleTree.Rules.Behaviors, config.CustomBehaviors...)
	}

	return ph.updateRuleTree(prop, config, ruleTree)
}

// Helper functions

func (ph *PropertyHelper) getRuleTree(prop *papi.Property, config PropertyConfig) (*papi.GetRuleTreeResponse, error) {
	return ph.client.GetRuleTree(ph.ctx, papi.GetRuleTreeRequest{
		PropertyID:      prop.PropertyID,
		PropertyVersion: prop.LatestVersion,
		ContractID:      config.ContractID,
		GroupID:         config.GroupID,
	})
}

func (ph *PropertyHelper) updateRuleTree(prop *papi.Property, config PropertyConfig, ruleTree *papi.GetRuleTreeResponse) error {
	_, err := ph.client.UpdateRuleTree(ph.ctx, papi.UpdateRulesRequest{
		PropertyID:      prop.PropertyID,
		PropertyVersion: prop.LatestVersion,
		ContractID:      config.ContractID,
		GroupID:         config.GroupID,
		Rules: papi.RulesUpdate{
			Rules: ruleTree.Rules,
		},
	})
	return err
}

// mergeBehaviors merges new behaviors with existing ones, avoiding duplicates
func mergeBehaviors(existing, new []papi.RuleBehavior) []papi.RuleBehavior {
	behaviorMap := make(map[string]papi.RuleBehavior)

	// Add existing behaviors to map
	for _, b := range existing {
		behaviorMap[strings.ToLower(b.Name)] = b
	}

	// Overwrite/add new behaviors
	for _, b := range new {
		behaviorMap[strings.ToLower(b.Name)] = b
	}

	// Convert back to slice
	result := make([]papi.RuleBehavior, 0, len(behaviorMap))
	for _, b := range behaviorMap {
		result = append(result, b)
	}

	return result
}

// AddHostnameToProperty adds a hostname to an existing property
func (ph *PropertyHelper) AddHostnameToProperty(prop *papi.Property, config PropertyConfig, edgeHostname string) error {
	fmt.Printf(">> Adding hostname %s to property %s\n", config.Domain, prop.PropertyName)

	// Get existing hostnames
	existingHostnames, err := ph.client.GetPropertyVersionHostnames(ph.ctx, papi.GetPropertyVersionHostnamesRequest{
		PropertyID:      prop.PropertyID,
		PropertyVersion: prop.LatestVersion,
		ContractID:      config.ContractID,
		GroupID:         config.GroupID,
	})
	if err != nil {
		return fmt.Errorf("failed to get existing hostnames: %w", err)
	}

	// Check if hostname already exists
	for _, h := range existingHostnames.Hostnames.Items {
		if h.CnameFrom == config.Domain {
			fmt.Printf(">> Hostname %s already configured\n", config.Domain)
			return nil
		}
	}

	// Add new hostname
	newHostnames := append(existingHostnames.Hostnames.Items, papi.Hostname{
		CnameType:            papi.HostnameCnameTypeEdgeHostname,
		CnameFrom:            config.Domain,
		CnameTo:              edgeHostname,
		CertProvisioningType: "CPS_MANAGED",
	})

	_, err = ph.client.UpdatePropertyVersionHostnames(ph.ctx, papi.UpdatePropertyVersionHostnamesRequest{
		PropertyID:      prop.PropertyID,
		PropertyVersion: prop.LatestVersion,
		ContractID:      config.ContractID,
		GroupID:         config.GroupID,
		Hostnames:       newHostnames,
	})
	if err != nil {
		return fmt.Errorf("failed to update hostnames: %w", err)
	}

	fmt.Printf(">> Hostname %s added successfully\n", config.Domain)
	return nil
}

// ActivateProperty activates a property to staging or production
func (ph *PropertyHelper) ActivateProperty(prop *papi.Property, config PropertyConfig, network papi.ActivationNetwork, notifyEmails []string) (*papi.CreateActivationResponse, error) {
	fmt.Printf(">> Activating property %s to %s\n", prop.PropertyName, network)

	activationResp, err := ph.client.CreateActivation(ph.ctx, papi.CreateActivationRequest{
		PropertyID: prop.PropertyID,
		ContractID: config.ContractID,
		GroupID:    config.GroupID,
		Activation: papi.Activation{
			PropertyVersion:        prop.LatestVersion,
			Network:                network,
			Note:                   fmt.Sprintf("Activation via property helper for %s", config.PropertyName),
			NotifyEmails:           notifyEmails,
			AcknowledgeAllWarnings: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create activation: %w", err)
	}

	fmt.Printf(">> Activation created: ID=%s, Link=%s\n", activationResp.ActivationID, activationResp.ActivationLink)
	return activationResp, nil
}
