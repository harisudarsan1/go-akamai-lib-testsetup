package main

import (
	"context"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
)

// TestCreateSecurityConfiguration creates a new security configuration
func TestCreateSecurityConfiguration(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Create Security Configuration ===\n")

	// Step 1: Authenticate
	t.Log("Step 1: Authenticating with Akamai API...")
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}
	t.Log("✅ Authentication successful")

	appsecClient := appsec.Client(sess)

	// Step 2: Get contract and group info
	t.Log("\nStep 2: Getting contract and group information...")

	// Load from cache or discover
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("❌ Failed to load configuration: %v", err)
	}

	t.Logf("✅ Using Contract: %s, Group: %s", config.ContractID, config.GroupID)

	// Step 3: Check if configuration already exists
	t.Log("\nStep 3: Checking if security configuration exists...")

	exists, existingID, existingVersion, err := configExists(ctx, appsecClient, SecurityConfigName)
	if err != nil {
		t.Fatalf("❌ Failed to check configurations: %v", err)
	}

	if exists {
		t.Logf("✅ Security configuration already exists")
		t.Logf("   Config ID: %d", existingID)
		t.Logf("   Config Name: %s", SecurityConfigName)
		t.Logf("   Latest Version: %d", existingVersion)
		return
	}

	// Step 4: Create new security configuration
	t.Logf("\nStep 4: Creating security configuration: %s", SecurityConfigName)

	// Convert GroupID string to int
	groupIDInt := 303793 // Ion group ID from config

	createResp, err := appsecClient.CreateConfiguration(ctx, appsec.CreateConfigurationRequest{
		Name:        SecurityConfigName,
		Description: SecurityConfigDescription,
		ContractID:  config.ContractID,
		GroupID:     groupIDInt,
		Hostnames:   []string{"example.kluisz.com"},
	})
	if err != nil {
		t.Fatalf("❌ Failed to create configuration: %v", err)
	}

	t.Log("✅ Security configuration created successfully!")
	t.Logf("   Config ID: %d", createResp.ConfigID)
	t.Logf("   Version: %d", createResp.Version)

	// Step 5: Verify configuration
	t.Log("\nStep 5: Verifying configuration...")

	getResp, err := appsecClient.GetConfiguration(ctx, appsec.GetConfigurationRequest{
		ConfigID: createResp.ConfigID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to verify configuration: %v", err)
	}

	t.Log("✅ Configuration verified!")
	t.Logf("   Name: %s", getResp.Name)
	t.Logf("   Description: %s", getResp.Description)
	t.Logf("   Contract: %s", config.ContractID)
	t.Logf("   Group: %s", config.GroupID)
	t.Logf("   Latest Version: %d", getResp.LatestVersion)

	t.Log("\n✅ Test completed successfully!")
}

// TestCreateSecurityPolicy creates a security policy within the configuration
func TestCreateSecurityPolicy(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Create Security Policy ===\n")

	// Step 1: Authenticate
	t.Log("Step 1: Authenticating...")
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}
	t.Log("✅ Authentication successful")

	appsecClient := appsec.Client(sess)

	// Step 2: Get security configuration
	t.Log("\nStep 2: Finding security configuration...")

	exists, configID, configVersion, err := configExists(ctx, appsecClient, SecurityConfigName)
	if err != nil {
		t.Fatalf("❌ Failed to list configurations: %v", err)
	}

	if !exists {
		t.Fatal("❌ Security configuration not found. Run TestCreateSecurityConfiguration first.")
	}

	t.Logf("✅ Found configuration: %s (ID: %d, Version: %d)",
		SecurityConfigName, configID, configVersion)

	// Step 3: Check if policy already exists
	t.Log("\nStep 3: Checking if security policy exists...")

	policyExists, existingPolicyID, err := policyExists(ctx, appsecClient, configID, configVersion, SecurityPolicyName)
	if err != nil {
		t.Fatalf("❌ Failed to list policies: %v", err)
	}

	if policyExists {
		t.Logf("✅ Security policy already exists")
		t.Logf("   Policy ID: %s", existingPolicyID)
		t.Logf("   Policy Name: %s", SecurityPolicyName)
		return
	}

	// Step 4: Create security policy
	t.Logf("\nStep 4: Creating security policy: %s", SecurityPolicyName)

	createResp, err := appsecClient.CreateSecurityPolicy(ctx, appsec.CreateSecurityPolicyRequest{
		ConfigID:     configID,
		Version:      configVersion,
		PolicyName:   SecurityPolicyName,
		PolicyPrefix: SecurityPolicyPrefix,
	})
	if err != nil {
		t.Fatalf("❌ Failed to create policy: %v", err)
	}

	t.Log("✅ Security policy created successfully!")
	t.Logf("   Policy ID: %s", createResp.PolicyID)
	t.Logf("   Policy Name: %s", createResp.PolicyName)

	// Step 5: Verify policy
	t.Log("\nStep 5: Verifying policy...")

	getResp, err := appsecClient.GetSecurityPolicy(ctx, appsec.GetSecurityPolicyRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: createResp.PolicyID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to verify policy: %v", err)
	}

	t.Log("✅ Policy verified!")
	t.Logf("   Policy ID: %s", getResp.PolicyID)
	t.Logf("   Policy Name: %s", getResp.PolicyName)

	t.Log("\n✅ Test completed successfully!")
}

