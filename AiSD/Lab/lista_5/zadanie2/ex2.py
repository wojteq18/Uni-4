import random
from collections import defaultdict
import heapq

class Edge:
    def __init__(self, u, v, weight):
        self.u = u
        self.v = v
        self.weight = weight

def generate_graph(n):
    edges = []
    for i in range(n):
        for j in range(i + 1, n):
            weight = random.random()
            edges.append(Edge(i, j, weight))
    return edges        

def prim_algorithm(n, edges):

    graph = [[] for _ in range(n)]
    for edge in edges:
        graph[edge.u].append((edge.v, edge.weight))
        graph[edge.v].append((edge.u, edge.weight))

    cost = [float('inf')] * n
    prev = [None] * n
    in_mst = [False] * n
    cost[0] = 0.0  # dowolny wierzchołek startowy

    # Kolejka priorytetowa: (waga, wierzchołek)
    heap = [(0.0, 0)]
    total_weight = 0.0
    mst_edges = []

    while heap:
        curr_cost, u = heapq.heappop(heap)
        if in_mst[u]:
            continue
        in_mst[u] = True
        total_weight += curr_cost

        if prev[u] is not None:
            mst_edges.append((prev[u], u, curr_cost))

        for v, w in graph[u]:
            if not in_mst[v] and w < cost[v]:
                cost[v] = w
                prev[v] = u
                heapq.heappush(heap, (w, v))

    return mst_edges, total_weight

def build_adjecency_list(edges):
    graph = defaultdict(list) #Tworzy specjalny słownik, gdzie klucze 
    for u, v, _ in edges:
        graph[u].append(v)
        graph[v].append(u)
    return graph    

def distance(node, parent, graph, time):
    times = []
    for neighbor in graph[node]:
        if neighbor != parent:
            distance(neighbor, node, graph, time)
            times.append(time[neighbor])

    times.sort(reverse=True)
    max_time = 0
    for i, t in enumerate(times): #dla liści ta pętla sie nie wykonuje
        max_time = max(max_time, t + i + 1) #t -> czas do poinformowania swojego całego poddrzewa, i -> numer kolejki w posortowanej liście

    time[node] = max_time  

def simulation(graph, root):
    time = {}
    distance(root, -1, graph, time)
    return time[root]          




n_min = 100
n_max = 5000
step = 100
rep = 4

for n in range(n_min, n_max + 1, step):
    max_value = 0
    min_value = float('inf')
    for _ in range(rep):
        edges = generate_graph(n)
        st_edges, _ = prim_algorithm(n, edges)
        graph = build_adjecency_list(st_edges)
        random_root = random.randint(0, n - 1) #losowy wierzchołek jako korzeń
        rounds = simulation(graph, random_root)
        print(f"n={n}, rounds: {rounds}, root: {random_root}")
        if max_value < rounds:
            max_value = rounds

        if min_value > rounds:
            min_value = rounds

    print(f"n={n}, max rounds: {max_value}")        
    print(f"n={n}, min rounds: {min_value}")
