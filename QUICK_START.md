# Quick Start Guide - Akamai Property Helpers

## TL;DR - Common Commands

```go
// Create Ion Standard (Website)
helper := NewPropertyHelper(ctx, papiClient)
prop, err := helper.CreateIonStandardProperty(PropertyConfig{
    PropertyName:   "my-website",
    ContractID:     "ctr_C-1234567",
    GroupID:        "grp_12345",
    Domain:         "www.example.com",
    OriginHostname: "origin.example.com",
    CPCode:         123456,
    CacheTTL:       86400,
})

// Create Download Delivery (Large Files)
prop, err := helper.CreateDownloadDeliveryProperty(config)

// Create Dynamic Site Delivery (APIs/Apps)
prop, err := helper.CreateDynamicSiteDeliveryProperty(config)

// Create Media Delivery (Video Streaming)
prop, err := helper.CreateMediaDeliveryProperty(config)
```

## Quick Property Type Selector

| **Your Use Case** | **Property Type** | **Method** |
|-------------------|-------------------|------------|
| Corporate website | Ion Standard | `CreateIonStandardProperty()` |
| Blog or news site | Ion Standard | `CreateIonStandardProperty()` |
| E-commerce site | Dynamic Site Delivery | `CreateDynamicSiteDeliveryProperty()` |
| REST API | Dynamic Site Delivery | `CreateDynamicSiteDeliveryProperty()` |
| Software downloads | Download Delivery | `CreateDownloadDeliveryProperty()` |
| Large file hosting | Download Delivery | `CreateDownloadDeliveryProperty()` |
| Video streaming | Media Delivery | `CreateMediaDeliveryProperty()` |
| Live streaming | Media Delivery | `CreateMediaDeliveryProperty()` |
| Mobile app backend | Dynamic Site Delivery | `CreateDynamicSiteDeliveryProperty()` |
| Static website | Ion Standard | `CreateIonStandardProperty()` |

## Common Workflows

### 1. Create and Activate a Website (5 steps)

```go
helper := NewPropertyHelper(ctx, papiClient)

// Step 1: Create property
config := PropertyConfig{
    PropertyName:      "production-website",
    ContractID:        ContractID,
    GroupID:           GroupID,
    Domain:            "www.example.com",
    OriginHostname:    "origin.example.com",
    CPCode:            123456,
    EnableCompression: true,
    EnableHTTP2:       true,
    CacheTTL:          86400,
}
prop, err := helper.CreateIonStandardProperty(config)

// Step 2: Add hostname
edgeHostname := "www.example.com.edgekey.net"
err = helper.AddHostnameToProperty(prop, config, edgeHostname)

// Step 3: Activate to staging
emails := []string{"ops@example.com"}
_, err = helper.ActivateProperty(prop, config, papi.ActivationNetworkStaging, emails)

// Step 4: Test on staging (manual step)
// Test your website: www.example.com.edgesuite-staging.net

// Step 5: Activate to production
_, err = helper.ActivateProperty(prop, config, papi.ActivationNetworkProduction, emails)
```

### 2. Create Multiple Property Types

```go
helper := NewPropertyHelper(ctx, papiClient)

// Website
webConfig := PropertyConfig{
    PropertyName:   "website",
    ContractID:     ContractID,
    GroupID:        GroupID,
    Domain:         "www.example.com",
    OriginHostname: "origin-www.example.com",
    CPCode:         100001,
    CacheTTL:       86400,
}
webProp, _ := helper.CreateIonStandardProperty(webConfig)

// Downloads
downloadConfig := PropertyConfig{
    PropertyName:   "downloads",
    ContractID:     ContractID,
    GroupID:        GroupID,
    Domain:         "downloads.example.com",
    OriginHostname: "origin-downloads.example.com",
    CPCode:         100002,
    CacheTTL:       604800, // 7 days
}
dlProp, _ := helper.CreateDownloadDeliveryProperty(downloadConfig)

// API
apiConfig := PropertyConfig{
    PropertyName:   "api",
    ContractID:     ContractID,
    GroupID:        GroupID,
    Domain:         "api.example.com",
    OriginHostname: "origin-api.example.com",
    CPCode:         100003,
    CacheTTL:       0, // No cache
}
apiProp, _ := helper.CreateDynamicSiteDeliveryProperty(apiConfig)
```

