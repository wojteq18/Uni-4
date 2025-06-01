import pandas as pd
import matplotlib.pyplot as plt
from io import StringIO

# Dane
data = """
Test,sign,n,avg_cmp,max_cmp,avg_point,max_point,avg_height,max_height
l;10000;14;29;30;59;25;30
l;20000;16;33;33;67;27;34
l;30000;17;39;36;79;31;40
l;40000;17;36;37;73;31;37
l;50000;19;41;39;83;34;42
l;60000;19;38;39;77;33;39
l;70000;19;38;39;77;34;39
l;80000;19;41;40;83;34;42
l;90000;19;44;39;89;35;45
l;100000;19;40;40;81;34;41
r;10000;4999;9999;9999;19999;5000;10000
r;20000;9999;19999;19999;39999;10000;20000
r;30000;14999;29999;29999;59999;15000;30000
r;40000;19999;39999;39999;79999;20000;40000
r;50000;24999;49999;49999;99999;25000;50000
r;60000;29999;59999;59999;119999;30000;60000
r;70000;34999;69999;69999;139999;35000;70000
r;80000;39999;79999;79999;159999;40000;80000
r;90000;44999;89999;89999;179999;45000;90000
r;100000;49999;99999;99999;199999;50000;100000
"""

# Parsowanie danych
data_no_header = "\n".join(data.strip().split("\n")[1:])
df = pd.read_csv(StringIO(data_no_header), sep=';', header=None)
df.columns = ['Test', 'n', 'avg_cmp', 'max_cmp', 'avg_point', 'max_point', 'avg_height', 'max_height']
df[['n', 'avg_cmp', 'max_cmp', 'avg_point', 'max_point', 'avg_height', 'max_height']] = df[[
    'n', 'avg_cmp', 'max_cmp', 'avg_point', 'max_point', 'avg_height', 'max_height'
]].apply(pd.to_numeric)

# Rozdzielenie danych i sortowanie
df_l = df[df['Test'] == 'l'].sort_values(by='n')
df_r = df[df['Test'] == 'r'].sort_values(by='n')

# Funkcja do rysowania oddzielnych wykresów
def plot_separate(df, y_label, column_avg, column_max, title, tag):
    plt.figure()
    plt.plot(df['n'], df[column_avg], 'o-', label=f'{tag} - Średnia')
    plt.plot(df['n'], df[column_max], 'o--', label=f'{tag} - Maksimum')
    plt.xlabel("n (rozmiar danych)")
    plt.ylabel(y_label)
    plt.title(title + f" ({tag})")
    plt.legend()
    plt.grid(True)
    plt.tight_layout()
    plt.show()

# Wykresy dla danych losowych
plot_separate(df_l, "Liczba porównań", "avg_cmp", "max_cmp", "Porównania kluczy", "Losowe")
plot_separate(df_l, "Odczyty/podstawienia wskaźników", "avg_point", "max_point", "Wskaźniki", "Losowe")
plot_separate(df_l, "Wysokość drzewa", "avg_height", "max_height", "Wysokość drzewa", "Losowe")

# Wykresy dla danych rosnących
plot_separate(df_r, "Liczba porównań", "avg_cmp", "max_cmp", "Porównania kluczy", "Rosnące")
plot_separate(df_r, "Odczyty/podstawienia wskaźników", "avg_point", "max_point", "Wskaźniki", "Rosnące")
plot_separate(df_r, "Wysokość drzewa", "avg_height", "max_height", "Wysokość drzewa", "Rosnące")
