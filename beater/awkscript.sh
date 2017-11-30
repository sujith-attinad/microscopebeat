#!/bin/bash
# awk script to parse nodetool cfstats output

nodetool $1 $2 | awk '
                    FNR == 2 { print $3 }
                    FNR == 4 { print $3 } '
