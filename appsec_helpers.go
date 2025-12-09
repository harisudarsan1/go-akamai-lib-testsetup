package main

import (
	"context"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/appsec"
)

// AppSec constants
const (
	// Security configuration settings
	SecurityConfigName        = "propertyname-security"
	SecurityConfigDescription = "Security configuration for propertyname property"
	SecurityPolicyName        = "production-policy"
	SecurityPolicyPrefix      = "prop"
)

// getSecurityConfig retrieves security configuration, version, and policy ID
// This is a helper function used across multiple AppSec tests
func getSecurityConfig(ctx context.Context, client appsec.APPSEC) (configID int, configVersion int, policyID string, err error) {
	// Get security configuration by name
	configsResp, err := client.GetConfigurations(ctx, appsec.GetConfigurationsRequest{})
	if err != nil {
		return 0, 0, "", fmt.Errorf("failed to list configurations: %w", err)
	}

	// Find our configuration
	for _, cfg := range configsResp.Configurations {
		if cfg.Name == SecurityConfigName {
			configID = cfg.ID
			configVersion = cfg.LatestVersion
			break
		}
	}

	if configID == 0 {
		return 0, 0, "", fmt.Errorf("security configuration '%s' not found", SecurityConfigName)
	}

	// Get security policies
	policiesResp, err := client.GetSecurityPolicies(ctx, appsec.GetSecurityPoliciesRequest{
		ConfigID: configID,
		Version:  configVersion,
	})
	if err != nil {
		return 0, 0, "", fmt.Errorf("failed to list security policies: %w", err)
	}

	// Find our policy
	for _, policy := range policiesResp.Policies {
		if policy.PolicyName == SecurityPolicyName {
			policyID = policy.PolicyID
			break
		}
	}

	if policyID == "" {
		return 0, 0, "", fmt.Errorf("security policy '%s' not found", SecurityPolicyName)
	}

	return configID, configVersion, policyID, nil
}

// configExists checks if a security configuration exists
func configExists(ctx context.Context, client appsec.APPSEC, configName string) (bool, int, int, error) {
	configsResp, err := client.GetConfigurations(ctx, appsec.GetConfigurationsRequest{})
	if err != nil {
		return false, 0, 0, err
	}

	for _, cfg := range configsResp.Configurations {
		if cfg.Name == configName {
			return true, cfg.ID, cfg.LatestVersion, nil
		}
	}

	return false, 0, 0, nil
}

// policyExists checks if a security policy exists within a configuration
func policyExists(ctx context.Context, client appsec.APPSEC, configID, configVersion int, policyName string) (bool, string, error) {
	policiesResp, err := client.GetSecurityPolicies(ctx, appsec.GetSecurityPoliciesRequest{
		ConfigID: configID,
		Version:  configVersion,
	})
	if err != nil {
		return false, "", err
	}

	for _, policy := range policiesResp.Policies {
		if policy.PolicyName == policyName {
			return true, policy.PolicyID, nil
		}
	}

	return false, "", nil
}

// matchTargetExistsForHostname checks if a match target exists for a specific hostname
func matchTargetExistsForHostname(ctx context.Context, client appsec.APPSEC, configID, configVersion int, hostname string) (bool, int, error) {
	targetsResp, err := client.GetMatchTargets(ctx, appsec.GetMatchTargetsRequest{
		ConfigID:      configID,
		ConfigVersion: configVersion,
	})
	if err != nil {
		return false, 0, err
	}

	// Check website targets for the hostname
	for _, target := range targetsResp.MatchTargets.WebsiteTargets {
		for _, h := range target.Hostnames {
			if h == hostname {
				return true, target.TargetID, nil
			}
		}
	}

	return false, 0, nil
}

// getAttackGroupByName finds an attack group by name or group code
func getAttackGroupByName(ctx context.Context, client appsec.APPSEC, configID, configVersion int, policyID, groupName string) (string, string, error) {
	groupsResp, err := client.GetAttackGroups(ctx, appsec.GetAttackGroupsRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to get attack groups: %w", err)
	}

	// Search for group by code (the Group field is what we need)
	for _, group := range groupsResp.AttackGroups {
		if group.Group == groupName {
			return group.Group, group.Action, nil
		}
	}

	return "", "", fmt.Errorf("attack group '%s' not found", groupName)
}

// ratePolicyExists checks if a rate policy exists by name
func ratePolicyExists(ctx context.Context, client appsec.APPSEC, configID int, policyName string) (bool, int, error) {
	policiesResp, err := client.GetRatePolicies(ctx, appsec.GetRatePoliciesRequest{
		ConfigID: configID,
	})
	if err != nil {
		return false, 0, err
	}

	for _, rp := range policiesResp.RatePolicies {
		if rp.MatchType == "path" && rp.Name == policyName {
			return true, rp.ID, nil
		}
	}

	return false, 0, nil
}

