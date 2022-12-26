# HashicorpVault
Vaullting newly generated certificates

Download vault:

    $ wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor | sudo tee /usr/share/keyrings/hashicorp-archive-keyring.gpg
    
    $ echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
    
    $ sudo apt update && sudo apt install vault

set the path to vault bin


Before writing any code, to get you up and running quickly, start Vault in “Dev” Server mode, by running the following command.

    $ vault server -dev

The command will write output similar to the following to the terminal.
You may need to set the following environment variable:

 $ export VAULT_ADDR='http://127.0.0.1:8200'

The unseal key and root token are displayed below in case you want to
seal/unseal the Vault or re-authenticate.

Unseal Key: OqoUeFDlv9PjqvJgBolLyevsj4y3gqPInNKvBubZTd0=
Root Token: hvs.2fWa1QeRWesGfjGeb2QqBYU4

Development mode should NOT be used in production installations!


From the output written to your terminal, copy the Root Token value written towards the end (e.g., hvs.2fWa1QeRWesGfjGeb2QqBYU4) and paste it in place of the placeholder <VAULT_TOKEN> in the second command below. 

Then, run the commands in a new terminal window.

    $ export VAULT_ADDR='http://127.0.0.1:8200'
    
    $ export VAULT_TOKEN=<VAULT_TOKEN>
    

Then setting the path to GO-PATH
export PATH=/usr/local/go/bin:$PATH


so first create and store the password:

    $ go run create.go
    
once done try to read and retrieve the secret

    $ go run read.go
    
 Now if you have created more than 1 secrets and you want to check, you can :
 
    $ go run readAll.go
    
 Now let us delete the latest saved secret
 
    $ go run delete.go
    
  Now let us delete all the saved secrets
  
    $ go run deleteAll.go
    
Closing the terminal in which the vault server is running (vault server -dev) will stop the running vault.

# Consul

Install Consul:

Add the HashiCorp GPG key.

    curl --fail --silent --show-error --location https://apt.releases.hashicorp.com/gpg | \
      gpg --dearmor | \
      sudo dd of=/usr/share/keyrings/hashicorp-archive-keyring.gpg

Add the official HashiCorp Linux repository.
    
    echo "deb [arch=amd64 signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | \
 sudo tee -a /etc/apt/sources.list.d/hashicorp.list

Update.

    sudo apt-get update
    
Find available versions to install

    sudo apt-cache policy consul

Install latest version.

    sudo apt-get install consul
    
Verify the installation

    consul
    
If you get an error that consul could not be found, your PATH environment variable was not set up properly. Make sure that your PATH variable contains the directory where you installed Consul.  

Generate the gossip encryption key

Gossip is encrypted with a symmetric key, since gossip between nodes is done over UDP. All agents must have the same encryption key.

You can create the encryption key via the Consul CLI even though no Consul agents are running yet. Generate the encryption key.

    consul keygen
    
Copy the output and keep it handy with you.

Generate TLS certificates for RPC encryption

Consul can use TLS to verify the authenticity of servers and clients. To enable TLS, Consul requires that all agents have certificates signed by a single Certificate Authority (CA).

Start by creating the CA on your admin instance, using the Consul CLI.

    consul tls ca create
    
Create the certificates

Next create a set of certificates, one for each Consul agent. You will need to select a name for your primary datacenter now, so that the certificates are named properly, as well as a domain for your Consul datacenter.

First, for your Consul servers, use the following command to create a certificate for each server. The file name increments automatically.

Option 1. For Single Consul Datacenter:

    consul tls cert create -server -dc dc1 -domain consul
   
 Option 2. For Federated Consul Datacenter
 
    consul tls cert create -server -dc dc1 -domain consul -additional-dnsname=*.dc2.consul
    
 Distribute the certificates to agents

You must distribute the CA certificate, consul-agent-ca.pem, to each of the Consul agents as well as the agent specific certificate and private key.

Below is an SCP example which will send the CA certificate, agent certificate and private key to the IP address you specify, and put it into the /etc/consul.d/ directory.

Following the earlier option 1

    sudo cp dc1-server-consul-0-key.pem /etc/consul.d/
    
    sudo cp consul-agent-ca.pem /etc/consul.d/
    
    sudo cp dc1-server-consul-0.pem /etc/consul.d/
    
    sudo cp dc1-server-consul-0-key.pem /etc/consul.d/
 
Consul server agents typically require a superset of configuration required by Consul client agents. You will specify common configuration used by all Consul agents in consul.hcl and server specific configuration in server.hcl.

Create a configuration file at /etc/consul.d/consul.hcl:

    sudo mkdir --parents /etc/consul.d
    
    sudo touch /etc/consul.d/consul.hcl
    
    sudo chown --recursive consul:consul /etc/consul.d
    
    sudo chmod 640 /etc/consul.d/consul.hcl
    
