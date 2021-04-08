package config

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"

	"github.com/imdario/mergo"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"

	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

// Rocket Pool config
type RocketPoolConfig struct {
    Rocketpool struct {
        StorageAddress string           `yaml:"storageAddress,omitempty"`
        RPLFaucetAddress string         `yaml:"rplFaucetAddress,omitempty"`
    }                                   `yaml:"rocketpool,omitempty"`
    Smartnode struct {
        ProjectName string              `yaml:"projectName,omitempty"`
        GraffitiVersion string          `yaml:"graffitiVersion,omitempty"`
        Image string                    `yaml:"image,omitempty"`
        PasswordPath string             `yaml:"passwordPath,omitempty"`
        WalletPath string               `yaml:"walletPath,omitempty"`
        ValidatorKeychainPath string    `yaml:"validatorKeychainPath,omitempty"`
        ValidatorRestartCommand string  `yaml:"validatorRestartCommand,omitempty"`
        GasPrice string                 `yaml:"gasPrice,omitempty"`
        GasLimit string                 `yaml:"gasLimit,omitempty"`
    }                                   `yaml:"smartnode,omitempty"`
    Chains struct {
        Eth1 Chain                      `yaml:"eth1,omitempty"`
        Eth2 Chain                      `yaml:"eth2,omitempty"`
    }                                   `yaml:"chains,omitempty"`
}
type Chain struct {
    Provider string                     `yaml:"provider,omitempty"`
    WsProvider string                   `yaml:"wsProvider,omitempty"`
    ChainID string                      `yaml:"chainID,omitempty"`
    Client struct {
        Options []ClientOption          `yaml:"options,omitempty"`
        Selected string                 `yaml:"selected,omitempty"`
        Params []UserParam              `yaml:"params,omitempty"`
    }                                   `yaml:"client,omitempty"`
}
type ClientOption struct {
    ID string                           `yaml:"id,omitempty"`
    Name string                         `yaml:"name,omitempty"`
    Desc string                         `yaml:"desc,omitempty"`
    Image string                        `yaml:"image,omitempty"`
    BeaconImage string                  `yaml:"beaconImage,omitempty"`
    ValidatorImage string               `yaml:"validatorImage,omitempty"`
    Link string                         `yaml:"link,omitempty"`
    Params []ClientParam                `yaml:"params,omitempty"`
}
type ClientParam struct {
    Name string                         `yaml:"name,omitempty"`
    Desc string                         `yaml:"desc,omitempty"`
    Env string                          `yaml:"env,omitempty"`
    Required bool                       `yaml:"required,omitempty"`
    Regex string                        `yaml:"regex,omitempty"`
    Type string                         `yaml:"type,omitempty"`
    Default string                      `yaml:"default,omitempty"`
    Max string                          `yaml:"max,omitempty"`
    BlankText string                    `yaml:"blankText,omitempty"`
}
type UserParam struct {
    Env string                          `yaml:"env,omitempty"`
    Value string                        `yaml:"value"`
}


// Get the selected clients from a config
func (config *RocketPoolConfig) GetSelectedEth1Client() *ClientOption {
    return config.Chains.Eth1.GetSelectedClient()
}
func (config *RocketPoolConfig) GetSelectedEth2Client() *ClientOption {
    return config.Chains.Eth2.GetSelectedClient()
}
func (chain *Chain) GetSelectedClient() *ClientOption {
    for _, option := range chain.Client.Options {
        if option.ID == chain.Client.Selected {
            return &option
        }
    }
    return nil
}


// Get the beacon & validator images for a client
func (client *ClientOption) GetBeaconImage() string {
    if client.BeaconImage != "" {
        return client.BeaconImage
    } else {
        return client.Image
    }
}
func (client *ClientOption) GetValidatorImage() string {
    if client.ValidatorImage != "" {
        return client.ValidatorImage
    } else {
        return client.Image
    }
}


// Serialize a config to yaml bytes
func (config *RocketPoolConfig) Serialize() ([]byte, error) {
    bytes, err := yaml.Marshal(config)
    if err != nil {
        return []byte{}, fmt.Errorf("Could not serialize config: %w", err)
    }
    return bytes, nil
}


// Parse a config from yaml bytes
func Parse(bytes []byte) (RocketPoolConfig, error) {
    var config RocketPoolConfig
    if err := yaml.Unmarshal(bytes, &config); err != nil {
        return RocketPoolConfig{}, fmt.Errorf("Could not parse config: %w", err)
    }

    // Validate the defaults
    if err := ValidateDefaults(config.Chains.Eth1, "eth1"); err != nil {
        return RocketPoolConfig{}, err
    }
    if err := ValidateDefaults(config.Chains.Eth2, "eth2"); err != nil {
        return RocketPoolConfig{}, err
    }

    return config, nil
}


