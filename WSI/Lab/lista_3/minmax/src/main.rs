
mod board;
mod bot;

use std::io::{Read, Write};
use std::net::TcpStream;
use board::{set_board, set_move, print_board};
use bot::{all_fields, delete_field, random_field, check_lose, choose_best_move, try_win, ME, ENEMY, EMPTY, field_to_coords};

fn main() -> std::io::Result<()> {
    let mut all_fields = all_fields();
    let mut my_fields: Vec<u32> = Vec::new();
    let mut enemy_fields: Vec<u32> = Vec::new();

    let args: Vec<String> = std::env::args().collect();
    if args.len() != 6 {
        eprintln!("Użycie: {} <ip> <port> <nr_gracza> <nazwa> <głębokość>", args[0]);
        std::process::exit(1);
    }
    let ip_adress = args[1].clone();
    let port = args[2].clone();
    let player_number_str  = args[3].clone();
    let username = args[4].clone();
    let deepness = args[5].clone();

    let player_number: usize = player_number_str.parse().expect("Numer gracza musi byc liczba calkowita!");
    let final_adress = format!("{}:{}", ip_adress, port);
    
    set_board();
    let mut end_game = false;

    // Łączenie z serwerem
    let mut stream = TcpStream::connect(final_adress)?;
    println!("Połączono z serwerem");

    // Wysyłanie identyfikacji
    let mut buffer = [0; 16];
    let bytes_read = stream.read(&mut buffer)?;
    let server_msg = String::from_utf8_lossy(&buffer[..bytes_read]);
    if server_msg.trim() == "700" {
        let ident = format!("{} {}", player_number, username);
        stream.write_all(ident.as_bytes())?;
        stream.flush()?;
    }

    while !end_game {

        let mut buffer = [0; 16];
        let bytes_read = match stream.read(&mut buffer) {
            Ok(0) => { println!("Serwer zakończył połączenie."); break; }
            Ok(n) => n,
            Err(e) => { eprintln!("Błąd odczytu z serwera: {}", e); break; }
        };

        let server_msg = String::from_utf8_lossy(&buffer[..bytes_read]);
        let server_msg = server_msg.trim_matches('\u{0}').trim();


        if server_msg.is_empty() {
            continue;
        }

        println!("Serwer: {}", server_msg);
        let server_usize: usize = server_msg.parse().unwrap_or(0);
        if server_usize != 0 && server_usize != 600 {
            enemy_fields.push(server_usize as u32);
            delete_field(&mut all_fields, server_usize as u32);
        }


        let num: i32 = server_msg.parse().unwrap_or(-1);
        let move_code = num % 100;
        let msg_code = num / 100;

        if move_code != 0 {
            set_move(move_code as usize, 3 - player_number);
            print_board();
        }

        if msg_code == 0 || msg_code == 6 {

            let mut board = [[EMPTY; 5]; 5];
            for &field in my_fields.iter() {
                let (r, c) = field_to_coords(field);
                board[r][c] = ME;
            }
            for &field in enemy_fields.iter() {
                let (r, c) = field_to_coords(field);
                board[r][c] = ENEMY;
            }

            if my_fields.len() + enemy_fields.len() == 0 {
                my_fields.push(33);
                delete_field(&mut all_fields, 33);
                set_move(33, player_number);
                print_board();
                stream.write_all("33".as_bytes())?;
                stream.flush()?;
                continue;
            }

            

            if let Some(win_move) = try_win(&board, ME) {

                my_fields.push(win_move);
                delete_field(&mut all_fields, win_move);
                set_move(win_move as usize, player_number);
                print_board();
                stream.write_all(win_move.to_string().as_bytes())?;
                stream.flush()?;
                continue;
            }

            if let Some(block_move) = try_win(&board, ENEMY) {
                my_fields.push(block_move);
                delete_field(&mut all_fields, block_move);
                set_move(block_move as usize, player_number);
                print_board();
                stream.write_all(block_move.to_string().as_bytes())?;
                stream.flush()?;
                continue;
            }
            
            match choose_best_move(&my_fields, &enemy_fields, deepness.parse::<u32>().unwrap_or(3)) {
                Some(best_move) => {
                    my_fields.push(best_move);
                    delete_field(&mut all_fields, best_move);
                    set_move(best_move as usize, player_number);
                    print_board();
                    stream.write_all(best_move.to_string().as_bytes())?;
                    stream.flush()?;
                }
                None => { 
                    loop {
                        let random = random_field(&all_fields).unwrap_or(0);
                        if random == 0 { println!("Brak ruchów!"); break; }
                        
                        if check_lose(&board, random, ME) {
                            continue; 
                        } else {
                            my_fields.push(random);
                            delete_field(&mut all_fields, random);
                            set_move(random as usize, player_number);
                            print_board();
                            stream.write_all(random.to_string().as_bytes())?;
                            stream.flush()?;
                            break;
                        }
                    }
                }     
            }
        } else {
            end_game = true;
            match msg_code {
                1 => println!("Wygrałeś."),
                2 => println!("Przegrałeś."),
                3 => println!("Remis."),
                4 => println!("Wygrałeś. Błąd przeciwnika."),
                5 => println!("Przegrałeś. Twój błąd."),
                _ => println!("Nieznany status zakończenia gry."),
            }
        }
    }
    Ok(())
}