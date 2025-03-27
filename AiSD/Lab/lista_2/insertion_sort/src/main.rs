use std::io::{self, BufRead};


fn insertion_sort (length: usize, arr: &mut [usize]) {
    for i in 1..arr.len() {
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


            insertion_sort(length, &mut numbers);
            println!("{:?}", numbers);
        }
    }    
}
// /home/wojteq18/Uni/AiSD/Lab/lista_2/generators/random_sequence/target/release/random_sequence 7 | /home/wojteq18/Uni/AiSD/Lab/lista_2/insertion_sort/target/release/insertion_sort
