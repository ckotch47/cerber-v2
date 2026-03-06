package recon

import (
	"sort"
	"strings"
	"sync"

	"cerber/internal/dns"
)

type SubdomainScanner struct {
	resolver func(string) bool
}

func NewSubdomainScanner(resolver func(string) bool) *SubdomainScanner {
	if resolver == nil {
		resolver = DNSResolver
	}
	return &SubdomainScanner{resolver: resolver}
}

func DNSResolver(host string) bool {
	return len(dns.CheckDomain(host)) > 0
}

func (s *SubdomainScanner) Collect(
	rootDomain string,
	wordlist []string,
	recurse bool,
	maxDepth int,
	concurrency int,
) []string {
	if maxDepth < 0 {
		maxDepth = 0
	}
	if concurrency <= 0 {
		concurrency = 1
	}

	visited := make(map[string]struct{})
	foundSet := make(map[string]struct{})
	found := make([]string, 0)
	semaphore := make(chan struct{}, concurrency)

	var scan func(domain string, depth int)
	scan = func(domain string, depth int) {
		if depth > maxDepth {
			return
		}

		candidates := make([]string, 0, len(wordlist))
		for _, prefix := range wordlist {
			prefix = strings.TrimSpace(prefix)
			if prefix == "" {
				continue
			}
			fqdn := prefix + "." + domain
			if _, ok := visited[fqdn]; ok {
				continue
			}
			visited[fqdn] = struct{}{}
			candidates = append(candidates, fqdn)
		}

		levelFound := make([]string, 0)
		var wg sync.WaitGroup
		results := make(chan string, len(candidates))

		for _, fqdn := range candidates {
			wg.Add(1)
			go func(candidate string) {
				defer wg.Done()
				semaphore <- struct{}{}
				ok := s.resolver(candidate)
				<-semaphore
				if ok {
					results <- candidate
				}
			}(fqdn)
		}

		wg.Wait()
		close(results)

		for fqdn := range results {
			if _, ok := foundSet[fqdn]; ok {
				continue
			}
			foundSet[fqdn] = struct{}{}
			found = append(found, fqdn)
			levelFound = append(levelFound, fqdn)
		}

		if !recurse || depth >= maxDepth {
			return
		}
		for _, sub := range levelFound {
			scan(sub, depth+1)
		}
	}

	scan(rootDomain, 0)
	sort.Strings(found)
	return found
}
