package go_ibm_db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	trc "github.com/ibmdb/go_ibm_db/log2"
)

// DBP struct type contains the timeout, dbinstance pointer and connection string
type DBP struct {
	*sql.DB
	con  string
	n    time.Duration
	pool *Pool // owning pool; avoids relying on the shared global pool pointer
}

// Pool struct contais the about the pool like size, used and available connections
type Pool struct {
	availablePool map[string][]*DBP
	usedPool      map[string][]*DBP
	poolSize      int
	maxLifetime   time.Duration
	closed        bool
	mu            sync.Mutex
}

const defaultMaxIdleConns = 2
const defaultConnMaxLifetime = 60

// Pconnect will return the pool instance
func Pconnect(poolSize string) *Pool {
	trc.Trace1("pooling.go: Pconnect() - ENTRY")
	trc.Trace1(fmt.Sprintf("poolSize=%s", poolSize))

	var size int
	count := len(poolSize)
	if count > 0 {
		opt := strings.Split(poolSize, "=")
		if opt[0] == "PoolSize" {
			size, _ = strconv.Atoi(opt[1])
			if size <= 0 {
				size = defaultMaxIdleConns
			}
		} else {
			fmt.Println("Not a valid parameter")
		}
	} else {
		size = defaultMaxIdleConns
	}
	p := &Pool{
		availablePool: make(map[string][]*DBP),
		usedPool:      make(map[string][]*DBP),
		poolSize:      size,
		maxLifetime:   time.Duration(defaultConnMaxLifetime) * time.Second,
	}

	trc.Trace1("pooling.go: Pconnect() - EXIT")
	return p
}

// size returns the total number of connections (used + available) currently
// owned by this pool. Callers must hold p.mu.
func (p *Pool) size() int {
	total := 0
	for _, v := range p.usedPool {
		total += len(v)
	}
	for _, v := range p.availablePool {
		total += len(v)
	}
	return total
}

func (p *Pool) takeAvailableLocked(connStr string, lifetime time.Duration) *DBP {
	val, ok := p.availablePool[connStr]
	if !ok || len(val) == 0 {
		return nil
	}

	dbpo := val[0]
	if len(val) == 1 {
		delete(p.availablePool, connStr)
	} else {
		copy(val[0:], val[1:])
		val[len(val)-1] = nil
		p.availablePool[connStr] = val[:len(val)-1]
	}
	p.usedPool[connStr] = append(p.usedPool[connStr], dbpo)
	dbpo.DB.SetConnMaxLifetime(lifetime)
	return dbpo
}

// Open will check for the connection in the pool
// If not opens a new connection and stores in the pool
func (p *Pool) Open(connStr string, options ...string) *DBP {
	trc.Trace1("pooling.go: Open() - ENTRY")
	trc.Trace1(fmt.Sprintf("connStr=%s", connStr))

	var Time time.Duration
	configuredLifetime := false
	count := len(options)
	if count > 0 {
		for i := 0; i < count; i++ {
			opt := strings.Split(options[i], "=")
			if opt[0] == "SetConnMaxLifetime" {
				lifetime, _ := strconv.Atoi(opt[1])
				if lifetime <= 0 {
					lifetime = defaultConnMaxLifetime
				}
				Time = time.Duration(lifetime) * time.Second
				configuredLifetime = true
			} else {
				fmt.Println("Not a valid parameter")
			}
		}
	}
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	if configuredLifetime {
		p.maxLifetime = Time
	} else {
		Time = p.maxLifetime
	}
	timeout := p.maxLifetime
	p.mu.Unlock()
	deadline := time.Now().Add(timeout)
	for {
		p.mu.Lock()
		if dbpo := p.takeAvailableLocked(connStr, Time); dbpo != nil {
			p.mu.Unlock()
			return dbpo
		}
		if p.size() < p.poolSize {
			db, err := sql.Open("go_ibm_db", connStr)
			if err != nil {
				p.mu.Unlock()
				return nil
			}
			dbi := &DBP{
				DB:   db,
				con:  connStr,
				n:    Time,
				pool: p,
			}
			p.usedPool[connStr] = append(p.usedPool[connStr], dbi)
			dbi.DB.SetConnMaxLifetime(Time)
			p.mu.Unlock()
			return dbi
		}
		p.mu.Unlock()

		if !time.Now().Before(deadline) {
			fmt.Println("Connection timeout")
			trc.Trace1("Connection timeout")
			trc.Trace1("pooling.go: Open() - EXIT")
			return nil
		}
		time.Sleep(3 * time.Second)
	}
}

