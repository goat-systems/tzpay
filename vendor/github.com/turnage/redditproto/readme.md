# redditproto

Protocol Buffer definitions for Reddit's JSON api types. For
more information about what these messages represent, see
[Reddit's docs](https://github.com/reddit/reddit/wiki/JSON).

To generate the code for the protobuffers for your language, you can most likely
do the following:

    [your package manager] install protobuf-compiler
    protoc --[your lang code]_out=. *.proto

Common examples:

    protoc --cpp_out=. *.proto
    protoc --java_out=. *.proto
    protoc --python_out=. *.proto

See ````protoc --help```` for further guidance.

# Update Policy

No breaking changes will be made to redditproto messages as of version 1.0. Any
message changes will be noted here in the readme.

* 0.9.0 -> August 25th, 2015
* 1.0.0 -> October 21st, 2015
* 2.0.0 -> Octover 29th, 2015 (No message changes; improve Go utilities).
* 2.1.0 -> December 22nd, 2015 (Added LinkSet message)

# Gophers

The golang generated code is included to make this package go gettable. There
are also some utility functions included for parsing Reddit's JSON responses. If
you're writing a golang utility for Reddit I'd be happy to add things for you!
