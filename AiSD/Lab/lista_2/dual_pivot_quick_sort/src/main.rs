use std::io::{self, BufRead};


fn dp_quick_sort(_length: usize, arr: &mut [usize], comparsion: &mut usize, swaps: &mut usize) {
    if arr.len() <= 1 {
        return;
    }

    let (first_pivot, last_pivot) = pivot_index(arr, comparsion, swaps);
    dp_quick_sort(first_pivot, &mut arr[..first_pivot], comparsion, swaps);
    dp_quick_sort(last_pivot - first_pivot - 1, &mut arr[first_pivot + 1..last_pivot], comparsion, swaps);
    dp_quick_sort(arr.len() - last_pivot - 1, &mut arr[last_pivot + 1..], comparsion, swaps);
}

fn pivot_index(arr: &mut [usize], comparsion: &mut usize, swaps: &mut usize) -> (usize, usize) {
    println!("State: {:?}", arr);
    let len = arr.len();

    if arr[0] > arr[len - 1] {
        arr.swap(0, len - 1);
        *swaps += 1;
    }

    let p1 = arr[0];
    let p2 = arr[len - 1];

    let mut less = 1; //wskaźnik na element mniejszy od p1
    let mut great = len - 2; //wskaźnik na element większy od p2
    let mut k = 1;

    while k <= great {
        *comparsion += 1;
        if arr[k] < p1 {
            arr.swap(k, less);
            less += 1;
            *swaps += 1;
        } else if arr[k] > p2 {
            while arr[great] > p2 && k < great {
                great -= 1;
                *comparsion += 1;
            }
            arr.swap(k, great);
            great -= 1;
            *swaps += 1;
            if arr[k] < p1 {
                arr.swap(k, less);
                less += 1;
                *swaps += 1;
            }
        }
        k += 1;
    }

    less -= 1;
    great += 1;

    arr.swap(0, less);
    arr.swap(len - 1, great);
    *swaps += 2;

    return (less, great);
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

            let copy = numbers.clone();
            println!("Nieposortowana: {:?}", numbers);
            dp_quick_sort(length, &mut numbers, &mut c, &mut s);
            println!("Posortowana: {:?}", numbers);
            println!("Nieposortowana: {:?}", copy);
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
//// /home/wojteq18/Uni/AiSD/Lab/lista_2/generators/random_sequence/target/release/random_sequence 30 | /home/wojteq18/Uni/AiSD/Lab/lista_2/dual_pivot_quick_sort/target/release/dual_pivot_quick_sort
