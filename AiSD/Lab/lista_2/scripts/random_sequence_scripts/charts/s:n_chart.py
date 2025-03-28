import matplotlib.pyplot as plt

def wczytaj_dane(plik):
    n, s, c = [], [], []
    with open(plik, 'r') as f:
        next(f)  
        for linia in f:
            dane = linia.strip().split(',')
            if len(dane) == 3:
                n.append(int(dane[0]))
                c.append(int(dane[1]))
                s.append(int(dane[2]))
    return n, c, s

n1, c1, s1 = wczytaj_dane('/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/random_sequence_scripts/insertion_sort.txt')
n2, c2, s2 = wczytaj_dane('/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/random_sequence_scripts/quick_sort.txt')
n3, c3, s3 = wczytaj_dane('/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/random_sequence_scripts/hybrid_sort.txt')

# pojedynczy wykres s/n
plt.plot(n1, [s / n for s, n in zip(s1, n1)], label='Insertion Sort', marker='o')
plt.plot(n2, [s / n for s, n in zip(s2, n2)], label='Quick Sort', marker='s')
plt.plot(n3, [s / n for s, n in zip(s3, n3)], label='Hybrid Sort', marker='^')

plt.xlabel('elements (n)')
plt.ylabel('s / n')
plt.title('Swaps per element')
plt.legend()
plt.grid(True)
plt.show()
