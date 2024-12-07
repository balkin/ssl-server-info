# ssl-server-info

**ssl-server-info** is a Go-based server that provides SSL/TLS information about incoming HTTPS requests. It supports client certificate verification, custom response delays, and status codes. This server is ideal for debugging and inspecting SSL/TLS-related configurations.

---

## Features

- **SSL/TLS Information**: Returns detailed information about the SSL/TLS certificate used for the connection (e.g., subject, issuer, validity).
- **Custom Response Delays**: Simulate delays in responses via the `delay` query parameter or path.
- **Custom Status Codes**: Modify the HTTP status code in the response via the `statusCode` query parameter or path.
- **Configurable SSL Certificates**: Use environment variables to configure SSL certificate files.
- **Port Customization**: The server can be run on any port (defaults to `443`).

---

## Endpoints

### `/ssl/json`
Returns SSL/TLS and request metadata in a JSON response.

### Optional Query Parameters
- `delay`: Simulate a delay in milliseconds before responding (e.g., `/ssl/json?delay=1000`).
- `statusCode`: Set a custom HTTP status code for the response (e.g., `/ssl/json?statusCode=500`).

### Path Variants
- `/ssl/json/delay/{value}`: Equivalent to `?delay={value}` (e.g., `/ssl/json/delay/1000`).
- `/ssl/json/statusCode/{value}`: Equivalent to `?statusCode={value}` (e.g., `/ssl/json/statusCode/500`).

### Combining and merging parameters 

In case both query parameter and path parameter are provided (`/ssl/json/delay/1000?delay=9999`), only query parameter will be considered.

Feel free to mix & match `statusCode` and `delay` parameters as you please. E.g., `/ssl/json/statusCode/222?delay=5000` or `/ssl/json/statusCode/222/delay/5000` and `/ssl/json/delay/5000/statusCode/222` are equal.

### Example Response

```json
{
  "server": "github.com/balkin/ssl-server-info",
  "https": "on",
  "headerContentType": "application/json",
  "headerAccept": "*/*",
  "headerUserAgent": "curl/7.68.0",
  "httpHost": "localhost",
  "httpServerAddr": "127.0.0.1:443",
  "requestProtocol": "HTTP/2.0",
  "requestMethod": "GET",
  "requestUri": "/ssl/json",
  "requestTmestamp": 1696523567,
  "sslSubject": "localhost",
  "sslIssuer": "localhost",
  "sslNotBefore": "2024-01-01T00:00:00Z",
  "sslNotAfter": "2025-01-01T00:00:00Z"
}