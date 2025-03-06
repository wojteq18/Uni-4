use std::collections::VecDeque; //dodanie dwustronnej kolejki

pub struct Queue<T> {
    elements: VecDeque<T>,
}

impl<T> Queue<T> {
    pub fn new() -> Self {
        Queue {
            elements: VecDeque::new(),
        }
    }

    pub fn push(&mut self, element: T) {
        self.elements.push_back(element);
    }

    pub fn out(&mut self) -> Option<T> {
        self.elements.pop_front()
    }

    pub fn len(&self) -> usize {
        self.elements.len()
    }
}