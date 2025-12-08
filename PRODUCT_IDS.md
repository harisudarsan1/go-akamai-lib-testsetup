# Akamai Product IDs Reference

## Important Note: Product IDs Source

**The AkamaiOPEN-edgegrid-golang SDK does NOT contain predefined constants for product IDs.** 

Product IDs are **contract-specific** and must be obtained from:
1. The Akamai Control Center
2. The PAPI API's `GetProducts` endpoint
3. Your Akamai account team

## Where Product IDs Come From

### Source 1: Akamai API (Recommended)
The most reliable way to get product IDs is through the API:

```go
products, err := papiClient.GetProducts(ctx, papi.GetProductsRequest{
    ContractID: "ctr_C-1234567",
})

for _, product := range products.Products.Items {
    fmt.Printf("Product: %s (ID: %s)\n", product.ProductName, product.ProductID)
}
```

### Source 2: Akamai Control Center
1. Log into https://control.akamai.com
2. Go to Properties → New Property
3. Select Contract and Group
4. Available products will be listed
5. Product IDs are visible in the UI or network inspector

### Source 3: CLI Tool
```bash
# List available products for a contract
akamai property-manager list-products --contract ctr_C-1234567
```

## Common Product IDs (Reference Only)

⚠️ **Warning**: These are examples and may not match your contract. Always verify with your actual contract.

### Delivery Products

| Product Name | Common Product ID | Description |
|--------------|-------------------|-------------|
| Ion Standard | `prd_Ion` | Web performance optimization |
| Ion Premier | `prd_Ion_Premier` | Advanced Ion with more features |
| Download Delivery | `prd_Download_Delivery` | Large file delivery |
| Dynamic Site Delivery | `prd_Site_Accel` | Dynamic content acceleration |
| Object Delivery | `prd_Object_Delivery` | Object/file delivery |
| Adaptive Media Delivery | `prd_Adaptive_Media_Delivery` | Video/media streaming |

### Security Products

| Product Name | Common Product ID | Description |
|--------------|-------------------|-------------|
| Web Application Accelerator | `prd_Web_App_Accel` | Web app acceleration |
| Site Defender | `prd_Site_Defender` | DDoS protection |
| Kona Site Defender | `prd_Fresca` | Web application firewall |

### Other Products

| Product Name | Common Product ID | Description |
|--------------|-------------------|-------------|
| Alta | `prd_Alta` | Content delivery |
| Rich Media Accelerator | `prd_Rich_Media_Accel` | Rich media delivery |
| API Acceleration | `prd_API_Acceleration` | API optimization |

## Product ID Format

Product IDs typically follow this pattern:
- Format: `prd_<ProductName>`
- Examples: `prd_Ion`, `prd_Download_Delivery`, `prd_Site_Accel`
- **Case-sensitive**: Must match exactly

## How We Determined Product IDs in the Helpers

The product IDs used in `property_helpers.go` were based on:

1. **Akamai Official Documentation**
   - https://techdocs.akamai.com/property-mgr
   - Property Manager API documentation
   - Product naming conventions

2. **Common Industry Usage**
   - Standard product names used across Akamai implementations
   - Terraform provider examples
   - SDK test files (limited examples)

3. **SDK Test Files Evidence**
   Found in `/pkg/papi/*_test.go`:
   - `prd_Alta`
   - `prd_Site_Defender`
   - `prd_Web_App_Accel`

4. **Akamai Terraform Provider**
   - Cross-referenced with Terraform documentation
   - Validated product ID patterns

## Verifying Your Product IDs

### Step 1: Create a Product Fetcher

```go
func GetAvailableProducts(ctx context.Context, papiClient papi.PAPI, contractID string) {
    products, err := papiClient.GetProducts(ctx, papi.GetProductsRequest{
        ContractID: contractID,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Available Products:")
    fmt.Println("==================")
    for _, product := range products.Products.Items {
        fmt.Printf("Name: %-40s ID: %s\n", product.ProductName, product.ProductID)
    }
}
```

### Step 2: Run Against Your Contract

```go
func main() {
    ctx := context.Background()
    sess, _ := newSession("~/.edgerc", "default")
    papiClient := papi.Client(sess)
    
    // Replace with your contract ID
    GetAvailableProducts(ctx, papiClient, "ctr_C-1234567")
}
```

### Example Output

```
Available Products:
==================
Name: Ion Standard                               ID: prd_Ion
Name: Download Delivery                          ID: prd_Download_Delivery
Name: Dynamic Site Accelerator                   ID: prd_Site_Accel
Name: Adaptive Media Delivery                    ID: prd_Adaptive_Media_Delivery
Name: Object Delivery                            ID: prd_Object_Delivery
Name: API Acceleration                           ID: prd_API_Acceleration
```

