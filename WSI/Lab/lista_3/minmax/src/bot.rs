use rand::Rng;

use crate::board::{WIN_MOVES, LOSE_MOVES, BEST_FIELD, WORST_FIELDS, SECOND_WORST_FIELDS, SECOND_BEST_FIELDS,
BEST_FIELD_SCORE, SECOND_BEST_SCORE, SECOND_WORST_SCORE, WORST_SCORE};

pub(crate) const ME: u8 = 1;
pub(crate) const ENEMY: u8 = 2;
pub(crate) const EMPTY: u8 = 0;

pub fn field_to_coords(field: u32) -> (usize, usize) {
    if field < 11 || field > 55 {
        panic!("Nieprawidłowy numer pola: {}", field);
    }
    let row  = ((field / 10) - 1) as usize;
    let col = ((field % 10) - 1) as usize;
    (row, col)
}

fn get_available_moves(board: &[[u8; 5]; 5]) -> Vec<u32> {
    let mut moves = Vec::new();
    for r in 0..5 {
        for c in 0..5 {
            if board[r][c] == EMPTY {
                moves.push(((r + 1) * 10 + (c + 1)) as u32);
            }
        }
    }
    moves
}

//funkcja oceny heurytycznej
fn evaluate_state(board: &[[u8; 5]; 5]) -> i32 {
    let mut score = 0;

    // Ocena gróźb wygranej
    if try_win(board, ME).is_some() { score += 40000; }
    if try_win(board, ENEMY).is_some() { score -= 40000; }
    
    // Ocena par w liniach do wygranej (2 w linii na 4)
    for line4 in WIN_MOVES.iter() {
        let mut my_markers = 0;
        let mut enemy_markers = 0;
        for &field in line4.iter() {
            let (r, c) = field_to_coords(field);
            if board[r][c] == ME { my_markers += 1; }
            else if board[r][c] == ENEMY { enemy_markers += 1; }
        }
        if my_markers == 2 && enemy_markers == 0 { score += 7000; }
        if enemy_markers == 2 && my_markers == 0 { score -= 7000; }
    }

    // Ocena groźby przegranej (2 w linii na 3)
    for line3 in LOSE_MOVES.iter() {
        let mut my_markers = 0;
        let mut enemy_markers = 0;
        for &field in line3.iter() {
            let (r, c) = field_to_coords(field);
            if board[r][c] == ME { my_markers += 1; }
            else if board[r][c] == ENEMY { enemy_markers += 1; }
        }
        if my_markers == 2 && enemy_markers == 0 { score -= 1000; } // Kara za groźbę samobója
        if enemy_markers == 2 && my_markers == 0 { score += 1000; } // Nagroda za groźbę samobója przeciwnika
    }

    // Ocena pozycyjna
    for r in 0..5 {
        for c in 0..5 {
            let field_val = ((r + 1) * 10 + (c + 1)) as u32;
            let field_score = if BEST_FIELD == field_val { BEST_FIELD_SCORE }
                              else if SECOND_BEST_FIELDS.contains(&field_val) { SECOND_BEST_SCORE }
                              else if SECOND_WORST_FIELDS.contains(&field_val) { SECOND_WORST_SCORE }
                              else if WORST_FIELDS.contains(&field_val) { WORST_SCORE }
                              else { 0 };
            
            if board[r][c] == ME { score += field_score; }
            else if board[r][c] == ENEMY { score -= field_score; }
        }
    }

    score
}

fn minmax(board: &mut [[u8; 5]; 5], player: u8, depth: i32, mut alpha: i32, mut beta: i32) -> (i32, Option<u32>) {
    if check_for_four_completed(board, ME) { return (100000, None); } 
    if check_for_four_completed(board, ENEMY) { return (-100000, None); }
    if check_for_three_completed(board, ME) { return (-50000, None); }
    if check_for_three_completed(board, ENEMY) { return (50000, None); }

    let available_moves = get_available_moves(board);
    if depth == 0 || available_moves.is_empty() {
        return (evaluate_state(board), None);
    }

    let mut best_move = None;
    let opponent = if player == ME { ENEMY } else { ME };

    if player == ME { // Tura gracza MAXymalizującego
        let mut max_eval = i32::MIN;
        for mv in available_moves {
            let (r, c) = field_to_coords(mv);
            
            board[r][c] = player; // Zrób ruch
            let (eval, _) = minmax(board, opponent, depth - 1, alpha, beta);
            board[r][c] = EMPTY;  // Cofnij ruch!

            if eval > max_eval {
                max_eval = eval;
                best_move = Some(mv);
            }
            alpha = alpha.max(eval);
            if beta <= alpha {
                break;
            }
        }
        (max_eval, best_move)
    } else { // Tura gracza MINimalizującego
        let mut min_eval = i32::MAX;
        for mv in available_moves {
            let (r, c) = field_to_coords(mv);

            board[r][c] = player; 
            let (eval, _) = minmax(board, opponent, depth - 1, alpha, beta);
            board[r][c] = EMPTY;  

            if eval < min_eval {
                min_eval = eval;
                best_move = Some(mv);
            }
            beta = beta.min(eval);
            if beta <= alpha {
                break;
            }
        }
        (min_eval, best_move)
    }
}

pub fn choose_best_move(my_fields: &Vec<u32>, enemy_fields: &Vec<u32>, depth: u32) -> Option<u32> {
    let mut board = [[EMPTY; 5]; 5];
    for &field in my_fields.iter() {
        let (r, c) = field_to_coords(field);
        board[r][c] = ME;
    }
    for &field in enemy_fields.iter() {
        let (r, c) = field_to_coords(field);
        board[r][c] = ENEMY;
    }

    let (_score, best_move) = minmax(&mut board, ME, depth as i32, i32::MIN, i32::MAX);
    best_move
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

pub fn try_win(board: &[[u8; 5]; 5], player: u8) -> Option<u32> {
    let opponent = if player == ME { ENEMY } else { ME };
    for win_line in WIN_MOVES.iter() {
        let mut count = 0;
        let mut empty_field = None;
        let mut possible = true;
        for &field in win_line.iter() {
            let (r, c) = field_to_coords(field);
            if board[r][c] == player {
                count += 1;
            } else if board[r][c] == opponent {
                possible = false;
                break;
            } else { // Puste pole
                if empty_field.is_none() {
                    empty_field = Some(field);
                } else { // Więcej niż jedno puste pole
                    possible = false;
                    break;
                }
            }
        }
        if possible && count == 3 && empty_field.is_some() {
            return empty_field;
        }
    }
    None
}

fn check_for_four_completed(board: &[[u8; 5]; 5], player: u8) -> bool {
    for win_line in WIN_MOVES.iter() {
        if win_line.iter().all(|&field| {
            let (r, c) = field_to_coords(field);
            board[r][c] == player
        }) {
            return true;
        }
    }
    false
}

fn check_for_three_completed(board: &[[u8; 5]; 5], player: u8) -> bool {
    for lose_line in LOSE_MOVES.iter() {
        if lose_line.iter().all(|&field| {
            let (r, c) = field_to_coords(field);
            board[r][c] == player
        }) {
            return true;
        }
    }
    false
}

//fn check_for_four_completed(pl)
pub fn check_lose(board: &[[u8; 5]; 5], mv: u32, player: u8) -> bool {
    let (r, c) = field_to_coords(mv);
    if board[r][c] != EMPTY {
        return false; 
    }

    let mut temp_board = *board;
    temp_board[r][c] = player;

    check_for_three_completed(&temp_board, player)
}


