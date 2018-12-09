# Démo Docker & isolation Bash

_Démo de l'isolation dans un container et prototype rudimentaire avec Bash_

### Démo Docker

```bash
docker run --name centos centos:latest bash
# Dans le docker:

# ps a
# yum check-update
# mount
```

### prototype Bash

```bash
# dans l'hote
# hostname toto
# sudo hostname toto
# unshare -u bash

## Dans le namespace
## readlink /proc/self/ns/uts
## hostname toto
## exit

# hostname
# man unshare
```
