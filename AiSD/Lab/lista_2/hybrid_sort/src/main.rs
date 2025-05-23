use std::io::{self, BufRead};


fn quick_sort(arr: &mut [usize], comparsion: &mut usize, swaps: &mut usize) {
    if arr.len() <= 1 {
        return;
    }

    let pivot_index = partition_lomuto(arr, comparsion, swaps);
    hybrid_sort(&mut arr[..pivot_index], comparsion, swaps);
    hybrid_sort(&mut arr[pivot_index + 1..], comparsion, swaps);
}

fn partition_lomuto(arr: &mut [usize], comparsion: &mut usize, swaps: &mut usize) -> usize {
    let pivot = arr[arr.len() - 1];
    let mut i = 0;

    for j in 0..arr.len() - 1 {
        *comparsion += 1;
        if arr[j] <= pivot {
            arr.swap(i, j);
            i += 1;
            *swaps += 1;
        }
    }
    arr.swap(i, arr.len() - 1);
    return i
}


fn insertion_sort(arr: &mut [usize], comparsion: &mut usize, swaps: &mut usize) {
    for i in 1..arr.len() {
        let key = arr[i];
        let mut j = i;

        // Szukamy miejsca dla key (to jest jedyne miejsce gdzie zliczamy porównania)
        while j > 0 {
            *comparsion += 1;
            if arr[j - 1] > key {
                j -= 1;
            } else {
                break;
            }
        }

        // Przesuwamy blok [j..i) o 1 w prawo
        if j != i {
            arr.copy_within(j..i, j + 1);
            arr[j] = key;
            *swaps += 1; // swap to u nas oznacza "przesunięcie bloku"
        }
    }
}

fn is_sorted(arr: &mut [usize]) -> bool {
    for i in 1..arr.len() {
        if arr[i] < arr[i - 1] {
            return false;
        }
    }
    true
}



fn hybrid_sort(arr: &mut [usize], comparsion: &mut usize, swaps: &mut usize) {
    if arr.len() < 20 {
        insertion_sort(arr, comparsion, swaps);
    } else {
        quick_sort(arr, comparsion, swaps);
    }
    println!("State: {:?}", arr);
}

fn main() {
    let mut c: usize = 0;
    let mut s: usize = 0;
    let stdin = io::stdin();
    let line = stdin.lock().lines().next().unwrap().unwrap();

    if let Some(start) = line.find('[') {
        if let Some(end) = line.find(']') {
            let numbers_str = &line[start + 1..end];

            let mut numbers: Vec<usize> = numbers_str
                .split(',')
                .filter_map(|s| s.trim().parse::<usize>().ok())
                .collect();

            let copy = numbers.clone();
            println!("Nieposortowana: {:?}", numbers);
            hybrid_sort(&mut numbers, &mut c, &mut s);
            println!("Posortowana: {:?}", numbers);
            println!("Nieposortowana: {:?}", copy);

            if is_sorted(&mut numbers) {
                println!("Tablica jest posortowana");
            } else {
                println!("Array is not sorted");
            }
            println!("s={}", s);
            println!("c={}", c);
        }
    } 
}
// /home/wojteq18/Uni/AiSD/Lab/lista_2/generators/random_sequence/target/release/random_sequence 7 | /home/wojteq18/Uni/AiSD/Lab/lista_2/hybrid_sort/target/release/hybrid_sort