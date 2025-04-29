import pandas as pd
import matplotlib.pyplot as plt

def load_and_prepare(path):
    df = pd.read_csv(path, delimiter=",", skipinitialspace=True)
    df.columns = [col.strip() for col in df.columns]
    df['n'] = pd.to_numeric(df['n'], errors='coerce')
    df['s'] = pd.to_numeric(df['s'], errors='coerce')
    return df.groupby('n')['s'].mean().reset_index()

# wczytywanie
quick_median = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_3/scripts/quick_sort.txt")
quick_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_3/scripts/1quick_sort.txt")
dp_quick_pivot = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_3/scripts/dp_quick_sort.txt")
dp_quick_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_3/scripts/1dp_quick_sort.txt")


# wykres


plt.plot(quick_median['n'], quick_median['s'], marker='o', label='Quick Sort MoM')
plt.plot(quick_avg['n'], quick_avg['s'], marker='o', label='Quick Sort')
plt.plot(dp_quick_pivot['n'], dp_quick_pivot['s'], marker='o', label='DP Quick Sort MoM')
plt.plot(dp_quick_avg['n'], dp_quick_avg['s'], marker='o', label='DP Quick Sort')


plt.title("Avarage number of swaps based on n (for k = 10)")
plt.xlabel("n ")
plt.ylabel("avg s")
plt.legend()
plt.grid(True)
plt.show()
