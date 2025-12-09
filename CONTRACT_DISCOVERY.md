# Contract Discovery Guide

This document explains how to automatically discover your Akamai Contract ID, Group ID, and Product IDs using the provided test tools.

---

## Overview

The contract discovery feature automatically retrieves your Akamai configuration by querying the API using your Group Name. This eliminates the need to manually find Contract IDs and Product IDs.

### What Gets Discovered

✅ **Contract ID** - Your Akamai contract identifier (e.g., `ctr_V-620VL0G`)  
✅ **Group ID** - Your property group identifier (e.g., `grp_304920`)  
✅ **Product IDs** - Available products for your contract (e.g., `prd_Download_Delivery`)  
✅ **Configuration Caching** - Results saved to `~/.akamai-config.json` for reuse

---

## Prerequisites

1. **Valid `.edgerc` file** with API credentials at `~/.edgerc`
2. **Group Name** - The exact name of your Akamai property group

### Setting up `.edgerc`

Your `.edgerc` file should look like this:

```ini
[default]
client_secret = YOUR_CLIENT_SECRET
host = YOUR_API_HOST.luna.akamaiapis.net
access_token = YOUR_ACCESS_TOKEN
client_token = YOUR_CLIENT_TOKEN
```

**File location**: `~/.edgerc` (or `/Users/yourusername/.edgerc` on macOS)  
**Permissions**: `chmod 600 ~/.edgerc` (readable only by you)

---

## Finding Your Group Name

### Option 1: Via Akamai Control Center

1. Login to https://control.akamai.com
2. Navigate to **Property Manager**
3. Look at your properties list
4. The **Group** column shows your group name
5. Example: `CTM LABS PRIVATE LIMITED (Kluisz)-V-620VL0G`

### Option 2: List All Groups via Test

Run the discovery test to see all available groups:

```bash
cd go-akamai-waf-test
go test -v -run TestListAllContractsAndGroups
```

**Output:**
```
=== Available Groups ===
1. CTM LABS PRIVATE LIMITED (Kluisz)-V-620VL0G (grp_304920)
   Contract(s): [ctr_V-620VL0G]
```

---

## Running the Discovery Test

### Test 1: Discover Contract by Group Name

This is the **primary test** that discovers your configuration:

```bash
cd go-akamai-waf-test
go test -v -run TestDiscoverContractByGroupName
```

### Expected Output

```
=== RUN   TestDiscoverContractByGroupName
    === Akamai Contract Discovery Test ===
    
    Step 1: Authenticating with Akamai API...
    ✅ Authentication successful
    
    Step 2: Discovering Contract for Group: "CTM LABS PRIVATE LIMITED (Kluisz)-V-620VL0G"
🔍 Searching for Group: "CTM LABS PRIVATE LIMITED (Kluisz)-V-620VL0G"
✅ Found matching group!
    
    Step 3: Validating discovered information...
    ✅ Discovered information validated
    
    Step 4: Discovering Product IDs...
🔍 Discovering Product IDs for Contract: ctr_V-620VL0G

📦 Available Products:
   1. Download_Delivery (prd_Download_Delivery)
    ✅ Found 1 product(s)
    
    Step 5: Saving configuration to cache...

💾 Configuration saved to: /Users/hari/.akamai-config.json
    ✅ Configuration saved successfully
    
    =================================

=== Discovered Information ===
Group Name:    CTM LABS PRIVATE LIMITED (Kluisz)-V-620VL0G
Group ID:      grp_304920
Contract ID:   ctr_V-620VL0G
Contract Name: DIRECT_CUSTOMER

Available Product IDs (1):
  1. prd_Download_Delivery

=== Usage in main.go ===
const ContractID = "ctr_V-620VL0G"
const GroupID    = "grp_304920"
const ProductID  = "prd_Download_Delivery"  // or choose from available products
    =================================
    
    ✅ All tests passed!
--- PASS: TestDiscoverContractByGroupName (3.35s)
```

