# kubehook
Monitor Kubernetes cluster operations on various webhook services
![N|Solid](https://s.phineas.io/share/ezgif-5-e34ac1716e.gif)

## Todo
- Add Slack webhook support
- Create modular components for services
- Add more Kubernetes operations to watchlist
- Add support for miscellaneous Kubernetes cluster services
- Create Discord bot for performing ops

## How to Run
This is designed to run as a single replica on a Google Kubernetes Engine cluster. 
1. Edit config.json, add your webhook URL
2. Modify enabled events to your needs
3. Build the image and submit it to container registry (takes a while due to Kube api):
`gcloud container builds submit --tag gcr.io/cx-network-204116/kubehook`
5. Edit kubehook-worker.yaml to support your image tag
6. Connect to your cluster
7. Create the role binding:
`kubectl create rolebinding kh --clusterrole=view --serviceaccount=default:default`
8. Create the deployment
`kubectl create -f kubehook-worker.yaml`