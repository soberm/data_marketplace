const UserContract = artifacts.require('UserContract');
const DeviceContract = artifacts.require('DeviceContract');
const ProductContract = artifacts.require('ProductContract');

module.exports = async function(callback) {
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
        let productContract = await ProductContract.deployed();

        await userContract.create(
            user.firstName, user.lastName, user.company, user.email, {from: accounts[0]}
        );

        for(let k = 1; k < 101; k++){
            let deviceID = web3.utils.randomHex(20);
            await deviceContract.create(
                deviceID, "Device", "asdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaa", web3.utils.hexToBytes(web3.utils.randomHex(64)), {from: accounts[0]}
            );

            for (let i = 0; i < k; i++) {
                await productContract.create(
                    deviceID, "Product", "asdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasasasaa", "light", 10, 5, {from: accounts[0]}
                );
            }

            let result = await deviceContract.remove(deviceID, {from: accounts[0]});

            console.log(k, ",", result.receipt.gasUsed)
        }

    }
    catch (e) {
        console.log(e)
    }
};
