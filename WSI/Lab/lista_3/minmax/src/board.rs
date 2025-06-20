pub const BOARD_SIZE: usize = 5;

static mut BOARD: [[u8; BOARD_SIZE]; BOARD_SIZE] = [[0; BOARD_SIZE]; BOARD_SIZE];

pub(crate) const WIN_MOVES: [[u32; 4]; 28] = [
    // Poziome 
    [11, 12, 13, 14],
    [12, 13, 14, 15],
    [21, 22, 23, 24],
    [22, 23, 24, 25],
    [31, 32, 33, 34],
    [32, 33, 34, 35],
    [41, 42, 43, 44],
    [42, 43, 44, 45],
    [51, 52, 53, 54],
    [52, 53, 54, 55],

    // Pionowe 
    [11, 21, 31, 41],
    [21, 31, 41, 51],
    [12, 22, 32, 42],
    [22, 32, 42, 52],
    [13, 23, 33, 43],
    [23, 33, 43, 53],
    [14, 24, 34, 44],
    [24, 34, 44, 54],
    [15, 25, 35, 45],
    [25, 35, 45, 55],

    // Skosy
    [11, 22, 33, 44],
    [12, 23, 34, 45],
    [21, 32, 43, 54],
    [22, 33, 44, 55],
    [14, 23, 32, 41],
    [15, 24, 33, 42],
    [24, 33, 42, 51],
    [25, 34, 43, 52],
];

pub(crate) const LOSE_MOVES: [[u32; 3]; 48] = [
    // Poziome
    [11, 12, 13], [12, 13, 14], [13, 14, 15],
    [21, 22, 23], [22, 23, 24], [23, 24, 25],
    [31, 32, 33], [32, 33, 34], [33, 34, 35],
    [41, 42, 43], [42, 43, 44], [43, 44, 45],
    [51, 52, 53], [52, 53, 54], [53, 54, 55],

    // Pionowe 
    [11, 21, 31], [21, 31, 41], [31, 41, 51],
    [12, 22, 32], [22, 32, 42], [32, 42, 52],
    [13, 23, 33], [23, 33, 43], [33, 43, 53],
    [14, 24, 34], [24, 34, 44], [34, 44, 54],
    [15, 25, 35], [25, 35, 45], [35, 45, 55],

    // Skosy 
    [13, 24, 35],
    [12, 23, 34], [23, 34, 45],
    [11, 22, 33], [22, 33, 44], [33, 44, 55],
    [21, 32, 43], [32, 43, 54],
    [31, 42, 53],

    // Skosy 
    [13, 22, 31],
    [14, 23, 32], [23, 32, 41],
    [15, 24, 33], [24, 33, 42], [33, 42, 51],
    [25, 34, 43], [34, 43, 52],
    [35, 44, 53],
];

//oceny ,,mocy" poszczególnych pól
pub(crate) const BEST_FIELD: u32 = 33;
pub(crate) const WORST_FIELDS: [u32; 4] = [11, 15, 51, 55];
pub(crate) const SECOND_WORST_FIELDS: [u32; 12] = [12, 13, 14, 21, 25, 31, 35, 41, 45, 52, 53, 54];
pub(crate) const SECOND_BEST_FIELDS: [u32; 8] = [22, 23, 24, 32, 34, 42, 43, 44];

//wagi pól
pub(crate) const BEST_FIELD_SCORE: i32 = 500;
pub(crate) const SECOND_BEST_SCORE: i32 = 250;
pub(crate) const SECOND_WORST_SCORE: i32 = -100; 
pub(crate) const WORST_SCORE: i32 = -300; 


pub fn set_board() {
    for i in 0..BOARD_SIZE {
        for j in 0..BOARD_SIZE {
            unsafe {
                BOARD[i][j] = 0;
            }
        }
    }
}

pub fn print_board() {
    print!(" 1 2 3 4 5\n");
    for i in 0..BOARD_SIZE {
        print!("{}\n", i + 1);
        for j in 0..BOARD_SIZE {
            match unsafe {BOARD[i][j]} {
                0 => print!(" . "),
                1 => print!(" X "),
                2 => print!(" O "),
                _ => print!(" ? "),
            }
        }
    }    
}

pub fn set_move(mv: usize, player: usize) -> bool {
    let i = (mv / 10) as isize - 1;
    let j = (mv % 10) as isize - 1;

    if i < 0 || i > 4 as isize || j < 0 || j >= BOARD_SIZE as isize {
        return false;
    }

    let (i, j) = (i as usize, j as usize);

    unsafe {
        if BOARD[i][j] != 0 {
            return false;
        }
        BOARD[i][j] = player as u8;
    }

    return true
}

