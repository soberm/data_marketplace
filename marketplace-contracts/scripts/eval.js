const UserContract = artifacts.require('UserContract');
const DeviceContract = artifacts.require('DeviceContract');
const BrokerContract = artifacts.require('BrokerContract');
const ProductContract = artifacts.require('ProductContract');
const NegotiationContract = artifacts.require('NegotiationContract');
const TradingContract = artifacts.require('TradingContract');

module.exports = async function(callback) {
    console.log("Creating test users");

    let user = {
        firstName: 'michael',
        lastName: 'sober',
        company: 'tu wien',
        email: 'e1326403@student.tuwien.ac.at',
    };

    try {
        let accounts = await web3.eth.getAccounts();

        let userContract = await UserContract.deployed();
        let deviceContract = await DeviceContract.deployed();
        let brokerContract = await BrokerContract.deployed();
        let productContract = await ProductContract.deployed();
        let negotiationContract = await NegotiationContract.deployed();
        let tradingContract = await TradingContract.deployed();

        let accountNumber = accounts.length / 10;
        for (let i = 0; i < accountNumber; i++) {
            let result = await userContract.create(
                user.firstName, user.lastName, user.company, user.email, {from: accounts[i]}
            );
            console.log("Created user ", accounts[i]);
            console.log("Gas: ", result.receipt.gasUsed)

            /*result = await userContract.update(
                user.firstName, user.lastName, user.company, user.email, {from: accounts[i]}
            );
            console.log("Updated user ", accounts[i]);
            console.log("Gas: ", result.receipt.gasUsed)*/
        }

        let brokerNumber = accountNumber*2;
        for (let i = accountNumber; i < brokerNumber; i++) {
            // let result = await brokerContract.create(
            //     accounts[i], "Broker", "127.0.0.1:25565", 0, {from: accounts[i-accountNumber]}
            // );
            // console.log("Created broker ", accounts[i]);
            // console.log("Gas: ", result.receipt.gasUsed);

/*            result = await brokerContract.remove(
                accounts[i], {from: accounts[i-accountNumber]}
            );
            console.log("Removed broker ", accounts[i]);
            console.log("Gas: ", result.receipt.gasUsed)*/

/*            result = await brokerContract.update(
                accounts[i], "Broker", "127.0.0.1:25565", 0, {from: accounts[i-accountNumber]}
            );
            console.log("Updated broker ", accounts[i]);
            console.log("Gas: ", result.receipt.gasUsed)*/
        }

        let j = 0;
        let productCount = 0;
        for (let i = brokerNumber; i < accounts.length; i++) {
            let result = await deviceContract.create(
                accounts[i], "Device", "asdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaa", web3.utils.hexToBytes(web3.utils.randomHex(64)), {from: accounts[j]}
            );
            console.log("Created Device ", accounts[i]);
            console.log("Gas: ", result.receipt.gasUsed);

            result = await deviceContract.remove(
                accounts[i], {from: accounts[j]}
            );
            console.log("Removed Device ", accounts[i]);
            console.log("Gas: ", result.receipt.gasUsed);

/*            result = await deviceContract.update(
                accounts[i], "Device", "asdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaa", web3.utils.hexToBytes(web3.utils.randomHex(64)), {from: accounts[j]}
            );
            console.log("Updated Device ", accounts[i]);
            console.log("Gas: ", result.receipt.gasUsed);*/

            for (let k = 0; k < 3; k++) {
                let result = await productContract.create(
                    accounts[i], "Product", "asdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaa", "light", 10, 5, {from: accounts[j]}
                );
                console.log("Created Product for ", accounts[i]);
                console.log("Gas: ", result.receipt.gasUsed);

                /*result = await productContract.update(
                    productCount, "Product", "asdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaaasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaa", "light", 10, 5, {from: accounts[j]}
                );
                console.log("Updated Product for ", accounts[i]);
                console.log("Gas: ", result.receipt.gasUsed)*/
                productCount++;
            }

            if (i % 10 === 0) {
                j++;
            }
        }

        /*for(let i = 0; i < productCount; i++) {
            let result = await productContract.remove(
                i, {from: accounts[0]}
            );
            console.log("Deleted Product for ", accounts[i]);
            console.log("Gas: ", result.receipt.gasUsed);
        }*/

        /*for (let k = 0; k < 10; k++) {
            let result = await tradingContract.requestTrading(
                0, accounts[brokerNumber-1], 1589265911, 1589265911, {from: accounts[brokerNumber+1]}
            );
            console.log("Created trading request");
            console.log("Gas: ", result.receipt.gasUsed);

            result = await tradingContract.acceptTradingRequest(k, {from: accounts[brokerNumber]});
            console.log("Accept trading request");
            console.log("Gas: ", result.receipt.gasUsed);
        }

        for (let k = 0; k < 10; k++) {
            let result = await negotiationContract.requestNegotiation(
                0, {from: accounts[brokerNumber+1]}
            );
            console.log("Created negotiation request");
            console.log("Gas: ", result.receipt.gasUsed);

            result = await negotiationContract.acceptNegotiationRequest(
                k, {from: accounts[brokerNumber]}
            );
            console.log("Accepted negotiation request");
            console.log("Gas: ", result.receipt.gasUsed);
        }*/

/*        let result = await userContract.remove({from: accounts[0]});
        console.log("Removed user ", accounts[0]);
        console.log("Gas: ", result.receipt.gasUsed)*/

    }
    catch (e) {
       console.log(e)
    }
};
