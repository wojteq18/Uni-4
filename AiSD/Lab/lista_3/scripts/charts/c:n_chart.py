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
select_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_3/scripts/select.txt")
random_select_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_3/scripts/random_select.txt")
#dp_quick_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/random_sequence_scripts/dp_quick_sort.txt")
#insertion_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/random_sequence_scripts/insertion_sort.txt")

# wykres
plt.figure(figsize=(10, 6))
#plt.yscale('log')


plt.plot(select_avg['n'], select_avg['c_per_n'], marker='o', label='Select')
plt.plot(random_select_avg['n'], random_select_avg['c_per_n'], marker='o', label='Random Select')
#lt.plot(dp_quick_avg['n'], dp_quick_avg['c_per_n'], marker='o', label='DP Quick Sort')
#plt.plot(insertion_avg['n'], insertion_avg['c_per_n'], marker='o', label='Insertion Sort')

plt.title("Average number of comparisons per element (c/n), m = 5")
plt.xlabel("n")
plt.ylabel("avg c / n")
plt.legend()
plt.grid(True)
plt.show()
