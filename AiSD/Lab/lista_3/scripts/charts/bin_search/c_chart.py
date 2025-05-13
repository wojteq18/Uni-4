import pandas as pd
import matplotlib.pyplot as plt

def load_and_prepare(path):
    df = pd.read_csv(path, delimiter=",", skipinitialspace=True)
    df.columns = [col.strip() for col in df.columns]
    df['n'] = pd.to_numeric(df['n'], errors='coerce')
    df['c'] = pd.to_numeric(df['c'], errors='coerce')
    return df.groupby('n')['c'].mean().reset_index()

# wczytywanie
bin_search = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_3/scripts/bin_search.txt")
#random_select_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_3/scripts/random_select.txt")
#dp_quick_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/random_sequence_scripts/dp_quick_sort.txt")
#insertion_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/random_sequence_scripts/insertion_sort.txt")

# wykres
plt.figure(figsize=(10, 6))

plt.plot(bin_search['n'], bin_search['c'], marker='o', label='BS')
#plt.plot(random_select_avg['n'], random_select_avg['c'], marker='o', label='Random Select')
#plt.plot(dp_quick_avg['n'], dp_quick_avg['c'], marker='o', label='DP Quick Sort')
#plt.plot(insertion_avg['n'], insertion_avg['c'], marker='o', label='Insertion Sort')



#plt.plot(insertion_avg['n'], insertion_avg['c'], marker='o', label='Insertion Sort')


plt.title("Time, k = 2")
plt.xlabel("n ")
plt.ylabel("avg c")
plt.yscale('log')
plt.legend()
plt.grid(True)
plt.show()
