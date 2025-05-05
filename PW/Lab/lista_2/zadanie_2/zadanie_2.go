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

var numRenters = 0

var startTime = time.Now()

var wg sync.WaitGroup

type Message struct {
	Type   string
	Answer chan interface{} //interfejs może przechowywać dowolną wartość
}

type Field struct {
	Entry chan Message
}

var Board [BoardWidth][BoardHeight]*Field
var renterPosition sync.Map
var renterSymbols sync.Map

func initBoard() {
	for i := 0; i < BoardWidth; i++ {
		for j := 0; j < BoardHeight; j++ {
			p := &Field{Entry: make(chan Message, 1)}
			Board[i][j] = p
			go initField(p)
		}
	}
}

func initField(p *Field) {
	isFree := true
	isRenterHere := false
	for {
		switch msg := <-p.Entry; msg.Type {
		case "enter":
			if isFree {
				msg.Answer <- true
				isFree = false
			} else {
				msg.Answer <- false
			}

		case "enterRenter":
			if isFree {
				msg.Answer <- true
				isFree = false
				isRenterHere = true
			} else {
				msg.Answer <- false
			}
		case "exit":
			isFree = true

		case "exitRenter":
			isFree = true
			isRenterHere = false

		case "status":
			if isRenterHere {
				msg.Answer <- true
			} else {
				msg.Answer <- false
			}
		}
	}
}

func EnterField(x, y int) bool {
	answer := make(chan interface{})
	msg := Message{
		Type:   "enter",
		Answer: answer,
	}

	Board[x][y].Entry <- msg
	res := <-answer
	return res.(bool)
}

func EnterFieldRenter(x, y int) bool {
	answer := make(chan interface{})
	msg := Message{
		Type:   "enterRenter",
		Answer: answer,
	}

	Board[x][y].Entry <- msg
	res := <-answer
	return res.(bool)
}

func ExitField(x, y int) {
	msg := Message{
		Type: "exit",
	}
	Board[x][y].Entry <- msg
}

func ExitFieldRenter(x, y int) {
	msg := Message{
		Type: "exitRenter",
	}
	Board[x][y].Entry <- msg
}

