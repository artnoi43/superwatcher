# Components

> Note: we also provide [`initsuperwatcher`](../initsuperwatcher/) so that
> users can initialize all components in one function call - eliminating
> possible confusion caused by the instability of the internals.

This public package provides functions for initializing superwatcher components.

The internals of superwatcher is not stable yet, so we provide a separate and more stable package
for creating new instances of the core components.

It is recommended that users only use code in these sub-packages only if they want to initialize
such components separately. Otherwise, use [`initsuperwatcher.New`](../initsuperwatcher/initsuperwatcher.go).
