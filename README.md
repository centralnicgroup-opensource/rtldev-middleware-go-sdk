# ispapi-lib-go

This module is a connector library for the insanely fast 1API Backend API. For further informations visit our [homepage](http://1api.net) and do not hesitate to contact us.

## Requirements

Installed GO on OS-side as described [here](https://golang.org/doc/install). Restart your machine after installing GO.
For developers: Visual Studio Code with installed plugin for Go Development described [here](https://code.visualstudio.com/docs/languages/go).
VS Studio Code will ask you to install some plugins when you start developing a .go file e.g.: gopkgs, goreturns, gocode. Just confirm!

## Getting Started

Clone the git repository by `git clone ssh://git@gitlab.hexonet.net:44447/hexonet-middleware/ispapi-lib-go.git`.

### For development purposes

Now you can already start working on the project.

### How to use this module in your project

Create a copy of our module on your local disk to ensure no updates will come in as it may break because of a new major release coming with breaking changes. GO doesn't support versioning up to now out of the box.
Import the archive in your project as shown in the examples below.

## Development

## Run Tests

Go to subfolder "test" and run `go test`.

### Release an Update

Simply make a PR / merge request.

## Contributing

Please read [CONTRIBUTING.md](https://gitlab.hexonet.net/hexonet-middleware/ispapi-lib-go/blob/master/CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

Our future plan:
We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://gitlab.hexonet.net/hexonet-middleware/ispapi-lib-go/tags).

As GO doesn't support versioning out of the box up to now, we suggest you save a copy of our module locally and use that copy.
That's the way google is also using it internally. We can not ensure that our repository module source code may not change and break (what we would basically call a major release).

## Authors

* **Kai Schwarz** - *lead development* - [PapaKai](https://github.com/papakai)

See also the list of [contributors](https://gitlab.hexonet.net/hexonet-middleware/ispapi-lib-go/graphs/master) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## How-to-use Examples

### Session based API Communication

```go
package main

import (
    "apiconnector/client"
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
    "apiconnector/client"
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

## Resources

... the above go documentation server that comes also with further informations around GO.
[Learn Go in Y Minutes](https://learnxinyminutes.com/docs/go/)
[An Introduction to Programming in GO](https://www.golang-book.com/books/intro)
[Go by Example](https://gobyexample.com/)