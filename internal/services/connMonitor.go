package services

import (
	"context"
	"fmt"
	"log"
	"time"
)

type ConnectionMonitor interface {
	Start()
	Stop()
	OnError() <-chan error
	ReportError(err error)
}

type monitorImpl struct {
	conn    SSHConnection
	errChan chan error
	cancel  context.CancelFunc
}

const _KEEPALIVETIMEOUT = 10 * time.Second

func NewMonitor(conn SSHConnection) ConnectionMonitor {
	return &monitorImpl{
		conn:    conn,
		errChan: make(chan error),
	}
}

func (m *monitorImpl) Start() {

	if m.cancel != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel

	go func() {
		ticker := time.NewTicker(_KEEPALIVETIMEOUT)

		defer func() {
			ticker.Stop()
			m.cancel = nil
		}()

		log.Println("Iniciado el Monitor")
		for {
			select {
			case <-ticker.C:
				if !m.conn.IsAlive() {
					m.ReportError(fmt.Errorf("connection lost, monitor stopped"))
					return
				}
			case <-ctx.Done():
				log.Println("Monitor stopped Manually")
				return
			}
		}
	}()
}

func (m *monitorImpl) Stop() {
	if m.cancel != nil {
		m.cancel()
	}
}

func (m *monitorImpl) OnError() <-chan error {
	return m.errChan
}

func (m *monitorImpl) ReportError(err error) {
	select {
	case m.errChan <- err:
	default:
	}
}
