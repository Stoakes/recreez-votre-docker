# Pivot root

_Une distribution Linux n'est finalement qu'un système de fichier au dessus d'un noyau_

```bash
# Récuperer un filesystem centos:
docker run --name fs centos
docker export fs -o centos.tar
docker rm fs
mkdir centos  && tar -xf centos.tar -C centos

# Editez main.go#L60 pour définir le bon emplacement de rootfsPath

go run *.go run bash

# Dans le container
# python version
# ps 
```
