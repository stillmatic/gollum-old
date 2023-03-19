# GoLLuM

An implementation of the [reAct](https://arxiv.org/pdf/2210.03629.pdf) paradigm in Golang. Heavily inspired by Simon Willison's [implementation in Python](https://til.simonwillison.net/llms/python-react-pattern).

```
Thought: I need to recall the geography of England and its neighboring countries.
Action: wikipedia: England
Observation: <span class="searchmatch">England</span> is a country that is part of the United Kingdom. It shares land borders with Wales to its west and Scotland to its north. The Irish Sea lies northwest
Answer: England shares land borders with Wales to its west and Scotland to its north.
```

