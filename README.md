# crechebot

The wonderful creche service at our gym will watch your children, for up to 2
hours per day, while you use the gym facilities. These creche spots must be 
booked in advance, via a clumsy web UI, and become available for booking 2
weeks in advance, at 10pm local time.

But, there is a catch...

There are, of course, more parents attempting to book these creche spots than
there are spots available. In order to improve one's chance of getting a spot,
you need to make the booking the moment they become available.

Before crechebot, every weeknight I had to:
- stop whatever I was doing at 9:55pm
- log into the webUI
- locate the spot I want to book
- wait until exactly 10pm
- click the "confirm booking" button in the UI
- hope that I got a spot

This daily ritual was less than ideal, so I made crechebot to do it for me.

Crechebot is a small Go program that makes the creche booking for me. The
booking service has no public API, so the creche bot mimics the set of requests
that the web UI sends when a real user makes a booking.

## Usage

	creche [-under2] [-early] USERNAME PASSWORD

		USERNAME, username for the web booking platform
			PASSWORD, password for the web booking platform
		-under2, book the slot for a child under 2 (default is age 2+)
		-early, book the slot for 9-11am (default is 11am-12pm)

## Full Automation

In order to fully replace my ritual, I run this program every evening via cron:

	56 21 * * 1,3,5 /some/path/bin/creche myuser mypass >> $HOME/logs/crechelogMe.txt 2>&1
	56 21 * * 2,4 /some/path/bin/creche -early myuser mypass >> $HOME/logs/crechelogMe.txt 2>&1
	56 21 * * 1,3,5 /some/path/bin/creche otheruser otherpass >> $HOME/logs/crechelogPartner.txt 2>&1
	56 21 * * 2,4 /some/path/bin/creche -early - under2 otheruser otherpass >> $HOME/logs/crechelogPartner.txt 2>&1

This manages the daily booking of creche for both of our children, without any
further intervention from us.
