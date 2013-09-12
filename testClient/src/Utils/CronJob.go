package Utils

import (
        "time"
)

// cronCommand operates on a crontab.
type cronCommand func() bool

// CheckFunc is the function type for checking if a job
// shall be performed now. It also returns if a job shall
// be deleted after execution.
type CheckFunc func(time.Time) (bool, bool)

// TaskFunc is the function type that will be performed 
// if a jobs check func returns true.
type TaskFunc func(string)

// job represents one cronological job.
type job struct {
        id    string
        check CheckFunc
        task  TaskFunc
}

// checkAndPerform checks, if a job shall be performed. If true the
// task function will be called.
func (j *job) checkAndPerform(time time.Time) bool {
        perform, delete := j.check(time)
        if perform {
                go j.task(j.id)
        }
        return perform && delete
}

// Crontab is one cron server. A system can run multiple in
// parallel.
type Crontab struct {
        jobs        map[string]*job
        commandChan chan cronCommand
        ticker      *time.Ticker
}

// NewCrontab creates a cron server.
func NewCrontab() *Crontab {
        c := &Crontab{
                jobs:        make(map[string]*job),
                commandChan: make(chan cronCommand),
                ticker:      time.NewTicker(1e9),
        }
        go c.backend()
        return c
}

// Stop terminates the server.
func (c *Crontab) Stop() {
        c.commandChan <- func() bool {
                return true
        }
}

// AddJob adds a new job to the server.
func (c *Crontab) AddJob(id string, cf CheckFunc, tf TaskFunc) {
        c.commandChan <- func() bool {
                c.jobs[id] = &job{id, cf, tf}
                return false
        }
}

// DeleteJob removes a job from the server.
func (c *Crontab) DeleteJob(id string) {
        c.commandChan <- func() bool {
                delete(c.jobs, id)
                return false
        }
}

// Crontab backend.
func (c *Crontab) backend() {
        for {
                select {
                case cmd := <-c.commandChan:
                        // A server command.
                        if cmd() {
                                c.ticker.Stop()
                                return
                        }
                case <-c.ticker.C:
                        // One tick every second.
                        c.tick()
                }
        }
}

// Handle one server tick.
func (c *Crontab) tick() {
        now := time.Now().UTC()
        deletes := make(map[string]*job)
        // Check and perform jobs.
        for id, job := range c.jobs {
                delete := job.checkAndPerform(now)

                if delete {
                        deletes[id] = job
                }
        }
        // Delete those marked for deletion.
        for id, _ := range deletes {
                delete(c.jobs, id)
        }
}
