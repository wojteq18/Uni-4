use std::io::{self, BufRead};

fn find_runs(arr: &mut [usize]) -> Vec<(usize)> {
    let mut runs = Vec::new();
    let length = arr.len();
    runs.push(0);
    let mut i = 0;
    while i < length - 1 {
        let mut start = i;
        if arr[i+1] >= arr[i] { 
            while i + 1 < length && arr[i+1] >= arr[i] {
                i += 1;
            }
        } else {
            while i + 1 < length && arr[i+1] <= arr[i] {
                i += 1;
            }
            arr[start..=i].reverse();
        }
        runs.push(i+1);
        i += 1;
    }

    if *runs.last().unwrap() != length {
        runs.push(length); 
    }
    return runs;
}

fn merge (arr1: &mut [usize], arr2: &mut [usize]) -> Vec<(usize)> {
    let mut merged = Vec::new();
    let mut i = 0;
    let mut j = 0;
    while i < arr1.len() && j < arr2.len() {
        if arr1[i] <= arr2[j] {
            merged.push(arr1[i]);
            i += 1;
        } else {
            merged.push(arr2[j]);
            j += 1;
        }
    }

    merged.extend_from_slice(&arr1[i..]); 
    merged.extend_from_slice(&arr2[j..]); 


    return merged;
}

fn wojteq_sort(arr: &mut [usize]) {
    loop {
        let runs = find_runs(arr);

        if runs.len() <= 2 {
            break;
        }

        let mut i = 0;
        while i + 2 < runs.len() {
            let left_start  = runs[i];
            let left_end    = runs[i + 1];
            let right_start = runs[i + 1];
            let right_end   = runs[i + 2];

            let (left_part, right_part) = arr.split_at_mut(right_start);

            let left_slice  = &mut left_part[left_start..left_end];
            let right_slice = &mut right_part[..(right_end - right_start)];

            let merged = merge(left_slice, right_slice);

            let left_len = left_slice.len();
            left_slice.copy_from_slice(&merged[..left_len]);
            right_slice.copy_from_slice(&merged[left_len..]);

            i += 2;
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


            wojteq_sort(&mut numbers);
            println!("{:?}", numbers);

            if is_sorted(&mut numbers) {
                println!("Sorted");
            } else {
                println!("Not sorted");
            }
        }
    }    
}
// /home/wojteq18/Uni/AiSD/Lab/lista_2/generators/random_sequence/target/release/random_sequence 7 | /home/wojteq18/Uni/AiSD/Lab/lista_2/wojteq_sort/target/release/wojteq_sort


