# errctx

Package errctx facilitates storing key-value data into contexts and then wrapping error values with that data so top-level error handlers have access to the data from the entire call chain.

You can call `WithMeta` as many times as you like during a chain of function calls to decorate that call chain with metadata such as user IDs, request IDs and other business domain information. Then, when an error occurs, you wrap the error with a contextual error which contains the key-value data that was stored in the `context.Context` value. Then when your error is handled, you can easily extract this metadata for logging or error message purposes.

Inspired by [fctx](https://github.com/Southclaws/fault/blob/bffb7f66aa9a5b6231b31594b23749eace65b9ac/fctx/fctx.go) package, but with modifications to support `any` value type, and handle odd number of key-value pairs like how `slog` package handle it.

TODO: add linter to check odd number of key and value pairs and key should be of type `string`.