// customRuleExists checks if a custom rule exists by name
func customRuleExists(ctx context.Context, client appsec.APPSEC, configID int, ruleName string) (bool, int, error) {
	rulesResp, err := client.GetCustomRules(ctx, appsec.GetCustomRulesRequest{
		ConfigID: configID,
	})
	if err != nil {
		return false, 0, err
	}

	for _, rule := range rulesResp.CustomRules {
		if rule.Name == ruleName {
			return true, rule.ID, nil
		}
	}

	return false, 0, nil
}

// getConfigurationDetails returns detailed information about a configuration
func getConfigurationDetails(ctx context.Context, client appsec.APPSEC, configID int) (*appsec.GetConfigurationResponse, error) {
	return client.GetConfiguration(ctx, appsec.GetConfigurationRequest{
		ConfigID: configID,
	})
}

// getActivationStatus checks the most recent activation status for a configuration
func getActivationStatus(ctx context.Context, client appsec.APPSEC, configID int) (string, error) {
	historyResp, err := client.GetActivationHistory(ctx, appsec.GetActivationHistoryRequest{
		ConfigID: configID,
	})
	if err != nil {
		return "", err
	}

	if len(historyResp.ActivationHistory) > 0 {
		// Return the most recent activation status
		return historyResp.ActivationHistory[0].Status, nil
	}

	return "NOT_ACTIVATED", nil
}

// printSecuritySummary prints a formatted summary of the security configuration
func printSecuritySummary(ctx context.Context, client appsec.APPSEC, configID, configVersion int, policyID string) error {
	fmt.Println("\n" + "=================================================")
	fmt.Println("         Security Configuration Summary")
	fmt.Println("=================================================")

	// Get configuration details
	config, err := client.GetConfiguration(ctx, appsec.GetConfigurationRequest{
		ConfigID: configID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("\nConfiguration: %s (ID: %d)\n", config.Name, config.ID)
	fmt.Printf("Description: %s\n", config.Description)
	fmt.Printf("Latest Version: %d\n", config.LatestVersion)
	fmt.Printf("Staging Version: %d\n", config.StagingVersion)
	fmt.Printf("Production Version: %d\n", config.ProductionVersion)

	// Get policy details
	policy, err := client.GetSecurityPolicy(ctx, appsec.GetSecurityPolicyRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("\nSecurity Policy: %s (ID: %s)\n", policy.PolicyName, policy.PolicyID)

	// Get WAF protection status
	wafProt, err := client.GetWAFProtection(ctx, appsec.GetWAFProtectionRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
	})
	if err == nil {
		fmt.Printf("\nWAF Protection: %v\n", wafProt.ApplyNetworkLayerControls)
	}

	// Get rate protection status
	rateProt, err := client.GetRateProtection(ctx, appsec.GetRateProtectionRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
	})
	if err == nil {
		fmt.Printf("Rate Protection: %v\n", rateProt.ApplyRateControls)
	}

	// Get match targets
	targets, err := client.GetMatchTargets(ctx, appsec.GetMatchTargetsRequest{
		ConfigID:      configID,
		ConfigVersion: configVersion,
	})
	if err == nil {
		websiteCount := len(targets.MatchTargets.WebsiteTargets)
		apiCount := len(targets.MatchTargets.APITargets)
		fmt.Printf("\nMatch Targets: %d website, %d API\n", websiteCount, apiCount)

		for _, target := range targets.MatchTargets.WebsiteTargets {
			fmt.Printf("  - Website Target %d: %v\n", target.TargetID, target.Hostnames)
		}
	}

	fmt.Println("\n=================================================\n")

	return nil
}

// validateSecuritySetup performs basic validation of security configuration
func validateSecuritySetup(ctx context.Context, client appsec.APPSEC, configID, configVersion int, policyID string) error {
	// Check configuration exists
	_, err := client.GetConfiguration(ctx, appsec.GetConfigurationRequest{
		ConfigID: configID,
	})
	if err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Check policy exists
	_, err = client.GetSecurityPolicy(ctx, appsec.GetSecurityPolicyRequest{
		ConfigID: configID,
		Version:  configVersion,
		PolicyID: policyID,
	})
	if err != nil {
		return fmt.Errorf("policy validation failed: %w", err)
	}

	// Check at least one match target exists
	targets, err := client.GetMatchTargets(ctx, appsec.GetMatchTargetsRequest{
		ConfigID:      configID,
		ConfigVersion: configVersion,
	})
	if err != nil {
		return fmt.Errorf("match target validation failed: %w", err)
	}

	websiteCount := len(targets.MatchTargets.WebsiteTargets)
	apiCount := len(targets.MatchTargets.APITargets)

	if websiteCount == 0 && apiCount == 0 {
		return fmt.Errorf("no match targets configured - hostname not linked to policy")
	}

	return nil
}
