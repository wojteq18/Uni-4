use std::env;

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


fn main() {
    let args: Vec<String> = env::args().skip(1).collect(); //pomiń nazwę programu

    if args.len() < 2 {
        println!("Za mało argumentów! ");
        return;
    }

    let length = args[0].parse::<usize>().unwrap();
    if args.len() != length + 1 {
        println!("Zła długość tablicy, oczekuje {} elementów, podano {}", length, args.len() - 1);
        return;
    }

    let mut arr: Vec<usize> = args[1..].iter().map(|s| s.parse().expect("Invalid number")).collect();

    insertion_sort(length, &mut arr);
    println!("{:?}", arr);
}
