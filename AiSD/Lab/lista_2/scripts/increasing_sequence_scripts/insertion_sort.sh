#!/bin/bash

GENERATOR="/home/wojteq18/Uni/AiSD/Lab/lista_2/generators/increasing_sequence/target/release/increasing_sequence"
SORT="/home/wojteq18/Uni/AiSD/Lab/lista_2/insertion_sort/target/release/insertion_sort"
OUTPUT="insertion_sort.txt"

echo "n,c,s" > "$OUTPUT"

for n in {8..32}; do
    result=$($GENERATOR $n | $SORT)
    c=$(echo "$result" | grep -oP 'c=\K[0-9]+')
    s=$(echo "$result" | grep -oP 's=\K[0-9]+')
    echo "$n,$c,$s" >> "$OUTPUT"
done
