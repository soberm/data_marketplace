const UserContract = artifacts.require('UserContract');
const truffleAssert = require('truffle-assertions');

contract("UserContractTest", async accounts => {

    let user = {
        firstName: 'michael',
        lastName: 'sober',
        company: 'company',
        email: 'michael.sober@ymail.com',
    };

    let contractInstance;

    beforeEach(function () {
        return UserContract.new()
            .then(function (instance) {
                contractInstance = instance;
            });
    });

    it("create success", async () => {
        let result = await contractInstance.create(
            user.firstName, user.lastName, user.company, user.email, user.image
        );
        truffleAssert.eventEmitted(result, 'CreatedUser');

        let foundUser = await contractInstance.findByAddress.call(accounts[0]);

        assert.equal(user.firstName, foundUser.firstName);
        assert.equal(user.lastName, foundUser.lastName);
        assert.equal(user.company, foundUser.company);
        assert.equal(user.email, foundUser.email);
    });

    it("create duplicate fails", async () => {
        await truffleAssert.passes(
            contractInstance.create(user.firstName, user.lastName, user.company, user.email)
        );
        await truffleAssert.fails(
            contractInstance.create(user.firstName, user.lastName, user.company, user.email),
            truffleAssert.ErrorType.REVERT,
            'user already exists'
        );
    });

});