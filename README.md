
# GoTOTP: A Simple Time-Based One-Time Password (TOTP) Library

[![Go](https://github.com/MrTuNNe/GoTOTP/actions/workflows/go.yml/badge.svg)](https://github.com/MrTuNNe/GoTOTP/actions/workflows/go.yml)

## Overview

GoTOTP is a simple, stable, and efficient Time-Based One-Time Password (TOTP) library written in Go, built with a focus on simplicity and stability. It aims to be a long-term solution for generating and verifying TOTPs, without reinventing the wheel.

This implementation is based on **[RFC 6238](https://datatracker.ietf.org/doc/html/rfc6238)**, ensuring compatibility with Google Authenticator and similar apps.

### Why GoTOTP?

- **Simplicity**: Prioritizes straightforward code and minimal changes once the library reaches maturity.
- **Stability**: Aims for long-term stability.
- **Standard-compliant**: Utilizes existing Go standard libraries wherever possible.
- **Personal Project**: Initially built out of curiosity and as a break from web development.

## Features

- Generates a 6-digit code, valid for 30 seconds (default).
- Built-in methods for secret generation, code verification, and URI generation for QR codes.
- QR code support for Google Authenticator and other similar apps (future enhancement).

## Project Goals & Roadmap

1. **Simplicity**: Maintain ease of use and clarity in design.
2. **QR Code Generation**: Implement a method to generate QR codes for TOTP setup.
3. **Long-Term Stability**: Keep the library rock-solid for years to come.

## Future Enhancements

- Potentially allow users to configure code length and expiration period.
  - *Note*: This is not a current priority.

## Installation

To get started, install the package:

```bash
go get github.com/MrTuNNe/GoTOTP
```

## Usage Example

Here is a basic example of how to generate and verify TOTPs using GoTOTP:

```go
package main

import (
    "fmt"
    "time"
    "github.com/MrTuNNe/GoTOTP"
)

func main() {
    // Generate a random secret (base32 encoded, without padding)
    secret, err := GoTOTP.GenerateRandomSecret(32) // 32 bytes length
    if err != nil {
        // Handle error
        fmt.Println("Error generating secret:", err)
        return
    }

    // Create a new TOTP instance
    totp := GoTOTP.TOTP{
        Key:      "OK6ZZOALZY6RNZBPM4QKD2ZFO5F3PTP56VIAXLDJLEHBPLJJIZNQ",
        Issuer:   "mrtunne.info",
        UserName: "admin@admin.test",
    }

    // Generate a TOTP based on the current timestamp
    code := totp.GenerateTOTP(time.Now().Unix())
    fmt.Println("Generated TOTP:", code)

    // Verify user input
    if totp.Verify("149425") {
        fmt.Println("Code verified successfully!")
    } else {
        fmt.Println("Invalid code.")
    }

    // Check code with a specific timestamp
    if totp.VerifyWithTimestamp(1723719527, "611626") {
        fmt.Println("Code valid for timestamp.")
    } else {
        fmt.Println("Code invalid for timestamp.")
    }

    // Generate a URI for QR code generation
    uri := totp.GenerateURI()
    fmt.Println("TOTP URI:", uri)
    // Example output:
    // otpauth://totp/mrtunne.info:%20admin@admin.test?algorithm=SHA256&digits=6&issuer=mrtunne.info&period=30&secret=OK6ZZOALZY6RNZBPM4QKD2ZFO5F3PTP56VIAXLDJLEHBPLJJIZNQ
}
```

### Notes:
- **Default Expiration**: Generated codes expire after 30 seconds.
- **Default Length**: Codes are 6 digits long.

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request if you have any improvements or suggestions.

## License

This project is licensed under the [MIT License](LICENSE).
