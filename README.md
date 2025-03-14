<p align="center">
  <img align="center" src="./imgs/aws-root-manager.jpg" width="13%" height="13%">
</p>

# AWS Root Manager

A CLI tool for easily manage [AWS Centralized Root Access](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_root-enable-root-access.html).

## Overview

This tool enables AWS Organization administrators to manage centralized root access, allowing you to:
- ✅ Check if Centralized Root Access is enabled in your AWS Organization.
- 🔒 Enable Centralized Root Access for better security and control.
- 📊 Audit root access status across all organization accounts.
- 🗑️ Delete root credentials to enforce security best practices.

## Features

- **Audit**: Get a detailed view of available root credentials in your organization member accounts.
- **Delete**: Remove root credentials with options for:
  - Login profiles.
  - Access keys.
  - MFA devices.
  - Signing certificates.
- **Enable**: Enable centralized root access.
- **Check**: Verify centralized root access settings.

Something missing? Open us a [feature request](https://github.com/unicrons/aws-root-manager/issues/new/choose)!

## Requirements

- AWS Organization management account access.
- AWS CLI configured with appropriate credentials.
- The following IAM permissions:
  ```
  iam:ListOrganizationsFeatures
  organizations:DescribeOrganization
  organizations:ListAccounts
  sts:AssumeRoot
  ```

  Additionally, if the centralized root access feature is not enabled, the following permissions are required to enable it:
  ```
  iam:EnableOrganizationsRootCredentialsManagement
  iam:EnableOrganizationsRootSessions (required only when working with resource policies)
  organizations:EnableAwsServiceAccess
  ```

## How to use it

### Installation

Download the latest version from [GitHub Releases](https://github.com/unicrons/aws-root-manager/releases), or build it from source:
```bash
git clone https://github.com/unicrons/aws-root-manager.git
cd aws-root-manager
go build
```

```bash
./aws-root-manager --help
```
```
Usage:
  aws-root-manager [command]

Available Commands:
  audit       Retrieve root credentials
  check       Check if centralized root access is enabled.
  completion  Generate the autocompletion script for the specified shell
  delete      Delete root credentials
  enable      Enable centralized root access
  help        Help about any command
```

## Examples

Get available root credentials for all member accounts in your AWS Organizations:
```bash
./aws-root-manager audit --accounts all
```

Get available root access keys for accounts `111111111111` and `222222222222`:
```bash
./aws-root-manager audit --accounts 111111111111,222222222222
```

Check if centralized root access is enabled:
```bash
./aws-root-manager check
```


Delete all organization member accounts root credentials:
```bash
./aws-root-manager delete all --accounts all
```

Delete root login profile for account `123456789012`:
```bash
./aws-root-manager delete login --accounts 123456789012
```

Enable centralized root access:
```bash
./aws-root-manager enable
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## References

- [AWS Organizations Documentation](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_introduction.html)
- [AWS Centralized Root Access Documentation](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_root-enable-root-access.html)

---

Made with ❤️ by unicrons.cloud 🦄
