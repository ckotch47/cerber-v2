package command

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"cerber/internal/style"
	"cerber/internal/utils"
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Выполняет поиск поддоменов по списку из файла",
	Run:   FindHost,
}

var commandBruteForce utils.BruteForceType

func init() {
	findCmd.Flags().StringVarP(
		&commandBruteForce.WorldList,
		"wordlist",
		"w",
		"",
		"Файл со списком",
	)
	findCmd.Flags().StringVar(
		&commandBruteForce.WorldList,
		"worldlis",
		"",
		"Устаревший алиас для --wordlist",
	)
	_ = findCmd.Flags().MarkDeprecated("worldlis", "use --wordlist instead")
	findCmd.Flags().BoolVarP(
		&commandBruteForce.Recurse,
		"recurse",
		"r",
		false,
		"Включить рекурсию для брутфорса",
	)
	findCmd.Flags().IntVar(
		&commandBruteForce.MaxDepth,
		"max-depth",
		2,
		"Максимальная глубина рекурсии для поиска поддоменов",
	)
	findCmd.Flags().IntVarP(
		&commandBruteForce.Concurrency,
		"concurrency",
		"c",
		20,
		"Количество параллельных DNS-запросов",
	)
}

func FindHost(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println(style.NotFoundStyle.Render("Не указан домен"))
		return
	}
	if commandBruteForce.WorldList == "" {
		fmt.Println(style.NotFoundStyle.Render("Файл со списком не найден"))
		return
	}

	domain := cleanDomain(args[0])
	domainList := utils.ReadFile(commandBruteForce.WorldList)
	if len(domainList) == 0 {
		fmt.Println(style.NotFoundStyle.Render("Файл со списком пустой или не удалось прочитать"))
		return
	}

	found := collectSubDomains(
		domain,
		domainList,
		commandBruteForce.Recurse,
		commandBruteForce.MaxDepth,
		commandBruteForce.Concurrency,
		hostExists,
	)
	for _, subdomain := range found {
		fmt.Println(style.SuccessStyle.Render(subdomain))
	}
}

func collectSubDomains(
	rootDomain string,
	wordlist []string,
	recurse bool,
	maxDepth int,
	concurrency int,
	resolver func(string) bool,
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
				ok := resolver(candidate)
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

var hostExists = func(host string) bool {
	return len(lookupHost(host)) > 0
}
