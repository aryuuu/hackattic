## Notes

from the website:
```
The pack contains, always in the following order:

a regular int (signed), to start off
an unsigned int
a short (signed) to make things interesting
a float because floating point is important
a double as well
another double but this time in big endian (network byte order)
```


the problem endpoint returned a base64 string which when decoded results in 32 bytes of data


- the first 4 is an int (signed)
- the next 4 bytes is unsigned int
- the next 2 bytes is signed short int
- the next 4 bytes is a float
- the next 8 bytes is a double int
- the next 8 bytes is a double int in big endian

for example:

-- signed int
180
172
145
133

-- unsigned int
18
172
74
216

--  short
29
223

-- padding
0
0

-- float
110
91
200
66

-- double (little endian)
50
147
123
25
168
93
95
64

-- double (big endian)
64
95
93
168
25
123
147
50
