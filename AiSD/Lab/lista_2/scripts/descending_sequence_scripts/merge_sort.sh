#!/bin/bash

GENERATOR="/home/wojteq18/Uni/AiSD/Lab/lista_2/generators/descending_sequence/target/release/descending_sequence"
SORT="/home/wojteq18/Uni/AiSD/Lab/lista_2/merge_sort/target/release/merge_sort"
OUTPUT="merge_sort.txt"

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
