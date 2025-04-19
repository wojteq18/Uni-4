mod board;
mod constants;
mod state;

use board::*;
use state::*;

fn main() { 
    let mut board = Board::new();
    board.test();
    //board.shuffle();
    board.print();
    let correct = board.how_many_correct();
    println!("Correct: {}", correct);
    let _ = fix(&mut board);
}