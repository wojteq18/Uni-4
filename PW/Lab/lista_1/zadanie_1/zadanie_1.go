package main

import (
	"fmt"
	"math/rand"
	"time"
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
	fmt.Printf("%v %d %d %d %c\n", t.Time_Stamp.Sub(startTime), t.Id, t.Position.X, t.Position.Y, t.Symbol)
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

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Hello, World!")
}
