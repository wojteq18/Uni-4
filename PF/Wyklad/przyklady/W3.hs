--SORTOWANIA

--w ghci :set +s podaje czas wykonania funkcji

--quick_sort
qSort :: Ord a => [a] -> [a] --ord to klasa typów, któ©e można porównywać między sobą
qSort [] = []
qSort (x:xs) = qSort [y | y <- xs, y < x] ++ [x] ++ qSort [y | y <- xs, y >= x]


--partition
partition :: (a -> Bool) -> [a] -> ([a], [a])
partition _ [] = ([], [])
partition p (x:xs) = if p x then (x:l, r)
                        else (l, x:r)
                        where (l, r) = partition p xs
--funkcja partition przyjmuje 2 argumenty i dzieli podaną tablicę na 2 podtablice, jedna, które spełnia podany warunek i druga, która nie spełnia.
--przykład:   in: partition even [1,2,3,4,5,6]  out: ([2,4,6],[1,3,5])  


--quick_sort bez list comprehension (bo to mocno spowalnia działanie funkcji)
qSort' :: Ord a => [a] -> [a]
qSort' [] = []
qSort' (x:xs) = qSort' l ++ [x] ++ qSort' r
                where (l, r) = partition (\t -> t < x)xs --dzielimy na dwie podtablice, jedna zawiera elementy mniejsze od pivota, druga większe


--insertion_sort
iSort :: Ord a => [a] -> [a]
iSort [] = []
iSort (x:xs) = l ++ [x] ++ r
                where sxs = iSort xs
                      (l,r) = partition (<x) sxs   


--zip
zip' :: [a] -> [b] -> [(a,b)] --przyjmuje 2 tablice i zwraca tablicę par elementów o tych samych indeksach (zip' [1,2] [3,4] = [(1,3),(2,4)])
zip' [] _ = []
zip' _ [] = []
zip' (x:xs) (y:ys) = (x, y): zip' xs ys


--add
add :: Num a => [a] -> a --sumuje elementy tablicy
add [] = 0
add (x:xs) = x + add xs


--myfoldr
myfoldr :: (t1 -> t2 -> t2) -> t2 -> [t1] -> t2
myfoldr op e [] = e
myfoldr op e (x:xs) = op x (myfoldr op e xs)
--myfoldr to rekurencyjna funkcja, która składa tablice od prawej strony; przyjmuje 3 argumenty: funkcję i 2 tablice

--isPrime
isPrime :: Int -> Bool
isPrime k = length [ x | x <- [2..k], k `mod` x == 0] == 1

