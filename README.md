# Akamai Property Management with Go SDK

A comprehensive toolkit for managing Akamai properties using the official AkamaiOPEN-edgegrid-golang v12 SDK.

## 📚 Documentation

- **[PRODUCT_IDS.md](PRODUCT_IDS.md)** - Understanding Akamai Product IDs (READ THIS FIRST!)
- **[QUICK_START.md](QUICK_START.md)** - Quick reference guide
- **[PROPERTY_HELPERS.md](PROPERTY_HELPERS.md)** - Detailed API documentation
- **[examples.go](examples.go)** - Working code examples

## 🎯 Features

### Property Type Helpers
Create optimized properties for different use cases:
- ✅ **Ion Standard** - Websites and web applications
- ✅ **Download Delivery** - Large file downloads
- ✅ **Dynamic Site Delivery** - APIs and dynamic content
- ✅ **Adaptive Media Delivery** - Video streaming

### Product Discovery Tools
- 🔍 List all available products for your contract
- 🔍 Find products by name
- 🔍 Validate product IDs
- 🔍 Automatic product ID mapping

### Complete Workflow Support
- Property creation with optimized behaviors
- Hostname management
- Edge hostname configuration
- Property activation (staging/production)
- CPS certificate integration
- AppSec/WAF onboarding

## 🚀 Quick Start

### 1. Install Dependencies

```bash
go mod download
```

### 2. Configure Credentials

Create `~/.edgerc`:
```ini
[default]
client_secret = your_client_secret
host = your_host.luna.akamaiapis.net
access_token = your_access_token
client_token = your_client_token
```

### 3. Discover Your Products (IMPORTANT!)

⚠️ **Before creating properties, discover your contract's product IDs:**

```go
package main

import (
    "context"
    "github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/papi"
)

func main() {
    ctx := context.Background()
    sess, _ := newSession("~/.edgerc", "default")
    papiClient := papi.Client(sess)
    
    // Replace with your contract ID
    DiscoverProducts(ctx, papiClient, "ctr_C-1234567")
}
```

**Output:**
```
Available Products for Contract: ctr_C-1234567
================================================================================
Product Name                                       Product ID
--------------------------------------------------------------------------------
Ion Standard                                       prd_Ion
Download Delivery                                  prd_Download_Delivery
Dynamic Site Accelerator                           prd_Site_Accel
Adaptive Media Delivery                            prd_Adaptive_Media_Delivery
================================================================================
Total Products: 4
```

### 4. Create a Property

```go
helper := NewPropertyHelper(ctx, papiClient)

config := PropertyConfig{
    PropertyName:      "my-website",
    ContractID:        "ctr_C-1234567",
    GroupID:           "grp_12345",
    ProductID:         "prd_Ion",  // Use your actual product ID from step 3
    Domain:            "www.example.com",
    OriginHostname:    "origin.example.com",
    CPCode:            123456,
    EnableCompression: true,
    EnableHTTP2:       true,
    CacheTTL:          86400,
}

prop, err := helper.CreateIonStandardProperty(config)
```

## 📦 Project Structure

```
go-akamai-waf-test/
├── main.go                 # Original onboarding workflow
├── property_helpers.go     # Property type helpers (635 lines)
├── product_utils.go        # Product discovery utilities (271 lines)
├── examples.go             # Working examples (294 lines)
│
├── README.md               # This file
├── PRODUCT_IDS.md          # Product ID reference (IMPORTANT!)
├── QUICK_START.md          # Quick reference
├── PROPERTY_HELPERS.md     # Detailed documentation
│
├── go.mod                  # Go module definition
└── .edgerc                 # Akamai credentials (not in repo)
```

## 🔑 Important: Product IDs

### The SDK Does NOT Have Product ID Constants!

The AkamaiOPEN-edgegrid-golang SDK **does not define product ID constants**. Product IDs are:
- ✅ Contract-specific
- ✅ Fetched via API (`GetProducts`)
- ✅ Visible in Akamai Control Center

