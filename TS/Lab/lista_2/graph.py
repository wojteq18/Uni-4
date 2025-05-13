import networkx as nx
import matplotlib.pyplot as plt
import random

# --- Parametry potrzebne dla generate_topology ---
NUM_VERTICES = 20
MAX_EDGES = 29 # Użyj tej samej wartości co w głównym skrypcie

# --- Skopiowana funkcja generate_topology z Twojego skryptu ---
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

    # Weryfikacja (wypisywanie informacji)
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
                 break

    print(f"Ostateczna topologia: |V|={G.number_of_nodes()}, |E|={G.number_of_edges()}")
    if G.number_of_edges() >= 30:
         print(f"Ostrzeżenie: Liczba krawędzi ({G.number_of_edges()}) nie jest mniejsza niż 30.")

    return G

# --- Generowanie topologii ---
print("Generowanie topologii do wizualizacji...")
# Ustawienie ziarna losowości dla powtarzalności (opcjonalne, ale pomocne)
# random.seed(42) # Możesz odkomentować i użyć stałego ziarna
G_to_visualize = generate_topology(NUM_VERTICES, MAX_EDGES)

# --- Wizualizacja ---
plt.figure(figsize=(12, 10)) # Możesz dostosować rozmiar obrazka

# Wybór algorytmu rozmieszczania węzłów (layoutu)
# 'spring_layout' jest popularny, ale możesz wypróbować inne:
# circular_layout, kamada_kawai_layout, random_layout
pos = nx.spring_layout(G_to_visualize, seed=42) # Użycie ziarna dla powtarzalnego layoutu

nx.draw_networkx_nodes(G_to_visualize, pos, node_size=400, node_color='lightblue', edgecolors='black')
nx.draw_networkx_edges(G_to_visualize, pos, width=1.5, alpha=0.7, edge_color='gray')
nx.draw_networkx_labels(G_to_visualize, pos, font_size=10, font_weight='bold')

plt.title(f"Wizualizacja Topologii Sieci (|V|={G_to_visualize.number_of_nodes()}, |E|={G_to_visualize.number_of_edges()})", fontsize=16)
# Można ukryć osie, bo w layoucie grafu nie mają one fizycznego znaczenia
plt.axis('off')

# --- Zapis do pliku ---
image_filename = "topology_visualization.png"
plt.savefig(image_filename, bbox_inches='tight') # bbox_inches='tight' pomaga przyciąć puste marginesy
print(f"Zapisano obraz topologii jako: {image_filename}")

# --- Wyświetlenie na ekranie ---
plt.show()