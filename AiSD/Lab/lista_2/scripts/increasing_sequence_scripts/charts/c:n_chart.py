import pandas as pd
import matplotlib.pyplot as plt

def load_and_prepare(path):
    df = pd.read_csv(path, delimiter=",", skipinitialspace=True)
    df.columns = [col.strip() for col in df.columns]
    df['n'] = pd.to_numeric(df['n'], errors='coerce')
    df['c'] = pd.to_numeric(df['c'], errors='coerce')  # <--- teraz bierzemy 'c'
    grouped = df.groupby('n')['c'].mean().reset_index()
    grouped['c_per_n'] = grouped['c'] / grouped['n']  # <<< c / n
    return grouped

# wczytywanie
hybrid_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/increasing_sequence_scripts/dp_quick_sort.txt")
quick_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/increasing_sequence_scripts/quick_sort.txt")
#insertion_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/increasing_sequence_scripts/insertion_sort.txt")

# wykres
plt.figure(figsize=(10, 6))


plt.plot(hybrid_avg['n'], hybrid_avg['c_per_n'] + 0.2 , marker='o', label='DP Quick Sort')
plt.plot(quick_avg['n'], quick_avg['c_per_n'], marker='o', label='Quick Sort')
#plt.plot(insertion_avg['n'], insertion_avg['c_per_n'], marker='o', label='Insertion Sort')

plt.title("Average number of comparisons per element (c/n) based on n (for k = 10)")
plt.xlabel("n")
plt.ylabel("avg c / n")
plt.legend()
plt.grid(True)
plt.show()
