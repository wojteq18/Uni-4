import re
import matplotlib.pyplot as plt
import os

def parse_one_experiment(all_lines, start_idx, n):
    """
    Parses data for one experiment (n inserts into H1, n inserts into H2, 
    1 union, 2n extract-min operations) from the given lines starting at start_idx.
    Returns a list of comparison counts and the index of the next line after this experiment.
    Returns (None, -1) if parsing fails.
    """
    comparisons = []
    current_idx = start_idx
    expected_total_ops = 4 * n + 1

    # Phase 1: Insert H1 (n operations)
    for i in range(n):
        if current_idx >= len(all_lines):
            # print(f"Debug (H1 Insert): Expected line {current_idx+1}, but EOF reached. Op {i+1}/{n}.")
            return None, -1
        line = all_lines[current_idx].strip()
        match = re.match(r"Inserted in h_1: \d+, comparisons: (\d+)", line)
        if not match:
            # print(f"Debug (H1 Insert): Line {current_idx+1} format error: '{line}'")
            return None, -1
        comparisons.append(int(match.group(1)))
        current_idx += 1

    # Phase 2: Insert H2 (n operations)
    for i in range(n):
        if current_idx >= len(all_lines):
            # print(f"Debug (H2 Insert): Expected line {current_idx+1}, but EOF reached. Op {i+1}/{n}.")
            return None, -1
        line = all_lines[current_idx].strip()
        match = re.match(r"Inserted in h_2: \d+, comparisons: (\d+)", line)
        if not match:
            # print(f"Debug (H2 Insert): Line {current_idx+1} format error: '{line}'")
            return None, -1
        comparisons.append(int(match.group(1)))
        current_idx += 1

    # Phase 3: Union (1 operation)
    if current_idx >= len(all_lines):
        # print(f"Debug (Union): Expected line {current_idx+1}, but EOF reached.")
        return None, -1
    line = all_lines[current_idx].strip()
    match = re.match(r"Union of h_1 and h_2, comparisons: (\d+)", line)
    if not match:
        # print(f"Debug (Union): Line {current_idx+1} format error: '{line}'")
        return None, -1
    comparisons.append(int(match.group(1)))
    current_idx += 1

    # Phase 4: Extract-Min (2n operations)
    for i in range(2 * n):
        if current_idx + 1 >= len(all_lines): # Need two lines for each extract op
            # print(f"Debug (Extract-Min): Expected lines {current_idx+1}-{current_idx+2}, but EOF reached. Op {i+1}/{2*n}.")
            return None, -1
        
        line_extract_val = all_lines[current_idx].strip()
        line_extract_cmp = all_lines[current_idx + 1].strip()

        match_val = re.match(r"Extracted min: .*", line_extract_val)
        match_cmp = re.match(r"Comparisons during extraction: (\d+)", line_extract_cmp)

        if not match_val or not match_cmp:
            # print(f"Debug (Extract-Min): Lines {current_idx+1}-{current_idx+2} format error: '{line_extract_val}' / '{line_extract_cmp}'")
            return None, -1
        
        comparisons.append(int(match_cmp.group(1)))
        current_idx += 2
    
    if len(comparisons) == expected_total_ops:
        return comparisons, current_idx
    else:
        # print(f"Debug: Collected {len(comparisons)} ops, expected {expected_total_ops}.")
        return None, -1


def plot_experiment_comparisons(experiment_data, experiment_num, n_val, output_dir_path):
    """
    Generates and saves a plot for a single experiment's comparison data.
    """
    plt.figure(figsize=(16, 8)) 
    operation_indices = list(range(1, len(experiment_data) + 1))
    plt.plot(operation_indices, experiment_data, marker='.', linestyle='-', markersize=2.5, linewidth=0.7, alpha=0.8)

    plt.title(f'Eksperyment {experiment_num} (n={n_val}): Liczba porównań na operację', fontsize=16)
    plt.xlabel('Indeks operacji (i-ta operacja)', fontsize=12)
    plt.ylabel('Liczba porównań (c)', fontsize=12)
    plt.grid(True, linestyle='--', alpha=0.6)
    
    ymin, ymax = plt.ylim() # Get y-limits after data is plotted for text placement

    # Add vertical lines and text to distinguish phases
    insert_h1_end = n_val
    plt.axvline(x=insert_h1_end, color='orangered', linestyle='--', linewidth=1.2)
    plt.text(insert_h1_end / 2, ymax * 0.95, 'Insert H1', horizontalalignment='center', color='orangered', fontsize=10)

    insert_h2_end = insert_h1_end + n_val
    plt.axvline(x=insert_h2_end, color='forestgreen', linestyle='--', linewidth=1.2)
    plt.text(insert_h1_end + (n_val / 2), ymax * 0.95, 'Insert H2', horizontalalignment='center', color='forestgreen', fontsize=10)

    union_end = insert_h2_end + 1
    plt.axvline(x=union_end, color='royalblue', linestyle='--', linewidth=1.2)
    plt.text(insert_h2_end + 0.5, ymax * 0.88, 'Union', horizontalalignment='center', color='royalblue', fontsize=10, rotation=0) # Adjusted y for visibility

    extract_min_start = union_end
    plt.text(extract_min_start + n_val, ymax * 0.95, 'Extract-Min', horizontalalignment='center', color='darkviolet', fontsize=10)
    
    plt.xlim(0, len(experiment_data) + 1)
    plt.tight_layout() # Adjust plot to ensure everything fits

    plot_filename_leaf = f'experiment_{experiment_num}_n{n_val}_comparisons_history.png'
    full_plot_path = os.path.join(output_dir_path, plot_filename_leaf)
    plt.savefig(full_plot_path)
    plt.close() 
    print(f"Wykres zapisano do: {full_plot_path}")

