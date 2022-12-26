# HashicorpVault
Vaullting newly generated certificates

Download vault:
wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor | sudo tee /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt update && sudo apt install vault

set the path to vault bin


Before writing any code, to get you up and running quickly, start Vault in “Dev” Server mode, by running the following command.
vault server -dev

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
go run create.go
