import pandas as pd
import matplotlib.pyplot as plt

def load_and_prepare(path):
    df = pd.read_csv(path, delimiter=",", skipinitialspace=True)
    df.columns = [col.strip() for col in df.columns]
    df['n'] = pd.to_numeric(df['n'], errors='coerce')
    df['c'] = pd.to_numeric(df['s'], errors='coerce')
    return df.groupby('n')['c'].mean().reset_index()

# wczytywanie
select_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_3/scripts/select.txt")
select3_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_3/scripts/select3.txt")
select7_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_3/scripts/select7.txt")
select9_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_3/scripts/select9.txt")

#random_select_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_3/scripts/random_select.txt")
#dp_quick_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/random_sequence_scripts/dp_quick_sort.txt")
#insertion_avg = load_and_prepare("/home/wojteq18/Uni/AiSD/Lab/lista_2/scripts/random_sequence_scripts/insertion_sort.txt")


# wykres


plt.plot(select_avg['n'], select_avg['c'], marker='o', label='Select')
plt.plot(select3_avg['n'], select3_avg['c'], marker='o', label='Select 3')
plt.plot(select7_avg['n'], select7_avg['c'], marker='o', label='Select 7')
plt.plot(select9_avg['n'], select9_avg['c'], marker='o', label='Select 9')

#plt.plot(dp_quick_avg['n'], dp_quick_avg['s'], marker='o', label='DP Quick Sort')
#plt.plot(insertion_avg['n'], insertion_avg['s'], marker='o', label='Insertion Sort')

plt.xlabel("n ")
plt.ylabel("avg c")
plt.legend()
plt.grid(True)
plt.show()