# Main script execution
if __name__ == "__main__":
    n_parameter = 500
    num_experiments_to_plot = 5
    
    # Paths are relative to the workspace root if the script is there,
    # or relative to the script's location if run directly.
    # Assuming 'zadanie3' is a subdirectory in the current working directory or workspace root.
    data_file_path = "out1.txt"
    plots_output_dir = "." # Save plots in the zadanie3 directory

    # Ensure the output directory for plots exists
    os.makedirs(plots_output_dir, exist_ok=True)

    try:
        with open(data_file_path, 'r') as f:
            all_lines_from_file = f.readlines()
    except FileNotFoundError:
        print(f"Błąd: Plik danych nie został znaleziony pod ścieżką: {data_file_path}")
        all_lines_from_file = [] # Avoid crashing if file not found

    if not all_lines_from_file:
        print("Brak danych do przetworzenia. Upewnij się, że plik out1.txt istnieje i zawiera dane.")
    else:
        current_line_in_file = 0
        successfully_plotted_count = 0
        
        for i in range(num_experiments_to_plot):
            print(f"\nPróba przetworzenia eksperymentu {i + 1} dla n={n_parameter}, zaczynając od linii ~{current_line_in_file + 1} pliku danych.")
            
            if current_line_in_file >= len(all_lines_from_file):
                print("Osiągnięto koniec pliku danych. Nie można przetworzyć więcej eksperymentów.")
                break

            # Attempt to parse one full experiment
            experiment_comparison_data, next_start_line = parse_one_experiment(all_lines_from_file, current_line_in_file, n_parameter)
            
            if experiment_comparison_data:
                lines_consumed = next_start_line - current_line_in_file
                print(f"Pomyślnie przetworzono eksperyment {i + 1} (n={n_parameter}). Liczba operacji: {len(experiment_comparison_data)}. Przetworzono linie: {lines_consumed}.")
                plot_experiment_comparisons(experiment_comparison_data, i + 1, n_parameter, plots_output_dir)
                successfully_plotted_count += 1
                current_line_in_file = next_start_line # Move cursor for the next experiment
            else:
                print(f"Nie udało się przetworzyć pełnego eksperymentu {i + 1} (n={n_parameter}) zaczynając od linii {current_line_in_file + 1}.")
                print("Może to być spowodowane nieprawidłowym formatem danych, niewystarczającą ilością danych lub końcem istotnych sekcji danych.")
                
                # Try to find the start of a potentially new experiment to resynchronize
                initial_search_idx = current_line_in_file + 1 # Start searching from the line after the current attempt
                found_next_start_candidate = False
                for search_idx in range(initial_search_idx, len(all_lines_from_file)):
                    # A common pattern for the start of an experiment's data
                    if re.match(r"Inserted in h_1: 1, comparisons: 0", all_lines_from_file[search_idx].strip()):
                        current_line_in_file = search_idx
                        print(f"Znaleziono potencjalny początek nowego eksperymentu na linii {current_line_in_file + 1}. Spróbuję w następnej iteracji.")
                        found_next_start_candidate = True
                        break 
                
                if not found_next_start_candidate:
                    print("Nie można znaleźć kolejnego potencjalnego początku eksperymentu. Zatrzymywanie.")
                    break # Stop if no new potential start is found after a failed parse

        # Final summary messages
        if successfully_plotted_count == 0:
            print(f"\nNie udało się wygenerować żadnych wykresów dla n={n_parameter} z pliku '{data_file_path}'.")
        elif successfully_plotted_count < num_experiments_to_plot:
            print(f"\nWygenerowano {successfully_plotted_count} z {num_experiments_to_plot} żądanych wykresów dla n={n_parameter}.")
        else:
            print(f"\nPomyślnie wygenerowano wszystkie {num_experiments_to_plot} żądane wykresy dla n={n_parameter}.")