func checkStatus(x, y int) bool {
	msg := Message{
		Type: "status",
	}
	answer := make(chan interface{})
	msg.Answer = answer
	Board[x][y].Entry <- msg
	res := <-answer
	return res.(bool)
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

var reportChannel = make(chan Traces_Sequence_Type, NrOfTravelers+50)

func printer(done chan struct{}) {
	defer close(done)
	for traces := range reportChannel {
		PrintTraces(traces)
	}
}

type Traveler struct {
	Id         int
	Symbol     rune
	Position   Position
	Deadlocked bool
}

type Renter struct {
	Id       int
	Symbol   rune
	Position Position
}

func wildRenter(Id int, symbol rune, seed int) {
	defer wg.Done()
	r := rand.New(rand.NewSource(int64(seed)))
	var traces Traces_Sequence_Type
	traces.Last = -1
	var renter Renter
	renter.Id = Id
	renter.Symbol = symbol

	for {
		renter.Position.X = r.Intn(BoardWidth)
		renter.Position.Y = r.Intn(BoardHeight)
		isAlright := EnterFieldRenter(renter.Position.X, renter.Position.Y)
		if isAlright {
			renterPosition.Store(renter.Id, renter.Position)
			renterSymbols.Store(renter.Id, renter.Symbol)
			break
		}
	}

	traces.Last++
	traces.TraceArray[traces.Last] = TraceType{
		Time_Stamp: time.Now(),
		Id:         renter.Id,
		Position:   renter.Position,
		Symbol:     renter.Symbol,
	}

	randTime := time.Duration(r.Intn(601)) * time.Millisecond
	time.Sleep(randTime)

	ExitFieldRenter(renter.Position.X, renter.Position.Y)
	renterPosition.Delete(renter.Id)
	renterSymbols.Delete(renter.Id)

	traces.Last++
	traces.TraceArray[traces.Last] = TraceType{
		Time_Stamp: time.Now(),
		Id:         renter.Id,
		Position:   renter.Position,
		Symbol:     '.',
	}
	reportChannel <- traces
}

func traveler(id int, sybol rune, seed int) {

	defer wg.Done()
	r := rand.New(rand.NewSource(int64(seed)))

	var traveler Traveler
	traveler.Id = id
	traveler.Symbol = sybol
	traveler.Position.X = r.Intn(BoardWidth)
	traveler.Position.Y = r.Intn(BoardHeight)

	var traces Traces_Sequence_Type
	traces.Last = -1

	for {
		isAlright := EnterField(traveler.Position.X, traveler.Position.Y)
		if isAlright {
			break
		} else {
			traveler.Position.X = r.Intn(BoardWidth)
			traveler.Position.Y = r.Intn(BoardHeight)
		}
	}
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

	deadlockTimeOut := (9*MaxDelay + MinDelay) / 2
	isDeadlock := false

	for i := 0; i < nrOfSteps; i++ {
		delay := MinDelay + time.Duration(r.Int63n(int64(MaxDelay-MinDelay)))
		time.Sleep(delay)

		oldPos := traveler.Position

		switch r.Intn(4) {
		case 0:
			traveler.Position.MoveUp()
		case 1:
			traveler.Position.MoveDown()
		case 2:
			traveler.Position.MoveLeft()
		case 3:
			traveler.Position.MoveRight()
		}

		newX := traveler.Position.X
		newY := traveler.Position.Y
		ok := false
		deadLockDeadLine := time.Now().Add(deadlockTimeOut)

		for time.Now().Before(deadLockDeadLine) {
			ok = EnterField(newX, newY)
			if ok {
				break
			}
			isRenter := checkStatus(newX, newY)
			switch isRenter {
			case true:
				var renterId int = -1
				renterPosition.Range(func(key, value interface{}) bool {
					pos := value.(Position)
					if pos.X == newX && pos.Y == newY {
						renterId = key.(int)
						return false //znaleziony renter, koniec iteracji
					}
					return true
				})
				if renterId != -1 {
					sym, okSym := renterSymbols.Load(renterId)
					if !okSym {
						sym = '?' // jeśli symbol się nie znajdzie
					}
					renter := Renter{
						Id:       renterId,
						Position: Position{X: newX, Y: newY},
						Symbol:   sym.(rune),
					}

					neighborPositions := []Position{
						{X: renter.Position.X, Y: (renter.Position.Y + 1) % BoardHeight},
						{X: renter.Position.X, Y: (renter.Position.Y - 1 + BoardHeight) % BoardHeight},
						{X: (renter.Position.X + 1) % BoardWidth, Y: renter.Position.Y},
						{X: (renter.Position.X - 1 + BoardWidth) % BoardWidth, Y: renter.Position.Y},
					}

					moved := false
					for _, newRenterPos := range neighborPositions {
						//fmt.Println("Wejscie w rentera")
						if EnterFieldRenter(newRenterPos.X, newRenterPos.Y) {
							ExitFieldRenter(renter.Position.X, renter.Position.Y)
							renterPosition.Store(renter.Id, newRenterPos)
							renter.Position = newRenterPos

							trace := Traces_Sequence_Type{
								Last: 0,
								TraceArray: TraceArray{
									0: TraceType{
										Time_Stamp: time.Now(),
										Id:         renter.Id,
										Position:   renter.Position,
										Symbol:     renter.Symbol,
									},
								},
							}
							reportChannel <- trace

							moved = true
							break
						}
					}

					if !moved {
						switch r.Intn(4) {
						case 0:
							traveler.Position.MoveUp()
						case 1:
							traveler.Position.MoveDown()
						case 2:
							traveler.Position.MoveLeft()
						case 3:
							traveler.Position.MoveRight()
						}
						newX = traveler.Position.X
						newY = traveler.Position.Y
					}

					// Traveler zajmuje pole rentera
					ok = EnterField(newX, newY)
					if ok {
						timeStamp = time.Since(startTime)
						traces.Last++
						traces.TraceArray[traces.Last] = TraceType{
							Time_Stamp: startTime.Add(timeStamp),
							Id:         traveler.Id,
							Position:   traveler.Position,
							Symbol:     traveler.Symbol,
						}
					}

				}
			case false:
				continue
			}

		}
		if !ok {
			isDeadlock = true
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

			//releaseSquare(oldPos.X, oldPos.Y)
			reportChannel <- traces
			break

		}

		ExitField(oldPos.X, oldPos.Y)

		timeStamp = time.Since(startTime)
		traces.Last++
		traces.TraceArray[traces.Last] = TraceType{
			Time_Stamp: startTime.Add(timeStamp),
			Id:         traveler.Id,
			Position:   traveler.Position,
			Symbol:     traveler.Symbol,
		}
	}

	if !isDeadlock {
		ExitField(traveler.Position.X, traveler.Position.Y)
		reportChannel <- traces
	}
}

/*func renter(id int, seed int) {
	r := rand.New(rand.NewSource(int64(seed)))

}*/

func main() {
	fmt.Printf("-1 %d %d %d\n", NrOfTravelers+500, BoardWidth, BoardHeight)

	initBoard()

	symbols := []rune{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
		'I', 'J', 'K', 'L', 'M', 'N', 'O',
	}

	renterSymbols := []rune{
		'1', '2', '3', '4', '5', '6', '7', '8', '9',
	}

	//reportChannel = make(chan Traces_Sequence_Type, NrOfTravelers+50)

	printerDone := make(chan struct{})
	go printer(printerDone)

	wg.Add(1)
	go func() {
		defer wg.Done()
		limit := 90
		for i := 0; i < limit; i++ {
			if rand.Float64() < 0.4 {
				renterIndex := rand.Intn(len(renterSymbols))
				wg.Add(1)
				go wildRenter(NrOfTravelers+numRenters+i, renterSymbols[renterIndex], rand.Int())
				numRenters++
			}
			randTime := MinDelay + time.Duration(rand.Int63n(int64(MaxDelay-MinDelay)))
			time.Sleep(randTime)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < NrOfTravelers; i++ {
			wg.Add(1)
			go traveler(i+1, symbols[i], rand.Int())
		}
	}()

	wg.Wait()
	close(reportChannel)

	<-printerDone
}