### 3. Add Custom Behaviors

```go
config := PropertyConfig{
    PropertyName:   "custom-site",
    ContractID:     ContractID,
    GroupID:        GroupID,
    Domain:         "www.example.com",
    OriginHostname: "origin.example.com",
    CPCode:         123456,
    CustomBehaviors: []papi.RuleBehavior{
        // Add CORS headers
        {
            Name: "modifyOutgoingResponseHeader",
            Options: papi.RuleOptionsMap{
                "action":                "ADD",
                "standardAddHeaderName": "OTHER",
                "customHeaderName":      "Access-Control-Allow-Origin",
                "headerValue":           "*",
            },
        },
        // Force GZIP
        {
            Name: "gzipResponse",
            Options: papi.RuleOptionsMap{
                "behavior": "ALWAYS",
            },
        },
    },
}

prop, _ := helper.CreateIonStandardProperty(config)
```

## Product ID Reference

| **Property Type** | **Product ID** | **Akamai Product Name** |
|-------------------|----------------|-------------------------|
| Ion Standard | `prd_Ion` | Ion Standard / Akamai Ion |
| Download Delivery | `prd_Download_Delivery` | Download Delivery |
| Dynamic Site Delivery | `prd_Site_Accel` | Dynamic Site Accelerator |
| Media Delivery | `prd_Adaptive_Media_Delivery` | Adaptive Media Delivery |
| Object Delivery | `prd_Object_Delivery` | Object Delivery |
| API Acceleration | `prd_API_Acceleration` | API Acceleration |

## Configuration Cheat Sheet

### Ion Standard (Websites)
```go
PropertyConfig{
    EnableCompression: true,    // ✓ Enable
    EnableHTTP2:       true,    // ✓ Enable
    EnableIPv6:        true,    // ✓ Enable
    CacheTTL:          86400,   // 1 day
}
// Includes: SureRoute, Prefetch, Advanced Caching
```

### Download Delivery (Large Files)
```go
PropertyConfig{
    EnableCompression: false,   // ✗ Disable
    EnableHTTP2:       true,    // ✓ Enable
    CacheTTL:          604800,  // 7 days
}
// Includes: Large File Optimization, Partial Object Caching
```

### Dynamic Site Delivery (APIs)
```go
PropertyConfig{
    EnableCompression: true,    // ✓ Enable
    EnableHTTP2:       true,    // ✓ Enable
    CacheTTL:          0,       // No cache / NO_STORE
}
// Includes: True Client IP, Real User Monitoring, All HTTP Methods
```

### Media Delivery (Streaming)
```go
PropertyConfig{
    EnableCompression: false,   // ✗ Disable
    EnableHTTP2:       true,    // ✓ Enable
    CacheTTL:          2592000, // 30 days
}
// Includes: Adaptive Media Delivery, Segment Protection
```

## Common Cache TTL Values

| **TTL** | **Seconds** | **Use Case** |
|---------|-------------|--------------|
| No cache | 0 | Dynamic/personalized content |
| 5 minutes | 300 | Frequently updated content |
| 1 hour | 3600 | Semi-dynamic content |
| 1 day | 86400 | Daily updated content |
| 7 days | 604800 | Weekly updated / downloads |
| 30 days | 2592000 | Rarely changed / media |
| 1 year | 31536000 | Static assets (versioned) |

## Before You Start Checklist

- [ ] Akamai contract ID (format: `ctr_C-1234567`)
- [ ] Group ID (format: `grp_12345`)
- [ ] CP Code created and assigned to contract
- [ ] Origin hostname accessible
- [ ] Certificate enrollment (for HTTPS)
- [ ] `.edgerc` file configured
- [ ] Product provisioned on contract

