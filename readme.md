# go-llm-reAct

An implementation of the [reAct](https://arxiv.org/pdf/2210.03629.pdf) paradigm in Golang. Heavily inspired by Simon Willison's [implementation in Python](https://til.simonwillison.net/llms/python-react-pattern).

```
Thought: I think I remember the neighboring countries of England, but let me double check to be sure
Action: wikipedia: England
Observation: <span class="searchmatch">England</span> is a country that is part of the United Kingdom. It shares land borders with Wales to its west and Scotland to its north. The Irish Sea lies northwestObservation: England shares borders with Wales to the west and Scotland to the north.
Answer: England shares borders with Wales and Scotland.
```