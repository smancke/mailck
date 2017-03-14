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
result, err := mailck.Check("noreply@mancke.net", "foo@example.com")

switch result.Result {

  case mailck.Valid:
    // valid!
    // the mailserver accepts mails for this mailbox.

  case mailck.Invalid:
    // invalid, e.g. bacause
    // the reason is contained in result.ResultDetail

  case mailck.Error:
  // something went wrong in the smtp communication
  // we can't say for sure if the address is valid or not


}
```

## License

MIT Licensed