### What the Test Does

1. **Authenticates** with Akamai API using your `.edgerc` credentials
2. **Searches** for your group by exact name match
3. **Retrieves** the associated Contract ID
4. **Discovers** available Product IDs for that contract
5. **Saves** configuration to `~/.akamai-config.json` (cached)
6. **Displays** discovered values in a ready-to-use format

---

## Configuration Cache

### Cache File Location

`~/.akamai-config.json`

### Cache File Format

```json
{
  "contractId": "ctr_V-620VL0G",
  "contractName": "DIRECT_CUSTOMER",
  "groupId": "grp_304920",
  "groupName": "CTM LABS PRIVATE LIMITED (Kluisz)-V-620VL0G",
  "productIds": [
    "prd_Download_Delivery"
  ]
}
```

### Cache Benefits

✅ **Faster Execution** - Subsequent runs use cached values (no API calls)  
✅ **Offline Development** - Work without active internet connection  
✅ **Consistency** - Same values across multiple runs  

### Clearing Cache

To force re-discovery, delete the cache file:

```bash
rm ~/.akamai-config.json
```

Next test run will query the API and create a fresh cache.

---

## Additional Tests

### Test 2: List All Contracts and Groups

View all contracts and groups accessible to your credentials:

```bash
go test -v -run TestListAllContractsAndGroups
```

**Use Case**: 
- Explore available resources
- Find correct group name if you're unsure
- Verify API access to multiple contracts

### Test 3: Test Cache Workflow

Verify caching works correctly:

```bash
go test -v -run TestDiscoverAndCache
```

**First run**: Queries API, saves to cache  
**Second run**: Loads from cache (faster)

---

## Using Discovered Values in `main.go`

The `main.go` file is already configured to **auto-discover** values on startup!

### How It Works

```go
const (
    // Group name to auto-discover Contract ID and Product IDs
    GroupName = "CTM LABS PRIVATE LIMITED (Kluisz)-V-620VL0G"
    
    // Property configuration (customize these)
    PropertyName = "my-api-gateway"
    UserDomain   = "api.example.com"
)

func main() {
    // ...
    
    // STEP 0: Auto-discover Contract ID, Group ID, and Product IDs
    config, err := DiscoverAndCache(ctx, papiClient, GroupName)
    if err != nil {
        log.Fatalf("Failed to discover contract information: %v", err)
    }
    
    // Set runtime configuration
    ContractID = config.ContractID
    GroupID = config.GroupID
    ProductID = config.ProductIDs[0]  // Use first available product
    
    // ... continue with onboarding
}
```

### Configuration Steps

1. **Edit `main.go`** - Update `GroupName` constant with your group name
2. **Run the program** - It will auto-discover and cache Contract/Product IDs
3. **Subsequent runs** - Uses cached values for faster startup

### Manual Override

If you prefer to hardcode values:

```go
// Option 1: Use auto-discovery (recommended)
const GroupName = "YOUR_GROUP_NAME"

// Option 2: Hardcode values (not recommended)
var (
    ContractID = "ctr_V-620VL0G"
    GroupID    = "grp_304920"
    ProductID  = "prd_Download_Delivery"
)
```

---

## Troubleshooting

### Error: "authentication failed"

**Cause**: Invalid or missing `.edgerc` credentials

**Solution**:
1. Verify `.edgerc` exists: `ls -la ~/.edgerc`
2. Check file format (must have `[default]` section header)
3. Verify permissions: `chmod 600 ~/.edgerc`
4. Test credentials in Akamai Control Center

### Error: "group not found"

**Cause**: Group name doesn't match exactly

**Solution**:
1. Run `go test -v -run TestListAllContractsAndGroups` to see available groups
2. Copy the **exact** group name (case-sensitive, including special characters)
3. Update `GroupName` constant in `main.go` or test file

