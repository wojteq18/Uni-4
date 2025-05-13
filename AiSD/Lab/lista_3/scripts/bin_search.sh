#!/bin/bash

GENERATOR="/home/wojteq18/Uni/AiSD/Lab/lista_3/random_generator/target/release/random_generator"
SELECT="/home/wojteq18/Uni/AiSD/Lab/lista_3/bin_search/bin_search_bin"
OUTPUT="bin_search.txt"

echo "n,c,trial" > "$OUTPUT"

for trial in {1..15}; do
    for (( n=100; n<=100000; n+=100 )); do
        result=$($GENERATOR $n | $SELECT)
        c=$(echo "$result" | grep -oP 'c\s*=\s*\K[0-9]+')
        echo "$n,$c,$trial" >> "$OUTPUT"
    done
    echo "Trial $trial done"
done
