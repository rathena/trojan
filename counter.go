package integration

import (
	"strconv"
	"sync"
)

// A counter used to safely track concurrent operations.
type Counter struct {
	sync.RWMutex
	account       int
	accountsIndex int
	accounts      []string
	parallelCount map[string]int
	failCount     map[string]int
}

// Capture the name of a benchmark and increment the count so we know how many
// parallel threads were created to run this benchmark.
func (c *Counter) Parallel(name string) {
	c.Lock()
	c.parallelCount[name] += 1
	c.Unlock()
}

// Get the next numeric index for account creation.
func (c *Counter) Next() string {
	c.Lock()
	next := strconv.Itoa(c.account)
	c.account++
	c.Unlock()
	return next
}

// Save a successfully created account name.
func (c *Counter) SaveAccount(name string) {
	c.Lock()
	c.accounts = append(c.accounts, name)
	c.Unlock()
}

// Get an account from the array.
func (c *Counter) GetAccount() string {
	c.RLock()
	name := c.accounts[c.accountsIndex]
	c.accountsIndex = (c.accountsIndex + 1) % len(c.accounts)
	c.RUnlock()
	return name
}

// Count the number of failed operations by benchmark name.
func (c *Counter) Fail(name string) {
	c.Lock()
	c.failCount[name] += 1
	c.Unlock()
}
