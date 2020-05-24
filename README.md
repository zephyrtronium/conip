# conip

conip produces a minimal-size string containing every IPv4 address.

With a recent version of Go installed and running in module mode, you can run this by doing `go run github.com/zephyrtronium/conip`. If you don't want it to print a many-gigabytes string to stdout, add `-help` to see options for output type, location, and buffer size. While running, the program prints progress updates to stderr indicating the most significant byte of its working memory; the time between such updates shortens quadratically.

The particular sequence printed is a de Bruijn sequence `B(256, 4)` beginning
with four zeros. With the default text output, the alphabet is the set
`{"0", "1", "2", ..., "255"}`. A `.` or newline character separates each
sequence term. The output is around 14.2 GiB.

With binary output, the alphabet is the set `{0, 1, 2, ..., 255}`, and each
term is written as a single byte with no separating characters. The output
is exactly 4 GiB plus three bytes.
