--funkcje collatz'a

coll :: Integer -> Integer
coll 1 = 1
coll n
    | even n    = coll (div n 2)
    | otherwise = coll (3 * n + 1)

collatz :: (Int, Int) -> (Int, Int)
collatz (n, s)
    | n == 1    = (n, s)
    | even n    = collatz (n `div` 2, s + 1)
    | otherwise = collatz (3 * n + 1, s + 1)

