use crate::constants::{SIZE, PUZZLE_SIZE};
use rand::Rng;
//use rand::seq::SliceRandom; // Importujemy SliceRandom dla metody shuffle

/*#[derive(Clone, Copy, Debug, Eq, PartialEq, Hash)]

pub struct Field {
    value: usize,
    index: usize,   
}*/

#[derive(Clone, Debug, Eq, PartialEq, Hash, Copy)]
pub struct Board {
    fields: [usize; PUZZLE_SIZE],
    pub zero_position: usize,
    //how_many_correct: usize,
}

impl Board {

    pub fn get_zero_element(&self) -> usize {
        return self.zero_position;
    }

// Poprawiona wersja Board::find_movable_piece:
pub fn find_movable_piece(&self) -> Vec<usize> {
    let mut movable_indices = Vec::new(); 
    let col = self.zero_position % SIZE;
    let row = self.zero_position / SIZE;

    if row > 0 {
        movable_indices.push(self.zero_position - SIZE);
    }
    if row < SIZE - 1 {
        movable_indices.push(self.zero_position + SIZE); 
    }
    if col > 0 {
        movable_indices.push(self.zero_position - 1);   
    }
    if col < SIZE - 1 {
        movable_indices.push(self.zero_position + 1);   
    }
    return movable_indices;
}

    pub fn new() -> Self {
        let mut fields = [0; PUZZLE_SIZE];
        let zero_position = PUZZLE_SIZE - 1; // The position of the empty space (0)

        for i in 0..PUZZLE_SIZE {
            fields[i] = (i + 1) % PUZZLE_SIZE; // Fill the board with numbers 1 to 15 and 0

        }

        Board { fields, zero_position }
    }

    pub fn print(&self) {
        for i in 0..SIZE {
            for j in 0..SIZE {
                print!("{:^5}", self.fields[i * SIZE + j]);
            }
            println!();
        }
        println!();
    }

    /*pub fn shuffle(&mut self) {
        let mut numbers: Vec<usize> = (1..PUZZLE_SIZE).collect();
        let mut rng = rand::rng();
        numbers.shuffle(&mut rng);
        
        for i in 0..PUZZLE_SIZE - 1 {
            self.fields[i] = numbers[i];   
        }
    }*/

    pub fn shuffle(&mut self, i: i32) {
        let mut rng = rand::rng();

        for _ in 0..i {
            let movable_piece = self.find_movable_piece();
            let length = movable_piece.len();
            let random_index = rng.random_range(0..length);
            let random_value = movable_piece[random_index];
            self.swap(self.zero_position, random_value);
        }
    }

    pub fn swap(&mut self, index1: usize, index2: usize) {
        let value1 = self.fields[index1];
        let value2 = self.fields[index2];
    
        self.fields.swap(index1, index2);
    
        if value1 == 0 {
            self.zero_position = index2;
        } else if value2 == 0 {
            self.zero_position = index1;
        }
    }

    pub fn how_many_correct(&self) -> usize {
        let mut count = 0;
        for i in 0..PUZZLE_SIZE{
            if self.fields[i] == (i + 1) % PUZZLE_SIZE {
                count += 1;
            }
        }
        return count;
    }

    pub fn manhattan_distance(&self) -> usize {
        let mut distance: usize = 0;
        for i in 0..PUZZLE_SIZE {
            let value = self.fields[i];
            if value != 0 {
                let target_row = (value + PUZZLE_SIZE - 1) % PUZZLE_SIZE / SIZE;
                let target_col = ((value + PUZZLE_SIZE - 1) % PUZZLE_SIZE) % SIZE;
                let current_row = i / SIZE;
                let current_col = i % SIZE;
                distance += (target_row as isize - current_row as isize).abs() as usize
                    + (target_col as isize - current_col as isize).abs() as usize;
            }
        }
        return distance
    }

    pub fn test(&mut self) {
        self.fields = [15, 14, 13, 9, 8, 3, 12, 1, 7, 11, 4, 2, 10, 5, 6, 0];
        self.zero_position = 15;
    }     
}

