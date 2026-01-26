package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	REQUEST int = 0
	RELEASE int = 1
	REPLY   int = 2
)

type Message struct {
	pid       int
	timestamp int
	msgtype   int
}

type Process struct {
	pid       int
	clock     int
	processes []*Process
	inbox     []Message
	queue     []Message
	replies   int
	mutex     sync.Mutex
	finished  bool
	ready     bool
}

var (
	pollInterval = 25 * time.Millisecond
	csDuration   = 200 * time.Millisecond
	mainSleep    = 5 * time.Second
)

// citire linie de la tastatură
func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

// citire int cu default
func readIntDefault(reader *bufio.Reader, prompt string, def int) int {
	fmt.Printf("%s [%d]: ", prompt, def)
	s := readLine(reader)
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

// mesaj de la un proces la altul (REQUEST-0, RELEASE-1, REPLY-2)
func (p *Process) sendMessage(msg Message, pid int) {
	proc := p.processes[pid]
	proc.mutex.Lock()
	proc.inbox = append(proc.inbox, msg)
	proc.mutex.Unlock()
}

// cere Replies pentru a intra in CS, verificare cu canRequest dupa.
func (p *Process) request() {
	p.mutex.Lock()

	fmt.Println("pid", p.pid, "request")
	p.clock++

	p.replies = 0
	req := Message{p.pid, p.clock, REQUEST} // request de trimis
	p.queue = append(p.queue, req)

	p.mutex.Unlock()

	// trimite mesaje
	for i, proc := range p.processes {
		if proc.pid == p.pid {
			continue
		}
		p.sendMessage(req, i)
	}
}

// nu mai foloseste CS si trimite un mesaj, elimina din queue
func (p *Process) release() {
	p.mutex.Lock()
	p.clock++
	req := Message{p.pid, p.clock, RELEASE}
	p.queue = p.queue[1:]
	p.mutex.Unlock()

	for i, proc := range p.processes {
		if proc.pid == p.pid {
			continue
		}
		p.sendMessage(req, i)
	}
}

func (p *Process) read() {
	for {
		if p.finished {
			return
		}
		p.mutex.Lock()

		// Vom trimite reply dupa eliberarea lock-ului pentru a nu produce un deadlock
		var replies []struct {
			msg Message
			pid int
		}

		for _, msg := range p.inbox {
			if p.clock < msg.timestamp {
				p.clock = msg.timestamp
			}
			p.clock++

			if msg.msgtype == REPLY {
				p.replies++
			}

			if msg.msgtype == REQUEST { // adauga la queue si da reply
				p.queue = append(p.queue, msg)

				// sortam dupa timestamp, altfel dupa pid
				sort.Slice(p.queue, func(i, j int) bool {
					a, b := p.queue[i], p.queue[j]
					if a.timestamp == b.timestamp {
						return a.pid < b.pid
					}
					return a.timestamp < b.timestamp
				})

				reply := Message{p.pid, p.clock, REPLY}
				replies = append(replies, struct {
					msg Message
					pid int
				}{reply, msg.pid})
			}

			// eliminam din queue daca a terminat
			if msg.msgtype == RELEASE {
				tempQueue := []Message{}
				for _, m := range p.queue {
					if m.pid != msg.pid {
						tempQueue = append(tempQueue, m)
					}
				}
				p.queue = tempQueue
			}
		}
		p.inbox = []Message{}
		p.mutex.Unlock()

		for _, r := range replies {
			p.sendMessage(r.msg, r.pid)
		}

		time.Sleep(pollInterval)
	}
}

// verificare pentru a incepe codul din CS
func (p *Process) canRequest() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return len(p.queue) > 0 &&
		p.queue[0].pid == p.pid &&
		p.replies == len(p.processes)-1
}

// Codul de executat pt fiecare "proces"
func csCode(p *Process) {
	p.request()

	for !p.canRequest() {
		time.Sleep(pollInterval)
	}

	fmt.Println("pid", p.pid, "entered CS")
	time.Sleep(csDuration)
	fmt.Println("pid", p.pid, "left CS")

	p.release()
	p.finished = true
}

func main() {
	// citire parametri
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Parametri (Enter = default):")

	p_count := readIntDefault(reader, "Numar procese (p_count)", 5)

	clocks := make([]int, p_count)

	for i := 0; i < p_count; i++ {
		clocks[i] = readIntDefault(
			reader,
			fmt.Sprintf("Clock initial pentru pid %d", i),
			0, // default value (change if you want)
		)
	}

	pollMs := readIntDefault(reader, "Poll interval (ms)", 25)
	//CS - critical section
	csMs := readIntDefault(reader, "CS duration (ms)", 200)
	runtimeSec := readIntDefault(reader, "Main runtime (sec)", 5)

	// setăm variabilele folosite de Sleep-uri
	if pollMs < 1 {
		pollMs = 1
	}
	if csMs < 0 {
		csMs = 0
	}
	if runtimeSec < 1 {
		runtimeSec = 1
	}

	pollInterval = time.Duration(pollMs) * time.Millisecond
	csDuration = time.Duration(csMs) * time.Millisecond
	mainSleep = time.Duration(runtimeSec) * time.Second

	processes := []*Process{}

	// initialiare toate p_count procese
	for i := 0; i < p_count; i++ {
		processes = append(processes, &Process{
			pid:      i,
			finished: false,
			clock:    clocks[i],
		})
	}

	// un goroutine pt a verifica "inbox"-ul la fiecare proces
	for i := 0; i < p_count; i++ {
		//fiecare proces primește o referință la toate celelalte procese.
		processes[i].processes = processes
		go processes[i].read()
	}

	// asteptam sa inceapa toate procesele
	time.Sleep(10 * time.Millisecond)
	fmt.Println()

	// fiecare proces va accesa CS o data
	for i := 0; i < p_count; i++ {
		go csCode(processes[i])
	}

	// lăsăm main să "doarmă" suficient cât toate procesele să termine
	time.Sleep(mainSleep)
}
