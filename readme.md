
# [trojan](https://github.com/rathena/trojan)

A golang package used to execute integration tests and benchmarks against running rathena server components.


## preparation

Start by compiling rathena source.

Add a user to the `login` table; ideally `ragnarok` with a password `ragnarok` and an `account_id` of 2.  _This will allow testing the registration of a second character server to the login server._

Adjust the following settings if you plan to run the benchmarks:

- set `enable_ip_rules` to `no` in `conf/import/packet_conf.txt`
- set `ipban_enable` to `no` in `conf/import/login_conf.txt`
- increase `allowed_regs` to something arbitrarily high (eg. 100000) in `conf/import/login_conf.txt`


## installation

Install [go](https://golang.org/), ideally using [gvm](https://github.com/moovweb/gvm) so you can easily upgrade versions.

You can download this package and all dependencies via:

	go get github.com/rathena/trojan

_This will download and build all dependencies, **but this project does not produce a binary.**_

Next navigate your terminal to `$GOPATH/src/github.com/rathena/trojan` to begin running commands.


## integration

Here is a list of parameters and the default values that can be adjusted to change execution behavior:

- `packetVer` is set to `55`
- `clientVer` is set to `20151029`
- `loginAddress` is set to `127.0.0.1:6900`
- `charAddress` is set to `127.0.0.1:6121`
- `charCredentials` is set to `ragnarok:ragnarok`
- `mapAddress` is set to `127.0.0.1:5121`

A username will be randomly generated per execution, and the password is a hard-coded value of `secret`.

To run the integration tests, use the following command:

	go test -v

_Add flags as needed to configure._


### notes

These tests are still very new and lacking in many ways.

- There is no support for alternative hexed client behavior, including the other supported login commands.
- There is no support for alternative packet formats created by different packet versions (namely 2017).
- There are no tests for map and character servers, yet.


## benchmarks

Use this command to run the benchmarks:

	go test -v -run=X -bench=.

_The same flags can be applied if addresses are different._
