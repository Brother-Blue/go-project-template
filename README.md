A Go project template to get started without needing to setup and rewrite all the underlying goodness.

### Features include
- A centralized config, if you don't want something to run in the project, disable it there.
- Configurable telemetry, the project is preconfigured for logs, metrics, and tracing. All you need to do is instrument your code or use an auto-instrumenting tool.
- [WIP] Preconfigured Dockerfile and compose.
- [WIP] Preconfigured pipelines for GitHub CI, Jenkins, etc. to run tests, builds, deploys. These will most likely need to be reconfigured a bit to fit your use case. 
- [WIP] Compose configurations for self-hosted CIs and analyzers such as Sonarcube.
- [WIP] Automated TLS bootstrapping and management.