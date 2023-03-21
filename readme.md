# GoLLuM

An implementation of the [ReAct](https://arxiv.org/pdf/2210.03629.pdf) paradigm in Golang. Heavily inspired by Simon Willison's [implementation in Python](https://til.simonwillison.net/llms/python-react-pattern).

Easily empower your LLM to use tools:

```
Thought: I think I remember the countries that England shares borders with, but I should double-check to be sure.
Action: wikipedia: England
Observation: England is a country that is part of the United Kingdom. It shares land borders with Wales to its west and Scotland to its north. The Irish Sea lies northwest
Thought: Ok, so England shares borders with Wales and Scotland, and has the Irish Sea to the northwest. I have the information I need.
Answer: England shares borders with Wales and Scotland, and has the Irish Sea to the northwest.
```

And to self recover from errors:

```
Question: what is 2 + 2?
Thought: This is a simple calculation that I can solve using the calculate action.
Action: calculate: 2 + 2
Oops: tool not found, available tools are: calculator, wikipedia
Sorry about that. Let me try using the calculator action instead.
Action: calculator: 2 + 2
Observation: 4
Answer: 2 + 2 is equal to 4.
```

We strive to use pure Go, allowing for cross-platform builds and portability. 

# Supported Tools

Implementations exist for:

1. `calc`: passes calculations to [`expr`](https://github.com/antonmedv/expr). These are sandboxed and safe to run, unlike `eval` in Python.
2. `wikipedia`: returns a snippet from the first Wikipedia search result for a given item.

Some tools require 'state' and are a bit more complicated to use.

1. `sql`: this is an interface in front of Go's [`database/sql`](https://pkg.go.dev/database/sql) interface. It requires a [`db`](https://pkg.go.dev/database/sql#DB) object to be instantiated. See the tests for an example - the overall result should work with [any Go package](https://github.com/golang/go/wiki/SQLDrivers) that implements the `driver` interface. Per [Rajkumar et al](https://arxiv.org/abs/2204.00498) (2023) - we provide the CTAS schema and sample data as input.

We also have some special built-in 'tools.' 

1. `help`: returns the description for the requested tool.