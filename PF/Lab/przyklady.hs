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
