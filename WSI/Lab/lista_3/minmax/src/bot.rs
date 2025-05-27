use rand::Rng;

use crate::board::WIN_MOVES;
use crate::board::LOSE_MOVES;

struct GameState {
    my_fields: Vec<u32>,
    enemy_fields: Vec<u32>,
    available_fields: Vec<u32>,
    depth: u32,
    maximizing_player: bool,
}

//funkcja oceny heurytycznej
fn evaluate_state(my_fields: &Vec<u32>, enemy_fields: &Vec<u32>) -> i32 {
    let mut score = 0;

    //Sprawdzamy warunki natychmiastowej wygranej dla nas
    if let Some(_win_move) = try_win(my_fields, enemy_fields) {
        score += 10000;
    }

    //Sprawdzamy warunki natychmiastowej przegranej dla przeciwnika
    if let Some(_lose_move) = try_win(enemy_fields, my_fields) {
        score -= 10000;
    }

    // Logika tworzenia lini3- elementowych
    for line in LOSE_MOVES.iter() {
        let in_my_line = line.iter().filter(|x| my_fields.contains(x)).count();
        let in_enemy_line = line.iter().filter(|x| enemy_fields.contains(x)).count();

        if in_my_line == 3 {
            score -= 1000; // Trzy pola zajęte przez mnie - przergana
        }

        if in_enemy_line == 3 {
            score += 1000; // Trzy pola zajęte przez przeciwnika - wygrana
        }

        //nagradzaj rozwój 2- elementowych linii
        match (in_my_line, in_enemy_line) {
            (2, 0) => score += 500,
            (0, 2) => score -= 500,
            _ => (),
        }
        return score;
    }

    //Ogólna ocena pozycji
    score += (my_fields.len() as i32 - enemy_fields.len() as i32) * 50;
    return score;
}

fn minmax(state: &GameState, depth: i32, mut alpha: i32, mut beta: i32) -> (i32, Option<u32>) {
    if state.depth == 0 || state.available_fields.is_empty() {
        return (evaluate_state(&state.my_fields, &state.enemy_fields), None);
    }
    let mut best_move = None;
    let mut best_value; //na początku moży być to max int albo mni int 
    if state.maximizing_player {
        best_value = i32::MIN;
    } else {
        best_value = i32::MAX;
    }

    for &mv in state.available_fields.iter() {
        let new_state = GameState {
            my_fields: if state.maximizing_player {
                let mut fields = state.my_fields.clone();
                fields.push(mv);
                fields
            } else {
                state.my_fields.clone()
            },
            enemy_fields: if !state.maximizing_player {
                let mut fields = state.enemy_fields.clone();
                fields.push(mv);
                fields
            } else {
                state.enemy_fields.clone()
            },
            available_fields: state.available_fields.iter().filter(|&&x| x != mv).cloned().collect(),
            depth: state.depth - 1,
            maximizing_player: !state.maximizing_player,
        };
        let (current_value, _) = minmax(&new_state, depth - 1, alpha, beta);
        
        if state.maximizing_player {
            if current_value > best_value {
                best_value = current_value;
                best_move = Some(mv);
                alpha = alpha.max(best_value); //najlepsze, dotychczas gwarantowane minimum
            }
        } else {
            if current_value < best_value {
                best_value = current_value;
                best_move = Some(mv);
                beta = beta.min(best_value); //najlepsze, dotychczas gwarantowane maksimum
            }
        }
        if beta <= alpha {
            break; //przycinanie alfa-beta
        }
    }

    return (best_value, best_move);
}

pub fn choose_best_move(my_fields: &Vec<u32>, enemy_fields: &Vec<u32>, available_fields: &Vec<u32>, depth: u32) -> Option<u32> {
    let game_state = GameState {
        my_fields: my_fields.clone(),
        enemy_fields: enemy_fields.clone(),
        available_fields: available_fields.clone(),
        depth,
        maximizing_player: true, // Zakładamy, że to my gramy
    };

    let (best_value, best_move) = minmax(&game_state, depth as i32, i32::MIN, i32::MAX);
    if best_value > 0 {
        return best_move;
    } else {
        return None;
    }
}


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


