# GRPC-UI
- You can use make grpc ui
```bash
make grpcui-call
```
- Or you can call directly
```bash
grpcui -cacert configs/ca.crt localhost:50051
```

# PVZ CRUD

## PVZ Model 
- PVZ_ID
- Name Not unique
- Address
- Contact
- CreatedAt
- UpdatedAt


Server - > Middleware -> handler localhost:9000/pvz/create

## CURL
We have to enter Homework-1 (CRUD PVZ) project then call this curl request from our cmd 
- Create PVZ
    ```bash
    curl -k --cert configs/ca.crt -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Basic $(echo -n 'Homework_3:test' | base64)" \
    -d '{
    "name": "aibek",
    "address": "berlin",
    "contact": "+213123131231"
    }' \
    http://localhost:9000/pvz_v1/create
    ```
- Get One PVZ
    ```bash
   curl -k --cert configs/ca.crt -X GET \
    -H "Content-Type: application/json" \
    -H "Authorization: Basic $(echo -n 'Homework_3:test' | base64)" \
    http://localhost:9000/pvz_v1/get/6
  ```
- List PVZ
  ```bash
    curl -k --cert configs/ca.crt -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Basic $(echo -n 'Homework_3:test' | base64)" \
    -d '{
    "currentPage": 1,
    "itemsPerPage": 10
    }' \
    http://localhost:9000/pvz_v1/list
    ```
- Update PVZ
    ```bash
  curl -k --cert configs/ca.crt -X PUT \
    -H "Content-Type: application/json" \
    -H "Authorization: Basic $(echo -n 'Homework_3:test' | base64)" \
    -d '{
    "id": 5,
    "name": "john",
    "address":"turnitin",
    "contact":"+2131312311"
    }' \
    http://localhost:9000/pvz_v1/update
  ```
- Delete PVZ
    ```bash
  curl -k --cert configs/ca.crt -X DELETE \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic $(echo -n 'Homework_3:test' | base64)" \
  http://localhost:9000/pvz_v1/delete/5
  ```


# Order CRUD

## Order Model
- ID
- Order_ID
- Weight
- Client_ID
- Returned_at
- Accepted_at
- Issued_at
- Expires_at
- Created_at
- Updated_at


- Receive Order
    ```bash
    curl -k --cert configs/ca.crt -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Basic $(echo -n 'Homework_3:test' | base64)" \
    -d '{
    "expireTimeDuration": 30,
    "orderID": 2,
    "clientID": 1,
    "weight": 9
    }' \
    http://localhost:9000/order_v1/receive
    ```

- Issue Slice of Orders
  ```bash
    curl -k --cert configs/ca.crt -X PUT \
    -H "Content-Type: application/json" \
    -H "Authorization: Basic $(echo -n 'Homework_3:test' | base64)" \
    -d '{
    "orderIDs": [1, 2]
    }' \
    http://localhost:9000/order_v1/issue
    ```


- Accept Order
    ```bash
    curl -k --cert configs/ca.crt -X PUT \
    -H "Content-Type: application/json" \
    -H "Authorization: Basic $(echo -n 'Homework_3:test' | base64)" \
    -d '{
    "orderID": 3,
    "clientID": 11
    }' \
    http://localhost:9000/order_v1/accept
    ```

- Turn In Order
    ```bash
    curl -k --cert configs/ca.crt -X DELETE \
    -H "Content-Type: application/json" \
    -H "Authorization: Basic $(echo -n 'Homework_3:test' | base64)" \
    http://localhost:9000/order_v1/turn_in/6
    ```


- Returned List of Order
    ```bash
    curl -k --cert configs/ca.crt -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Basic $(echo -n 'Homework_3:test' | base64)" \
    -d '{
    "currentPage": 1,
    "itemsPerPage": 10
    }' \
    http://localhost:9000/order_v1/returns
    ```


- List of Orders
    ```bash
    curl -k --cert configs/ca.crt -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Basic $(echo -n 'Homework_3:test' | base64)" \
    -d '{
    "currentPage": 1,
    "itemsPerPage": 10
    }' \
    http://localhost:9000/order_v1/list
    ```

- Issue order with Box
    ```bash
    curl -k --cert configs/ca.crt -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Basic $(echo -n 'Homework_3:test' | base64)" \
    -d '{
    "orderID": 3,
    "boxID": 2
    }' \
    http://localhost:9000/order_v1/issue_with_box
    ```

