import pandas as pd
import matplotlib.pyplot as plt

def load_and_prepare(path):
    df = pd.read_csv(path, delimiter=",", skipinitialspace=True)
    df.columns = [col.strip() for col in df.columns]
    df['n'] = pd.to_numeric(df['n'], errors='coerce')
    df['s'] = pd.to_numeric(df['s'], errors='coerce')
    return df.groupby('n')['s'].mean().reset_index()

# wczytywanie
select_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_3/scripts/select.txt")
random_select_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_3/scripts/random_select.txt")
#dp_quick_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/random_sequence_scripts/dp_quick_sort.txt")
#insertion_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/random_sequence_scripts/insertion_sort.txt")


# wykres


plt.plot(select_avg['n'], select_avg['s'], marker='o', label='Select')
plt.plot(random_select_avg['n'], random_select_avg['s'], marker='o', label='Random Select')
#plt.plot(dp_quick_avg['n'], dp_quick_avg['s'], marker='o', label='DP Quick Sort')
#plt.plot(insertion_avg['n'], insertion_avg['s'], marker='o', label='Insertion Sort')


plt.title("Avarage number of swaps based on n (for k = 10)")
plt.xlabel("n ")
plt.ylabel("avg s")
plt.legend()
plt.grid(True)
plt.show()
