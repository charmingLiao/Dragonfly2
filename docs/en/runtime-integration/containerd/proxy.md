# Use dfget daemon as HTTP proxy for containerd

Currently, `ctr` command of containerd doesn't support private registries with `registry-mirrors`,
in order to do so, we need to use HTTP proxy for containerd.

## Quick Start

### Step 1: Generate CA certificate for HTTP proxy

Generate a CA certificate private key.

```bash
openssl genrsa -out ca.key 2048
```

Open openssl config file `openssl.conf`.
Note set `basicConstraints` to true, that you can modify the values.

```text
[ req ]
#default_bits = 2048
#default_md = sha256
#default_keyfile = privkey.pem
distinguished_name = req_distinguished_name
attributes = req_attributes
extensions               = v3_ca
req_extensions           = v3_ca

[ req_distinguished_name ]
countryName = Country Name (2 letter code)
countryName_min = 2
countryName_max = 2
stateOrProvinceName = State or Province Name (full name)
localityName = Locality Name (eg, city)
0.organizationName = Organization Name (eg, company)
organizationalUnitName = Organizational Unit Name (eg, section)
commonName = Common Name (eg, fully qualified host name)
commonName_max = 64
emailAddress = Email Address
emailAddress_max = 64

[ req_attributes ]
challengePassword = A challenge password
challengePassword_min = 4
challengePassword_max = 20

[ v3_ca ]
basicConstraints         = CA:TRUE
```

Generate the CA certificate.

```bash
openssl req -new -key ca.key -nodes -out ca.csr -config openssl.conf
openssl x509 -req -days 36500 -extfile openssl.conf \
    -extensions v3_ca -in ca.csr -signkey ca.key -out ca.crt
```

### Step 2: Configure dfget daemon

To use dfget daemon as HTTP proxy, first you need to append a proxy rule in
`/etc/dragonfly/dfget.yaml`, This will proxy
`your.private.registry`'s requests for image layers:

```yaml
proxy:
  security:
    insecure: true
  tcpListen:
    listen: 0.0.0.0
    port: 65001
  proxies:
    - regx: blobs/sha256.*
  hijackHTTPS:
    # CA certificate's path used to hijack https requests
    cert: ca.crt
    key: ca.key
    hosts:
      - regx: your.private.registry
```

### Step 3: Configure containerd

Set dfget damone as `HTTP_PROXY` and `HTTPS_PROXY` for containerd in
`/etc/systemd/system/containerd.service.d/http-proxy.conf`:

```toml
[Service]
Environment="HTTP_PROXY=http://127.0.0.1:65001"
Environment="HTTPS_PROXY=http://127.0.0.1:65001"
```

### Step 4: Pull images with proxy

Through the above steps, we can start to validate if Dragonfly works as expected.

And you can pull the image as usual, for example:

```bash
ctr image pull your.private.registry/namespace/image:latest
```

## Custom assets

### Registry uses a self-signed certificate

If your registry uses a self-signed certificate, you can either choose to
ignore the certificate error with:

```yaml
proxy:
  security:
    insecure: true
  tcpListen:
    listen: 0.0.0.0
    port: 65001
  proxies:
    - regx: blobs/sha256.*
  hijackHTTPS:
    # CA certificate's path used to hijack https requests
    cert: ca.crt
    key: ca.key
    hosts:
      - regx: your.private.registry
        insecure: true
```

Or provide a certificate with:

```yaml
proxy:
  security:
    insecure: true
  tcpListen:
    listen: 0.0.0.0
    port: 65001
  proxies:
    - regx: blobs/sha256.*
  hijackHTTPS:
    # CA certificate's path used to hijack https requests
    cert: ca.crt
    key: ca.key
    hosts:
      - regx: your.private.registry
        certs: ["server.crt"]
```

You can get the certificate of your server with:

```bash
openssl x509 -in <(openssl s_client -showcerts \
    -servername your.domain.com -connect your.domain.com:443 -prexit 2>/dev/null)
```
