mod board;
mod bot;

use std::io::{Read, Write}; //odczyt i zapis w gnieździe sieciowym
use std::net::TcpStream;
use board::{set_board, set_move, print_board};
use bot::{all_fields, delete_field, random_field, try_win, check_lose};

fn main() -> std::io::Result<()> {
    let mut all_fields = all_fields(); // Inicjalizacja wszystkich pól
    let mut my_fields: Vec<u32> = Vec::new();
    let mut enemy_fields: Vec<u32> = Vec::new();

    let args: Vec<String> = std::env::args().collect();
    if args.len() != 6 {
        eprintln!("Użycie: {} <nazwa_użytkownika>", args[0]);
        std::process::exit(1);
    }
    let ip_adress = args[1].clone();
    let port = args[2].clone();
    let player_number_str  = args[3].clone();
    let username = args[4].clone();
    let _deepness = args[5].clone();

    let player_number: usize = player_number_str.parse().expect("Numer gracza musi byc liczba calkowita!");
    let final_adress = format!("{}:{}", ip_adress, port);
    
    set_board(); // Inicjalizacja planszy
    let mut end_game = false;



    // Łączenie z serwerem
    let mut stream = TcpStream::connect(final_adress)?;
    println!("Połączono z serwerem");

    //Wysyałanie identyfikacji
    let mut buffer = [0; 16];
    let bytes_read = stream.read(&mut buffer)?;
    let server_msg = String::from_utf8_lossy(&buffer[..bytes_read]);
    if server_msg.trim() == "700" {
        let ident = format!("{} {}", player_number, username);
        stream.write_all(ident.as_bytes())?;
        stream.flush()?;
    }



    while !end_game {
        let mut buffer = [0; 16]; // Utwórz bufor podobny do tego w C
        let bytes_read = match stream.read(&mut buffer) {
            Ok(0) => {
                println!("Serwer zakończył połączenie.");
                break;
            }
            Ok(n) => n,
            Err(e) => {
                eprintln!("Błąd odczytu z serwera: {}", e);
                break;
            }
        };

        // Konwertuj tylko odczytane bajty na string, usuwając puste bajty (null bytes)
        let server_msg = String::from_utf8_lossy(&buffer[..bytes_read]); //konwertuje format bajtów na string, usuwa niepotrzebne znaki
        let server_msg = server_msg.trim_matches('\u{0}').trim();

        if server_msg.is_empty() {
            continue;
        }

        println!("Serwer: {}", server_msg);
        //parsowanie wiadomości od serwera
        let server_usize: usize = server_msg.parse().unwrap_or(0);
        enemy_fields.push(server_usize as u32);
        delete_field(&mut all_fields, server_usize as u32);

        let num: i32 = server_msg.parse().unwrap_or(-1);
        let move_code = num % 100;
        let msg_code = num / 100;

        if move_code != 0 {
            set_move(move_code as usize, 3 - player_number);
            print_board();
        }

        if msg_code == 0 || msg_code == 6 {
            match try_win(&my_fields, &enemy_fields) {
                Some(win_move) => {
                    my_fields.push(win_move);
                    delete_field(&mut all_fields, win_move);
                    set_move(win_move as usize, player_number);
                    print_board();
                    stream.write_all(win_move.to_string().as_bytes())?;
                    stream.flush()?;
                }
                None => {
                    loop {
                        let random = random_field(&all_fields).unwrap_or(0);
                        if check_lose(&my_fields, random) {
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