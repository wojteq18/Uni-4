import matplotlib.pyplot as plt
import re
import os

def parse_n_range_data(file_path):
    """
    Parses data from the file containing n vs. average comparisons.
    Expected format: "n=VALUE, comparisons: AVG_COMPARISONS"
    Returns two lists: n_values and avg_comparisons_values.
    """
    n_values = []
    avg_comparisons_values = []
    
    try:
        with open(file_path, 'r') as f:
            for line in f:
                line = line.strip()
                match = re.match(r"n=(\d+), comparisons: ([\d\.]+)", line)
                if match:
                    n_values.append(int(match.group(1)))
                    avg_comparisons_values.append(float(match.group(2)))
                else:
                    print(f"Ostrzeżenie: Pomijanie linii o nieoczekiwanym formacie: '{line}'")
    except FileNotFoundError:
        print(f"Błąd: Plik danych nie został znaleziony pod ścieżką: {file_path}")
        return None, None
        
    return n_values, avg_comparisons_values

def plot_n_range_comparisons(n_values, avg_comparisons_values, output_dir_path):
    """
    Generates and saves a plot for n vs. average comparisons.
    """
    if not n_values or not avg_comparisons_values:
        print("Brak danych do wygenerowania wykresu.")
        return

    plt.figure(figsize=(12, 7))
    plt.plot(n_values, avg_comparisons_values, marker='o', linestyle='-', markersize=4, color='dodgerblue')
    
    plt.title('Zależność średniego kosztu operacji od n', fontsize=16)
    plt.xlabel('Rozmiar danych wejściowych (n)', fontsize=12)
    plt.ylabel('Średni koszt operacji (łączna liczba porównań / n)', fontsize=12)
    plt.grid(True, linestyle='--', alpha=0.7)
    plt.tight_layout()
    
    plot_filename = "average_cost_vs_n_plot.png"
    full_plot_path = os.path.join(output_dir_path, plot_filename)
    
    try:
        plt.savefig(full_plot_path)
        print(f"Wykres zapisano do: {full_plot_path}")
    except Exception as e:
        print(f"Nie udało się zapisać wykresu: {e}")
    
    plt.close()

# Main script execution
if __name__ == "__main__":
    # Assuming out2.txt is in the same directory as the script.
    # Plots will also be saved in this directory.
    data_file_path = "out2.txt" 
    plots_output_dir = "." # Save plots in the current directory

    # Ensure the output directory for plots exists (it's the current dir, so it does)
    # os.makedirs(plots_output_dir, exist_ok=True) # Not strictly necessary for "."

    n_vals, avg_comps = parse_n_range_data(data_file_path)
    
    if n_vals and avg_comps:
        plot_n_range_comparisons(n_vals, avg_comps, plots_output_dir)
    else:
        print("Nie udało się przetworzyć danych lub plik nie istnieje.")
