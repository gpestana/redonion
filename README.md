# Redonion

[![Build Status](https://travis-ci.org/gpestana/redonion.svg?branch=master)](https://travis-ci.org/gpestana/redonion)

Deep web scanning done right

![redonion logo](https://raw.githubusercontent.com/gpestana/redonion/master/redonion.png)

### How to use

To scan a list of onion sites and output it to a file:

`./redonion -list=./sample/short_list.txt > output.json`


#### Fetchers

Fetchers fetch data from a list of onion sites.

- **fetcher.go**: (`/fetchers/fetch.go`) will get all the HTML data from a 
certain onion site.

#### Processors

Processors make transformations on the input data. The input data may come from
fetchers or from other processors.

- **text**: (`/processors/text.go`) pass-through processor which will pass the
output content from input channel to output channel.

- **image**: (`/processors/image.go`) fetches all website images and inspects
various types of metadata

### Contribute

Contributions for new fetchers, processors and outputs welcome.

gpestana © MIT
