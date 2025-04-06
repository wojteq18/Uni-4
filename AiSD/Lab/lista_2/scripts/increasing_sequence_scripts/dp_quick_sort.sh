#!/bin/bash

GENERATOR="/home/wojteq18/Uni/AiSD/Lab/lista_2/generators/increasing_sequence/target/release/increasing_sequence"
SORT="/home/wojteq18/Uni/AiSD/Lab/lista_2/dual_pivot_quick_sort/target/release/dual_pivot_quick_sort"
OUTPUT="dp_quick_sort.txt"

echo "n,c,s,trial" > "$OUTPUT"

for trial in {1..10}; do
    for (( n=1000; n<=50000; n+=1000 )); do
        result=$($GENERATOR $n | $SORT)
        c=$(echo "$result" | grep -oP 'c=\K[0-9]+')
        s=$(echo "$result" | grep -oP 's=\K[0-9]+')
        echo "$n,$c,$s,$trial" >> "$OUTPUT"
    done
    echo "Trial $trial done"
done