### Product IDs Used in This Project

The product IDs in `property_helpers.go` (like `prd_Ion`, `prd_Download_Delivery`) are:
- ❌ **NOT** from the SDK
- ✅ Based on common Akamai naming patterns
- ✅ Referenced from Akamai documentation
- ✅ Examples from Terraform providers
- ⚠️ **May not match your contract**

### Before Using This Code

1. **Run product discovery** for YOUR contract
2. **Verify** product IDs match your contract
3. **Update** helper functions with your product IDs, or
4. **Pass ProductID** explicitly in PropertyConfig

See [PRODUCT_IDS.md](PRODUCT_IDS.md) for complete details.

## 🛠️ Core Components

### 1. PropertyHelper

Main class for property management:

```go
helper := NewPropertyHelper(ctx, papiClient)

// Create different property types
helper.CreateIonStandardProperty(config)
helper.CreateDownloadDeliveryProperty(config)
helper.CreateDynamicSiteDeliveryProperty(config)
helper.CreateMediaDeliveryProperty(config)

// Manage properties
helper.AddHostnameToProperty(prop, config, edgeHostname)
helper.ActivateProperty(prop, config, network, emails)
```

### 2. ProductFetcher

Discover and validate products:

```go
fetcher := NewProductFetcher(ctx, papiClient)

// List all products
fetcher.ListProducts(contractID)

// Find by name
productID, _ := fetcher.FindProductByName(contractID, "ion")

// Validate ID
valid, _ := fetcher.ValidateProductID(contractID, "prd_Ion")
```

### 3. ProductMapper

Map property types to product IDs:

```go
mapper := NewProductMapper(ctx, papiClient, contractID)

// Automatic mapping
productID, _ := mapper.GetProductID(PropertyTypeIonStandard)

// Manual mapping
mapper.SetProductID(PropertyTypeIonStandard, "prd_SPM")
```

## 📖 Usage Examples

### Example 1: Simple Property Creation

```go
helper := NewPropertyHelper(ctx, papiClient)
config := PropertyConfig{
    PropertyName:   "website",
    ContractID:     "ctr_C-1234567",
    GroupID:        "grp_12345",
    Domain:         "www.example.com",
    OriginHostname: "origin.example.com",
    CPCode:         123456,
}

prop, err := helper.CreateIonStandardProperty(config)
```

### Example 2: With Product Discovery

```go
// Step 1: Discover products
fetcher := NewProductFetcher(ctx, papiClient)
productID, _ := fetcher.FindProductByName(contractID, "ion")

// Step 2: Use discovered product ID
config.ProductID = productID
prop, err := helper.CreateProperty(config)
```

### Example 3: Complete Onboarding

```go
helper := NewPropertyHelper(ctx, papiClient)

// 1. Create property
prop, _ := helper.CreateIonStandardProperty(config)

// 2. Add hostname
helper.AddHostnameToProperty(prop, config, "www.example.com.edgekey.net")

// 3. Activate to staging
helper.ActivateProperty(prop, config, papi.ActivationNetworkStaging, emails)
```

See [examples.go](examples.go) for more complete examples.

## 🎨 Customization

### Custom Behaviors

Add custom behaviors to any property:

```go
config.CustomBehaviors = []papi.RuleBehavior{
    {
        Name: "gzipResponse",
        Options: papi.RuleOptionsMap{
            "behavior": "ALWAYS",
        },
    },
}

prop, _ := helper.CreateIonStandardProperty(config)
```

### Custom Product IDs

Override default product IDs:

```go
config := PropertyConfig{
    PropertyName: "my-property",
    ProductID:    "prd_YourCustomProduct",  // Your actual product ID
    // ... other fields
}

prop, _ := helper.CreateProperty(config)  // Generic creation
```

## 🔧 Configuration

