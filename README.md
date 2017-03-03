# mailck
golang library for email validation

[![Build Status](https://api.travis-ci.org/smancke/mailck.svg?branch=master)](https://travis-ci.org/smancke/mailck)
[![Go Report Card](https://goreportcard.com/badge/github.com/smancke/mailck)](https://goreportcard.com/report/github.com/smancke/mailck)
[![Coverage Status](https://coveralls.io/repos/github/smancke/mailck/badge.svg?branch=master)](https://coveralls.io/github/smancke/mailck?branch=master)


This library allows you to check if an email address realy valid,
by connecting to the mailserver.

## Preconditions
Make sure, that the ip address you are calling from is not
black listed. This is e.g. the case if the ip is a dynamic IP
or does not have a valid SPF record entry for the from domain.
Also make sure, that you have a correct reverse dns lookup for
your ip address.

In case of a blacklisting, the target mailserver may respond with an `SMTP 554`
or just let you run into a timout.