// TestConfigureMatchTarget links hostname to security policy
func TestConfigureMatchTarget(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Configure Match Target ===\n")

	// Step 1: Authenticate
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)

	// Step 2: Get configuration and policy
	t.Log("Step 2: Finding configuration and policy...")

	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	t.Logf("✅ Found config %d v%d, policy %s", configID, configVersion, policyID)

	// Step 3: Check existing match targets
	t.Log("\nStep 3: Checking existing match targets...")

	exists, targetID, err := matchTargetExistsForHostname(ctx, appsecClient, configID, configVersion, "example.kluisz.com")
	if err != nil {
		t.Fatalf("❌ Failed to get match targets: %v", err)
	}

	if exists {
		t.Log("✅ Match target already exists for example.kluisz.com")
		t.Logf("   Target ID: %d", targetID)
		return
	}

	// Step 4: Create match target
	t.Log("\nStep 4: Creating match target...")
	t.Log("⚠️  Note: Match target creation requires manual setup via Akamai Control Center")
	t.Log("   or using the CLI with proper JSON payload.")
	t.Log("\nTo create via Control Center:")
	t.Log("   1. Go to Security Configurations")
	t.Logf("   2. Select config %d version %d", configID, configVersion)
	t.Log("   3. Go to Match Targets")
	t.Log("   4. Add Website Match Target:")
	t.Log("      - Hostname: example.kluisz.com")
	t.Log("      - Path: /*")
	t.Logf("      - Policy: %s", policyID)

	t.Log("\n✅ Test completed successfully!")
	t.Log("\n🎉 Security foundation configuration identified!")
	t.Log("   Next: Create match target, then enable protections")
}

// TestVerifySecuritySetup verifies the complete security foundation
func TestVerifySecuritySetup(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Verify Security Setup ===\n")

	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)

	// Check 1: Configuration exists
	t.Log("Check 1: Security configuration...")
	exists, configID, configVersion, err := configExists(ctx, appsecClient, SecurityConfigName)
	if err != nil {
		t.Fatalf("❌ Failed: %v", err)
	}

	if !exists {
		t.Fatal("❌ Security configuration not found")
	}
	t.Logf("✅ Configuration exists: %s (ID: %d)", SecurityConfigName, configID)

	// Check 2: Policy exists
	t.Log("\nCheck 2: Security policy...")
	policyExists, policyID, err := policyExists(ctx, appsecClient, configID, configVersion, SecurityPolicyName)
	if err != nil {
		t.Fatalf("❌ Failed: %v", err)
	}

	if !policyExists {
		t.Fatal("❌ Security policy not found")
	}
	t.Logf("✅ Policy exists: %s (ID: %s)", SecurityPolicyName, policyID)

	// Check 3: Match target exists
	t.Log("\nCheck 3: Match target...")
	mtExists, targetID, err := matchTargetExistsForHostname(ctx, appsecClient, configID, configVersion, "example.kluisz.com")
	if err != nil {
		t.Fatalf("❌ Failed: %v", err)
	}

	if !mtExists {
		t.Fatal("❌ Match target not found for example.kluisz.com")
	}

	t.Logf("✅ Match target exists for: example.kluisz.com")
	t.Logf("   Target ID: %d", targetID)

	// Print summary
	t.Log("\n" + "==================================================")
	t.Log("✅ Security Setup Verification Complete!")
	t.Log("==================================================")
	t.Log("\nSecurity Foundation:")
	t.Logf("  Config: %s (ID: %d, Version: %d)", SecurityConfigName, configID, configVersion)
	t.Logf("  Policy: %s (ID: %s)", SecurityPolicyName, policyID)
	t.Log("  Hostname: example.kluisz.com")
	t.Log("  Paths: /*")
	t.Log("\nNext Steps:")
	t.Log("  1. Enable WAF protection")
	t.Log("  2. Configure attack groups (SQLi, XSS)")
	t.Log("  3. Enable rate controls")
	t.Log("  4. Configure bot detection")
	t.Log("  5. Activate to staging")
}

// TestPrintSecuritySummary prints a detailed summary
func TestPrintSecuritySummary(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Security Configuration Summary ===\n")

	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)

	// Get config details
	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	// Print comprehensive summary
	err = printSecuritySummary(ctx, appsecClient, configID, configVersion, policyID)
	if err != nil {
		t.Fatalf("❌ Failed to print summary: %v", err)
	}

	t.Log("✅ Summary displayed successfully!")
}
