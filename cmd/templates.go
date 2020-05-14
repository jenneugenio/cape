package main

var configMapTemplate = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: coordinator-config-map
data:
  coordinator-config.yaml: |
    {{ . | nindent 4 }}
`

var secretTemplate = `
apiVersion: v1
kind: Secret
metadata:
  name: coordinator-config-secret
type: Opaque
data:
  coordinator-config.yaml: {{ . }}
`
