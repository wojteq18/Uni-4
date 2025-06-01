import pandas as pd
import matplotlib.pyplot as plt
from io import StringIO

# Trzeci zestaw danych
data = """
Test,sign,n,avg_cmp,max_cmp,avg_point,max_point,avg_height,max_height
r;10000;19;10700;133;58848;2520;10000
r;20000;21;37706;146;207377;5021;20000
r;30000;22;58916;154;324036;7522;30000
r;40000;23;59070;159;324879;10025;40000
r;50000;23;73828;163;406052;12524;50000
r;60000;24;69476;166;382116;15024;60000
r;70000;24;112556;169;619056;17525;70000
r;80000;24;151180;171;831488;20025;80000
r;90000;25;137480;174;756138;22526;90000
r;100000;25;138278;176;760523;25029;100000
l;10000;25;84;211;650;30;51
l;20000;28;88;232;622;34;61
l;30000;29;90;244;684;35;59
l;40000;30;94;252;711;37;60
l;50000;31;100;259;717;38;61
l;60000;31;102;264;743;39;67
l;70000;32;106;269;786;39;61
l;80000;32;108;273;841;40;67
l;90000;33;114;276;819;41;66
l;100000;33;106;280;821;41;67
"""

# Wczytywanie i konwersja
data_no_header = "\n".join(data.strip().split("\n")[1:])
df = pd.read_csv(StringIO(data_no_header), sep=';', header=None)
df.columns = ['Test', 'n', 'avg_cmp', 'max_cmp', 'avg_point', 'max_point', 'avg_height', 'max_height']
df[['n', 'avg_cmp', 'max_cmp', 'avg_point', 'max_point', 'avg_height', 'max_height']] = df[[
    'n', 'avg_cmp', 'max_cmp', 'avg_point', 'max_point', 'avg_height', 'max_height'
]].apply(pd.to_numeric)

# Oddzielenie danych
df_l = df[df['Test'] == 'l'].sort_values(by='n')
df_r = df[df['Test'] == 'r'].sort_values(by='n')

# Funkcja rysująca oddzielne wykresy
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

# Wykresy dla losowego
plot_separate(df_l, "Liczba porównań", "avg_cmp", "max_cmp", "Porównania kluczy", "Losowe")
plot_separate(df_l, "Odczyty/podstawienia wskaźników", "avg_point", "max_point", "Wskaźniki", "Losowe")
plot_separate(df_l, "Wysokość drzewa", "avg_height", "max_height", "Wysokość drzewa", "Losowe")

# Wykresy dla rosnącego
plot_separate(df_r, "Liczba porównań", "avg_cmp", "max_cmp", "Porównania kluczy", "Rosnące")
plot_separate(df_r, "Odczyty/podstawienia wskaźników", "avg_point", "max_point", "Wskaźniki", "Rosnące")
plot_separate(df_r, "Wysokość drzewa", "avg_height", "max_height", "Wysokość drzewa", "Rosnące")
