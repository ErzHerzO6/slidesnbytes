# Welcome to Slides
A terminal based presentation tool

## Everything is markdown
In fact this entire presentation is a markdown file.

Press `n` to go to the next slide.

---

# Display Code

```go
package main

import "fmt"

func main() {
  // You can show code in slides
  // Press ctrl+e to execute this code directly in slides
  fmt.Println("Tada!")
}
```

---

# h1

You can use everything in markdown!
* Like bulleted list
* You know the deal

1. Numbered lists too

## h2

| Tables | Too    |
| ------ | ------ |
| Even   | Tables |


### h3

#### h4
##### h5
###### h6

---


# Progressive Reveal

Text can appear step by step using `<!-- #break -->`:

- This point is visible immediately

<!-- #break -->

- This appears after pressing next

<!-- #break -->

- And this is the final reveal

Great for building up ideas incrementally!

---

# Graphs

```
digraph {
    rankdir = LR;
    a -> b;
    b -> c;
}
```
```
┌───┐     ┌───┐     ┌───┐
│ a │ ──▶ │ b │ ──▶ │ c │
└───┘     └───┘     └───┘
```
---

All you need to do is separate slides with triple dashes
`---` on a separate line, like so:

```markdown
# Slide 1
Some stuff

--- 

# Slide 2
Some other stuff
```
