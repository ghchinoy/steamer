# Porkbun API V3 Documentation Reference

Source of truth: [Porkbun API V3 Documentation](https://porkbun.com/api/json/v3/documentation)

## General Information
- **API Hostname:** `api.porkbun.com`
- **Protocol:** HTTP POST (All endpoints use POST)
- **Authentication:** API Key and Secret are passed in the JSON body.
- **Error Handling:** Any HTTP response code other than 200 is an error. Errors are also indicated by the `status: "ERROR"` field in the JSON response.

## Authentication (Ping)
Verify API credentials and retrieve your IP address.

- **Endpoint:** `https://api.porkbun.com/api/json/v3/ping`
- **Request:**
```json
{
  "secretapikey": "YOUR_SECRET_API_KEY",
  "apikey": "YOUR_API_KEY"
}
```
- **Response:**
```json
{
  "status": "SUCCESS",
  "yourIp": "66.249.72.71"
}
```

## Domain Functionality

### List All Domains
- **Endpoint:** `https://api.porkbun.com/api/json/v3/domain/listAll`
- **Request:**
```json
{
  "secretapikey": "YOUR_SECRET_API_KEY",
  "apikey": "YOUR_API_KEY",
  "start": "0", // Optional index to start at, increment by 1000 to paginate.
  "includeLabels": "yes" // Optional, set to "yes" to include label information.
}
```

### Domain Check (Availability)
- **Endpoint:** `https://api.porkbun.com/api/json/v3/domain/checkDomain/DOMAIN`
- **Rate Limits:** Domain checks are rate-limited. Example: 1 check per 10 seconds.

### Glue Records
- **Endpoints:**
    - Create: `https://api.porkbun.com/api/json/v3/domain/createGlue/DOMAIN/HOST`
    - Update: `https://api.porkbun.com/api/json/v3/domain/updateGlue/DOMAIN/HOST`
    - Delete: `https://api.porkbun.com/api/json/v3/domain/deleteGlue/DOMAIN/HOST`
    - Get: `https://api.porkbun.com/api/json/v3/domain/getGlue/DOMAIN`
- **Request (Create/Update):**
```json
{
  "secretapikey": "YOUR_SECRET_API_KEY",
  "apikey": "YOUR_API_KEY",
  "ips": ["192.168.1.1", "2001:db8::1"]
}
```

## DNS Functionality

### Retrieve Records
- **Endpoint:** `https://api.porkbun.com/api/json/v3/dns/retrieve/DOMAIN/[ID]`
- **Request:**
```json
{
  "secretapikey": "YOUR_SECRET_API_KEY",
  "apikey": "YOUR_API_KEY"
}
```

### Create Record
- **Endpoint:** `https://api.porkbun.com/api/json/v3/dns/create/DOMAIN`
- **Request:**
```json
{
  "secretapikey": "YOUR_SECRET_API_KEY",
  "apikey": "YOUR_API_KEY",
  "name": "www", // Optional subdomain, blank for root, * for wildcard.
  "type": "CNAME", // A, MX, CNAME, ALIAS, TXT, NS, AAAA, SRV, TLSA, CAA, HTTPS, SVCB, SSHFP
  "content": "ghs.google.com",
  "ttl": "600", // Optional, min/default 600.
  "prio": "0", // Optional priority for MX/SRV.
  "notes": "" // Optional notes.
}
```

### Edit Record by Domain and ID
- **Endpoint:** `https://api.porkbun.com/api/json/v3/dns/edit/DOMAIN/ID`
- **Request:** Same fields as Create Record.
- **Notes Behavior:** Passing an empty string `""` will clear the notes; passing `null` will make no changes.

### Delete Record by Domain and ID
- **Endpoint:** `https://api.porkbun.com/api/json/v3/dns/delete/DOMAIN/ID`
- **Request:**
```json
{
  "secretapikey": "YOUR_SECRET_API_KEY",
  "apikey": "YOUR_API_KEY"
}
```

## DNSSEC Functionality

### Create/Get/Delete DNSSEC Records
- **Endpoints:**
    - Create: `https://api.porkbun.com/api/json/v3/dns/createDnssecRecord/DOMAIN`
    - Get: `https://api.porkbun.com/api/json/v3/dns/getDnssecRecords/DOMAIN`
    - Delete: `https://api.porkbun.com/api/json/v3/dns/deleteDnssecRecord/DOMAIN/KEYTAG`
- **Create Request:**
```json
{
  "secretapikey": "YOUR_SECRET_API_KEY",
  "apikey": "YOUR_API_KEY",
  "keyTag": "64087",
  "alg": "13",
  "digestType": "2",
  "digest": "...",
  "maxSigLife": "", // Optional
  "keyDataFlags": "", // Optional
  "keyDataProtocol": "", // Optional
  "keyDataAlgo": "", // Optional
  "keyDataPubKey": "" // Optional
}
```

## SSL Functionality

### Retrieve SSL Bundle
- **Endpoint:** `https://api.porkbun.com/api/json/v3/ssl/retrieve/DOMAIN`
- **Request:**
```json
{
  "secretapikey": "YOUR_SECRET_API_KEY",
  "apikey": "YOUR_API_KEY"
}
```
- **Response:**
```json
{
  "status": "SUCCESS",
  "certificatechain": "-----BEGIN CERTIFICATE-----\n...",
  "privatekey": "-----BEGIN PRIVATE KEY-----\n...",
  "publickey": "-----BEGIN PUBLIC KEY-----\n..."
}
```
