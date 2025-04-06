use std::io::{self, BufRead};

fn merge(arr1: &mut [usize], arr2: &mut [usize], s: &mut usize, c: &mut usize) -> Vec<usize> {
    let mut i = 0;
    let mut j = 0;
    let mut merged = Vec::new();

    while i < arr1.len() && j < arr2.len() {
        *c += 1;
        if arr1[i] < arr2[j] {
            merged.push(arr1[i]);
            *s += 1;
            i += 1;
        } else {
            merged.push(arr2[j]);
            *s += 1;
            j += 1;
        }
    }

    merged.extend_from_slice(&arr1[i..]);
    merged.extend_from_slice(&arr2[j..]);
    println!("State: {:?}", merged);
    return merged;

}

fn merge_sort(arr: &mut [usize], s: &mut usize, c: &mut usize) {
    if arr.len() <= 1 {
        return;
    }
    let mid = arr.len() / 2;
    let mut left = arr[0..mid].to_vec();
    let mut right = arr[mid..].to_vec();
    merge_sort(&mut left, s, c);
    merge_sort(&mut right, s, c);
    let merged = merge(&mut left, &mut right, s, c);
    arr.copy_from_slice(&merged);
}

fn is_sorted(arr: &mut [usize]) -> bool {
    for i in 1..arr.len() {
        if arr[i] < arr[i - 1] {
            return false;
        }
    }
    true
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
            merge_sort(&mut numbers, &mut c, &mut s);
            println!("Posortowana: {:?}", numbers);
            println!("Nieposortowana: {:?}", copy);
            println!("s={}", s);
            println!("c={}", c);
            if is_sorted(&mut numbers) {
                println!("Tablica jest posortowana");
            } else {
                println!("Nieposortowane");
            }
        }
    }    
}
//// /home/wojteq18/Uni/AiSD/Lab/lista_2/generators/random_sequence/target/release/random_sequence 30 | /home/wojteq18/Uni/AiSD/Lab/lista_2/merge_sort/target/release/merge_sort

