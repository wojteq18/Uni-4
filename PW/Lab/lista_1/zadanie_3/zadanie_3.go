package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	NrOfTravelers = 26
	MinSteps      = 10
	MaxSteps      = 100

	MinDelay = 10 * time.Millisecond
	MaxDelay = 50 * time.Millisecond

	BoardWidth  = 9
	BoardHeight = 9
)

var startTime = time.Now()

var Board [BoardWidth][BoardHeight]bool

func occupy(x, y int) {
	Board[x][y] = true
}

func release(x, y int) {
	Board[x][y] = false
}

func isOccupied(x, y int) bool {
	return Board[x][y]
}

func deadlock(x, y, nextX, previousX, nextY, previousY int) bool {
	return isOccupied(x, nextY) && isOccupied(x, previousY) && isOccupied(nextX, y) && isOccupied(previousX, y)
}

type Position struct {
	X int
	Y int
}

func (p *Position) MoveDown() {
	release(p.X, p.Y)
	p.Y = (p.Y + 1) % BoardHeight
	occupy(p.X, p.Y)
}

func (p *Position) MoveUp() {
	release(p.X, p.Y)
	p.Y = (p.Y - 1 + BoardHeight) % BoardHeight
	occupy(p.X, p.Y)
}

func (p *Position) MoveRight() {
	release(p.X, p.Y)
	p.X = (p.X + 1) % BoardWidth
	occupy(p.X, p.Y)
}

func (p *Position) MoveLeft() {
	release(p.X, p.Y)
	p.X = (p.X - 1 + BoardWidth) % BoardWidth
	occupy(p.X, p.Y)
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
}

type Traveler struct {
	Id       int
	Symbol   rune
	Position Position
}

func traveler(id int, sybol rune, seed int) {
	r := rand.New(rand.NewSource(int64(seed)))
	var isDedlocked bool = false
	var deadlockStart time.Time

	var traveler Traveler
	traveler.Id = id
	traveler.Symbol = sybol
	traveler.Position.X = r.Intn(BoardWidth)
	traveler.Position.Y = r.Intn(BoardHeight)

	var traces Traces_Sequence_Type
	traces.Last = -1

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

	for i := 0; i < nrOfSteps; i++ {
		delay := MinDelay + time.Duration(r.Int63n(int64(MaxDelay-MinDelay)))
		time.Sleep(delay)

		if isDedlocked {
			timeStamp = time.Since(startTime)
			traces.Last++
			traces.TraceArray[traces.Last] = TraceType{
				Time_Stamp: startTime.Add(timeStamp),
				Id:         traveler.Id,
				Position:   traveler.Position,
				Symbol:     traveler.Symbol,
			}
			continue
		}

		nextY := (traveler.Position.Y - 1 + BoardHeight) % BoardHeight
		previousY := (traveler.Position.Y + 1) % BoardHeight

		nextX := (traveler.Position.X + 1) % BoardWidth
		previousX := (traveler.Position.X - 1 + BoardWidth) % BoardWidth

		if deadlock(traveler.Position.X, traveler.Position.Y, nextX, previousX, nextY, previousY) {
			if deadlockStart.IsZero() {
				deadlockStart = time.Now()
			} else if time.Since(deadlockStart) > MaxDelay {
				traveler.Symbol = 'x'
				isDedlocked = true
			}
			continue
		} else {
			deadlockStart = time.Time{}
		}

		switch r.Intn(4) {
		case 0:
			if isOccupied(traveler.Position.X, nextY) {
				continue
			} else {
				traveler.Position.MoveUp()
			}
		case 1:
			if isOccupied(traveler.Position.X, previousY) {
				continue
			} else {
				traveler.Position.MoveDown()
			}
		case 2:
			if isOccupied(nextX, traveler.Position.Y) {
				continue
			} else {
				traveler.Position.MoveRight()
			}
		case 3:
			if isOccupied(previousX, traveler.Position.Y) {
				continue
			} else {
				traveler.Position.MoveLeft()
			}
		}

		timeStamp = time.Since(startTime)
		traces.Last++
		traces.TraceArray[traces.Last] = TraceType{
			Time_Stamp: startTime.Add(timeStamp),
			Id:         traveler.Id,
			Position:   traveler.Position,
			Symbol:     traveler.Symbol,
		}
	}
	reportChannel <- traces
}

func main() {
	fmt.Printf("-1 %d %d %d\n", NrOfTravelers, BoardWidth, BoardHeight)

	go printer()

	symbols := []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

	for i := 0; i < NrOfTravelers; i++ {
		go traveler(i, symbols[i], rand.Int())
	}

	time.Sleep(3 * time.Second)
}
