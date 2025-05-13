import networkx as nx
import numpy as np
import random
import matplotlib.pyplot as plt
import math

# --- Parametry globalne ---
NUM_VERTICES = 20
MAX_EDGES = 29  # Mniej niż 30
AVG_PACKET_SIZE_BITS = 1000 * 8 # Średnia wielkość pakietu w bitach (np. 1 kB)
PROB_EDGE_OK = 0.95           # Prawdopodobieństwo p, że krawędź działa
T_MAX = 0.1                   # Maksymalne dopuszczalne średnie opóźnienie (w sekundach)
NUM_SIMULATIONS = 5000        # Liczba symulacji Monte Carlo do oszacowania niezawodności

# --- Krok 1: Generowanie Topologii G ---
def generate_topology(num_vertices, max_edges):
    """Generuje graf G=(V,E) spełniający warunki |V|, |E| i braku izolowanych wierzchołków."""
    if max_edges < num_vertices - 1:
        raise ValueError("Maksymalna liczba krawędzi musi być co najmniej |V|-1, aby graf mógł być spójny.")

    G = nx.Graph()
    G.add_nodes_from(range(num_vertices))

    # Zapewnienie braku izolowanych wierzchołków przez stworzenie ścieżki
    edges_added = 0
    for i in range(num_vertices - 1):
        G.add_edge(i, i + 1)
        edges_added += 1

    # Dodawanie losowych krawędzi, aż do osiągnięcia limitu lub blisko niego
    possible_edges = list(nx.non_edges(G))
    random.shuffle(possible_edges)

    edges_to_add = max_edges - edges_added
    added_count = 0
    for u, v in possible_edges:
        if added_count >= edges_to_add:
            break
        if not G.has_edge(u, v):
             # Sprawdzenie czy dodanie krawędzi nie przekroczy max_edges
            if G.number_of_edges() < max_edges:
                G.add_edge(u, v)
                added_count += 1
            else:
                break # Osiągnęliśmy już maksymalną liczbę krawędzi

    # Weryfikacja
    print(f"Wygenerowana topologia: |V|={G.number_of_nodes()}, |E|={G.number_of_edges()}")
    isolated = list(nx.isolates(G))
    if isolated:
        print(f"Ostrzeżenie: Znaleziono izolowane wierzchołki: {isolated}. Poprawka może być potrzebna.")
        # Prosta poprawka - połącz izolowany wierzchołek z losowym innym
        for node in isolated:
             if G.number_of_edges() < max_edges:
                 target = random.choice([n for n in G.nodes() if n != node])
                 G.add_edge(node, target)
                 print(f"Dodano krawędź ({node}, {target}) aby usunąć izolację.")
             else:
                 print(f"Nie można dodać krawędzi dla {node} - osiągnięto limit {max_edges}")
                 break # Przerwij, jeśli nie możemy już dodawać krawędzi

    print(f"Ostateczna topologia: |V|={G.number_of_nodes()}, |E|={G.number_of_edges()}")
    if G.number_of_edges() >= 30:
         print(f"Ostrzeżenie: Liczba krawędzi ({G.number_of_edges()}) nie jest mniejsza niż 30.")

    return G

# --- Krok 2: Generowanie Macierzy Natężeń N ---
def generate_traffic_matrix(num_vertices, max_packets_per_pair=10):
    """Generuje macierz N=[n(i,j)] z losowymi natężeniami."""
    N = np.zeros((num_vertices, num_vertices))
    for i in range(num_vertices):
        for j in range(num_vertices):
            if i != j:
                # Losowa liczba pakietów (można użyć innych rozkładów)
                N[i, j] = random.randint(0, max_packets_per_pair)
    print(f"Wygenerowano macierz N. Całkowity ruch G = {np.sum(N):.2f} pakietów/s")
    return N