Add this configuration to the consul.hcl configuration file:

    # Full configuration options can be found at https://www.consul.io/docs/agent/config
    # datacenter
    # This flag controls the datacenter in which the agent is running. If not provided,
    # it defaults to "dc1". Consul has first-class support for multiple datacenters, but 
    # it relies on proper configuration. Nodes in the same datacenter should be on a 
    # single LAN.
    #datacenter = "my-dc-1"
    # data_dir
    # This flag provides a data directory for the agent to store state. This is required
    # for all agents. The directory should be durable across reboots. This is especially
    # critical for agents that are running in server mode as they must be able to persist
    # cluster state. Additionally, the directory must support the use of filesystem
    # locking, meaning some types of mounted folders (e.g. VirtualBox shared folders) may
    # not be suitable.
    data_dir = "/opt/consul"
    # client_addr
    # The address to which Consul will bind client interfaces, including the HTTP and DNS
    # servers. By default, this is "127.0.0.1", allowing only loopback connections. In
    # Consul 1.0 and later this can be set to a space-separated list of addresses to bind
    # to, or a go-sockaddr template that can potentially resolve to multiple addresses.
    #client_addr = "0.0.0.0"
    # ui
    # Enables the built-in web UI server and the required HTTP routes. This eliminates
    # the need to maintain the Consul web UI files separately from the binary.
    # Version 1.10 deprecated ui=true in favor of ui_config.enabled=true
    #ui_config{
    #  enabled = true
    #}
    # server
    # This flag is used to control if an agent is in server or client mode. When provided,
    # an agent will act as a Consul server. Each Consul cluster must have at least one
    # server and ideally no more than 5 per datacenter. All servers participate in the Raft
    # consensus algorithm to ensure that transactions occur in a consistent, linearizable
    # manner. Transactions modify cluster state, which is maintained on all server nodes to
    # ensure availability in the case of node failure. Server nodes also participate in a
    # WAN gossip pool with server nodes in other datacenters. Servers act as gateways to
    # other datacenters and forward traffic as appropriate.
    #server = true
    # Bind addr
    # You may use IPv4 or IPv6 but if you have multiple interfaces you must be explicit.
    #bind_addr = "[::]" # Listen on all IPv6
    #bind_addr = "0.0.0.0" # Listen on all IPv4
    #
    # Advertise addr - if you want to point clients to a different address than bind or LB.
    #advertise_addr = "127.0.0.1"
    # Enterprise License
    # As of 1.10, Enterprise requires a license_path and does not have a short trial.
    #license_path = "/etc/consul.d/consul.hclic"
    # bootstrap_expect
    # This flag provides the number of expected servers in the datacenter. Either this value
    # should not be provided or the value must agree with other servers in the cluster. When
    # provided, Consul waits until the specified number of servers are available and then
    # bootstraps the cluster. This allows an initial leader to be elected automatically.
    # This cannot be used in conjunction with the legacy -bootstrap flag. This flag requires
    # -server mode.
    #bootstrap_expect=3
    # encrypt
    # Specifies the secret key to use for encryption of Consul network traffic. This key must
    # be 32-bytes that are Base64-encoded. The easiest way to create an encryption key is to
    # use consul keygen. All nodes within a cluster must share the same encryption key to
    # communicate. The provided key is automatically persisted to the data directory and loaded
    # automatically whenever the agent is restarted. This means that to encrypt Consul's gossip
    # protocol, this option only needs to be provided once on each agent's initial startup
    # sequence. If it is provided after Consul has been initialized with an encryption key,
    # then the provided key is ignored and a warning will be displayed.
    #encrypt = "..."
    # retry_join
    # Similar to -join but allows retrying a join until it is successful. Once it joins 
    # successfully to a member in a list of members it will never attempt to join again.
    # Agents will then solely maintain their membership via gossip. This is useful for
    # cases where you know the address will eventually be available. This option can be
    # specified multiple times to specify multiple agents to join. The value can contain
    # IPv4, IPv6, or DNS addresses. In Consul 1.1.0 and later this can be set to a go-sockaddr
    # template. If Consul is running on the non-default Serf LAN port, this must be specified
    # as well. IPv6 must use the "bracketed" syntax. If multiple values are given, they are
    # tried and retried in the order listed until the first succeeds. Here are some examples:
    #retry_join = ["consul.domain.internal"]
    #retry_join = ["10.0.4.67"]
    #retry_join = ["[::1]:8301"]
    #retry_join = ["consul.domain.internal", "10.0.4.67"]
    # Cloud Auto-join examples:
    # More details - https://www.consul.io/docs/agent/cloud-auto-join
    #retry_join = ["provider=aws tag_key=... tag_value=..."]
    #retry_join = ["provider=azure tag_name=... tag_value=... tenant_id=... client_id=... subscription_id=... secret_access_key=..."]
    #retry_join = ["provider=gce project_name=... tag_value=..."]
    datacenter = "dc1"
    encrypt = "ebjTffRGKzWWiae87v9wSUmjuLhbDzkkts+ccF0OWY0="
    verify_incoming = true
    verify_outgoing = true
    verify_server_hostname = true
    ca_file = "/home/hp/vault-go/consul-agent-ca.pem"
    cert_file = "/home/hp/vault-go/dc1-server-consul-0.pem"
    key_file = "/home/hp/vault-go/dc1-server-consul-0-key.pem"
    auto_encrypt {
        allow_tls = true
    }
    acl {
        enabled = true
        default_policy = "allow"
        enable_token_persistence = true
    }
    audit {
        enabled = true
        sink "My sink" {
            type   = "file"
            format = "json"
            path   = "data/audit/audit.json"
            delivery_guarantee = "best-effort"
            rotate_duration = "24h"
            rotate_max_files = 15
            rotate_bytes = 25165824
            }
        }
