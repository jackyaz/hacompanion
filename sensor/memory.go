package sensor

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"hacompanion/entity"
	"hacompanion/util"
)

var reMemory = regexp.MustCompile(`(?mi)^\s?(?P<name>[^:]+):\s+(?P<value>\d+)`)

type Memory struct{}

func NewMemory() *Memory {
	return &Memory{}
}

func (m Memory) Run(ctx context.Context) (*entity.Payload, error) {
	b, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	return m.process(string(b))
}

func (m Memory) process(output string) (*entity.Payload, error) {
	p := entity.NewPayload()
	matches := reMemory.FindAllStringSubmatch(output, -1)
	for _, match := range matches {
		var err error
		var kb int
		if len(match) != 3 {
			continue
		}
		kb, err = strconv.Atoi(strings.TrimSpace(match[2]))
		if err != nil {
			continue
		}
		// Convert kb to MB.
		mb := util.RoundToTwoDecimals(float64(kb) / 1024)
		switch strings.TrimSpace(match[1]) {
		case "MemTotal":
			fallthrough
		case "MemAvailable":
			fallthrough
		case "MemFree":
			fallthrough
		case "SwapFree":
			fallthrough
		case "Buffers":
			fallthrough
		case "Cached":
			fallthrough
		case "SwapTotal":
			p.Attributes[util.ToSnakeCase(match[1])] = mb
		}
	}
	p.State = util.RoundToTwoDecimals(p.Attributes["mem_total"].(float64) - p.Attributes["mem_free"].(float64)  - p.Attributes["buffers"].(float64)  - p.Attributes["cached"].(float64))
	if p.State == "" {
		return nil, fmt.Errorf("could not determine memory state based on /proc/meminfo: %s", output)
	}
	return p, nil
}
