import pandas as pd
import matplotlib.pyplot as plt

def load_and_prepare(path):
    df = pd.read_csv(path, delimiter=",", skipinitialspace=True)
    df.columns = [col.strip() for col in df.columns]
    df['n'] = pd.to_numeric(df['n'], errors='coerce')
    df['c'] = pd.to_numeric(df['c'], errors='coerce')
    grouped = df.groupby('n')['c'].mean().reset_index()
    grouped = grouped[grouped['c'] > 0]  # <<< usuwamy problematyczne zera
    return grouped

# wczytanie
hybrid_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/increasing_sequence_scripts/dp_quick_sort.txt")
quick_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/increasing_sequence_scripts/quick_sort.txt")
#insertion_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/increasing_sequence_scripts/insertion_sort.txt")

# wykres
plt.figure(figsize=(10, 6))
plt.yscale('log')

plt.plot(hybrid_avg['n'], hybrid_avg['c'] + 0.5, marker='o', label='DP Quick Sort')
plt.plot(quick_avg['n'], quick_avg['c'], marker='o', label='Quick Sort')
#plt.plot(insertion_avg['n'], insertion_avg['c'], marker='o', label='Insertion Sort')

plt.title("Average number of comparisons based on n (for k = 10)")
plt.xlabel("n")
plt.ylabel("avg c")
plt.legend()
plt.grid(True)
plt.show()
