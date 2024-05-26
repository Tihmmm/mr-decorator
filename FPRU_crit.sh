#!/bin/bash

FPRUtility -information -search -query "[fortify priority order]: critical suppressed:false" -listIssues -project "$1/current.fpr" -outputFormat CSV > "$1/critical.csv"
tail -n 2 "$1/critical.csv" | cut -d '"' -f 2 | head -n 1 > "$1/critical_count.txt"