## Getting Your IDs

```bash
# List contracts
akamai property-manager list-contracts

# List groups
akamai property-manager list-groups --contract ctr_C-1234567

# List CP Codes
akamai property-manager list-cpcodes --contract ctr_C-1234567 --group grp_12345

# List products
akamai property-manager list-products --contract ctr_C-1234567
```

Or via API:

```go
// Get contracts
contracts, err := papiClient.GetContracts(ctx)

// Get groups
groups, err := papiClient.GetGroups(ctx)

// Get CP codes
cpcodes, err := papiClient.GetCPCodes(ctx, papi.GetCPCodesRequest{
    ContractID: "ctr_C-1234567",
    GroupID:    "grp_12345",
})
```

## Error Handling

```go
prop, err := helper.CreateIonStandardProperty(config)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "already exists"):
        fmt.Println("Property exists, continuing...")
        // Property already exists, fetch it
        
    case strings.Contains(err.Error(), "not found"):
        fmt.Println("Resource not found (CP code, contract, etc.)")
        return err
        
    case strings.Contains(err.Error(), "validation"):
        fmt.Println("Validation error in configuration")
        return err
        
    default:
        fmt.Printf("Unexpected error: %v\n", err)
        return err
    }
}
```

## Testing Your Configuration

### Staging Network
After activating to staging, test using the staging hostname:
```
curl -H "Host: www.example.com" https://www.example.com.edgesuite-staging.net
```

### Production Network
After activating to production:
```
curl https://www.example.com
```

### Check Akamai Headers
```bash
curl -I https://www.example.com
# Look for:
# X-Cache: TCP_HIT from ... (cache hit)
# X-Cache-Key: ... (cache key used)
```

## Next Steps

1. ✅ Created property helpers
2. 📖 Read `PROPERTY_HELPERS.md` for detailed docs
3. 💡 Check `examples.go` for complete examples
4. 🚀 Run your onboarding: `./akamai-onboard`
5. 📊 Monitor in Akamai Control Center
6. 🔧 Customize behaviors as needed

## Common Issues

**"Property already exists"**
- This is expected on re-runs
- Helper will return existing property
- Safe to continue

**"CP Code not found"**
- Create CP code first via Control Center or API
- Ensure CP code is assigned to contract

**"Invalid product ID"**
- Check product is provisioned on contract
- Verify product ID spelling

**"Certificate error"**
- Ensure certificate enrollment exists
- Check enrollment ID is correct
- Verify certificate is deployed

**"Validation errors"**
- Check all required fields in PropertyConfig
- Verify origin hostname is accessible
- Review behavior options

## Support

- 📚 Full Documentation: `PROPERTY_HELPERS.md`
- 💻 Code Examples: `examples.go`
- 🔗 Akamai Docs: https://techdocs.akamai.com/property-mgr
- 🐛 SDK Issues: https://github.com/akamai/AkamaiOPEN-edgegrid-golang

## Performance Tips

1. **Reuse PropertyHelper**: Create once, use for multiple operations
2. **Batch Operations**: Create multiple properties before activating
3. **Parallel Activations**: Activate to staging in parallel for multiple properties
4. **Monitor Progress**: Use `GetActivation()` to check activation status

```go
// Efficient: Reuse helper
helper := NewPropertyHelper(ctx, papiClient)
prop1, _ := helper.CreateIonStandardProperty(config1)
prop2, _ := helper.CreateIonStandardProperty(config2)

// Check activation status
activation, err := papiClient.GetActivation(ctx, papi.GetActivationRequest{
    PropertyID:   prop.PropertyID,
    ActivationID: activationResp.ActivationID,
    ContractID:   ContractID,
    GroupID:      GroupID,
})
fmt.Printf("Status: %s\n", activation.Activation.Status)
```
