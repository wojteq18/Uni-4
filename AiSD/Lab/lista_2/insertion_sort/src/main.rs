use std::io::{self, BufRead};


fn insertion_sort (arr: &mut [usize]) {
    let mut s = 0;
    let mut c = 0;

    for i in 1..arr.len() {
        let key = arr[i];
        let mut j = i;

        while j > 0 {
            c += 1;
            if arr[j - 1] > key {
                j -= 1;
            } else {
                break;
            }
        }

        if j != i {
            arr.copy_within(j..i, j + 1);
            arr[j] = key;
            s += 1;
            println!("State: {:?}", arr);
        }
    }
    println!("s={}", s);
    println!("c={}", c);
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

    if let Some(start) = line.find('[') {
        if let Some(end) = line.find(']') {
            let numbers_str = &line[start + 1..end];

            let mut numbers: Vec<usize> = numbers_str
                .split(',')
                .filter_map(|s| s.trim().parse::<usize>().ok())
                .collect();

            let copy = numbers.clone();
            println!("Nieposortowana: {:?}", numbers);
            insertion_sort(&mut numbers);
            println!("Posortowana: {:?}", numbers);
            println!("Nieposortowana: {:?}", copy);

            if is_sorted(&mut numbers) {
                println!("Tablica jest posortowana");
            } else {
                println!("Not sorted");
            }
        }
    }    
}
// /home/wojteq18/Uni/AiSD/Lab/lista_2/generators/random_sequence/target/release/random_sequence 7 | /home/wojteq18/Uni/AiSD/Lab/lista_2/insertion_sort/target/release/insertion_sort
