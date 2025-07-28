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
	OnError() error
	ReportError(err error)
}

type monitorImpl struct {
	conn    SSHConnection
	errChan chan error
	cancel  context.CancelFunc
	ctx     context.Context
}

const _KEEPALIVETIMEOUT = 10 * time.Second

func NewMonitor(conn SSHConnection) ConnectionMonitor {

	ctx, cancel := context.WithCancel(context.Background())
	return &monitorImpl{
		conn:    conn,
		errChan: make(chan error),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (m *monitorImpl) Start() {
	// no compruebo que no se pueda ejecutar multiples veces pero bueno

	go func() {
		ticker := time.NewTicker(_KEEPALIVETIMEOUT)

		defer func() {
			ticker.Stop()
		}()

		log.Println("Iniciado el Monitor")
		for {
			select {
			case <-ticker.C:
				if !m.conn.IsAlive() {
					m.ReportError(fmt.Errorf("connection lost, monitor stopped"))
					return
				}
			case <-m.ctx.Done():
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

func (m *monitorImpl) OnError() (err error) {

	select {

	case err = <-m.errChan:
		return err
	case <-m.ctx.Done():

	}

	return
}

func (m *monitorImpl) ReportError(err error) {
	select {
	case m.errChan <- err:
	default:
	}
}
