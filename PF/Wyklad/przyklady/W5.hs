--Struktura osobaa, deriving show automatycznie tworzy funkcję show
data Osoba = Osoba
    { imie :: String
    , nazwisko :: String
    , wiek :: Int
    } deriving (Show)


zmienWiek ::Osoba -> Int -> Osoba
zmienWiek (Osoba imie nazwisko wiek) nowyWiek = Osoba imie nazwisko nowyWiek    
