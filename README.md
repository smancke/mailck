# mailck - SMTP mail validation
golang library for email validation

[![Build Status](https://api.travis-ci.org/smancke/mailck.svg?branch=master)](https://travis-ci.org/smancke/mailck)
[![Go Report Card](https://goreportcard.com/badge/github.com/smancke/mailck)](https://goreportcard.com/report/github.com/smancke/mailck)
[![Coverage Status](https://coveralls.io/repos/github/smancke/mailck/badge.svg?branch=master)](https://coveralls.io/github/smancke/mailck?branch=master)

This library allows you to check if an email address is realy valid:

* Syntax check
* Blacklist of disposable mailservers (e.g. mailinator.com)
* SMTP mailbox check

## Preconditions
Make sure, that the ip address you are calling from is not
black listed. This is e.g. the case if the ip is a dynamic IP.
Also make sure, that you have a correct reverse dns lookup for
your ip address, matching the hostname of your *from* adress.
Alternatively use a SPF DNS record entry matching the host part
of the *from* address.

In case of a blacklisting, the target mailserver may respond with an `SMTP 554`
or just let you run into a timout.

## Usage

[![GoDoc](https://godoc.org/github.com/smancke/mailck?status.png)](https://godoc.org/github.com/smancke/mailck)

Do all checks at once:

```go
result, _ := mailck.Check("noreply@mancke.net", "foo@example.com")
switch {

  case result.IsValid():
  // valid!
  // the mailserver accepts mails for this mailbox.

  case result.IsError():
  // something went wrong in the smtp communication
  // we can't say for sure if the address is valid or not

  case result.IsInvalid():
  // invalid for some reason
  // the reason is contained in result.ResultDetail
  // or we can check for different reasons:
  switch (result) {
    case mailck.InvalidDomain:
    // domain is invalid
    case mailck.InvalidSyntax:
    // e-mail address syntax is invalid
  }
}
```

## License

MIT Licensed
