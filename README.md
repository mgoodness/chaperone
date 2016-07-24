# chaperone
[![Docker Repository on Quay](https://quay.io/repository/mgoodness/chaperone/status "Docker Repository on Quay")](https://quay.io/repository/mgoodness/chaperone)

Chaperone is a Golang file/directory watcher that POSTs to a URL when a file change is detected.

## Kubernetes pod manifest for Prometheus
```
---
kind: Pod
apiVersion: v1
metadata:
  name: prometheus-server
spec:
  containers:
    - name: chaperone
      args:
        - --dir=/etc/prometheus
        - --url=http://localhost:9090/-/reload
      image: quay.io/mgoodness/chaperone:v0.1
      volumeMounts:
        - name: config-volume
          mountPath: /etc/prometheus
          readOnly: true

    - name: prometheus
      args:
        - --config.file=/etc/prometheus/prometheus.yml
        - --storage.local.path=/prometheus
        - --web.console.libraries=/etc/prometheus/console_libraries
        - --web.console.templates=/etc/prometheus/consoles
      image: prom/prometheus:v1.0.1
      livenessProbe:
        httpGet:
          path: /status
          port: 9090
        initialDelaySeconds: 30
        timeoutSeconds: 30
      ports:
        - containerPort: 9090
          protocol: TCP
      volumeMounts:
        - name: config-volume
          mountPath: /etc/prometheus
          readOnly: true
        - name: data
          mountPath: /prometheus
  volumes:
    - name: config-volume
      configMap:
        name: prometheus-server

    - name: data
      persistentVolumeClaim:
        claimName: prometheus-server
```
