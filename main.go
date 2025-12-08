package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cps"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/edgegrid"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/session"
)

// --- Configuration (replace with your real values) ---
const (
	EdgercPath    = "~/.edgerc"
	EdgercSection = "default"

	ContractID   = "ctr_1-12345" // Your Contract ID
	GroupID      = "grp_12345"   // Your Group ID
	ProductID    = "ion"     // Example product
	PropertyName = "my-api-gateway"
	UserDomain   = "api.example.com"
)

func main() {
	ctx := context.Background()

	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		log.Fatalf("edgegrid/session init failed: %v", err)
	}

	fmt.Println(">> Starting Onboarding Flow for:", UserDomain)

	// --- Clients ---
	cpsClient := cps.Client(sess)
	papiClient := papi.Client(sess)
	appsecClient := appsec.Client(sess)

	// ---------------------------------------------------------
	// STEP 1: Handle SSL / CPS (mocked)
	// ---------------------------------------------------------
	enrollmentID, err := getOrCreateEnrollment(ctx, cpsClient, UserDomain)
	if err != nil {
		log.Fatalf("CPS Error: %v", err)
	}
	fmt.Printf(">> Using Certificate Enrollment ID: %d\n", enrollmentID)

	// ---------------------------------------------------------
	// STEP 2: Configure Delivery (PAPI) – mostly mocked
	// ---------------------------------------------------------

	// A. Get or create a property (mocked Property structure for now)
	prop, err := getOrCreateProperty(ctx, papiClient, PropertyName)
	if err != nil {
		log.Fatalf("PAPI Property Error: %v", err)
	}

	// B. Create / ensure Edge Hostname (CreateEdgeHostname call mocked)
	edgeHostname, err := ensureEdgeHostname(ctx, papiClient, enrollmentID, UserDomain)
	if err != nil {
		log.Fatalf("EdgeHostname Error: %v", err)
	}
	fmt.Printf(">> Akamai Edge Hostname: %s\n", edgeHostname)

	// C. Update Property rules & hostnames (mocked)
	if err := updatePropertyRules(ctx, papiClient, prop, UserDomain, edgeHostname); err != nil {
		log.Fatalf("Rule Update Error: %v", err)
	}

	// D. Activate property to staging (mocked)
	if err := activateToStaging(ctx, papiClient, prop); err != nil {
		log.Fatalf("Activation Error: %v", err)
	}

	// ---------------------------------------------------------
	// STEP 3: Fake DNS Handler (CNAME Update)
	// ---------------------------------------------------------
	configureCNAMEInDNS(UserDomain, edgeHostname)

	// ---------------------------------------------------------
	// STEP 4: Configure WAF (AppSec) – mocked
	// ---------------------------------------------------------
	if err := onboardToWAF(ctx, appsecClient, UserDomain); err != nil {
		log.Fatalf("AppSec Error: %v", err)
	}

	fmt.Println(">> ✅ FULL ONBOARDING COMPLETE (with mocked Akamai API calls)")
}

// =========================================================
// Session / config helpers
// =========================================================

func newSession(edgercPath, section string) (session.Session, error) {
	expanded, err := expandTilde(edgercPath)
	if err != nil {
		return nil, fmt.Errorf("expand edgerc path: %w", err)
	}

	signer, err := edgegrid.New(
		edgegrid.WithFile(expanded),
		edgegrid.WithSection(section),
	)
	if err != nil {
		return nil, fmt.Errorf("edgegrid.New: %w", err)
	}

	sess, err := session.New(session.WithSigner(signer))
	if err != nil {
		return nil, fmt.Errorf("session.New: %w", err)
	}

	return sess, nil
}

func expandTilde(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if path == "~" {
		return home, nil
	}
	return filepath.Join(home, path[2:]), nil
}

// =========================================================
// 1. SSL / CPS Logic (mocked)
// =========================================================