# --- Krok 3: Obliczanie Przepływów 'a(e)' ---
# Założenie: Routing po najkrótszej ścieżce (pod względem liczby krawędzi)
def calculate_flows(G, N):
    """Oblicza przepływ a(e) dla każdej krawędzi, realizując macierz N przez najkrótsze ścieżki."""
    flows = {edge: 0.0 for edge in G.edges()}
    num_vertices = G.number_of_nodes()

    for i in range(num_vertices):
        for j in range(num_vertices):
            if i != j and N[i, j] > 0:
                try:
                    # Znajdź najkrótszą ścieżkę
                    path = nx.shortest_path(G, source=i, target=j)
                    # Dodaj przepływ N[i,j] do każdej krawędzi na ścieżce
                    for k in range(len(path) - 1):
                        u, v = path[k], path[k+1]
                        # Klucz krawędzi w 'flows' jest krotką (min(u,v), max(u,v))
                        edge = tuple(sorted((u, v)))
                        if edge in flows:
                            flows[edge] += N[i, j]
                        else:
                             # To nie powinno się zdarzyć w grafie nieskierowanym, ale dla pewności
                            print(f"Ostrzeżenie: Krawędź {edge} ze ścieżki nie znaleziona w G.edges()")

                except nx.NetworkXNoPath:
                    print(f"Ostrzeżenie: Brak ścieżki między {i} a {j} dla N[{i},{j}]={N[i, j]}. Ruch nie może być zrealizowany.")
                    # W praktyce może to oznaczać problem z topologią lub potrzebę innego routingu
                    return None # Zwracamy None, sygnalizując problem

    print("Obliczono przepływy a(e) na krawędziach.")
    return flows

# --- Krok 4: Ustalanie Przepustowości 'c(e)' ---
def set_capacities(flows, m, min_capacity_margin=1.1):
    """Ustawia przepustowość c(e) tak, aby c(e)/m > a(e)."""
    capacities = {}
    if flows is None:
        return None

    for edge, flow_a in flows.items():
        # Przepustowość w pakietach/s musi być większa niż przepływ a(e)
        # c(e)/m > a(e) => c(e) > a(e) * m
        # Ustawiamy c(e)/m = a(e) * margin => c(e) = a(e) * m * margin
        # Dodajmy też małą stałą, aby uniknąć zerowej przepustowości gdy a(e)=0
        required_capacity_packets_sec = flow_a + 1 # Dodajemy 1 pakiet/s marginesu minimalnego
        capacity_packets_sec = required_capacity_packets_sec * min_capacity_margin
        capacities[edge] = capacity_packets_sec * m # c(e) w bitach/s

    print("Ustalono przepustowości c(e) dla krawędzi.")
    return capacities

# --- Krok 5: Obliczanie Średniego Opóźnienia T ---
def calculate_T(G_current, flows, capacities, N, m):
    """Oblicza średnie opóźnienie T dla danego (potencjalnie uszkodzonego) grafu G_current."""
    total_traffic_G = np.sum(N)
    if total_traffic_G == 0:
        return 0.0 # Brak ruchu, brak opóźnienia

    sum_delay_terms = 0.0
    # Sumujemy tylko po krawędziach istniejących w G_current
    for edge in G_current.edges():
        # Upewnijmy się, że krawędź jest w naszym formacie (u,v) gdzie u<v
        edge_key = tuple(sorted(edge))

        if edge_key not in flows or edge_key not in capacities:
            # To nie powinno się zdarzyć, jeśli G_current jest podgrafem G
            print(f"Ostrzeżenie: Brak danych o przepływie/przepustowości dla krawędzi {edge_key} w G_current.")
            continue

        a_e = flows[edge_key]
        c_e = capacities[edge_key]

        capacity_packets_sec = c_e / m
        denominator = capacity_packets_sec - a_e

        # Warunek c(e)/m > a(e) powinien zapobiec <= 0, ale sprawdzamy dla pewności
        if denominator <= 1e-9: # Używamy małej tolerancji zamiast ścisłego > 0
             # Jeśli warunek nie jest spełniony (np. przez błędy zaokrągleń lub zbyt mały margines)
             # lub jeśli krawędź jest przeciążona, opóźnienie dąży do nieskończoności.
             # W symulacji oznacza to, że sieć jest niestabilna/nieniezawodna.
             # print(f"Krawędź {edge_key}: c/m ({capacity_packets_sec:.2f}) <= a ({a_e:.2f}). Opóźnienie -> inf.")
             return float('inf') # Zwracamy nieskończoność

        delay_term = a_e / denominator
        sum_delay_terms += delay_term

    T = (1 / total_traffic_G) * sum_delay_terms
    return T

