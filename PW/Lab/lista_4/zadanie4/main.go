package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"zadanie4/monitor"
)

var WaitGroup sync.WaitGroup

const (
	NrOfWriters = 5
	NrOfReaders = 10

	MinSteps = 50
	MaxSteps = 100

	MinDelay = 10 * time.Millisecond
	MaxDelay = 50 * time.Millisecond

	BoardWidth  = NrOfWriters + NrOfReaders
	BoardHeight = 4
)

type ProcessState int

const (
	LocalSection ProcessState = iota
	Start
	ReadingRoom
	Stop
)

// elementy monitora
var (
	rwMonitor      *monitor.Monitor
	okToRead       *monitor.Condition
	okToWrite      *monitor.Condition
	activeReaders  int
	activeWriters  int
	waitingWriters int
)

func initializeMonitor() {
	rwMonitor = monitor.NewMonitor()
	okToRead = rwMonitor.NewCondition()
	okToWrite = rwMonitor.NewCondition()
	activeReaders = 0
	activeWriters = 0
	waitingWriters = 0
}

func startWrite() {
	rwMonitor.Lock()
	waitingWriters++
	for activeReaders > 0 || activeWriters > 0 {
		okToWrite.Wait()
	}
	waitingWriters--
	activeWriters++
	rwMonitor.Unlock()
}

func endWrite() {
	rwMonitor.Lock()
	activeWriters--
	if waitingWriters > 0 {
		okToWrite.Signal()
	} else {
		okToRead.SignalAll()
	}
	rwMonitor.Unlock()
}

func startRead() {
	rwMonitor.Lock()
	for activeWriters > 0 || waitingWriters > 0 {
		okToRead.Wait()
	}
	activeReaders++
	rwMonitor.Unlock()
}

func endRead() {
	rwMonitor.Lock()
	activeReaders--
	if activeReaders == 0 && waitingWriters > 0 {
		okToWrite.Signal()
	}
	rwMonitor.Unlock()
}

type UserState int

const (
	R UserState = iota
	W
)

func (ur *UserState) Symbol() rune {
	switch *ur {
	case R:
		return 'R'
	case W:
		return 'W'
	default:
		return '?'
	}
}

var startTime = time.Now()

type Position struct {
	X int
	Y int
}

type TraceType struct {
	Time_Stamp time.Time
	Id         int
	Position   Position
	Symbol     rune
}

type TraceArray [MaxSteps + 111]TraceType

type Traces_Sequence_Type struct {
	Last       int
	TraceArray TraceArray
}

func PrintTrace(t TraceType) {
	elapsed := t.Time_Stamp.Sub(startTime).Seconds()
	fmt.Printf("%.6f %d %d %d %c\n", elapsed, t.Id, t.Position.X, t.Position.Y, t.Symbol)
}

func PrintTraces(t Traces_Sequence_Type) {
	for i := 0; i <= t.Last; i++ {
		PrintTrace(t.TraceArray[i])
	}
}

var reportChannel = make(chan Traces_Sequence_Type)

func printer() {
	for i := 0; i < NrOfWriters+NrOfReaders; i++ {
		traces := <-reportChannel
		PrintTraces(traces)
	}

	fmt.Printf("-1 %d %d %d ", BoardWidth, BoardWidth, BoardHeight)

	states := []string{"LocalSection", "Start", "ReadingRoom", "Stop"}
	for _, state := range states {
		fmt.Printf("%s;", state)
	}

	fmt.Println("EXTRA_LABEL;")
	WaitGroup.Done()
}

type Process struct {
	Id       int
	Role     UserState
	Position Position
}

func process(id int, symbol UserState, seed int) {
	defer WaitGroup.Done()
	r := rand.New(rand.NewSource(int64(seed)))

	var state ProcessState = LocalSection

	var process Process
	process.Id = id
	process.Role = symbol
	process.Position.X = id
	process.Position.Y = int(state)

	var traces Traces_Sequence_Type
	traces.Last = -1

	storeTrace := func() {
		process.Position.Y = int(state)
		timeStamp := time.Since(startTime)
		traces.Last++
		traces.TraceArray[traces.Last] = TraceType{
			Time_Stamp: startTime.Add(timeStamp),
			Id:         process.Id,
			Position:   process.Position,
			Symbol:     process.Role.Symbol(),
		}
	}

	storeTrace()

	nrOfSteps := MinSteps + r.Intn(MaxSteps-MinSteps+1)

	for i := 0; i < (nrOfSteps/4)-1; i++ {
		delay := MinDelay + time.Duration(r.Int63n(int64(MaxDelay-MinDelay)))
		state = LocalSection
		storeTrace()
		time.Sleep(delay)

		//próba wejścia
		if process.Role == W {
			startWrite()
		} else {
			startRead()
		}
		state = Start
		storeTrace()
		time.Sleep(delay)

		state = ReadingRoom
		storeTrace()
		time.Sleep(delay)

		//koniec
		if process.Role == W {
			endWrite()
		} else {
			endRead()
		}
		state = Stop
		storeTrace()
		time.Sleep(delay)

		state = LocalSection
		storeTrace()
	}
	reportChannel <- traces
}

func main() {
	initializeMonitor()
	WaitGroup.Add(1)
	go printer()

	for j := 0; j < NrOfReaders; j++ {
		WaitGroup.Add(1)
		go process(j, 0, j)
	}

	for i := NrOfReaders; i < NrOfReaders+NrOfWriters; i++ {
		WaitGroup.Add(1)
		go process(i, 1, i)
	}
	WaitGroup.Wait()
}
