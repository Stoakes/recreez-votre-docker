# CGroups

_Limites sur les consommations de resources systèmes_

```bash
# Récuperer un filesystem centos:
docker run --name fs centos
docker export fs -o centos.tar
docker rm fs
mkdir centos  && tar -xf centos.tar -C centos
cp hungry.py ./centos/root/hungry.py

# Editez main.go#L81 pour définir le bon emplacement de rootfsPath

# Definition des CGroups
sudo mkdir -p /sys/fs/cgroup/memory/demo
AA=$USER
sudo chgrp -R $AA /sys/fs/cgroup/memory/demo
sudo chmod -R g+w /sys/fs/cgroup/memory/demo

go run *.go run bash

# Dans le container:
# python /root/hungry.py
```
