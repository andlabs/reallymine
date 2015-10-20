>The [thing] is not revealed to comply with the “responsible disclosure” model.
>- the paper

Yeah, fuck that. My data is more important than your ethics code.

# Symwave 6316 hardcoded key (page 22)
This is stored in the upgrade blobs. Amazingly it's not stored directly, but rather embedded straight in the firmware source:

```
ROM:04001E84                 move.l  #$29A2607A,d0
ROM:04001E8A                 move.l  d0,-$28(a6)
ROM:04001E8E                 move.l  #$EA0B64AB,d0
ROM:04001E94                 move.l  d0,-$24(a6)
ROM:04001E98                 move.l  #$7BB3B9AB,d0
ROM:04001E9E                 move.l  d0,-$20(a6)
ROM:04001EA2                 move.l  #$A5698B40,d0
ROM:04001EA8                 move.l  d0,-$1C(a6)
ROM:04001EAC                 move.l  #$2E4793A6,d0
ROM:04001EB2                 move.l  d0,-$18(a6)
ROM:04001EB6                 move.l  #$8145C9CC,d0
ROM:04001EBC                 move.l  d0,-$14(a6)
ROM:04001EC0                 move.l  #$79946A01,d0
ROM:04001EC6                 move.l  d0,-$10(a6)
ROM:04001ECA                 move.l  #$840B34FE,d0
ROM:04001ED0                 move.l  d0,-$C(a6)
```

Enjoy.

I don't have any Symwave blocks to test this with, though :(
