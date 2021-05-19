# xk6-output-template
Is template for k6 output extensions 

You should make a repo from this template and go through the code and replace everywhere where it says `template` in order to use it
TODO: more instructions and comment inline
</div>

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

1. Build with `xk6`:

```bash
xk6 build --with github.com/k6io/xk6-output-template
```

This will result in a `k6` binary in the current directory.

2. Run with the just build `k6:

```bash
./k6 run -o xk6-template <script.js>
```