func (p *Pool) Init(numConn int, connStr string) bool {
	trc.Trace1("pooling.go: Init() - ENTRY")
	trc.Trace1(fmt.Sprintf("numConn=%d, connStr=%s", numConn, connStr))

	var Time time.Duration

	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return false
	}
	Time = p.maxLifetime
	p.mu.Unlock()

	for i := 0; i < numConn; i++ {
		db, err := sql.Open("go_ibm_db", connStr)
		if err != nil {
			trc.Trace1("pooling.go: Init() - return false")
			return false
		}
		dbi := &DBP{
			DB:   db,
			con:  connStr,
			n:    Time,
			pool: p,
		}
		p.mu.Lock()
		if p.closed {
			p.mu.Unlock()
			db.Close()
			trc.Trace1("pooling.go: Init() - return false")
			return false
		}
		p.availablePool[connStr] = append(p.availablePool[connStr], dbi)
		dbi.DB.SetConnMaxLifetime(Time)
		p.mu.Unlock()
	}
	trc.Trace1("pooling.go: Init() - EXIT")
	return true
}

// Close will make the connection available for the next release
func (d *DBP) Close() {
	trc.Trace1("pooling.go: Close() - ENTRY")

	p := d.pool
	if p == nil {
		d.DB.Close()
		return
	}
	found := false
	var pos int
	p.mu.Lock()
	if valc, okc := p.usedPool[d.con]; okc {
		for i, v := range valc {
			if v == d {
				pos = i
				found = true
				break
			}
		}
		if found {
			if len(valc) > 1 {
				dbpc := valc[pos]
				copy(valc[pos:], valc[pos+1:])
				valc[len(valc)-1] = nil
				valc = valc[:len(valc)-1]
				p.usedPool[d.con] = valc
				p.availablePool[d.con] = append(p.availablePool[d.con], dbpc)
			} else {
				dbpc := valc[0]
				p.availablePool[d.con] = append(p.availablePool[d.con], dbpc)
				delete(p.usedPool, d.con)
			}
			go d.Timeout()
		} else {
			d.DB.Close()
		}
	} else {
		d.DB.Close()
	}
	p.mu.Unlock()
	trc.Trace1("pooling.go: Close() - EXIT")
}

// Timeout for closing the connection in pool
func (d *DBP) Timeout() {
	trc.Trace1("pooling.go: Timeout() - ENTRY")

	p := d.pool
	<-time.After(d.n)
	if p == nil {
		return
	}
	found := false
	var pos int
	p.mu.Lock()
	if valt, okt := p.availablePool[d.con]; okt {
		for i, v := range valt {
			if v == d {
				pos = i
				found = true
				break
			}
		}
		if found {
			dbpt := valt[pos]
			copy(valt[pos:], valt[pos+1:])
			valt[len(valt)-1] = nil
			valt = valt[:len(valt)-1]
			if len(valt) == 0 {
				delete(p.availablePool, d.con)
			} else {
				p.availablePool[d.con] = valt
			}
			dbpt.DB.Close()
		}
	}
	p.mu.Unlock()
	trc.Trace1("pooling.go: Timeout() - EXIT")
}

// Release will close all the connections in the pool
func (p *Pool) Release() {
	trc.Trace1("pooling.go: Release() - ENTRY")

	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	if p.availablePool != nil {
		for _, vala := range p.availablePool {
			for _, dbpr := range vala {
				dbpr.DB.Close()
			}
		}
		p.availablePool = nil
	}
	if p.usedPool != nil {
		for _, valu := range p.usedPool {
			for _, dbpr := range valu {
				dbpr.DB.Close()
			}
		}
		p.usedPool = nil
	}
	p.mu.Unlock()
	trc.Trace1("pooling.go: Release() - EXIT")
}

// Set the connMaxLifetime
func (p *Pool) SetConnMaxLifetime(num int) {
	trc.Trace1("pooling.go: SetConnMaxLifetime()")
	trc.Trace1(fmt.Sprintf("connMaxLifetime=%d", num))

	if num <= 0 {
		num = defaultConnMaxLifetime
	}
	p.mu.Lock()
	p.maxLifetime = time.Duration(num) * time.Second
	p.mu.Unlock()
}

// Display will print the  values in the map
func (p *Pool) Display() {
	fmt.Println(p.availablePool)
	fmt.Println(p.usedPool)
	fmt.Println(p.poolSize)
}
