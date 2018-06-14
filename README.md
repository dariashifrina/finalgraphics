# Welcome to Gilvir and Dasha's Graphics Engine!

The following is a list of our additional **features**:


## Multiple Light Sources

Add additional light sources, in accordance to MDL spec;

```
light r g b x y z
```
``

You can also change the ambient light source, again in accordance to the MDL spec.

```
ambient r g b
```

## Meshes
In progress:
Functionality, can interpret "v" and "f" commands from Wavefront .obj format properly (ie the pear and teapot), however cannot do texture mappings or interpret vertex normals.

One problem we had was that OBJ files are 1-indexed and our function assumed 0-indexing. Please don't look at how we fixed it. It's kind of disgusting.
## Written (Mostly) in Go!

Our engine is a Go-Python Hybrid. Running `make` runs a the Go `parse` function, which in turn execvps the Python parsing code. Python 2.7 is used with yacc and flex implementations to parse the MDL files. Once the MDL files are parsed, we also deal with animation steps primarily in Python. The `mdl.py` file was slightly modified to ensure that any args that needed to be added to go were added and passed. A custom scripting language is used in place of serialization between Python and Go, and
the standard out of Python is piped to `parse.go`. The formatting of this language is similar to MDL, except it returns to having the command on a seperate line from it's argument. 

For animation, we simply have the Python code print out the commands with the knob adjusted values. For example, if you have an MDL file with a single `move` command accomponied by a knob (ie `move 100 0 0 knob`), then for every frame, Python will print out the `sphere` command with the {100} value modified based off the current knob value, then it pass a `save` command to Go that saves the current frame to a PNG. On the first pass, for example, it might be `move 0 0 0`, then `move 20 0 0`, then `move 40 0 0` etc...

tl;dr -- the Python code includes the parser and determines the order of operations, the Go code runs the Python code and actually executes the operations passed to it.
### Gotta Go fast
Writing the engine in Go provides a lot of benefits to us as the programmer and allows our engine to be one of the fastest MDL gif creating interpreters in the class, faster than most C implementations even. Similar to languages like Rust and Java, Go provides speeds comprable to C, with at least some of the ease of use found in Python (for those used to statically typed languages, they might even find Go to be easier than Python). Additionally, the Go language was built from the ground up for concurrency, something that cannot really be said for any of the previously mentioned languages other than Rust. 

Using Go has enabled us, for example, to use coroutines for our Matrix multiplication and polygon drawing, with very few complications and only about ~6 extra lines of code for each. This means that performance scales with the amount of threads available, something which can be seen if one was to run `htop` or `top` whilst generating an animation using this engine. 

Go's syntactic sugar and use of slices (similar to Python lists) also made it really easy to use column major arrays instead of 2d arrays for our matrices (something that works out really well since we will seldom be appending rows, only columns). This makes accessing values much faster, since you do not need to 2 dereferences, only 1 + some modular arithmetic.

Obviously, the engine isn't fast. That wasn't the point of this course. However, using Go has allowed us to easily make our engine much faster while gaining first hand experience with useful tools for concurrency, such as mutex locks and semaphores. It's also allowed us to see an example of how something like this might scale with dedicated GPU's. 

Go is also a language with a growing community and is one that is definitely worth learning (Dasha learned the language while doing this project, and Gilvir learned it at a slower but still reasonable pace when initially starting this engine at the beginning of the semester). It's a ease of use and documentation, and inclusion of modern features, makes it a very powerful tool for all developers.
