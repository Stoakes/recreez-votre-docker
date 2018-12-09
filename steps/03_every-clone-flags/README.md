# Tous les clones flags

_Passer tous les clones flags ne suffit pas à tout isoler. Par contre on peut lancer des namespace sans être root._

```bash
go run *.go run bash

# Dans le container:
# id        # Pas root
# ls -alh
```

