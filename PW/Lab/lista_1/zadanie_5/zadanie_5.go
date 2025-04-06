package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"unicode"
)

const (
	NrOfTravelers = 15
	MinSteps      = 10
	MaxSteps      = 100

	MinDelay = 10 * time.Millisecond
	MaxDelay = 50 * time.Millisecond

	BoardWidth  = 15
	BoardHeight = 15
)

var startTime = time.Now()
var wg sync.WaitGroup

var Board [BoardWidth][BoardHeight]chan struct{}

func initBoard() {
	for i := 0; i < BoardWidth; i++ {
		for j := 0; j < BoardHeight; j++ {
			Board[i][j] = make(chan struct{}, 1) //deklarujemy kanał 1
		}
	}
}

type Position struct {
	X int
	Y int
}

func (p *Position) MoveDown() {
	p.Y = (p.Y + 1) % BoardHeight
}

func (p *Position) MoveUp() {
	p.Y = (p.Y - 1 + BoardHeight) % BoardHeight
}

func (p *Position) MoveRight() {
	p.X = (p.X + 1) % BoardWidth
}

func (p *Position) MoveLeft() {
	p.X = (p.X - 1 + BoardWidth) % BoardWidth
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
	for i := 0; i < NrOfTravelers; i++ {
		traces := <-reportChannel
		PrintTraces(traces)
	}
	wg.Done()
}

type Traveler struct {
	Id         int
	Symbol     rune
	Position   Position
	Deadlocked bool
}

func acquireSquareWithTimeout(x, y int, timeout time.Duration) bool {
	select {
	case Board[x][y] <- struct{}{}: //kanał był pusty, teraz jest zajęty
		return true
	case <-time.After(timeout): //kanał był zajęty, nie udało się dostać w czasie 'timeout' -> deadlock
		return false
	}
}

func releaseSquare(x, y int) {
	<-Board[x][y] //uwalniamy kanał
}

func traveler(id int, sybol rune, seed int) {
	defer wg.Done()
	r := rand.New(rand.NewSource(int64(seed)))

	var isDeadlocked bool

	var traveler Traveler
	traveler.Id = id
	traveler.Symbol = sybol
	traveler.Position.X = id
	traveler.Position.Y = id

	var traces Traces_Sequence_Type
	traces.Last = -1

	Board[traveler.Position.X][traveler.Position.Y] <- struct{}{} //blokujemy kanał z board'a

	timeStamp := time.Since(startTime)
	traces.Last++
	traces.TraceArray[traces.Last] = TraceType{
		Time_Stamp: startTime.Add(timeStamp),
		Id:         traveler.Id,
		Position:   traveler.Position,
		Symbol:     traveler.Symbol,
	}

	nrOfSteps := MinSteps + r.Intn(MaxSteps-MinSteps+1)

	time.Sleep(100 * time.Millisecond)

	deadlockTimeOut := (MaxDelay + MinDelay) / 2
	for i := 0; i < nrOfSteps; i++ {
		delay := MinDelay + time.Duration(r.Int63n(int64(MaxDelay-MinDelay)))
		time.Sleep(delay)

		oldPos := traveler.Position

		if traveler.Id%2 == 0 {
			switch seed % 2 {
			case 0:
				traveler.Position.MoveDown()
			case 1:
				traveler.Position.MoveUp()
			}
		} else {
			switch seed % 2 {
			case 0:
				traveler.Position.MoveRight()
			case 1:
				traveler.Position.MoveLeft()
			}
		}

		newX := traveler.Position.X
		newY := traveler.Position.Y

		ok := acquireSquareWithTimeout(newX, newY, deadlockTimeOut)
		if !ok {
			isDeadlocked = true
			traveler.Symbol = unicode.ToLower(traveler.Symbol)

			traveler.Position = oldPos

			timeStamp = time.Since(startTime)
			traces.Last++
			traces.TraceArray[traces.Last] = TraceType{
				Time_Stamp: startTime.Add(timeStamp),
				Id:         traveler.Id,
				Position:   traveler.Position,
				Symbol:     traveler.Symbol,
			}

			releaseSquare(oldPos.X, oldPos.Y)
			reportChannel <- traces
			break

		}

		releaseSquare(oldPos.X, oldPos.Y)

		timeStamp = time.Since(startTime)
		traces.Last++
		traces.TraceArray[traces.Last] = TraceType{
			Time_Stamp: startTime.Add(timeStamp),
			Id:         traveler.Id,
			Position:   traveler.Position,
			Symbol:     traveler.Symbol,
		}
	}

	if !isDeadlocked {
		releaseSquare(traveler.Position.X, traveler.Position.Y)
		reportChannel <- traces
	}
}

func main() {
	fmt.Printf("-1 %d %d %d\n", NrOfTravelers, BoardWidth, BoardHeight)

	initBoard()
	wg.Add(1 + NrOfTravelers)

	go printer()

	symbols := []rune{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
		'I', 'J', 'K', 'L', 'M', 'N', 'O',
	}

	for i := 0; i < NrOfTravelers; i++ {
		go traveler(i, symbols[i], rand.Int())
	}

	wg.Wait()
	time.Sleep(3 * time.Second)
}
