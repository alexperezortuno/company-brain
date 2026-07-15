package health

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Status struct {
	Status       string            `json:"status"`
	Instance     string            `json:"instance"`
	Dependencies map[string]string `json:"dependencies"`
	CheckedAt    time.Time         `json:"checked_at"`
}

type Checker struct {
	timeout         time.Duration
	postgresAddress string
	redisAddress    string
	redisPassword   string
	qdrantURL       string
	minioHealthURL  string
	workerHealthURL string
}

func New(timeout time.Duration, postgresAddress, redisAddress, redisPassword, qdrantURL, minioHealthURL, workerHealthURL string) *Checker {
	return &Checker{
		timeout: timeout, postgresAddress: postgresAddress, redisAddress: redisAddress,
		redisPassword: redisPassword, qdrantURL: qdrantURL,
		minioHealthURL: minioHealthURL, workerHealthURL: workerHealthURL,
	}
}

func (c *Checker) Check(ctx context.Context, instance string) Status {
	checks := map[string]func(context.Context) error{
		"postgres":         func(ctx context.Context) error { return tcpCheck(ctx, c.postgresAddress) },
		"redis":            func(ctx context.Context) error { return redisCheck(ctx, c.redisAddress, c.redisPassword) },
		"qdrant":           func(ctx context.Context) error { return httpCheck(ctx, c.qdrantURL) },
		"minio":            func(ctx context.Context) error { return httpCheck(ctx, c.minioHealthURL) },
		"knowledge_worker": func(ctx context.Context) error { return httpCheck(ctx, c.workerHealthURL) },
	}

	deps := make(map[string]string, len(checks))
	var mu sync.Mutex
	var wg sync.WaitGroup
	for name, check := range checks {
		name, check := name, check
		wg.Add(1)
		go func() {
			defer wg.Done()
			child, cancel := context.WithTimeout(ctx, c.timeout)
			defer cancel()
			result := "healthy"
			if err := check(child); err != nil {
				result = "unhealthy: " + err.Error()
			}
			mu.Lock()
			deps[name] = result
			mu.Unlock()
		}()
	}
	wg.Wait()

	overall := "healthy"
	for _, result := range deps {
		if !strings.HasPrefix(result, "healthy") {
			overall = "degraded"
			break
		}
	}

	return Status{Status: overall, Instance: instance, Dependencies: deps, CheckedAt: time.Now().UTC()}
}

func tcpCheck(ctx context.Context, address string) error {
	dialer := net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return err
	}
	return conn.Close()
}

func httpCheck(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 0}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return nil
}

func redisCheck(ctx context.Context, address, password string) error {
	dialer := net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(2 * time.Second))

	if password != "" {
		if _, err := fmt.Fprintf(conn, "*2\r\n$4\r\nAUTH\r\n$%d\r\n%s\r\n", len(password), password); err != nil {
			return err
		}
		line, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil || !strings.HasPrefix(line, "+OK") {
			return fmt.Errorf("redis auth failed")
		}
	}
	if _, err := fmt.Fprint(conn, "*1\r\n$4\r\nPING\r\n"); err != nil {
		return err
	}
	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return err
	}
	if !strings.HasPrefix(line, "+PONG") {
		return fmt.Errorf("unexpected response: %s", strings.TrimSpace(line))
	}
	return nil
}
