# scanner
Simple golang scanner to test domains for data bouncing

## Usage

The tool requires a file of domains and a callback domain (e.g., burp suite collaborator or from [Interactsh](https://app.interactsh.com/)).

```bash
./scanner -f <domains-file> -e <your-callback-domain> [options]
```

#### Example

To scan domains from `domains.txt` and test for `X-Forwarded-For` and `Referer` injection, using `xyz.oast.fun` as the callback domain:
```bash
./scanner -f domains.txt -e xyz.oast.fun -H xff -H ref
```

### Flags

| Flag             | Alias | Description                                                                | Required | Default |
|------------------|-------|----------------------------------------------------------------------------|----------|---------|
| `--file`         | `-f`  | File with domains to scan (one per line).                                  | **Yes**  | N/A     |
| `--exfil-domain` | `-e`  | Your callback domain for DNS resolution.                                   | **Yes**  | N/A     |
| `--header`       | `-H`  | Header alias to inject. Use multiple times.                                | No       | `All`   |
| `--workers`      | `-w`  | Number of concurrent workers.                                              | No       | `50`    |

### Supported Header Aliases
Use these short aliases with the `-H` flag:

`host`, `origin`, `xff`, `xwp`, `contact`, `rip`, `tcip`, `xclip`, `fwd`, `xoip`, `clip`, `ref`, `from`.

Run `./scanner --help` for a full mapping of aliases to header names.

