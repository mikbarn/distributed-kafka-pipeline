# Custom KRaft Kafka Cluster & Concurrent Go Load-Tester

A bare-metal, cloud-native workspace built to deploy a multi-node, ZooKeeperless Kafka (KRaft) cluster natively inside Kubernetes and stress-test it using a custom-engineered, highly concurrent Go ingestion pipeline. 

## 🛠️ The Architecture & What I Actually Built

### 1. The Distributed Ingestion Broker (Kubernetes & KRaft)
Instead of relying on out-of-the-box Helm charts or legacy ZooKeeper topologies, this configuration structures a raw **Kafka KRaft cluster running natively in Kubernetes namespaces** using individual `StatefulSets`:
* **Decoupled Roles:** Deploys dedicated controller nodes (`kafka-kraft`) and data broker worker nodes (`kafka`).
* **Headless Internal Routing:** Utilizes a `clusterIP: None` configuration for localized, identity-stable internal DNS broker discovery.
* **Storage Ingress:** Employs volume claim templates mapping directly to custom local directories (`/kafkamnt`) on the host system to ensure persistence across pod cycles.

### 2. High-Throughput Load Generator (Go / Sarama)
The `producer` directory holds a native Go service built on the IBM/Sarama library that behaves like a custom performance harness:
* **Decoupled Concurrency:** Spawns a parallel matrix of isolated workers (`runner` goroutines) executing synchronous queue pushes, randomized back-offs, and state evaluations.
* **Asynchronous Metrics Aggregation:** Rather than blocking execution loops, workers feed individual performance logs across decoupled channels (`cWin`, `cSent`, `cFail`) into a standalone pipeline metric `collector`.
* **Graceful Lifecycle Management:** Uses deep `context.Context` cancellation patterns and structural OS interrupt capture (`syscall.SIGINT`/`SIGTERM`) to guarantee that active producers cleanly flush buffers and exit without leaking connections.

### 3. Scala Stream Consumer (Scala / SBT)
* An isolated event-processing microservice packaged inside a lightweight JVM container to continually scale, pull partition batches, and drain active data queues.

---

## 📂 Repository Layout
```
├── producer/              # High-throughput Go performance pipeline
│   └── main.go            # Goroutine architecture, channels, and signal logic
├── consumer/              # Scala event processor (SBT environment)
└── k8s/yamls/             # The Infrastructure Manifesto
├── kafka-ss.yml       # KRaft broker node StatefulSet scaling deployment
├── kafka-kraft.yml    # Controller node configurations
├── kafka-svc.yml      # Headless discovery network configuration
└── admin.yml          # Cluster role bindings and service permissions
```

## 🚀 Spin Up & Operations

### 1. Launch Cluster
Apply the infrastructure configuration natively to the cluster configuration tool:
```bash
kubectl apply -f k8s/yamls/
```

### 2. Run the Load Generator
Execute the Go harness from the command line, routing traffic directly to your ingress server or local host map:
```bash
go run producer/main.go -hostname=192.168.49.2:30000 -topic=t1
```
The console will stream aggregate metrics on worker offset shifts, target partition balancing, and transactional failures until a standard terminal termination signal (`Ctrl+C`) triggers a graceful lifecycle close down.
