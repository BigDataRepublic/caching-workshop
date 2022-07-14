#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

while IFS="," read -r rec_column1 rec_column2 rec_remaining
do
  curl -s -H "Content-type: application/json" -d "{\"username\": \"$rec_column1\", \"points\": $rec_column2}" localhost:8080/points
  echo ""
done < <(tail -n +2 $SCRIPT_DIR/users.csv)