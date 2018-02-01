package utils

import (
    "time"
)

type TimerHandler func(...interface{})

type Timer struct {
    ticker  *time.Ticker
    Times   uint64
    Count   uint64
    Handler TimerHandler
    Args    []interface{}
}

func (this *Timer)clearTicker() {
    if this.ticker != nil {
        this.ticker.Stop()
        this.ticker = nil
    }
}

func (this *Timer)Listen(handler TimerHandler, args ...interface{}) {
    this.Handler = handler
    this.Args = args
}

func (this *Timer)Start(interval uint64) {
    this.clearTicker()

    this.ticker = time.NewTicker(time.Duration(interval) * time.Second)
    go func() {
        for {
            if this.ticker == nil {
                break
            }
            if this.Times > 0 && this.Count == this.Times {
                this.Stop()
                break
            }
            select {
            case <-this.ticker.C:
                if (this.Handler != nil) {
                    this.Handler(this.Args...)
                }
                this.Count++
            }
        }
    }()
}

func (this *Timer)Stop() {
    this.clearTicker()

    this.Handler = nil
    this.Args = nil
    this.Count = 0;
    this.Times = 0
}

func SetInterval(interval uint64, handler TimerHandler, args ...interface{}) *Timer {
    timer := new(Timer)
    timer.Listen(handler, args...)
    timer.Start(interval)
    return timer
}

func ClearInterval(t *Timer) {
    if t != nil {
        t.Stop()
    }
}

func SetTimeout(interval uint64, handler TimerHandler, args ...interface{}) *Timer {
    timer := new(Timer)
    timer.Times = 1
    timer.Listen(handler, args...)
    timer.Start(interval)
    return timer
}

func ClearTimeout(t *Timer) {
    if t != nil {
        t.Stop()
    }
}