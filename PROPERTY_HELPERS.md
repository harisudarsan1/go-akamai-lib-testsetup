# Akamai Property Helpers

This package provides helper functions for creating and configuring different types of Akamai properties using the EdgeGrid Golang SDK v12.

## Overview

The property helpers simplify the process of creating, configuring, and managing Akamai properties by providing pre-configured templates for common use cases.

## Supported Property Types

### 1. Ion Standard (Web Performance)
- **Product ID**: `prd_Ion`
- **Best For**: General websites, web applications, dynamic content
- **Key Features**:
  - SureRoute for optimal routing
  - Prefetch for faster page loads
  - HTTP/2 support
  - Advanced caching
  - Compression

### 2. Download Delivery
- **Product ID**: `prd_Download_Delivery`
- **Best For**: Large file downloads, software distribution, media files
- **Key Features**:
  - Large file optimization
  - Partial object caching
  - Aggressive caching (7 days default)
  - Prefetchable content
  - POST support for downloads

### 3. Dynamic Site Delivery (DSA)
- **Product ID**: `prd_Site_Accel`
- **Best For**: Highly dynamic websites, personalized content, e-commerce
- **Key Features**:
  - True Client IP forwarding
  - Conservative caching
  - Real User Monitoring (RUM)
  - All HTTP methods support
  - HTTP/2 support

### 4. Adaptive Media Delivery
- **Product ID**: `prd_Adaptive_Media_Delivery`
- **Best For**: Video streaming, HLS/DASH delivery, live streaming
- **Key Features**:
  - Adaptive media delivery
  - Segment delivery support
  - Long cache TTL (30 days)
  - No compression (preserves media quality)
  - HTTP/2 for better streaming

## Installation

The property helpers are part of the main package. No additional installation required.

## Usage

### Basic Setup

```go
import (
    "context"
    "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
)

// Initialize the helper
ctx := context.Background()
papiClient := papi.Client(sess)
helper := NewPropertyHelper(ctx, papiClient)
```

### Creating an Ion Standard Property

```go
config := PropertyConfig{
    PropertyName:      "my-website",
    ContractID:        "ctr_C-1234567",
    GroupID:           "grp_12345",
    Domain:            "www.example.com",
    OriginHostname:    "origin.example.com",
    CPCode:            123456,
    EnableCompression: true,
    EnableHTTP2:       true,
    EnableIPv6:        true,
    CacheTTL:          86400, // 1 day in seconds
}

prop, err := helper.CreateIonStandardProperty(config)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created property: %s (ID: %s)\n", prop.PropertyName, prop.PropertyID)
```

### Creating a Download Delivery Property

```go
config := PropertyConfig{
    PropertyName:      "downloads-site",
    ContractID:        "ctr_C-1234567",
    GroupID:           "grp_12345",
    Domain:            "downloads.example.com",
    OriginHostname:    "origin-downloads.example.com",
    CPCode:            123457,
    EnableCompression: false, // Typically disabled for downloads
    EnableHTTP2:       true,
    CacheTTL:          604800, // 7 days
}

prop, err := helper.CreateDownloadDeliveryProperty(config)
```

### Creating a Dynamic Site Delivery Property

```go
config := PropertyConfig{
    PropertyName:      "dynamic-app",
    ContractID:        "ctr_C-1234567",
    GroupID:           "grp_12345",
    Domain:            "app.example.com",
    OriginHostname:    "origin-app.example.com",
    CPCode:            123458,
    EnableCompression: true,
    EnableHTTP2:       true,
    CacheTTL:          0, // No caching for dynamic content
}

prop, err := helper.CreateDynamicSiteDeliveryProperty(config)
```

### Creating a Media Delivery Property

```go
config := PropertyConfig{
    PropertyName:      "video-streaming",
    ContractID:        "ctr_C-1234567",
    GroupID:           "grp_12345",
    Domain:            "media.example.com",
    OriginHostname:    "origin-media.example.com",
    CPCode:            123459,
    EnableCompression: false, // Don't compress media
    EnableHTTP2:       true,
    CacheTTL:          2592000, // 30 days
}

prop, err := helper.CreateMediaDeliveryProperty(config)
```

### Adding Custom Behaviors

```go
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

config.CustomBehaviors = customBehaviors
prop, err := helper.CreateIonStandardProperty(config)
```

### Adding Hostnames to Properties

```go
edgeHostname := "www.example.com.edgekey.net"
err := helper.AddHostnameToProperty(prop, config, edgeHostname)
if err != nil {
    log.Fatal(err)
}
```

### Activating Properties

```go
// Activate to staging
notifyEmails := []string{"ops@example.com"}
activationResp, err := helper.ActivateProperty(
    prop,
    config,
    papi.ActivationNetworkStaging,
    notifyEmails,
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Activation ID: %s\n", activationResp.ActivationID)
fmt.Printf("Activation Link: %s\n", activationResp.ActivationLink)

// After testing on staging, activate to production
activationResp, err = helper.ActivateProperty(
    prop,
    config,
    papi.ActivationNetworkProduction,
    notifyEmails,
)
```

## Complete Onboarding Example

