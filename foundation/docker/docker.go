package docker

import (
	"bytes"
	"encoding/json"
	"net"
	"os/exec"
	"testing"
)

type Container struct {
	ID   string
	Host string
}

func StartContainer(t *testing.T, image string, port string, args ...string) *Container {

	arg := []string{"run", "-P", "-d"}
	arg = append(arg, args...)
	arg = append(arg, image)

	cmd := exec.Command("docker", arg...)

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to start container: %s: %v", image, err)
	}

	id := out.String()[:12]

	cmd = exec.Command("docker", "inspect", id)
	out.Reset()
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to inspect container: %s: %v", id, err)
	}

	var doc []map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatalf("failed to unmarshal json: %v", err)
	}

	ip, randPort := extractIPPort(t, doc, port)

	c := Container{
		ID:   id,
		Host: net.JoinHostPort(ip, randPort),
	}

	t.Logf("Image:			%s", image)
	t.Logf("ContainerID:    %s", id)
	t.Logf("Host:			%s", c.Host)

	return &c
}

func StopContainer(t *testing.T, id string) {

	if err := exec.Command("docker", "stop", id).Run(); err != nil {
		t.Fatalf("failed to stop container: %v", err)
	}

	t.Log("Stopped:", id)
	if err := exec.Command("docker", "rm", id, "-v").Run(); err != nil {
		t.Fatalf("failed to remove container: %v", err)
	}
	t.Log("Removed:", id)

}

/*
"5432/tcp": [
                    {
                        "HostIp": "0.0.0.0",
                        "HostPort": "5432"
                    }
            ],

			Extract this information from the output of the docker inspect command

*/

func extractIPPort(t *testing.T, doc []map[string]interface{}, port string) (string, string) {

	nw, exists := doc[0]["NetworkSettings"]
	if !exists {
		t.Fatal("failed to find network settings")
	}
	ports, exists := nw.(map[string]interface{})["Ports"]
	if !exists {
		t.Fatal("failed to find network ports Settings")
	}
	tcp, exists := ports.(map[string]interface{})[port+"/tcp"]
	if !exists {
		t.Fatal("failed to find network port/tcp settings")
	}

	list, exists := tcp.([]interface{})
	if !exists {
		t.Fatal("failed to find network port/tcp list settings")
	}

	var hostIP string
	var hostPort string

	for _, l := range list {
		data, exists := l.(map[string]interface{})
		if !exists {
			t.Fatal("failed to find network port/tcp list data")
		}
		hostPort = data["HostIP"].(string)
		if hostIP != "::" {
			hostPort = data["HostPort"].(string)
			break // Need to Remove this ??
		}
	}
	return hostIP, hostPort

}

func DumpContainerLogs(t *testing.T, id string) {

	out, err := exec.Command("docker", "logs", id).CombinedOutput()
	if err != nil {
		t.Fatalf("failed to log container: %v", err)
	}

	t.Logf("Logs for %s\n%s:", id, out)
}