// Make sure the default parameter values can be parsed into the parameter types
func ValidateDefaults(Chain Chain, ChainName string) (error) {
    for _, option := range Chain.Client.Options {
        for _, param := range option.Params {
            if param.Default != "" {
                var err error

                switch param.Type {
                case "", "string":
                    continue

                case "uint":
                    _, err = strconv.ParseUint(param.Default, 0, 0)

                case "uint16":
                    _, err = strconv.ParseUint(param.Default, 0, 16)
                }

                if err != nil {
                    return fmt.Errorf("Could not parse config - " +
                        "parameter '%s' in %s client option '%s' " +
                        "is a %s but has a default value of '%s' which failed parsing: %w", 
                        param.Name, ChainName, option.Name, param.Type, param.Default, err)
                }
            }
        }
    }

    return nil
}


// Merge configs
func Merge(configs ...*RocketPoolConfig) RocketPoolConfig {
    var merged RocketPoolConfig
    for i := len(configs) - 1; i >= 0; i-- {
        mergo.Merge(&merged, configs[i])
    }
    return merged
}


// Load merged config from files
func Load(c *cli.Context) (RocketPoolConfig, error) {

    // Load configs
    globalConfig, err := loadFile(os.ExpandEnv(c.GlobalString("config")), true)
    if err != nil {
        return RocketPoolConfig{}, err
    }
    userConfig, err := loadFile(os.ExpandEnv(c.GlobalString("settings")), false)
    if err != nil {
        return RocketPoolConfig{}, err
    }
    cliConfig := getCliConfig(c)

    // Merge and return
    return Merge(&globalConfig, &userConfig, &cliConfig), nil

}


// Load config from a file
func loadFile(path string, required bool) (RocketPoolConfig, error) {

    // Read file; squelch not found errors if file is optional
    bytes, err := ioutil.ReadFile(path)
    if err != nil {
        if required {
            return RocketPoolConfig{}, fmt.Errorf("Could not find config file at %s: %w", path, err)
        } else {
            return RocketPoolConfig{}, nil
        }
    }

    // Parse config
    var config RocketPoolConfig
    if err := yaml.Unmarshal(bytes, &config); err != nil {
        return RocketPoolConfig{}, fmt.Errorf("Could not parse config file at %s: %w", path, err)
    }

    // Return
    return config, nil

}


// Create config from CLI arguments
func getCliConfig(c *cli.Context) RocketPoolConfig {
    var config RocketPoolConfig
    config.Rocketpool.StorageAddress = c.GlobalString("storageAddress")
    config.Rocketpool.RPLFaucetAddress = c.GlobalString("rplFaucetAddress")
    config.Smartnode.PasswordPath = c.GlobalString("password")
    config.Smartnode.WalletPath = c.GlobalString("wallet")
    config.Smartnode.ValidatorKeychainPath = c.GlobalString("validatorKeychain")
    config.Smartnode.GasPrice = c.GlobalString("gasPrice")
    config.Smartnode.GasLimit = c.GlobalString("gasLimit")
    config.Chains.Eth1.Provider = c.GlobalString("eth1Provider")
    config.Chains.Eth2.Provider = c.GlobalString("eth2Provider")
    return config
}


// Parse and return the gas price in wei
func (config *RocketPoolConfig) GetGasPrice() (*big.Int, error) {

    // No gas price specified
    if config.Smartnode.GasPrice == "" {
        return nil, nil
    }

    // Parse gas price in gwei
    gasPriceGwei, err := strconv.ParseFloat(config.Smartnode.GasPrice, 64)
    if err != nil {
        return nil, fmt.Errorf("Invalid gas price '%s': %w", config.Smartnode.GasPrice, err)
    }

    // Return nil if gas price is set to zero
    if gasPriceGwei == 0 {
        return nil, nil
    }

    // Return gas price in wei
    return eth.GweiToWei(gasPriceGwei), nil

}


// Parse and return the gas limit
func (config *RocketPoolConfig) GetGasLimit() (uint64, error) {

    // No gas limit specified
    if config.Smartnode.GasLimit == "" {
        return 0, nil
    }

    // Parse gas limit
    gasLimit, err := strconv.ParseUint(config.Smartnode.GasLimit, 10, 64)
    if err != nil {
        return 0, fmt.Errorf("Invalid gas limit '%s': %w", config.Smartnode.GasLimit, err)
    }

    // Return
    return gasLimit, nil

}

