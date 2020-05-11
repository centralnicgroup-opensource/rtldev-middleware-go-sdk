# go-sdk

[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release)
[![Build Status](https://travis-ci.com/hexonet/go-sdk.svg?branch=master)](https://travis-ci.com/hexonet/go-sdk)
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

## Features

* Automatic IDN Domain name conversion to punycode (our API accepts only punycode format in commands)
* Allows nested associative arrays in API commands to improve for bulk parameters
* Connecting and communication with our API
* Several ways to access and deal with response data
* Getting the command again returned together with the response
* Sessionless communication
* Session based communication
* Possibility to save API session identifier in session
* Configure a Proxy for API communication
* Configure a Referer for API communication
* High Performance Proxy Setup

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

### High Performance Proxy Setup

Long distances to our main data center in Germany may result in high network latencies. If you encounter such problems, we highly recommend to use this setup, as it uses persistent connections to our API server and the overhead for connection establishments is omitted.

#### Step 1: Required Apache2 packages / modules

*At least Apache version 2.2.9* is required.

The following Apache2 modules must be installed and activated:

```bash
proxy.conf
proxy.load
proxy_http.load
ssl.conf # for HTTPs connection to our API server
ssl.load # for HTTPs connection to our API server
```

#### Step 2: Apache configuration

An example Apache configuration with binding to localhost:

```bash
<VirtualHost 127.0.0.1:80>
    ServerAdmin webmaster@localhost
    ServerSignature Off
    SSLProxyEngine on
    ProxyPass /api/call.cgi https://api.ispapi.net/api/call.cgi min=1 max=2
    <Proxy *>
        Order Deny,Allow
        Deny from none
        Allow from all
    </Proxy>
</VirtualHost>
```

After saving your configuration changes please restart the Apache webserver.

#### Step 3: Using this setup

```go
package main

import (
    "fmt"

    CL "github.com/hexonet/go-sdk/apiclient"
)

func main() {
    cl := CL.NewAPIClient()
    //Default Connection Setup would be used otherwise by default
    cl.UseHighPerformanceConnectionSetup()
    //LIVE System would be used otherwise by default
    cl.UseOTESystem()
    cl.SetCredentials("test.user", "test.passw0rd")

    r := cl.Request(map[string]interface{}{
        "COMMAND": "StatusAccount"
    })
}
```

So, what happens in code behind the scenes? We communicate with localhost (so our proxy setup) that passes the requests to the HEXONET API.
Of course we can't activate this setup by default as it is based on Steps 1 and 2. Otherwise connecting to our API wouldn't work.

Just in case the above port or ip address can't be used, use function setURL instead to set a different URL / Port.
`http://127.0.0.1/api/call.cgi` is the default URL for the High Performance Proxy Setup.
e.g. `$cl->setURL("http://127.0.0.1:8765/api/call.cgi");` would change the port. Configure that port also in the Apache Configuration (-> Step 2)!

Don't use `https` for that setup as it leads to slowing things down as of the https `overhead` of securing the connection. In this setup we just connect to localhost, so no direct outgoing network traffic using `http`. The apache configuration finally takes care passing it to `https` for the final communication to the HEXONET API.

### Customize Logging / Outputs

When having the debug mode activated `github.com/hexonet/logger` will be used for doing outputs.
Of course it could be of interest for integrators to look for a way of getting this replaced by a custom mechanism like forwarding things to a 3rd-party software, logging into file or whatever.

```php
package main

import (
    "fmt"

    CL "github.com/hexonet/go-sdk/apiclient"
    LG "github.com/myspace/customlogger"
)

func main() {
    cl := CL.NewAPIClient()
    cl.SetCredentials("test.user", "test.passw0rd")//username, password
    cl.UseOTESystem()//LIVE System would be used otherwise by default
    cl.enableDebugMode()//activate debug outputs / logging
    cl.setCustomLogger(new LG.NewCustomerLogger())//set your custom mechanism for debug outputs/logging

    r := cl.Request(map[string]interface{}{
        "COMMAND": "StatusAccount"
    })
}
```

NOTE: Find an example for a custom logger class implementation in `customlogger/customlogger.go`. If you have questions, feel free to open a github issue. Follow the interface `ILogger` defined in `logger/logger.go`.

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
        r = cl.Request(map[string]interface{}{
            "COMMAND": "StatusAccount"
        })
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
    r := cl.Request(map[string]interface{}{
        "COMMAND": "StatusAccount"
    })
    if r.IsSuccess() {
        fmt.Println("Command succeeded.")
    } else {
        fmt.Println("Command failed.")
    }
}
```

#### Using Bulk Parameters in API Command

Of course, you could do the following:

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
    cl.UseOTESystem()
    r := cl.Request(map[string]interface{}{
        "COMMAND": "QueryDomainOptions",
        "DOMAIN0": "example1.com";
        "DOMAIN1": "example2.com";
    })
    if r.IsSuccess() {
        fmt.Println("Command succeeded.")
    } else {
        fmt.Println("Command failed.")
    }
}
```

but probably better:

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
    cl.UseOTESystem()
    r := cl.Request(map[string]interface{}{
        "COMMAND": "QueryDomainOptions",
        "DOMAIN": []string{
            "example1.com",
            "example2.com"
        }
    })
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
