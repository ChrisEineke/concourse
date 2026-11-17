# Concourse

Concourse is a libp2p-based and delegation-oriented framework for distributed
deployment of workloads.

## Overview

* Concourse Clusters have *no* centrally-managed control plane.
* A single distributed work queue holds instructions for workload deployment.
* Concourse Nodes cooperatively schedule workloads.

## Framwork Modules

### Concourse Node Orchestrator

The Concourse Node Orchestrator bootstraps nodes to be used in a Concourse
Cluster and establishes a close-to-optimal network topology.


