# go-sdk

[![GoDoc](https://godoc.org/github.com/hexonet/go-sdk?status.svg)](https://godoc.org/github.com/hexonet/go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/hexonet/go-sdk)](https://goreportcard.com/report/github.com/hexonet/go-sdk)
[![cover.run](https://cover.run/go/github.com/hexonet/go-sdk.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=github.com%2Fhexonet%2Fgo-sdk)

This module is a connector library for the insanely fast HEXONET Backend API. For further informations visit our [homepage](http://hexonet.net) and do not hesitate to contact us.

## Requirements

Installed GO on OS-side as described [here](https://golang.org/doc/install). Restart your machine after installing GO.
For developers: Visual Studio Code with installed plugin for Go Development described [here](https://code.visualstudio.com/docs/languages/go).
VS Studio Code will ask you to install some plugins when you start developing a .go file e.g.: gopkgs, goreturns, gocode. Just confirm!

## Getting Started

Clone the git repository into your standard go folder structure by  `go get github.com/hexonet/go-sdk`.
We have also a demo app available showing how to integrate and use our SDK. See [here](https://github.com/hexonet/go-sdk-demo).

### For development purposes

Now you can already start working on the project.

### How to use this module in your project

Use [govendor](https://github.com/kardianos/govendor) for the dependency installation: `govendor fetch github.com/hexonet/go-sdk@<tag id>`. You can update this dependency later on by `govendor sync github.com/hexonet/go-sdk@<new tag id>`.
The dependencies will be installed in subfolder "vendor". Import the module in your project as shown in the examples below.

For more details on govendor, please read the [CheatSheet](https://github.com/kardianos/govendor/wiki/Govendor-CheatSheet) and also the [developer guide](https://github.com/kardianos/govendor/blob/master/doc/dev-guide.md). Knowing about the latter one is very important.

## Development

## Run Tests

Go to subfolder "test" and run `go test`.

### Release an Update

Simply make a PR / merge request.

## Contributing

Please read [CONTRIBUTING.md](https://github.com/hexonet/go-sdk/blob/master/CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/hexonet/go-sdk/tags).

## Authors

* **Kai Schwarz** - *lead development* - [PapaKai](https://github.com/papakai)

See also the list of [contributors](https://github/hexonet/go-sdk/graphs/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## How-to-use Examples

### Session based API Communication

```go
package main

import (
    "github.com/hexonet/go-sdk/client"
    "fmt"
)

func main() {
    cl := client.NewClient()
    cl.SetCredentials("test.user", "test.passw0rd", "")//username, password, otp code (2FA)
    cl.UseOTESystem()
    r := cl.Login()
    if r.IsSuccess() {
        fmt.Println("Login succeeded.")
        cmd := map[string]string{
            "COMMAND": "StatusAccount",
        }
        r = cl.Request(cmd)
        if r.IsSuccess() {
            fmt.Println("Command succeeded.")
            r = cl.Logout()
            if r.IsSuccess() {
                fmt.Println("Logout succeeded.")
            } else {
                fmt.Println("Logout failed.")
            }
        } else {
            fmt.Println("Command failed.")
        }
    } else {
        fmt.Println("Login failed.")
    }
}
```

### Sessionless API Communication

```go
    package main

import (
    "github.com/hexonet/go-sdk/client"
    "fmt"
)

func main() {
    cl := client.NewClient()
    cl.SetCredentials("test.user", "test.passw0rd", "")
    cl.UseOTESystem()
    cmd := map[string]string{
        "COMMAND": "StatusAccount",
    }
    r := cl.Request(cmd)
    if r.IsSuccess() {
        fmt.Println("Command succeeded.")
    } else {
        fmt.Println("Command failed.")
    }
}
```

## Documentation

Run `godoc -http=:6060` on command line and access it via [http://localhost:6060](http://localhost:6060).
Navigate to "packages" > "apiconnector".

Alternative: See our SDK @ [GoDoc.org](https://godoc.org/github.com/hexonet/go-sdk)

## Resources

... the above go documentation server that comes also with further informations around GO.
[Learn Go in Y Minutes](https://learnxinyminutes.com/docs/go/)
[An Introduction to Programming in GO](https://www.golang-book.com/books/intro)
[Go by Example](https://gobyexample.com/)
