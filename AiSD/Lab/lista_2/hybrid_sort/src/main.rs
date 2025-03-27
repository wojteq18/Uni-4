use std::io::{self, BufRead};


fn quick_sort (length: usize, arr: &mut [usize]) {
    if arr.len() <= 1 {
        return;
    }

    let pivot_index = partition_lomuto(arr);
    
    hybrid_sort(arr[..pivot_index].len(), &mut arr[..pivot_index]);
    hybrid_sort(arr[pivot_index + 1..].len(), &mut arr[pivot_index + 1..]);
}

fn partition_lomuto(arr: &mut [usize]) -> usize {
    let pivot = arr[arr.len() - 1];
    let mut i = 0;

    for j in 0..arr.len() - 1 {
        if arr[j] <= pivot {
            arr.swap(i, j);
            i = i + 1;
        }
    }
    arr.swap(i, arr.len() - 1);
    return i;
}

fn insertion_sort (length: usize, arr: &mut [usize]) {
    for i in 1..arr.len() {
        let key = arr[i];
        let mut j: usize = i;
        while j > 0 && arr[j - 1] > key {
            arr[j] = arr[j - 1];
            j = j - 1;
        }
        arr[j] = key;
    }
}

fn hybrid_sort(length: usize, arr: &mut [usize]) {
    if length < 16 {
        insertion_sort(length, arr);
    } else {
        quick_sort(length, arr);
    }
}

fn main() {
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

            quick_sort(length, &mut numbers);
            println!("{:?}", numbers);
        }
    }    
}
// /home/wojteq18/Uni/AiSD/Lab/lista_2/generators/random_sequence/target/release/random_sequence 7 | /home/wojteq18/Uni/AiSD/Lab/lista_2/hybrid_sort/target/release/hybrid_sort
