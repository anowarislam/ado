# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. We appreciate your efforts to responsibly disclose your findings.

### How to Report

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please use **GitHub Security Advisories**:

1. Go to the [Security tab](https://github.com/anowarislam/ado/security/advisories) of this repository
2. Click "Report a vulnerability"
3. Fill out the form with details about the vulnerability

This is the fastest way to reach us and allows for secure, private communication about the issue.

### What to Include

Please include the following information in your report:

- **Type of vulnerability** (e.g., buffer overflow, command injection, etc.)
- **Full paths of source file(s)** related to the vulnerability
- **Location of the affected source code** (tag/branch/commit or direct URL)
- **Step-by-step instructions** to reproduce the issue
- **Proof-of-concept or exploit code** (if possible)
- **Impact assessment** of the vulnerability
- **Suggested fix** (if you have one)

### Response Timeline

- **Initial Response**: Within 48 hours of receiving your report
- **Status Update**: Within 7 days with our assessment
- **Resolution Target**: Within 90 days for confirmed vulnerabilities

### What to Expect

1. **Acknowledgment**: We will acknowledge receipt of your report within 48 hours
2. **Assessment**: We will investigate and assess the severity of the issue
3. **Updates**: We will keep you informed of our progress
4. **Resolution**: Once fixed, we will notify you before public disclosure
5. **Credit**: With your permission, we will credit you in the security advisory

### Safe Harbor

We consider security research conducted in accordance with this policy to be:

- Authorized concerning any applicable anti-hacking laws
- Authorized concerning any relevant anti-circumvention laws
- Exempt from restrictions in our Terms of Service that would interfere with conducting security research

We will not pursue civil action or initiate a complaint to law enforcement for accidental, good-faith violations of this policy.

### Scope

This security policy applies to the `ado` CLI tool and its associated code in this repository.

**In Scope:**
- The `ado` binary and its functionality
- Go source code in `cmd/` and `internal/`
- Build and release infrastructure (CI/CD pipelines)
- Container images published to ghcr.io

**Out of Scope:**
- Python lab code (`lab/py/`) - experimental, not shipped to users
- Third-party dependencies (report to upstream maintainers)
- Social engineering attacks
- Physical attacks

## Security Best Practices for Users

### Verifying Downloads

All releases include SHA256 checksums. Verify your download:

```bash
# Download the checksums file
curl -LO https://github.com/anowarislam/ado/releases/download/vX.Y.Z/checksums.txt

# Verify your downloaded archive
sha256sum -c checksums.txt --ignore-missing
```

### Container Image Verification

Container images are signed and can be verified:

```bash
# Pull with digest verification
docker pull ghcr.io/anowarislam/ado:vX.Y.Z@sha256:...
```

### Reporting Other Issues

For non-security bugs, please use the [GitHub Issues](https://github.com/anowarislam/ado/issues) page.

## Acknowledgments

We thank the following security researchers for their responsible disclosures:

*No acknowledgments yet. Be the first!*
