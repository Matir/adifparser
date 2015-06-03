# ADIFParser #

`ADIFParser` is a Go library for reading and writing [Amateur Data Interchange
Format](http://www.adif.org/) files.  It wraps `io.Reader` and `io.Writer`
interfaces to handle I/O and attempts to handle the irregularities of parsing
files as much as possible.

### Shortcomings ###

Currently, no validation of the content of fields is done.  Also, no
transformations are done, the fields are all handled as strings.

### License ###

This library is released under a 2-clause BSD license.  See COPYING for the
license.

### Bugs ###

Please report bugs on GitHub at https://github.com/Matir/adifparser, and please
include an ADIF test case to demonstrate your bug.
