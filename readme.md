# Configurable SVG calendar generator

Build: 

    go build

Usage:

    paracal --help

Example:
 
    paracal -s ./configs/all.yaml -y 2021 -m 12 -o cal.svg -b '#ffdab9'

Convert to printable PNG:

    rsvg-convert cal.svg -d 300 -p 300 -w 3555 -h 2527 -o cal.png

OR you can use inkscape, but it slower:

    inkscape cal.svg --export-dpi=300 -w 3555 -h 2527 -o cal.png 

