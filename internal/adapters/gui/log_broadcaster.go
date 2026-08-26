package gui

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

// LogBroadcaster gerencia a retenção em buffer circular e o streaming de logs em tempo real.
type LogBroadcaster struct {
	mu          sync.RWMutex
	capacity    int
	buffer      []string
	subscribers map[chan string]struct{}
}

// NewLogBroadcaster inicializa um broadcaster com capacidade máxima definida.
func NewLogBroadcaster(capacity int) *LogBroadcaster {
	if capacity <= 0 {
		capacity = 500
	}
	return &LogBroadcaster{
		capacity:    capacity,
		buffer:      make([]string, 0, capacity),
		subscribers: make(map[chan string]struct{}),
	}
}

// Ensure io.Writer interface is implemented
var _ io.Writer = (*LogBroadcaster)(nil)

// Write implementa io.Writer para permitir que o broadcaster receba saídas do log padrão do Go.
func (b *LogBroadcaster) Write(p []byte) (n int, err error) {
	lines := strings.Split(strings.TrimRight(string(p), "\r\n"), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			b.Broadcast(trimmed)
		}
	}
	return len(p), nil
}

// Broadcast adiciona uma mensagem ao buffer circular e encaminha para todos os assinantes ativos.
func (b *LogBroadcaster) Broadcast(msg string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	timestamp := time.Now().Format("15:04:05")
	formatted := fmt.Sprintf("[%s] %s", timestamp, msg)

	// Adiciona ao buffer circular mantendo a capacidade
	if len(b.buffer) >= b.capacity {
		b.buffer = b.buffer[1:]
	}
	b.buffer = append(b.buffer, formatted)

	// Envia aos assinantes sem bloquear caso o canal esteja cheio
	for ch := range b.subscribers {
		select {
		case ch <- formatted:
		default:
		}
	}
}

// Subscribe registra um novo canal ouvinte e retorna todo o histórico de logs acumulado.
func (b *LogBroadcaster) Subscribe() (history []string, ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch = make(chan string, 100)
	b.subscribers[ch] = struct{}{}

	history = make([]string, len(b.buffer))
	copy(history, b.buffer)

	return history, ch
}

// Unsubscribe remove o canal de ouvinte e libera recursos.
func (b *LogBroadcaster) Unsubscribe(ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.subscribers[ch]; exists {
		delete(b.subscribers, ch)
		close(ch)
	}
}

// GetHistory retorna uma cópia de todas as mensagens no buffer circular.
func (b *LogBroadcaster) GetHistory() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	history := make([]string, len(b.buffer))
	copy(history, b.buffer)
	return history
}

// Clear limpa o buffer de mensagens acumuladas.
func (b *LogBroadcaster) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buffer = make([]string, 0, b.capacity)
}
