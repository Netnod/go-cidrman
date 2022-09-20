# go-cidrman
Golang CIDR block management module/package.

Inspired by the Python netaddr library:
* https://github.com/drkjam/netaddr/

This is a fork from [EvilSuperstars/go-cidrman](https://github.com/EvilSuperstars/go-cidrman).

## Background

*Note:* This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it
safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH).
The instructions that follow assume a directory in your home directory
outside of the standard GOPATH (i.e `$HOME/git/GitHub/Netnod/go-cidrman/`).

# Build instructions

## Clone repository

Clone repository to `$HOME/git/GitHub/Netnod/go-cidrman/`:

```sh
$ mkdir -p $HOME/git/GitHub/Netnod/; cd $HOME/git/GitHub/Netnod/
$ git clone git@github.com:Netnod/go-cidrman.git
```

## Build

```sh
$ cd $HOME/git/GitHub/Netnod/go-cidrman/
$ make build
```

## Test

```sh
$ cd $HOME/git/GitHub/Netnod/go-cidrman/
$ make test
```

## Clean

Clean the testcache.

```sh
$ cd $HOME/git/GitHub/Netnod/go-cidrman/
$ make clean
```

# Project status and progress

## Findings about the original project

A lot of work was done in the `ipv6-experimental` branch in February and December 2017.
In June 2019, some of the *metafiles* in the `master` branch were created or updated,
without changes to the code itself.

## Where are we now?

IPv6 merge support was added to the `main` branch in March 2022.
Due to the timeline *issue* above and the `experimental` name, some of the code in this fork was
written from scratch or manaually copied from the `experimental` branch.
The `ipv6-additions` branch was created to create a **PR** in the original project, to merge the IPv6 support.

In April 2022, a `merge-experimental` branch was temporarily used to merge the old work in `ipv6-experimental`
together with the IPv6 support developed in this fork.
With IPv6 support now in `main`, the `ipv6-experimental` branch was removed as it's no longer relevant in this fork.

New `RemoveCIDRs` functions added (in June 2022) to remove/exclude CIDR blocks or IP ranges.

New `SubsetCIDRs` functions added (in September 2022) to select specific CIDR blocks or IP ranges to keep.
A bit like the oposite of RemoveCIDRs, `SubsetCIDRs` selects what to keep instead of what to remove.

## The future

For now, all necessary features have been added to support Netnod's internal development needs.
