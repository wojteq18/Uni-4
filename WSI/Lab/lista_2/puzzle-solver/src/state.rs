use crate::constants::PUZZLE_SIZE;
use crate::board::Board;
use priority_queue::PriorityQueue;
use std::cmp::Reverse;
use std::rc::Rc;
use std::collections::HashMap;

#[derive(Debug, Eq, PartialEq, Hash)]
pub struct State {
    pub board: Board,
    pub cost: usize,             // liczba ruchów wykonanych od początku
    pub estimated_cost: usize,   // heurystyka + cost, ustala priorytet
    pub parent: Option<Rc<State>>, // wskaźnik na poprzedni stan
}

impl State {
    pub fn new(board: Board, cost: usize, parent: Option<Rc<State>>) -> Self {
        let heuristic = board.manhattan_distance();
        let estimated_cost = cost + heuristic;
        State {
            board,
            cost,
            estimated_cost,
            parent,
        }
    }
}

pub fn fix(board: &mut Board) {
    let mut cost_map = HashMap::new();
    let mut queue = PriorityQueue::<Rc<State>, Reverse<usize>>::new();

    let initial_state = Rc::new(State::new(board.clone(), 0, None));
    queue.push(initial_state.clone(), Reverse(initial_state.estimated_cost));
    cost_map.insert(board.clone(), 0);

    while let Some((state_rc, _)) = queue.pop() {
        let state = state_rc.as_ref();

        // Sprawdź, czy aktualny stan ma minimalny znany koszt
        if let Some(&saved_cost) = cost_map.get(&state.board) {
            if state.cost > saved_cost {
                continue;
            }
        }

        if state.board.how_many_correct() == PUZZLE_SIZE {
            println!("Found solution with cost: {}", state.cost);
            println!("States: {}", cost_map.len());
            print_solution_path(&state_rc);
            return;
        }

        for next_move in state.board.find_movable_piece() {
            let mut next_board = state.board.clone();
            next_board.swap(state.board.get_zero_element(), next_move);

            let new_cost = state.cost + 1;

            // Sprawdź, czy nowy stan jest lepszy od dotychczasowego
            match cost_map.get(&next_board) {
                Some(&existing_cost) if new_cost >= existing_cost => continue,
                _ => {
                    // Dodaj lub aktualizuj koszt w mapie
                    cost_map.insert(next_board.clone(), new_cost);
                }
            }

            let next_state = Rc::new(State::new(
                next_board.clone(),
                new_cost,
                Some(state_rc.clone()),
            ));

            queue.push(next_state.clone(), Reverse(next_state.estimated_cost));
        }
    }

    println!("No solution found.");
}

fn print_solution_path(end_state: &Rc<State>) {
    let mut path = Vec::new();
    let mut current = Some(end_state.clone());

    while let Some(state_rc) = current {
        path.push(state_rc.board.clone());
        current = state_rc.parent.clone();
    }

    for board in path.iter().rev() {
        board.print();
        println!();
    }
}
