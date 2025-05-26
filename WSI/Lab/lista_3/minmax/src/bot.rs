use rand::Rng;

use crate::board::WIN_MOVES;
use crate::board::LOSE_MOVES;

pub fn all_fields() -> Vec<u32> {
    let mut fields = Vec::new();
    fields.extend([
        11, 12, 13, 14, 15,
        21, 22, 23, 24, 25,
        31, 32, 33, 34, 35,
        41, 42, 43, 44, 45,
        51, 52, 53, 54, 55,
    ]);
    fields
}

pub fn delete_field(fields: &mut Vec<u32>, value: u32) {
    if let Some(index) = fields.iter().position(|&x| x == value) {
        fields.remove(index);
    }
}


pub fn random_field(fields: &Vec<u32>) -> Option<u32> {
    if fields.is_empty() {
        None
    } else {
        let index = rand::rng().random_range(0..fields.len());
        Some(fields[index])
    }
}

pub fn try_win(my_fields: &Vec<u32>, enemy_fields: &Vec<u32>) -> Option<u32> {
    for win_move in WIN_MOVES.iter() {
        let (is_win, remember) = check_tree(my_fields, enemy_fields, *win_move);
        if is_win {
            return Some(remember);
        }
    }
    None
}

fn check_tree(my_fields: &Vec<u32>, enemy_fields: &Vec<u32>, win_move: [u32; 4]) -> (bool, u32) {
    let mut count = 0;
    let mut remember = 0;
    for line in win_move.iter() {
        if my_fields.contains(line) {
            count += 1;
        } else {
            remember = *line;
        }
    }
    if count == 3 && !enemy_fields.contains(&remember) {
        return (true, remember);
    }
    (false, 0)
}

pub fn check_lose(my_fields: &Vec<u32>, player_move: u32) -> bool {
    for lose_move in LOSE_MOVES.iter() {
        let mut count = 0;
        for line in lose_move.iter() {
            if my_fields.contains(line) || *line == player_move {
                count += 1;
            }
        }
        if count == 3 {
            return true; // Przeciwnik może wygrać w następnym ruchu
        }
    }
    false
}


