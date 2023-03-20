# GoLLuM

An implementation of the [reAct](https://arxiv.org/pdf/2210.03629.pdf) paradigm in Golang. Heavily inspired by Simon Willison's [implementation in Python](https://til.simonwillison.net/llms/python-react-pattern).

```
Thought: I think I remember the countries that England shares borders with, but I should double-check to be sure.
Action: wikipedia: England
Observation: England is a country that is part of the United Kingdom. It shares land borders with Wales to its west and Scotland to its north. The Irish Sea lies northwest
Thought: Ok, so England shares borders with Wales and Scotland, and has the Irish Sea to the northwest. I have the information I need.
Answer: England shares borders with Wales and Scotland, and has the Irish Sea to the northwest.
```

Error handling is supported:
```
Action: search map of England on Google
Oops: tool not found, available tools are: calculator, wikipedia
Action: wikipedia: England
Observation: England is a country that is part of the United Kingdom. It shares land borders with Wales to its west and Scotland to its north. The Irish Sea lies northwest
```