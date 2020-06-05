{{/*
Expand the name of the chart.
*/}}
{{- define "cape.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a fully qualified app name.
*/}}
{{- define "cape.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "cape.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create the name of the service account
*/}}
{{- define "cape.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
  {{ default (include "cape.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
  {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Allow for overriding namespace
*/}}
{{- define "cape.namespace" -}}
  {{- if .Values.namespaceOverride -}}
    {{- .Values.namespaceOverride -}}
  {{- else -}}
    {{- .Release.Namespace -}}
  {{- end -}}
{{- end -}}

{{/*
Create a name for the configuration secret
*/}}
{{- define "cape.configSecret" -}}
{{- default (printf "%s-cm" (include "cape.fullname" .)) .Values.token_secret -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "cape.labels" -}}
helm.sh/chart: {{ include "cape.chart" . }}
{{ include "cape.selectorLabels" . }}
{{- if or .Chart.AppVersion .Values.image.tag }}
app.kubernetes.io/version: {{  default .Chart.AppVersion .Values.image.tag | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "cape.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cape.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}
