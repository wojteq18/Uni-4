use std::io::{self, BufRead};


fn sort_fives(arr: &mut [usize]) {
    arr.sort();
}

fn find_median_in_five(arr: [usize; 5]) -> usize {
    let mut temp = arr;
    temp.sort();
    return temp[2];
}

fn median_of_medians(arr: &mut [usize], comparsion: &mut usize, swaps: &mut usize) -> usize {
    let len = arr.len();
    if len < 5 {
        sort_fives(arr);
        return arr[len / 2];
    }

    let mut medians = Vec::new();
    for i in (0..len).step_by(5) {
        let end = std::cmp::min(i + 5, len);
        let mut sub_arr = [0; 5];
        sub_arr[..end - i].copy_from_slice(&arr[i..end]);
        medians.push(find_median_in_five(sub_arr));
    }

    median_of_medians(&mut medians, comparsion, swaps)
}

fn quick_sort(arr: &mut [usize], comparsion: &mut usize, swaps: &mut usize) {
    if arr.len() <= 1 {
        return;
    }

    let pivot_index = partition_lomuto(arr, comparsion, swaps);
    quick_sort( &mut arr[..pivot_index], comparsion, swaps);
    quick_sort( &mut arr[pivot_index + 1..], comparsion, swaps);
}

fn partition_lomuto(arr: &mut [usize], comparsion: &mut usize, swaps: &mut usize) -> usize {
    println!("State: {:?}", arr);
    let pivot = median_of_medians(arr, comparsion, swaps);
    arr[arr.len() - 1] = pivot;
    let mut i = 0;

    for j in 0..arr.len() - 1 {
        *comparsion += 1;
        if arr[j] <= pivot {
            if i != j {
                arr.swap(i, j);
                *swaps += 1;
            }
            i += 1;
        }
    }

    if i != arr.len() - 1 {
        arr.swap(i, arr.len() - 1);
        *swaps += 1;
    }

    i
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
            println!("Nieposortowana: {:?}", numbers);
            quick_sort( &mut numbers, &mut c, &mut s);
            println!("Posortowana: {:?}", numbers);
            println!("Nieposortowana: {:?}", copy);
            println!("{:?}", numbers);
            println!("s={}", s);
            println!("c={}", c);
            if is_sorted(&mut numbers) {
                println!("Posortowane");
            } else {
                println!("Nieposortowane");
            }
        }
    }    
}
//// /home/wojteq18/Uni/AiSD/Lab/lista_2/generators/random_sequence/target/release/random_sequence 7 | /home/wojteq18/Uni/AiSD/Lab/lista_2/quick_sort/target/release/quick_sort
