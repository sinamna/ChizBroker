# Welcome project for newcomers!

# Introduction
In this project you have to implement a message broker, based on `broker.Broker`
interface. There are unit tests to specify requirements and also validate your implementation.

# Roadmap
- [ ] Implement `broker.Broker` interface and pass all tests
- [ ] Add basic logs and prometheus metrics
  - Metrics for each RPCs:
    - `method_count` to show count of failed/successful RPC calls
    - `method_duration` for latency of each call, in 99, 95, 50 quantiles
    - `active_subscribers` to display total active subscriptions
  - Env metrics:
    - Metrics for your application memory, cpu load, cpu utilization, GCs
- [ ] Implement gRPC API for the broker and main functionalities
- [ ] Create *dockerfile* and *docker-compose* files for your deployment
- [ ] Deploy your app with the previous `docker-compose` on a remote machine
- [ ] Deploy your app on K8

# Phase 2 Evaluation
We run our gRPC client that implemented the `broker.proto` against your deployed broker application.

As it should function properly ( like the unit tests ), we expect the provided metrics to display a good observation, and if
anything unexpected happened, you could diagnose your app, using the logs and other tools.