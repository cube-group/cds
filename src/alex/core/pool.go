package core

import (
    "errors"
    "time"
    "sync"
    "io"
    "fmt"
)

//连接池生产
type PoolFactory func() (io.Closer, error)

//连接池脚手架
type Pool struct {
    mu      sync.Mutex

    MaxIdle int
    MaxOpen int
    IsDebug bool
    Name    string

    busy    bool
    factory PoolFactory
    stack   chan io.Closer
}

//new MysqlPool
func NewPool(factory PoolFactory, maxIdle int, maxOpen int, name string, debug bool) *Pool {
    pool := new(Pool)
    pool.Name = name
    pool.factory = factory
    pool.IsDebug = debug
    pool.setInit(maxIdle, maxOpen)
    return pool
}

//log
func (this *Pool)log(value ...interface{}) {
    if this.IsDebug {
        fmt.Println("[pool]", this.Name, value)
    }
}

//set init
func (this *Pool)setInit(maxIdle int, maxOpen int) error {
    if maxOpen < maxIdle {
        return errors.New("maxOpen can not less than maxIdle")
    }

    this.stack = make(chan io.Closer, maxOpen)

    //init maxIdle start
    var wg sync.WaitGroup
    for i := 0; i < maxIdle; i++ {
        wg.Add(1)
        go func() {
            if created := this.factoryCreate(); created == true {
                wg.Done()
            }
        }()
    }
    wg.Wait()
    //init maxIdle end

    //factory listening
    go func() {
        for {
            busyState := this.busy && len(this.stack) < maxOpen
            idleState := len(this.stack) < maxIdle
            if maxIdle <= 0 || busyState || idleState {
                this.factoryCreate()
            }
            time.Sleep(time.Microsecond * 10)
        }
    }()
    return nil
}

//factory create
func (this *Pool)factoryCreate() bool {
    one, err := this.factory()
    if err == nil {
        this.stack <- one
        this.log("create one", len(this.stack))
        return true
    }
    return false
}


//back to pool
func (this *Pool)Back(one io.Closer) error {
    if one != nil {
        return one.Close()
    }
    return errors.New("back instance is nil")
}

//get instance
func (this *Pool)Get() (io.Closer, error) {
    this.mu.Lock()
    defer this.mu.Unlock()

    if this.MaxIdle > 0 && len(this.stack) < this.MaxIdle {
        this.busy = true
    } else {
        this.busy = false
    }

    select {
    case one := <-this.stack:
        this.log("use one")
        return one, nil
    case <-time.After(time.Microsecond * 10):
        this.busy = true
        return nil, errors.New("pool timeout")
    }
}