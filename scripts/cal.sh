#!/bin/bash
set -e

# Usage:
# ./cal.sh [noimg] [01-12]
#
# Expecting:
# ./cfg/*.yaml
# ./img/[01-12].jpg
#
# Result:
# ./svg/*.svg
# ./png/*.png

YEAR=25
DPI=300
WIDTH=3555
HEIGH=2527

MON=$1
NOIMG=

if [[ $1 == "noimg" ]]; then
    MON=$2
    NOIMG=1
fi

from=1
to=12
if [[ $MON ]]; then from=$MON; to=$MON; fi

mkdir -p ./svg
mkdir -p ./png
for i in $(seq $from $to); do
    i=$(printf "%02d" $i)
    if [[ -z $NOIMG ]]; then
        IMG="${i}_tmp_.jpg"
        back="--back $IMG"
        cp ./img/$i.jpg ./svg/$IMG
    fi

    echo "Generate $i"
    paracal -s ./cfg/$i.yaml -y $YEAR -m $i -o "./svg/$i.svg" $back
    # Use inkscape or rsvg-convert (it faster).
    # inkscape $1.svg -d $DPI -w $WIDTH -h $HEIGH -o $1.png && rm $1.svg
    rsvg-convert ./svg/$i.svg -d $DPI -p $DPI -w $WIDTH -h $HEIGH -o ./png/$i.png

    if [[ -z $NOIMG ]]; then rm ./svg/$IMG; fi
done

if [[ $MON && $3 == "open" ]]; then open ./png/$MON.png; fi
