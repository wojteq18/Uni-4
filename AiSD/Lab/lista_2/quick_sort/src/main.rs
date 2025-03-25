use std::env;

fn quick_sort (length: usize, arr: &mut [usize]) {
    if arr.len() <= 1 {
        return;
    }

    let pivot_index = partition_lomuto(arr);
    
    quick_sort(arr[..pivot_index].len(), &mut arr[..pivot_index]);
    quick_sort(arr[pivot_index + 1..].len(), &mut arr[pivot_index + 1..]);
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


fn main() {
    let mut arr: [usize; 10] = [2, 2, 3222, 43, 51, 6, 17, 18, 9, 0];

    quick_sort(10, &mut arr);

    println!("{:?}", arr);
}
