package main

import (
    "context"
    "log"
    "os"

    vault "github.com/hashicorp/vault/api"
)

const password string = "projat"

func main() {
    config := vault.DefaultConfig()
    config.Address = os.Getenv("VAULT_ADDR")

    client, err := vault.NewClient(config)
    if err != nil {
        log.Fatalf("Unable to initialize a Vault client: %v", err)
    }

    client.SetToken(os.Getenv("VAULT_TOKEN"))

    secretData := map[string]interface{}{
        "password": password,
    }

    ctx := context.Background()
    
    _, err = client.KVv2("secret").Put(ctx, "my-secret-password", secretData)
if err != nil {
    log.Fatalf("Unable to write secret: %v to the vault", err)
}
log.Println("Super secret password written successfully to the vault.")

versions, err := client.KVv2("secret").GetVersionsAsList(ctx, "my-secret-password")
if err != nil {
    log.Fatalf(
        "Unable to retrieve all versions of the super secret password from the vault. Reason: %v", 
        err,
    )
}

for _, version := range versions {
    deleted := "Not deleted"
    if !version.DeletionTime.IsZero() {
        deleted = version.DeletionTime.Format("HH:MM:SS")
    }

    secret, err := client.KVv2("secret").
        GetVersion(ctx, "my-secret-password", version.Version)
    if err != nil {
        log.Fatalf(
            "Unable to retrieve version %d of the super secret password from the vault. Reason: %v",
            err,
        )
    }
    value, ok := secret.Data["password"].(string)

    if ok {
        log.Printf(
            "Version: %d. Deleted at: %s. Destroyed: %t. Value: '%s'.\n",
            version.Version,
            //version.CreatedTime.Format("HH:MM:SS"),
            deleted,
            version.Destroyed,
            value,
       )
    }
}

}

