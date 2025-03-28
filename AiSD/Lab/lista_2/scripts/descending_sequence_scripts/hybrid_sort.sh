#!/bin/bash

GENERATOR="/home/wojteq18/Uni/AiSD/Lab/lista_2/generators/descending_sequence/target/release/descending_sequence"
SORT="/home/wojteq18/Uni/AiSD/Lab/lista_2/hybrid_sort/target/release/hybrid_sort"
OUTPUT="hybrid_sort.txt"

echo "n,c,s" > "$OUTPUT"

for n in {8..32}; do
    result=$($GENERATOR $n | $SORT)
    c=$(echo "$result" | grep -oP 'c=\K[0-9]+')
    s=$(echo "$result" | grep -oP 's=\K[0-9]+')
    echo "$n,$c,$s" >> "$OUTPUT"
done
