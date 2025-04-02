
--or
or' xs = foldl (||) False xs

--and
and' xs = foldl (&&) True xs

--concat
concat' xxs = foldl (++) [] xxs

--concatMap
concatMap' f = foldr ((++) . f) []
