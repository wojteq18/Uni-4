import pandas as pd
import matplotlib.pyplot as plt

def load_and_prepare(path):
    df = pd.read_csv(path, delimiter=",", skipinitialspace=True)
    df.columns = [col.strip() for col in df.columns]
    df['n'] = pd.to_numeric(df['n'], errors='coerce')
    df['s'] = pd.to_numeric(df['s'], errors='coerce')
    grouped = df.groupby('n')['s'].mean().reset_index()
    grouped['s_per_n'] = grouped['s'] / grouped['n']  # <<< tutaj s / n
    return grouped

# wczytywanie
hybrid_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/increasing_sequence_scripts/dp_quick_sort.txt")
quick_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/increasing_sequence_scripts/quick_sort.txt")
#insertion_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/increasing_sequence_scripts/insertion_sort.txt")

# wykres
plt.figure(figsize=(10, 6))
#plt.yscale('log')  # jeśli różnice są duże

# poprawne wartości s/n
plt.plot(hybrid_avg['n'], hybrid_avg['s'] / hybrid_avg['n'], marker='o', label='DP Quick Sort (s/n)')
plt.plot(quick_avg['n'], quick_avg['s'] / quick_avg['n'], marker='o', label='Quick Sort (s/n)')

plt.title("Average number of swaps per element (s/n) based on n")
plt.xlabel("n")
plt.ylabel("avg swaps per element (s/n)")
plt.legend()
plt.grid(True)
plt.show()
