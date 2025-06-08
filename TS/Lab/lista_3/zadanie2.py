import random
from dataclasses import dataclass # Skraca definicje klsay nośnej (Wave) - automatcznie generuje __init__, __repr__ itd.
from typing import List, Optional, Dict # Dokumnetuje, jakie struktury danych trzymamy

@dataclass
class Wave:
    pos: int    # pozycja fali
    val: int    # >0: ID węzła, -1: fala kolizji
    ttl: int    # czas życia fali

class Medium:
    def __init__(self, size: int):
        self.size = size
        self.waves: List[Wave] = [] # Lista aktywnych fal w danym momencie

    def propagate(self, transmit_positions: List[tuple]):
        # 1. Propagacja istniejących fal
        candidates: List[Wave] = []
        for w in self.waves:
            if w.ttl > 1:
                # Propagacja do samej siebie i sąsiadów
                candidates.append(Wave(w.pos, w.val, w.ttl - 1))
                if w.pos > 0:
                    candidates.append(Wave(w.pos - 1, w.val, w.ttl - 1))
                if w.pos < self.size - 1:
                    candidates.append(Wave(w.pos + 1, w.val, w.ttl - 1))

        # 2. Wstrzykiwanie nowych fal TYLKO w pozycji węzła
        for node_id, pos, ttl in transmit_positions:
            candidates.append(Wave(pos, node_id, ttl))

        # 3. Grupowanie fal i wykrywanie kolizji
        grouped: Dict[int, List[Wave]] = {}
        for w in candidates:
            grouped.setdefault(w.pos, []).append(w)

        new_waves: List[Wave] = []
        for pos, waves in grouped.items():
            vals = {w.val for w in waves if w.val > 0}
            has_collision = any(w.val == -1 for w in waves)
            best_ttl = max(w.ttl for w in waves)

            if has_collision or len(vals) > 1:
                new_waves.append(Wave(pos, -1, best_ttl))
            elif len(vals) == 1:
                new_waves.append(Wave(pos, next(iter(vals)), best_ttl))

        self.waves = new_waves # Kończy propagację na jeden krok symulacji

    def get_cells(self) -> List[int]: #0 = cisza, >0 = sygnał od danego węzła, -1 = kolizja
        cells = [0] * self.size
        for w in self.waves:
            cells[w.pos] = w.val
        return cells

    def __str__(self):
        return ' '.join(f'{x:2}' for x in self.get_cells())

class Node: # Reprezentacje ojedynczego węzła sieciowego
    def __init__(self, node_id: int, position: int, medium: Medium,
                 transmit_prob: float = 0.3,
                 transmit_duration: Optional[int] = None):
        self.id = node_id # Unikalne id węzła
        self.position = position # indeks w medium
        self.medium = medium # referencja do wspólnego medium
        self.transmit_prob = transmit_prob # szansa na transmisję w stanie idle
        self.transmit_duration = transmit_duration or 2 * (medium.size - 1)

        self.state = 'idle' # state nalezy do {idle, transitting, waiting}
        self.transmit_timer = 0 #Liczy pozostałe kroki nadawania
        self.collision_count = 0 # Ile kolizji dotychczas w tej transmisji
        self.backoff_timer = 0 # Ile kroków jeszcze czeka węzeł przed ponowną próbą
        self.new_transmission = False  # Flaga nowej transmisji

    def sense(self) -> bool: #Czy medium pod wezlem jest wolne
        return self.medium.get_cells()[self.position] == 0

    def attempt_transmit(self): # Gdy wezel jest idle i medium wolne - z prawd. transmit_prob przechodzi do stanu transmitting
        if self.state == 'idle' and self.sense():
            if random.random() < self.transmit_prob:
                self.state = 'transmitting'
                self.transmit_timer = self.transmit_duration
                self.new_transmission = True  # Oznacz nową transmisję

    def handle_collision(self):
        self.collision_count += 1
        self.state = 'waiting'
        k = min(self.collision_count, 10)
        self.backoff_timer = random.randint(0, 2**k - 1)
        self.transmit_timer = 0
        print(f"  ⚠️  Node {self.id}: collision! backoff={self.backoff_timer}")

    def update(self) -> Optional[str]: # Gdy dojdzie do 0, wracamy do idle
        if self.state == 'waiting':
            self.backoff_timer -= 1
            if self.backoff_timer <= 0:
                self.state = 'idle'
            return None

        if self.state == 'transmitting': # Odczytujemy co pozycj "słýszy" - jesli to nie nasze echo (id) ani cisza to wykryto kozlije
            seen = self.medium.get_cells()[self.position]
            if seen != self.id and seen != 0:
                self.handle_collision()
                return 'collision'
            
            self.transmit_timer -= 1
            if self.transmit_timer <= 0:
                print(f"  ✅  Node {self.id}: transmission completed")
                self.state = 'idle'
                self.collision_count = 0
                return 'success'
        return None

class Simulation:
    def __init__(self, num_steps: int, medium_size: int,
                 node_positions: Optional[List[int]] = None,
                 seed: Optional[int] = None):
        if seed is not None:
            random.seed(seed)
            
        self.medium = Medium(medium_size)
        node_positions = node_positions or [2, medium_size - 2]
        self.nodes = [
            Node(i+1, pos, self.medium)
            for i, pos in enumerate(node_positions)
        ]
        self.num_steps = num_steps
        self.stats = {'collisions': 0, 'successes': 0}

    def step(self, step_num: int):
        # Krok 1: Nowe transmisje
        tx_positions = []
        for node in self.nodes:
            node.attempt_transmit()
            if node.new_transmission:
                tx_positions.append((node.id, node.position, node.transmit_duration))
                node.new_transmission = False

        # Krok 2: Propagacja w medium
        self.medium.propagate(tx_positions)

        # Krok 3: Aktualizacja stanu węzłów
        collisions = 0
        successes = 0
        for node in self.nodes:
            result = node.update()
            if result == 'collision':
                collisions += 1
            elif result == 'success':
                successes += 1
        
        self.stats['collisions'] += collisions
        self.stats['successes'] += successes

        # Logowanie
        print(f"\n--- Krok {step_num} ---")
        print("Medium:", self.medium)
        for node in self.nodes:
            state = f"{node.state} (backoff: {node.backoff_timer})" if node.state == 'waiting' else node.state
            print(f"  Node {node.id}@{node.position}: {state}")

    def run(self):
        print("=== Rozpoczęcie symulacji CSMA/CD ===")
        print(f"Medium (rozmiar {self.medium.size}):")
        print("Węzły:", ', '.join(f"Node {n.id}@{n.position}" for n in self.nodes))
        print("=" * 40)
        
        for step in range(self.num_steps):
            self.step(step)
        
        print("\n=== Podsumowanie ===")
        print(f"Wykryte kolizje: {self.stats['collisions']}")
        print(f"Sukcesy transmisji: {self.stats['successes']}")

def main():
    sim = Simulation(
        num_steps=36,
        medium_size=5,
        node_positions=[2, 4],
        seed=421
    )
    sim.run()

if __name__ == "__main__":
    main()