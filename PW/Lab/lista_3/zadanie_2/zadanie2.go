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

	MinSteps = 50
	MaxSteps = 100

	MinDelay = 10 * time.Millisecond
	MaxDelay = 50 * time.Millisecond

	BoardWidth  = NrOfProcess
	BoardHeight = 4
)

var ticketMutex sync.Mutex
var maxTicket int32 = 0

// zmienne potrzebne do algorytmu Piekarnianego
var Flag [NrOfProcess]int32
var Number [NrOfProcess]int32

func findMaxAtomic(arr []int32) int32 {
	maxVal := int32(0)
	for i := range arr {
		val := atomic.LoadInt32(&arr[i])
		if val > maxVal {
			maxVal = val
		}
	}
	return maxVal
}

type ProcessState int

const (
	LocalSection ProcessState = iota
	EntryProtocol
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

var reportChannel = make(chan Traces_Sequence_Type, NrOfProcess)

func printer() {
	defer WaitGroup.Done()
	for i := 0; i < NrOfProcess; i++ {
		traces := <-reportChannel
		PrintTraces(traces)
	}

	fmt.Printf("-1 %d %d %d ", NrOfProcess, BoardWidth, BoardHeight)

	states := []string{"Local_Section", "Entry_Protocol", "Critical_Section", "Exit_Protocol"}
	for _, state := range states {
		fmt.Printf("%s;", state)
	}

	fmt.Println("EXTRA_LABEL = ;", maxTicket)
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

	for i := 0; i < (nrOfSteps/4)-1; i++ {
		//Local Section
		state = LocalSection
		storeTrace()
		time.Sleep(MinDelay + time.Duration(r.Int63n(int64(MaxDelay-MinDelay))))

		//Entry Protocol
		state = EntryProtocol
		storeTrace()
		atomic.StoreInt32(&Flag[id], 1)
		ticketMutex.Lock()
		newTicket := findMaxAtomic(Number[:]) + 1
		atomic.StoreInt32(&Number[id], newTicket)
		if newTicket > maxTicket {
			atomic.StoreInt32(&maxTicket, newTicket)
		}
		ticketMutex.Unlock()
		atomic.StoreInt32(&Flag[id], 0)

		for j := 0; j < NrOfProcess; j++ {
			if j != id {
				// Czekaj, aż proces j skończy wybierać swój numer
				for atomic.LoadInt32(&Flag[j]) == 1 {
					fmt.Println("Process", id, "waiting for process", j, "to finish choosing number")
					time.Sleep(1 * time.Millisecond)
				}

				// Czekaj tylko jeśli proces j próbuje wejść do sekcji krytycznej (Number[j] > 0)
				// i ma niższy numer lub równy, ale niższy ID
				for atomic.LoadInt32(&Number[j]) != 0 &&
					((atomic.LoadInt32(&Number[j]) < atomic.LoadInt32(&Number[id])) ||
						(atomic.LoadInt32(&Number[j]) == atomic.LoadInt32(&Number[id]) && j < id)) {
					time.Sleep(10 * time.Millisecond)
				}
			}
		}

		//Critical Section
		state = CriticalSection
		storeTrace()
		time.Sleep(MinDelay + time.Duration(r.Int63n(int64(MaxDelay-MinDelay))))

		//Exit Protocol
		state = ExitProtocol
		storeTrace()
		atomic.StoreInt32(&Flag[id], 0)
		atomic.StoreInt32(&Number[id], 0)
	}

	if state != LocalSection {
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
