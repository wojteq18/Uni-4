import matplotlib.pyplot as plt

def read_data(filename):
    lengths = []
    comparisons = []
    
    with open(filename, 'r') as f:
        lines = f.readlines()
    
    for line in lines:
        line = line.strip()
        line = line.replace(" ", "")
        
        if line.startswith('length='):
            val = int(line.split('=')[1])
            lengths.append(val)
        elif line.startswith('c='):
            val = int(line.split('=')[1])
            comparisons.append(val)
    
    return lengths, comparisons

def main():
    filename = 'outjj'  
    lengths, comps = read_data(filename)
    
    if len(lengths) != len(comps):
        print("Liczba elementów w 'lengths' i 'comparisons' nie jest taka sama!")
        print("lengths:", len(lengths), "comps:", len(comps))
        return
    
    plt.figure()
    plt.plot(lengths, comps, marker='o')  
    plt.xlabel('Threshold (length)')
    plt.ylabel('Comparisons (c)')
    plt.title('Comparisons vs. Threshold')
    plt.grid(True)
    plt.show()

if __name__ == '__main__':
    main()
