use rand::Rng;

use crate::board::{WIN_MOVES, LOSE_MOVES, BEST_FIELD, WORST_FIELDS, SECOND_WORST_FIELDS, SECOND_BEST_FIELDS,
BEST_FIELD_SCORE, SECOND_BEST_SCORE, SECOND_WORST_SCORE, WORST_SCORE};

struct GameState {
    my_fields: Vec<u32>,
    enemy_fields: Vec<u32>,
    available_fields: Vec<u32>,
    depth: u32,
    maximizing_player: bool,
}

//funkcja oceny heurytycznej
fn evaluate_state(my_fields: &Vec<u32>, enemy_fields: &Vec<u32>) -> i32 {

    //Sprawdzamy warunki natychmiastowej wygranej dla nas
    if check_for_four_completed(my_fields) {
        return 100000;
    }

    //Sprawdzamy warunki natychmiastowej przegranej dla przeciwnika
    if check_for_four_completed(enemy_fields) {
        return -100000;
    }

    if check_for_three_completed(my_fields) {
        return -50000
    }

    if check_for_three_completed(enemy_fields) {
        return 50000;
    }

    let mut score = 0;

    if try_win(my_fields, enemy_fields).is_some() {
        score += 40000; // Nagroda za możliwość wygranej
    }

    if try_win(enemy_fields, my_fields).is_some() {
        score -= 40000; // Kara za możliwość przegranej
    }

    for line4 in WIN_MOVES.iter() {
        let my_markers_in_line4 = line4.iter().filter(|f| my_fields.contains(f)).count();
        let enemy_markers_in_line4 = line4.iter().filter(|f| enemy_fields.contains(f)).count();

        if my_markers_in_line4 == 2 && enemy_markers_in_line4 == 0 {
            // Dwa moje symbole i dwa puste pola w linii do wygranej
            score += 7000;
        }
        if enemy_markers_in_line4 == 2 && my_markers_in_line4 == 0 {
            // Dwa symbole przeciwnika i dwa puste pola w linii do jego wygranej
            score -= 7000;
        }
    }
    
    for line3 in LOSE_MOVES.iter() {
        let my_markers_in_line3 = line3.iter().filter(|f| my_fields.contains(f)).count();
        let enemy_markers_in_line3 = line3.iter().filter(|f| enemy_fields.contains(f)).count();

        if my_markers_in_line3 == 2 && enemy_markers_in_line3 == 0 {
             let empty_spots_count = line3.iter().filter(|f| !my_fields.contains(f) && !enemy_fields.contains(f)).count();
             if empty_spots_count == 1 { // Dokładnie jedno puste miejsce do utworzenia trójki
                score -= 1000;
             }
        }
        if enemy_markers_in_line3 == 2 && my_markers_in_line3 == 0 {
            let empty_spots_count = line3.iter().filter(|f| !my_fields.contains(f) && !enemy_fields.contains(f)).count();
            if empty_spots_count == 1 {
                score += 1000;
            }
        }
    }

    // Ocena moich pól
    let my_worst_fields = my_fields.iter().filter(|&&f| WORST_FIELDS.contains(&f)).count();
    let my_second_worst_fields = my_fields.iter().filter(|&&f| SECOND_WORST_FIELDS.contains(&f)).count();
    let my_second_best_fields = my_fields.iter().filter(|&&f| SECOND_BEST_FIELDS.contains(&f)).count();

    if my_fields.contains(&BEST_FIELD.clone()) {
        score += BEST_FIELD_SCORE;
    }
    score += (my_worst_fields as i32) * WORST_SCORE;
    score += (my_second_worst_fields as i32) * SECOND_WORST_SCORE;
    score += (my_second_best_fields as i32) * SECOND_BEST_SCORE;

    // Ocena pól przeciwnika
    let enemy_worst_fields = enemy_fields.iter().filter(|&&f| WORST_FIELDS.contains(&f)).count();
    let enemy_second_worst_fields = enemy_fields.iter().filter(|&&f| SECOND_WORST_FIELDS.contains(&f)).count();
    let enemy_second_best_fields = enemy_fields.iter().filter(|&&f| SECOND_BEST_FIELDS.contains(&f)).count();
    if enemy_fields.contains(&BEST_FIELD.clone()) {
        score -= BEST_FIELD_SCORE;
    }
    score -= (enemy_worst_fields as i32) * WORST_SCORE;
    score -= (enemy_second_worst_fields as i32) * SECOND_WORST_SCORE;
    score -= (enemy_second_best_fields as i32) * SECOND_BEST_SCORE;

    score
}

fn minmax(state: &GameState, depth: i32, mut alpha: i32, mut beta: i32) -> (i32, Option<u32>) {
    if check_for_three_completed(&state.my_fields) {
        return (-50000, None);
    }

    if check_for_three_completed(&state.enemy_fields) {
        return (50000, None);
    }

    if check_for_four_completed(&state.my_fields) {
        return (100000, None);
    }

    if state.depth == 0 || state.available_fields.is_empty() {
        return (evaluate_state(&state.my_fields, &state.enemy_fields), None);
    }

    if check_for_four_completed(&state.enemy_fields) {
        return (-100000, None);
    }
    
    let mut best_move = None;
    let mut best_value; //na początku moży być to max int albo min int 
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

    let (_best_value, best_move) = minmax(&game_state, depth as i32, i32::MIN, i32::MAX);
    return best_move; 
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
    for win_line in WIN_MOVES.iter() {
        let mut count = 0;
        let mut empty_field = None;
        let mut possible = true;
        for &field in win_line.iter() {
            if my_fields.contains(&field) {
                count += 1;
            } else if enemy_fields.contains(&field) {
                possible = false;
                break;
            } else {
                if empty_field.is_none() {
                    empty_field = Some(field);
                } else {
                    possible = false;
                    break;
                }
            }
        }
        if possible && count == 3 && empty_field.is_some() {
            return empty_field; // Zwracamy pole, które można zająć, aby wygrać
        }
    }
    return None;
}

fn check_for_four_completed(player_fields: &Vec<u32>) -> bool {
    for win_line in WIN_MOVES.iter() {
        if win_line.iter().all(|&x| player_fields.contains(&x)) {
            return true; 
        }
    }
    return false;
}

fn check_for_three_completed(player_fields: &Vec<u32>) -> bool {
    for loose_line in LOSE_MOVES.iter() {
        if loose_line.iter().all(|&x| player_fields.contains(&x)) {
            return true; 
        }
    }
    return false;
}

//fn check_for_four_completed(pl)
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