```go
func CompleteOnboarding(ctx context.Context, sess session.Session) error {
    papiClient := papi.Client(sess)
    helper := NewPropertyHelper(ctx, papiClient)
    
    // Step 1: Create property
    config := PropertyConfig{
        PropertyName:      "production-site",
        ContractID:        "ctr_C-1234567",
        GroupID:           "grp_12345",
        Domain:            "www.example.com",
        OriginHostname:    "origin.example.com",
        CPCode:            123456,
        EnableCompression: true,
        EnableHTTP2:       true,
        CacheTTL:          86400,
    }
    
    prop, err := helper.CreateIonStandardProperty(config)
    if err != nil {
        return err
    }
    
    // Step 2: Add hostname
    edgeHostname := "www.example.com.edgekey.net"
    err = helper.AddHostnameToProperty(prop, config, edgeHostname)
    if err != nil {
        return err
    }
    
    // Step 3: Activate to staging
    notifyEmails := []string{"ops@example.com"}
    _, err = helper.ActivateProperty(
        prop,
        config,
        papi.ActivationNetworkStaging,
        notifyEmails,
    )
    if err != nil {
        return err
    }
    
    fmt.Println("Property created and activated to staging!")
    return nil
}
```

## PropertyConfig Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `PropertyName` | string | Yes | Unique name for the property |
| `ContractID` | string | Yes | Akamai contract ID (e.g., "ctr_C-1234567") |
| `GroupID` | string | Yes | Akamai group ID (e.g., "grp_12345") |
| `ProductID` | string | Auto | Set automatically based on property type |
| `Domain` | string | Yes | The hostname to serve (e.g., "www.example.com") |
| `OriginHostname` | string | Yes | Your origin server hostname |
| `CPCode` | int | Yes | CP Code for reporting and billing |
| `EdgeHostnameID` | string | No | Edge hostname ID (if known) |
| `CertEnrollmentID` | int | No | Certificate enrollment ID for HTTPS |
| `EnableCompression` | bool | No | Enable gzip compression (default varies by type) |
| `EnableHTTP2` | bool | No | Enable HTTP/2 protocol |
| `EnableIPv6` | bool | No | Enable IPv6 support |
| `CacheTTL` | int | No | Cache time-to-live in seconds |
| `CustomBehaviors` | []RuleBehavior | No | Additional custom behaviors |

## Property Type Defaults

### Ion Standard Defaults
- Compression: Enabled
- HTTP/2: Enabled
- Cache TTL: 1 day (86400s)
- SureRoute: Enabled
- Prefetch: Enabled

### Download Delivery Defaults
- Compression: Disabled
- HTTP/2: Enabled
- Cache TTL: 7 days (604800s)
- Large File Optimization: Enabled
- Partial Object Caching: Enabled

### Dynamic Site Delivery Defaults
- Compression: Enabled
- HTTP/2: Enabled
- Cache TTL: No Store (0)
- True Client IP: Enabled
- Real User Monitoring: Enabled

### Media Delivery Defaults
- Compression: Disabled
- HTTP/2: Enabled
- Cache TTL: 30 days (2592000s)
- Adaptive Media Delivery: Enabled

## Advanced Features

### Behavior Merging

The helpers automatically merge behaviors, preventing duplicates:
- Existing behaviors are preserved
- New behaviors override existing ones with the same name
- Custom behaviors are added last

### Idempotency

All operations are idempotent:
- Creating a property that already exists returns the existing property
- Adding a hostname that already exists skips the operation
- Safe to run multiple times

### Error Handling

All functions return detailed errors:
```go
prop, err := helper.CreateIonStandardProperty(config)
if err != nil {
    if strings.Contains(err.Error(), "already exists") {
        // Handle existing property
    } else {
        // Handle other errors
        log.Fatal(err)
    }
}
```

## Best Practices

1. **CP Codes**: Create separate CP codes for different properties or traffic types for better reporting
2. **Cache TTL**: Set appropriate cache TTLs based on your content update frequency
3. **Compression**: Enable for text-based content, disable for already-compressed formats
4. **Testing**: Always test on staging before activating to production
5. **Notifications**: Include relevant email addresses for activation notifications
6. **Origin**: Use meaningful origin hostnames that match your infrastructure

## Common Behaviors Reference

Here are some common behaviors you can add via `CustomBehaviors`:

```go
// Add custom response header
{
    Name: "modifyOutgoingResponseHeader",
    Options: papi.RuleOptionsMap{
        "action": "ADD",
        "standardAddHeaderName": "OTHER",
        "customHeaderName": "X-Custom-Header",
        "headerValue": "value",
    },
}

// CORS headers
{
    Name: "modifyOutgoingResponseHeader",
    Options: papi.RuleOptionsMap{
        "action": "ADD",
        "standardAddHeaderName": "OTHER",
        "customHeaderName": "Access-Control-Allow-Origin",
        "headerValue": "*",
    },
}

// Redirect HTTP to HTTPS
{
    Name: "redirectplus",
    Options: papi.RuleOptionsMap{
        "enabled": true,
        "destination": "https://[Host][URI]",
        "responseCode": 301,
    },
}

// Set cache control headers
{
    Name: "cacheKeyQueryParams",
    Options: papi.RuleOptionsMap{
        "behavior": "INCLUDE_ALL_PRESERVE_ORDER",
    },
}
```

## Troubleshooting

### Property Creation Fails
- Verify contract ID and group ID are correct
- Ensure CP code exists and is assigned to the contract
- Check that product is available on your contract

### Hostname Addition Fails
- Verify edge hostname exists and is correct
- Check certificate enrollment is valid
- Ensure hostname format is correct

### Activation Fails
- Review property validation errors
- Check that all required behaviors are present
- Verify notification emails are valid

## See Also

- `examples.go` - Complete working examples
- `main.go` - Integration with onboarding flow
- [Akamai Property Manager API](https://techdocs.akamai.com/property-mgr/reference/api)
- [Akamai EdgeGrid SDK](https://github.com/akamai/AkamaiOPEN-edgegrid-golang)
