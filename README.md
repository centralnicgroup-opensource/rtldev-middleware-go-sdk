# GO SDK Library - CNIC - 3rd Party Software Integrations

[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release)
[![Build Status](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/workflows/Release/badge.svg?branch=master)](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/workflows/Release/badge.svg?branch=master)
[![GoDoc](https://godoc.org/github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk?status.svg)](https://godoc.org/github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk)](https://goreportcard.com/report/github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![PRs welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/blob/master/CONTRIBUTING.md)

This module is a connector library for the insanely fast CentralNic Reseller Backend API. For further informations visit our [homepage](https://www.centralnicreseller.com) and do not hesitate to [contact us](https://www.centralnicreseller.com/contact).

## Deprecation Notice: Hexonet Go SDK

This SDK succeeds the deprecated Hexonet Go SDK. It is an enhanced version that builds upon the foundation laid by the Hexonet SDK, offering improved features and performance. Hexonet is migrating to CentralNic Reseller, ensuring continued support and development under the new branding.

## Resources

- [Documentation](https://support.centralnicreseller.com/hc/en-gb/articles/5714403954333-Self-Development-Kit-for-Go-Golang)

## Release Notes

For detailed release notes, please visit the [Release Notes](https://github.com/centralnicgroup-opensource/rtldev-middleware-go-sdk/releases) page.

## Running the Demo Application

To run the demo application, follow these steps:

1. **Set Your Credentials**: Ensure your credentials are available. You can either:
  - Replace them directly in the application file.
  - Set the environment variables `CNR_TEST_USER` and `CNR_TEST_PASSWORD` in your terminal.

2. **Execute the Demo**: Once your credentials are set, run the following command in the terminal:

   ```sh
   npm run test-demo
   ```

3. **Update Demo Contents**: If you need to update the contents of the demo file, you can find it at:

   ```plaintext
   demo/demo.go
   ```

By following these steps, you can successfully run and update the demo application.

## Authors

- **Kai Schwarz** - _lead development_ - [KaiSchwarz-cnic](https://github.com/kaischwarz-cnic)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
