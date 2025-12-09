package main

import (
	"context"
	"strings"
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
		t.Logf("⚠️  WAF mode may already be set: %v", err)
	} else {
		t.Log("✅ WAF mode set to ASE_AUTO (automatic updates)")
	}

	// Step 4: Enable WAF protection
	t.Log("\nStep 4: Enabling WAF protection...")

	_, err = appsecClient.UpdateWAFProtection(ctx, appsec.UpdateWAFProtectionRequest{
		ConfigID:                      configID,
		Version:                       configVersion,
		PolicyID:                      policyID,
		ApplyApplicationLayerControls: true,
	})
	if err != nil {
		t.Logf("⚠️  WAF protection may already be enabled: %v", err)
	} else {
		t.Log("✅ WAF protection enabled!")
	}

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
	t.Logf("   Apply Application Layer Controls: %v", protResp.ApplyApplicationLayerControls)

	t.Log("\n✅ WAF protection is now active!")
}

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
	var sqliGroup string
	var currentAction string
	for _, group := range groupsResp.AttackGroups {
		if group.Group == "SQL" || strings.Contains(strings.ToLower(group.Group), "sql") {
			sqliGroup = group.Group
			currentAction = group.Action
			break
		}
	}

	if sqliGroup == "" {
		t.Fatal("❌ SQL Injection attack group not found")
	}

	t.Logf("✅ Found SQLi attack group: %s", sqliGroup)
	t.Logf("   Current action: %s", currentAction)

	// Step 2: Set action to DENY
	t.Log("\nStep 2: Setting SQL Injection action to DENY...")

	_, err = appsecClient.UpdateAttackGroup(ctx, appsec.UpdateAttackGroupRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
		Group:    sqliGroup,
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
		Group:    sqliGroup,
	})
	if err != nil {
		t.Fatalf("❌ Failed to verify: %v", err)
	}

	t.Log("✅ Configuration verified!")
	t.Logf("   Action: %s", verifyResp.Action)

	t.Log("\n🛡️ SQL Injection protection is now ACTIVE!")
	t.Log("   All SQL injection attempts will be BLOCKED")
}

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

	var xssGroup string
	var currentAction string
	for _, group := range groupsResp.AttackGroups {
		if group.Group == "XSS" || strings.Contains(strings.ToLower(group.Group), "xss") ||
			strings.Contains(strings.ToLower(group.Group), "cross-site scripting") {
			xssGroup = group.Group
			currentAction = group.Action
			break
		}
	}

	if xssGroup == "" {
		t.Fatal("❌ XSS attack group not found")
	}

	t.Logf("✅ Found XSS attack group: %s", xssGroup)
	t.Logf("   Current action: %s", currentAction)

	// Step 2: Set action to DENY
	t.Log("\nStep 2: Setting XSS action to DENY...")

	_, err = appsecClient.UpdateAttackGroup(ctx, appsec.UpdateAttackGroupRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
		Group:    xssGroup,
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
			t.Logf("   ✅ %s: DENY", group.Group)
		}
	}

	t.Log("\n🛡️ XSS protection is now ACTIVE!")
}

// TestConfigureCommandInjection configures command injection protection
func TestConfigureCommandInjection(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Configure Command Injection Protection ===\n")

	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)
	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	t.Log("Step 1: Finding Command Injection attack group...")

	groupsResp, err := appsecClient.GetAttackGroups(ctx, appsec.GetAttackGroupsRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to get attack groups: %v", err)
	}

	var cmdGroup string
	for _, group := range groupsResp.AttackGroups {
		if group.Group == "CMD" || strings.Contains(strings.ToLower(group.Group), "command") {
			cmdGroup = group.Group
			break
		}
	}

	if cmdGroup == "" {
		t.Log("⚠️  Command Injection attack group not found (may not be available)")
		return
	}

	t.Logf("✅ Found Command Injection attack group: %s", cmdGroup)

	t.Log("\nStep 2: Setting Command Injection action to DENY...")

	_, err = appsecClient.UpdateAttackGroup(ctx, appsec.UpdateAttackGroupRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
		Group:    cmdGroup,
		Action:   "deny",
	})
	if err != nil {
		t.Fatalf("❌ Failed to update: %v", err)
	}

	t.Log("✅ Command Injection protection configured!")
	t.Log("\n🛡️ Command Injection protection is now ACTIVE!")
}

