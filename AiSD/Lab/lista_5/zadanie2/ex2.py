import heapq
import random
import time
from collections import defaultdict

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

def build_adjecency_list(edges):
    graph = defaultdict(list) #Tworzy specjalny słownik, gdzie klucze 
    for u, v, _ in edges:
        graph[u].append(v)
        graph[v].append(u)
    return graph    


graph = generate_graph(10)  
mst_edges, _ = kruskal_algorithm(10, graph)

graph = build_adjecency_list(mst_edges)

for i in range (len(graph)):
    print(f"Node {i}: ", end="")
    for neighbor in graph[i]:
        print(f"{neighbor} ", end="")
    print()