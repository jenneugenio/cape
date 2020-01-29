k8s_yaml(helm('charts/connector', values=['charts/connector/values.yaml']))
k8s_yaml(helm('charts/controller', values=['charts/connector/values.yaml']))

docker_build('dropoutlabs/privacyai:latest', '.')