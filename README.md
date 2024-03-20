# Simple Kubernetes Cluster demo

This is a simple Kubernetes cluster demo that spins up a web application and
demonstrates how to deploy to a Kubernetes cluster.

It is composed of three different parts:

1. A simple frontend application that queries the backend intermittently.
2. A backend application that serves a simple JSON response by querying a database.
3. A [PostgreSQL](https://www.postgresql.org/) database that stores a simple table.


# Running it

## docker-compose

To run the application, you can use `docker-compose`:

```bash
docker-compose up
```

> **_NOTE:_**  Remember to set up the `DOCKER_DEFAULT_PLATFORM` env var. For a
> mac M1/M2/M3, it would be `linux/arm64`.

You should be able to open the web interface in

http://localhost

and `curl` the backend in

http://localhost:8080


## Kubernetes deployment

The Kubernetes deployment will spin up a local Kubernetes cluster and create
two services, one for the frontend and one for the backend.

The backend service will spawn three replicas, and you will be able to see how
the load is balanced between them by refreshing the frontend.

> **_NOTE:_** As the frontend will keep the connection open, all requests will
> be sent to the same backend instance. To see the load balancing in action, you
> must hard refresh the browser (Ctrl+Shift+R/Cmd+Shift+R).

To deploy the application to a local Kubernetes cluster, you can use
[minikube](https://minikube.sigs.k8s.io/docs/).

```bash
minikube start
minikube mount frontend:/frontend
kubectl apply -f kubernetes/
```

Then you must tunnel to both the backend and frontend services:

```bash
minikube service backend
minikube service frontend
```

This should open both the web interface and even query the backend in your browser.

### Cleaning up

To clean up the deployment, you can use the following command:

```bash
kubectl delete -f kubernetes/
minikube stop
```


# Todo

- Use [Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets) to safely
  store the secrets in the repository.
- Set the network policy to `deny` all traffic and only allow the backend to connect
  to the database.
- Explore [Kustomize](https://kustomize.io/) to manage the Kubernetes manifests.
- Explore [ArgoCD](https://argoproj.github.io/argo-cd/) to manage the deployments.
