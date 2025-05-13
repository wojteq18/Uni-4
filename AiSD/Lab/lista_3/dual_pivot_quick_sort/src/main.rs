use std::io::{self, BufRead};

fn sort_fives(arr: &mut [usize]) {
    arr.sort();
}

fn find_median_in_five(arr: [usize; 5]) -> usize {
    let mut temp = arr;
    temp.sort();
    temp[2]
}

fn median_of_medians(arr: &[usize]) -> usize {
    let len = arr.len();
    if len < 5 {
        let mut temp = arr.to_vec();
        temp.sort();
        return temp[len / 2];
    }

    let mut medians = Vec::new();
    for i in (0..len).step_by(5) {
        let end = std::cmp::min(i + 5, len);
        let mut sub_arr = arr[i..end].to_vec();
        sub_arr.sort();
        medians.push(sub_arr[sub_arr.len() / 2]);
    }

    median_of_medians(&medians)
}

fn dp_quick_sort(arr: &mut [usize], comparsion: &mut usize, swaps: &mut usize) {
    if arr.len() <= 1 {
        return;
    }

    let (first_pivot, last_pivot) = pivot_index(arr, comparsion, swaps);

    dp_quick_sort(&mut arr[..first_pivot], comparsion, swaps);
    dp_quick_sort(&mut arr[first_pivot + 1..last_pivot], comparsion, swaps);
    dp_quick_sort(&mut arr[last_pivot + 1..], comparsion, swaps);
}

fn pivot_index(arr: &mut [usize], comparsion: &mut usize, swaps: &mut usize) -> (usize, usize) {
    println!("State: {:?}", arr);
    let len = arr.len();

    let p1 = median_of_medians(&arr[..len / 2]);
    let p2 = median_of_medians(&arr[len / 2..]);

    let p1_index = arr.iter().position(|&x| x == p1).unwrap();
    arr.swap(0, p1_index);
    *swaps += 1;

    let p2_index = if p1_index == len - 1 {
        arr.iter().rposition(|&x| x == p2).unwrap()
    } else {
        arr.iter().rposition(|&x| x == p2).unwrap()
    };
    arr.swap(len - 1, p2_index);
    *swaps += 1;

    *comparsion += 1;
    if arr[0] > arr[len - 1] {
        arr.swap(0, len - 1);
        *swaps += 1;
    }

    let p1 = arr[0];
    let p2 = arr[len - 1];

    let mut less = 1;
    let mut great = len - 2;
    let mut k = 1;

    while k <= great {
        *comparsion += 1;
        if arr[k] < p1 {
            arr.swap(k, less);
            *swaps += 1;
            less += 1;
            k += 1;
        } else {
            *comparsion += 1;
            if arr[k] > p2 {
                while arr[great] > p2 && k < great {
                    *comparsion += 1;
                    great -= 1;
                }
                arr.swap(k, great);
                *swaps += 1;
                great -= 1;

                *comparsion += 1;
                if arr[k] < p1 {
                    arr.swap(k, less);
                    *swaps += 1;
                    less += 1;
                }
                k += 1;
            } else {
                k += 1;
            }
        }
    }

    less -= 1;
    great += 1;

    arr.swap(0, less);
    arr.swap(len - 1, great);
    *swaps += 2;

    (less, great)
}

fn is_sorted(arr: &[usize]) -> bool {
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
            dp_quick_sort(&mut numbers, &mut c, &mut s);
            println!("Posortowana: {:?}", numbers);
            println!("Nieposortowana: {:?}", copy);
            println!("s={}", s);
            println!("c={}", c);
            if is_sorted(&numbers) {
                println!("Posortowane");
            } else {
                println!("Nieposortowane");
            }
        }
    }
}
