## Usage
Prepare a list of zones to provision, as defined in Go package
[`github.com/mittwald/go-powerdns/apis/zones`](https://pkg.go.dev/github.com/mittwald/go-powerdns/apis/zones#Zone).
YAML keys are names of struct fields converted to lowercase.

```yaml
- name: "example.com."
  resourcerecordsets:
  - name: "www.example.com."
    type: "A"
    ttl: 60
    records:
    - content: "192.0.2.1"
```

Create a Kubernetes ConfigMap containing the YAML list:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: powerdns-zone-provisioner-zone-file
data:
  zones.yml: |
    ---
    - name: "example.com."
      resourcerecordsets:
      - name: "www.example.com."
        type: "A"
        ttl: 60
        records:
        - content: "192.0.2.1"
```

Execute `powerdns-zone-provisioner` in a Kubernetes Job:
* Mount the ConfigMap as a volume
* Supply the path to the YAML via the `ZONE_FILE` environment variable:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: powerdns-zone-provisioner
spec:
  backoffLimit: 2048 # arbitrary large number
  template:
    metadata:
    spec:
      containers:
      - name: powerdns-zone-provisioner
        image: ghcr.io/parthenogen/powerdns-zone-provisioner:0
        imagePullPolicy: Always
        env:
        - name: SERVER_HOST
          value: auth-auth-api.powerdns-primary.svc.cluster.local
        - name: SERVER_PORT
          value: "8081"
        - name: API_KEY
          value: GGn7XHbLi1oJ5wSLb3qk
        - name: SERVER_ID
          value: localhost
        - name: ZONE_FILE
          value: /etc/powerdns-zone-provisioner/zones.yml
        - name: TIMEOUT
          value: 30s # per HTTP request
        volumeMounts:
        - mountPath: /etc/powerdns-zone-provisioner
          name: powerdns-zone-provisioner-zone-file
      restartPolicy: Never
      volumes:
      - name: powerdns-zone-provisioner-zone-file
        configMap:
          name: powerdns-zone-provisioner-zone-file
```

The following environment variables correspond to Authoritative Server settings:
| Environment Variable | Authoritative Server setting |
|---|---|
| `SERVER_PORT` | [`webserver-port`](https://doc.powerdns.com/authoritative/settings.html#webserver-port) |
| `API_KEY` | [`api-key`](https://doc.powerdns.com/authoritative/settings.html#setting-api-key) |
| `SERVER_ID` | [`server-id`](https://doc.powerdns.com/authoritative/settings.html#server-id) |
