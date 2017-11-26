# Redonion

i**Deep web scanning done right**

![redonion logo](https://raw.githubusercontent.com/gpestana/redonion/master/redonion.png)


Redonion is a plugable scanner for the tor network. You can write your own data i
fetchers (`/fetchers`), processors (`/processors`) and poutput components 
(`/outputs`) or use the pre built components. All components run on goroutines
and results of each components are piped using go channels. This architecture
allows for fast scanning and easy development of new features.

## How to use

To scan a list of onion sites and output it to a file:

`./redonion -list=./sample/short_list.txt > output.json`


### Fetchers

Fetchers fetch data from a list of onion sites.

- **fetcher.go**: (`/fetchers/fetch.go`) will get all the HTML data from a 
certain onion site.

### Processors

Processors make transformations on the input data. The input data may come from
fetchers or from other processors.

- **text**: (`/processors/text.go`) pass-through processor which will pass the
output content from input channel to output channel.

### Outputs 

- **stdout**: (`/outputs/stdout.go`) writes result of the pipeline to standart 
output


## Contribute

Contributions for new fetchers, processors and outputs welcome.

gpestana © MIT or whatever
