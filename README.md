# mailck
golang library for email validation

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

