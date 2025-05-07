--zadanie 1
power x y = y ^ x

p2 = power 4
p3 = power 3

--zadanie 3
f :: Int -> Int
f x = x ^ 2

g :: Int -> Int -> Int
g x y = x + 2 * y

h x y = f (g x y)

--zadanie 13 a)
phi :: Int -> Int
phi n = length [k | k <- [1..n], gcd k n == 1]

--zadanie 13 b)
phi2 :: Int -> Int
phi2 n = sum [phi k | k <- [1..n], n `mod` k == 0]

--zadanie 14
isPerfect :: Int -> Bool
isPerfect n = n == sum [k | k <- [1..n-1], n `mod` k == 0]

allPerfect :: Int -> [Int]
allPerfect n = [k | k <- [1..n], isPerfect k]

--zadanie 15
sumOfDivisors :: Int -> Int
sumOfDivisors n = sum [k | k <- [1..n-1], n `mod` k == 0]

areSociable :: Int -> Int -> Bool
areSociable a b = sumOfDivisors a == b && sumOfDivisors b == a && a /= b

socialPairs :: Int -> [(Int, Int)]
socialPairs limit = [(a, b) | a <- [1..limit], let b = sumOfDivisors a, areSociable a b && a < b && b < limit]

--zadanie 16 a)
dcp1 :: Int -> Double
dcp1 n = fromIntegral a / fromIntegral b 
    where 
            a = length [(k, l) | k <- [1..n], l <- [1..n], gcd k l == 1]
            b = n ^ 2

--zadanie 16 b)
dcp' :: Int -> Double
dcp' n = fromIntegral (countCoprimes n n 1 1) / fromIntegral (n ^ 2)

countCoprimes :: Int -> Int -> Int -> Int -> Int
countCoprimes n m i j
    | i > m = 0 --koniec
    | j > m = countCoprimes n m (i + 1) 1 --przejscie do nowego weirsza
    | gcd i j == 1 = 1 + countCoprimes n m i (j + 1) --liczmy wzglednie pierwsza pare
    | otherwise = countCoprimes n m i (j + 1) --nie jest wzglednie 

--zadanie 16 c)        
dcpList :: [Double]
dcpList = [dcp1 k | k <- [100, 200..10000]]

--zadanie 17
nub' :: Eq a => [a] -> [a] --Eq a oznacza, że zawiera elementy porównywalne
nub' [] = [] --jeżeli lista jest pusta, to zwracamy pustą tablicę
nub' (x:xs) --wyciągamy pierwszy element z listy
    | x `elem` xs = nub' xs --jeżeli element już występuje w tablicy, to go pomijamy i rekurencyjnie wywołujemy funkcję dla reszty tablicy
    | otherwise = x : nub' xs --jeżeli element nie występuje w tablicy, to dodajemy go do tablicy wynikowej i rekurencyjnie wywołujemy funkcję dla reszty tablicy

--zadanie 18
inits' :: [a] -> [[a]] --zwraca tablicę tablic
inits' [] = [[]] --jeżeli tablica jest pusta, to zwracamy tablicę z jednym elementem, którym jest pusta tablica
inits' (x:xs) = [] : map (x:) (inits' xs) --dla każdego elementu w tablicy wywołujemy funkcję rekurencyjnie
--dla reszty tablicy, a następnie dodajemy ten element do każdego wyniku

--zadanie 19
tails' :: [a] -> [[a]] --zwraca tablicę tablic
tails' [] = [[]] --jeżeli tablica jest pusta, to zwracamy tablicę z jednym elementem, którym jest pusta tablica
tails' xs = tails' (tail xs) ++ [xs] --tail zwraca tablicę, bez pierwszego elementu.
--dla każdej tablicy wywołujemy rekurencyjnie funkcję, a jej ostatni element dodajemy na końcu tablicy (++)
-- tails' [1,2,3] = tails' [2,3] ++ [[1,2,3]] = tails' [3] ++ [[2,3]] ++ [[1,2,3]] = 
-- =tails' [] ++ [[3]] ++ [[2,3]] ++ [[1,2,3]] = [[]] ++ [[3]] ++ [[2,3]] ++ [[1,2,3]] = [[], [3], [2,3], [1,2,3]]

--zadanie 20
splits :: [a] -> [([a], [a])] --zwraca tablicę par tablic
splits xs = [(take n xs, drop n xs) | n <- [0..length xs]] --take n xs pobiera n pierwszych elementów, drop n pomija pierwsze n elementów

--zadanie 21
partition' :: (a -> Bool) -> [a] -> ([a], [a])
partition' _ [] = ([], []) --jeżeli lista wejściowa jest pusta, to zwracamy dwie puste tablice
partition' p (x:xs) = if p x then (x:l, r) --jeżeli element spełnia warunek, to trafia do pierwszej tablicy
                        else (l, x:r) --jeżeli element nie spełnia warunku, to trafia do drugiej tablicy
                        where (l, r) = partition' p xs --reszte tablicy dzielimy rekurencyjnie


