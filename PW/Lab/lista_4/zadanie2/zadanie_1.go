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
	NrOfProcess = 15

	MinSteps = 30
	MaxSteps = 70

	MinDelay = 10 * time.Millisecond
	MaxDelay = 50 * time.Millisecond

	BoardWidth  = NrOfProcess
	BoardHeight = 7
)

type ProcessState int

var Flag [NrOfProcess]int32

const (
	LocalSection ProcessState = iota
	Entry_Protocol1
	Entry_Protocol2
	Entry_Protocol3
	Entry_Protocol4
	CriticalSection
	ExitProtocol
)

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
	for i := 0; i < NrOfProcess; i++ {
		traces := <-reportChannel
		PrintTraces(traces)
	}

	fmt.Printf("-1 %d %d %d ", NrOfProcess, BoardWidth, BoardHeight)

	states := []string{"LocalSection", "Entry_Protocol1", "Entry_Protocol2", "Entry_Protocol3", "Entry_Protocol4", "Critical_Section", "Exit_Protocol"}
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

	for i := 0; i < (nrOfSteps / 7); i++ {
		state = LocalSection
		storeTrace()
		time.Sleep(MinDelay + time.Duration(r.Int63n(int64(MaxDelay-MinDelay))))

		state = Entry_Protocol1
		storeTrace()
		atomic.StoreInt32(&Flag[id], 1)

		keepWaiting := true
		for keepWaiting {
			//fmt.Println("Jestem 1")
			keepWaiting = false
			for j := 0; j < NrOfProcess; j++ {
				if atomic.LoadInt32(&Flag[j]) == 3 || atomic.LoadInt32(&Flag[j]) == 4 {
					keepWaiting = true
					break
				}
			}
			if !keepWaiting {
				break
			}
			time.Sleep(1 * time.Millisecond)
		}
		state = Entry_Protocol3
		storeTrace()
		atomic.StoreInt32(&Flag[id], 3)

		foundWithFlag1 := false
		idx := -1
		for k := 0; k < NrOfProcess; k++ {
			if k == id {
				continue
			}
			if atomic.LoadInt32(&Flag[k]) == 1 {
				foundWithFlag1 = true
				idx = k
				break
			}
		}

		if foundWithFlag1 {
			state = Entry_Protocol2
			storeTrace()
			atomic.StoreInt32(&Flag[id], 2)
			//fmt.Println("Byłem w entry protocol 2!!!")

			for atomic.LoadInt32(&Flag[idx]) == 4 {
				time.Sleep(1 * time.Millisecond)
			}
		}

		state = Entry_Protocol4
		storeTrace()
		atomic.StoreInt32(&Flag[id], 4)
		waiting := true
		for waiting {
			waiting = false
			for j := 0; j < id; j++ {
				if atomic.LoadInt32(&Flag[j]) == 2 || atomic.LoadInt32(&Flag[j]) == 3 || atomic.LoadInt32(&Flag[j]) == 4 {
					waiting = true
					break
				}
			}
			if waiting == false {
				break
			}
			time.Sleep(1 * time.Millisecond)
		}

		keepWaiting = true
		for keepWaiting {
			processNeedToWait := false
			for k := id + 1; k < NrOfProcess; k++ {
				if fk := atomic.LoadInt32(&Flag[k]); fk == 2 || fk == 3 {
					processNeedToWait = true
					break
				}
			}
			if !processNeedToWait {
				keepWaiting = false
			} else {
				time.Sleep(1 * time.Millisecond)
				keepWaiting = true
			}
		}
		state = CriticalSection
		storeTrace()
		time.Sleep(2 * time.Millisecond)

		state = ExitProtocol
		storeTrace()
		atomic.StoreInt32(&Flag[id], 0)

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
