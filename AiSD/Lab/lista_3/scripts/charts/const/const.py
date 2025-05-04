import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

def load_and_prepare(path):
    df = pd.read_csv(path, delimiter=",", skipinitialspace=True)
    df.columns = [col.strip() for col in df.columns]
    df['n'] = pd.to_numeric(df['n'], errors='coerce')
    df['c'] = pd.to_numeric(df['c'], errors='coerce')
    df = df.groupby('n')['c'].mean().reset_index()
    df['c_over_log_n'] = df['c'] / np.log2(df['n'])
    return df

# wczytywanie danych
bin_search = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_3/scripts/bin_search.txt")

# wykres c / log n
plt.figure(figsize=(10, 6))

plt.plot(bin_search['n'], bin_search['c_over_log_n'], marker='o', label='Select (c / log₂(n))')

plt.title("Average number of comparisons per log₂(n)")
plt.xlabel("n")
plt.ylabel("avg c / log₂(n)")
plt.legend()
plt.grid(True)
plt.show()
