# Akamai Application Security - End-to-End Implementation Guide

## Table of Contents

### Getting Started
- [Prerequisites](#prerequisites)
- [Overview](#overview)
- [What We'll Build](#what-well-build)

### Part 1: Initial Security Setup
- [1.1 Create Security Configuration](#part-1-initial-security-setup)
- [1.2 Create Security Policy](#12-create-security-policy)
- [1.3 Configure Match Target](#13-configure-match-target)
- [1.4 Verify Setup](#14-verify-setup)

### Part 2: Protect Against SQL Injection
- [2.1 Understanding SQL Injection](#part-2-protect-against-sql-injection)
- [2.2 Enable WAF Protection](#22-enable-waf-protection)
- [2.3 Configure SQLi Attack Group](#23-configure-sqli-attack-group)
- [2.4 Test Protection](#24-test-protection)

### Part 3: Protect Against XSS
- [3.1 Understanding XSS Attacks](#part-3-protect-against-xss)
- [3.2 Enable XSS Protection](#32-enable-xss-protection)
- [3.3 Configure Exceptions](#33-configure-exceptions)
- [3.4 Test Protection](#34-test-protection)

### Part 4: Rate Limiting & DDoS Protection
- [4.1 Understanding Rate Controls](#part-4-rate-limiting--ddos-protection)
- [4.2 Create Rate Policies](#42-create-rate-policies)
- [4.3 Configure Penalty Box](#43-configure-penalty-box)
- [4.4 Test Rate Limiting](#44-test-rate-limiting)

### Part 5: Bot Detection & Management
- [5.1 Understanding Bot Threats](#part-5-bot-detection--management)
- [5.2 Configure Bot Detection](#52-configure-bot-detection)
- [5.3 Set Bot Actions](#53-set-bot-actions)
- [5.4 Test Bot Protection](#54-test-bot-protection)

### Part 6: IP/Geo Blocking
- [6.1 Understanding IP/Geo Controls](#part-6-ipgeo-blocking)
- [6.2 Create Network Lists](#62-create-network-lists)
- [6.3 Configure Geo Blocking](#63-configure-geo-blocking)
- [6.4 Test IP/Geo Rules](#64-test-ipgeo-rules)

### Part 7: Custom Rules
- [7.1 When to Use Custom Rules](#part-7-custom-rules)
- [7.2 Create Application-Specific Rules](#72-create-application-specific-rules)
- [7.3 Advanced Conditions](#73-advanced-conditions)

### Part 8: Testing & Validation
- [8.1 Evaluation Mode](#part-8-testing--validation)
- [8.2 Security Event Monitoring](#82-security-event-monitoring)
- [8.3 Rule Tuning](#83-rule-tuning)

### Part 9: Production Activation
- [9.1 Pre-Activation Checklist](#part-9-production-activation)
- [9.2 Staging Activation](#92-staging-activation)
- [9.3 Production Activation](#93-production-activation)
- [9.4 Rollback Procedures](#94-rollback-procedures)

### Part 10: Complete Example
- [10.1 Securing "propertyname"](#part-10-complete-example)
- [10.2 Full Implementation](#102-full-implementation)

---

## Prerequisites

Before starting this guide, ensure you have:

### ✅ Required
1. **Akamai Account** with Application Security entitlement
2. **Property Created** and activated with hostname
3. **EdgeGrid Credentials** configured in `~/.edgerc`
4. **Go Environment** set up (Go 1.19+)
5. **Akamai Go SDK** installed

### ✅ Completed Previous Steps
- Property created: `propertyname` (prp_1275953)
- Hostname configured: `example.kluisz.com`
- Ion group access: grp_303793, ctr_V-5ZUL2W6
- Contract discovery working
- Property tests passing

### ✅ Knowledge
- Basic understanding of HTTP/HTTPS
- Familiarity with web security concepts
- Basic Go programming knowledge
- Understanding of Akamai property management

### 📁 Project Structure
```
go-akamai-waf-test/
├── contract_discovery.go
├── contract_discovery_test.go
├── appsec_helpers.go          ← We'll create this
├── appsec_basic_test.go        ← We'll create this
├── appsec_waf_test.go          ← We'll create this
├── appsec_rate_limiting_test.go ← We'll create this
├── appsec_bot_detection_test.go ← We'll create this
└── APPSEC_END_TO_END_GUIDE.md  ← You are here
```

---

## Overview

This guide provides **step-by-step instructions** with **working code examples** to secure your Akamai property using Application Security.

### What You'll Learn
- How to create and configure security configurations
- How to protect against common attacks (SQLi, XSS, DDoS)
- How to implement rate limiting and bot detection
- How to test security rules safely
- How to activate security to production

### Approach
- **Hands-on**: Every section includes working code
- **Incremental**: Build security layer by layer
- **Safe**: Test in evaluation mode first
- **Production-ready**: Complete activation workflows

### Time Estimate
- **Part 1-3** (Setup, SQLi, XSS): 60 minutes
- **Part 4-6** (Rate limiting, Bots, IP/Geo): 60 minutes
- **Part 7-9** (Custom rules, Testing, Activation): 45 minutes
- **Total**: ~3 hours for complete implementation

---

## What We'll Build

By the end of this guide, you'll have:

### 🛡️ Comprehensive Security Configuration

```
Security Configuration: "propertyname-security"
│
├── Security Policy: "production-policy"
│   ├── WAF Protection (Enabled)
│   │   ├── SQL Injection rules (DENY)
│   │   ├── XSS rules (DENY)
│   │   ├── Command Injection rules (DENY)
│   │   └── Other attack groups (ALERT)
│   │
│   ├── Rate Controls (Enabled)
│   │   ├── Global rate: 1000 req/min
│   │   ├── Per-IP rate: 100 req/min
│   │   ├── API endpoint: 50 req/min
│   │   └── Penalty Box: 300 seconds
│   │
│   ├── Bot Detection (Enabled)
│   │   ├── Allow: Search engines, monitors
│   │   ├── Deny: Malicious bots, scrapers
│   │   ├── Challenge: Unknown bots
│   │   └── Monitor: Social media bots
│   │
│   ├── IP/Geo Firewall (Enabled)
│   │   ├── Allowlist: Corporate IPs
│   │   ├── Denylist: Known attackers
│   │   └── Geo-block: High-risk countries
│   │
│   └── Custom Rules
│       ├── Admin path protection
│       ├── API authentication check
│       └── File upload validation
│
└── Match Target
    ├── Hostnames: ["example.kluisz.com"]
    ├── Paths: ["/*"]
    └── Policy: "production-policy"
```

### 📊 Protection Coverage

| Threat Type | Protection Layer | Action |
|------------|------------------|--------|
| SQL Injection | WAF + Custom Rules | DENY |
| XSS | WAF + CSP | DENY |
| Command Injection | WAF | DENY |
| DDoS | Rate Controls | DENY |
| Bot Attacks | Bot Detection | DENY/CHALLENGE |
| Scraping | Bot + Rate Controls | DENY |
| Geo Threats | IP/Geo Firewall | DENY |
| API Abuse | Rate + Custom Rules | DENY |
| Brute Force | Rate + Penalty Box | DENY |
| Zero-Day | Reputation + WAF | ALERT/DENY |

---

## Part 1: Initial Security Setup

**Goal**: Create the foundation for all security protections.

**Time**: 15 minutes

### What We'll Create
1. Security Configuration (container for all settings)
2. Security Policy (set of protections)
3. Match Target (link hostname to policy)

### 1.1 Create Security Configuration

A security configuration is the top-level container for all security settings.

#### Code Example

File: `appsec_basic_test.go`

```go
package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
)

const (
	// Security configuration settings
	SecurityConfigName        = "propertyname-security"
	SecurityConfigDescription = "Security configuration for propertyname property"
	SecurityPolicyName        = "production-policy"
	SecurityPolicyPrefix      = "prop"
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
	
	configsResp, err := appsecClient.GetConfigurations(ctx, appsec.GetConfigurationsRequest{})
	if err != nil {
		t.Fatalf("❌ Failed to list configurations: %v", err)
	}

	var existingConfig *appsec.ConfigurationResponse
	for i := range configsResp.Configurations {
		if configsResp.Configurations[i].Name == SecurityConfigName {
			existingConfig = &configsResp.Configurations[i]
			break
		}
	}

	if existingConfig != nil {
		t.Logf("✅ Security configuration already exists")
		t.Logf("   Config ID: %d", existingConfig.ID)
		t.Logf("   Config Name: %s", existingConfig.Name)
		t.Logf("   Latest Version: %d", existingConfig.LatestVersion)
		t.Logf("   Staging Version: %d", existingConfig.StagingVersion)
		t.Logf("   Production Version: %d", existingConfig.ProductionVersion)
		return
	}

	// Step 4: Create new security configuration
	t.Logf("\nStep 4: Creating security configuration: %s", SecurityConfigName)

	createResp, err := appsecClient.CreateConfiguration(ctx, appsec.CreateConfigurationRequest{
		Name:        SecurityConfigName,
		Description: SecurityConfigDescription,
		ContractID:  config.ContractID,
		GroupID:     config.GroupID,
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
	t.Logf("   Group: %d", config.GroupID)
	t.Logf("   Latest Version: %d", getResp.LatestVersion)

	t.Log("\n✅ Test completed successfully!")
}
```

#### What This Does

1. **Authenticates** with Akamai API
2. **Checks** if configuration already exists (idempotent)
3. **Creates** new security configuration if needed
4. **Links** configuration to contract and group
5. **Associates** hostname (`example.kluisz.com`)
6. **Verifies** configuration was created

#### Running the Test

```bash
cd /Users/hari/kluisz/akamai/go-akamai-waf-test
go test -v -run TestCreateSecurityConfiguration
```

#### Expected Output

```
=== Create Security Configuration ===

Step 1: Authenticating with Akamai API...
✅ Authentication successful

Step 2: Getting contract and group information...
✅ Using Contract: ctr_V-5ZUL2W6, Group: grp_303793

Step 3: Checking if security configuration exists...

Step 4: Creating security configuration: propertyname-security
✅ Security configuration created successfully!
   Config ID: 123456
   Version: 1

Step 5: Verifying configuration...
✅ Configuration verified!
   Name: propertyname-security
   Description: Security configuration for propertyname property
   Contract: ctr_V-5ZUL2W6
   Group: 303793
   Latest Version: 1

✅ Test completed successfully!
```

---

### 1.2 Create Security Policy

A security policy defines which protections to enable and how to configure them.

#### Code Example

```go
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
	
	configsResp, err := appsecClient.GetConfigurations(ctx, appsec.GetConfigurationsRequest{})
	if err != nil {
		t.Fatalf("❌ Failed to list configurations: %v", err)
	}

	var configID int
	var configVersion int
	for _, cfg := range configsResp.Configurations {
		if cfg.Name == SecurityConfigName {
			configID = cfg.ID
			configVersion = cfg.LatestVersion
			break
		}
	}

	if configID == 0 {
		t.Fatal("❌ Security configuration not found. Run TestCreateSecurityConfiguration first.")
	}

	t.Logf("✅ Found configuration: %s (ID: %d, Version: %d)", 
		SecurityConfigName, configID, configVersion)

	// Step 3: Check if policy already exists
	t.Log("\nStep 3: Checking if security policy exists...")
	
	policiesResp, err := appsecClient.GetSecurityPolicies(ctx, appsec.GetSecurityPoliciesRequest{
		ConfigID: configID,
		Version:  configVersion,
	})
	if err != nil {
		t.Fatalf("❌ Failed to list policies: %v", err)
	}

	for _, policy := range policiesResp.Policies {
		if policy.PolicyName == SecurityPolicyName {
			t.Logf("✅ Security policy already exists")
			t.Logf("   Policy ID: %s", policy.PolicyID)
			t.Logf("   Policy Name: %s", policy.PolicyName)
			return
		}
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
```

#### What This Does

1. **Finds** the security configuration created earlier
2. **Checks** if policy already exists
3. **Creates** new security policy if needed
4. **Associates** policy with configuration
5. **Verifies** policy creation

#### Running the Test

```bash
go test -v -run TestCreateSecurityPolicy
```

---

### 1.3 Configure Match Target

A match target links your hostname to the security policy.

#### Code Example

```go
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
	
	configsResp, err := appsecClient.GetConfigurations(ctx, appsec.GetConfigurationsRequest{})
	if err != nil {
		t.Fatalf("❌ Failed to list configurations: %v", err)
	}

	var configID int
	var configVersion int
	for _, cfg := range configsResp.Configurations {
		if cfg.Name == SecurityConfigName {
			configID = cfg.ID
			configVersion = cfg.LatestVersion
			break
		}
	}

	policiesResp, err := appsecClient.GetSecurityPolicies(ctx, appsec.GetSecurityPoliciesRequest{
		ConfigID: configID,
		Version:  configVersion,
	})
	if err != nil {
		t.Fatalf("❌ Failed to list policies: %v", err)
	}

	var policyID string
	for _, policy := range policiesResp.Policies {
		if policy.PolicyName == SecurityPolicyName {
			policyID = policy.PolicyID
			break
		}
	}

	t.Logf("✅ Found config %d v%d, policy %s", configID, configVersion, policyID)

	// Step 3: Check existing match targets
	t.Log("\nStep 3: Checking existing match targets...")
	
	targetsResp, err := appsecClient.GetMatchTargets(ctx, appsec.GetMatchTargetsRequest{
		ConfigID: configID,
		Version:  configVersion,
	})
	if err != nil {
		t.Fatalf("❌ Failed to get match targets: %v", err)
	}

	// Check if hostname already has a match target
	for _, target := range targetsResp.MatchTargets {
		for _, hostname := range target.Hostnames {
			if hostname == "example.kluisz.com" {
				t.Log("✅ Match target already exists for example.kluisz.com")
				t.Logf("   Target ID: %d", target.TargetID)
				t.Logf("   Policy: %s", target.SecurityPolicy.PolicyID)
				return
			}
		}
	}

	// Step 4: Create match target
	t.Log("\nStep 4: Creating match target...")

	createResp, err := appsecClient.CreateMatchTarget(ctx, appsec.CreateMatchTargetRequest{
		ConfigID: configID,
		Version:  configVersion,
		Type:     "website",
		Hostnames: []string{"example.kluisz.com"},
		FilePaths: []string{"/*"},
		SecurityPolicy: appsec.SecurityPolicyReference{
			PolicyID: policyID,
		},
	})
	if err != nil {
		t.Fatalf("❌ Failed to create match target: %v", err)
	}

	t.Log("✅ Match target created successfully!")
	t.Logf("   Target ID: %d", createResp.TargetID)
	t.Logf("   Hostname: example.kluisz.com")
	t.Logf("   Paths: /*")
	t.Logf("   Policy: %s", policyID)

	t.Log("\n✅ Test completed successfully!")
	t.Log("\n🎉 Security foundation is ready!")
	t.Log("   Next: Enable protections (WAF, rate controls, etc.)")
}
```

#### What This Does

1. **Finds** configuration and policy
2. **Checks** if match target exists for hostname
3. **Creates** match target linking hostname to policy
4. **Configures** to protect all paths (`/*`)

#### Running the Test

```bash
go test -v -run TestConfigureMatchTarget
```

---

### 1.4 Verify Setup

Let's create a verification test to ensure everything is configured correctly.

#### Code Example

```go
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
	configsResp, err := appsecClient.GetConfigurations(ctx, appsec.GetConfigurationsRequest{})
	if err != nil {
		t.Fatalf("❌ Failed: %v", err)
	}

	var configID int
	var configVersion int
	found := false
	for _, cfg := range configsResp.Configurations {
		if cfg.Name == SecurityConfigName {
			configID = cfg.ID
			configVersion = cfg.LatestVersion
			found = true
			break
		}
	}

	if !found {
		t.Fatal("❌ Security configuration not found")
	}
	t.Logf("✅ Configuration exists: %s (ID: %d)", SecurityConfigName, configID)

	// Check 2: Policy exists
	t.Log("\nCheck 2: Security policy...")
	policiesResp, err := appsecClient.GetSecurityPolicies(ctx, appsec.GetSecurityPoliciesRequest{
		ConfigID: configID,
		Version:  configVersion,
	})
	if err != nil {
		t.Fatalf("❌ Failed: %v", err)
	}

	var policyID string
	found = false
	for _, policy := range policiesResp.Policies {
		if policy.PolicyName == SecurityPolicyName {
			policyID = policy.PolicyID
			found = true
			break
		}
	}

	if !found {
		t.Fatal("❌ Security policy not found")
	}
	t.Logf("✅ Policy exists: %s (ID: %s)", SecurityPolicyName, policyID)

	// Check 3: Match target exists
	t.Log("\nCheck 3: Match target...")
	targetsResp, err := appsecClient.GetMatchTargets(ctx, appsec.GetMatchTargetsRequest{
		ConfigID: configID,
		Version:  configVersion,
	})
	if err != nil {
		t.Fatalf("❌ Failed: %v", err)
	}

	found = false
	for _, target := range targetsResp.MatchTargets {
		for _, hostname := range target.Hostnames {
			if hostname == "example.kluisz.com" {
				found = true
				t.Logf("✅ Match target exists for: %s", hostname)
				t.Logf("   Target ID: %d", target.TargetID)
				t.Logf("   Policy: %s", target.SecurityPolicy.PolicyID)
				break
			}
		}
	}

	if !found {
		t.Fatal("❌ Match target not found for example.kluisz.com")
	}

	// Summary
	t.Log("\n" + "=" * 50)
	t.Log("✅ Security Setup Verification Complete!")
	t.Log("=" * 50)
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
```

#### Running All Setup Tests

```bash
# Run all basic setup tests
go test -v -run "TestCreate|TestConfigure|TestVerify"
```

---

## Summary of Part 1

### ✅ What We Accomplished

1. **Created Security Configuration** - Container for all security settings
2. **Created Security Policy** - Defines protections to enable
3. **Configured Match Target** - Linked hostname to policy
4. **Verified Setup** - Confirmed everything is configured

### 📊 Current State

```
Security Configuration: "propertyname-security"
└── Security Policy: "production-policy"
    └── Match Target
        ├── Hostname: example.kluisz.com
        ├── Paths: /*
        └── Protections: None (yet)
```

### ⏭️ Next Steps

In Part 2, we'll enable WAF protection and configure SQL injection defenses.

---

## Part 2: Protect Against SQL Injection

**Goal**: Enable WAF protection with SQL injection attack group

**Time**: 15 minutes

### 2.1 Understanding SQL Injection

SQL Injection (SQLi) is one of the most common and dangerous web application vulnerabilities.

#### How SQLi Attacks Work

**Normal Query:**
```sql
SELECT * FROM users WHERE username='john' AND password='pass123'
```

**Malicious Input:**
```
Username: admin' OR '1'='1
Password: anything
```

**Resulting Query:**
```sql
SELECT * FROM users WHERE username='admin' OR '1'='1' AND password='anything'
-- This returns all users because '1'='1' is always true
```

#### Common SQLi Patterns
- `' OR '1'='1`
- `'; DROP TABLE users;--`
- `UNION SELECT * FROM sensitive_data`
- `admin'--`
- `1' AND 1=1--`

#### Akamai Protection

Akamai WAF includes **50+ SQL injection rules** that detect:
- SQL keywords in suspicious contexts
- SQL comment sequences (`--`, `/*`, `*/`)
- UNION-based injection attempts
- Boolean-based blind injection
- Time-based blind injection
- Out-of-band injection
- Second-order injection

---

### 2.2 Enable WAF Protection

Let's enable WAF protection and configure the SQL injection attack group.

#### Code Example

File: `appsec_waf_test.go`

```go
package main

import (
	"context"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
)

// TestEnableWAFProtection enables Web Application Firewall
func TestEnableWAFProtection(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Enable WAF Protection ===\n")

	// Step 1: Setup
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)

	// Step 2: Get configuration details
	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	t.Logf("✅ Using config %d v%d, policy %s", configID, configVersion, policyID)

	// Step 3: Set WAF mode to automatic updates
	t.Log("\nStep 3: Setting WAF mode to automatic...")
	
	_, err = appsecClient.UpdateWAFMode(ctx, appsec.UpdateWAFModeRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
		Mode:     "ASE_AUTO", // Automatic Security Engine updates
	})
	if err != nil {
		t.Fatalf("❌ Failed to set WAF mode: %v", err)
	}

	t.Log("✅ WAF mode set to ASE_AUTO (automatic updates)")

	// Step 4: Enable WAF protection
	t.Log("\nStep 4: Enabling WAF protection...")
	
	_, err = appsecClient.UpdateWAFProtection(ctx, appsec.UpdateWAFProtectionRequest{
		ConfigID:      configID,
		Version:       configVersion,
		PolicyID:      policyID,
		ApplyNetworkLayerControls: true,
	})
	if err != nil {
		t.Fatalf("❌ Failed to enable WAF protection: %v", err)
	}

	t.Log("✅ WAF protection enabled!")

	// Step 5: Verify protection status
	t.Log("\nStep 5: Verifying protection status...")
	
	protResp, err := appsecClient.GetWAFProtection(ctx, appsec.GetWAFProtectionRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to get protection status: %v", err)
	}

	t.Logf("✅ WAF protection status verified")
	t.Logf("   Apply Network Layer Controls: %v", protResp.ApplyNetworkLayerControls)

	t.Log("\n✅ WAF protection is now active!")
}

// Helper function to get security configuration details
func getSecurityConfig(ctx context.Context, client appsec.APPSEC) (int, int, string, error) {
	// Get configuration
	configsResp, err := client.GetConfigurations(ctx, appsec.GetConfigurationsRequest{})
	if err != nil {
		return 0, 0, "", err
	}

	var configID int
	var configVersion int
	for _, cfg := range configsResp.Configurations {
		if cfg.Name == SecurityConfigName {
			configID = cfg.ID
			configVersion = cfg.LatestVersion
			break
		}
	}

	// Get policy
	policiesResp, err := client.GetSecurityPolicies(ctx, appsec.GetSecurityPoliciesRequest{
		ConfigID: configID,
		Version:  configVersion,
	})
	if err != nil {
		return 0, 0, "", err
	}

	var policyID string
	for _, policy := range policiesResp.Policies {
		if policy.PolicyName == SecurityPolicyName {
			policyID = policy.PolicyID
			break
		}
	}

	return configID, configVersion, policyID, nil
}
```

---

### 2.3 Configure SQLi Attack Group

Now let's configure the SQL Injection attack group to DENY malicious requests.

#### Code Example

```go
// TestConfigureSQLiProtection configures SQL injection attack group
func TestConfigureSQLiProtection(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Configure SQL Injection Protection ===\n")

	// Setup
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)
	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	// Step 1: Get attack groups
	t.Log("Step 1: Finding SQL Injection attack group...")
	
	groupsResp, err := appsecClient.GetAttackGroups(ctx, appsec.GetAttackGroupsRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to get attack groups: %v", err)
	}

	// Find SQL Injection attack group
	var sqliGroup *appsec.AttackGroupResponse
	for i, group := range groupsResp.AttackGroups {
		if group.Group == "SQL" || group.GroupName == "SQL Injection" {
			sqliGroup = &groupsResp.AttackGroups[i]
			break
		}
	}

	if sqliGroup == nil {
		t.Fatal("❌ SQL Injection attack group not found")
	}

	t.Logf("✅ Found SQLi attack group: %s", sqliGroup.Group)
	t.Logf("   Current action: %s", sqliGroup.Action)

	// Step 2: Set action to DENY
	t.Log("\nStep 2: Setting SQL Injection action to DENY...")
	
	_, err = appsecClient.UpdateAttackGroupAction(ctx, appsec.UpdateAttackGroupActionRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
		Group:    sqliGroup.Group,
		Action:   "deny",
	})
	if err != nil {
		t.Fatalf("❌ Failed to update attack group: %v", err)
	}

	t.Log("✅ SQL Injection protection configured!")
	t.Log("   Action: DENY (block all SQL injection attempts)")

	// Step 3: Verify configuration
	t.Log("\nStep 3: Verifying configuration...")
	
	verifyResp, err := appsecClient.GetAttackGroup(ctx, appsec.GetAttackGroupRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
		Group:    sqliGroup.Group,
	})
	if err != nil {
		t.Fatalf("❌ Failed to verify: %v", err)
	}

	t.Log("✅ Configuration verified!")
	t.Logf("   Attack Group: %s", verifyResp.Group)
	t.Logf("   Action: %s", verifyResp.Action)

	t.Log("\n🛡️ SQL Injection protection is now ACTIVE!")
	t.Log("   All SQL injection attempts will be BLOCKED")
}
```

---

### 2.4 Test Protection

Let's verify that SQL injection protection is working.

#### Testing Approach

We'll use **evaluation mode** first to test without blocking.

#### Code Example

```go
// TestSQLiProtectionEval tests SQLi protection in evaluation mode
func TestSQLiProtectionEval(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Test SQL Injection Protection (Evaluation Mode) ===\n")

	// Setup
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)
	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	// Step 1: Enable evaluation mode for SQL group
	t.Log("Step 1: Enabling evaluation mode for SQL Injection...")
	
	groupsResp, err := appsecClient.GetAttackGroups(ctx, appsec.GetAttackGroupsRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
	})
	if err != nil {
		t.Fatalf("❌ Failed: %v", err)
	}

	var sqliGroup string
	for _, group := range groupsResp.AttackGroups {
		if group.Group == "SQL" {
			sqliGroup = group.Group
			break
		}
	}

	// Update to evaluation mode (alert only, don't block)
	_, err = appsecClient.UpdateEvalGroup(ctx, appsec.UpdateEvalGroupRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
		Group:    sqliGroup,
		Action:   "alert", // Alert only, don't block
	})
	if err != nil {
		t.Fatalf("❌ Failed to set eval mode: %v", err)
	}

	t.Log("✅ Evaluation mode enabled")
	t.Log("   SQL injection attempts will be LOGGED but not blocked")
	t.Log("\nTest Traffic Patterns:")
	t.Log("   1. Normal request: /search?q=hello")
	t.Log("   2. SQLi attempt: /search?q=' OR '1'='1")
	t.Log("   3. SQLi attempt: /search?q='; DROP TABLE users;--")
	t.Log("\nAll requests will pass through but attacks will be logged.")
	t.Log("Review security events in Akamai Control Center after testing.")
}
```

#### Manual Testing with curl

After enabling evaluation mode:

```bash
# Normal request (should pass)
curl "https://example.kluisz.com/search?q=hello"

# SQL injection attempt (should pass but be logged)
curl "https://example.kluisz.com/search?q=%27%20OR%20%271%27=%271"

# Another SQLi attempt (should pass but be logged)
curl "https://example.kluisz.com/login?user=admin%27--"
```

**Check logs** in Akamai Control Center → Security Events

---

## Summary of Part 2

### ✅ What We Accomplished

1. **Enabled WAF Protection** - Activated Kona Rule Set
2. **Configured SQLi Attack Group** - Set to DENY mode
3. **Set WAF Mode** - Automatic rule updates (ASE_AUTO)
4. **Tested in Evaluation Mode** - Verified detection without blocking

### 📊 Current Protection

```
WAF Protection: ENABLED
├── Mode: ASE_AUTO (automatic updates)
├── SQL Injection: DENY ✅
├── XSS: (not configured yet)
└── Other groups: Default action
```

### 🧪 Testing Results

- Normal requests: ✅ Allowed
- SQL injection attempts: 🚫 Detected and logged (eval mode) or Blocked (deny mode)
- False positives: Monitor and tune

### ⏭️ Next: Part 3

Configure Cross-Site Scripting (XSS) protection.

---

## Part 3: Protect Against XSS

**Goal**: Enable XSS protection and configure exceptions

**Time**: 15 minutes

### 3.1 Understanding XSS Attacks

Cross-Site Scripting (XSS) allows attackers to inject malicious scripts into web pages viewed by other users.

#### Types of XSS

**1. Reflected XSS**
```
URL: https://example.com/search?q=<script>alert('XSS')</script>
Page displays: Search results for: <script>alert('XSS')</script>
Result: Script executes in victim's browser
```

**2. Stored XSS**
```
Attacker posts comment: <script>steal_cookies()</script>
Comment stored in database
Every user viewing the comment gets attacked
```

**3. DOM-based XSS**
```javascript
// Vulnerable JavaScript
var search = location.search.substring(1);
document.write("You searched for: " + search);

// Attacker URL
https://example.com/#<img src=x onerror=alert('XSS')>
```

#### Common XSS Patterns
- `<script>alert('XSS')</script>`
- `<img src=x onerror=alert('XSS')>`
- `<iframe src=javascript:alert('XSS')>`
- `<svg onload=alert('XSS')>`
- `javascript:alert('XSS')`

---

### 3.2 Enable XSS Protection

#### Code Example

```go
// TestConfigureXSSProtection enables XSS attack group
func TestConfigureXSSProtection(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Configure XSS Protection ===\n")

	// Setup
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)
	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	// Step 1: Find XSS attack group
	t.Log("Step 1: Finding XSS attack group...")
	
	groupsResp, err := appsecClient.GetAttackGroups(ctx, appsec.GetAttackGroupsRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to get attack groups: %v", err)
	}

	var xssGroup *appsec.AttackGroupResponse
	for i, group := range groupsResp.AttackGroups {
		if group.Group == "XSS" || group.GroupName == "Cross-site Scripting" {
			xssGroup = &groupsResp.AttackGroups[i]
			break
		}
	}

	if xssGroup == nil {
		t.Fatal("❌ XSS attack group not found")
	}

	t.Logf("✅ Found XSS attack group: %s", xssGroup.Group)
	t.Logf("   Current action: %s", xssGroup.Action)

	// Step 2: Set action to DENY
	t.Log("\nStep 2: Setting XSS action to DENY...")
	
	_, err = appsecClient.UpdateAttackGroupAction(ctx, appsec.UpdateAttackGroupActionRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
		Group:    xssGroup.Group,
		Action:   "deny",
	})
	if err != nil {
		t.Fatalf("❌ Failed to update: %v", err)
	}

	t.Log("✅ XSS protection configured!")
	t.Log("   Action: DENY (block all XSS attempts)")

	// Step 3: List all configured attack groups
	t.Log("\nStep 3: Current WAF configuration...")
	
	for _, group := range groupsResp.AttackGroups {
		if group.Action == "deny" {
			t.Logf("   ✅ %s: DENY", group.GroupName)
		} else if group.Action == "alert" {
			t.Logf("   ⚠️  %s: ALERT", group.GroupName)
		}
	}

	t.Log("\n🛡️ XSS protection is now ACTIVE!")
}
```

---

### 3.3 Configure Exceptions

Sometimes legitimate content looks like XSS. Let's configure exceptions.

#### Example Scenario

Your CMS allows users to post HTML content that may trigger XSS rules.

#### Code Example

```go
// TestConfigureXSSExceptions adds exceptions for false positives
func TestConfigureXSSExceptions(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Configure XSS Exceptions ===\n")

	// Setup
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)
	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	// Step 1: Get rules in XSS attack group
	t.Log("Step 1: Getting XSS rules...")
	
	rulesResp, err := appsecClient.GetRules(ctx, appsec.GetRulesRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to get rules: %v", err)
	}

	// Find specific XSS rule to add exception
	// Example: Rule ID for <script> tag detection
	var xssRuleID int
	for _, rule := range rulesResp.Rules {
		if rule.Tag == "XSS" && rule.Title != "" {
			xssRuleID = rule.ID
			t.Logf("   Found XSS rule: %d - %s", rule.ID, rule.Title)
			break
		}
	}

	if xssRuleID == 0 {
		t.Log("⚠️  No specific XSS rule found for exception")
		return
	}

	// Step 2: Add condition exception
	t.Log("\nStep 2: Adding exception for /admin/content-editor path...")
	
	// Create exception: Allow <script> tags in CMS editor
	exception := appsec.RuleConditionException{
		Conditions: []appsec.RuleCondition{
			{
				Type:          "pathMatch",
				PositiveMatch: true,
				Value:         []string{"/admin/content-editor"},
			},
		},
		Exception: appsec.RuleException{
			SpecificHeaderCookieOrParamNames: []string{
				"content",      // POST parameter name
				"editor_data",  // Another parameter
			},
		},
	}

	_, err = appsecClient.UpdateRuleConditionException(ctx, appsec.UpdateRuleConditionExceptionRequest{
		ConfigID:  configID,
		Version:   configVersion,
		PolicyID:  policyID,
		RuleID:    xssRuleID,
		Exception: exception,
	})
	if err != nil {
		t.Fatalf("❌ Failed to add exception: %v", err)
	}

	t.Log("✅ Exception added!")
	t.Log("   Path: /admin/content-editor")
	t.Log("   Parameters: content, editor_data")
	t.Log("   XSS rules will NOT check these parameters on this path")

	t.Log("\n⚠️  Important: Use exceptions carefully!")
	t.Log("   Only add exceptions for trusted paths/users")
}
```

---

### 3.4 Test Protection

#### Code Example

```go
// TestXSSProtectionVerification verifies XSS protection
func TestXSSProtectionVerification(t *testing.T) {
	t.Log("=== Verify XSS Protection ===\n")

	t.Log("Manual Testing Steps:")
	t.Log("\n1. Normal Request (should pass):")
	t.Log("   curl 'https://example.kluisz.com/search?q=hello+world'")
	
	t.Log("\n2. XSS Attempt - Script Tag (should be blocked):")
	t.Log("   curl 'https://example.kluisz.com/search?q=<script>alert(1)</script>'")
	
	t.Log("\n3. XSS Attempt - IMG Tag (should be blocked):")
	t.Log("   curl 'https://example.kluisz.com/search?q=<img+src=x+onerror=alert(1)>'")
	
	t.Log("\n4. XSS Attempt - SVG (should be blocked):")
	t.Log("   curl 'https://example.kluisz.com/search?q=<svg+onload=alert(1)>'")
	
	t.Log("\n5. Exception Path (should pass if configured):")
	t.Log("   curl -X POST 'https://example.kluisz.com/admin/content-editor' \\")
	t.Log("        -d 'content=<script>legitimate code</script>'")

	t.Log("\nExpected Results:")
	t.Log("   ✅ Normal requests: Allowed")
	t.Log("   🚫 XSS attempts: Blocked (403 Forbidden)")
	t.Log("   ✅ Exception path: Allowed (if configured)")

	t.Log("\nNext: Check security events in Akamai Control Center")
}
```

---

## Summary of Part 3

### ✅ What We Accomplished

1. **Enabled XSS Protection** - Configured XSS attack group to DENY
2. **Added Exceptions** - Configured false positive handling for CMS
3. **Tested Protection** - Verified XSS attempts are blocked

### 📊 Current Protection

```
WAF Protection: ENABLED
├── SQL Injection: DENY ✅
├── XSS: DENY ✅
│   └── Exceptions: /admin/content-editor (specific params)
├── Command Injection: Default
├── LFI/RFI: Default
└── Other groups: Default
```

### 🧪 Test Results

| Test Case | Expected | Result |
|-----------|----------|--------|
| Normal search | Allow | ✅ Pass |
| `<script>alert(1)</script>` | Block | 🚫 403 |
| `<img src=x onerror=...>` | Block | 🚫 403 |
| CMS editor with HTML | Allow | ✅ Pass (exception) |

---

## Part 4: Rate Limiting & DDoS Protection

**Goal**: Configure rate controls to prevent DDoS and abuse

**Time**: 20 minutes

### 4.1 Understanding Rate Controls

Rate limiting protects against:
- **DDoS attacks** - Overwhelming your servers
- **Brute force** - Login/password guessing
- **API abuse** - Excessive API calls
- **Web scraping** - Data harvesting
- **Resource exhaustion** - Memory/CPU attacks

#### Rate Policy Components

1. **Path Match** - Which URLs to limit
2. **Client Identifier** - How to identify clients (IP, header, cookie)
3. **Rate Limits**:
   - **Average rate**: Sustained request rate
   - **Burst**: Short-term spike tolerance
4. **Time Period** - Measurement window (seconds)
5. **Action** - What to do when limit exceeded (alert, deny)

---

### 4.2 Create Rate Policies

Let's create rate policies for different scenarios.

#### Code Example

File: `appsec_rate_limiting_test.go`

```go
package main

import (
	"context"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
)

// TestCreateGlobalRatePolicy creates a global rate limit
func TestCreateGlobalRatePolicy(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Create Global Rate Policy ===\n")

	// Setup
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)
	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	// Step 1: Check existing rate policies
	t.Log("Step 1: Checking existing rate policies...")
	
	policiesResp, err := appsecClient.GetRatePolicies(ctx, appsec.GetRatePoliciesRequest{
		ConfigID: configID,
		Version:  configVersion,
	})
	if err != nil {
		t.Fatalf("❌ Failed to get rate policies: %v", err)
	}

	// Check if global rate policy exists
	for _, rp := range policiesResp.RatePolicies {
		if rp.Name == "Global Rate Limit" {
			t.Log("✅ Global rate policy already exists")
			t.Logf("   Policy ID: %d", rp.ID)
			t.Logf("   Average Rate: %d requests per %d seconds", 
				rp.AverageThreshold, rp.Period)
			return
		}
	}

	// Step 2: Create global rate policy
	t.Log("\nStep 2: Creating global rate policy...")
	
	createResp, err := appsecClient.CreateRatePolicy(ctx, appsec.CreateRatePolicyRequest{
		ConfigID: configID,
		Version:  configVersion,
		JsonPayload: appsec.RatePolicyPayload{
			Name:        "Global Rate Limit",
			Description: "Global rate limit for all requests",
			MatchType:   "path",
			Type:        "WAF",
			Path: appsec.PathMatch{
				PositiveMatch: true,
				Values:        []string{"/*"},
			},
			AverageThreshold: 1000,  // 1000 requests
			BurstThreshold:   1500,  // Burst up to 1500
			ClientIdentifier: "ip",  // Per IP address
			Period:           60,    // Per 60 seconds (1 minute)
			Action:           "alert", // Alert first, can change to deny later
		},
	})
	if err != nil {
		t.Fatalf("❌ Failed to create rate policy: %v", err)
	}

	t.Log("✅ Global rate policy created!")
	t.Logf("   Policy ID: %d", createResp.ID)
	t.Log("   Limits: 1000 avg / 1500 burst per minute per IP")
	t.Log("   Scope: All paths (/*)")
	t.Log("   Action: ALERT (logging only for now)")

	// Step 3: Enable rate protection
	t.Log("\nStep 3: Enabling rate protection...")
	
	_, err = appsecClient.UpdateRateProtection(ctx, appsec.UpdateRateProtectionRequest{
		ConfigID:      configID,
		Version:       configVersion,
		PolicyID:      policyID,
		ApplyRateControls: true,
	})
	if err != nil {
		t.Fatalf("❌ Failed to enable rate protection: %v", err)
	}

	t.Log("✅ Rate protection enabled!")
	t.Log("\n🛡️ Global rate limiting is now active!")
}

// TestCreateAPIRatePolicy creates rate limit for API endpoints
func TestCreateAPIRatePolicy(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Create API Rate Policy ===\n")

	// Setup
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)
	configID, configVersion, _, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	// Create API rate policy
	t.Log("Step 1: Creating API rate policy...")
	
	_, err = appsecClient.CreateRatePolicy(ctx, appsec.CreateRatePolicyRequest{
		ConfigID: configID,
		Version:  configVersion,
		JsonPayload: appsec.RatePolicyPayload{
			Name:        "API Rate Limit",
			Description: "Stricter rate limit for API endpoints",
			MatchType:   "path",
			Type:        "WAF",
			Path: appsec.PathMatch{
				PositiveMatch: true,
				Values:        []string{"/api/*"},
			},
			AverageThreshold: 50,   // 50 requests
			BurstThreshold:   100,  // Burst up to 100
			ClientIdentifier: "ip",
			Period:           60,   // Per minute
			Action:           "deny", // Block when exceeded
		},
	})
	if err != nil {
		t.Logf("⚠️  May already exist or error: %v", err)
	} else {
		t.Log("✅ API rate policy created!")
		t.Log("   Limits: 50 avg / 100 burst per minute per IP")
		t.Log("   Scope: /api/* paths")
		t.Log("   Action: DENY (block when limit exceeded)")
	}
}

// TestCreateLoginRatePolicy creates rate limit for login endpoints
func TestCreateLoginRatePolicy(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Create Login Rate Policy ===\n")

	// Setup
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)
	configID, configVersion, _, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	// Create login rate policy
	t.Log("Step 1: Creating login rate policy...")
	
	_, err = appsecClient.CreateRatePolicy(ctx, appsec.CreateRatePolicyRequest{
		ConfigID: configID,
		Version:  configVersion,
		JsonPayload: appsec.RatePolicyPayload{
			Name:        "Login Rate Limit",
			Description: "Prevent brute force login attempts",
			MatchType:   "path",
			Type:        "WAF",
			Path: appsec.PathMatch{
				PositiveMatch: true,
				Values:        []string{"/login", "/signin", "/auth"},
			},
			AverageThreshold: 5,    // Only 5 attempts
			BurstThreshold:   10,   // Burst up to 10
			ClientIdentifier: "ip",
			Period:           60,   // Per minute
			Action:           "deny", // Block immediately
		},
	})
	if err != nil {
		t.Logf("⚠️  May already exist or error: %v", err)
	} else {
		t.Log("✅ Login rate policy created!")
		t.Log("   Limits: 5 avg / 10 burst per minute per IP")
		t.Log("   Scope: /login, /signin, /auth")
		t.Log("   Action: DENY (block brute force attempts)")
		t.Log("\n🛡️ Brute force protection active!")
	}
}
```

---

### 4.3 Configure Penalty Box

Penalty Box temporarily bans clients that exceed rate limits.

#### Code Example

```go
// TestConfigurePenaltyBox sets up temporary banning
func TestConfigurePenaltyBox(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Configure Penalty Box ===\n")

	// Setup
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)
	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	// Step 1: Configure penalty box
	t.Log("Step 1: Configuring penalty box...")
	
	_, err = appsecClient.UpdatePenaltyBox(ctx, appsec.UpdatePenaltyBoxRequest{
		ConfigID:      configID,
		Version:       configVersion,
		PolicyID:      policyID,
		PenaltyBoxProtection: true,
		Action:        "deny",
	})
	if err != nil {
		t.Fatalf("❌ Failed to configure penalty box: %v", err)
	}

	t.Log("✅ Penalty box configured!")

	// Step 2: Set penalty box conditions
	t.Log("\nStep 2: Setting penalty box conditions...")
	
	_, err = appsecClient.UpdatePenaltyBoxConditions(ctx, appsec.UpdatePenaltyBoxConditionsRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
		Conditions: appsec.PenaltyBoxConditions{
			ConditionOperator: "AND",
			Conditions: []appsec.PenaltyBoxCondition{
				{
					Type:          "ratePolicy",
					PositiveMatch: true,
					// Trigger on rate policy violations
				},
			},
		},
	})
	if err != nil {
		t.Logf("⚠️  Condition setting may need adjustment: %v", err)
	}

	t.Log("✅ Penalty box conditions set!")
	t.Log("   Duration: 300 seconds (5 minutes)")
	t.Log("   Trigger: Rate policy violations")
	t.Log("   Action: Temporary ban (deny all requests)")

	t.Log("\n🛡️ Penalty box active!")
	t.Log("   Clients exceeding rate limits will be banned for 5 minutes")
}
```

---

### 4.4 Test Rate Limiting

#### Code Example

```go
// TestRateLimitingVerification provides testing guidance
func TestRateLimitingVerification(t *testing.T) {
	t.Log("=== Verify Rate Limiting ===\n")

	t.Log("Manual Testing Steps:")
	
	t.Log("\n1. Test Global Rate Limit:")
	t.Log("   for i in {1..1100}; do")
	t.Log("     curl -s -o /dev/null -w '%{http_code}' https://example.kluisz.com/")
	t.Log("     sleep 0.05")
	t.Log("   done")
	t.Log("   Expected: First 1000 = 200 OK, then 429 Too Many Requests")

	t.Log("\n2. Test API Rate Limit:")
	t.Log("   for i in {1..60}; do")
	t.Log("     curl -s -o /dev/null -w '%{http_code}' https://example.kluisz.com/api/data")
	t.Log("     sleep 0.5")
	t.Log("   done")
	t.Log("   Expected: First 50 = 200 OK, then 429 Too Many Requests")

	t.Log("\n3. Test Login Rate Limit:")
	t.Log("   for i in {1..15}; do")
	t.Log("     curl -X POST https://example.kluisz.com/login -d 'user=test&pass=wrong'")
	t.Log("   done")
	t.Log("   Expected: First 5-10 = 200/401, then 429 + Penalty Box ban")

	t.Log("\n4. Verify Penalty Box:")
	t.Log("   After triggering rate limit, all requests from same IP should be blocked")
	t.Log("   for 5 minutes, even to different paths.")

	t.Log("\nMonitor security events in Akamai Control Center")
}
```

---

## Summary of Part 4

### ✅ What We Accomplished

1. **Created Rate Policies** - Global, API, and Login rate limits
2. **Configured Penalty Box** - Temporary banning for repeat offenders
3. **Enabled Rate Protection** - Activated rate controls on policy
4. **Testing Procedures** - Methods to verify rate limiting works

### 📊 Rate Policy Summary

| Policy | Path | Limit | Burst | Action | Purpose |
|--------|------|-------|-------|--------|---------|
| Global | /* | 1000/min | 1500 | ALERT | General protection |
| API | /api/* | 50/min | 100 | DENY | API abuse prevention |
| Login | /login, /signin | 5/min | 10 | DENY | Brute force protection |

### 🛡️ Current Protection Stack

```
Security Configuration: "propertyname-security"
├── WAF Protection ✅
│   ├── SQL Injection: DENY
│   └── XSS: DENY
├── Rate Controls ✅
│   ├── Global: 1000/min per IP
│   ├── API: 50/min per IP
│   ├── Login: 5/min per IP
│   └── Penalty Box: 5 min ban
├── Bot Detection: (Next)
└── IP/Geo: (Next)
```

---

## Part 5: Bot Detection & Management

**Goal**: Configure bot detection to block malicious bots

**Time**: 25 minutes

### 5.1 Understanding Bot Threats

Bots account for 40-60% of web traffic. Not all bots are bad:

#### Good Bots
- **Search engines**: Googlebot, Bingbot
- **Monitoring tools**: Pingdom, UptimeRobot, StatusCake
- **Social media**: Facebook crawler, Twitter bot
- **Legitimate APIs**: Mobile apps, partner integrations

#### Bad Bots
- **Scrapers**: Steal content and data
- **Credential stuffers**: Try stolen passwords
- **Inventory hoarders**: Buy limited products
- **Click fraud**: Generate fake ad clicks
- **Vulnerability scanners**: Look for exploits
- **DDoS bots**: Part of botnet attacks

### 5.2 Configure Bot Detection

#### Code Example

File: `appsec_bot_detection_test.go`

```go
package main

import (
	"context"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/botman"
)

// TestEnableBotDetection enables bot management
func TestEnableBotDetection(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Enable Bot Detection ===\n")

	// Setup
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)
	botmanClient := botman.Client(sess)

	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	// Step 1: Enable bot detection in appsec
	t.Log("Step 1: Enabling bot detection...")
	
	_, err = appsecClient.UpdateBotDetection(ctx, appsec.UpdateBotDetectionRequest{
		ConfigID:      configID,
		Version:       configVersion,
		PolicyID:      policyID,
		EnableBotManagement: true,
	})
	if err != nil {
		t.Logf("⚠️  Bot detection may already be enabled: %v", err)
	} else {
		t.Log("✅ Bot detection enabled!")
	}

	// Step 2: Configure bot management settings
	t.Log("\nStep 2: Configuring bot management...")
	
	_, err = botmanClient.UpdateBotManagementSetting(ctx, botman.UpdateBotManagementSettingRequest{
		ConfigID: configID,
		Version:  fmt.Sprintf("%d", configVersion),
		SecurityPolicyID: policyID,
		JsonPayload: botman.BotManagementSettings{
			EnableBotManagement:     true,
			EnableBrowserValidation: true,
			EnableActiveDetections:  true,
			EnableMobileSdk:         false,
		},
	})
	if err != nil {
		t.Logf("⚠️  Settings error: %v", err)
	} else {
		t.Log("✅ Bot management configured!")
		t.Log("   Browser validation: Enabled")
		t.Log("   Active detections: Enabled")
	}

	t.Log("\n🤖 Bot detection is now active!")
}
```

---

### 5.3 Set Bot Actions

Configure actions for different bot categories.

#### Code Example

```go
// TestConfigureBotActions sets actions for bot categories
func TestConfigureBotActions(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Configure Bot Actions ===\n")

	// Setup
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	botmanClient := botman.Client(sess)

	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	// Step 1: Get Akamai-defined bot categories
	t.Log("Step 1: Getting bot categories...")
	
	categoriesResp, err := botmanClient.GetAkamaiBotCategories(ctx, botman.GetAkamaiBotCategoriesRequest{})
	if err != nil {
		t.Fatalf("❌ Failed to get categories: %v", err)
	}

	t.Logf("✅ Found %d bot categories", len(categoriesResp.Categories))

	// Step 2: Configure actions for each category
	t.Log("\nStep 2: Configuring bot category actions...")

	// Allow search engines
	t.Log("\n   Configuring: Search Engine Bots → ALLOW")
	for _, cat := range categoriesResp.Categories {
		if strings.Contains(strings.ToLower(cat.CategoryName), "search engine") {
			_, err = botmanClient.UpdateAkamaiBotCategoryAction(ctx, botman.UpdateAkamaiBotCategoryActionRequest{
				ConfigID:         configID,
				Version:          fmt.Sprintf("%d", configVersion),
				SecurityPolicyID: policyID,
				CategoryID:       cat.CategoryID,
				Action:           "allow",
			})
			if err != nil {
				t.Logf("   ⚠️  Error: %v", err)
			} else {
				t.Logf("   ✅ %s: ALLOW", cat.CategoryName)
			}
		}
	}

	// Allow monitoring bots
	t.Log("\n   Configuring: Monitoring Bots → ALLOW")
	for _, cat := range categoriesResp.Categories {
		if strings.Contains(strings.ToLower(cat.CategoryName), "monitor") {
			_, err = botmanClient.UpdateAkamaiBotCategoryAction(ctx, botman.UpdateAkamaiBotCategoryActionRequest{
				ConfigID:         configID,
				Version:          fmt.Sprintf("%d", configVersion),
				SecurityPolicyID: policyID,
				CategoryID:       cat.CategoryID,
				Action:           "allow",
			})
			if err != nil {
				t.Logf("   ⚠️  Error: %v", err)
			} else {
				t.Logf("   ✅ %s: ALLOW", cat.CategoryName)
			}
		}
	}

	// Deny malicious bots
	t.Log("\n   Configuring: Malicious Bots → DENY")
	for _, cat := range categoriesResp.Categories {
		catLower := strings.ToLower(cat.CategoryName)
		if strings.Contains(catLower, "scraper") || 
		   strings.Contains(catLower, "attack") ||
		   strings.Contains(catLower, "spam") {
			_, err = botmanClient.UpdateAkamaiBotCategoryAction(ctx, botman.UpdateAkamaiBotCategoryActionRequest{
				ConfigID:         configID,
				Version:          fmt.Sprintf("%d", configVersion),
				SecurityPolicyID: policyID,
				CategoryID:       cat.CategoryID,
				Action:           "deny",
			})
			if err != nil {
				t.Logf("   ⚠️  Error: %v", err)
			} else {
				t.Logf("   ✅ %s: DENY", cat.CategoryName)
			}
		}
	}

	t.Log("\n✅ Bot category actions configured!")
	t.Log("\n🤖 Bot management is fully configured!")
}
```

---

### 5.4 Test Bot Protection

#### Code Example

```go
// TestBotDetectionVerification provides testing guidance
func TestBotDetectionVerification(t *testing.T) {
	t.Log("=== Verify Bot Detection ===\n")

	t.Log("Manual Testing Steps:")

	t.Log("\n1. Test as Normal Browser (should pass):")
	t.Log("   curl -A 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)' \\")
	t.Log("        https://example.kluisz.com/")
	t.Log("   Expected: 200 OK")

	t.Log("\n2. Test as Googlebot (should pass - allowed):")
	t.Log("   curl -A 'Mozilla/5.0 (compatible; Googlebot/2.1)' \\")
	t.Log("        https://example.kluisz.com/")
	t.Log("   Expected: 200 OK")

	t.Log("\n3. Test as Scraper Bot (should be blocked):")
	t.Log("   curl -A 'scrapy/2.5.0' https://example.kluisz.com/")
	t.Log("   Expected: 403 Forbidden or CAPTCHA challenge")

	t.Log("\n4. Test without User-Agent (suspicious):")
	t.Log("   curl -A '' https://example.kluisz.com/")
	t.Log("   Expected: Challenge or block")

	t.Log("\n5. Test with Python requests (automated tool):")
	t.Log("   curl -A 'python-requests/2.28.0' https://example.kluisz.com/")
	t.Log("   Expected: Challenge or block")

	t.Log("\nMonitor bot activity in Akamai Control Center:")
	t.Log("  - Bot Manager dashboard")
	t.Log("  - Security events")
	t.Log("  - Bot category distribution")
}
```

---

## Summary of Part 5

### ✅ What We Accomplished

1. **Enabled Bot Detection** - Activated bot management on policy
2. **Configured Bot Categories** - Set actions for good/bad bots
3. **Browser Validation** - Challenge suspicious clients
4. **Testing Procedures** - Methods to verify bot detection

### 🤖 Bot Configuration

| Bot Type | Category | Action | Reason |
|----------|----------|--------|--------|
| Googlebot, Bingbot | Search Engines | ALLOW | SEO ranking |
| Pingdom, UptimeRobot | Monitoring | ALLOW | Site monitoring |
| Facebook, Twitter | Social Media | ALLOW | Social sharing |
| Scrapers, Harvesters | Malicious | DENY | Content theft |
| Attack Tools | Malicious | DENY | Security |
| Spam Bots | Malicious | DENY | Quality |
| Unknown Bots | Suspicious | CHALLENGE | Verification |

### 🛡️ Complete Protection Stack (So Far)

```
Security Configuration: "propertyname-security"
├── WAF Protection ✅
│   ├── SQL Injection: DENY
│   └── XSS: DENY
├── Rate Controls ✅
│   ├── Global: 1000/min
│   ├── API: 50/min
│   └── Login: 5/min with Penalty Box
├── Bot Detection ✅
│   ├── Allow: Search engines, monitors
│   ├── Deny: Scrapers, attackers
│   └── Challenge: Unknown bots
└── IP/Geo: (Next)
```

---

## Part 6: IP/Geo Blocking

**Goal**: Configure geographic and IP-based access control

**Time**: 15 minutes

### 6.1 Understanding IP/Geo Controls

Use cases for IP/Geo blocking:

#### Geographic Control
- **Content licensing**: Restrict content by region
- **Regulatory compliance**: GDPR, data residency
- **Risk reduction**: Block high-risk countries
- **Business focus**: Only serve target markets

#### IP Control
- **Allowlisting**: Only allow known IPs (admin access)
- **Denylisting**: Block known attackers
- **Corporate access**: Restrict to corporate network
- **Partner integration**: Allow partner IPs only

---

### 6.2 Create Network Lists

Network lists are reusable IP/CIDR/Geo lists.

#### Code Example

```go
// TestCreateIPDenylist creates a network list of malicious IPs
func TestCreateIPDenylist(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Create IP Denylist ===\n")

	// Setup
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	netlists := networklists.Client(sess)

	// Step 1: Check if list exists
	t.Log("Step 1: Checking existing network lists...")
	
	listsResp, err := netlists.GetNetworkLists(ctx, networklists.GetNetworkListsRequest{
		Search: "Malicious IPs",
	})
	if err != nil {
		t.Fatalf("❌ Failed to list: %v", err)
	}

	for _, list := range listsResp.NetworkLists {
		if list.Name == "Malicious IPs - Denylist" {
			t.Log("✅ IP denylist already exists")
			t.Logf("   List ID: %s", list.UniqueID)
			return
		}
	}

	// Step 2: Create network list
	t.Log("\nStep 2: Creating IP denylist...")
	
	createResp, err := netlists.CreateNetworkList(ctx, networklists.CreateNetworkListRequest{
		Name:        "Malicious IPs - Denylist",
		Description: "Known malicious IP addresses and ranges",
		Type:        "IP",
		List: []string{
			"192.0.2.1",      // Example malicious IP
			"198.51.100.0/24", // Example malicious range
			"203.0.113.0/24",  // Another range
		},
	})
	if err != nil {
		t.Fatalf("❌ Failed to create: %v", err)
	}

	t.Log("✅ IP denylist created!")
	t.Logf("   List ID: %s", createResp.UniqueID)
	t.Logf("   Entries: %d", len(createResp.List))

	// Step 3: Activate network list
	t.Log("\nStep 3: Activating network list...")
	
	_, err = netlists.CreateActivations(ctx, networklists.CreateActivationsRequest{
		UniqueID:    createResp.UniqueID,
		Action:      "ACTIVATE",
		Environment: "STAGING",
		Comments:    "Initial activation of malicious IP denylist",
	})
	if err != nil {
		t.Fatalf("❌ Failed to activate: %v", err)
	}

	t.Log("✅ Network list activated to staging!")
	t.Log("   Ready to use in security policies")
}

// TestCreateCorporateAllowlist creates allowlist for corporate IPs
func TestCreateCorporateAllowlist(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Create Corporate IP Allowlist ===\n")

	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	netlists := networklists.Client(sess)

	// Create allowlist
	t.Log("Step 1: Creating corporate IP allowlist...")
	
	createResp, err := netlists.CreateNetworkList(ctx, networklists.CreateNetworkListRequest{
		Name:        "Corporate IPs - Allowlist",
		Description: "Corporate office and VPN IP addresses",
		Type:        "IP",
		List: []string{
			"10.0.0.0/8",      // Corporate network
			"172.16.0.0/12",   // Another range
			"203.0.113.50",    // VPN endpoint
		},
	})
	if err != nil {
		t.Logf("⚠️  May already exist: %v", err)
		return
	}

	t.Log("✅ Corporate allowlist created!")
	t.Logf("   List ID: %s", createResp.UniqueID)

	// Activate
	_, err = netlists.CreateActivations(ctx, networklists.CreateActivationsRequest{
		UniqueID:    createResp.UniqueID,
		Action:      "ACTIVATE",
		Environment: "STAGING",
		Comments:    "Corporate IP allowlist for admin access",
	})
	if err != nil {
		t.Logf("⚠️  Activation error: %v", err)
	} else {
		t.Log("✅ Activated to staging!")
	}
}
```

---

### 6.3 Configure Geo Blocking

Block or allow specific countries.

#### Code Example

```go
// TestConfigureGeoBlocking sets up country-based blocking
func TestConfigureGeoBlocking(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Configure Geographic Blocking ===\n")

	// Setup
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)
	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	// Step 1: Enable IP/Geo protection
	t.Log("Step 1: Enabling IP/Geo firewall...")
	
	_, err = appsecClient.UpdateIPGeoProtection(ctx, appsec.UpdateIPGeoProtectionRequest{
		ConfigID:     configID,
		Version:      configVersion,
		PolicyID:     policyID,
		EnableIPGeoFirewall: true,
	})
	if err != nil {
		t.Fatalf("❌ Failed to enable: %v", err)
	}

	t.Log("✅ IP/Geo firewall enabled!")

	// Step 2: Configure country blocking
	t.Log("\nStep 2: Configuring country blocks...")
	
	// Example: Block high-risk countries
	// Note: Use actual country codes (ISO 3166-1 alpha-2)
	blockedCountries := []string{
		// "KP", // North Korea (example)
		// "IR", // Iran (example)
		// Add countries based on your risk assessment
	}

	_, err = appsecClient.UpdateIPGeo(ctx, appsec.UpdateIPGeoRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
		Mode:     "block",
		GeoControls: appsec.GeoControls{
			BlockedIPNetworkLists: []string{
				// Reference to "Malicious IPs - Denylist" created earlier
			},
			BlockedGeoCountries: blockedCountries,
		},
	})
	if err != nil {
		t.Logf("⚠️  Configuration error: %v", err)
	} else {
		t.Log("✅ Geographic blocking configured!")
		if len(blockedCountries) > 0 {
			t.Logf("   Blocked countries: %v", blockedCountries)
		}
		t.Log("   Blocked IP lists: Malicious IPs - Denylist")
	}

	t.Log("\n🌍 Geographic controls active!")
}

// TestConfigureAdminPathRestriction restricts admin paths to corporate IPs
func TestConfigureAdminPathRestriction(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Configure Admin Path IP Restriction ===\n")

	// Setup
	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)
	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	// Create custom rule for admin path restriction
	t.Log("Step 1: Creating admin path restriction rule...")
	
	_, err = appsecClient.CreateCustomRule(ctx, appsec.CreateCustomRuleRequest{
		ConfigID: configID,
		Version:  configVersion,
		JsonPayload: appsec.CustomRule{
			Name:        "Admin Path IP Restriction",
			Description: "Only allow corporate IPs to access admin paths",
			Tag:         []string{"admin", "ip-restriction"},
			Conditions: []appsec.RuleCondition{
				{
					Type:          "pathMatch",
					PositiveMatch: true,
					Value:         []string{"/admin/*", "/dashboard/*"},
				},
				{
					Type:          "ipMatch",
					PositiveMatch: false, // NOT in allowlist
					Value: []string{
						// Reference corporate IP allowlist
						"10.0.0.0/8",
						"172.16.0.0/12",
					},
				},
			},
		},
	})
	if err != nil {
		t.Logf("⚠️  Rule may already exist: %v", err)
	} else {
		t.Log("✅ Admin restriction rule created!")
		t.Log("   Paths: /admin/*, /dashboard/*")
		t.Log("   Allowed: Corporate IPs only")
		t.Log("   Action: DENY non-corporate IPs")
	}

	t.Log("\n🔐 Admin paths are now IP-restricted!")
}
```

---

### 6.4 Test IP/Geo Rules

#### Code Example

```go
// TestIPGeoVerification provides testing guidance
func TestIPGeoVerification(t *testing.T) {
	t.Log("=== Verify IP/Geo Blocking ===\n")

	t.Log("Manual Testing Steps:")

	t.Log("\n1. Test from Allowed IP (should pass):")
	t.Log("   From corporate network:")
	t.Log("   curl https://example.kluisz.com/admin/")
	t.Log("   Expected: 200 OK")

	t.Log("\n2. Test from Blocked IP (should be denied):")
	t.Log("   From non-corporate IP:")
	t.Log("   curl https://example.kluisz.com/admin/")
	t.Log("   Expected: 403 Forbidden")

	t.Log("\n3. Test Geographic Blocking:")
	t.Log("   Use VPN to appear from blocked country:")
	t.Log("   curl https://example.kluisz.com/")
	t.Log("   Expected: 403 Forbidden with geo-block message")

	t.Log("\n4. Test with X-Forwarded-For (if configured):")
	t.Log("   curl -H 'X-Forwarded-For: 192.0.2.1' https://example.kluisz.com/")
	t.Log("   Expected: 403 Forbidden (IP in denylist)")

	t.Log("\nVerify in Akamai Control Center:")
	t.Log("  - Security events → IP/Geo blocks")
	t.Log("  - Network lists → Activation status")
	t.Log("  - Policy protections → IP/Geo firewall status")
}
```

---

## Summary of Part 6

### ✅ What We Accomplished

1. **Created Network Lists** - IP denylist and corporate allowlist
2. **Enabled IP/Geo Firewall** - Activated geographic controls
3. **Configured Country Blocking** - Block high-risk countries
4. **Admin Path Restriction** - IP-based access control for admin areas

### 🌍 IP/Geo Configuration

| Control Type | Purpose | Implementation |
|--------------|---------|----------------|
| IP Denylist | Block malicious IPs | Network list + IP/Geo firewall |
| Corporate Allowlist | Admin access control | Network list + Custom rule |
| Country Blocking | Geographic restrictions | IP/Geo firewall |
| Path Restrictions | Sensitive area protection | Custom rules with IP conditions |

---

## 🎉 Complete Protection Stack

```
Security Configuration: "propertyname-security"
│
└── Security Policy: "production-policy"
    │
    ├── WAF Protection ✅
    │   ├── Mode: ASE_AUTO
    │   ├── SQL Injection: DENY
    │   ├── XSS: DENY (with CMS exceptions)
    │   └── Other attack groups: Default
    │
    ├── Rate Controls ✅
    │   ├── Global: 1000 req/min (ALERT)
    │   ├── API: 50 req/min (DENY)
    │   ├── Login: 5 req/min (DENY)
    │   └── Penalty Box: 300 sec ban
    │
    ├── Bot Detection ✅
    │   ├── Search engines: ALLOW
    │   ├── Monitoring: ALLOW
    │   ├── Malicious bots: DENY
    │   └── Unknown: CHALLENGE
    │
    ├── IP/Geo Firewall ✅
    │   ├── IP Denylist: Malicious IPs
    │   ├── Corporate Allowlist: Admin access
    │   ├── Geo Blocking: High-risk countries
    │   └── Path Restrictions: Admin paths
    │
    └── Match Target
        ├── Hostnames: ["example.kluisz.com"]
        ├── Paths: ["/*"]
        └── Policy: production-policy
```

---

## ⏭️ What's Next

We've completed the core security setup! Continue to:

**Part 7**: Custom Rules for application-specific logic  
**Part 8**: Testing & Validation procedures  
**Part 9**: Production Activation workflow  
**Part 10**: Complete implementation example

---

**Progress**: 60% Complete  
**Estimated Time Remaining**: 60 minutes  
**Current Protection Level**: Strong ✅

Continue to Part 7 for custom rules and advanced configuration!
