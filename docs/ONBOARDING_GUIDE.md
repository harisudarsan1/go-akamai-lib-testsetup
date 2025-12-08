# Akamai Property Onboarding Guide

A comprehensive, step-by-step tutorial for onboarding a domain to Akamai's CDN platform using the EdgeGrid Go SDK.

---

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [The Onboarding Flow](#the-onboarding-flow)
4. [Step 0: Understanding the Architecture](#step-0-understanding-the-architecture)
5. [Step 1: Authentication - Creating a Session](#step-1-authentication---creating-a-session)
6. [Step 2: SSL Certificate Enrollment (CPS)](#step-2-ssl-certificate-enrollment-cps)
7. [Step 3: Property Creation (PAPI)](#step-3-property-creation-papi)
8. [Step 4: Edge Hostname Creation](#step-4-edge-hostname-creation)
9. [Step 5: Hostname Mapping](#step-5-hostname-mapping)
10. [Step 6: Rule Tree & Origin Configuration](#step-6-rule-tree--origin-configuration)
11. [Step 7: Property Activation](#step-7-property-activation)
12. [Step 8: DNS Configuration](#step-8-dns-configuration)
13. [Step 9: Application Security (AppSec/WAF)](#step-9-application-security-appsecwaf)
14. [Complete Code Reference](#complete-code-reference)
15. [Troubleshooting](#troubleshooting)

---

## Overview

### What This Guide Covers

This guide walks you through the complete process of onboarding a domain to Akamai's Content Delivery Network (CDN). By the end, you'll understand:

- **How to authenticate** with Akamai APIs using EdgeGrid
- **How SSL certificates** are managed through CPS
- **How properties** define delivery configurations
- **How traffic flows** from users to your origin through Akamai
- **How to activate** configurations safely
- **How to secure** your domain with AppSec/WAF

### What is "Onboarding"?

**Onboarding** is the process of configuring Akamai to deliver content for your domain. It involves:

1. ✅ **Certificate Setup** - Get SSL/TLS certificates for HTTPS
2. ✅ **Property Configuration** - Define how content is delivered
3. ✅ **DNS Configuration** - Point your domain to Akamai
4. ✅ **Security Setup** - Protect against attacks

### Why Akamai?

Akamai provides:
- **Global CDN** - Servers in 135+ countries
- **Performance** - Faster page loads via edge caching
- **Security** - DDoS protection, WAF, bot management
- **Reliability** - 100% uptime SLA
- **Scale** - Handle traffic spikes automatically

### Time Estimate

- **Reading this guide**: 2-3 hours
- **First implementation**: 4-6 hours
- **Subsequent onboardings**: 1-2 hours

---

## Prerequisites

### Required Knowledge

- ✅ Basic understanding of DNS (A records, CNAME records)
- ✅ Basic understanding of HTTP/HTTPS
- ✅ Familiarity with Go programming language
- ✅ Command line comfort

### Required Access

| Requirement | How to Get It |
|-------------|---------------|
| **Akamai Account** | Contact Akamai sales team |
| **Control Center Access** | Request from account admin |
| **API Credentials** | Create in Control Center → Identity & Access |
| **Contract & Group IDs** | Find in Control Center or via API |
| **Product Provisioning** | Ensure products are added to contract |

### Required Information

Before starting, gather:

```yaml
Contract Information:
  - Contract ID: "ctr_C-1234567"
  - Group ID: "grp_12345"
  - Product ID: "prd_Ion" (or your product)

Domain Information:
  - User-Facing Domain: "www.example.com"
  - Origin Server: "origin.example.com" or "203.0.113.42"

DNS Access:
  - Ability to create CNAME records
  - Access to DNS provider (Route53, Cloudflare, etc.)

Notification:
  - Email for activation notifications
```

### Development Environment

```bash
# Go 1.21 or later
go version

# Git (for cloning repository)
git --version

# Text editor or IDE
# VS Code, GoLand, Vim, etc.
```

---

## The Onboarding Flow

### High-Level Steps

```
┌─────────────────────────────────────────────────────────────┐
│                  AKAMAI ONBOARDING FLOW                      │
└─────────────────────────────────────────────────────────────┘

Step 0: Prerequisites
├─ Gather credentials
├─ Identify contract/group
└─ Determine product

Step 1: Authentication ──────────────┐
├─ Load .edgerc credentials          │
├─ Create EdgeGrid signer            │  API Setup
└─ Initialize session                │
                                     ┘
Step 2: Certificate (CPS) ───────────┐
├─ List enrollments                  │
├─ Find/create for domain            │  Certificate
├─ Get enrollment ID                 │  Management
└─ Certificate deploys to edge       │
                                     ┘
Step 3: Property (PAPI) ─────────────┐
├─ List/search properties            │
├─ Create if not exists              │  Content
├─ Get property ID & version         │  Delivery
└─ Property configuration ready      │  Configuration
                                     ┘
Step 4: Edge Hostname ───────────────┤
├─ List edge hostnames               │
├─ Create with cert enrollment       │
├─ Links domain → certificate        │
└─ Returns *.edgekey.net address     │
                                     ┘
Step 5: Hostname Mapping ────────────┤
├─ Map user domain → edge hostname   │
├─ Configure certificate type        │
├─ Set up Host header forwarding     │
└─ Enable for property               │
                                     ┘
Step 6: Origin Configuration ────────┤
├─ Get current rule tree             │
├─ Add origin behavior               │
├─ Configure caching                 │
└─ Update rule tree                  │
                                     ┘
Step 7: Activation ──────────────────┤
├─ Activate to STAGING network       │
├─ Test configuration                │
├─ Activate to PRODUCTION            │
└─ Monitor deployment                │
                                     ┘
Step 8: DNS Configuration ───────────┤
├─ Create CNAME record               │
├─ Wait for propagation              │
└─ Verify resolution                 │
                                     ┘
Step 9: Security (AppSec) ───────────┤
├─ List security configurations      │  Optional
├─ Add hostname protection           │
├─ Configure WAF policies            │
└─ Activate security config          │
                                     ┘

Result: ✅ Domain serving via Akamai
```

### Dependencies Between Steps

```
Authentication (1) ─┬─> Certificate (2) ──> Edge Hostname (4) ──┐
                    │                                            │
                    └─> Property (3) ──────────────────────────>├─> Hostname Mapping (5)
                                                                 │
                                                                 ▼
                                                         Origin Config (6)
                                                                 │
                                                                 ▼
                                                         Activation (7)
                                                                 │
                                                                 ▼
                                                         DNS Setup (8)
                                                                 │
                                                                 ▼
                                                         AppSec (9) [Optional]
```

---

## Step 0: Understanding the Architecture

### What Happens When a User Visits Your Site?

```
1. User Browser
   │ User types: https://www.example.com
   │
   ▼
2. DNS Resolution
   ┌─────────────────────────────────┐
   │ User's DNS Resolver             │
   │ Queries: www.example.com        │
   └─────────────────────────────────┘
   │
   │ Returns: CNAME → www.example.com.edgekey.net
   │
   ▼
3. Akamai DNS
   ┌─────────────────────────────────┐
   │ Akamai's DNS System             │
   │ - Determines optimal edge server│
   │ - Based on:                     │
   │   • User location               │
   │   • Server load                 │
   │   • Network conditions          │
   └─────────────────────────────────┘
   │
   │ Returns: Edge Server IP (23.x.x.x)
   │
   ▼
4. Akamai Edge Server
   ┌─────────────────────────────────┐
   │ Edge Server                     │
   │ ┌─────────────────────────┐     │
   │ │ TLS Handshake           │     │
   │ │ - Presents certificate  │     │
   │ │ - Validates client      │     │
   │ └─────────────────────────┘     │
   │ ┌─────────────────────────┐     │
   │ │ Cache Check             │     │
   │ │ Is content cached?      │     │
   │ └─────────────────────────┘     │
   │          │                       │
   │     ┌────┴────┐                  │
   │     │  Cache  │                  │
   │   Hit ▼    Miss ▼                │
   │ Serve    Fetch from              │
   │ Cache    Origin                  │
   └─────────────────────────────────┘
   │                    │
   │                    ▼
   │            5. Origin Server
   │            ┌─────────────────────────────────┐
   │            │ origin.example.com              │
   │            │ - Receives request with:        │
   │            │   Host: www.example.com         │
   │            │ - Processes request             │
   │            │ - Returns response              │
   │            └─────────────────────────────────┘
   │                    │
   └────────◀───────────┘
   │
   │ Response (HTML, images, etc.)
   │
   ▼
6. User Browser
   └─ Renders page
```

### Key Components

| Component | Purpose | Example |
|-----------|---------|---------|
| **User Domain** | What users type | `www.example.com` |
| **Edge Hostname** | Akamai's CDN entry point | `www.example.com.edgekey.net` |
| **Edge Server** | Akamai's cache server | Closest to user |
| **Origin Server** | Your actual server | `origin.example.com` |
| **Property** | Delivery configuration | Rules, behaviors, caching |
| **Certificate** | SSL/TLS for HTTPS | Managed by CPS |
| **Security Config** | WAF/DDoS protection | AppSec configuration |

---

## Step 1: Authentication - Creating a Session

### What is Authentication?

Before you can make API calls to Akamai, you need to authenticate using **EdgeGrid**, Akamai's authentication protocol. EdgeGrid uses:

- **Client Token** - Identifies your API client
- **Client Secret** - Proves you own the client
- **Access Token** - Grants permissions to specific APIs
- **Host** - The API endpoint base URL

### Why EdgeGrid?

EdgeGrid is a custom authentication protocol designed for:
- ✅ **Security** - Uses HMAC-SHA256 signatures
- ✅ **Request Integrity** - Signs entire request (headers + body)
- ✅ **Replay Protection** - Uses timestamp and nonce
- ✅ **No session state** - Each request is independently authenticated

### The `.edgerc` File

Credentials are stored in an `.edgerc` file (similar to `.aws/credentials`):

```ini
# ~/.edgerc

[default]
host = akaa-xxxxxxxxx.luna.akamaiapis.net
client_token = akab-xxxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxxx
client_secret = xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
access_token = akab-xxxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxxx
```

**File Structure**:
- **Sections** - Multiple credential sets (like AWS profiles)
- **host** - API base URL (region-specific)
- **Credentials** - Token/secret pairs

### Creating API Credentials

#### Step-by-step in Control Center:

```
1. Login to Akamai Control Center
   └─ https://control.akamai.com

2. Navigate to Identity & Access Management
   ├─ Top right menu (profile icon)
   └─ Identity & Access → API User

3. Create New API Client
   ├─ Click "Create API Client"
   ├─ Name: "Property Onboarding Script"
   ├─ Set credential expiration (e.g., 2 years)
   └─ Grant API Permissions (see table below)

4. Download Credentials
   ├─ Copy generated credentials
   ├─ Save to ~/.edgerc
   └─ Set permissions: chmod 600 ~/.edgerc
```

#### Required API Permissions

| API Service | Permission Level | Used For |
|-------------|------------------|----------|
| **PAPI (Property Manager)** | READ-WRITE | Create properties, hostnames, rules |
| **CPS (Certificate)** | READ-WRITE | Manage SSL certificates |
| **Edge Hostnames** | READ-WRITE | Create edge hostnames |
| **AppSec** | READ-WRITE | Configure WAF/security |
| **Contract** | READ | Fetch contract/group info |

### Code Walkthrough

#### File: `main.go` Lines 100-134

```go
func newSession(edgercPath, section string) (session.Session, error) {
    // 1. Expand tilde (~) in path to full home directory
    expanded, err := expandTilde(edgercPath)
    if err != nil {
        return nil, fmt.Errorf("expand edgerc path: %w", err)
    }

    // 2. Create EdgeGrid signer from .edgerc file
    signer, err := edgegrid.New(
        edgegrid.WithFile(expanded),    // Path to .edgerc
        edgegrid.WithSection(section),  // Section name (e.g., "default")
    )
    if err != nil {
        return nil, fmt.Errorf("edgegrid.New: %w", err)
    }

    // 3. Create session with the signer
    sess, err := session.New(session.WithSigner(signer))
    if err != nil {
        return nil, fmt.Errorf("session.New: %w", err)
    }

    return sess, nil
}
```

**What Happens**:
1. **Path Expansion** - Convert `~/.edgerc` to `/home/user/.edgerc`
2. **Signer Creation** - Load credentials, create signing function
3. **Session Creation** - Create HTTP client with auto-signing

#### Using the Session

```go
// In main()
sess, err := newSession(EdgercPath, EdgercSection)
if err != nil {
    log.Fatalf("Session init failed: %v", err)
}

// Create API clients
cpsClient := cps.Client(sess)
papiClient := papi.Client(sess)
appsecClient := appsec.Client(sess)
```

**Each API client**:
- Uses the same session
- Auto-signs requests with EdgeGrid
- Handles retries and rate limiting

### Configuration Parameters

| Parameter | Type | Example | Description |
|-----------|------|---------|-------------|
| `EdgercPath` | string | `"~/.edgerc"` | Path to credentials file |
| `EdgercSection` | string | `"default"` | Section in .edgerc |
| `ContractID` | string | `"ctr_C-1234567"` | Your Akamai contract |
| `GroupID` | string | `"grp_12345"` | Property group |

**Finding Contract/Group IDs**:

```go
// Option 1: Via API
contracts, _ := papiClient.GetContracts(ctx)
for _, contract := range contracts.Contracts.Items {
    fmt.Printf("Contract: %s (ID: %s)\n", 
        contract.ContractTypeName, contract.ContractID)
}

// Option 2: In Control Center
// Property Manager → Any Property → Details Panel
```

### Common Issues

| Error | Cause | Solution |
|-------|-------|----------|
| `edgegrid.New: file not found` | .edgerc doesn't exist | Create file at ~/.edgerc |
| `401 Unauthorized` | Invalid credentials | Regenerate API credentials |
| `403 Forbidden` | Insufficient permissions | Grant required API access |
| `client_token is empty` | Missing credentials | Check .edgerc file format |
| `invalid section` | Section doesn't exist | Verify section name in .edgerc |

### Testing Authentication

```go
// Simple test: List properties
properties, err := papiClient.GetProperties(ctx, papi.GetPropertiesRequest{})
if err != nil {
    log.Fatal("Authentication failed:", err)
}
fmt.Printf("✅ Authenticated! Found %d properties\n", len(properties.Properties.Items))
```

### Security Best Practices

```bash
# 1. Restrict file permissions
chmod 600 ~/.edgerc

# 2. Never commit .edgerc to git
echo ".edgerc" >> .gitignore

# 3. Use separate credentials for different environments
[production]
host = akaa-prod.luna.akamaiapis.net
...

[staging]
host = akaa-staging.luna.akamaiapis.net
...

# 4. Rotate credentials regularly
# Set expiration when creating API client

# 5. Use minimal permissions
# Only grant access to APIs you need
```

### What's Next?

Now that you have an authenticated session, you can:
1. ✅ Call any Akamai API
2. ➡️ **Next Step**: Get or create SSL certificate (CPS)

---

## Step 2: SSL Certificate Enrollment (CPS)

### What is CPS?

**CPS (Certificate Provisioning System)** manages SSL/TLS certificates for your domains on Akamai's edge servers.

### Why SSL Certificates?

```
Without Certificate:
User → https://www.example.com
       └─ ❌ Browser Error: "Your connection is not private"

With Certificate:
User → https://www.example.com
       └─ ✅ Secure connection (padlock icon)
```

### Certificate Types

| Type | Validation | Issuance Time | Use Case |
|------|------------|---------------|----------|
| **DV** (Domain Validation) | Email/DNS | Minutes | Most common, fastest |
| **OV** (Organization) | Business verification | Days | Corporate sites |
| **EV** (Extended Validation) | Legal verification | Weeks | Banking, high-security |

### Certificate Authorities

Akamai partners with:
- **Let's Encrypt** - Free DV certificates (90-day validity)
- **DigiCert** - Commercial certificates (1-2 year validity)
- **Third-Party** - Bring your own certificate

### Enrollment Lifecycle

```
1. Create Enrollment
   ├─ Specify domains (CN + SANs)
   ├─ Choose validation type
   ├─ Select certificate authority
   └─ Provide org details (OV/EV only)

2. Domain Validation
   ├─ DNS: Add TXT record
   │   └─ _acme-challenge.example.com → "token"
   │
   └─ HTTP: Place file on origin
       └─ /.well-known/acme-challenge/token

3. Certificate Issuance
   ├─ CA validates domain ownership
   ├─ CA issues certificate
   └─ Akamai receives certificate

4. Deployment to Edge
   ├─ Certificate deployed to Akamai edge servers
   ├─ Propagation across global network
   └─ Status: DEPLOYED (ready to use)

5. Auto-Renewal (DV only)
   └─ Akamai auto-renews before expiration
```

### Code Walkthrough

#### File: `main.go` Lines 140-176

```go
func getOrCreateEnrollment(ctx context.Context, cpsClient cps.CPS, domain string) (int, error) {
    fmt.Println(">> Checking for existing SSL Enrollment for", domain)

    // 1. List all enrollments for your contract
    enrollmentsResp, err := cpsClient.ListEnrollments(ctx, cps.ListEnrollmentsRequest{
        ContractID: ContractID,
    })
    if err != nil {
        return 0, fmt.Errorf("failed to list enrollments: %w", err)
    }

    // 2. Check if enrollment exists for this domain
    for _, enrollment := range enrollmentsResp.Enrollments {
        // Check Common Name (CN)
        if enrollment.CSR != nil && enrollment.CSR.CN == domain {
            fmt.Printf(">> Found existing enrollment ID: %d\n", enrollment.ID)
            return enrollment.ID, nil
        }
        
        // Check Subject Alternative Names (SANs)
        if enrollment.CSR != nil {
            for _, san := range enrollment.CSR.SANS {
                if san == domain {
                    fmt.Printf(">> Found enrollment ID: %d (in SANs)\n", enrollment.ID)
                    return enrollment.ID, nil
                }
            }
        }
    }

    // 3. No enrollment found - need manual creation
    return 0, fmt.Errorf("no enrollment found - manual creation required")
}
```

**What Happens**:
1. **List Enrollments** - Fetch all certificates for your contract
2. **Search by Domain** - Check CN (primary domain) and SANs (additional domains)
3. **Return Enrollment ID** - Used to link edge hostname to certificate

### Enrollment Structure

```go
type Enrollment struct {
    ID              int                 // Unique enrollment ID
    Status          string              // DEPLOYED, PENDING, FAILED
    CSR             *CSR               // Certificate Signing Request
    NetworkConfig   *NetworkConfig     // Edge deployment settings
    ValidationType  string             // DV, OV, EV
    CertificateType string             // THIRD_PARTY, SAN
}

type CSR struct {
    CN    string   // Common Name: www.example.com
    SANS  []string // Subject Alternative Names: [api.example.com, cdn.example.com]
    Org   string   // Organization name (OV/EV)
    ...
}
```

### Creating Enrollment (Manual Process)

**Why Manual?** 
- Enrollment requires organization details, validation preferences
- Typically done once per domain/wildcard
- Reused across multiple properties

**Steps**:
1. **Control Center** → Certificate Provisioning → Enrollments
2. **Create Enrollment**:
   ```
   Certificate Type: SAN (multiple domains)
   Common Name (CN): www.example.com
   SANs: api.example.com, cdn.example.com
   Validation Type: DV (Domain Validation)
   Certificate Authority: Let's Encrypt
   Network: Enhanced TLS
   ```
3. **Domain Validation**:
   - Choose DNS or HTTP validation
   - Add validation token to DNS/origin
4. **Wait for Deployment**:
   - Status changes: PENDING → VALIDATED → DEPLOYED
   - Deployment takes 15-60 minutes

### Common Enrollment Configurations

#### Single Domain (DV)
```yaml
Type: SAN Certificate
CN: www.example.com
SANs: []
Validation: Domain Validation (DNS)
CA: Let's Encrypt
Network: Enhanced TLS
Auto-Renewal: Yes
```

#### Multiple Domains
```yaml
Type: SAN Certificate
CN: www.example.com
SANs: [api.example.com, cdn.example.com, images.example.com]
Validation: Domain Validation (DNS)
CA: DigiCert
Network: Enhanced TLS
Validity: 1 year
```

#### Wildcard Certificate
```yaml
Type: SAN Certificate
CN: *.example.com
SANs: [example.com]
Validation: Domain Validation (DNS required)
CA: DigiCert
Note: Wildcard only works with DNS validation
```

### Enrollment Status

| Status | Meaning | Action Needed |
|--------|---------|---------------|
| `PENDING` | Waiting for validation | Complete domain validation |
| `VALIDATED` | CA verified, deploying | Wait for deployment |
| `DEPLOYED` | ✅ Live on edge | Ready to use |
| `FAILED` | Validation/issuance failed | Check error, retry validation |
| `EXPIRING_SOON` | Expires in < 30 days | Renew certificate |

### Certificate Parameters

| Parameter | Options | Description |
|-----------|---------|-------------|
| `certificateType` | SAN, THIRD_PARTY | Akamai-managed or your own |
| `validationType` | DV, OV, EV | Validation level |
| `networkConfiguration` | ENHANCED_TLS, STANDARD_TLS | TLS version support |
| `signatureAlgorithm` | SHA-256, SHA-384 | Signing algorithm |
| `keyAlgorithm` | RSA-2048, RSA-4096, ECDSA | Key type |

### DNS Validation Example

When you create a DV enrollment with DNS validation:

```bash
# Akamai provides validation token
Token Name:  _acme-challenge.www.example.com
Token Value: xyzABC123...

# Add TXT record to your DNS
Type:  TXT
Name:  _acme-challenge.www
Value: xyzABC123...
TTL:   300

# Verify DNS propagation
dig TXT _acme-challenge.www.example.com

# Akamai automatically validates and issues certificate
```

### HTTP Validation Example

```bash
# Akamai provides challenge file
URL:  http://www.example.com/.well-known/acme-challenge/token123
Content: "xyzABC123..."

# Place file on your origin server
mkdir -p /var/www/.well-known/acme-challenge/
echo "xyzABC123..." > /var/www/.well-known/acme-challenge/token123

# Verify accessibility
curl http://www.example.com/.well-known/acme-challenge/token123

# Akamai validates and issues certificate
```

### Common Issues

| Error | Cause | Solution |
|-------|-------|----------|
| `No enrollment found` | Domain not enrolled | Create enrollment in Control Center |
| `Enrollment PENDING` | Validation incomplete | Complete DNS/HTTP validation |
| `Enrollment FAILED` | Validation failed | Check DNS records, retry |
| `Certificate expired` | Past expiration date | Renew or create new enrollment |
| `Domain mismatch` | Domain not in CN or SANs | Add domain to enrollment |

### Enrollment Reuse

**Best Practice**: One enrollment can be used for multiple properties.

```go
// Property 1: www.example.com (uses enrollment 12345)
// Property 2: api.example.com (uses same enrollment if api.example.com in SANs)

enrollmentID := 12345  // Shared across properties
```

### What's Next?

Now that you have a certificate enrollment:
1. ✅ SSL/TLS certificate is deployed to Akamai edge
2. ➡️ **Next Step**: Create a property configuration (PAPI)

---

## Step 3: Property Creation (PAPI)

### What is a Property?

A **property** is Akamai's configuration container that defines:
- ✅ **Delivery settings** - Caching, compression, optimization
- ✅ **Hostnames** - Which domains use this configuration
- ✅ **Origin** - Where to fetch content
- ✅ **Rules** - Conditional behaviors based on requests
- ✅ **Security** - Request/response modifications

**Think of it as**: A configuration file for how Akamai delivers your content.

### PAPI (Property Manager API)

**PAPI** is the API for managing properties. It provides:
- Property CRUD operations
- Version control (like Git for configs)
- Rule tree management
- Hostname assignments
- Activation to staging/production

### Property Hierarchy

```
Contract (ctr_C-1234567)
  └─ Group (grp_12345)
      ├─ Property: www-example-com
      │   ├─ Version 1 (production)
      │   ├─ Version 2 (staging)
      │   └─ Version 3 (latest, editable)
      │
      └─ Property: api-example-com
          └─ Version 1 (production)
```

### Property Versions

Properties use **version control**:

| Version Type | Description | Editable? |
|--------------|-------------|-----------|
| **Latest** | Newest version | ✅ Yes |
| **Staging** | Active on staging network | ❌ Read-only |
| **Production** | Active on production | ❌ Read-only |

**Workflow**:
```
1. Edit latest version (v3)
2. Activate v3 to staging
3. Test staging
4. Activate v3 to production
5. Create new version (v4) for next change
```

### Code Walkthrough

#### File: `main.go` Lines 182-236

```go
func getOrCreateProperty(ctx context.Context, papiClient papi.PAPI, name string) (*papi.Property, error) {
    fmt.Printf(">> Checking for existing property: %s\n", name)

    // 1. List all properties in your contract/group
    properties, err := papiClient.GetProperties(ctx, papi.GetPropertiesRequest{
        ContractID: ContractID,
        GroupID:    GroupID,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to list properties: %w", err)
    }

    // 2. Search for property by name
    for _, prop := range properties.Properties.Items {
        if prop.PropertyName == name {
            fmt.Printf(">> Found existing property: %s (ID: %s, Version: %d)\n",
                prop.PropertyName, prop.PropertyID, prop.LatestVersion)
            return prop, nil
        }
    }

    // 3. Property doesn't exist, create it
    fmt.Printf(">> Creating new property: %s\n", name)
    createResp, err := papiClient.CreateProperty(ctx, papi.CreatePropertyRequest{
        ContractID: ContractID,
        GroupID:    GroupID,
        Property: papi.PropertyCreate{
            ProductID:    ProductID,    // e.g., "prd_Ion"
            PropertyName: name,
            RuleFormat:   "latest",     // Use latest rule format
        },
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create property: %w", err)
    }

    propertyID := createResp.PropertyID
    fmt.Printf(">> Property created with ID: %s\n", propertyID)

    // 4. Fetch the newly created property details
    propResp, err := papiClient.GetProperty(ctx, papi.GetPropertyRequest{
        ContractID: ContractID,
        GroupID:    GroupID,
        PropertyID: propertyID,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to fetch created property: %w", err)
    }

    return propResp.Property, nil
}
```

**What Happens**:
1. **List Properties** - Get all properties in your group
2. **Search by Name** - Check if property already exists
3. **Create if Missing** - Initialize new property with product
4. **Fetch Details** - Get full property object with version info

### Property Structure

```go
type Property struct {
    PropertyID      string  // prp_12345
    PropertyName    string  // www-example-com
    ContractID      string  // ctr_C-1234567
    GroupID         string  // grp_12345
    ProductID       string  // prd_Ion
    LatestVersion   int     // 3 (editable)
    StagingVersion  int     // 2 (on staging)
    ProductionVersion int   // 1 (on production)
    RuleFormat      string  // v2023-10-30
}
```

### Products and Property Types

**Product** determines the feature set and capabilities.

| Product ID | Type | Use Case | Features |
|------------|------|----------|----------|
| `prd_Ion` | Ion Standard | General web delivery | Caching, compression, image optimization |
| `prd_Download_Delivery` | Download | Software, large files | Adaptive bitrate, resumable downloads |
| `prd_Dynamic_Site_Delivery` | DSD | Dynamic content, APIs | Low TTL, origin offload, real-time |
| `prd_Media_Delivery` | Media | Video, audio streaming | Adaptive streaming, DRM support |

**⚠️ Product IDs are contract-specific** - see `PRODUCT_IDS.md` for details.

### Finding Your Product ID

```go
// Option 1: Fetch from API
productsResp, _ := papiClient.GetProducts(ctx, papi.GetProductsRequest{
    ContractID: ContractID,
})
for _, product := range productsResp.Products.Items {
    fmt.Printf("Product: %s (ID: %s)\n", product.ProductName, product.ProductID)
}

// Option 2: Use helper (see product_utils.go)
mapper := NewProductMapper()
productID, err := mapper.FindProductID(ctx, papiClient, ContractID, "Ion Standard")
```

### Property Naming Conventions

**Best Practices**:
```go
// ✅ Good names
"www-example-com"        // Domain-based
"api-gateway-production" // Descriptive
"images-cdn-v2"          // Purpose-based

// ❌ Avoid
"property1"              // Not descriptive
"test"                   // Too generic
"www.example.com"        // Dots can cause issues
```

### Rule Format

**Rule Format** specifies the version of PAPI rule syntax:

| Format | Description | Use When |
|--------|-------------|----------|
| `"latest"` | Most recent stable format | ✅ New properties |
| `"v2023-10-30"` | Specific version | Compatibility needed |
| `"frozen"` | Locked format | Legacy properties |

**Recommendation**: Always use `"latest"` for new properties.

### Property Configuration Options

```go
Property: papi.PropertyCreate{
    ProductID:    "prd_Ion",        // Delivery product
    PropertyName: "my-cdn",          // Unique name
    RuleFormat:   "latest",          // Rule syntax version
    
    // Optional fields
    CloneFrom: &papi.PropertyClone{
        PropertyID:      "prp_11111",  // Clone existing property
        CloneFromVersion: 1,
    },
}
```

### Property Lifecycle

```
1. Create Property
   └─ PropertyID: prp_12345, Version: 1

2. Configure Property
   ├─ Add hostnames
   ├─ Set origin
   └─ Configure rules

3. Activate to Staging
   └─ Version 1 → Staging

4. Test Staging
   └─ curl -H "Host: www.example.com" staging-edge-server

5. Activate to Production
   └─ Version 1 → Production

6. Make Changes
   ├─ Create version 2 (copy of v1)
   ├─ Edit version 2
   └─ Repeat activation cycle
```

### Idempotent Operations

The code is **idempotent** - safe to run multiple times:

```go
// First run: Creates property
// Second run: Finds existing property

if property exists {
    return existing property
} else {
    create new property
}
```

### Common Issues

| Error | Cause | Solution |
|-------|-------|----------|
| `Property name already exists` | Duplicate name | Use different name or fetch existing |
| `Invalid product ID` | Product not in contract | Use `GetProducts()` to find valid IDs |
| `Contract not found` | Wrong contract ID | Verify contract ID |
| `Group not found` | Wrong group ID | Check group ID in Control Center |
| `Insufficient permissions` | Missing PAPI write access | Grant PAPI permissions to API client |

### What Properties Do NOT Include

Properties are **just configurations**. They don't handle:
- ❌ SSL certificates (that's CPS)
- ❌ Edge hostname creation (separate API)
- ❌ DNS configuration (you manage that)
- ❌ Origin server (you manage that)

### What's Next?

Now that you have a property:
1. ✅ Property configuration container created
2. ✅ Version 1 exists (editable)
3. ➡️ **Next Step**: Create edge hostname (links certificate to CDN entry point)

---

## Step 4: Edge Hostname Creation

### What is an Edge Hostname?

An **edge hostname** is Akamai's entry point for your domain. It's the CNAME target that connects your user-facing domain to Akamai's CDN.

```
User Domain          Edge Hostname           Akamai Network
www.example.com  →  www.example.com.edgekey.net  →  Edge Servers
```

### Why Edge Hostnames?

```
Without Edge Hostname:
www.example.com → [No Akamai]
                  └─ Direct to origin (slow, no protection)

With Edge Hostname:
www.example.com → www.example.com.edgekey.net → Akamai Edge → Origin
                  └─ Cached, optimized, protected
```

### Edge Hostname Components

```
www.example.com.edgekey.net
├─────────────┬────────────┬──
│             │            │
Domain Prefix │            Domain Suffix
              │
              Edge Network Identifier
```

| Component | Value | Description |
|-----------|-------|-------------|
| **Domain Prefix** | `www.example.com` | Your domain (without TLD) |
| **Domain Suffix** | `.edgekey.net` | Akamai's network suffix |
| **Full Edge Hostname** | `www.example.com.edgekey.net` | CNAME target |

### Domain Suffixes

| Suffix | Network | SSL Support | Use Case |
|--------|---------|-------------|----------|
| `.edgekey.net` | Standard TLS | ✅ HTTPS | Most common, shared certificate |
| `.edgesuite.net` | Legacy | ❌ HTTP only | Deprecated |
| `.akamaized.net` | Enhanced TLS | ✅ HTTPS | Modern, better performance |

**Recommendation**: Use `.edgekey.net` (default for most properties).

### Secure Network Types

| Network | Description | TLS Versions | Compatibility |
|---------|-------------|--------------|---------------|
| **STANDARD_TLS** | Legacy | TLS 1.0+ | Older browsers |
| **ENHANCED_TLS** | Modern | TLS 1.2+ | ✅ Recommended |
| **SHARED_CERT** | Shared certificate | TLS 1.2+ | Multi-tenant |

### IP Version Behavior

| Option | Description | Use When |
|--------|-------------|----------|
| **IPV4** | IPv4 only | Legacy systems |
| **IPV6_COMPLIANCE** | IPv4 + IPv6 | ✅ Modern web (recommended) |
| **IPV6_PERFORMANCE** | IPv6 preferred | IPv6-heavy traffic |

### Code Walkthrough

#### File: `main.go` Lines 238-288

```go
func ensureEdgeHostname(ctx context.Context, papiClient papi.PAPI, 
                        certEnrollmentID int, domain string) (string, error) {
    fmt.Printf(">> Checking for existing edge hostname for domain: %s\n", domain)

    // 1. List existing edge hostnames
    edgeHostnamesResp, err := papiClient.GetEdgeHostnames(ctx, papi.GetEdgeHostnamesRequest{
        ContractID: ContractID,
        GroupID:    GroupID,
    })
    if err != nil {
        return "", fmt.Errorf("failed to list edge hostnames: %w", err)
    }

    // 2. Extract domain prefix (first part before first dot)
    domainPrefix := domain
    if idx := strings.Index(domain, "."); idx > 0 {
        domainPrefix = domain[:idx]  // "www.example.com" → "www"
    }

    // 3. Check if edge hostname already exists
    expectedEdgeHostname := domainPrefix + ".edgekey.net"
    for _, eh := range edgeHostnamesResp.EdgeHostnames.Items {
        if eh.Domain == expectedEdgeHostname || eh.DomainPrefix == domainPrefix {
            fmt.Printf(">> Found existing edge hostname: %s (ID: %s)\n", 
                eh.Domain, eh.ID)
            return eh.Domain, nil
        }
    }

    // 4. Create new edge hostname
    fmt.Printf(">> Creating edge hostname with prefix: %s\n", domainPrefix)
    createResp, err := papiClient.CreateEdgeHostname(ctx, papi.CreateEdgeHostnameRequest{
        ContractID: ContractID,
        GroupID:    GroupID,
        EdgeHostname: papi.EdgeHostnameCreate{
            ProductID:         ProductID,           // Same as property
            DomainPrefix:      domainPrefix,        // "www"
            DomainSuffix:      "edgekey.net",       // Network type
            SecureNetwork:     papi.EHSecureNetworkEnhancedTLS,
            IPVersionBehavior: papi.EHIPVersionV4,  // IPv4 + IPv6
            CertEnrollmentID:  certEnrollmentID,    // Links certificate
        },
    })
    if err != nil {
        return "", fmt.Errorf("failed to create edge hostname: %w", err)
    }

    edgeHostname := domainPrefix + ".edgekey.net"
    fmt.Printf(">> Edge hostname created: %s (ID: %s)\n", 
        edgeHostname, createResp.EdgeHostnameID)

    return edgeHostname, nil
}
```

**What Happens**:
1. **List Edge Hostnames** - Check if one already exists for this domain
2. **Extract Prefix** - Get domain prefix (e.g., "www" from "www.example.com")
3. **Search Existing** - Look for matching edge hostname
4. **Create if Missing** - Initialize edge hostname with certificate link
5. **Return Edge Hostname** - Full `.edgekey.net` address

### Edge Hostname Structure

```go
type EdgeHostname struct {
    EdgeHostnameID    string  // ehn_12345
    Domain            string  // www.example.com.edgekey.net
    DomainPrefix      string  // www.example.com
    DomainSuffix      string  // edgekey.net
    SecureNetwork     string  // ENHANCED_TLS
    IPVersionBehavior string  // IPV4
    CertEnrollmentID  int     // Links to CPS enrollment
    ProductID         string  // prd_Ion
    Status            string  // CREATED, PENDING, ACTIVE
}
```

### Certificate Linking

**Critical**: Edge hostname MUST be linked to a certificate enrollment.

```go
EdgeHostname: papi.EdgeHostnameCreate{
    CertEnrollmentID: certEnrollmentID,  // From Step 2 (CPS)
    ...
}
```

**Why?**
- Edge hostname serves HTTPS traffic
- Needs certificate to present during TLS handshake
- Without certificate: ❌ "Certificate error" for users

### Edge Hostname Deployment

```
1. Create Edge Hostname Request
   └─ API returns immediately (async operation)

2. Edge Hostname Provisioning
   ├─ Certificate linked
   ├─ Deployed to edge servers globally
   └─ Status: CREATED → PENDING → ACTIVE

3. Ready to Use
   └─ Can now map user domains to this edge hostname
```

**Deployment Time**: 5-15 minutes (async, check status with `GetEdgeHostname`).

### Domain Prefix Extraction

```go
// Examples of domain → prefix mapping
"www.example.com"      → "www"
"api.example.com"      → "api"
"cdn.example.com"      → "cdn"
"example.com"          → "example"
"a.b.c.example.com"    → "a"  // Only first segment
```

**Edge Hostname Result**:
```
"www" → "www.edgekey.net"
"api" → "api.edgekey.net"
```

### Multiple Edge Hostnames

**Can I create multiple edge hostnames?**
- ✅ Yes, one per domain prefix
- ✅ All can share the same certificate (if domains in SANs)
- ✅ All can belong to same property

**Example**:
```go
// Enrollment: CN=www.example.com, SANs=[api.example.com, cdn.example.com]
// Edge Hostnames:
eh1 := "www.edgekey.net"      // For www.example.com
eh2 := "api.edgekey.net"      // For api.example.com
eh3 := "cdn.edgekey.net"      // For cdn.example.com
```

### Configuration Parameters

| Parameter | Type | Example | Description |
|-----------|------|---------|-------------|
| `ProductID` | string | `"prd_Ion"` | Same product as property |
| `DomainPrefix` | string | `"www"` | First part of domain |
| `DomainSuffix` | string | `"edgekey.net"` | Akamai network suffix |
| `SecureNetwork` | enum | `ENHANCED_TLS` | TLS configuration |
| `IPVersionBehavior` | enum | `IPV4` | IPv4/IPv6 support |
| `CertEnrollmentID` | int | `12345` | Certificate to use |

### Edge Hostname Reuse

**Idempotent**: Safe to run multiple times.

```go
// First run: Creates edge hostname
// Second run: Finds existing edge hostname

if edge hostname exists {
    return existing edge hostname
} else {
    create new edge hostname
}
```

### Common Issues

| Error | Cause | Solution |
|-------|-------|----------|
| `Edge hostname already exists` | Prefix collision | Different domain prefix |
| `Invalid certificate enrollment` | Wrong enrollment ID | Verify enrollment ID from Step 2 |
| `Certificate not deployed` | Enrollment status != DEPLOYED | Wait for certificate deployment |
| `Domain prefix in use` | Prefix used by another property | Use different prefix or share edge hostname |
| `Invalid product` | Product mismatch | Use same product as property |

### Edge Hostname Status

| Status | Meaning | Action |
|--------|---------|--------|
| `CREATED` | API request accepted | Wait for provisioning |
| `PENDING` | Deploying to edge | Wait 5-15 minutes |
| `ACTIVE` | ✅ Ready to use | Can map user domains |
| `FAILED` | Provisioning error | Check error, recreate |

### Testing Edge Hostname

```bash
# Check DNS resolution (won't work until DNS configured in Step 8)
dig www.example.com.edgekey.net

# Should return Akamai edge server IPs (23.x.x.x range)
```

### What's Next?

Now that you have an edge hostname:
1. ✅ Edge hostname created (e.g., www.example.com.edgekey.net)
2. ✅ Linked to SSL certificate
3. ✅ Deployed to Akamai edge servers
4. ➡️ **Next Step**: Map user domain to edge hostname (hostname mapping)

---

## Step 5: Hostname Mapping

### What is Hostname Mapping?

**Hostname mapping** connects your user-facing domain to the edge hostname within a property configuration.

```
Property Configuration:
├─ Hostname Mapping
│  ├─ CnameFrom: www.example.com       (what users type)
│  └─ CnameTo: www.example.com.edgekey.net  (where it goes)
├─ Rules
└─ Behaviors
```

### Why Hostname Mapping?

**Without Mapping**:
- Edge hostname exists but isn't associated with property
- Requests to edge hostname have no configuration
- No caching, no origin, no behaviors

**With Mapping**:
- Property knows which domains it handles
- Requests to mapped hostnames use property rules
- Caching, origin, behaviors all applied

### Hostname Structure

```go
type Hostname struct {
    CnameType            string   // "EDGE_HOSTNAME"
    EdgeHostnameID       string   // "ehn_12345"
    CnameFrom            string   // "www.example.com"
    CnameTo              string   // "www.example.com.edgekey.net"
    CertProvisioningType string   // "CPS_MANAGED"
    CertStatus           string   // "DEPLOYED"
}
```

### Hostname Fields Explained

| Field | Type | Example | Description |
|-------|------|---------|-------------|
| **CnameType** | enum | `EDGE_HOSTNAME` | Type of target (always edge hostname) |
| **EdgeHostnameID** | string | `ehn_12345` | Reference to edge hostname |
| **CnameFrom** | string | `www.example.com` | User-facing domain |
| **CnameTo** | string | `www.example.com.edgekey.net` | Edge hostname target |
| **CertProvisioningType** | enum | `CPS_MANAGED` | Certificate type |
| **CertStatus** | string | `DEPLOYED` | Certificate readiness |
| **CCMCertStatus** | string | `DEPLOYED` | CCM certificate status |
| **CCMCertificates** | []string | `[12345]` | CCM certificate IDs |
| **MTLS** | object | `{...}` | Mutual TLS settings |
| **TLSConfiguration** | object | `{...}` | Custom TLS config |

### Certificate Provisioning Types

| Type | Description | Use When |
|------|-------------|----------|
| **CPS_MANAGED** | ✅ Akamai-managed certificate | Standard (recommended) |
| **DEFAULT** | Shared Akamai certificate | Quick setup, no custom cert |
| **CCM** | Custom Certificate Manager | Bring your own certificate |

### Code Walkthrough

#### File: `main.go` Lines 290-340

```go
func updatePropertyRules(ctx context.Context, papiClient papi.PAPI, 
                         prop *papi.Property, domain, edgeHostname string) error {
    fmt.Println(">> Updating property rules and hostnames")

    // PART 1: Update Hostnames
    
    // 1. Get existing hostnames for the property version
    existingHostnames, err := papiClient.GetPropertyVersionHostnames(ctx, 
        papi.GetPropertyVersionHostnamesRequest{
            PropertyID:      prop.PropertyID,
            PropertyVersion: prop.LatestVersion,
            ContractID:      ContractID,
            GroupID:         GroupID,
        })
    if err != nil {
        return fmt.Errorf("failed to get existing hostnames: %w", err)
    }

    // 2. Check if hostname already exists
    hostnameExists := false
    for _, h := range existingHostnames.Hostnames.Items {
        if h.CnameFrom == domain {
            hostnameExists = true
            fmt.Printf(">> Hostname %s already configured\n", domain)
            break
        }
    }

    // 3. Add new hostname if it doesn't exist
    if !hostnameExists {
        newHostnames := append(existingHostnames.Hostnames.Items, papi.Hostname{
            CnameType:            papi.HostnameCnameTypeEdgeHostname,
            CnameFrom:            domain,                // www.example.com
            CnameTo:              edgeHostname,         // www.example.com.edgekey.net
            CertProvisioningType: "CPS_MANAGED",        // Use CPS certificate
        })

        // 4. Update property with new hostname list
        _, err = papiClient.UpdatePropertyVersionHostnames(ctx, 
            papi.UpdatePropertyVersionHostnamesRequest{
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

    // PART 2: Update Origin (see Step 6)
    ...
}
```

**What Happens**:
1. **Get Existing Hostnames** - Fetch current hostname mappings for property version
2. **Check if Exists** - Avoid duplicate hostname entries
3. **Append New Hostname** - Add new mapping to list
4. **Update Property** - Save hostname list to property version

### Hostname Mapping Flow

```
User Request: https://www.example.com/page.html

1. DNS Resolution
   └─ www.example.com → www.example.com.edgekey.net → Edge Server IP

2. Request Arrives at Edge Server
   └─ Host: www.example.com

3. Hostname Lookup
   ├─ Check property configurations
   └─ Find property with CnameFrom = "www.example.com"

4. Apply Property Rules
   ├─ Use rule tree from matched property
   ├─ Apply caching behaviors
   ├─ Check origin configuration
   └─ Process request

5. Return Response
   └─ Cached content or fetch from origin
```

### Multiple Hostnames per Property

**Single property can handle multiple domains**:

```go
// Property: www-api-cdn

Hostnames:
1. CnameFrom: www.example.com   → CnameTo: www.example.com.edgekey.net
2. CnameFrom: api.example.com   → CnameTo: api.example.com.edgekey.net
3. CnameFrom: cdn.example.com   → CnameTo: cdn.example.com.edgekey.net

// All use same rules, origin, behaviors
// Can differentiate with conditional rules (match on hostname)
```

### Host Header Forwarding

When Akamai fetches from origin, which `Host` header is sent?

| Option | Value Sent to Origin | Use When |
|--------|---------------------|----------|
| **REQUEST_HOST_HEADER** | Original user domain (www.example.com) | ✅ Most common |
| **ORIGIN_HOSTNAME** | Origin hostname (origin.example.com) | Origin expects specific host |
| **CUSTOM** | Custom value | Special origin requirements |

**Configured in origin behavior (Step 6)**, not in hostname mapping.

### Certificate Status

Check certificate status in hostname configuration:

| CertStatus | Meaning | Action |
|------------|---------|--------|
| `DEPLOYED` | ✅ Certificate ready | Can activate property |
| `PENDING` | Deploying certificate | Wait for deployment |
| `FAILED` | Certificate error | Check CPS enrollment |
| `NOT_DEPLOYED` | No certificate | Link enrollment to edge hostname |

### Hostname Validation

```go
// Valid hostnames
"www.example.com"       ✅
"api.example.com"       ✅
"cdn.example.com"       ✅
"example.com"           ✅

// Invalid hostnames
"www.example.com."      ❌ (trailing dot)
"http://example.com"    ❌ (protocol)
"example.com/path"      ❌ (path)
"*.example.com"         ❌ (wildcard - use in certificate, not hostname)
```

### Idempotent Hostname Updates

```go
// First run: Adds hostname
// Second run: Detects existing hostname, skips

if hostname exists in property {
    skip addition
} else {
    append to hostname list
    update property
}
```

### Hostname Mapping vs DNS

**Important**: These are separate steps.

```
Step 5 (Now): Hostname Mapping in Property
└─ Tell Akamai: "www.example.com uses this edge hostname"

Step 8 (Later): DNS Configuration
└─ Tell Internet: "www.example.com points to edge hostname"
```

### Common Issues

| Error | Cause | Solution |
|-------|-------|----------|
| `Hostname already in use` | Another property uses this domain | Remove from other property first |
| `Edge hostname not found` | Invalid CnameTo value | Verify edge hostname exists |
| `Certificate mismatch` | Domain not in certificate SANs | Add domain to CPS enrollment |
| `Invalid CnameFrom` | Malformed domain name | Check domain format |
| `Version not editable` | Trying to edit activated version | Create new property version |

### Testing Hostname Mapping

```bash
# After mapping (before DNS)
# Edge hostname works, but user domain doesn't yet

curl -H "Host: www.example.com" https://www.example.com.edgekey.net/
# Should work if activated

# After DNS (Step 8)
curl https://www.example.com/
# Should work
```

### What's Next?

Now that hostnames are mapped:
1. ✅ User domain linked to edge hostname in property
2. ✅ Certificate provisioning type set
3. ✅ Property knows which domains it handles
4. ➡️ **Next Step**: Configure origin server (where to fetch content)

---

## Step 6: Rule Tree & Origin Configuration

### What is a Rule Tree?

A **rule tree** is the core configuration of a property. It defines:
- ✅ **Origin** - Where to fetch content
- ✅ **Caching** - What to cache and for how long
- ✅ **Behaviors** - How to process requests/responses
- ✅ **Conditions** - When to apply rules

**Think of it as**: Apache/Nginx configuration, but managed via API.

### Rule Tree Structure

```
Root Rule (Default Rule)
├─ Behaviors (applied to all requests)
│  ├─ Origin
│  ├─ Caching
│  ├─ Compression
│  └─ CP Code
│
└─ Child Rules (conditional)
   ├─ Rule: "Static Content"
   │  ├─ Match: File extension = .jpg, .png, .css, .js
   │  └─ Behavior: Cache for 7 days
   │
   ├─ Rule: "API Requests"
   │  ├─ Match: Path = /api/*
   │  └─ Behavior: No caching, forward cookies
   │
   └─ Rule: "Redirect HTTP to HTTPS"
      ├─ Match: Protocol = HTTP
      └─ Behavior: Redirect to HTTPS
```

### Origin Behavior

The **origin** behavior tells Akamai where your actual server is.

```go
type OriginBehavior struct {
    Name: "origin",
    Options: {
        originType:         "CUSTOMER",              // Your origin
        hostname:           "origin.example.com",    // Origin server
        forwardHostHeader:  "REQUEST_HOST_HEADER",   // Host header value
        cacheKeyHostname:   "REQUEST_HOST_HEADER",   // Cache key component
        compress:           true,                     // Compress responses
        enableTrueClientIp: false,                   // Add client IP header
        verificationMode:   "PLATFORM_SETTINGS",     // SSL verification
        httpPort:           80,                      // HTTP port
        httpsPort:          443,                     // HTTPS port
    }
}
```

### Origin Types

| Type | Description | When to Use |
|------|-------------|-------------|
| **CUSTOMER** | ✅ Your own server | Most common |
| **NET_STORAGE** | Akamai's storage | Files stored on Akamai |
| **MEDIA_SERVICE_LIVE** | Live video streaming | Live content |

### Origin Hostname Options

**What can you use as origin hostname?**

| Option | Example | Use Case |
|--------|---------|----------|
| **Domain name** | `origin.example.com` | ✅ Standard setup |
| **IP address** | `203.0.113.42` | Direct IP access |
| **Load balancer** | `lb.example.com` | Multiple backend servers |
| **Cloud service** | `my-app.us-east-1.elb.amazonaws.com` | AWS/GCP/Azure |

**Recommendation**: Use domain name for flexibility (easier to change IP).

### Host Header Forwarding

When Akamai requests from your origin, which `Host` header does it send?

| Option | Value Sent | Use When |
|--------|------------|----------|
| **REQUEST_HOST_HEADER** | ✅ User's domain (www.example.com) | Origin serves multiple domains |
| **ORIGIN_HOSTNAME** | Origin hostname (origin.example.com) | Origin expects specific host |
| **CUSTOM** | Custom value | Special requirements |

**Example**:
```
User requests: https://www.example.com/page.html
Origin receives: GET /page.html HTTP/1.1
                 Host: www.example.com  (if REQUEST_HOST_HEADER)
                 Host: origin.example.com  (if ORIGIN_HOSTNAME)
```

### Cache Key Hostname

Determines which hostname is used in the cache key.

| Option | Cache Key Uses | Effect |
|--------|----------------|--------|
| **REQUEST_HOST_HEADER** | User's domain | ✅ Standard (separate cache per domain) |
| **ORIGIN_HOSTNAME** | Origin hostname | Shared cache across domains |

**Example**:
```
With REQUEST_HOST_HEADER:
- www.example.com/page.html → Cache key: www.example.com/page.html
- api.example.com/page.html → Cache key: api.example.com/page.html
(Separate cache entries)

With ORIGIN_HOSTNAME:
- www.example.com/page.html → Cache key: origin.example.com/page.html
- api.example.com/page.html → Cache key: origin.example.com/page.html
(Same cache entry)
```

### Code Walkthrough

#### File: `main.go` Lines 341-399

```go
func updatePropertyRules(ctx context.Context, papiClient papi.PAPI,
                         prop *papi.Property, domain, edgeHostname string) error {
    
    // PART 1: Hostname Mapping (covered in Step 5)
    ...

    // PART 2: Origin Configuration

    fmt.Println(">> Updating rule tree with origin configuration")

    // 1. Get current rule tree
    ruleTree, err := papiClient.GetRuleTree(ctx, papi.GetRuleTreeRequest{
        PropertyID:      prop.PropertyID,
        PropertyVersion: prop.LatestVersion,
        ContractID:      ContractID,
        GroupID:         GroupID,
    })
    if err != nil {
        return fmt.Errorf("failed to get rule tree: %w", err)
    }

    // 2. Check if origin behavior already exists
    originHostname := "server:IP"  // Replace with your origin
    originExists := false
    for _, behavior := range ruleTree.Rules.Behaviors {
        if behavior.Name == "origin" {
            originExists = true
            fmt.Println(">> Origin behavior already configured")
            break
        }
    }

    // 3. Add origin behavior if not present
    if !originExists {
        originBehavior := papi.RuleBehavior{
            Name: "origin",
            Options: papi.RuleOptionsMap{
                "originType":         "CUSTOMER",              // Your server
                "hostname":           originHostname,          // Origin address
                "forwardHostHeader":  "REQUEST_HOST_HEADER",  // Send user's Host
                "cacheKeyHostname":   "REQUEST_HOST_HEADER",  // Cache by user's Host
                "compress":           true,                    // Enable compression
                "enableTrueClientIp": false,                  // No client IP header
                "verificationMode":   "PLATFORM_SETTINGS",    // Default SSL verify
                "httpPort":           80,                      // HTTP port
                "httpsPort":          443,                     // HTTPS port
            },
        }
        
        // Append to behaviors list
        ruleTree.Rules.Behaviors = append(ruleTree.Rules.Behaviors, originBehavior)

        // 4. Update the rule tree
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
```

**What Happens**:
1. **Get Rule Tree** - Fetch current rule configuration
2. **Check Origin** - See if origin behavior already exists
3. **Add Origin Behavior** - Configure where to fetch content
4. **Update Rule Tree** - Save changes to property version

### Origin Configuration Parameters

| Parameter | Type | Options | Description |
|-----------|------|---------|-------------|
| `originType` | string | CUSTOMER, NET_STORAGE | Origin location |
| `hostname` | string | domain or IP | Origin server address |
| `forwardHostHeader` | string | REQUEST_HOST_HEADER, ORIGIN_HOSTNAME, CUSTOM | Host header to send |
| `cacheKeyHostname` | string | REQUEST_HOST_HEADER, ORIGIN_HOSTNAME | Host for cache key |
| `compress` | bool | true, false | Enable gzip compression |
| `enableTrueClientIp` | bool | true, false | Add True-Client-IP header |
| `verificationMode` | string | PLATFORM_SETTINGS, CUSTOM | SSL verification level |
| `httpPort` | int | 1-65535 | HTTP port (default: 80) |
| `httpsPort` | int | 1-65535 | HTTPS port (default: 443) |

### Origin SSL Verification

Controls how Akamai validates your origin's SSL certificate.

| Mode | Description | Use When |
|------|-------------|----------|
| **PLATFORM_SETTINGS** | ✅ Standard validation | Origin has valid cert |
| **CUSTOM** | Custom cert/CA | Self-signed cert |
| **THIRD_PARTY** | Specific CA | Non-standard CA |

### True Client IP

**What**: Adds the end user's IP address to origin requests.

```
Without True-Client-IP:
Origin sees: X-Forwarded-For: 23.x.x.x (Akamai edge server)

With True-Client-IP:
Origin sees: True-Client-IP: 198.51.100.50 (actual user)
```

**When to enable**:
- ✅ Need user IP for rate limiting, geolocation, logging
- ❌ Disable if not needed (reduces request size)

### Origin Ports

**Default ports** (80/443) work for most scenarios.

**Custom ports**:
```go
"httpPort":  8080,  // Origin HTTP on port 8080
"httpsPort": 8443,  // Origin HTTPS on port 8443
```

**Use cases**:
- Development/staging servers on non-standard ports
- Backend services with specific port requirements
- Containerized applications

### Origin Connection

```
Edge Server → Origin Request Flow:

1. DNS Resolution
   └─ origin.example.com → 203.0.113.42

2. Protocol Selection
   ├─ User used HTTPS? → Connect to origin:443
   └─ User used HTTP?  → Connect to origin:80

3. TLS Handshake (if HTTPS)
   ├─ Verify origin certificate
   ├─ Establish encrypted connection
   └─ Continue if valid

4. HTTP Request
   GET /page.html HTTP/1.1
   Host: www.example.com (or origin.example.com)
   X-Forwarded-For: 203.0.113.100
   Via: 1.1 akamai (ghost)

5. Origin Response
   HTTP/1.1 200 OK
   Content-Type: text/html
   Cache-Control: max-age=3600
   ...

6. Cache & Serve
   ├─ Cache response on edge
   └─ Serve to user
```

### Caching Behavior

**Default Caching** (if no explicit cache behavior):

```
Akamai respects origin's Cache-Control headers:

Origin Response:
Cache-Control: max-age=3600  →  Akamai caches for 1 hour
Cache-Control: no-cache      →  Akamai doesn't cache
Cache-Control: public        →  Akamai can cache
```

**Override Caching** (add caching behavior to rule tree):

```go
cachingBehavior := papi.RuleBehavior{
    Name: "caching",
    Options: papi.RuleOptionsMap{
        "behavior":     "MAX_AGE",
        "mustRevalidate": false,
        "ttl":          "7d",  // Cache for 7 days regardless of origin headers
    },
}
```

### Advanced Origin Behaviors

#### Load Balancing (Multiple Origins)

```go
// Use origin behavior with failover
{
    "originType": "CUSTOMER",
    "hostname":   "origin1.example.com",
    "failover": {
        "hostname": "origin2.example.com",  // Backup origin
        "port":     443,
    },
}
```

#### Path-Based Origin

```go
// Child rule: Route /api to different origin
{
    "name": "API Routing",
    "criteria": [
        {
            "name": "path",
            "options": {
                "matchOperator": "MATCHES_ONE_OF",
                "values": ["/api/*"],
            },
        },
    ],
    "behaviors": [
        {
            "name": "origin",
            "options": {
                "hostname": "api-backend.example.com",
            },
        },
    ],
}
```

### Rule Tree Best Practices

```go
// ✅ Good: Explicit, maintainable
Origin: "origin.example.com"
Host Header: REQUEST_HOST_HEADER
Cache Key: REQUEST_HOST_HEADER

// ❌ Avoid: Using IP directly (hard to change)
Origin: "203.0.113.42"

// ❌ Avoid: Mixing host header settings without reason
Host Header: ORIGIN_HOSTNAME
Cache Key: REQUEST_HOST_HEADER
```

### Common Issues

| Error | Cause | Solution |
|-------|-------|----------|
| `Origin not reachable` | Invalid hostname/IP | Verify origin is accessible |
| `Origin SSL handshake failed` | Certificate error | Check origin cert, adjust verificationMode |
| `5xx errors from origin` | Origin server error | Check origin logs |
| `Caching too aggressive` | Wrong TTL settings | Adjust caching behavior |
| `Host header mismatch` | Wrong forwardHostHeader | Set to REQUEST_HOST_HEADER |

### Testing Origin Configuration

```bash
# Test origin directly (before Akamai)
curl -H "Host: www.example.com" http://origin.example.com/

# Test through Akamai edge (after activation)
curl -H "Host: www.example.com" https://www.example.com.edgekey.net/

# Test through user domain (after DNS)
curl https://www.example.com/

# Check Akamai debug headers
curl -H "Pragma: akamai-x-cache-on" https://www.example.com/
# Returns: X-Cache: TCP_HIT or TCP_MISS
```

### What's Next?

Now that the origin is configured:
1. ✅ Akamai knows where to fetch content
2. ✅ Host header and caching are configured
3. ✅ Rule tree is complete
4. ➡️ **Next Step**: Activate property to staging network

---

## Step 7: Property Activation

### What is Activation?

**Activation** deploys your property configuration to Akamai's edge network. Until activated, your configuration exists only as a draft.

```
Before Activation:
Property Version 1: [Draft] ──────┐
                                  │  Not deployed
                                  └─ Changes not live

After Activation to Staging:
Property Version 1: [Staging] ────┐
                                  │  Deployed to staging network
                                  └─ Can test at *.edgekey.net

After Activation to Production:
Property Version 1: [Production] ─┐
                                  │  Deployed to production network
                                  └─ Live for all users
```

### Why Two Networks?

| Network | Purpose | Use For |
|---------|---------|---------|
| **STAGING** | ✅ Testing environment | Validate changes before production |
| **PRODUCTION** | ✅ Live environment | Serve actual user traffic |

**Best Practice**: Always activate to staging first, test, then activate to production.

### Activation Lifecycle

```
1. Create Activation Request
   ├─ Specify property version
   ├─ Choose network (staging/production)
   ├─ Add activation note
   └─ Provide notification emails

2. Validation
   ├─ Akamai validates configuration
   ├─ Checks for errors/warnings
   └─ Returns validation results

3. Deployment
   ├─ Configuration pushed to edge servers
   ├─ Propagates across global network
   └─ Status: PENDING → ACTIVE

4. Completion
   └─ ✅ Configuration live on network

Time: 5-15 minutes for staging, 5-30 minutes for production
```

### Activation Status

| Status | Meaning | Action |
|--------|---------|--------|
| `PENDING` | Deployment in progress | Wait for completion |
| `ACTIVE` | ✅ Deployed successfully | Configuration is live |
| `FAILED` | Deployment failed | Check errors, fix, retry |
| `ABORTED` | Manually cancelled | Create new activation |
| `DEACTIVATED` | Removed from network | Property not serving |

### Code Walkthrough

#### File: `main.go` Lines 402-429

```go
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
            PropertyVersion:        prop.LatestVersion,  // Version to activate
            Network:                papi.ActivationNetworkStaging,  // Target network
            Note:                   "Auto-onboard via Go SDK",  // Change description
            NotifyEmails:           []string{"admin@example.com"},  // Notification emails
            AcknowledgeAllWarnings: true,  // Auto-accept warnings
        },
    })
    if err != nil {
        return fmt.Errorf("failed to create activation: %w", err)
    }

    fmt.Printf(">> Activation created with ID: %s\n", activationResp.ActivationID)
    fmt.Println(">> Note: Activation is now PENDING. Monitor status with GetActivation")
    fmt.Printf(">>       ActivationLink: %s\n", activationResp.ActivationLink)

    return nil
}
```

**What Happens**:
1. **Create Activation** - Submit property version for deployment
2. **Specify Network** - Choose staging or production
3. **Set Notifications** - Email addresses for status updates
4. **Acknowledge Warnings** - Auto-accept non-critical issues
5. **Return Activation ID** - Track deployment status

### Activation Structure

```go
type Activation struct {
    ActivationID           string    // act_12345
    PropertyVersion        int       // Version to activate
    Network                string    // STAGING or PRODUCTION
    Status                 string    // PENDING, ACTIVE, FAILED
    SubmitDate             string    // When submitted
    UpdateDate             string    // Last status change
    Note                   string    // Change description
    NotifyEmails           []string  // Email notifications
    AcknowledgeAllWarnings bool      // Auto-accept warnings
    ComplianceRecord       *ComplianceRecord  // For compliance requirements
}
```

### Activation Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `PropertyID` | string | ✅ Yes | Property to activate |
| `PropertyVersion` | int | ✅ Yes | Version number |
| `Network` | enum | ✅ Yes | STAGING or PRODUCTION |
| `Note` | string | ✅ Yes | Change description (audit trail) |
| `NotifyEmails` | []string | ✅ Yes | Email for status notifications |
| `AcknowledgeAllWarnings` | bool | No | Auto-accept warnings (default: false) |
| `ComplianceRecord` | object | No | For regulated industries |

### Network Types

```go
// Staging Network
papi.ActivationNetworkStaging
// Use for: Testing, validation, QA

// Production Network
papi.ActivationNetworkProduction
// Use for: Live traffic, actual users
```

### Activation Notes

**Purpose**: Audit trail, change tracking, compliance.

**Best Practices**:
```go
// ✅ Good notes
"Add new API endpoint /v2/users"
"Fix caching for images - increase TTL to 7d"
"Emergency: Revert origin to old server"
"Security: Add rate limiting to /api/*"

// ❌ Poor notes
"Update"
"Changes"
"Testing"
"Fix stuff"
```

### Validation Warnings

Akamai may return warnings during activation:

| Warning Type | Example | Action |
|--------------|---------|--------|
| **Informational** | "CP code usage noted" | Safe to ignore |
| **Performance** | "Caching disabled for path" | Review settings |
| **Security** | "HTTPS not enforced" | Consider fixing |
| **Best Practice** | "Origin timeout too high" | Recommended to adjust |

**AcknowledgeAllWarnings**:
- `true` - Auto-accept all warnings, activation proceeds
- `false` - Manual review required in Control Center

### Monitoring Activation

```go
// Check activation status
activationResp, err := papiClient.GetActivation(ctx, papi.GetActivationRequest{
    PropertyID:   prop.PropertyID,
    ActivationID: activationID,
})

fmt.Printf("Status: %s\n", activationResp.Activation.Status)
// PENDING → Wait
// ACTIVE → Deployment complete
// FAILED → Check error message
```

### Activation Timing

| Scenario | Expected Time |
|----------|---------------|
| **First staging activation** | 5-15 minutes |
| **Subsequent staging activations** | 5-10 minutes |
| **Production activation (small property)** | 5-15 minutes |
| **Production activation (large property)** | 15-30 minutes |
| **Production activation (many hostnames)** | 30-60 minutes |

**Factors affecting time**:
- Number of hostnames in property
- Number of edge servers in network
- Current network load
- Configuration complexity

### Testing Staging Activation

```bash
# Option 1: Use staging-specific hostname
curl -v https://www-example-com.edgesuite-staging.net/

# Option 2: Override Host header with edge hostname
curl -v -H "Host: www.example.com" https://www.example.com.edgekey.net/

# Option 3: Override Host header with edge server IP
# (Find staging edge server IP)
dig www.example.com.edgekey-staging.net

# Make request to staging IP
curl -v -H "Host: www.example.com" https://23.x.x.x/ --insecure

# Check Akamai debug headers
curl -H "Pragma: akamai-x-cache-on" \
     -H "Pragma: akamai-x-get-request-id" \
     -H "Host: www.example.com" \
     https://www.example.com.edgekey.net/
```

### Staging vs Production Differences

| Aspect | Staging | Production |
|--------|---------|------------|
| **Edge Servers** | Dedicated staging servers | Global production network |
| **Traffic** | Testing only | Real user traffic |
| **Certificate** | May use test certificate | Uses production certificate |
| **DNS** | Not in public DNS | Requires DNS configuration |
| **Performance** | May be slower | Optimized for speed |
| **Availability** | Lower SLA | High availability SLA |

### Activation to Production

```go
// After testing staging, activate to production
activationResp, err := papiClient.CreateActivation(ctx, papi.CreateActivationRequest{
    PropertyID: prop.PropertyID,
    ContractID: ContractID,
    GroupID:    GroupID,
    Activation: papi.Activation{
        PropertyVersion:        prop.LatestVersion,
        Network:                papi.ActivationNetworkProduction,  // Production network
        Note:                   "Deploying new origin configuration",
        NotifyEmails:           []string{"oncall@example.com"},
        AcknowledgeAllWarnings: false,  // Manual review for production
    },
})
```

### Simultaneous Activations

**Can you activate to both networks at once?**
- ✅ Yes, you can activate the same version to staging and production
- ✅ Activations are independent
- ❌ But wait for staging validation first (best practice)

**Workflow**:
```
1. Activate Version 3 to Staging → Wait → Test
2. Activate Version 3 to Production → Wait → Monitor
```

### Rollback

**How to rollback?**
- Activate a previous version to production

```go
// Rollback to version 2
activationResp, err := papiClient.CreateActivation(ctx, papi.CreateActivationRequest{
    PropertyID: prop.PropertyID,
    Activation: papi.Activation{
        PropertyVersion: 2,  // Previous version
        Network:        papi.ActivationNetworkProduction,
        Note:          "Emergency rollback - origin issues detected",
        NotifyEmails:  []string{"oncall@example.com"},
    },
})
```

### Compliance Records

For regulated industries (finance, healthcare):

```go
Activation: papi.Activation{
    ...
    ComplianceRecord: &papi.ComplianceRecord{
        NonComplianceReason: "EMERGENCY",  // Why bypassing approval
        PeerReviewer:        "john.doe@example.com",
        CustomerEmail:       "compliance@example.com",
        UnitTested:          true,
    },
}
```

### Activation Email Notifications

When you provide `NotifyEmails`, Akamai sends:

```
Subject: Akamai Property Activation - ACTIVE
To: admin@example.com

Property: www-example-com (prp_12345)
Version: 3
Network: STAGING
Status: ACTIVE
Activated: 2025-12-08 10:30:45 UTC
Note: Auto-onboard via Go SDK

Link: https://control.akamai.com/apps/property-manager/#/property-version/prp_12345/3/...
```

### Common Issues

| Error | Cause | Solution |
|-------|-------|----------|
| `Version not editable` | Trying to activate unsaved changes | Save changes first |
| `Validation failed` | Configuration errors | Review validation errors, fix |
| `Activation failed` | Edge deployment error | Check activation logs, retry |
| `Hostname not provisioned` | Edge hostname not ready | Wait for edge hostname deployment |
| `Certificate not deployed` | CPS enrollment pending | Wait for certificate deployment |
| `Warning requires acknowledgment` | Critical warning | Review warning, acknowledge |

### Activation Best Practices

```
1. ✅ Always test in staging first
2. ✅ Use descriptive activation notes
3. ✅ Monitor activation status
4. ✅ Test thoroughly before production
5. ✅ Activate during low-traffic windows
6. ✅ Have rollback plan ready
7. ✅ Set up monitoring/alerts
8. ❌ Don't activate to production untested
9. ❌ Don't ignore validation warnings
10. ❌ Don't activate without notification emails
```

### Idempotency Note

Activations are **NOT idempotent**:
- Each activation request creates a new activation
- Multiple calls create multiple activations (all will succeed)
- Check current activation status before creating new one

```go
// Check current activation status first
activations, _ := papiClient.GetActivations(ctx, papi.GetActivationsRequest{
    PropertyID: prop.PropertyID,
})

// Check if version is already active on staging
for _, act := range activations.Activations.Items {
    if act.PropertyVersion == prop.LatestVersion && 
       act.Network == "STAGING" && 
       act.Status == "ACTIVE" {
        fmt.Println("Version already active on staging")
        return nil
    }
}

// Create activation if not already active
...
```

### What's Next?

Now that the property is activated:
1. ✅ Configuration deployed to staging network
2. ✅ Can test at edge hostname
3. ✅ Ready for DNS configuration
4. ➡️ **Next Step**: Configure DNS to point your domain to Akamai

---

## Step 8: DNS Configuration

### What is DNS Configuration?

**DNS configuration** is the final step that directs user traffic to Akamai. You create a CNAME record that points your domain to the edge hostname.

```
Before DNS:
User types: www.example.com
DNS returns: Your origin IP (203.0.113.42)
Traffic flow: User → Origin (no Akamai)

After DNS:
User types: www.example.com
DNS returns: CNAME → www.example.com.edgekey.net → Edge Server IP (23.x.x.x)
Traffic flow: User → Akamai Edge → Origin
```

### Why CNAME Records?

**CNAME** (Canonical Name) is an alias record:

```
www.example.com  →  CNAME  →  www.example.com.edgekey.net
                               │
                               └─ Resolves to optimal edge server IP
```

**Benefits**:
- ✅ Akamai dynamically updates edge IPs (no DNS changes needed)
- ✅ Geo-routing automatically selects nearest server
- ✅ Load balancing across edge servers
- ✅ DDoS mitigation and failover

### DNS Record Types

| Record Type | Purpose | Example |
|-------------|---------|---------|
| **A Record** | Maps domain to IP | www.example.com → 203.0.113.42 |
| **CNAME** | Maps domain to another domain | www.example.com → www.example.com.edgekey.net |
| **ALIAS** | Like CNAME but for apex | example.com → www.example.com.edgekey.net |

### CNAME Configuration

```
Type:  CNAME
Name:  www
Value: www.example.com.edgekey.net
TTL:   300 (5 minutes)
```

**Explanation**:
- **Type**: CNAME (alias record)
- **Name**: Subdomain (`www`, `api`, `cdn`, etc.)
- **Value**: Edge hostname from Step 4
- **TTL**: Time-to-live (how long to cache DNS response)

### TTL (Time To Live)

| TTL Value | Duration | Use Case |
|-----------|----------|----------|
| **300** | 5 minutes | ✅ Initial deployment (easy to change) |
| **3600** | 1 hour | Standard production |
| **86400** | 24 hours | Very stable configuration |

**Recommendation**: Start with 300, increase after deployment stabilizes.

### DNS Providers

**Common DNS providers**:

| Provider | Management |
|----------|------------|
| **Route 53** (AWS) | AWS Console or CLI |
| **Cloudflare** | Cloudflare Dashboard |
| **GoDaddy** | GoDaddy DNS Manager |
| **Namecheap** | Namecheap Dashboard |
| **Azure DNS** | Azure Portal |
| **Google Cloud DNS** | GCP Console |

### Example: Route 53 (AWS)

```bash
# Using AWS CLI
aws route53 change-resource-record-sets \
  --hosted-zone-id Z1234567890ABC \
  --change-batch '{
    "Changes": [{
      "Action": "UPSERT",
      "ResourceRecordSet": {
        "Name": "www.example.com",
        "Type": "CNAME",
        "TTL": 300,
        "ResourceRecords": [{
          "Value": "www.example.com.edgekey.net"
        }]
      }
    }]
  }'
```

### Example: Cloudflare

```bash
# Using Cloudflare API
curl -X POST "https://api.cloudflare.com/client/v4/zones/{zone_id}/dns_records" \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  -H "Content-Type: application/json" \
  --data '{
    "type": "CNAME",
    "name": "www",
    "content": "www.example.com.edgekey.net",
    "ttl": 300,
    "proxied": false
  }'
```

**⚠️ Cloudflare Note**: Set `"proxied": false` to prevent Cloudflare from proxying (you want direct Akamai access).

### Apex Domain (example.com)

**Problem**: CNAME records don't work for apex/root domains.

```
❌ This doesn't work:
example.com  →  CNAME  →  example.com.edgekey.net
(CNAME not allowed at zone apex per RFC)

✅ Solutions:
1. ALIAS record (Route 53, Cloudflare)
2. ANAME record (DNS Made Easy)
3. Flatten CNAME (Cloudflare)
4. Use www.example.com and redirect apex
```

#### Route 53 ALIAS

```bash
aws route53 change-resource-record-sets \
  --hosted-zone-id Z1234567890ABC \
  --change-batch '{
    "Changes": [{
      "Action": "UPSERT",
      "ResourceRecordSet": {
        "Name": "example.com",
        "Type": "A",
        "AliasTarget": {
          "HostedZoneId": "Z2BJ6VIAURAKBL",
          "DNSName": "example.com.edgekey.net.",
          "EvaluateTargetHealth": false
        }
      }
    }]
  }'
```

### DNS Propagation

After creating CNAME, DNS changes take time to propagate:

```
1. Update DNS Record
   └─ Change made at DNS provider

2. Local Propagation
   └─ Your DNS servers update (seconds)

3. Global Propagation
   ├─ ISP DNS servers fetch new record
   ├─ Public DNS servers (8.8.8.8) update
   └─ Worldwide propagation (minutes to hours)

4. Full Propagation
   └─ All resolvers have new record
```

| Timeframe | What's Updated |
|-----------|----------------|
| **Immediate** | Your DNS provider |
| **5-10 minutes** | Major DNS servers (Google, Cloudflare) |
| **1 hour** | Most ISP resolvers |
| **24-48 hours** | All global resolvers |

### Testing DNS

```bash
# Check DNS resolution locally
dig www.example.com

# Expected output:
# www.example.com.    300    IN    CNAME    www.example.com.edgekey.net.
# www.example.com.edgekey.net.    20    IN    A    23.x.x.x

# Check from Google's DNS
dig @8.8.8.8 www.example.com

# Check from Cloudflare's DNS
dig @1.1.1.1 www.example.com

# Check propagation globally
# Use online tools:
# - https://www.whatsmydns.net/
# - https://dnschecker.org/
```

### Testing End-to-End

```bash
# Test HTTP
curl -v http://www.example.com/

# Test HTTPS
curl -v https://www.example.com/

# Check certificate
curl -v https://www.example.com/ 2>&1 | grep "subject:"
# Should show certificate CN matching your domain

# Check Akamai headers
curl -v https://www.example.com/ 2>&1 | grep -i "x-cache"
# Should show: X-Cache: TCP_HIT or TCP_MISS (from Akamai)

# Check server header
curl -I https://www.example.com/ | grep -i "server:"
# May show: Server: AkamaiGHost (Akamai's edge server)
```

### DNS Configuration in Code (Mock)

#### File: `main.go` Lines 435-443

```go
func configureCNAMEInDNS(domain, edgeHostname string) {
    fmt.Println("---------------------------------------------------------")
    fmt.Println(">> 📡 MOCK DNS PROVIDER API")
    fmt.Printf(">> ACTION: UPSERT CNAME RECORD\n")
    fmt.Printf(">> KEY:    %s\n", domain)
    fmt.Printf(">> VALUE:  %s\n", edgeHostname)
    fmt.Println(">> STATUS: 200 OK (Propagating...)")
    fmt.Println("---------------------------------------------------------")
}
```

**Why Mock?**
- DNS providers have different APIs (Route 53, Cloudflare, etc.)
- Manual DNS configuration is common
- Automated DNS requires provider-specific SDK

### Multiple Hostnames

```bash
# If property handles multiple domains:
www.example.com   →  CNAME  →  www.example.com.edgekey.net
api.example.com   →  CNAME  →  api.example.com.edgekey.net
cdn.example.com   →  CNAME  →  cdn.example.com.edgekey.net
```

### Redirect Apex to WWW

Common pattern: Redirect `example.com` to `www.example.com`

```bash
# Option 1: DNS provider's redirect service
example.com  →  HTTP 301 Redirect  →  https://www.example.com

# Option 2: Akamai edge redirect
# Configure in property rule tree:
{
  "name": "Redirect Apex",
  "criteria": [
    {"name": "hostname", "options": {"values": ["example.com"]}}
  ],
  "behaviors": [
    {
      "name": "redirect",
      "options": {
        "destinationHostname": "www.example.com",
        "responseCode": 301
      }
    }
  ]
}
```

### DNSSEC

**What**: DNS Security Extensions (prevents DNS spoofing)

**Akamai Support**: ✅ Supported

**Configuration**:
1. Enable DNSSEC at your DNS provider
2. Add DS records from Akamai to parent zone
3. Akamai signs DNS responses

### DNS Failover

**Automatic failover** if Akamai edge is unavailable:

```bash
# Configure multiple CNAME records (depends on DNS provider)
www.example.com  →  CNAME  →  www.example.com.edgekey.net (primary)
                    CNAME  →  backup.example.com (failover)
```

### Common Issues

| Error | Cause | Solution |
|-------|-------|----------|
| `CNAME resolution failed` | Wrong edge hostname | Verify edge hostname spelling |
| `DNS propagation slow` | High TTL on old record | Wait or flush DNS cache |
| `Certificate name mismatch` | Domain not in certificate | Add domain to CPS enrollment |
| `Too many DNS lookups` | CNAME chain too long | Check for circular references |
| `Apex domain CNAME error` | RFC violation | Use ALIAS/ANAME or redirect |

### DNS Cache Flushing

```bash
# Flush local DNS cache

# macOS
sudo dscacheutil -flushcache
sudo killall -HUP mDNSResponder

# Windows
ipconfig /flushdns

# Linux
sudo systemd-resolve --flush-caches

# Chrome browser
chrome://net-internals/#dns → Clear host cache
```

### Monitoring DNS

```bash
# Monitor DNS continuously
watch -n 5 'dig +short www.example.com'

# Expected progression:
# (empty)                          ← Old record expired
# www.example.com.edgekey.net.     ← CNAME appears
# 23.x.x.x                         ← Edge server IP appears
```

### What's Next?

Now that DNS is configured:
1. ✅ User domain points to Akamai
2. ✅ Traffic flows through Akamai edge
3. ✅ Content is cached and accelerated
4. ➡️ **Next Step (Optional)**: Add security with AppSec/WAF

---

## Step 9: Application Security (AppSec/WAF)

### What is AppSec?

**AppSec (Application Security)** is Akamai's security product that provides:
- ✅ **WAF** (Web Application Firewall) - Blocks malicious requests
- ✅ **DDoS Protection** - Mitigates denial-of-service attacks
- ✅ **Bot Management** - Identifies and blocks bots
- ✅ **API Security** - Protects API endpoints
- ✅ **Rate Limiting** - Prevents abuse

### Why AppSec?

```
Without AppSec:
User Request → Akamai Edge → Origin
               └─ All requests passed through (vulnerable)

With AppSec:
User Request → Akamai Edge → [Security Check] → Origin
                              ├─ SQL injection? → Block (403)
                              ├─ DDoS attack? → Rate limit
                              ├─ Bad bot? → Challenge
                              └─ Legitimate? → Allow
```

### Security Threats Blocked

| Threat Type | Description | Example |
|-------------|-------------|---------|
| **SQL Injection** | Database attacks | `?id=1' OR '1'='1` |
| **XSS** | Cross-site scripting | `<script>alert('xss')</script>` |
| **Path Traversal** | File system access | `../../etc/passwd` |
| **Command Injection** | OS command execution | `; rm -rf /` |
| **DDoS** | Denial of service | 100k req/sec flood |
| **Bots** | Automated scrapers | Content theft, price scraping |
| **Brute Force** | Password guessing | 1000 login attempts |

### AppSec Hierarchy

```
Security Configuration
├─ Configuration ID: 12345
├─ Version: 3 (latest), 2 (staging), 1 (production)
│
├─ Match Targets (which hostnames to protect)
│  ├─ Target 1: www.example.com, api.example.com
│  └─ Target 2: admin.example.com
│
├─ Security Policies
│  ├─ Policy: "Web Application Protection"
│  │  ├─ WAF Mode: KRS (Akamai Kona Rule Set)
│  │  ├─ Attack Groups: SQL Injection, XSS, etc.
│  │  └─ Rate Limiting: 1000 req/min per IP
│  │
│  └─ Policy: "API Protection"
│     ├─ WAF Mode: Custom rules
│     └─ Rate Limiting: 100 req/sec per API key
│
└─ Advanced Settings
   ├─ IP/Geo Blocking
   ├─ Custom Rules
   └─ Exception Rules
```

### Security Configuration Structure

```go
type SecurityConfiguration struct {
    ID                    int       // Config ID: 12345
    Name                  string    // "Production Security"
    LatestVersion         int       // 3 (editable)
    StagingVersion        int       // 2 (active on staging)
    ProductionVersion     int       // 1 (active on production)
    ProductionHostnames   []string  // Protected hostnames
    StagingHostnames      []string  // Staging protected hostnames
}
```

### Code Walkthrough

#### File: `main.go` Lines 449-502

```go
func onboardToWAF(ctx context.Context, appsecClient appsec.APPSEC, domain string) error {
    fmt.Println(">> Onboarding hostname into WAF/AppSec")
    fmt.Printf("   Hostname: %s\n", domain)

    // 1. List available security configurations
    configsResp, err := appsecClient.GetConfigurations(ctx, 
        appsec.GetConfigurationsRequest{})
    if err != nil {
        return fmt.Errorf("failed to list AppSec configurations: %w", err)
    }

    // 2. Check if configurations exist
    if len(configsResp.Configurations) == 0 {
        fmt.Println(">> Note: No AppSec configurations found")
        fmt.Println(">>       You need to create a security configuration first")
        fmt.Println(">>       This typically includes:")
        fmt.Println(">>       1. Creating a security configuration")
        fmt.Println(">>       2. Creating security policies")
        fmt.Println(">>       3. Configuring WAF rules and protections")
        return fmt.Errorf("no AppSec configurations available")
    }

    // 3. Use the first available configuration
    config := configsResp.Configurations[0]
    fmt.Printf(">> Using AppSec Configuration: %s (ID: %d, Version: %d)\n",
        config.Name, config.ID, config.LatestVersion)

    // 4. Check if hostname is already in the configuration
    hostnameExists := false
    for _, h := range config.ProductionHostnames {
        if h == domain {
            hostnameExists = true
            fmt.Printf(">> Hostname %s already protected by AppSec\n", domain)
            break
        }
    }

    // 5. If hostname not protected, show manual steps
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
```

**What Happens**:
1. **List Configurations** - Get available security configurations
2. **Select Configuration** - Use existing security config
3. **Check Hostname** - See if domain is already protected
4. **Manual Steps** - Show how to add hostname protection

### Why Manual AppSec Setup?

**AppSec configuration is complex** and requires:
- Security policy selection (WAF rules, rate limits)
- Match target configuration (which paths/hostnames)
- Attack group tuning (prevent false positives)
- Custom rule creation
- Testing and validation

**Best Practice**: Configure AppSec in Control Center first, then reference in code.

### Creating Security Configuration (Manual)

```
1. Control Center → Security → Application Security

2. Create Configuration
   ├─ Name: "Production Web Protection"
   ├─ Contract: Select your contract
   └─ Hostnames: Select property hostnames

3. Create Security Policy
   ├─ Name: "Default Web Protection"
   ├─ Policy Type: Website
   └─ WAF Mode: KRS (Kona Rule Set)

4. Configure Attack Groups
   ├─ SQL Injection: Alert (or Deny)
   ├─ Cross-Site Scripting: Alert
   ├─ Command Injection: Deny
   └─ ... (30+ attack groups)

5. Create Match Target
   ├─ Hostnames: www.example.com
   ├─ Paths: /*
   └─ Apply Policy: "Default Web Protection"

6. Test in Staging
   ├─ Activate to staging
   ├─ Test malicious requests
   └─ Review security events

7. Deploy to Production
   └─ Activate to production
```

### WAF Modes

| Mode | Description | Use When |
|------|-------------|----------|
| **KRS** | Kona Rule Set (Akamai-managed) | ✅ Standard web apps (recommended) |
| **ASE** | Adaptive Security Engine | Auto-learning, less tuning |
| **Custom** | Your own rules | Specific app requirements |

### Attack Groups

**KRS includes 30+ attack groups**:

| Attack Group | Blocks |
|--------------|--------|
| SQL Injection | `' OR 1=1`, `UNION SELECT`, etc. |
| Cross-Site Scripting | `<script>`, `javascript:`, etc. |
| Command Injection | `; ls`, `| cat /etc/passwd` |
| Path Traversal | `../../../etc/passwd` |
| Protocol Violations | Malformed HTTP requests |
| Trojans | Known malware signatures |
| Outbound | Sensitive data leakage |

**Actions per group**:
- **Alert** - Log but allow (for tuning)
- **Deny** - Block with 403 response
- **None** - Disable group

### Rate Limiting

Protect against abuse and DDoS:

```go
// Rate Policy Example
{
    "name": "API Rate Limit",
    "requests": 1000,          // Max requests
    "interval": 60,            // Per 60 seconds
    "action": "ALERT",         // Or DENY
    "burstWindow": 5,          // Allow burst
    "pathMatchType": "PREFIX",
    "path": "/api/",
}
```

### IP/Geo Blocking

```go
// Block specific countries
IPGeoFirewall: {
    "block": "GEO",
    "geoControls": {
        "blockedIPNetworkLists": ["China", "Russia"],
    },
}

// Block specific IP ranges
IPGeoFirewall: {
    "block": "IP",
    "ipControls": {
        "blockedIPNetworkLists": ["malicious_ips"],
    },
}
```

### Match Targets

**Match targets** define which requests get security policies:

```go
type MatchTarget struct {
    ConfigID:       12345,
    TargetID:       123,
    Type:           "website",  // or "api"
    Hostnames:      []string{"www.example.com", "api.example.com"},
    FilePaths:      []string{"/*"},
    SecurityPolicy: "default_policy",
}
```

**Examples**:

```go
// Protect entire website
{
    "hostnames": ["www.example.com"],
    "filePaths": ["/*"],
    "policy": "Web Protection",
}

// Protect API endpoints only
{
    "hostnames": ["api.example.com"],
    "filePaths": ["/api/*"],
    "policy": "API Protection",
}

// Protect admin area with stricter rules
{
    "hostnames": ["www.example.com"],
    "filePaths": ["/admin/*", "/dashboard/*"],
    "policy": "Admin Protection",
}
```

### Exception Rules

Allow legitimate traffic that triggers WAF:

```go
// Exception: Allow specific query parameter
{
    "conditions": [
        {"path": "/search"},
        {"query": "q=<script>"},  // Triggers XSS
    ],
    "exception": "XSS_ATTACK_GROUP",
    "action": "ALLOW",
}
```

### Bot Management

Identify and manage bots:

| Bot Type | Action |
|----------|--------|
| **Search Engine** (Google, Bing) | Allow |
| **Legitimate Bot** (monitoring) | Allow |
| **Scraper Bot** | Challenge or block |
| **Malicious Bot** | Block |

### Security Events

View attacks blocked:

```bash
# In Control Center → Security → Events

Event Log:
- Time: 2025-12-08 10:45:23
- Client IP: 198.51.100.50
- Hostname: www.example.com
- Path: /search?q=<script>alert(1)</script>
- Attack: XSS (Cross-Site Scripting)
- Action: DENIED (403)
- Policy: Default Web Protection
```

### Testing WAF

```bash
# Test SQL injection block
curl "https://www.example.com/search?q=1' OR '1'='1"
# Expected: 403 Forbidden

# Test XSS block
curl "https://www.example.com/search?q=<script>alert(1)</script>"
# Expected: 403 Forbidden

# Test legitimate request
curl "https://www.example.com/search?q=akamai"
# Expected: 200 OK
```

### Tuning WAF

**Initial deployment**:
1. Start in **Alert Mode** (log but don't block)
2. Monitor security events for false positives
3. Create exceptions for legitimate traffic
4. Switch to **Deny Mode** after tuning

**Timeframe**: 1-2 weeks of monitoring recommended.

### Activation (AppSec)

Similar to property activation:

```go
// Activate security configuration
activationResp, err := appsecClient.CreateActivation(ctx, 
    appsec.CreateActivationRequest{
        ConfigID: config.ID,
        Version:  config.LatestVersion,
        Network:  appsec.ActivationNetworkStaging,
        Note:     "Add WAF protection for www.example.com",
        NotifyEmails: []string{"security@example.com"},
    })
```

### Common Issues

| Error | Cause | Solution |
|-------|-------|----------|
| `No configurations found` | AppSec not provisioned | Contact Akamai to enable AppSec |
| `Hostname not in property` | Domain not in PAPI property | Add to property first |
| `False positives` | Legitimate traffic blocked | Create exception rules |
| `Activation failed` | Configuration errors | Review validation errors |

### Security Best Practices

```
1. ✅ Start with KRS (Kona Rule Set)
2. ✅ Use Alert mode initially, then Deny mode
3. ✅ Monitor security events regularly
4. ✅ Create exceptions for false positives
5. ✅ Enable rate limiting for APIs
6. ✅ Block high-risk countries if not needed
7. ✅ Test malicious requests in staging
8. ❌ Don't skip tuning period
9. ❌ Don't block all requests accidentally
10. ❌ Don't ignore security alerts
```

### What's Next?

Now that security is configured:
1. ✅ Domain protected by WAF
2. ✅ DDoS mitigation active
3. ✅ Bot detection enabled
4. ✅ **ONBOARDING COMPLETE!**

---

## Complete Code Reference

### Full Onboarding Flow

Here's the complete `main.go` with all steps integrated:

```go
package main

import (
    "context"
    "fmt"
    "log"
)

func main() {
    ctx := context.Background()

    // Step 1: Authentication
    sess, err := newSession(EdgercPath, EdgercSection)
    if err != nil {
        log.Fatalf("Authentication failed: %v", err)
    }

    // Initialize clients
    cpsClient := cps.Client(sess)
    papiClient := papi.Client(sess)
    appsecClient := appsec.Client(sess)

    fmt.Println(">> Starting Onboarding Flow for:", UserDomain)

    // Step 2: SSL Certificate (CPS)
    enrollmentID, err := getOrCreateEnrollment(ctx, cpsClient, UserDomain)
    if err != nil {
        log.Fatalf("CPS Error: %v", err)
    }
    fmt.Printf("✅ Certificate Enrollment ID: %d\n", enrollmentID)

    // Step 3: Property Creation
    prop, err := getOrCreateProperty(ctx, papiClient, PropertyName)
    if err != nil {
        log.Fatalf("Property Error: %v", err)
    }
    fmt.Printf("✅ Property ID: %s (Version: %d)\n", prop.PropertyID, prop.LatestVersion)

    // Step 4: Edge Hostname
    edgeHostname, err := ensureEdgeHostname(ctx, papiClient, enrollmentID, UserDomain)
    if err != nil {
        log.Fatalf("Edge Hostname Error: %v", err)
    }
    fmt.Printf("✅ Edge Hostname: %s\n", edgeHostname)

    // Step 5 & 6: Hostname Mapping + Origin Configuration
    if err := updatePropertyRules(ctx, papiClient, prop, UserDomain, edgeHostname); err != nil {
        log.Fatalf("Configuration Error: %v", err)
    }
    fmt.Println("✅ Hostname mapping and origin configured")

    // Step 7: Activation to Staging
    if err := activateToStaging(ctx, papiClient, prop); err != nil {
        log.Fatalf("Activation Error: %v", err)
    }
    fmt.Println("✅ Activated to STAGING network")

    // Step 8: DNS Configuration (Manual)
    configureCNAMEInDNS(UserDomain, edgeHostname)
    fmt.Println("✅ DNS instructions provided")

    // Step 9: AppSec/WAF (Optional)
    if err := onboardToWAF(ctx, appsecClient, UserDomain); err != nil {
        log.Printf("AppSec Note: %v", err)
    }
    fmt.Println("✅ Security configuration checked")

    fmt.Println("\n🎉 ONBOARDING COMPLETE!")
    fmt.Printf("\nNext Steps:\n")
    fmt.Printf("1. Wait for staging activation (5-15 min)\n")
    fmt.Printf("2. Test: curl -H 'Host: %s' https://%s/\n", UserDomain, edgeHostname)
    fmt.Printf("3. Configure DNS: CNAME %s → %s\n", UserDomain, edgeHostname)
    fmt.Printf("4. Activate to production\n")
}
```

### Configuration Constants

```go
const (
    EdgercPath    = "~/.edgerc"
    EdgercSection = "default"
    
    ContractID   = "ctr_1-12345"      // Your Contract ID
    GroupID      = "grp_12345"        // Your Group ID
    ProductID    = "prd_Ion"          // Your Product ID
    PropertyName = "my-api-gateway"   // Property name
    UserDomain   = "api.example.com"  // Your domain
)
```

**🔧 Customize these values** before running!

### Helper Functions Summary

| Function | Purpose | Returns |
|----------|---------|---------|
| `newSession()` | Create authenticated session | Session |
| `getOrCreateEnrollment()` | Find/create SSL certificate | Enrollment ID |
| `getOrCreateProperty()` | Find/create property | Property object |
| `ensureEdgeHostname()` | Find/create edge hostname | Edge hostname string |
| `updatePropertyRules()` | Configure hostname + origin | Error |
| `activateToStaging()` | Deploy to staging | Error |
| `configureCNAMEInDNS()` | DNS instructions (mock) | - |
| `onboardToWAF()` | Check security config | Error |

### Running the Code

```bash
# 1. Update configuration constants
vim main.go
# Set ContractID, GroupID, ProductID, PropertyName, UserDomain

# 2. Ensure .edgerc exists
cat ~/.edgerc
# Should contain valid credentials

# 3. Run onboarding
go run main.go

# 4. Monitor output
# Look for ✅ checkmarks and any errors

# 5. Check activation status
# Use Control Center or GetActivation API

# 6. Configure DNS manually
# Create CNAME: www.example.com → www.example.com.edgekey.net

# 7. Test end-to-end
curl https://www.example.com/
```

### Code Organization

```
project/
├── main.go                  # Main onboarding flow
├── property_helpers.go      # Property type helpers
├── product_utils.go         # Product ID discovery
├── examples.go              # Usage examples
├── go.mod                   # Dependencies
├── .edgerc                  # Credentials (never commit!)
│
├── docs/
│   ├── ONBOARDING_GUIDE.md  # This guide
│   ├── CONCEPTS.md          # Glossary (to be created)
│   ├── ARCHITECTURE.md      # Flow diagrams (to be created)
│   └── TROUBLESHOOTING.md   # Common issues (to be created)
│
├── README.md                # Project overview
├── PRODUCT_IDS.md           # Product ID explanation
├── PROPERTY_HELPERS.md      # Helper API docs
└── QUICK_START.md           # Quick reference
```

### Property Helpers

For production use, consider using `property_helpers.go`:

```go
// Example: Create Ion Standard property
helper := NewPropertyHelper(
    PropertyTypeIonStandard,
    "www-example-com",
    "www.example.com",
    "origin.example.com",
)

// Idempotent operations
err := helper.EnsureProperty(ctx, papiClient, contractID, groupID)
err = helper.EnsureEdgeHostname(ctx, papiClient, certEnrollmentID)
err = helper.MapHostname(ctx, papiClient)
err = helper.ConfigureOrigin(ctx, papiClient)
err = helper.ActivateStaging(ctx, papiClient)
```

See `PROPERTY_HELPERS.md` for full documentation.

---

## Troubleshooting

### Authentication Issues

#### `401 Unauthorized`

**Cause**: Invalid credentials or expired access token.

**Solution**:
```bash
# 1. Verify .edgerc file exists and has correct format
cat ~/.edgerc

# 2. Check credentials in Control Center
# Login → Identity & Access → API User → Your Client

# 3. Regenerate credentials if needed

# 4. Verify section name matches
# Code: EdgercSection = "default"
# File: [default]
```

#### `403 Forbidden`

**Cause**: Insufficient API permissions.

**Solution**:
```bash
# Grant required permissions:
# - PAPI: READ-WRITE
# - CPS: READ-WRITE
# - Edge Hostnames: READ-WRITE
# - AppSec: READ-WRITE (if using)
```

#### `edgegrid.New: no such file`

**Cause**: `.edgerc` file doesn't exist.

**Solution**:
```bash
# Create .edgerc file
mkdir -p ~/.edgerc
chmod 600 ~/.edgerc

# Add credentials:
cat > ~/.edgerc << EOF
[default]
host = akaa-xxxxxxxxx.luna.akamaiapis.net
client_token = akab-xxxxx
client_secret = xxxxx
access_token = akab-xxxxx
EOF
```

---

### Certificate (CPS) Issues

#### `No enrollment found`

**Cause**: No certificate enrollment exists for domain.

**Solution**:
```bash
# Manual enrollment required
# 1. Control Center → Certificate Provisioning
# 2. Create enrollment with your domain
# 3. Complete domain validation
# 4. Wait for deployment (15-60 min)
# 5. Run script again
```

#### `Enrollment status: PENDING`

**Cause**: Domain validation not complete.

**Solution**:
```bash
# Complete domain validation:

# Option 1: DNS Validation
# Add TXT record provided by Akamai
dig TXT _acme-challenge.www.example.com

# Option 2: HTTP Validation
# Place challenge file on origin
curl http://www.example.com/.well-known/acme-challenge/token

# Check enrollment status
# Control Center → Certificate Provisioning → Enrollments
```

#### `Enrollment status: FAILED`

**Cause**: Validation failed or CA error.

**Solution**:
```bash
# 1. Check validation method completed correctly
# 2. Review error message in Control Center
# 3. Retry validation or create new enrollment
# 4. Contact Akamai support if persistent
```

---

### Property Issues

#### `Property name already exists`

**Cause**: Property with that name already exists.

**Solution**:
```go
// Option 1: Use different name
PropertyName = "my-api-gateway-v2"

// Option 2: Fetch existing property
// Code automatically finds existing property
```

#### `Invalid product ID`

**Cause**: Product not available in contract.

**Solution**:
```go
// List available products
productsResp, _ := papiClient.GetProducts(ctx, papi.GetProductsRequest{
    ContractID: ContractID,
})
for _, product := range productsResp.Products.Items {
    fmt.Printf("%s: %s\n", product.ProductID, product.ProductName)
}

// Use correct product ID
ProductID = "prd_Ion"  // Or whatever is available
```

#### `Contract not found`

**Cause**: Invalid contract ID.

**Solution**:
```bash
# Find contract ID:
# Control Center → Property Manager → Any Property → Details

# Or via API:
contracts, _ := papiClient.GetContracts(ctx)
```

---

### Edge Hostname Issues

#### `Edge hostname already exists`

**Cause**: Edge hostname with that prefix already exists.

**Solution**:
```go
// Option 1: Use existing edge hostname (code does this automatically)

// Option 2: Use different domain prefix
UserDomain = "www2.example.com"  // Creates www2.edgekey.net
```

#### `Invalid certificate enrollment`

**Cause**: Enrollment ID doesn't exist or isn't deployed.

**Solution**:
```bash
# Verify enrollment:
# 1. Check enrollment ID from Step 2
# 2. Verify status is DEPLOYED
# 3. Wait for deployment if PENDING
```

#### `Certificate not deployed`

**Cause**: CPS enrollment not ready.

**Solution**:
```bash
# Wait for certificate deployment:
# Status: PENDING → VALIDATED → DEPLOYED (15-60 min)

# Check status via API or Control Center
```

---

### Hostname Mapping Issues

#### `Hostname already in use`

**Cause**: Another property already uses this hostname.

**Solution**:
```bash
# Option 1: Remove from other property first
# Control Center → Property Manager → Other Property → Hostnames

# Option 2: Use different hostname
UserDomain = "www2.example.com"
```

#### `Version not editable`

**Cause**: Trying to edit an activated version.

**Solution**:
```bash
# Create new version:
# Control Center → Property Manager → Property → Create New Version

# Or via API:
newVersion, _ := papiClient.CreatePropertyVersion(ctx, ...)
```

---

### Origin Issues

#### `Origin not reachable`

**Cause**: Origin hostname doesn't resolve or isn't accessible.

**Solution**:
```bash
# Test origin accessibility
dig origin.example.com
curl -v http://origin.example.com/

# Verify:
# - DNS resolves correctly
# - Firewall allows Akamai edge IPs
# - Origin server is running
# - Origin ports (80/443) are open
```

#### `Origin SSL handshake failed`

**Cause**: Origin certificate is invalid or self-signed.

**Solution**:
```go
// Option 1: Fix origin certificate
// Ensure valid certificate on origin

// Option 2: Adjust verification mode
"verificationMode": "CUSTOM",  // Allow self-signed
```

#### `5xx errors from origin`

**Cause**: Origin server returning errors.

**Solution**:
```bash
# Check origin server logs
# Test origin directly (bypass Akamai)
curl -H "Host: www.example.com" http://origin.example.com/

# Verify origin can handle requests
```

---

### Activation Issues

#### `Validation failed`

**Cause**: Configuration errors in property.

**Solution**:
```bash
# Review validation errors:
# Control Center → Property Manager → Property → Activate

# Common issues:
# - Missing origin behavior
# - Invalid hostname mapping
# - Certificate not deployed
# - CP code not assigned
```

#### `Activation failed`

**Cause**: Edge deployment error.

**Solution**:
```bash
# Check activation logs
# Control Center → Property Manager → Activation History → View Logs

# Common fixes:
# - Wait and retry
# - Fix validation errors
# - Contact Akamai support
```

#### `Hostname not provisioned`

**Cause**: Edge hostname not ready.

**Solution**:
```bash
# Wait for edge hostname deployment (5-15 min)
# Verify edge hostname status:
edgeHostname, _ := papiClient.GetEdgeHostname(ctx, ...)
fmt.Println(edgeHostname.Status)  // Should be ACTIVE
```

---

### DNS Issues

#### `CNAME resolution failed`

**Cause**: DNS record not configured or typo.

**Solution**:
```bash
# Verify CNAME record
dig www.example.com

# Should return:
# www.example.com.  300  IN  CNAME  www.example.com.edgekey.net.

# If not, create CNAME in DNS provider
```

#### `DNS propagation slow`

**Cause**: High TTL on old record.

**Solution**:
```bash
# Check TTL on old record
dig www.example.com

# Wait for TTL expiration, then:
# - Flush local DNS cache
# - Test from different locations
# - Use online DNS checkers (whatsmydns.net)
```

#### `Certificate name mismatch`

**Cause**: Domain not in certificate SANs.

**Solution**:
```bash
# Check certificate includes your domain:
# Control Center → Certificate Provisioning → Enrollment

# If missing:
# - Add domain to enrollment SANs
# - Wait for certificate reissue
# - Update edge hostname
```

#### `CNAME at apex domain error`

**Cause**: CNAME records not allowed at zone apex per RFC.

**Solution**:
```bash
# Option 1: Use ALIAS record (Route 53, Cloudflare)
example.com  ALIAS  example.com.edgekey.net

# Option 2: Use ANAME record (DNS Made Easy)
example.com  ANAME  example.com.edgekey.net

# Option 3: Redirect apex to www
example.com  →  301 redirect  →  www.example.com
www.example.com  CNAME  www.example.com.edgekey.net
```

---

### AppSec Issues

#### `No configurations found`

**Cause**: AppSec not provisioned on contract.

**Solution**:
```bash
# Contact Akamai account team to:
# - Enable AppSec product
# - Add to your contract
# - Create initial security configuration
```

#### `Hostname not in property`

**Cause**: Domain not added to PAPI property first.

**Solution**:
```bash
# Add hostname to property first (Step 5)
# Then add to AppSec configuration
```

#### `False positives blocking traffic`

**Cause**: Legitimate requests triggering WAF rules.

**Solution**:
```bash
# 1. Review security events
# Control Center → Security → Events

# 2. Identify false positives
# Look for blocked legitimate requests

# 3. Create exception rules
# Security → Configuration → Exception Rules

# 4. Test exceptions in staging first
```

---

### General Debugging

#### Enable Debug Logging

```go
// Add to main.go
import "os"

func main() {
    // Enable HTTP request logging
    os.Setenv("AKAMAI_DEBUG", "true")
    
    // Your code...
}
```

#### Check API Response

```go
// Capture detailed error information
activationResp, err := papiClient.CreateActivation(ctx, req)
if err != nil {
    // Print full error
    fmt.Printf("Error: %+v\n", err)
    
    // Check if API-specific error
    if apiErr, ok := err.(*papi.Error); ok {
        fmt.Printf("Status: %d\n", apiErr.StatusCode)
        fmt.Printf("Detail: %s\n", apiErr.Detail)
    }
}
```

#### Test in Stages

```go
// Test each step independently

// Step 1: Auth
sess, _ := newSession(...)
fmt.Println("✅ Auth successful")

// Step 2: CPS
enrollment, _ := getOrCreateEnrollment(...)
fmt.Printf("✅ Enrollment: %d\n", enrollment)

// Continue testing each step...
```

---

### Getting Help

| Resource | URL |
|----------|-----|
| **Akamai Documentation** | https://techdocs.akamai.com/ |
| **API Reference** | https://developer.akamai.com/ |
| **SDK Issues** | https://github.com/akamai/AkamaiOPEN-edgegrid-golang/issues |
| **Support Portal** | https://control.akamai.com/apps/support/ |
| **Community Forums** | https://community.akamai.com/ |

### Support Ticket Information

When creating a support ticket, include:

```
Subject: [Property Onboarding Issue] <Brief description>

Details:
- Contract ID: ctr_1-12345
- Group ID: grp_12345
- Property ID: prp_12345
- Domain: www.example.com
- Error Message: <Full error>
- Steps to Reproduce: <What you did>
- Expected vs Actual: <What should happen vs what happened>
- Timestamp: 2025-12-08 10:30:45 UTC
- Request ID: <From Pragma: akamai-x-get-request-id>
```

---

## Summary

### What You've Learned

✅ **Authentication**: How to connect to Akamai APIs with EdgeGrid  
✅ **Certificates**: Managing SSL/TLS with CPS enrollments  
✅ **Properties**: Creating and configuring delivery properties  
✅ **Edge Hostnames**: Setting up CDN entry points  
✅ **Hostname Mapping**: Linking user domains to edge hostnames  
✅ **Origin Configuration**: Telling Akamai where your server is  
✅ **Activation**: Deploying configurations to staging/production  
✅ **DNS**: Pointing domains to Akamai with CNAME records  
✅ **Security**: Protecting with AppSec/WAF  

### Onboarding Checklist

```
Prerequisites:
☐ Akamai account and contract
☐ API credentials created
☐ Contract ID, Group ID, Product ID identified
☐ Domain to onboard (e.g., www.example.com)
☐ Origin server accessible (e.g., origin.example.com)
☐ DNS access to create CNAME records

Step-by-Step:
☐ Create .edgerc with credentials
☐ Update constants in main.go
☐ Create or find certificate enrollment
☐ Create property configuration
☐ Create edge hostname with certificate
☐ Map user domain to edge hostname
☐ Configure origin server
☐ Activate to staging network
☐ Test staging deployment
☐ Configure DNS CNAME record
☐ Wait for DNS propagation
☐ Test production with user domain
☐ Activate to production (optional)
☐ Configure AppSec/WAF (optional)

Validation:
☐ curl https://www.example.com/ returns 200
☐ Certificate matches domain
☐ X-Cache headers present (Akamai)
☐ Content serves correctly
☐ Origin receives requests with correct Host header
☐ Security events logging (if AppSec enabled)
```

### Next Steps

1. **Production Deployment**
   - Test staging thoroughly
   - Activate to production
   - Monitor edge traffic
   - Set up alerts

2. **Optimization**
   - Tune caching behaviors
   - Add compression
   - Enable image optimization
   - Configure adaptive acceleration

3. **Advanced Features**
   - Multi-origin failover
   - Path-based routing
   - API acceleration
   - Video streaming

4. **Monitoring**
   - Set up CloudMonitor
   - Configure alerts
   - Review analytics
   - Track performance metrics

5. **Security Hardening**
   - Tune WAF rules
   - Enable rate limiting
   - Configure IP/Geo blocking
   - Set up bot management

### Additional Resources

**Documentation**:
- `README.md` - Project overview
- `PRODUCT_IDS.md` - Product ID reference
- `PROPERTY_HELPERS.md` - Helper function API
- `QUICK_START.md` - Quick reference guide
- `examples.go` - Working code examples

**Code Files**:
- `main.go` - Core onboarding flow
- `property_helpers.go` - Property type helpers
- `product_utils.go` - Product discovery utilities

### Feedback

This guide is a living document. Suggestions welcome:
- GitHub Issues: https://github.com/akamai/AkamaiOPEN-edgegrid-golang/issues
- Email: developer@akamai.com

---

## Appendix

### A. Akamai Network Overview

```
Global Network:
├─ 300,000+ servers
├─ 135+ countries
├─ 1,600+ networks
└─ 4,100+ points of presence

Traffic Handled:
├─ 15-30% of global web traffic
├─ 100+ Tbps capacity
└─ Trillions of requests per day

Key Regions:
├─ North America: 40+ clusters
├─ Europe: 35+ clusters
├─ Asia Pacific: 30+ clusters
├─ Latin America: 10+ clusters
├─ Middle East/Africa: 8+ clusters
└─ Oceania: 5+ clusters
```

### B. Glossary

| Term | Definition |
|------|------------|
| **PAPI** | Property Manager API - manages delivery configurations |
| **CPS** | Certificate Provisioning System - manages SSL certificates |
| **AppSec** | Application Security - WAF, DDoS, bot management |
| **Edge Server** | Akamai cache server closest to users |
| **Edge Hostname** | CNAME target for CDN access (*.edgekey.net) |
| **Property** | Configuration container for delivery settings |
| **Rule Tree** | Hierarchical rules defining behaviors |
| **Behavior** | Action to perform (cache, redirect, modify, etc.) |
| **Origin** | Your actual server where content lives |
| **Enrollment** | SSL certificate configuration in CPS |
| **CN** | Common Name - primary domain in certificate |
| **SAN** | Subject Alternative Name - additional domains |
| **DV** | Domain Validation - basic SSL validation |
| **OV** | Organization Validation - business verified SSL |
| **EV** | Extended Validation - highest trust SSL |
| **WAF** | Web Application Firewall - blocks attacks |
| **KRS** | Kona Rule Set - Akamai's WAF rules |
| **Match Target** | Defines which requests get security policies |
| **TTL** | Time To Live - cache duration |
| **CP Code** | Content Provider Code - billing/reporting identifier |

### C. API Rate Limits

| API | Rate Limit | Burst |
|-----|------------|-------|
| **PAPI** | 1000 req/min | 50 req/sec |
| **CPS** | 100 req/min | 10 req/sec |
| **AppSec** | 500 req/min | 25 req/sec |

**Note**: SDK handles rate limiting automatically.

### D. Support Contacts

| Issue Type | Contact |
|------------|---------|
| **Account** | Account Manager |
| **Technical** | support@akamai.com |
| **Billing** | billing@akamai.com |
| **Security Incident** | security@akamai.com |
| **Emergency** | 24/7 NOC Hotline |

---

**Document Version**: 1.0  
**Last Updated**: 2025-12-08  
**Maintainer**: Development Team  

---

