mod fifo_queue;
mod lifo_queue;
use fifo_queue::Queue as FifoQueue;
use lifo_queue::Queue as LifoQueue;
use std::process;

fn main() {
    let mut fifo_queue = FifoQueue::new();
    let mut lifo_queue = LifoQueue::new();

    for i in 0..50 {
        fifo_queue.push(i);
        lifo_queue.push(i);
    }
    let a: usize = fifo_queue.len();
    for _i in 0..a + 10 {
        if let Some(i) = lifo_queue.out() {
            println!("Element lifo: {}", i);

            if let Some(i) = fifo_queue.out() {
                println!("Element fifo: {}", i);
            }
        } else {
            println!("Error: Queue is empty");
            process::exit(1);
        }
    }
}