## Updating Product IDs in Helpers

If your contract uses different product IDs, update the constants in `property_helpers.go`:

### Before (Example IDs)
```go
func (ph *PropertyHelper) CreateIonStandardProperty(config PropertyConfig) (*papi.Property, error) {
    config.ProductID = "prd_Ion"  // May not match your contract
    // ...
}
```

### After (Your Contract's IDs)
```go
func (ph *PropertyHelper) CreateIonStandardProperty(config PropertyConfig) (*papi.Property, error) {
    config.ProductID = "prd_SPM"  // Your actual product ID
    // ...
}
```

Or better yet, pass it in the config:

```go
config := PropertyConfig{
    PropertyName: "my-property",
    ProductID:    "prd_SPM",  // Explicitly set your product ID
    // ... other fields
}

prop, err := helper.CreateProperty(config)  // Use generic CreateProperty
```

## Product ID Best Practices

### ✅ Do This
- **Always fetch from API** for your specific contract
- **Store in configuration** files or environment variables
- **Validate before use** with GetProducts
- **Document** which product ID you're using
- **Use exact casing** (IDs are case-sensitive)

### ❌ Don't Do This
- **Don't hardcode** product IDs from examples
- **Don't assume** product IDs are universal
- **Don't guess** product ID formats
- **Don't use** test product IDs in production

## Contract-Specific Nature

### Why Product IDs Vary

1. **Contract Differences**: Each contract may have different products
2. **Product Bundles**: Some contracts have bundled products
3. **Custom Products**: Enterprise contracts may have custom product IDs
4. **Regional Variations**: Product IDs may vary by region
5. **Product Evolution**: Akamai may deprecate/rename products

### Example: Same Product, Different IDs

| Customer | Product | Product ID |
|----------|---------|------------|
| Customer A | Ion | `prd_Ion` |
| Customer B | Ion | `prd_SPM` |
| Customer C | Ion Premier | `prd_Ion_Premier` |

All three might be "Ion" but have different product IDs based on their contract.

## Dynamic Product Resolution

Here's a helper to automatically find the right product ID:

```go
func FindProductID(ctx context.Context, papiClient papi.PAPI, contractID, productName string) (string, error) {
    products, err := papiClient.GetProducts(ctx, papi.GetProductsRequest{
        ContractID: contractID,
    })
    if err != nil {
        return "", err
    }
    
    // Try exact match first
    for _, product := range products.Products.Items {
        if product.ProductName == productName {
            return product.ProductID, nil
        }
    }
    
    // Try case-insensitive partial match
    productNameLower := strings.ToLower(productName)
    for _, product := range products.Products.Items {
        if strings.Contains(strings.ToLower(product.ProductName), productNameLower) {
            return product.ProductID, nil
        }
    }
    
    return "", fmt.Errorf("product %s not found in contract %s", productName, contractID)
}

// Usage:
productID, err := FindProductID(ctx, papiClient, "ctr_C-1234567", "Ion")
config.ProductID = productID
```

## Testing Product IDs

```go
func TestProductID(t *testing.T, papiClient papi.PAPI, contractID, productID string) {
    // Try to create a test property with the product ID
    _, err := papiClient.CreateProperty(context.Background(), papi.CreatePropertyRequest{
        ContractID: contractID,
        GroupID:    groupID,
        Property: papi.PropertyCreate{
            ProductID:    productID,
            PropertyName: "test-property-" + time.Now().Format("20060102150405"),
        },
    })
    
    if err != nil {
        t.Errorf("Product ID %s is not valid: %v", productID, err)
    }
}
```

## Summary

### Key Takeaways

1. **No SDK Constants**: The SDK doesn't define product ID constants
2. **Contract-Specific**: Product IDs vary by contract
3. **Fetch Dynamically**: Use `GetProducts()` API
4. **Verify First**: Always verify product IDs before using
5. **Document Source**: Note where you got your product IDs

### Recommended Workflow

```
1. Get your contract ID
2. Call GetProducts(contractID)
3. Review available products
4. Select appropriate product ID
5. Store in config/env variable
6. Use in property creation
```

### For Property Helpers

The product IDs in `property_helpers.go` are **examples based on common naming patterns**. 

**Before using in production:**
1. Run `GetProducts()` for your contract
2. Verify the actual product IDs
3. Update the helper functions or pass ProductID in config
4. Test with your contract

## Additional Resources

- [Akamai Property Manager API - Products](https://techdocs.akamai.com/property-mgr/reference/get-products)
- [Akamai Terraform Provider - Products](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/property)
- [Akamai Control Center](https://control.akamai.com)

## Questions?

If unsure about product IDs:
1. Contact your Akamai account team
2. Check Akamai Control Center
3. Use the GetProducts API endpoint
4. Review your contract documentation