func getOrCreateEnrollment(ctx context.Context, cpsClient cps.CPS, domain string) (int, error) {
	fmt.Println(">> Checking for existing SSL Enrollment for", domain)

	// List existing enrollments for the contract
	enrollmentsResp, err := cpsClient.ListEnrollments(ctx, cps.ListEnrollmentsRequest{
		ContractID: ContractID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to list enrollments: %w", err)
	}

	// Check if enrollment already exists for this domain
	for _, enrollment := range enrollmentsResp.Enrollments {
		// Check if the enrollment contains our domain in its common name or SANs
		if enrollment.CSR != nil && enrollment.CSR.CN == domain {
			fmt.Printf(">> Found existing enrollment ID: %d for domain: %s\n", enrollment.ID, domain)
			return enrollment.ID, nil
		}
		// Also check in SANs if available
		if enrollment.CSR != nil {
			for _, san := range enrollment.CSR.SANS {
				if san == domain {
					fmt.Printf(">> Found existing enrollment ID: %d for domain: %s (in SANs)\n", enrollment.ID, domain)
					return enrollment.ID, nil
				}
			}
		}
	}

	fmt.Printf(">> No existing enrollment found for %s\n", domain)
	fmt.Println(">> Note: Certificate enrollment creation requires additional configuration")
	fmt.Println(">>       (validation type, org details, etc.) and is typically done manually")
	fmt.Println(">>       or through a separate enrollment process.")
	fmt.Println(">> For this demo, please create an enrollment manually or provide an existing enrollment ID")

	return 0, fmt.Errorf("no enrollment found for domain %s - manual enrollment creation required", domain)
}

// =========================================================
// 2. PAPI Logic – mostly mocked so it compiles cleanly
// =========================================================

func getOrCreateProperty(ctx context.Context, papiClient papi.PAPI, name string) (*papi.Property, error) {
	fmt.Printf(">> Checking for existing property: %s\n", name)

	// Try to find existing property by name
	properties, err := papiClient.GetProperties(ctx, papi.GetPropertiesRequest{
		ContractID: ContractID,
		GroupID:    GroupID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list properties: %w", err)
	}

	// Search for property by name
	for _, prop := range properties.Properties.Items {
		if prop.PropertyName == name {
			fmt.Printf(">> Found existing property: %s (ID: %s, Version: %d)\n",
				prop.PropertyName, prop.PropertyID, prop.LatestVersion)
			return prop, nil
		}
	}

	// Property doesn't exist, create it
	fmt.Printf(">> Creating new property: %s\n", name)
	createResp, err := papiClient.CreateProperty(ctx, papi.CreatePropertyRequest{
		ContractID: ContractID,
		GroupID:    GroupID,
		Property: papi.PropertyCreate{
			ProductID:    ProductID,
			PropertyName: name,
			RuleFormat:   "latest",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create property: %w", err)
	}

	propertyID := createResp.PropertyID
	fmt.Printf(">> Property created with ID: %s\n", propertyID)

	// Fetch the newly created property to get full details
	propResp, err := papiClient.GetProperty(ctx, papi.GetPropertyRequest{
		ContractID: ContractID,
		GroupID:    GroupID,
		PropertyID: propertyID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created property: %w", err)
	}

	if propResp.Property == nil {
		return nil, fmt.Errorf("property response is nil")
	}

	return propResp.Property, nil
}

func ensureEdgeHostname(ctx context.Context, papiClient papi.PAPI, certEnrollmentID int, domain string) (string, error) {
	fmt.Printf(">> Checking for existing edge hostname for domain: %s\n", domain)

	// List existing edge hostnames
	edgeHostnamesResp, err := papiClient.GetEdgeHostnames(ctx, papi.GetEdgeHostnamesRequest{
		ContractID: ContractID,
		GroupID:    GroupID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to list edge hostnames: %w", err)
	}

	// Generate expected edge hostname (domain prefix)
	// Extract the first part of the domain as prefix
	domainPrefix := domain
	if idx := strings.Index(domain, "."); idx > 0 {
		domainPrefix = domain[:idx]
	}

	// Check if edge hostname already exists
	expectedEdgeHostname := domainPrefix + ".edgekey.net"
	for _, eh := range edgeHostnamesResp.EdgeHostnames.Items {
		if eh.Domain == expectedEdgeHostname || eh.DomainPrefix == domainPrefix {
			fmt.Printf(">> Found existing edge hostname: %s (ID: %s)\n", eh.Domain, eh.ID)
			return eh.Domain, nil
		}
	}

	// Create new edge hostname
	fmt.Printf(">> Creating edge hostname with prefix: %s\n", domainPrefix)
	createResp, err := papiClient.CreateEdgeHostname(ctx, papi.CreateEdgeHostnameRequest{
		ContractID: ContractID,
		GroupID:    GroupID,
		EdgeHostname: papi.EdgeHostnameCreate{
			ProductID:         ProductID,
			DomainPrefix:      domainPrefix,
			DomainSuffix:      "edgekey.net",
			SecureNetwork:     papi.EHSecureNetworkEnhancedTLS,
			IPVersionBehavior: papi.EHIPVersionV4,
			CertEnrollmentID:  certEnrollmentID,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create edge hostname: %w", err)
	}

	edgeHostname := domainPrefix + ".edgekey.net"
	fmt.Printf(">> Edge hostname created: %s (ID: %s)\n", edgeHostname, createResp.EdgeHostnameID)

	return edgeHostname, nil
}

func updatePropertyRules(ctx context.Context, papiClient papi.PAPI, prop *papi.Property, domain, edgeHostname string) error {
	fmt.Println(">> Updating property rules and hostnames")
	fmt.Printf("   PropertyID:   %s\n", prop.PropertyID)
	fmt.Printf("   Version:      %d\n", prop.LatestVersion)

	// Step 1: Update hostnames for the property version
	fmt.Printf(">> Adding hostname: %s -> %s\n", domain, edgeHostname)

	// First, get existing hostnames
	existingHostnames, err := papiClient.GetPropertyVersionHostnames(ctx, papi.GetPropertyVersionHostnamesRequest{
		PropertyID:      prop.PropertyID,
		PropertyVersion: prop.LatestVersion,
		ContractID:      ContractID,
		GroupID:         GroupID,
	})
	if err != nil {
		return fmt.Errorf("failed to get existing hostnames: %w", err)
	}

	// Check if hostname already exists
	hostnameExists := false
	for _, h := range existingHostnames.Hostnames.Items {
		if h.CnameFrom == domain {
			hostnameExists = true
			fmt.Printf(">> Hostname %s already configured\n", domain)
			break
		}
	}

	// Add new hostname if it doesn't exist
	if !hostnameExists {
		newHostnames := append(existingHostnames.Hostnames.Items, papi.Hostname{
			CnameType:            papi.HostnameCnameTypeEdgeHostname,
			CnameFrom:            domain,
			CnameTo:              edgeHostname,
			CertProvisioningType: "CPS_MANAGED",
		})

		_, err = papiClient.UpdatePropertyVersionHostnames(ctx, papi.UpdatePropertyVersionHostnamesRequest{
			PropertyID:      prop.PropertyID,
			PropertyVersion: prop.LatestVersion,
			ContractID:      ContractID,
			GroupID:         GroupID,
			Hostnames:       newHostnames,
		})
		if err != nil {
			return fmt.Errorf("failed to update hostnames: %w", err)
		}
		fmt.Printf(">> Hostname %s added successfully\n", domain)
	}

	// Step 2: Update rule tree to add origin behavior
	fmt.Println(">> Updating rule tree with origin configuration")

	ruleTree, err := papiClient.GetRuleTree(ctx, papi.GetRuleTreeRequest{
		PropertyID:      prop.PropertyID,
		PropertyVersion: prop.LatestVersion,
		ContractID:      ContractID,
		GroupID:         GroupID,
	})
	if err != nil {
		return fmt.Errorf("failed to get rule tree: %w", err)
	}

	// Add origin behavior if not already present
	originHostname := "server:IP"
	originExists := false
	for _, behavior := range ruleTree.Rules.Behaviors {
		if behavior.Name == "origin" {
			originExists = true
			fmt.Println(">> Origin behavior already configured")
			break
		}
	}

	if !originExists {
		// Add origin behavior to the default rule
		originBehavior := papi.RuleBehavior{
			Name: "origin",
			Options: papi.RuleOptionsMap{
				"originType":         "CUSTOMER",
				"hostname":           originHostname,
				"forwardHostHeader":  "REQUEST_HOST_HEADER",
				"cacheKeyHostname":   "REQUEST_HOST_HEADER",
				"compress":           true,
				"enableTrueClientIp": false,
				"verificationMode":   "PLATFORM_SETTINGS",
				"httpPort":           80,
				"httpsPort":          443,
			},
		}
		ruleTree.Rules.Behaviors = append(ruleTree.Rules.Behaviors, originBehavior)

		// Update the rule tree
		_, err = papiClient.UpdateRuleTree(ctx, papi.UpdateRulesRequest{
			PropertyID:      prop.PropertyID,
			PropertyVersion: prop.LatestVersion,
			ContractID:      ContractID,
			GroupID:         GroupID,
			Rules: papi.RulesUpdate{
				Rules: ruleTree.Rules,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to update rule tree: %w", err)
		}
		fmt.Printf(">> Origin configured: %s (HTTP 80 / HTTPS 443)\n", originHostname)
	}

	return nil
}

func activateToStaging(ctx context.Context, papiClient papi.PAPI, prop *papi.Property) error {
	fmt.Println(">> Activating property to STAGING")
	fmt.Printf("   PropertyID:   %s\n", prop.PropertyID)
	fmt.Printf("   Version:      %d\n", prop.LatestVersion)

	// Create activation request
	activationResp, err := papiClient.CreateActivation(ctx, papi.CreateActivationRequest{
		PropertyID: prop.PropertyID,
		ContractID: ContractID,
		GroupID:    GroupID,
		Activation: papi.Activation{
			PropertyVersion:        prop.LatestVersion,
			Network:                papi.ActivationNetworkStaging,
			Note:                   "Auto-onboard via Go SDK",
			NotifyEmails:           []string{"admin@example.com"},
			AcknowledgeAllWarnings: true,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create activation: %w", err)
	}

	fmt.Printf(">> Activation created with ID: %s\n", activationResp.ActivationID)
	fmt.Println(">> Note: Activation is now PENDING. You can monitor its status using GetActivation")
	fmt.Printf(">>       ActivationLink: %s\n", activationResp.ActivationLink)

	return nil
}

// =========================================================
// 3. DNS Mock Handler
// =========================================================

func configureCNAMEInDNS(domain, edgeHostname string) {
	fmt.Println("---------------------------------------------------------")
	fmt.Println(">> 📡 MOCK DNS PROVIDER API")
	fmt.Printf(">> ACTION: UPSERT CNAME RECORD\n")
	fmt.Printf(">> KEY:    %s\n", domain)
	fmt.Printf(">> VALUE:  %s\n", edgeHostname)
	fmt.Println(">> STATUS: 200 OK (Propagating...)")
	fmt.Println("---------------------------------------------------------")
}

// =========================================================
// 4. AppSec / WAF Logic (mocked)
// =========================================================

func onboardToWAF(ctx context.Context, appsecClient appsec.APPSEC, domain string) error {
	fmt.Println(">> Onboarding hostname into WAF/AppSec")
	fmt.Printf("   Hostname: %s\n", domain)

	// List available security configurations
	configsResp, err := appsecClient.GetConfigurations(ctx, appsec.GetConfigurationsRequest{})
	if err != nil {
		return fmt.Errorf("failed to list AppSec configurations: %w", err)
	}

	if len(configsResp.Configurations) == 0 {
		fmt.Println(">> Note: No AppSec configurations found")
		fmt.Println(">>       You need to create a security configuration first")
		fmt.Println(">>       This typically includes:")
		fmt.Println(">>       1. Creating a security configuration")
		fmt.Println(">>       2. Creating security policies")
		fmt.Println(">>       3. Configuring WAF rules and protections")
		return fmt.Errorf("no AppSec configurations available")
	}

	// Use the first available configuration
	config := configsResp.Configurations[0]
	fmt.Printf(">> Using AppSec Configuration: %s (ID: %d, Version: %d)\n",
		config.Name, config.ID, config.LatestVersion)

	// Check if hostname is already in the configuration
	hostnameExists := false
	for _, h := range config.ProductionHostnames {
		if h == domain {
			hostnameExists = true
			fmt.Printf(">> Hostname %s already protected by AppSec\n", domain)
			break
		}
	}

	if !hostnameExists {
		fmt.Printf(">> Note: To add hostname %s to AppSec configuration:\n", domain)
		fmt.Println(">>       1. Add hostname to selected hostnames in the configuration version")
		fmt.Println(">>       2. Create or update match targets to apply security policies")
		fmt.Println(">>       3. Activate the configuration version to staging/production")
		fmt.Println(">>       These steps typically require:")
		fmt.Println(">>          - UpdateSelectableHostnames")
		fmt.Println(">>          - CreateMatchTarget or UpdateMatchTarget")
		fmt.Println(">>          - CreateActivation")
	}

	fmt.Println(">> AppSec configuration check complete")
	fmt.Printf(">>    Config ID: %d\n", config.ID)
	fmt.Printf(">>    Latest Version: %d\n", config.LatestVersion)
	fmt.Printf(">>    Staging Version: %d\n", config.StagingVersion)
	fmt.Printf(">>    Production Version: %d\n", config.ProductionVersion)

	return nil
}
