use std::io::{self, BufRead};


fn quick_sort(_length: usize, arr: &mut [usize], comparsion: &mut usize, swaps: &mut usize) {
    if arr.len() <= 1 {
        return;
    }

    let pivot_index = partition_lomuto(arr, comparsion, swaps);
    quick_sort(pivot_index, &mut arr[..pivot_index], comparsion, swaps);
    quick_sort(arr.len() - pivot_index - 1, &mut arr[pivot_index + 1..], comparsion, swaps);
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

fn main() {
    let mut c: usize = 0;
    let mut s: usize = 0;
    let stdin = io::stdin();
    let line = stdin.lock().lines().next().unwrap().unwrap();

    let length = line.split_whitespace()
        .next()
        .and_then(|s| s.parse::<usize>().ok())
        .unwrap_or(0);

    if let Some(start) = line.find('[') {
        if let Some(end) = line.find(']') {
            let numbers_str = &line[start + 1..end];

            let mut numbers: Vec<usize> = numbers_str
                .split(',')
                .filter_map(|s| s.trim().parse::<usize>().ok())
                .collect();

            quick_sort(length, &mut numbers, &mut c, &mut s);
            println!("{:?}", numbers);
            println!("s={}", s);
            println!("c={}", c);
        }
    }    
}
//// /home/wojteq18/Uni/AiSD/Lab/lista_2/generators/random_sequence/target/release/random_sequence 7 | /home/wojteq18/Uni/AiSD/Lab/lista_2/quick_sort/target/release/quick_sort
