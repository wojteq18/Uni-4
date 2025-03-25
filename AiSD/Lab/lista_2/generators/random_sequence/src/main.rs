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

    let limit = 64;

    let arr: Vec<usize> = (0..length).map(|_| rng.gen_range(0..limit)).collect();
    println!("{:?}", arr);
}