# Extending Gorilla Mux Router Tutorial #

I found this to be a nice pattern for centralizing your handler declarations and providing context to your handlers.  This allows you to supply resources to your Handlers; e.g. DB connection pool, Redis connection, etc...

There's also a few examples of how to wrap your handlers in middlewares to provide authentication or GZIP compression.

The master branch is the complete tutorial.  There are some branches which have incremental changes along the way.