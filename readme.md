# Configurable SVG calendar generator

Build: 

    go build

Usage:

    paracal --help

Example:
 
    paracal -s ./configs/all.yaml -y 2021 -m 12 -o cal.svg -b '#ffdab9'

Convert to printable PNG:

    inkscape cal.svg --export-dpi=300 -o cal.png 
