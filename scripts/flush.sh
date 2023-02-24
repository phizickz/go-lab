kubectl delete -f ../manifests
docker system prune -f
docker build ../. -t go-web:lab
kubectl apply -f ../manifests