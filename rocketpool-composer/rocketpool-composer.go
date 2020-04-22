package main

import(
    "errors"
    "fmt"
    "log"
    "os"
    "path/filepath"

    "github.com/urfave/cli"

    "github.com/rocket-pool/smartnode/rocketpool-composer/config"
    "github.com/rocket-pool/smartnode/rocketpool-composer/service"
)


// Run
func main() {

    // Add logo to application help template
    cli.AppHelpTemplate = fmt.Sprintf(`
______           _        _    ______           _ 
| ___ \         | |      | |   | ___ \         | |
| |_/ /___   ___| | _____| |_  | |_/ /__   ___ | |
|    // _ \ / __| |/ / _ \ __| |  __/ _ \ / _ \| |
| |\ \ (_) | (__|   <  __/ |_  | | | (_) | (_) | |
\_| \_\___/ \___|_|\_\___|\__| \_|  \___/ \___/|_|

%s`, cli.AppHelpTemplate)

    // Initialise application
    app := cli.NewApp()

    // Set application info
    app.Name = "rocketpool"
    app.Usage = "Rocket Pool CLI"
    app.Version = "0.0.1"
    app.Authors = []cli.Author{
        cli.Author{
            Name:  "David Rugendyke",
            Email: "david@rocketpool.net",
        },
        cli.Author{
            Name:  "Jake Pospischil",
            Email: "jake@rocketpool.net",
        },
    }
    app.Copyright = "(c) 2020 Rocket Pool Pty Ltd"

    // Register commands
            config.RegisterCommands(app, "config",  []string{"c"})
    service.RegisterServiceCommands(app, "service", []string{"s"})
        service.RegisterRunCommands(app, "run",     []string{"r"})

    // Check environment conditions before run
    app.Before = func(c *cli.Context) error {
        if err := checkEnv(); err != nil { log.Fatal(err) }
        return nil
    }

    // Run application
    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }

}


// Check environment conditions
func checkEnv() error {

    // Check RP_PATH environment variable
    rpPath := os.Getenv("RP_PATH")
    if rpPath == "" {
        return errors.New("The RP_PATH environment variable is not set. If you've just installed Rocket Pool, please start a new terminal session and try again.")
    }

    // Check RP_PATH exists
    if _, err := os.Stat(rpPath); os.IsNotExist(err) {
        return errors.New(fmt.Sprintf("The RP_PATH directory (%s) does not exist. Please create the directory and try again.", rpPath))
    }

    // Check config file exists
    configFilePath := filepath.Join(rpPath, "config.yml")
    if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
        return errors.New(fmt.Sprintf("The Rocket Pool configuration file (%s) does not exist. Please create the file and try again.", configFilePath))
    }

    // Check docker-compose file exists
    composeFilePath := filepath.Join(rpPath, "docker-compose.yml")
    if _, err := os.Stat(composeFilePath); os.IsNotExist(err) {
        return errors.New(fmt.Sprintf("The Rocket Pool docker-compose configuration file (%s) does not exist. Please create the file and try again.", composeFilePath))
    }

    // Return
    return nil

}
