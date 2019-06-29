# rdma-ds

This dameon-set will run on every node and will be a RESTful API endpoint for querying data about the node that it is currently running on.

## Test API
To test the API run the following commands:
```
cd src
go test
```

## Environmental Variables
The following variables are able to be modefied to change the server running:
  - PORT - default 54005 - the port number for the server to run on

## Building Binaries (for testing or for use)

### Client
To use the client library import `github.com/swrap/rdma-ds/src` and utilize the built in supported functions.

### Server Build and Run
To start the server run the script below:
```
./build_run_server.sh
```
If you want to specify a port number, set the environment variable `PORT` like below:
```
PORT=40007 ./build_run_server.sh
```

### Server Build for all Architectures
To build binaries for all architures and place them in a bin directory run:
```
./build_server.sh
```

## Docker
To build the docker image run the following in the root directory of this repo:
```
docker build -t rdma-ds:latest .
```
This will build an image with `rdma-ds` as its name.

To run this image run the following command:
```
docker run -it --rm --name test -e PORT=5000 --network host rdma-ds-v1
```
This will do the following:
  - `-it` - runs in interactive mode
  - `--rm` - removes the container when stopped
  - `-e PORT=5000` - specifies the port for the server to run on, if none is specified it will default to port 54005
  - `--network host` - will share the network with your host OS, so you can access the api by going to localhost:5000
  - `rdma-ds` - the name of the image. NOTE if you build with a different image name, you will need to change this

### Avoiding Dockerhub
If you want to avoid docker hub completely, you can save the image in a tar and than load it in.
Save command:
```
docker save <image-name>:<image-tag> > <save-name>.tar
```
Load command:
```
docker load < <save-name>.tar
```
Ex:
```
docker save rdma-ds:latest > rdma-ds.tar
docker load < rdma-ds.tar
```

## Kubernetes
In order to run this application in kuberentes you must execute the following command:
```
kubectl apply -f <dameonset-config-file>.yaml
```
This will spin up this repo with the dameonset. Note the current implementation you must first put the docker images on the remote nodes before running because this is a private repo.

Deleting dameonset:
```
kubectl delete ds/<dameonset-name> --namespace kube-system
```

Updating dameonset by changing the image:
```
kubectl set image ds/<daemonset-name> <container-name>=<container-new-image> --namespace=kube-system
```

Check updating status
```
kubectl rollout status ds/<daemonset-name> --namespace=kube-system
```

## How it works
Assumptions:
 - VFs in use will move the `/sys/class/net/<pf-name>/device/virtfn<vf-number>/net` file out of namespace when the VF is in use
 - Must have config-map set for both pfNetDevices and maxPfBandwidth