# --- Krok 6: Szacowanie Niezawodności (Symulacja Monte Carlo) ---
def estimate_reliability(G_orig, N, flows, capacities, m, p_ok, t_max, num_sim):
    """Szacuje niezawodność P(T < T_max) przez symulację Monte Carlo."""
    successful_simulations = 0
    total_traffic_pairs = np.sum(N > 0) # Liczba par (i,j) z niezerowym ruchem

    if not flows or not capacities:
         print("Brak danych o przepływach lub przepustowościach. Nie można oszacować niezawodności.")
         return 0.0

    for _ in range(num_sim):
        # 1. Generuj stan sieci (usuń krawędzie z prawd. 1-p_ok)
        G_current = G_orig.copy()
        edges_to_remove = []
        for edge in G_orig.edges():
            if random.random() > p_ok: # Prawdopodobieństwo 1-p_ok, że krawędź ulegnie awarii
                edges_to_remove.append(edge)
        G_current.remove_edges_from(edges_to_remove)

        # 2. Sprawdź, czy sieć nadal obsługuje wymagany ruch N
        all_paths_exist = True
        num_vertices = G_orig.number_of_nodes()
        for i in range(num_vertices):
            for j in range(num_vertices):
                if N[i, j] > 0:
                    if not nx.has_path(G_current, source=i, target=j):
                        all_paths_exist = False
                        break
            if not all_paths_exist:
                break

        if not all_paths_exist:
            # Sieć rozspójniona dla potrzeb ruchu N -> nieniezawodna w tej symulacji
            continue # Przejdź do następnej symulacji

        # 3. Oblicz T dla działającej sieci
        # Ważne: Używamy ORYGINALNYCH a(e) i c(e) dla krawędzi, które PRZETRWAŁY
        current_flows = {edge: flows[tuple(sorted(edge))] for edge in G_current.edges()}
        current_capacities = {edge: capacities[tuple(sorted(edge))] for edge in G_current.edges()}

        # Sprawdzenie warunku c(e)/m > a(e) dla działających krawędzi (czy nie zostały przeciążone)
        stable = True
        for edge in G_current.edges():
             edge_key = tuple(sorted(edge))
             if capacities[edge_key] / m <= flows[edge_key] + 1e-9:
                  stable = False
                  break
        if not stable:
             # Jeśli jakaś działająca krawędź jest przeciążona, T -> inf
             continue # Nieniezawodna

        T = calculate_T(G_current, flows, capacities, N, m) # Używamy oryginalnych flows/caps

        # 4. Sprawdź warunek T < T_max
        if T < t_max:
            successful_simulations += 1

    # 6. Oblicz oszacowaną niezawodność
    reliability = successful_simulations / num_sim
    return reliability

