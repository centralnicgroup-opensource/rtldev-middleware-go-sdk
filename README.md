# go-sdk

[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release)
[![Build Status](https://travis-ci.org/hexonet/go-sdk.svg?branch=master)](https://travis-ci.org/hexonet/go-sdk)
[![GoDoc](https://godoc.org/github.com/hexonet/go-sdk?status.svg)](https://godoc.org/github.com/hexonet/go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/hexonet/go-sdk)](https://goreportcard.com/report/github.com/hexonet/go-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![PRs welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/hexonet/go-sdk/blob/master/CONTRIBUTING.md)

This module is a connector library for the insanely fast HEXONET Backend API. For further informations visit our [homepage](http://hexonet.net) and do not hesitate to [contact us](https://www.hexonet.net/contact).

## Resources

* [Usage Guide](https://github.com/hexonet/go-sdk/blob/master/README.md#how-to-use-this-module-in-your-project)
* [SDK Documenation](https://godoc.org/github.com/hexonet/go-sdk)
* [HEXONET Backend API Documentation](https://github.com/hexonet/hexonet-api-documentation/tree/master/API)
* [Release Notes](https://github.com/hexonet/go-sdk/releases)
* [Development Guide](https://github.com/hexonet/go-sdk/wiki/Development-Guide)

## How to use this module in your project

We have also a demo app available showing how to integrate and use our SDK. See [here](https://github.com/hexonet/go-sdk-demo).

### Requirements

* Installed [GO/GOLANG](https://golang.org/doc/install). Restart your machine after installing GO.
* Installed [govendor](https://github.com/kardianos/govendor).

NOTE: Make sure you add the go binary path to your PATH environment variable. Add the below lines for a standard installation into your profile configuration file (~/.profile).

```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

Then reload the profile configuration by `source ~/.profile`.

### Using govendor

Use [govendor](https://github.com/kardianos/govendor) for the dependency installation by `govendor fetch -tree github.com/hexonet/go-sdk@<tag id>` where *tag id* corresponds to a [release version tag](https://github.com/hexonet/go-sdk/releases). You can update this dependency later on by `govendor sync github.com/hexonet/go-sdk@<new tag id>`. The dependencies will be installed in your project's subfolder "vendor". Import the module in your project as shown in the examples below.

For more details on govendor, please read the [CheatSheet](https://github.com/kardianos/govendor/wiki/Govendor-CheatSheet) and also the [developer guide](https://github.com/kardianos/govendor/blob/master/doc/dev-guide.md).

### Usage Examples

Please have an eye on our [HEXONET Backend API documentation](https://github.com/hexonet/hexonet-api-documentation/tree/master/API). Here you can find information on available Commands and their response data.

#### Session based API Communication

```go
package main

import (
    "fmt"
    
    CL "github.com/hexonet/go-sdk/apiclient"
)

func main() {
    cl := CL.NewAPIClient()
    cl.SetCredentials("test.user", "test.passw0rd")//username, password
    // or cl.SetRoleCredentials("test.user", "testrole", "test.passw0rd")
    // for role user credentials
    cl.UseOTESystem()

    // use this to provide your outgoing ip address for api communication
    // to be used in case you have ip filter settings active
    cl.SetRemoteIPAddress("1.2.3.4");

    // cl.EnableDebugMode() // to activate debug outputs of the API communication
    r := cl.Login()
    // or r := cl.Login("12345678") // provide here your 2FA otp code
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

#### Sessionless API Communication

```go
package main

import (
    "fmt"
    
    CL "github.com/hexonet/go-sdk/apiclient"
)

func main() {
    cl := CL.NewAPIClient()
    cl.SetCredentials("test.user", "test.passw0rd")
    cl.SetRemoteIPAddress("1.2.3.4")
    //cl.SetOTP("12345678") to provide your 2FA otp code
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

## Contributing

Please read [our development guide](https://github.com/hexonet/go-sdk/wiki/Development-Guide) for details on our code of conduct, and the process for submitting pull requests to us.

## Authors

* **Kai Schwarz** - *lead development* - [PapaKai](https://github.com/papakai)

See also the list of [contributors](https://github.com/hexonet/go-sdk/graphs/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
