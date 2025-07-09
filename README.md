# bounce
Simple implentation of a data bouncing scanner and exfiltrator 

## What it is
Databouncing provides to the attacker the possability of DNS exfiltration even if the DNS provider is controlled by the victim. Through custom request headers the attacker can trigger a dns resolution on the side of the GET request target, this can make defenders missinterpret the request or even ignore it. Deeper research is present at this [url](https://databouncing.io/).

## Tools
This repo provides two tools to better take advantage of databouncing. The first one is a [scanner](/scanner) built in go that allows the operator to scan a list of domains to see if their are valid points of bouncing. The second one is a simple [bouncer](/bouncer) rust program that takes advantage of this feature to bounce data through a valid point of bouncing.