- List of Unique Clients
    ```bash
    curl -k --cert configs/ca.crt -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Basic $(echo -n 'Homework_3:test' | base64)" \
    -d '{
    "currentPage": 1,
    "itemsPerPage": 10
    }' \
    http://localhost:9000/order_v1/unique_clients 
    ```

# Package CRUD

## Package Model
- ID
- Name
- Cost
- Is_check
- Weight
- Deleted_at
- Created_at
- Updated_at


- Create Package
    ```bash
    curl -k --cert configs/ca.crt -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Basic $(echo -n 'Homework_3:test' | base64)" \
    -d '{
    "name": "tico",
    "cost": 100.0,
    "isCheck": true,
    "weight": 200
    }' \
    http://localhost:9000/box_v1/create
    ```

- Delete Package
    ```bash
    curl -k --cert configs/ca.crt -X DELETE \
    -H "Content-Type: application/json" \
    -H "Authorization: Basic $(echo -n 'Homework_3:test' | base64)" \
    http://localhost:9000/box_v1/delete/2
    ```


- List Packages
    ```bash
    curl -k --cert configs/ca.crt -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Basic $(echo -n 'Homework_3:test' | base64)" \
    -d '{
    "currentPage": 1,
    "itemsPerPage": 10
    }' \
    http://localhost:9000/box_v1/list
    ```

In order for the TLS certificate to work perfectly when used by our gRPC server, it needs to include "localhost" as a valid hostname. This hostname (Common Name, CN, or part of the Subject Alternative Names, SAN) must be specified during the certificate generation process. If it's missing, clients (such as grpcui) that attempt to verify the server's identity against the hostname they connect to (in this case, localhost) will fail.

### Create Credentials Steps

We need to ensure that when we generate our server certificate, "localhost" is included as either the CN or as a SAN. Here's how to do it:

1. **Generate the Server Certificate with "localhost" as a SAN**: The most secure and recommended way to specify hostnames and IP addresses in certificates is through the use of SANs (Subject Alternative Names). Here's how to create a new CSR (Certificate Signing Request) that includes "localhost" as a SAN:

   **Step 1**: Create a new configuration file for OpenSSL to include the SAN. Let's call this file `san.cnf` in configs folder:

    ```conf
    [req]
    default_bits       = 2048
    prompt             = no
    default_md         = sha256
    req_extensions     = req_ext
    distinguished_name = dn

    [dn]
    C=IT
    ST=Turin
    L=Turin
    O=Ozon
    OU=Backend
    emailAddress=aibek@gmail.com
    CN=localhost

    [req_ext]
    subjectAltName = @alt_names

    [alt_names]
    DNS.1   = localhost
    DNS.2   = my_app
    IP.1    = 127.0.0.1
    ```
   **Step 2**: Generate the CA private key:
    ```bash
      openssl genrsa -out configs/ca.key 2048
   ```
   **Step 3**: Create a self-signed CA certificate:
    ```bash
      openssl req -x509 -new -nodes -key configs/ca.key -sha256 -days 1024 -out configs/ca.crt -config configs/san.cnf
   ```
   **Step 4**: Generate the server private key:
    ```bash
      openssl genrsa -out configs/server.key 2048
   ```
   **Step 5**: Generate the CSR using the server key and SAN details:
    ```bash
      openssl req -new -key configs/server.key -out configs/server.csr -config configs/san.cnf
   ```
   **Step 6**: Sign the new CSR with your CA:
    ```bash
      openssl x509 -req -in configs/server.csr -CA configs/ca.crt -CAkey configs/ca.key -CAcreateserial -out configs/server.crt -days 365 -sha256 -extfile configs/san.cnf -extensions req_ext
   ```
2. **Verify the Certificate**: After generating the new certificate, verify that it includes the SAN:
    **Step 1**: Verify the certificate chain:
    ```bash
   openssl verify -CAfile configs/ca.crt configs/server.crt
    ```
   **Step 2**: Inspect the certificate details.:
    ```bash
    openssl x509 -in configs/server.crt -text -noout
   ```

   Look for the `X509v3 extensions` section and confirm that `Subject Alternative Name` includes both `DNS:localhost` and `IP Address:127.0.0.1`.

3. **Restart our gRPC Server**: Use the newly created `server.crt` and `server.key` to restart our gRPC server.

4. **Reattempt Connection with `grpcui`**: Now, try connecting again with `grpcui`:

    ```bash
    grpcui -cacert configs/ca.crt localhost:50051
    ```

By following these steps, we ensure that our server's certificate is correctly set up to be verified against "localhost". This setup allows `grpcui` and other clients to securely connect and verify the server's identity.