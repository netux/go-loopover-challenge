# Loopover Challenge in Go

This is a port of spdskatr's loopover_oo.py script.

## How to play
```
 18 22 23  4  3
  7 24  1 16  8
 17 21 20 12 25
  2  5  9  6 15
 10 19 13 11 14
Move: $
```

You are presented with a grid like the one above, your goal is to get all the numbers ordered from lower to higher, starting from the top left and ending on the bottom right, like this:

```
  1  2  3  4  5
  6  7  8  9 10
 11 12 13 14 15
 16 17 18 19 20
 21 22 23 24 25
Solved!
```

But here is the catch: you can only move a row or a column of numbers, for example:

```bash
# Moving the third row from the top by 1
  1  2  3  4  5        1  2  3  4  5
  6  7  8  9 10        6  7  8  9 10
[11 12 13 14 15] --> [15 11 12 13 14]
 16 17 18 19 20       16 17 18 19 20
 21 22 23 24 25       21 22 23 24 25

# Moving the first column from the left by 2
| 1| 2  3  4  5      |16| 2  3  4  5
| 6| 7  8  9 10      |21| 7  8  9 10
|11|12 13 14 15  --> | 1|12 13 14 15
|16|17 18 19 20      | 6|17 18 19 20
|21|22 23 24 25      |11|22 23 24 25
```

## Possible moves

- `shuffle`: shuffles the board. You probably want to do this before anything else.
- `reset`: resets the board to its original state and sets the moves done back to 0.
- programmer's notation: allows to modify the board. See the "Programmer's Notation" section below to learn more about it.

## Programmer's Notation
```js
// TODO(netux): add programer's notation syntax
```

You can see the Wirth syntax notation of the Programmer's Notation in line 21 of loopover.go