### PropertyConfig Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| PropertyName | string | ✅ | Unique property name |
| ContractID | string | ✅ | Akamai contract ID |
| GroupID | string | ✅ | Akamai group ID |
| ProductID | string | ⚠️ | Auto-set by helpers (verify first!) |
| Domain | string | ✅ | Hostname to serve |
| OriginHostname | string | ✅ | Origin server |
| CPCode | int | ✅ | CP code for reporting |
| EnableCompression | bool | ⚠️ | Varies by property type |
| EnableHTTP2 | bool | ⚠️ | Usually true |
| CacheTTL | int | ⚠️ | Cache duration in seconds |
| CustomBehaviors | []RuleBehavior | ❌ | Additional behaviors |

## 🧪 Testing

### Verify Product IDs

```bash
go run . --discover-products --contract ctr_C-1234567
```

### Test Property Creation

```bash
go run . --test-property --contract ctr_C-1234567 --group grp_12345
```

## 📊 Property Type Comparison

| Feature | Ion Standard | Download | Dynamic Site | Media |
|---------|--------------|----------|--------------|-------|
| Compression | ✅ | ❌ | ✅ | ❌ |
| HTTP/2 | ✅ | ✅ | ✅ | ✅ |
| Cache (default) | 1 day | 7 days | None | 30 days |
| SureRoute | ✅ | ❌ | ❌ | ❌ |
| Large Files | ❌ | ✅ | ❌ | ❌ |
| True Client IP | ❌ | ❌ | ✅ | ❌ |
| Adaptive Delivery | ❌ | ❌ | ❌ | ✅ |

## ⚠️ Common Issues

### "Product ID not found"
**Solution:** Run product discovery for your contract first.

### "Property already exists"
**Solution:** This is normal. Helpers return existing property.

### "Invalid product ID"
**Solution:** Verify product ID with `GetProducts()` API.

### "CP Code not found"
**Solution:** Create CP code first via Control Center or API.

## 🔗 Resources

- [Akamai Property Manager API](https://techdocs.akamai.com/property-mgr)
- [EdgeGrid SDK Documentation](https://github.com/akamai/AkamaiOPEN-edgegrid-golang)
- [Akamai Control Center](https://control.akamai.com)
- [Akamai CLI](https://developer.akamai.com/cli)

## 📝 Development

### Build

```bash
go build
```

### Run

```bash
./akamai-onboard
```

### Add New Property Type

1. Add constant to `property_helpers.go`:
   ```go
   PropertyTypeNewType PropertyType = "New_Type"
   ```

2. Create helper function:
   ```go
   func (ph *PropertyHelper) CreateNewTypeProperty(config PropertyConfig) (*papi.Property, error) {
       config.ProductID = "prd_NewType"  // Your product ID
       // ... configuration
   }
   ```

3. Add to product mapper in `product_utils.go`

## 🤝 Contributing

When adding new features:
1. Update relevant documentation
2. Add examples to `examples.go`
3. Test with your contract
4. Verify product IDs

## 📄 License

Apache 2.0 - See LICENSE file

## 🙋 Support

For questions about:
- **Product IDs**: See [PRODUCT_IDS.md](PRODUCT_IDS.md)
- **Quick Usage**: See [QUICK_START.md](QUICK_START.md)
- **Detailed Docs**: See [PROPERTY_HELPERS.md](PROPERTY_HELPERS.md)
- **Examples**: See [examples.go](examples.go)

## 🎯 Next Steps

1. ✅ Read [PRODUCT_IDS.md](PRODUCT_IDS.md) (IMPORTANT!)
2. ✅ Run product discovery for your contract
3. ✅ Verify product IDs match your contract
4. ✅ Review [QUICK_START.md](QUICK_START.md)
5. ✅ Check [examples.go](examples.go)
6. ✅ Create your first property!

---

**Note:** The product IDs in this codebase are examples. Always verify with your Akamai contract before using in production.
