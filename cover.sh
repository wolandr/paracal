#!/bin/bash
#
# Geneate image 3508x2480.
# 
# Expecting:
#  ./cover/00.jpg 1724 x 1210 (center image)
#  ./cover/[01 - 12].jpg 852 x 595   (12 icon image)
#  ./logo/logo_stroke.png
#
# Usage:
# ./cover.sh cover_2024.png
#

# Parameters
W=3508 # Image width 
H=2480 # Image heigh
T=20   # Frame width
EXT=47 # Extra canvas border

OFFSET=$(( EXT/2 ))
iW=$(( (W-T*5)/4 )) # Icon width
iH=$(( (H-T*5)/4 )) # icon heigh

echo "Expected icon size: ${iW} x ${iH}"
echo "Actual icons  size: $(identify -format '%W x %H' ./cover/01.jpg)"

echo "Expected center sz: $(( iW*2+T )) x $(( iH*2+T ))"
echo "Actual center size: $(identify -format '%W x %H' ./cover/00.jpg)"

# Positions for images
X0=$(( T + OFFSET ))
X1=$(( X0 + iW + T ))
X2=$(( X1 + iW + T ))
X3=$(( X2 + iW + T ))
Y0=$(( T + OFFSET )) 
Y1=$(( Y0 + iH + T )) 
Y2=$(( Y1 + iH + T )) 
Y3=$(( Y2 + iH + T ))

function img() {
    IMG+="image over $1,$2 ${iW},${iH} cover/$3.jpg  "
}

img $X0 $Y0 01  #01
img $X1 $Y0 02  #02 
img $X2 $Y0 03  #03
img $X3 $Y0 04  #04

img $X0 $Y1 05  #05 
img $X3 $Y1 06  #06

img $X0 $Y2 07  #07
img $X3 $Y2 08  #08

img $X0 $Y3 09  #09
img $X1 $Y3 10  #10
img $X2 $Y3 11  #11
img $X3 $Y3 12  #12

# Draw center image
IMG+="image over $X1,$Y1 $((iW*2+T)),$((iH*2+T)) cover/00.jpg  "

LOGO+="image over 0,-250 1050,432 logo/logo_stroke.png  "
TEXT+="text 0,250 \"2024\""

xW=$((W+EXT))
xH=$((H+EXT))
echo "Drawing ${xW} x ${xH} ..."

export LC_CTYPE=
convert -size ${xW}x${xH} xc:none \
    -fill "#ffffff" -stroke none -draw "rectangle 0,0,${xW},${xH}" \
    -draw "${IMG}" \
    -gravity center -draw "${LOGO}" \
    -pointsize 300 -font Lobster-regular -fill "#0054ff" -stroke white -strokewidth 1 -gravity center -draw "${TEXT}" \
    -colorspace sRGB PNG32:"$1"
