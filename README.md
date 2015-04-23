Japanese to English dictionary API for Golang
---------------------------------------------

This is a basic dictionary storage and lookup for Golang for the electronic
dictionary compiled by Jim Breen and The Electronic Dictionary Research and
Development Group.

### Installation
#### Download dictionary file
First download the JMdict dictionary file. This project only supports the JMdict
XML dictionary format, not the original EDICT format. The easiest way to download
this is to use rsync:

    rsync -v -z --progress ftp.monash.edu.au::nihongo/JMdict JMdict.xml

The `-z` option tells rsync to compress over the wire which will hugely speed
up transfer.

#### Set up your database
Import the postgre.sql file into your database to set up the tables needed.

#### Import the dictionary
Now your database should be ready and you have the dictionary file.  First
install this package using `go install`.  Run the following command, replacing
the appropriate parts of the DB connection string:

    go-jedict -db=postgres://username:password@hostname/database -import JMdict.xml

Make a cup of miso soup, the import will take some time.

### Usage
#### Look up words via the command line

To make sure everything works, look up a word via the command line:

    go-jedict -db=postgres://username:password@hostname/database -kanji 辛い

You should see:

    辛い
    からい
    ----
    spicy, chilly (hot)

#### Using the library API

