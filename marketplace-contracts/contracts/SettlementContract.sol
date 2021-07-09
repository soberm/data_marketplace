pragma solidity >=0.4.21 <0.6.0;

import "./RatingContract.sol";

contract SettlementContract {

    event Deposited(
        address indexed depositor,
        uint amount
    );

    event CounterSet(
        address indexed setter,
        uint counter
    );

    event Settled(
        uint actualCost,
        uint provider,
        uint consumer,
        uint broker
    );

    event Dispute(
        uint providerCounter,
        uint consumerCounter
    );

    struct Settlement {
        uint actualCost;
        uint provider;
        uint consumer;
        uint broker;
    }

    struct Counter {
        uint value;
        bool set;
    }

    uint constant timeout = 3600;

    address payable private provider;
    address payable private consumer;
    address payable private broker;

    uint private startTime;
    uint private endTime;
    uint private validTime;

    uint private cost;

    uint private maxCost;
    uint private maxCounter;

    Counter private providerCounter;
    Counter private consumerCounter;
    Counter private brokerCounter;

    Settlement private settlement;

    RatingContract private ratingContract;

    constructor(
        address _provider,
        address _consumer,
        address _broker,
        uint _startTime,
        uint _endTime,
        uint _cost,
        uint _frequency,
        address _ratingContractAddress
    )
    public
    {
        provider = address(uint160(_provider));
        consumer = address(uint160(_consumer));
        broker = address(uint160(_broker));
        startTime = _startTime;
        endTime = _endTime;
        cost = _cost;
        maxCounter = (endTime - startTime) / _frequency;
        maxCost = maxCounter * cost;
        validTime = endTime + timeout;
        ratingContract = RatingContract(_ratingContractAddress);
    }

    function deposit()
    public
    payable
    {
        require(block.timestamp < startTime, 'trade already started');
        require(msg.value >= maxCost, 'amount does not cover total cost');
        emit Deposited(msg.sender, msg.value);
    }

    function settleTrade(uint counter)
    public
    {
        require(block.timestamp >= endTime, 'trade not finished');
        require(counter <= maxCounter, 'counter is greater than max possible counter');

        setCounters(counter);

        if (settled()) {
            transferFunds(counter);
            ratingContract.handleSettlement(provider);
            ratingContract.handleSettlement(consumer);
        } else if (dispute()) {
            emit Dispute(providerCounter.value, consumerCounter.value);
        }
    }

    function resolveDispute(uint counter)
    public
    {
        require(dispute(), 'no dispute');
        require(msg.sender == broker, 'not authorized to resolve dispute');

        setBrokerCounter(counter);
        transferFunds(counter);

        if (providerCounter.value != brokerCounter.value) {
            ratingContract.handleDispute(provider);
        }

        if (consumerCounter.value != brokerCounter.value) {
            ratingContract.handleDispute(provider);
        }
    }

    function resolveTimeout()
    public
    {
        require(block.timestamp > validTime, 'no timeout');

        if (onlyProviderCounterSet()) {
            transferFunds(providerCounter.value);
        } else if (onlyConsumerCounterSet()) {
            transferFunds(consumerCounter.value);
        } else if (!countersSet()) {
            transferFunds(0);
        }
    }

    function transferFunds(uint counter)
    private
    {
        settlement.actualCost = counter * cost;
        settlement.broker = settlement.actualCost / 100;
        settlement.provider = settlement.actualCost - settlement.broker;
        settlement.consumer = address(this).balance - settlement.actualCost;

        provider.transfer(settlement.provider);
        consumer.transfer(settlement.consumer);
        broker.transfer(settlement.broker);

        emit Settled(settlement.actualCost, settlement.provider, settlement.consumer, settlement.broker);
    }

    function setCounters(uint counter)
    private
    {
        setProviderCounter(counter);
        setConsumerCounter(counter);
    }

    function setProviderCounter(uint counter)
    private
    {
        if (msg.sender == provider && !providerCounter.set) {
            providerCounter.value = counter;
            providerCounter.set = true;
            emit CounterSet(msg.sender, counter);
        }
    }

    function setConsumerCounter(uint counter)
    private
    {
        if (msg.sender == consumer && !consumerCounter.set) {
            consumerCounter.value = counter;
            consumerCounter.set = true;
            emit CounterSet(msg.sender, counter);
        }
    }

    function setBrokerCounter(uint counter)
    private
    {
        if (msg.sender == broker && !brokerCounter.set) {
            brokerCounter.value = counter;
            brokerCounter.set = true;
            emit CounterSet(msg.sender, counter);
        }
    }

    function getProviderCounter()
    public
    view
    returns (uint value, bool set)
    {
        value = providerCounter.value;
        set = providerCounter.set;
    }

    function getConsumerCounter()
    public
    view
    returns (uint value, bool set)
    {
        value = consumerCounter.value;
        set = consumerCounter.set;
    }

    function getBrokerCounter()
    public
    view
    returns (uint value, bool set)
    {
        value = brokerCounter.value;
        set = brokerCounter.set;
    }

    function getSettlement()
    public
    view
    returns (uint _actualCost, uint _provider, uint _consumer, uint _broker)
    {
        _actualCost = settlement.actualCost;
        _provider = settlement.provider;
        _consumer = settlement.consumer;
        _broker = settlement.broker;
    }

    function settled()
    public
    view
    returns (bool _settled)
    {
        return countersSet() && countersMatch();
    }

    function dispute()
    public
    view
    returns (bool _dispute)
    {
        return countersSet() && !countersMatch();
    }

    function countersSet()
    private
    view
    returns (bool set)
    {
        return providerCounter.set && consumerCounter.set;
    }

    function onlyProviderCounterSet()
    private
    view
    returns (bool set)
    {
        return providerCounter.set && !consumerCounter.set;
    }

    function onlyConsumerCounterSet()
    private
    view
    returns (bool set)
    {
        return !providerCounter.set && consumerCounter.set;
    }

    function countersMatch()
    private
    view
    returns (bool matched)
    {
        return providerCounter.value == consumerCounter.value;
    }

}
