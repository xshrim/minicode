{{/*
Generates a templated object storage config.

Usage:
{{ include "gitlab.appConfig.objectStorage.configuration" ( \
     dict                                                   \
         "name" "STORAGE_NAME"                              \
         "config" .Values.path.to.objectstorage.config      \
         "context" $                                        \
     ) }}
*/}}
{{- define "gitlab.appConfig.objectStorage.configuration" -}}
object_store:
  enabled: {{ ne (default false .config.enabled) false }}
  {{- if ne .name "object_store" }}
  remote_directory: {{ .config.bucket }}
  {{- end }}
  direct_upload: true
  background_upload: false
  proxy_download: {{ or (not (kindIs "bool" .config.proxy_download)) .config.proxy_download }}
  {{- if and .config.enabled .config.connection }}
  connection: <%= YAML.load_file("/etc/gitlab/objectstorage/{{ .name }}").to_json() %>
  {{- else if .context.Values.global.minio.enabled }}
  {{-   include "gitlab.appConfig.objectStorage.connection.minio" . | nindent 2 }}
  {{- end -}}
{{- end -}}{{/* "gitlab.appConfig.objectStorage.configuration" */}}

{{/*
Generates a templated object storage connection settings for Minio.

Usage:
{{ include "gitlab.appConfig.objectStorage.connection.minio" ( \
     dict                                                  \
         "name" "STORAGE_NAME"                             \
         "config" .Values.path.to.objectstorage.config     \
     ) }}
*/}}
{{- define "gitlab.appConfig.objectStorage.connection.minio" -}}
connection:
  provider: AWS
  region: us-east-1
  host: {{ template "gitlab.minio.hostname" .context }}
  endpoint: {{ template "gitlab.minio.endpoint" .context }}
  path_style: true
  aws_access_key_id: "<%= File.read('/etc/gitlab/minio/accesskey').strip.dump[1..-2] %>"
  aws_secret_access_key: "<%= File.read('/etc/gitlab/minio/secretkey').strip.dump[1..-2] %>"
{{- end }}


{{/*
Generates a templated object storage secrets mounts.

Usage:
{{ include "gitlab.appConfig.objectStorage.mountSecrets" ( \
     dict                                                  \
         "name" "STORAGE_NAME"                             \
         "config" .Values.path.to.objectstorage.config     \
     ) }}
*/}}
{{- define "gitlab.appConfig.objectStorage.mountSecrets" -}}
# mount secret for {{ .name }}
{{- if .config.connection }}
- secret:
    name: {{ .config.connection.secret }}
    items:
      - key: {{ default "connection" .config.connection.key }}
        path: objectstorage/{{ .name }}
{{- end -}}
{{- end -}}{{/* "gitlab.appConfig.objectStorage.mountSecrets" */}}


{{/*
Generates a templated object storage bucket config

Usage:
{{ include "gitlab.appConfig.objectStorage.object" ( \
     dict                                                  \
         "name" "STORAGE_NAME"                             \
         "config" .Values.path.to.objectstorage.config     \
     ) }}
*/}}
{{- define "gitlab.appConfig.objectStorage.object" -}}
{{ .name }}:
  bucket: {{ .config.bucket }}
{{- end }}