--zadanie 26
isSorted :: Ord a => [a] -> Bool --Ord a oznacza, że elementy są porównywalne
isSorted [] = True
isSorted [_] = True
isSorted (x:y:xs)
    | x <= y = isSorted (y:xs)
    | otherwise = False

--zadanie 29
rev :: [a] -> [a]
rev [] = []
rev (x:xs) = rev xs ++ [x] 
--operator ++ ma w haskellu złożoność obliczeniową O(n), bo musi przejść przez całą lewą listę, żeby dokleić x
--Zatem złożoność naszego całego algorytmu to: 1 + 2 + 3 + ... + n-1 + n = O(n^2)

--zadanie 30
filter' p = concat . map box
    where box x = if p x then [x] else []
--funkcja box działa w następujący sposób: jeżeli x spełnia dany predykat to zostaje utworzona tablica z tym jednym elementem, 
--jeżel nie spełnia, to zostaje utworzona pusta tablica []. concat łączy wszystkie tablice w jedną, odrzucając puste tablice

--zadanie 31
--a)
takeWhile' :: (a -> Bool) -> [a] -> [a]
takeWhile' _ [] = []
takeWhile' p (x:xs)
    | p x = x : takeWhile' p xs
    | otherwise = []
--w przypadku gdy do takeWhile wrzucimy pustą tablicę, to zwróci pustą tablicę, w innym przypadku będzie zwracał elementy, dopóki nie natrafi
--na pierwszy element nie spełniający p

--b)
dropWhile' :: (a -> Bool) -> [a] -> [a]
dropWhile' _ [] = []
dropWhile' p (x:xs)
    | p x = dropWhile' p xs
    | otherwise = (x:xs)
--w przypadku gdy x spełnia warunek jest pomijany. gdy dropWhile' natrafi na pierwszy element nie spełiający warunku,
--zwraca go wraz z pozostałą częscią listy

--zadanie 34
-- :t sum: sum :: (Foldable t, Num a) => t a -> a | sumuje wszystkie liczby w strukturze
-- :t product: product :: (Foldable t, Num a) => t a -> a | mnoży wszystkie liczby w strukturze
-- :t all: all :: Foldable t => (a -> Bool) -> t a -> Bool | sprawdza czy wszystkie elementy w strukturze spełniają dany warunek
-- :t any: any :: Foldable t => (a -> Bool) -> t a -> Bool | sprawdza czy przynajmniej jeden element w strukturze spełnia dany warunek
--foldable oznacza, że struktura może być przekształcona w jedną wartość za pomocą folda

--zadanie 35
--let x = [1..100000]

--foldl (+) 0 xs -> zwraca sumę wszystkich elementów w tablicy
--foldr (+) 0 xs -> zwraca sumę wszystkich elementów w tablicy (liczy od końca)
--foldl1 (+) xs -> zwraca sumę wszystkich elementów w tablicy (liczy od początku, bez podawania wartości początkowej, działą dla list niepustych)
--foldr1 (+) xs -> zwraca sumę wszystkich elementów w tablicy (liczy od końca, bez podawania wartości początkowej, działą dla list niepustych)
--foldl' i foldr' są szybsze i nie używają lazy evaluation

--zadanie 36
reverse' :: [a] -> [a]
reverse' = foldl (flip(:)) [] --foldl przechodzi przez wszystkie elementy w tablicy i dodaje je do nowej tablicy, flip zamienia miejscami argumenty funkcji (:)

--zadanie 37
count_even' :: [Integer] -> Integer
count_even' = foldr (\x acc -> if even x then acc + 1 else acc) 0 --foldr przechodzi przez wszystkie elementy w tablicy i 
--zlicza parzyste liczby, acc to akumulator, który przechowuje liczbę parzystych liczb

--zadanie 42
approx :: Int -> Double
approx n = foldr (+) 0 [1 / fromIntegral (product [1..k]) | k <- [1..n]] --tworzymy tablicę liczb [1..n]
--dla każdej liczby k w tablicy obliczamy 1 / k!, i każdy następny element to suma poprzednich dodać ten 1/k!
--zaczynamy od 0, kończymy na n

--zadanie 47
subsetlist :: Eq a => [a] -> [a] -> Bool
subsetlist [] _ = True
subsetlist _ [] = False
subsetlist (x:xs) ys = elem x ys && subsetlist xs ys
--elementy muszą być porównywalne (klasa Eq)

--zadanie 48
mmap f = map (map f)
mmmap f = map (map (map f))

--zadanie 49
dotprod :: (Num a) => [a] -> [a] -> a
dotprod [] _ = 0
dotprod _ [] = 0
dotprod (x:xs) (y:ys) = x * y + dotprod xs ys
--klasa Num definiuje, że możemy na elementach wykonywać operacje arytmetyczne



