package network

import (
	"fmt"
	"os"
	"path"
	"ssemu/internal/utils"
	"strings"
	"time"
)

func (p *Packet) DumpPacket(origin string) {
	var sb strings.Builder
	spacer := strings.Repeat("-", 75)
	dateStr := time.Now().Format("20060102-150405")
	sb.WriteString(spacer + utils.EnvNewLine())
	sb.WriteString(fmt.Sprintf("date: %s | packet: %02X | origin: %s", dateStr, p.id, origin))
	sb.WriteString(utils.EnvNewLine() + spacer + utils.EnvNewLine())
	sb.WriteString(utils.BufferToString(p.buffer))
	wd, _ := os.Getwd()
	fileName := fmt.Sprintf("packetdump_%02X_%s.txt", p.id, dateStr)
	output := []byte(sb.String())
	os.WriteFile(path.Join(wd, fileName), output, 0o755)
}
