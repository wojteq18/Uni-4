#!/bin/bash

GENERATOR="/home/wojteq18/Uni/AiSD/Lab/lista_3/random_generator/target/release/random_generator"
SELECT="/home/wojteq18/Uni/AiSD/Lab/lista_3/random_select/random_select_bin"
OUTPUT="random_select.txt"

echo "n,c,s,trial" > "$OUTPUT"

for trial in {1..5}; do
    for (( n=100; n<=50000; n+=100 )); do
        result=$($GENERATOR $n | $SELECT)
        c=$(echo "$result" | grep -oP 'c\s*=\s*\K[0-9]+')
        s=$(echo "$result" | grep -oP 's\s*=\s*\K[0-9]+')
        echo "$n,$c,$s,$trial" >> "$OUTPUT"
    done
    echo "Trial $trial done"
done
