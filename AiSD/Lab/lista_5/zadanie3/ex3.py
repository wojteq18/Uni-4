import random

class Node:

    def __init__(self, key):
        self.key = key
        self.parent = None
        self.child = None
        self.sibling = None
        self.degree = 0

class BinomialHeap:

    def __init__(self):
        self.head = None

    def is_empty(self):
        return self.head is None

    def _link_trees(self, y, z): # Łączy dwa drzewa dwumiaowe, o jedny stopniu
        # Zakładamy, że y.key <= z.key
        if y.key > z.key:
            y, z = z, y # Zamiana, aby w y był mniejszy klucz
            
        z.parent = y
        z.sibling = y.child
        y.child = z
        y.degree += 1
        
    def union(self, other_heap):
        cmp = 0

        new_heap = BinomialHeap()
        new_heap.head = self._merge_root_lists(self, other_heap)

        # Zwolnienie pamięci starych kopców (już nie są potrzebne)
        self.head = None
        other_heap.head = None

        if new_heap.head is None:
            return new_heap, cmp

        prev_x = None
        x = new_heap.head
        next_x = x.sibling

        while next_x is not None:
            # Przypadek 1 i 2: Stopnie są różne lub jest więcej niż 2 drzew o tym samym stopniu
            if (x.degree != next_x.degree) or \
               (next_x.sibling is not None and next_x.sibling.degree == x.degree):
                prev_x = x
                x = next_x
            # Przypadek 3 i 4: Dwa kolejne drzewa mają ten sam stopień
            else:
                cmp += 1
                if x.key <= next_x.key:
                    # Łączymy next_x z xf
                    x.sibling = next_x.sibling
                    self._link_trees(x, next_x)
                else:
                    # Łączymy x z next_x
                    if prev_x is None:
                        new_heap.head = next_x
                    else:
                        prev_x.sibling = next_x
                    self._link_trees(next_x, x)
                    x = next_x # Kontynuujemy od nowego korzenia
            
            next_x = x.sibling
            
        return new_heap, cmp

    def _merge_root_lists(self, h1, h2):

        if h1.head is None:
            return h2.head
        if h2.head is None:
            return h1.head

        new_head = None
        
        h1_ptr, h2_ptr = h1.head, h2.head

        # Ustawienie głowy nowej listy
        if h1_ptr.degree <= h2_ptr.degree:
            new_head = h1_ptr
            h1_ptr = h1_ptr.sibling
        else:
            new_head = h2_ptr
            h2_ptr = h2_ptr.sibling

        current = new_head
        
        # Scalanie pozostałych elementów
        while h1_ptr is not None and h2_ptr is not None:
            if h1_ptr.degree <= h2_ptr.degree:
                current.sibling = h1_ptr
                h1_ptr = h1_ptr.sibling
            else:
                current.sibling = h2_ptr
                h2_ptr = h2_ptr.sibling
            current = current.sibling

        # Dołączenie reszty z jednej z list
        if h1_ptr is not None:
            current.sibling = h1_ptr
        else:
            current.sibling = h2_ptr
            
        return new_head
        

    def insert(self, key):

        temp_heap = BinomialHeap()
        temp_heap.head = Node(key)
        
        # Zastąpienie bieżącego kopca wynikiem unii
        # `self` jest niszczony wewnątrz `union`
        new_heap, comparsion = self.union(temp_heap)
        self.head = new_heap.head
        return comparsion


    def get_minimum(self):

        if self.head is None:
            return None

        cmp = 0
        min_node = self.head
        current = self.head.sibling
        while current is not None:
            cmp += 1
            if current.key < min_node.key:
                min_node = current
            current = current.sibling
        return min_node

    def extract_min(self):

        if self.head is None:
            return None

        cmp = 0
        # 1. Znajdź korzeń x z minimalnym kluczem
        min_node = None
        min_prev = None
        current = self.head
        prev = None
        
        min_key = float('inf')

        while current is not None:
            cmp += 1
            if current.key < min_key:
                min_key = current.key
                min_node = current
                min_prev = prev
            prev = current
            current = current.sibling
            
        # 2. Usuń x z listy korzeni kopca
        if min_prev is None: # Minimalny był pierwszym elementem
            self.head = min_node.sibling
        else:
            min_prev.sibling = min_node.sibling

        # 3. Utwórz nowy kopiec z dzieci usuniętego węzła
        child_heap = BinomialHeap()
        
        # Dzieci są połączone w kolejności malejących stopni.
        # Musimy odwrócić tę listę, aby stworzyć poprawny kopiec.
        child_head = min_node.child
        
        # Odwracanie listy dzieci
        reversed_child_head = None
        current = child_head
        while current is not None:
            next_node = current.sibling
            current.sibling = reversed_child_head
            current.parent = None # Dzieci stają się korzeniami
            reversed_child_head = current
            current = next_node
        
        child_heap.head = reversed_child_head

        # 4. Połącz (union) bieżący kopiec z kopcem dzieci
        new_heap, ucmp = self.union(child_heap)
        self.head = new_heap.head

        total_cmp = cmp + ucmp
        
        return min_node, total_cmp
    
n = 500
h_1 = BinomialHeap()
h_2 = BinomialHeap()

n_min = 100
n_max = 10000
n_step = 100

for n in range(n_min, n_max + 1, n_step):
    cmp = 0
    h_1 = BinomialHeap()
    h_2 = BinomialHeap()


    for i in range(n):
        h1_cmp = h_1.insert(random.randint(1, n**2))
        cmp += h1_cmp

    # Insert n elements into h_2
    for i in range(n):
        h2_cmp = h_2.insert(random.randint(1, n**2))
        cmp += h2_cmp
    # Union of h_1 and h_2
    h_3, cmp1 = h_1.union(h_2)
    cmp += cmp1

    # Extract min from the unioned heap
    while h_3.head is not None:
        min_node, comparsions = h_3.extract_min()
        cmp += comparsions
    print(f"n={n}, comparisons: {cmp / n:.2f}")    