**Example**:
```go
// ❌ Wrong - missing parentheses
const GroupName = "CTM LABS PRIVATE LIMITED Kluisz-V-620VL0G"

// ✅ Correct - exact match
const GroupName = "CTM LABS PRIVATE LIMITED (Kluisz)-V-620VL0G"
```

### Error: "no products found"

**Cause**: Contract has no provisioned products

**Solution**:
1. Contact your Akamai account manager
2. Request product provisioning for your contract
3. Common products:
   - `prd_Ion` - Ion Standard (web delivery)
   - `prd_Download_Delivery` - Large file downloads
   - `prd_Dynamic_Site_Delivery` - Dynamic content/APIs

### Error: "failed to read config file"

**Cause**: Cache file corrupted or permissions issue

**Solution**:
```bash
# Delete corrupted cache
rm ~/.akamai-config.json

# Re-run discovery test
go test -v -run TestDiscoverContractByGroupName
```

### Warning: "Failed to discover products"

**Cause**: Products API call failed (non-fatal)

**Impact**: Test continues, but ProductIDs list will be empty

**Solution**:
- Check API permissions (ensure PAPI access)
- Verify contract has products provisioned
- Contact Akamai support if persistent

---

## API Permissions Required

Your API credentials must have these permissions:

| API | Permission | Used For |
|-----|------------|----------|
| **PAPI** | READ | List contracts, groups, products |
| **Property Manager** | READ | Access property configuration |

**Note**: **WRITE** permissions are NOT required for discovery tests (read-only operations).

---

## Workflow Summary

```
┌─────────────────────────────────────────┐
│  1. Run Discovery Test                  │
│     go test -v -run TestDiscover...     │
└─────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────┐
│  2. Authenticate with .edgerc           │
│     Load credentials from ~/.edgerc     │
└─────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────┐
│  3. Query Akamai API                    │
│     GET /papi/v1/groups                 │
│     GET /papi/v1/contracts              │
│     GET /papi/v1/products?contractId=X  │
└─────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────┐
│  4. Save to Cache                       │
│     Write to ~/.akamai-config.json      │
└─────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────┐
│  5. Display Results                     │
│     Contract ID, Group ID, Products     │
└─────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────┐
│  6. Use in main.go                      │
│     Auto-load from cache on startup     │
└─────────────────────────────────────────┘
```

---

## Files Created

### Discovery Implementation

- **`contract_discovery.go`** - Core discovery functions
  - `DiscoverContractByGroupName()` - Find contract by group name
  - `DiscoverContractByGroupID()` - Find contract by group ID
  - `DiscoverProductIDs()` - Find available products
  - `DiscoverAndCache()` - Full discovery + caching workflow
  - `SaveConfig()` / `LoadConfig()` - Cache management

### Test Files

- **`contract_discovery_test.go`** - Go tests
  - `TestDiscoverContractByGroupName` - Main discovery test
  - `TestListAllContractsAndGroups` - List all resources
  - `TestDiscoverAndCache` - Cache workflow test

### Updated Files

- **`main.go`** - Updated to use auto-discovery
  - Removed hardcoded Contract/Group/Product IDs
  - Added `DiscoverAndCache()` call on startup
  - Uses runtime configuration variables

---

## Next Steps

After running the discovery test:

1. ✅ **Verify discovered values** match your Akamai account
2. ✅ **Update `main.go`** with your Group Name (if different)
3. ✅ **Customize property settings** (PropertyName, UserDomain)
4. ✅ **Run the onboarding flow** with auto-discovered configuration

---

## Support

### Discovered Values

If you successfully ran the test, your values are:

```
Contract ID: ctr_V-620VL0G
Group ID:    grp_304920
Product ID:  prd_Download_Delivery
```

You can now use these in your onboarding workflow!

### Questions?

- Check `ONBOARDING_GUIDE.md` for detailed onboarding steps
- See `README.md` for project overview
- Review `PRODUCT_IDS.md` for product information

---

**Document Version**: 1.0  
**Last Updated**: 2025-12-09
