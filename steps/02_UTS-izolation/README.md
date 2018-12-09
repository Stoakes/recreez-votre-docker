# Isolation UTS

_Pareil qu'avec unshare & Bash, mais avec Golang_

```bash
# Namespace accessibles uniquement à root (pour le moment), 
# donc compilation + exécution pour éviter de configurer l'environnement Golang pour root
go build .
sudo ./02_UTS-izolation run bash

# Dans le bash isolé:
# readlink /proc/self/ns/uts
# hostname toto
# hostname
# ps a
# exit
# hostname
```




