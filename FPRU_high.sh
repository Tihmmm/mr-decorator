#!/bin/bash

FPRUtility -information -search -query "[fortify priority order]: high suppressed:false" -listIssues -project "$1/current.fpr" -outputFormat CSV > "$1/high.csv"
tail -n 2 "$1/high.csv" | cut -d '"' -f 2 | head -n 1 > "$1/high_count.txt"