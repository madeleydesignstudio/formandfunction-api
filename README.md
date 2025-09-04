# Form & Function - Go API Service

A high-performance Go API service that provides both HTTP REST and gRPC interfaces for steel beam data management and structural engineering applications.

## 🏗️ Architecture

This service operates as part of a microservice architecture:

```
Frontend Clients → HTTP REST API (Port 8080) → Go API Service
                                            ↑
Python Calc Engine → gRPC API (Port 9090) ──┘
```

## 🚀 Features

- **Dual Protocol Support**: HTTP REST for frontend clients, gRPC for backend services
- **Steel Beam Database**: Comprehensive UK steel beam section data
- **CORS Enabled**: Ready for web frontend integration
- **Type-Safe gRPC**: Protobuf-based service definitions
- **Health Monitoring**: Built-in health check endpoints
- **Production Ready**: Optimized for Railway deployment

## 📋 API Endpoints

### HTTP REST API (Port 8080)

| Method | Endpoint | Description | Response |
|--------|----------|-------------|----------|
| `GET` | `/` | Service information | Service details |
| `GET` | `/health` | Health check | Service status |
| `GET` | `/beams` | Get all steel beams | Array of beam objects |
| `GET` | `/beams/{section}` | Get specific beam | Single beam object |
| `POST` | `/beams` | Create new beam | Created beam object |
| `PUT` | `/beams/{section}` | Update existing beam | Updated beam object |
| `DELETE` | `/beams/{section}` | Delete beam | Success/error message |

### gRPC API (Port 9090)

| Service | Method | Description |
|---------|--------|-------------|
| `SteelBeamService` | `GetBeams()` | Retrieve all beams |
| `SteelBeamService` | `GetBeam(section)` | Get specific beam |
| `SteelBeamService` | `CreateBeam(data)` | Create new beam |

## 🛠️ Local Development

### Prerequisites

- Go 1.21 or later
- Protocol Buffers compiler (`protoc`)
- Go protobuf plugins

### Setup

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd formandfunction-api
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Generate protobuf files**
   ```bash
   chmod +x build_proto.sh
   ./build_proto.sh
   ```

4. **Run the service**
   ```bash
   export PORT=8080
   export GRPC_PORT=9090
   go run .
   ```

### Testing

```bash
# Test HTTP endpoints
curl http://localhost:8080/health
curl http://localhost:8080/beams

# Test gRPC (requires grpcurl)
grpcurl -plaintext localhost:9090 list
grpcurl -plaintext localhost:9090 steelbeam.SteelBeamService/GetBeams
```

## 🚀 Production Deployment (Railway)

### Environment Variables

```
PORT=8080
GRPC_PORT=9090
GO_ENV=production
```

### Custom Domain

Set up your custom domain in Railway to point to: `api.itsformfunction.com`

### Build Process

Railway automatically:
1. Runs `build_proto.sh` to generate protobuf files
2. Executes `go build -o main .`
3. Starts the service with `./main`

### Port Configuration

Ensure Railway exposes both ports:
- **8080**: HTTP REST API (for frontend clients)
- **9090**: gRPC API (for backend services)

## 🔧 Configuration

### CORS Configuration

For production, update CORS settings in `main.go`:

```go
app.Use(cors.New(cors.Config{
    AllowOrigins: []string{
        "https://yourfrontend.com",
        "https://app.yourfrontend.com",
    },
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders: []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
}))
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `GRPC_PORT` | gRPC server port | `9090` |
| `GO_ENV` | Environment mode | `development` |

## 📊 Monitoring

### Health Checks

- **HTTP**: `GET /health`
- **gRPC**: Service reflection enabled

### Logging

All requests and gRPC calls are logged with:
- Request/response times
- Status codes
- Error details
- Client information

## 🔒 Security

### Production Security Checklist

- [ ] Enable HTTPS/TLS for HTTP endpoints
- [ ] Enable TLS for gRPC communication
- [ ] Restrict CORS origins to your frontend domains
- [ ] Implement API rate limiting
- [ ] Add request validation middleware
- [ ] Set up proper error handling (don't expose internal errors)

## 🧪 Testing Your Deployment

Use the verification script:

```bash
python ../verify_production_deployment.py --api-only
```

Or test manually:

```bash
# Health check
curl https://api.itsformfunction.com/health

# Get all beams
curl https://api.itsformfunction.com/beams

# Test specific beam
curl https://api.itsformfunction.com/beams/UB406x178x74
```

## 📚 Steel Beam Data

The service includes comprehensive UK steel beam data with properties:

- Section designation (e.g., "UB406x178x74")
- Mass per metre (kg/m)
- Dimensional properties (depth, width, thickness)
- Section properties (moments of inertia, elastic moduli)
- Local buckling ratios
- Torsional properties

## 🔧 Development Commands

```bash
# Format code
go fmt ./...

# Run tests
go test ./...

# Build binary
go build -o main .

# Generate protobuf files
./build_proto.sh

# Run with specific ports
PORT=8080 GRPC_PORT=9090 go run .
```

## 🤝 Frontend Integration

### TanStack Query Example

```typescript
const API_BASE_URL = 'https://api.itsformfunction.com'

export const beamApi = {
  getBeams: () => 
    fetch(`${API_BASE_URL}/beams`).then(res => res.json()),
  
  getBeam: (section: string) =>
    fetch(`${API_BASE_URL}/beams/${section}`).then(res => res.json()),
    
  createBeam: (beam: SteelBeam) =>
    fetch(`${API_BASE_URL}/beams`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(beam)
    }).then(res => res.json()),
}
```

### React Hook Example

```tsx
import { useQuery } from '@tanstack/react-query'

export function useBeams() {
  return useQuery({
    queryKey: ['beams'],
    queryFn: beamApi.getBeams,
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}
```

## 📈 Performance

- **HTTP REST**: Optimized for frontend clients with JSON responses
- **gRPC**: High-performance binary protocol for service communication
- **Memory**: Efficient in-memory beam data storage
- **Concurrency**: Go's goroutines handle multiple concurrent requests

## 🐛 Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure ports 8080 and 9090 are available
2. **Protobuf errors**: Run `./build_proto.sh` to regenerate files
3. **CORS issues**: Check frontend domain is in allowed origins
4. **gRPC connection**: Verify gRPC port is exposed in Railway

### Debug Mode

Set environment variable for verbose logging:

```bash
export LOG_LEVEL=debug
go run .
```

## 📄 License

[Your License Here]

---

🚀 **Production URL**: https://api.itsformfunction.com

For the complete microservice architecture documentation, see `../MICROSERVICE_DEPLOYMENT_GUIDE.md`
