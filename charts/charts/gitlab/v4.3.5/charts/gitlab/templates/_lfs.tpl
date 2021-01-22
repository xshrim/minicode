{{/*
Generates a templated config for lfs key in gitlab.yml.

Usage:
{{ include "gitlab.appConfig.lfs.configuration" ( \
     dict                                         \
         "config" .Values.path.to.lfs.config      \
         "context" $                              \
     ) }}
*/}}
{{- define "gitlab.appConfig.lfs.configuration" -}}
lfs:
  enabled: {{ if kindIs "bool" .config.enabled }}{{ eq .config.enabled true }}{{ end }}
  {{- if .config.connection }}
  {{-   include "gitlab.appConfig.objectStorage.configuration" (dict "name" "lfs" "config" .config "context" .context) | nindent 2 }}
  {{- end -}}
{{- end -}}{{/* "gitlab.appConfig.lfs.configuration" */}}
