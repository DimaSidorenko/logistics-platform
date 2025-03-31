package ratelimiter

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Limiter struct {
	ch      chan struct{}
	closeCh chan struct{}
	once    sync.Once
}

func New(rpsLimit int) *Limiter {
	if rpsLimit <= 0 {
		panic("rpsLimit must be positive")
	}

	l := &Limiter{
		ch:      make(chan struct{}, rpsLimit),
		closeCh: make(chan struct{}),
	}

	for i := 0; i < rpsLimit; i++ {
		l.ch <- struct{}{}
	}

	go l.refillTokens(rpsLimit)

	return l
}

func (l *Limiter) refillTokens(rpsLimit int) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Пополняем токены до rpsLimit
			for i := 0; i < rpsLimit; i++ {
				select {
				case l.ch <- struct{}{}:
				default:
					// Бак уже полон
				}
			}
		case <-l.closeCh:
			return
		}
	}
}

func (l *Limiter) Wait(ctx context.Context) error {
	select {
	case <-l.closeCh:
		return fmt.Errorf("ratelimiter is closed")
	case <-ctx.Done():
		return ctx.Err()
	default:
		select {
		case <-l.ch:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case <-l.closeCh:
			return fmt.Errorf("ratelimiter is closed")
		}
	}
}

func (l *Limiter) Stop() {
	l.once.Do(func() {
		close(l.closeCh)
		close(l.ch)
	})
}
