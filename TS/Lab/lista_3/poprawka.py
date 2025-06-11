from __future__ import annotations
import random
from enum import Enum, auto
from dataclasses import dataclass
from typing import List, Optional, Tuple


MEDIUM_LENGTH: int = 80
MAX_ATTEMPTS:  int = 10


class NodeState(Enum):
    IDLE         = auto()
    TRANSMITTING = auto()
    BACKOFF      = auto()
    JAMMING      = auto()
    SUCCESS      = auto()

@dataclass
class Node:
    name: str
    position: int
    transmission_tick: int = -1
    state: NodeState = NodeState.IDLE
    backoff: int = 0
    attempts: int = 0
    jam_start: int = 0

    def __post_init__(self) -> None:
        if self.transmission_tick < 0:
            self.transmission_tick = random.randint(0, 1000)

    def __str__(self) -> str:
        return f"{self.name}(x={self.position})"

@dataclass
class Signal:
    pos: int
    direction: int
    source: str
    created: int

Cell = Optional[Tuple[str, int]]


def propagate_signal(
    medium: List[Cell],
    signals: List[Signal],
    tick: int
) -> List[Signal]:
    new_signals: List[Signal] = []
    for sig in signals:
        new_pos = sig.pos + sig.direction
        if 0 <= new_pos < MEDIUM_LENGTH: #sprawdza czy pozycja sygnału mieści się w medium
            if medium[new_pos] is None: #pole puste - wpisujemy swój sygnał
                medium[new_pos] = (sig.source, 0)
            elif medium[new_pos][0] != sig.source: #ktoś tam był - kolizja
                medium[new_pos] = ('x', 0)
            new_signals.append(
                Signal(new_pos, sig.direction, sig.source, sig.created) #dodajemy nowy sygnał do listy propagujących
            )
    return new_signals


def is_medium_idle(medium: List[Cell], pos: int) -> bool:
    return medium[pos] is None


def run_simulation(nodes: List[Node]) -> None:
    medium:  List[Cell]   = [None] * MEDIUM_LENGTH
    signals: List[Signal] = [] #lista aktywnych sygnałów
    log:     List[str]    = []

    success = 0
    tick    = 0 #licznik kroków czasowych

    while success < len(nodes): #dopóki wszystkie nody nie osiągną sukcesu
        tick += 1
        # linia do logu
        line = [' '] * MEDIUM_LENGTH

        # propagacja sygnałów
        signals = propagate_signal(medium, signals, tick) #jeli pole jest wolne, zostawia swój source, jeśli jest zajęte przez innego nadawce: x

        # obsługa każdego węzła
        for node in nodes:
            # nowa próba wysyłki
            if tick == node.transmission_tick and node.state == NodeState.IDLE: #jesli medium jest wolne, to zaczyna nadawać
                if is_medium_idle(medium, node.position):
                    medium[node.position] = (node.name, 0)
                    signals.append(Signal(node.position, -1, node.name, tick))
                    signals.append(Signal(node.position,  1, node.name, tick))
                    node.state = NodeState.TRANSMITTING
                else:
                    # kanał zajęty – przesuwamy zaplanowany tick o 1
                    node.transmission_tick += 1

            elif node.state == NodeState.TRANSMITTING:
                cell = medium[node.position]
                if cell is not None and cell[0] != node.name: #sprawdza, czy w jego pozycji w medium pojawił się inny sygnał
                    # kolizja!
                    node.state = NodeState.JAMMING
                    if node.attempts < MAX_ATTEMPTS:
                        factor = 2 ** node.attempts #node.attempts to licznik dotychczasowych kolizji kabla
                        node.backoff = 2 * MEDIUM_LENGTH * random.randint(0, factor - 1)
                        node.transmission_tick = (
                            tick + 2 * MEDIUM_LENGTH + 1 + node.backoff
                        )
                        print(factor)
                    elif node.attempts > MAX_ATTEMPTS + 6:
                        node.state = NodeState.IDLE
                        print(
                            f"Node {node.name} failed after {node.attempts} attempts."
                        )
                    node.attempts  += 1
                    node.jam_start  = tick
                    medium[node.position] = ('x', 0)
                else:
                    # wysyłamy dalej lub uznajemy sukces
                    if tick >= node.transmission_tick + 2 * MEDIUM_LENGTH: #jeżeli wysyłał sygnał przez ten czas, to uznaje sukces
                        node.state = NodeState.SUCCESS
                        success   += 1
                        medium[node.position] = (node.name, 0)
                    else:
                        medium[node.position] = (node.name, 0)
                        signals.append(Signal(node.position, -1, node.name, tick))
                        signals.append(Signal(node.position,  1, node.name, tick))

            elif node.state == NodeState.JAMMING:
                if tick >= node.jam_start + 2 * MEDIUM_LENGTH:
                    node.state = NodeState.IDLE #po określonym czasie wraca do idle
                else:
                    medium[node.position] = ('x', 0)
                    signals.append(Signal(node.position, -1, 'x', tick))
                    signals.append(Signal(node.position,  1, 'x', tick))

        # zapisujemy widok medium do loga
        for i in range(MEDIUM_LENGTH):
            if medium[i] is not None:
                line[i] = medium[i][0]
        log.append(''.join(line))

        medium[:] = [None] * MEDIUM_LENGTH

    with open("output.txt", "w", encoding="utf-8") as fout:
        for l in log:
            fout.write(l + "\n")

    for l in log:
        print(l)


def main() -> None:
    rng = random.Random()
    n        = rng.randint(3, 3)         
    nodes: List[Node] = []

    for i in range(n):
        position = rng.randint(0, MEDIUM_LENGTH - 1)
        tx_tick  = rng.randint(0, 100) #kiedy planuje zacząć nadawanie
        nodes.append(Node(chr(ord('a') + i), position, tx_tick))

    run_simulation(nodes)


if __name__ == "__main__":
    main()
