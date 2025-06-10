import heapq
import random
import time

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

    # Tworzymy listę sąsiedztwa z wagami
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
  

def kruskal_algorithm(n, edges):
    parent = list(range(n))
    rank = [0] * n

    def find(x):
        if parent[x] != x:
            parent[x] = find(parent[x])
        return parent[x]

    def union(x, y):
        root_x = find(x)
        root_y = find(y)
        if root_x != root_y:
            if rank[root_x] > rank[root_y]:
                parent[root_y] = root_x
            elif rank[root_x] < rank[root_y]:
                parent[root_x] = root_y
            else:
                parent[root_y] = root_x
                rank[root_x] += 1

    edges.sort(key=lambda e: e.weight) #sortowanie po krawędziach według wagi
    mst_edges = []
    total_weight = 0.0

    for edge in edges:
        if find(edge.u) != find(edge.v):
            union(edge.u, edge.v)
            mst_edges.append((edge.u, edge.v, edge.weight))
            total_weight += edge.weight

    return mst_edges, total_weight




n_min = 1000
n_max = 10000
step = 100
rep = 3

for n in range(n_min, n_max + 1, step):
    for _ in range(rep):
        edges = generate_graph(n)
        
        start_time = time.time()
        prim_edges, prim_weight = prim_algorithm(n, edges)
        prim_time = time.time() - start_time

        start_time = time.time()
        kruskal_edges, kruskal_weight = kruskal_algorithm(n, edges)
        kruskal_time = time.time() - start_time

        print(f"n={n}, Prim: {prim_time:.4f}s, Kruskal: {kruskal_time:.4f}")