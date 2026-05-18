
# GoTOTP: A Simple Time-Based One-Time Password (TOTP) Library

[![Go](https://github.com/MrTuNNe/GoTOTP/actions/workflows/go.yml/badge.svg)](https://github.com/MrTuNNe/GoTOTP/actions/workflows/go.yml)

## Overview

GoTOTP is a simple, stable, and efficient Time-Based One-Time Password (TOTP) library written in Go, built with a focus on simplicity and stability. It aims to be a long-term solution for generating and verifying TOTPs, without reinventing the wheel.

This implementation is based on **[RFC 6238](https://datatracker.ietf.org/doc/html/rfc6238)**, ensuring compatibility with Google Authenticator and similar apps.

> **Note**: Most authenticator apps (Google Authenticator, Authy, etc.) default to SHA-1. Set `Algorithm: GoTOTP.SHA1` on your `TOTP` struct when targeting those apps. Omitting the field defaults to SHA-256.

### Why GoTOTP?

- **Simplicity**: Prioritizes straightforward code and minimal changes once the library reaches maturity.
- **Stability**: Aims for long-term stability.
- **Standard-compliant**: Utilizes existing Go standard libraries wherever possible.
- **Personal Project**: Initially built out of curiosity and as a break from web development.

## Features

- Generates a 6-digit code, valid for 30 seconds.
- Supports **SHA-1**, **SHA-256**, and **SHA-512** HMAC algorithms via the `Algorithm` field.
- Built-in methods for secret generation, code verification, and URI generation for QR codes.

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
    secret, err := GoTOTP.GenerateRandomSecret(32)
    if err != nil {
        fmt.Println("Error generating secret:", err)
        return
    }

    // Create a new TOTP instance (Algorithm defaults to SHA-256 when omitted)
    totp := GoTOTP.TOTP{
        Key:      secret,
        Issuer:   "mrtunne.info",
        UserName: "admin@admin.test",
    }

    // Generate a TOTP based on the current timestamp
    code, err := totp.GenerateTOTP(time.Now().Unix())
    if err != nil {
        fmt.Println("Error generating TOTP:", err)
        return
    }
    fmt.Println("Generated TOTP:", code)

    // Verify user input against the current timestamp
    if totp.Verify(code) {
        fmt.Println("Code verified successfully!")
    } else {
        fmt.Println("Invalid code.")
    }

    // Generate a URI suitable for QR code scanners
    uri, err := totp.GenerateURI()
    if err != nil {
        fmt.Println("Error generating URI:", err)
        return
    }
    fmt.Println("TOTP URI:", uri)
}
```

### Choosing an algorithm

```go
// SHA-1 — required for Google Authenticator and most mobile apps
totp := GoTOTP.TOTP{
    Key:       secret,
    Algorithm: GoTOTP.SHA1,
}

// SHA-512 — for services that explicitly require it
totp := GoTOTP.TOTP{
    Key:       secret,
    Algorithm: GoTOTP.SHA512,
}
```

### Notes:
- **Default algorithm**: SHA-256 (set `Algorithm: GoTOTP.SHA1` for most authenticator apps).
- **Default expiration**: Generated codes expire after 30 seconds.
- **Default length**: Codes are 6 digits long.

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request if you have any improvements or suggestions.

## License

This project is licensed under the [MIT License](LICENSE).
