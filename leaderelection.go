package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var electionRunning atomic.Bool
var leaderID atomic.Int32

// Process struct representing a process in the election
type Process struct {
	ID      int
	Next    chan int
	Leader  chan int
	IsAlive bool
}

// StartElection starts an election only if no other is running
func (p *Process) StartElection(wg *sync.WaitGroup) {
	defer wg.Done()

	if !electionRunning.CompareAndSwap(false, true) || !p.IsAlive {
		return
	}

	electionID := p.ID
	fmt.Printf("Process %d started election.\n", p.ID)
	p.Next <- electionID

	for {
		select {
		case receivedID := <-p.Next:
			if receivedID > p.ID {
				p.Next <- receivedID
			} else if receivedID < p.ID {
				p.Next <- p.ID
			} else {
				fmt.Printf("Process %d is the new leader!\n", p.ID)
				leaderID.Store(int32(p.ID))
				broadcastLeader(p.ID)
				electionRunning.Store(false)
				return
			}
		case <-time.After(3 * time.Second):
			electionRunning.Store(false)
			return
		}
	}
}

// Broadcast leader information to all processes
func broadcastLeader(leader int) {
	fmt.Printf("Broadcasting Leader: %d\n", leader)
	leaderID.Store(int32(leader))
	for i := 0; i < 5; i++ {
		fmt.Printf("Process %d acknowledges Leader: %d\n", i, leader)
	}
}

// ListenForLeader ensures processes wait for a leader announcement before starting a new election
func (p *Process) ListenForLeader(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case leader := <-p.Leader:
			fmt.Printf("Process %d acknowledges Leader: %d\n", p.ID, leader)
			return
		case <-time.After(5 * time.Second):
			if !electionRunning.Load() && leaderID.Load() == -1 {
				wg.Add(1)
				go p.StartElection(wg)
			}
			return
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	numProcesses := 5
	var wgElection sync.WaitGroup
	var wgLeader sync.WaitGroup

	processes := make([]*Process, numProcesses)
	channels := make([]chan int, numProcesses)
	leaderChannels := make([]chan int, numProcesses)

	// Initialize channels
	for i := range channels {
		channels[i] = make(chan int, 1)
		leaderChannels[i] = make(chan int, 1)
	}

	leaderID.Store(-1) // No leader initially
	var highestProcess *Process

	// Create and initialize processes
	for i := 0; i < numProcesses; i++ {
		id := rand.Intn(100)
		isAlive := rand.Float32() > 0.2
		processes[i] = &Process{
			ID:      id,
			Next:    channels[(i+1)%numProcesses],
			Leader:  leaderChannels[i],
			IsAlive: isAlive,
		}
		fmt.Printf("Process %d has ID %d (Alive: %t)\n", i, id, isAlive)

		if isAlive && (highestProcess == nil || id > highestProcess.ID) {
			highestProcess = processes[i]
		}
	}

	if highestProcess != nil {
		wgElection.Add(1)
		go highestProcess.StartElection(&wgElection)
	} else {
		fmt.Println("No active processes available. Election cannot start.")
		return
	}

	for i := 0; i < numProcesses; i++ {
		wgLeader.Add(1)
		go processes[i].ListenForLeader(&wgLeader)
	}

	wgElection.Wait()
	wgLeader.Wait()
}