--zadanie 67
data MyTree a = Empty | Leaf a | Node (MyTree a) a (MyTree a)
    deriving (Eq)

--1) 
instance Show a => Show (MyTree a) where
  show Empty = "Empty"
  show (Leaf x) = "Leaf " ++ show x
  show (Node left x right) =
   "Node (" ++ show left ++ ") " ++ show x ++ " (" ++ show right ++ ")"

--przykładowa deklaracja drzewa: let drzewo = Node (Leaf 1) 2 (Node Empty 3 (Leaf 4))

--2) 
--Funktor to typ, do którego można zastosować funkcję, nie zmieniając jego struktury
--formalnie, funktor to typ f, dla którego można zdefiniować funkcję:
-- fmap :: (a -> b) -> f a -> f b

instance Functor MyTree where
    fmap _ Empty = Empty
    fmap f (Leaf x) = Leaf (f x)
    fmap f (Node left x right) = Node (fmap f left ) (f x) (fmap f right)

--fmap (*5) drzewko

--3)
instance Foldable MyTree where
    foldr _ acc Empty = acc
    foldr f acc (Leaf x) = f x acc
    foldr f acc (Node left x right) = foldr f (f x (foldr f acc right)) left

--ewnetualna implementacja foldl:
--foldl _ acc Empty = acc
--foldl f acc (Leaf x) = f acc x
--foldl f acc (Node left x right) = foldl f (f (foldl f acc left) x) right

--foldr (+) 4 drzewko


--4)
height :: MyTree a -> Int
height Empty = 0
height (Leaf _) = 1
height (Node left _ right) = 1 + max (height left) (height right)

--5)
nbLeafs :: MyTree a -> Int 
nbLeafs Empty = 0
nbLeafs (Leaf _) = 1
nbLeafs (Node left _ right) = nbLeafs left + nbLeafs right

--6)
nbNodes :: MyTree a -> Int
nbNodes Empty = 0
nbNodes (Leaf _) = 0
nbNodes (Node left _ right) = 1 + nbNodes left + nbNodes right


--zadanie 68
--a)
--Ta naturalna transformacja jest w haskellu zaimplementowana za pomocą funkcji id (:i id)
-- id :: a -> a

--zadanie 73
choose :: Ord a => [a] -> [a]
choose (x:y:zs)
    | x > y = x : choose (y:zs)
    | otherwise = choose (y:zs)
choose _ = []   
--zapis (x:y:zs) oznacza, że tablica musi mieć co najmniej dwa elementy - x jest pierwszym, y drugim, a zs to reszta tablicy 


--zadanie 74

maxSum :: (Ord a, Num a) => [a] -> [a]
maxSum (x:xs)
    |x > sum xs = x : maxSum xs
    | otherwise = maxSum xs
maxSum [] = []

--sum :: (Num a) => [a] -> a
--sum [] = 0
--sum (x:xs) = x + sum xs

--Obecna implementacja funkcji maxSum ma złożoność obliczeniową O(n^2), spróbujmy ulepszyć ją do złożoności liniowej:
--na początku przygotujmy tablicę, która przechowuje kolejne zsumowane elementy tablicy wejściowej: np. dla tablicy [1, 3, 5, 7] przechowuje ona [1, 4, 9, 16]
--następnie możemy sprawdzać, czy dany element jest większy od sumy pozostałych elementów tablicy odejmując liczbę będącą na tym samym indeksie 
--w przygotowanej tablicy od ostatniego elementu. Jeżeli wynik różnicy będzie mniejszy niż ta liczba, to znacz, że jest ona większa od sumy pozostałych elementów


maxSumInN :: (Ord a, Num a) => [a] -> [a]
maxSumInN xs = go xs 0 totalSum []
  where
    totalSum = sum xs
    go [] _ _ acc = reverse acc
    go (y:ys) prefixSum remaining acc
        | y > (remaining - y) = go ys (prefixSum + y) (remaining - y) (y:acc)
        | otherwise           = go ys (prefixSum + y) (remaining - y) acc





--sums to tablica, która przechowuje kolejne zsumowane elementy tablicy wejściowej, scanl działa podobnie do foldl, ale zachowuje wszystkie pośrednie wyniki
--total to ostani element tablicy sums, czyli suma wszystkich elementów tablicy wejściowej
--zip xs (init sums) łączy oryginalną tablicę z listą sums, ale bez ostatniego elementu, żeby dopasowanie miało tę samą długość
--przykład działania zip: zip [1, 3, 5, 7] [0, 1, 4, 9] = [(1,0), (3,1), (5,4), (7,9)]
--x > total - s oznacza, że dany element jest większy od sumy pozostałych elementów tablicy


--test
prop_sameResults :: [Integer] -> Bool
prop_sameResults xs = maxSum xs == maxSumInN xs
--do odpalenia: import Test.QuickCheck, quickCheck prop_sameResults
--as