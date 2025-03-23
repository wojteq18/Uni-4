use std::env;

fn quick_sort (length: usize, arr: &mut Vec<usize>)  -> Vec<usize> {
    if arr.len() <= 1 {
        return arr.clone();
    }

    let pivot = arr[0];
    let mut left_side = vec![];
    let mut middle_side = vec![];
    let mut right_side = vec![];

    for &item in arr.iter() { 
        if item < pivot {
            left_side.push(item);
        } else if item == pivot {
            middle_side.push(item);
        } else {
            right_side.push(item);
        }
    }

    let mut sorted_left = quick_sort(left_side.len(), &mut left_side);
    let mut sorted_right = quick_sort(right_side.len(), &mut right_side);

    sorted_left.append(&mut middle_side);
    sorted_left.append(&mut sorted_right);

    return sorted_left;
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

    let sorted = quick_sort(length, &mut arr);
    println!("{:?}", sorted);
}