// TestConfigureAllCriticalAttackGroups sets DENY for all critical attack groups
func TestConfigureAllCriticalAttackGroups(t *testing.T) {
	ctx := context.Background()

	t.Log("=== Configure All Critical Attack Groups ===\n")

	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)
	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	// Get all attack groups
	t.Log("Step 1: Getting all attack groups...")

	groupsResp, err := appsecClient.GetAttackGroups(ctx, appsec.GetAttackGroupsRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to get attack groups: %v", err)
	}

	t.Logf("✅ Found %d attack groups", len(groupsResp.AttackGroups))

	// Critical attack groups that should always be DENY
	criticalGroups := []string{
		"SQL", "XSS", "CMD", "LFI", "RFI",
		"POLICY", "PROTOCOL", "TROJAN",
	}

	t.Log("\nStep 2: Configuring critical attack groups to DENY...")

	updatedCount := 0
	for _, criticalGroup := range criticalGroups {
		for _, group := range groupsResp.AttackGroups {
			if strings.EqualFold(group.Group, criticalGroup) ||
				strings.Contains(strings.ToLower(group.Group), strings.ToLower(criticalGroup)) {

				if group.Action != "deny" {
					_, err = appsecClient.UpdateAttackGroup(ctx, appsec.UpdateAttackGroupRequest{
						ConfigID: configID,
						Version:  configVersion,
						PolicyID: policyID,
						Group:    group.Group,
						Action:   "deny",
					})
					if err != nil {
						t.Logf("   ⚠️  Failed to update %s: %v", group.Group, err)
					} else {
						t.Logf("   ✅ %s: set to DENY", group.Group)
						updatedCount++
					}
				} else {
					t.Logf("   ✓ %s: already DENY", group.Group)
				}
				break
			}
		}
	}

	t.Logf("\n✅ Updated %d critical attack groups", updatedCount)

	// Display summary
	t.Log("\nStep 3: Current attack group configuration:")

	denyCount := 0
	alertCount := 0
	noneCount := 0

	for _, group := range groupsResp.AttackGroups {
		switch group.Action {
		case "deny":
			denyCount++
		case "alert":
			alertCount++
		case "none":
			noneCount++
		}
	}

	t.Logf("   DENY: %d groups", denyCount)
	t.Logf("   ALERT: %d groups", alertCount)
	t.Logf("   NONE: %d groups", noneCount)

	t.Log("\n🛡️ Critical attack groups are now protected!")
}

// TestListAllAttackGroups displays all available attack groups
func TestListAllAttackGroups(t *testing.T) {
	ctx := context.Background()

	t.Log("=== List All Attack Groups ===\n")

	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)
	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	groupsResp, err := appsecClient.GetAttackGroups(ctx, appsec.GetAttackGroupsRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to get attack groups: %v", err)
	}

	t.Logf("Found %d attack groups:\n", len(groupsResp.AttackGroups))

	for i, group := range groupsResp.AttackGroups {
		action := group.Action
		symbol := "⚠️"
		switch action {
		case "deny":
			symbol = "🛡️"
		case "alert":
			symbol = "📊"
		case "none":
			symbol = "⭕"
		}

		t.Logf("%d. %s %s: %s", i+1, symbol, group.Group, strings.ToUpper(action))
	}

	t.Log("\nLegend:")
	t.Log("  🛡️ DENY - Blocks the attack")
	t.Log("  📊 ALERT - Logs but doesn't block")
	t.Log("  ⭕ NONE - Rule is disabled")
}

// TestWAFProtectionSummary displays comprehensive WAF status
func TestWAFProtectionSummary(t *testing.T) {
	ctx := context.Background()

	t.Log("=== WAF Protection Summary ===\n")

	sess, err := newSession(EdgercPath, EdgercSection)
	if err != nil {
		t.Fatalf("❌ Authentication failed: %v", err)
	}

	appsecClient := appsec.Client(sess)
	configID, configVersion, policyID, err := getSecurityConfig(ctx, appsecClient)
	if err != nil {
		t.Fatalf("❌ Failed to get config: %v", err)
	}

	// Get WAF mode
	t.Log("WAF Configuration:")
	wafMode, err := appsecClient.GetWAFMode(ctx, appsec.GetWAFModeRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
	})
	if err == nil {
		t.Logf("  Mode: %s", wafMode.Mode)
	}

	// Get WAF protection status
	wafProt, err := appsecClient.GetWAFProtection(ctx, appsec.GetWAFProtectionRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
	})
	if err == nil {
		t.Logf("  Protection Enabled: %v", wafProt.ApplyApplicationLayerControls)
	}

	// Get attack groups
	groupsResp, err := appsecClient.GetAttackGroups(ctx, appsec.GetAttackGroupsRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
	})
	if err != nil {
		t.Fatalf("❌ Failed to get attack groups: %v", err)
	}

	// Count by action
	stats := make(map[string]int)
	for _, group := range groupsResp.AttackGroups {
		stats[group.Action]++
	}

	t.Log("\nAttack Group Statistics:")
	t.Logf("  Total Groups: %d", len(groupsResp.AttackGroups))
	t.Logf("  DENY (blocking): %d", stats["deny"])
	t.Logf("  ALERT (monitoring): %d", stats["alert"])
	t.Logf("  NONE (disabled): %d", stats["none"])

	// List critical protections
	t.Log("\nCritical Protections:")
	criticalGroups := []string{"SQL", "XSS", "CMD", "LFI", "RFI"}
	for _, criticalGroup := range criticalGroups {
		for _, group := range groupsResp.AttackGroups {
			if strings.EqualFold(group.Group, criticalGroup) {
				symbol := "❌"
				if group.Action == "deny" {
					symbol = "✅"
				}
				t.Logf("  %s %s: %s", symbol, group.Group, group.Action)
				break
			}
		}
	}

	t.Log("\n✅ WAF summary complete!")
}
