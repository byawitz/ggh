package history

import (
	"testing"
)

var historyFile = `
[
  {
    "time": "2024-01-01 00:00:00 +0000 UTC",
    "connection": {
      "name": "stage",
      "host": "host.name",
      "key": "~/.ssh/id_rsa"
    }
  },
  {
    "time": "2022-01-01 00:00:00 +0000 UTC",
    "connection": {
      "name": "production",
      "host": "host2.name",
      "port": "5412",
      "user": "ubuntu"
    }
  }
]
`

func TestParsing(t *testing.T) {
	history, err := Fetch([]byte(historyFile))

	if err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	if len(history) != 2 {
		t.Errorf("Parsing config file failed: got %v, want %v\n", len(history), 3)
	}

	if history[0].Connection.Host != "host.name" {
		t.Errorf("Parsing config file failed: got %v, want %v\n", history[0].Connection.Host, "host.name")
	}

	if history[0].Connection.Port != "" {
		t.Errorf("Parsing config file failed: got %v, want %v\n", history[0].Connection.Port, "")
	}

	if history[0].Connection.User != "" {
		t.Errorf("Parsing config file failed: got %v, want %v\n", history[0].Connection.Port, "")
	}

	if history[1].Connection.Host != "host2.name" {
		t.Errorf("Parsing config file failed: got %v, want %v\n", history[1].Connection.Host, "host.name")
	}
}
