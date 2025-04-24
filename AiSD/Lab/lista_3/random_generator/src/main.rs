use rand_mt::Mt64; 
use rand::Rng;
use std::env;

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

    let random_index = rng.gen_range(0..length-1);

    let arr: Vec<usize> = (0..length).map(|_| rng.gen_range(0..limit)).collect();
    println!("{:?} {}", arr, random_index);
}

//../random_generator/target/release/random_generator 44 | go run random_select.go 