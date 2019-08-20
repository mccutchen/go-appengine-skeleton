package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

const (
	testProjectID = "testproject"
	testDBPort    = "8081"
	testDBImage   = "google/cloud-sdk"
)

func startTestDb(ctx context.Context) (string, func(), error) {
	containerName := fmt.Sprintf("%s-testdb-%d", testProjectID, time.Now().UnixNano())

	cmd := exec.CommandContext(
		ctx,
		"docker", "run", "-d", "--rm",
		"--name", containerName,
		"-p", testDBPort,
		testDBImage,
		"gcloud", "beta", "emulators", "datastore", "start", "--project", testProjectID, "--no-store-on-disk", "--host-port", "0.0.0.0:"+testDBPort,
	)

	out, err := cmd.Output()
	if err != nil {
		return "", nil, fmt.Errorf("docker failure: %s", err)
	}

	cid := strings.TrimSpace(string(out))

	stop := func() {
		log.Printf("stopping testdb container %s", cid)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "docker", "rm", "-f", cid)
		if err := cmd.Run(); err != nil {
			log.Panicf("failed to stop container %q: %s", cid, err)
		}
	}

	addr, err := getTestDBAddr(ctx, cid)
	if err != nil {
		stop()
		return "", nil, fmt.Errorf("failed to get testdb addr: %s", err)
	}

	maxWait := 60 * time.Second
	log.Printf("waiting up to %s for %s (%s) to come online ...", maxWait, addr, cid)
	if err := waitForServer(ctx, maxWait, fmt.Sprintf("http://%s", addr)); err != nil {
		stop()
		return "", nil, fmt.Errorf("timed out waiting for testdb to come online: %s", err)
	}

	log.Printf("testdb cid=%q host=%q", cid, addr)

	return addr, stop, nil
}

func getTestDBAddr(ctx context.Context, cid string) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", "port", cid)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// out = "8081/tcp -> 0.0.0.0:32768"
	parts := strings.Split(string(out), " ")
	if len(parts) != 3 {
		return "", fmt.Errorf("could not parse `docker port` output: %v", parts)
	}

	return strings.TrimSpace(parts[2]), nil
}

func waitForServer(ctx context.Context, maxWait time.Duration, url string) error {
	stepDelay := 250 * time.Millisecond

	ctx, cancel := context.WithTimeout(ctx, maxWait)
	defer cancel()

	client := &http.Client{
		Timeout: 250 * time.Millisecond,
	}

	for {
		select {
		case <-time.After(stepDelay):
			resp, err := client.Get(url)
			if err != nil {
				continue
			}
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
