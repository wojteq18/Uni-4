use rand;
use rand::Rng;
use core::time;
use std::{ops::{AddAssign, Sub}, time::Duration};
use waitgroup::WaitGroup;
use once_cell::sync::Lazy;
use std::sync::Mutex;
use std::time::Instant;
use array_init::array_init;
use std::time::{SystemTime, UNIX_EPOCH};
use crossbeam::channel::{unbounded, Sender, Receiver, select};
use std::thread;


//const 
const NR_OF_TRAVELERS: i32 = 15;
const MIN_STEP: usize = 10;
const MAX_STEP: usize = 100;

const MIN_DELAY: Duration = Duration::from_millis(10);
const MAX_DELAY: Duration = Duration::from_millis(50);

const BOARD_WIDTH: usize = 15;
const BOARD_HEIGHT: usize = 15;


static START_TIME: Lazy<Instant> = Lazy::new(|| {
    Instant::now()
});
static BOARD: Lazy<[[Mutex<()>; BOARD_HEIGHT]; BOARD_WIDTH]> = Lazy::new(|| {
    array_init(|_| array_init(|_| Mutex::new(())))
});

#[derive(Clone, Copy)]
struct Position {
    x: usize,
    y: usize,
}

impl Position {
    fn move_down(&mut self) {
        self.y = (self.y + 1) % BOARD_HEIGHT;
    }

    fn move_up(&mut self) {
        self.y = (self.y + BOARD_HEIGHT - 1) % BOARD_HEIGHT;
    }

    fn move_right(&mut self) {
        self.x = (self.x + 1) % BOARD_WIDTH;
    }

    fn move_left(&mut self) {
        self.x = (self.x + BOARD_WIDTH - 1) % BOARD_WIDTH;
    }

    fn copy(&self) -> Position {
        Position {
            x: self.x,
            y: self.y,
        }
    }
}
#[derive(Clone, Copy)]
struct Trace_Type{
    time_stamp: Duration,
    id: usize,
    position: Position,
    symbol: char,
}

impl Trace_Type {
    /*fn copy(&self) -> Trace_Type {
        Trace_Type {
            time_stamp: self.time_stamp,
            id: self.id,
            position: Position { x: self.position.x, y: self.position.y },
            symbol: self.symbol,
        }
    }*/

    fn default() -> Trace_Type {
        Trace_Type {
            time_stamp: Duration::new(0, 0),
            id: 0,
            position: Position { x: 0, y: 0 },
            symbol: ' ',
        }
    }
}

type Trace_Array = [Trace_Type; MAX_STEP + 1];
#[derive(Clone)]
struct Traces_Sequence_Type {
    last: Option<usize>,
    trace_array: Trace_Array,
}

impl Traces_Sequence_Type {
    fn new() -> Self {
        Traces_Sequence_Type {
            last: None,
            trace_array: array_init(|_| Trace_Type::default()),
        }
    }

    fn push(&mut self, trace: Trace_Type) {
        let next = self.last.map_or(0, |last| last + 1);
        self.trace_array[next] = trace;
        self.last = Some(next);
    }
}



fn print_trace(t: Trace_Type) {
    let elapsed_secs: f64 = t.time_stamp.as_secs_f64();
    print!("{:.6} {} {} {} {} \n", elapsed_secs, t.id, t.position.x, t.position.y, t.symbol);
}

fn print_traces(t: Traces_Sequence_Type) {
    if let Some(last) = t.last {
        for i in 0..=last {
            print_trace(t.trace_array[i]);
        }
    }
}


static REPORT_CHANNEL: Lazy<(Sender<Traces_Sequence_Type>, Receiver<Traces_Sequence_Type>)> = Lazy::new(|| {
    unbounded() //unbounded narzuca brak limitu bufora
});

fn printer() {
    for _ in 0..(NR_OF_TRAVELERS) as usize {
        let traces = REPORT_CHANNEL.1.recv().unwrap(); //odwołujemy się do receivera
        print_traces(traces);
    }
}

struct Traveler {
    id: usize,
    symbol: char,
    position: Position,
}

fn acquire_squere_with_timeout(x: usize, y: usize, time: Duration) -> Option<std::sync::MutexGuard<'static, ()>> {
    let ceil = &BOARD[x][y];
    let start = Instant::now();
    loop {
        if let Ok(guard) = ceil.try_lock() {
            return Some(guard);
        } else if start.elapsed() > time {
            return None;
        }
        thread::sleep(Duration::from_millis(1));
    }
}

fn traveler(id: usize, symbol: char, seed: usize) {
    let mut rng = rand::thread_rng();
    let mut traveler = Traveler {
        id,
        symbol,
        position: Position { x: seed % BOARD_WIDTH, y: seed % BOARD_HEIGHT },
    };

    let mut traces = Traces_Sequence_Type::new();
    //let ceil = &BOARD[traveler.position.x][traveler.position.y];
    //let _lock = ceil.lock().unwrap();

    traces.push(Trace_Type {
        time_stamp: START_TIME.elapsed(),
        id: traveler.id,
        position: traveler.position.copy(),
        symbol: traveler.symbol,
    });

    let nr_of_steps = rng.gen_range(MIN_STEP..=MAX_STEP);

    let deadlok_timeout = 2 * MAX_DELAY;
    let mut is_deadlock = false;
    let mut guard = BOARD[traveler.position.x][traveler.position.y].lock().unwrap();


    for _ in 0..nr_of_steps {
        let range = (MAX_DELAY - MIN_DELAY).as_millis() as u64;
        let random_offset = Duration::from_millis(rng.gen_range(0..range));
        let delay = MIN_DELAY + random_offset;


        std::thread::sleep(delay);
        let old_position = traveler.position.copy();

        match rng.gen_range(0..=3) {
            0 => traveler.position.move_down(),
            1 => traveler.position.move_up(),
            2 => traveler.position.move_right(),
            _ => traveler.position.move_left(),
        }

        let new_x = traveler.position.x;
        let new_y = traveler.position.y;

        if let Some(new_guard) = acquire_squere_with_timeout(new_x, new_y, deadlok_timeout) {
            drop(guard); 
            guard = new_guard;
        } else {
            is_deadlock = true;
            traveler.symbol = 'x';
            traveler.position = old_position;
            //BOARD[new_x][new_y].lock().unwrap();

            traces.push(Trace_Type {
                time_stamp: START_TIME.elapsed(),
                id: traveler.id,
                position: traveler.position.copy(),
                symbol: traveler.symbol,
            });
            REPORT_CHANNEL.0.send(traces.clone()).unwrap();
            return;
        }
        traces.push(Trace_Type {
            time_stamp: START_TIME.elapsed(),
            id: traveler.id,
            position: traveler.position.copy(),
            symbol: traveler.symbol,
        });
    }
    if !is_deadlock {
        REPORT_CHANNEL.0.send(traces).unwrap();
    }
}

fn main() {
    println!("-1 {} {} {}", NR_OF_TRAVELERS, BOARD_WIDTH, BOARD_HEIGHT);

    let symbols: Vec<char> = vec![
        'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
        'I', 'J', 'K', 'L', 'M', 'N', 'O',
    ];

    let printer_handle = thread::spawn(|| {
        printer();
    });

    for i in 0..(NR_OF_TRAVELERS as usize) {
        let symbol = symbols[i];
        let seed = rand::random::<usize>();
        thread::spawn(move || {
            traveler(i, symbol, seed);
        });
    }

    printer_handle.join().unwrap();
}
