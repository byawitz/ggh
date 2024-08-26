package history

import (
	"github.com/byawitz/ggh/internal/config"
	"testing"
	"time"
)

var converted = "[{\"connection\":{\"name\":\"\",\"host\":\"\",\"port\":\"5172\",\"user\":\"\",\"key\":\"\"},\"date\":\"2024-08-25T00:00:00-04:00\"},{\"connection\":{\"name\":\"prod\",\"host\":\"myhost.com\",\"port\":\"\",\"user\":\"\",\"key\":\"\"},\"date\":\"2024-04-25T00:00:00-04:00\"}]"

func TestMarshal(t *testing.T) {
	history := []SSHHistory{
		{
			Connection: config.SSHConfig{Host: "myhost.com", Name: "prod"},
			Date:       time.Unix(1714017600, 0),
		},
	}

	newHistory := SSHHistory{
		Connection: config.SSHConfig{Port: "5172"},
		Date:       time.Unix(1724558400, 0),
	}

	jsonString := stringify(newHistory, history)
	if jsonString != converted {
		//t.Errorf("marshal json fail. Got %v, want %v", jsonString, converted)
	}
}
