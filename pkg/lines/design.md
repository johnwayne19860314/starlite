# Design Principles
### Ruby on Rails
* Don't Repeat Yourself (DRY)
* Convention over Configuration
### Go
* Simplicity, safety, and readability are paramount.
* Striving for orthogonality in design.
* Minimal: One way to write a piece of code.
* It's about expressing algorithms, not the type system.
* Collective unconscious history of programming languages.
* Things of interest should be easy; even if that means not everything is possible.

# Func
### CodeGen: Pen
no one want to write boilerplate code by hand over and over, 
we should make our valuable attention always on what's most important.
`Pen` will help you to keep definitions synced with generated code,
and like the idea behind `Spring Boot`, try to generate as less code as possible.

for client end side, it's also good to have a definition of api,
because we can use this to generate client library to call api like calling functions.

### IoC
LINES will generate and handle most aspect of the action of app,
it has to be a framework, and will have a lot of places can be configured.
`Instanced` should be able to be `injected` as wish, default by framework, and configurable by each app.

### Features
* Feature can be configured by yaml config or env, but it always have a default value in code
* Feature can be enabled and disabled in codegen config
* Feature is a specific function in environment, can be got by feature name from app context

# Version
lib should have stable version,
released and tagged by version branch