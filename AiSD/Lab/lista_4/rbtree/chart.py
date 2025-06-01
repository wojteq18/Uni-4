import pandas as pd
import matplotlib.pyplot as plt
from io import StringIO

# Nowy zestaw danych
data = """
Test,sign,n,avg_cmp,max_cmp,avg_point,max_point,avg_height,max_height
r;10000;21;29;50;95;16;24
r;20000;23;35;53;102;17;26
r;30000;24;33;55;107;18;27
r;40000;25;35;56;109;19;28
r;50000;26;37;57;111;19;29
r;60000;26;37;58;114;19;29
r;70000;27;39;59;114;20;30
r;80000;27;45;59;116;20;30
r;90000;27;39;60;118;20;30
r;100000;28;41;60;118;21;31
l;10000;17;31;38;80;14;16
l;20000;18;33;40;85;15;18
l;30000;19;35;41;83;16;18
l;40000;20;37;42;86;17;19
l;50000;20;37;43;89;17;20
l;60000;21;39;43;85;17;20
l;70000;21;39;43;93;18;20
l;80000;21;39;44;89;18;20
l;90000;22;39;44;91;18;20
l;100000;22;39;45;90;18;21
"""

# Parsowanie danych
data_no_header = "\n".join(data.strip().split("\n")[1:])
df = pd.read_csv(StringIO(data_no_header), sep=';', header=None)
df.columns = ['Test', 'n', 'avg_cmp', 'max_cmp', 'avg_point', 'max_point', 'avg_height', 'max_height']
df[['n', 'avg_cmp', 'max_cmp', 'avg_point', 'max_point', 'avg_height', 'max_height']] = df[[
    'n', 'avg_cmp', 'max_cmp', 'avg_point', 'max_point', 'avg_height', 'max_height'
]].apply(pd.to_numeric)

# Rozdzielenie i sortowanie
df_l = df[df['Test'] == 'l'].sort_values(by='n')
df_r = df[df['Test'] == 'r'].sort_values(by='n')

# Funkcja do rysowania
def plot_separate(df, y_label, column_avg, column_max, title, tag):
    plt.figure()
    plt.plot(df['n'], df[column_avg], 'o-', label=f'{tag} - Średnia')
    plt.plot(df['n'], df[column_max], 'o--', label=f'{tag} - Maksimum')
    plt.xlabel("n (rozmiar danych)")
    plt.ylabel(y_label)
    plt.title(f"{title} ({tag})")
    plt.legend()
    plt.grid(True)
    plt.tight_layout()
    plt.show()

# Wykresy – losowe
plot_separate(df_l, "Liczba porównań", "avg_cmp", "max_cmp", "Porównania kluczy", "Losowe")
plot_separate(df_l, "Odczyty/podstawienia wskaźników", "avg_point", "max_point", "Wskaźniki", "Losowe")
plot_separate(df_l, "Wysokość drzewa", "avg_height", "max_height", "Wysokość drzewa", "Losowe")

# Wykresy – rosnące
plot_separate(df_r, "Liczba porównań", "avg_cmp", "max_cmp", "Porównania kluczy", "Rosnące")
plot_separate(df_r, "Odczyty/podstawienia wskaźników", "avg_point", "max_point", "Wskaźniki", "Rosnące")
plot_separate(df_r, "Wysokość drzewa", "avg_height", "max_height", "Wysokość drzewa", "Rosnące")