# --- Główny program ---
if __name__ == "__main__":
    # --- Inicjalizacja modelu ---
    print("--- Inicjalizacja modelu sieci ---")
    G_initial = generate_topology(NUM_VERTICES, MAX_EDGES)
    N_initial = generate_traffic_matrix(NUM_VERTICES, max_packets_per_pair=5) # Mniejszy ruch na start
    flows_initial = calculate_flows(G_initial, N_initial)
    capacities_initial = set_capacities(flows_initial, AVG_PACKET_SIZE_BITS, min_capacity_margin=1.2) # Margines 20%

    if flows_initial is None or capacities_initial is None:
        print("Błąd inicjalizacji modelu. Zakończenie.")
        exit()

    # --- Początkowa niezawodność ---
    print(f"\n--- Obliczanie początkowej niezawodności (p={PROB_EDGE_OK}, T_max={T_MAX}) ---")
    initial_reliability = estimate_reliability(G_initial, N_initial, flows_initial, capacities_initial, AVG_PACKET_SIZE_BITS, PROB_EDGE_OK, T_MAX, NUM_SIMULATIONS)
    print(f"Oszacowana początkowa niezawodność: {initial_reliability:.4f}")


    # --- Analiza Parametryczna ---

    print("\n--- Analiza: Wpływ zwiększania natężeń N ---")
    traffic_factors = np.linspace(1.0, 3.0, 10) # Zwiększamy ruch od 1x do 3x
    reliability_vs_N = []
    for factor in traffic_factors:
        print(f"Analiza dla N * {factor:.2f}")
        N_current = N_initial * factor
        # UWAGA: Zwiększenie N wymaga ponownego obliczenia a(e) i potencjalnie c(e)
        flows_current = calculate_flows(G_initial, N_current)
        # Ustalmy przepustowości na nowo, aby *nadal* spełniały warunek c/m > a
        # To ważne założenie - zakładamy, że możemy dostosować sieć (zwiększyć c) do rosnącego ruchu
        capacities_current = set_capacities(flows_current, AVG_PACKET_SIZE_BITS, min_capacity_margin=1.2)
        if flows_current is None or capacities_current is None:
             print("Nie można obliczyć przepływów/przepustowości dla zwiększonego ruchu. Pomijanie.")
             reliability_vs_N.append(0.0) # Traktujemy jako awarię
             continue

        # A co jeśli *nie* dostosowujemy C? Wtedy warunek c/m > a może paść.
        # Sprawdźmy to: użyjmy ORYGINALNYCH przepustowości
        # To bardziej realistyczny scenariusz krótkoterminowy
        rel = estimate_reliability(G_initial, N_current, flows_current, capacities_initial, AVG_PACKET_SIZE_BITS, PROB_EDGE_OK, T_MAX, NUM_SIMULATIONS // 5) # Mniej symulacji dla szybkości
        reliability_vs_N.append(rel)
        print(f"  Niezawodność (stałe C): {rel:.4f}")

    plt.figure()
    plt.plot(traffic_factors, reliability_vs_N, marker='o')
    plt.title("Niezawodność vs Mnożnik Natężeń N (stałe C)")
    plt.xlabel("Mnożnik Natężeń N")
    plt.ylabel("Oszacowana Niezawodność P(T < T_max)")
    plt.grid(True)
    #plt.ylim(0, 1.1) # Ustawienie zakresu osi Y

    print("\n--- Analiza: Wpływ zwiększania przepustowości C ---")
    capacity_factors = np.linspace(1.0, 3.0, 10) # Zwiększamy C od 1x do 3x
    reliability_vs_C = []
    for factor in capacity_factors:
        print(f"Analiza dla C * {factor:.2f}")
        capacities_current = {edge: c * factor for edge, c in capacities_initial.items()}
        # N i a(e) pozostają bez zmian (N_initial, flows_initial)
        rel = estimate_reliability(G_initial, N_initial, flows_initial, capacities_current, AVG_PACKET_SIZE_BITS, PROB_EDGE_OK, T_MAX, NUM_SIMULATIONS // 5)
        reliability_vs_C.append(rel)
        print(f"  Niezawodność: {rel:.4f}")

    plt.figure()
    plt.plot(capacity_factors, reliability_vs_C, marker='s', color='green')
    plt.title("Niezawodność vs Mnożnik Przepustowości C")
    plt.xlabel("Mnożnik Przepustowości C")
    plt.ylabel("Oszacowana Niezawodność P(T < T_max)")
    plt.grid(True)
    #plt.ylim(0, 1.1)

    print("\n--- Analiza: Wpływ dodawania krawędzi ---")
    num_edges_to_add = 5 # Ile krawędzi chcemy dodać
    reliability_vs_Edges = [initial_reliability] # Zaczynamy od początkowej
    G_evolving = G_initial.copy()
    edges_added_count_list = [G_evolving.number_of_edges()]

    # Oblicz średnią przepustowość istniejących krawędzi
    avg_capacity = np.mean(list(capacities_initial.values())) if capacities_initial else 1e6 # Jakaś domyślna wartość

    # Przechowajmy początkowe przepływy i przepustowości
    current_flows = flows_initial.copy()
    current_capacities = capacities_initial.copy()

    for i in range(num_edges_to_add):
        print(f"Analiza po dodaniu {i+1} krawędzi...")
        # Znajdź parę wierzchołków bez krawędzi
        possible_new_edges = list(nx.non_edges(G_evolving))
        if not possible_new_edges:
            print("Brak możliwości dodania nowych krawędzi (graf pełny?). Zakończenie analizy topologii.")
            break
        u, v = random.choice(possible_new_edges)
        G_evolving.add_edge(u, v)
        new_edge_key = tuple(sorted((u, v)))
        print(f"Dodano krawędź: {new_edge_key}")

        # Dodaj przepustowość dla nowej krawędzi
        current_capacities[new_edge_key] = avg_capacity
        # Dodaj zerowy przepływ początkowy dla nowej krawędzi
        current_flows[new_edge_key] = 0.0

        # UWAGA: Dodanie krawędzi MOŻE zmienić najkrótsze ścieżki, a więc i przepływy a(e)
        # Trzeba by przeliczyć przepływy od nowa dla G_evolving i N_initial
        print("Przeliczanie przepływów po dodaniu krawędzi...")
        flows_new_topo = calculate_flows(G_evolving, N_initial)

        if flows_new_topo is None:
             print("Błąd podczas przeliczania przepływów dla nowej topologii. Pomijanie.")
             # Można by spróbować dodać inną krawędź, lub przerwać
             G_evolving.remove_edge(u, v) # Cofnij dodanie krawędzi
             del current_capacities[new_edge_key]
             continue # Przejdź do następnej iteracji (może inna losowa krawędź zadziała)


        # Upewnijmy się, że wszystkie krawędzie G_evolving są w flows_new_topo
        for edge_evolved in G_evolving.edges():
             edge_evolved_key = tuple(sorted(edge_evolved))
             if edge_evolved_key not in flows_new_topo:
                  flows_new_topo[edge_evolved_key] = 0.0 # Dodaj z zerowym przepływem, jeśli jakimś cudem brakuje

        # Teraz mamy nowe przepływy. Czy przepustowości są wystarczające?
        # Sprawdźmy warunek c/m > a dla WSZYSTKICH krawędzi w G_evolving
        # używając current_capacities i flows_new_topo
        capacities_ok = True
        temp_capacities_for_check = current_capacities.copy() # Użyj bieżących C
        for edge_key_check, flow_val in flows_new_topo.items():
             if edge_key_check not in temp_capacities_for_check:
                  # Jeśli nowa krawędź pojawiła się w przepływach, a nie ma jej w C (nie powinno się zdarzyć)
                  temp_capacities_for_check[edge_key_check] = avg_capacity # Dodaj ze średnią

             if temp_capacities_for_check[edge_key_check] / AVG_PACKET_SIZE_BITS <= flow_val + 1e-9:
                  print(f"  Ostrzeżenie: Po dodaniu krawędzi i rerutingu, krawędź {edge_key_check} stała się przeciążona (c/m={temp_capacities_for_check[edge_key_check]/AVG_PACKET_SIZE_BITS:.2f} <= a={flow_val:.2f}). Zwiększanie C...")
                  # Zwiększmy C dla tej krawędzi, aby spełnić warunek
                  new_c = (flow_val + 1) * 1.2 * AVG_PACKET_SIZE_BITS
                  current_capacities[edge_key_check] = new_c # Zaktualizuj w głównym słowniku
                  print(f"  Nowa przepustowość dla {edge_key_check}: {new_c / (1024*1024):.2f} Mbps")


        # Po ewentualnych korektach C, obliczamy niezawodność
        rel = estimate_reliability(G_evolving, N_initial, flows_new_topo, current_capacities, AVG_PACKET_SIZE_BITS, PROB_EDGE_OK, T_MAX, NUM_SIMULATIONS // 5)
        reliability_vs_Edges.append(rel)
        edges_added_count_list.append(G_evolving.number_of_edges())
        print(f"  Niezawodność (|E|={G_evolving.number_of_edges()}): {rel:.4f}")
        # Zaktualizuj przepływy na potrzeby następnej iteracji
        current_flows = flows_new_topo


    plt.figure()
    plt.plot(edges_added_count_list, reliability_vs_Edges, marker='^', color='red')
    plt.title("Niezawodność vs Liczba Krawędzi")
    plt.xlabel("Liczba Krawędzi w Sieci")
    plt.ylabel("Oszacowana Niezawodność P(T < T_max)")
    plt.grid(True)
    #plt.ylim(0, 1.1)

    plt.tight_layout()
    plt.show() # Pokaż wszystkie wykresy