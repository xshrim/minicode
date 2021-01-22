# DEPRECATION NOTICE

This chart is **DEPRECATED**.

### Replacement

We have built a set of fully cloud native charts in [gitlab/gitlab](https://gitlab.com/charts/gitlab).
These new charts are designed from the ground up to be performant, flexible, scalable, and resilient.

We _very strongly_ recommend transitioning, if you are currently using these charts. If you have
never used these charts, _do not now_.

### Availability

This project remains visible as an example of how to convert a full monolith application to Kubernetes capable.
[Monolith to Microservice: Pitchforks not included](https://youtu.be/rIUth_KrJdw?list=PLj6h78yzYM2PZf9eA7bhWnIh_mK1vyOfU) (video)
details the work done to break this monolithic container into component parts.

# GitLab-Omnibus Helm Chart

This chart is an easy way to get started with GitLab on Kubernetes. It includes everything needed to run GitLab, including: a Runner, Container Registry, automatic SSL, and an Ingress.

For more information, please review [our documentation](http://docs.gitlab.com/ee/install/kubernetes/gitlab_omnibus.html).

