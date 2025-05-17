# logistics-platform

## 🚀 What is it?

This is a simplified version of a logistics platform developed by me as part of the **Advanced Go Development** course by [Ozon Tech](https://ozon.tech/?__rr=1&abt_att=1&origin_referer=www.google.com).

The project demonstrates clean architecture, modular design, and interaction between multiple microservices written in Go.

## 🧑‍💻 Code Ownership

**95% of the codebase was written entirely by me** — including the overall architecture, domain logic, service interactions (via gRPC, Kafka), and infrastructure setup. The remaining 5% consists of preconfigured CI wrappers provided by the course (e.g., host URLs and service ports).

## 🧱 Project Structure

This monorepo consists of the following services:

- `cart` – shopping cart service
- `loms` – logistics order management system
- `notifier` – notification service
- `comments` – product comment system
- `product-service` – product catalog and management

> [TODO]: Add detailed description for each service.

## 🛠 Technologies Used

- Go  
- gRPC / HTTP  
- PostgreSQL  
- Kafka  
- Docker  
- Prometheus (metrics collection)  
- Jaeger (distributed tracing)  
- Grafana (visualization and monitoring)  
- [planned CI/CD] GitHub Actions 

Feel free to explore the code and follow the progress as I refine the project and CI workflows.
