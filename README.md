consul-backdown
===============

Pipe-able backup tool for Consul's K/V.

HOW DO I USE IT?
----------------

First, install it with:

    go get github.com/ghostbar/consul-backdown

Or just go to the [releases](https://github.com/ghostbar/consul-backdown/releases)
page and download the binary for your system.

Then, just run it like:

    $GOPATH/bin/consul-backdown backup > my-kv-backup.txt

More help can be found on `$GOPATH/bin/consul-backdown --help`.

HOW IS THIS DIFFERENT TO CONSUL-BACKUP?
---------------------------------------

There's another older tool,
[kailunshi/consul-backup](https://github.com/kailunshi/consul-backup) that may
work pretty well for you. Go and try it. This tool is actually based on that.

So why build this?:

- The code is not very well written, not even correctly formatted with gofmt.
- It's not pipe-able.
- It's not fully configurable.
- Author is not responsive to pull-requests (reason why I did not sent my
  improvements over there in the first time but instead wrote this).
- No binary releases.

AUTHOR AND LICENSE
------------------

© 2016, Jose-Luis Rivas `<me@ghostbar.co>`.

This software is licensed under the MIT terms, a copy of the license can be
found in the `LICENSE` file in this repository.
