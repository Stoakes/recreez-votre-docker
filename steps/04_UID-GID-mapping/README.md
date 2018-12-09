# Mapping UID GID

_Lien entre utilisateur dans le container et un utilisateur de la machine hôte_

```bash
go run main.go run bash

# Dans le container:
# id
# ls -alh
# touch hello.md
# exit

# Dans l'hôte:
# ls -alh     # hello.md appartient à l'utilisateur courant
```

