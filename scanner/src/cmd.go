package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
)

type InjectableHeader struct {
	Name, Format string
	Aliases      []string
}

var allInjectableHeaders = []InjectableHeader{
	{Name: "Host", Format: "host.%s.%s", Aliases: []string{"host"}},
	{Name: "Origin", Format: "https://origin.%s.%s", Aliases: []string{"origin"}},
	{Name: "X-Forwarded-For", Format: "xff.%s.%s", Aliases: []string{"xff", "x-forwarded-for"}},
	{Name: "X-Wap-Profile", Format: "https://wap.%s.%s/wap.xml", Aliases: []string{"xwp", "x-wap-profile"}},
	{Name: "Contact", Format: "root@contact.%s.%s", Aliases: []string{"contact"}},
	{Name: "X-Real-IP", Format: "rip.%s.%s", Aliases: []string{"rip", "x-real-ip"}},
	{Name: "True-Client-IP", Format: "tcip.%s.%s", Aliases: []string{"tcip", "true-client-ip"}},
	{Name: "X-Client-IP", Format: "xclip.%s.%s", Aliases: []string{"xclip", "x-client-ip"}},
	{Name: "Forwarded", Format: "for=fwd.%s.%s;proto=https", Aliases: []string{"fwd", "forwarded"}},
	{Name: "X-Originating-IP", Format: "origip.%s.%s", Aliases: []string{"xoip", "x-originating-ip"}},
	{Name: "Client-IP", Format: "clip.%s.%s", Aliases: []string{"clip", "client-ip"}},
	{Name: "Referer", Format: "https://ref.%s.%s", Aliases: []string{"ref", "referer"}},
	{Name: "From", Format: "root@from.%s.%s", Aliases: []string{"from"}},
}

func NewApp() *cli.Command {
	return &cli.Command{
		Name:        "bounce-scanner",
		Usage:       "Scans domains for databouncing by injecting headers that trigger DNS lookups on a GET request",
		Description: generateHeaderHelp(),
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "file", Aliases: []string{"f"}, Usage: "File with domains to be scanned", Required: true},
			&cli.StringFlag{Name: "exfil-domain", Aliases: []string{"e"}, Usage: "Domain for DNS resolution", Required: true},
			&cli.IntFlag{Name: "workers", Aliases: []string{"w"}, Usage: "Number of concurrent workers", Value: defaultWorkers},
			&cli.StringSliceFlag{
				Name: "header", Aliases: []string{"H"},
				Usage: "Header alias to inject (e.g., 'xff'). Use multiple times. Defaults to all.",
			},
		},
		Action: runScan,
	}
}

func runScan(ctx context.Context, cmd *cli.Command) error {
	domainsFile := cmd.String("file")
	exfilDomain := cmd.String("exfil-domain")
	workers := cmd.Int("workers")
	requestedHeaders := cmd.StringSlice("header")

	domains, err := readDomainsFromFile(domainsFile)
	if err != nil {
		return err
	}

	headersToInject, showOpsecWarning, err := getHeadersToInject(requestedHeaders)
	if err != nil {
		return err
	}

	config := ScannerConfig{
		Domains:          domains,
		ExfilDomain:      exfilDomain,
		HeadersToInject:  headersToInject,
		Workers:          workers,
		ShowOpsecWarning: showOpsecWarning,
	}

	scanner := NewScanner(config)
	PrintInfo("Starting scan...")
	scanner.Run(ctx)

	PrintSuccess("Scan complete.")
	return nil
}

func getHeadersToInject(requestedAliases []string) (map[string]string, bool, error) {
	headersToInject := make(map[string]string)
	showWarning := false

	if len(requestedAliases) == 0 {
		showWarning = true
		for _, h := range allInjectableHeaders {
			headersToInject[h.Name] = h.Format
		}
		return headersToInject, showWarning, nil
	}

	aliasMap := make(map[string]InjectableHeader)
	for _, h := range allInjectableHeaders {
		for _, alias := range h.Aliases {
			aliasMap[alias] = h
		}
	}

	for _, alias := range requestedAliases {
		if headerDef, ok := aliasMap[strings.ToLower(alias)]; ok {
			headersToInject[headerDef.Name] = headerDef.Format
		} else {
			PrintInfo("Unknown header or alias %q requested and will be ignored.", alias)
		}
	}

	if len(headersToInject) == 0 {
		return nil, false, fmt.Errorf("none of the requested headers are supported. Run with --help to see available options")
	}
	return headersToInject, showWarning, nil
}

func generateHeaderHelp() string {
	var builder strings.Builder
	builder.WriteString("\nAvailable headers and their aliases (case-insensitive):\n")
	for _, h := range allInjectableHeaders {
		builder.WriteString(fmt.Sprintf("  - %-20s (Aliases: %s)\n", h.Name, strings.Join(h.Aliases, ", ")))
	}
	return builder.String()
}

const defaultWorkers = 50
