# Security Policy

## Supported Versions

We actively support the following versions of StudioSpeech with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of StudioSpeech seriously. If you discover a security vulnerability, please follow these steps:

### 1. Do Not Create Public Issues

**Please do not report security vulnerabilities through public GitHub issues.**

### 2. Use GitHub Security Advisories

Report security vulnerabilities using GitHub's Security Advisory feature:

1. Go to the [StudioSpeech repository](https://github.com/g-laliotis/StudioSpeech)
2. Click on the "Security" tab
3. Click "Report a vulnerability"
4. Fill out the advisory form with details

### 3. What to Include

Please include the following information in your report:

- **Description** - A clear description of the vulnerability
- **Impact** - What could an attacker accomplish?
- **Reproduction** - Step-by-step instructions to reproduce the issue
- **Affected Versions** - Which versions are affected?
- **Suggested Fix** - If you have ideas for how to fix the issue

### 4. Response Timeline

We will acknowledge receipt of your vulnerability report within **48 hours** and provide a more detailed response within **7 days** indicating the next steps in handling your report.

After the initial reply, we will keep you informed of the progress towards a fix and may ask for additional information or guidance.

## Security Considerations

### Local Processing
StudioSpeech is designed for local, offline processing to minimize security risks:

- **No Cloud APIs** - All processing happens locally
- **No Data Transmission** - Input files never leave your system
- **No Telemetry** - No usage data is collected or transmitted

### Input Validation
- All input files are validated before processing
- File type restrictions prevent execution of malicious content
- Memory limits prevent resource exhaustion attacks

### Dependencies
- We regularly update dependencies to patch known vulnerabilities
- Minimal dependency footprint reduces attack surface
- All dependencies are vetted for security issues

### File System Access
- StudioSpeech only accesses files explicitly provided by the user
- No automatic file discovery or scanning
- Temporary files are cleaned up after processing

## Best Practices for Users

### Safe Usage
- Only process files from trusted sources
- Keep StudioSpeech updated to the latest version
- Run with minimal required permissions
- Use in isolated environments for sensitive content

### System Security
- Keep your operating system updated
- Use antivirus software
- Regularly backup important data
- Monitor system resources during processing

## Vulnerability Disclosure Policy

### Coordinated Disclosure
We follow a coordinated disclosure process:

1. **Report received** - We acknowledge the report
2. **Investigation** - We investigate and validate the issue
3. **Fix development** - We develop and test a fix
4. **Release** - We release the fix in a new version
5. **Public disclosure** - We publicly disclose the vulnerability after users have had time to update

### Timeline
- **Day 0** - Vulnerability reported
- **Day 1-7** - Initial assessment and acknowledgment
- **Day 7-30** - Investigation and fix development
- **Day 30-90** - Testing and release preparation
- **Day 90+** - Public disclosure (may be extended if needed)

## Security Updates

Security updates will be:
- Released as patch versions (e.g., 1.0.1)
- Clearly marked in release notes
- Announced through GitHub releases
- Documented in the changelog

## Contact

For security-related questions or concerns that are not vulnerabilities, you can contact the maintainers through:

- GitHub Issues (for general security questions)
- GitHub Discussions (for security best practices)

## Acknowledgments

We appreciate the security research community's efforts to improve the security of open source software. Contributors who responsibly disclose security vulnerabilities will be acknowledged in our security advisories (with their permission).

## Legal

This security policy is provided "as is" without warranty of any kind. The maintainers reserve the right to modify this policy at any time.