const HDWalletProvider = require('truffle-hdwallet-provider');
const mnemonic = "when injury nature opinion october hair weekend pledge renew isolate fog giraffe";

module.exports = {
  networks: {
    development: {
      host: "127.0.0.1",
      port: 7545,
      network_id: "*"
    },
    ropsten: {
      provider: function() {
        return new HDWalletProvider(mnemonic, "https://ropsten.infura.io/v3/8a88efacc95d4b04a999672591c5feef")
      },
      network_id: 3
    }
  },
  solc: {
    optimizer: {
      enabled: true,
      runs: 200
    }
  }
};
