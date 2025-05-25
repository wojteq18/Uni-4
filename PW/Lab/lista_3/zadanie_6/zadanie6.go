package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var WaitGroup sync.WaitGroup

const (
	NrOfProcess = 2

	MinSteps = 50
	MaxSteps = 100

	MinDelay = 10 * time.Millisecond
	MaxDelay = 50 * time.Millisecond

	BoardWidth  = NrOfProcess
	BoardHeight = 4
)

type ProcessState int

const (
	LocalSection ProcessState = iota
	EntryProtocol
	CriticalSection
	ExitProtocol
)

//zmienne globalne potrzebne do algorytmu Petersona

var isInterested [NrOfProcess]int32
var Last int32 = -1

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

type TraceArray [MaxSteps + 1]TraceType

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
	for i := 0; i < NrOfProcess; i++ {
		traces := <-reportChannel
		PrintTraces(traces)
	}

	fmt.Printf("-1 %d %d %d ", NrOfProcess, BoardWidth, BoardHeight)

	states := []string{"Local_Section", "Entry_Protocol", "Critical_Section", "Exit_Protocol"}
	for _, state := range states {
		fmt.Printf("%s;", state)
	}

	fmt.Println("EXTRA_LABEL;")
	WaitGroup.Done()
}

type Process struct {
	Id       int
	Symbol   rune
	Position Position
}

func process(id int, symbol rune, seed int) {
	defer func() {
		atomic.StoreInt32(&isInterested[id], 0)
	}()

	defer WaitGroup.Done()
	r := rand.New(rand.NewSource(int64(seed)))

	var state ProcessState = LocalSection

	var process Process
	process.Id = id
	process.Symbol = symbol
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
			Symbol:     process.Symbol,
		}
	}

	storeTrace()

	nrOfSteps := MinSteps + r.Intn(MaxSteps-MinSteps+1)

	for i := 0; i < nrOfSteps; i++ {
		delay := MinDelay + time.Duration(r.Int63n(int64(MaxDelay-MinDelay)))
		time.Sleep(delay)
		state = ProcessState((int(state) + 1) % BoardHeight)

		if state == EntryProtocol {
			atomic.StoreInt32(&isInterested[id], 1)
			atomic.StoreInt32(&Last, int32(1-id))
			for (atomic.LoadInt32(&isInterested[1-id]) == 1) && atomic.LoadInt32(&Last) == int32(1-id) {
				time.Sleep((1 * time.Millisecond))
			}
		}

		if state == ExitProtocol {
			atomic.StoreInt32(&isInterested[id], 0)
		}
		storeTrace()
	}
	if state == CriticalSection {
		state = LocalSection
		storeTrace()
	}
	reportChannel <- traces
}

func main() {
	WaitGroup.Add(1)
	go printer()
	symbols := []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O'}
	for i := 0; i < NrOfProcess; i++ {
		WaitGroup.Add(1)
		go process(i, symbols[i], i)
	}
	WaitGroup.Wait()
}
