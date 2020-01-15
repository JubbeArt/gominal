# Gominal

Terminal like UI framework. 
* Grid based
* Single binary executable
* Communication done with JSON over stdin / stdout
* A box in the grid can contain either a character or some image data

## Requests

Sent to gominal on stdin

### char - draw a character to screen

```
{
    "type": "char"
    "char": string
    "col":  int
    "row":  int
    "style": "normal" | "bold" (optional, default "normal")
    "color": (optional, default white)
    {
        "r": int (0 to 255)
        "g": int
        "b": int
    } 
    "background": (optinal, default black)
    {
        "r": int
        "g": int
        "b": int
    }
}
```

**Example**

```json
{
    "type": "char",
    "char": "ö",
    "col": 5,
    "row": 10,
    "background": {
        "r": 255
    }
}
```


### image - draw image to screen
Will draw the image starting from col & row sent with the request. Gominal will use as many columns and rows as
is needed to draw the full image.

```
{
    "type": "image"
    "image": string, base64 encoded jpg or png image
    "col": int
    "row": int
}
```

**Example**
```json
{
    "type": "image",
    "image": "data...",
    "col": 5,
    "row": 5
}
```

### clear - clears screen from previous draw calls
**Example**

```json
{
    "type": "clear"
}
```

### title - set title of window
**Example**

```json
{
    "type": "title",
    "title": "The Gominal"
}
```

### close - closes window
**Example**

```json
{
    "type": "close"
}
```


## Events

Sent from gominal to stdout

### key - key press from window
Will run each time a key is pressed on the keyboard.

```
{
    "type":  "key"
    "key":   string
    "ctrl":  bool
    "shift": bool
    "super": bool
    "alt":   bool
}
```

**Example**

```json
{
    "type": "key",
    "key": "backspace",
    "ctrl": true,
    "shift": false,  
    "meta": false,
    "alt": false
}
```

### char - each character generated from window
Will run on each character generated from the user keyboard. 
This might be different from key event, since certain keyboard layouts can require 
multiple key presses to generate a single unicode character.

```
{
    "type": "char",
    "char": string (single unicode character)
}
```

**Example**
```json
{
    "type": "char",
    "char": "ö"
}
```

### size - columns & rows info
Guaranteed to always be the first thing sent on startup. Will then be sent each time the number of rows or columns change.
Also contains info about the size of each box in the grid.

```
{
    "type": "size"
    "cols": int
    "rows": int
    "colWidth": int
    "rowHeight": int
}
```

### error - errors related to sent requests
```
{
    "type": "error"
    "error": string
}
```

**Example**
```json
{
    "type": "error",
    "error": "char request was sent with empty char"
}
```
