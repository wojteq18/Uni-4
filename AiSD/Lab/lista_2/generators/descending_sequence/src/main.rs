use rand_mt::Mt64; 
use rand::Rng;
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

    if args.len() != 1 {
        println!("Zła ilość argumentów! ");
        return;
    }

    let length = args[0].parse::<usize>().unwrap();

    let seed = rand::thread_rng().gen(); 
    let mut rng = Mt64::new(seed);

    let limit = 2 * length - 1;

    let mut arr: Vec<usize> = (0..length).map(|_| rng.gen_range(0..limit)).collect();
    insertion_sort(length, &mut arr);
    arr.reverse();
    println!("{:?}", arr);
}