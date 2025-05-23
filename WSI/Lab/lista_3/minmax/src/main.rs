mod board;

use std::io::{BufReader, BufRead, Write};
use std::net::TcpStream;
use std::thread;
use std::sync::{Arc, Mutex};
use board::{set_board, set_move, print_board};

fn main() -> std::io::Result<()> {
    let args: Vec<String> = std::env::args().collect();
    if args.len() != 6 {
        eprintln!("Użycie: {} <nazwa_użytkownika>", args[0]);
        std::process::exit(1);
    }
    let ip_adress = args[1].clone();
    let port = args[2].clone();
    let player_number  = args[3].clone();
    let username = args[4].clone();
    let _deepness = args[5].clone();

    let final_adress = format!("{}:{}", ip_adress, port);
    
    set_board(); // Inicjalizacja planszy
    let mut end_game = false;



    // Łączenie z serwerem
    let mut stream = TcpStream::connect(final_adress)?;
    println!("Połączono z serwerem");

    //Wysyałanie identyfikacji
    let ident = format!("{} {}\n", player_number, username);
    stream.write_all(ident.as_bytes())?;

    let mut reader = BufReader::new(stream.try_clone()?);
    let stdin = std::io::stdin();

    while !end_game {
        let mut server_msg = String::new();
        reader.read_line(&mut server_msg)?;

        if server_msg.trim().is_empty() {
            println!("Serwer zakończył połączenie");
            break;
        }

        println!("Serwer: {}", server_msg.trim());

        let num: i32 = server_msg.trim().parse().unwrap_or(-1);
        let move_code = num % 100;
        let msg_code = num / 100;

        if move_code != 0 {
            set_move(move_code as usize, 3 - player_number.parse::<usize>().unwrap());
            print_board();
        }

        if msg_code == 0 || msg_code == 6 {
            print!("Twój ruch: ");
            std::io::stdout().flush()?;

            let mut input = String::new();
            stdin.lock().read_line(&mut input)?;
            let input = input.trim();

            let mv: usize = input.parse().unwrap_or(0);
            set_move(mv, player_number.parse().unwrap());
            print_board();

            write!(stream, "{}", mv)?;
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