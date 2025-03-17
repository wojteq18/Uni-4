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