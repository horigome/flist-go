# flist go

2017-07-05 M.Horigome

## Overview
This is, to get the file list, is a small tool to be output to the CSV.


## How to make

    > cd flist-go
    > make

or

    > cd flist-go
    > go run main.go


## How to usage

    > cd _build
    > flist

#### -no:

No output CSV file

#### -nd:

No display datail

#### -m "pattern"

Match filename pattern (regurar expression)  

#### -md "pattern"

Match directory name pattern (regurar expression)  

#### -s "pattern"

Skip filename pattern (regurar expression)  

#### -sd "pattern"

Skip directory name pattern (regurar expression)  

#### -f "csv filename"

Specify csv file name

#### -sjij

use SJIS encoding

#### -version

Show version no.

## License

MIT.
