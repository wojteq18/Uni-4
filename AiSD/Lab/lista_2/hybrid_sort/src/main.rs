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
    for i in 1..length {
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
    let mut arr: [usize; 10] = [2, 2, 3222, 43, 51, 6, 17, 18, 9, 0];

    hybrid_sort(10, &mut arr);

    println!("{:?}", arr);
}