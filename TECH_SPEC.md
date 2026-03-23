# Technical Specification: rmm23 (Remote Monitoring and Management)

## 1. Executive Summary

The `rmm23` system is a modern, high-performance Remote Monitoring and Management (RMM) platform designed for
centralized infrastructure control, granular user accounting, and automated network management. It emphasizes a flexible
data model, high-reliability backend storage, and multi-protocol frontend access (LDAP, Radius, API).

## 2. Core Architecture

### 2.1 Backend Storage (Redis)

The system utilizes Redis as its primary persistent and transient storage, partitioned by database index:

- **DB 0 (Main Data & Indexing):** Primary storage for UUID-based entities (Users, Devices, Groups, Domains).
- **DB 1 (Configuration):** Storage for dynamic configurations and template metadata.
- **DB 2 (MQTT/Messaging):** Real-time messaging and event distribution for remote monitoring.

### 2.2 Frontend Interfaces (AAA & API)

The application provides multiple entry points for external systems:

- **Built-in LDAP Frontend:** Provides directory services for user accounting and authentication.
- **Built-in RADIUS Frontend:** Handles network access requests for VPN and infrastructure equipment.
- **RESTful API:** Management interface for CRUD operations on entities and system monitoring.
- **SSH/API I/O:** Specialized modules for interacting with remote networking equipment (CSCO, JNPR).

### 2.3 Data Model & Identification

All infrastructure entities are uniquely identified through a hierarchical scheme:

- **UUID:** Primary unique identifier. For legacy data, null-based UUIDs are derived from Distinguished Names (DN).
- **Network Addressing:**
    - **ASN (Autonomous System Number):** `uint32` identifier (mapped to `uidNumber`).
    - **ipHostNumber:** Dedicated IPv4 subnet (typically `/27`) assigned per user/entity.
    - **Entity Byte Mapping (0x00-0x1F):** Specific offsets within subnets representing gateways, primary user devices (
      mobile, notebook, tablet), and broadcast addresses.

## 3. Key Functional Modules

### 3.1 Security & PKI (Public Key Infrastructure)

- **Certificate Management:** Automated building of certificate chains and local verification of trust.
- **ACME Integration:** Automated issuance and renewal of SSL/TLS certificates.
- **AAA Support:** Unified authentication, authorization, and accounting across LDAP, RADIUS, and SSH keys.

### 3.2 Network Management & ACLs

- **JunOS-style Policy Engine:** Implements security policies, policy options, and firewalls using a structure inspired
  by Juniper Networks' JunOS.
- **ACL Application Hierarchy:** `Infra` → `Domain` → `ACL-Groups` → `User`.
- **VPN Service Integration:** Automated configuration generation for OpenVPN, Cisco AnyConnect (CSCO), and Shadowsocks.

### 3.3 Automation & Configuration Engine

- **Template System:** Leverages Go `text/template` for generating service configurations.
- **Virtual File System (VFS):** Backend abstraction for configuration storage before deployment to the physical File
  System (FS).

## 4. Technical Standards & Quality of Development (QoD)

### 4.1 Modern Go Idioms

- **Language Version:** Go 1.26+.
- **Standard Library:** Mandatory use of modern packages (e.g., `slices`, `maps`, `cmp`).
- **Context Management:** Strict propagation of `context.Context` for cancellation and timeouts.
- **JSON Tags:** Use of `omitzero` for precise control over empty field serialization.

### 4.2 Advanced Error Handling (Error Arrays)

Instead of standard wrapped errors, `rmm23` implements an "Error Array" pattern:

1. **Multi-Error Returns:** Functions return `[]error` to capture multiple failure points.
2. **Severity Evaluation:** A dedicated checker evaluates the error array against Syslog-style severity levels.
3. **Flow Control:** The calling function proceeds based on the highest returned severity level rather than a simple
   boolean error check.

## 5. Deployment & Scalability

- **Daemon Operation:** Native support for running as a background service with cluster awareness.
- **Clustering Modes:**
    - **Slave:** Read-only/caching nodes for distributed frontend performance.
    - **Multi-master:** High-availability configuration for primary data synchronization.
- **Monitoring Integration:** Built-in support for Prometheus-style metrics and MQTT-based real-time telemetry.
