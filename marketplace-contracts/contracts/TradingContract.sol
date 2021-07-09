pragma solidity >=0.4.21 <0.6.0;

import "./BiddingContract.sol";
import "./DeviceContract.sol";
import "./ProductContract.sol";
import "./BrokerContract.sol";
import "./SettlementContract.sol";
import "./NegotiationContract.sol";

contract TradingContract {

    event RequestedTrading(
        uint id,
        address indexed requester,
        uint indexed product
    );

    event AcceptedTradingRequest(
        uint indexed id,
        uint trade
    );

    event CreatedTrade(
        address indexed broker,
        uint tradeId
    );

    event DeclinedTradingRequest(
        uint indexed id
    );

    struct TradingRequest {
        uint id;
        uint index;
        uint product;
        uint cost;
        uint frequency;
        uint startTime;
        uint endTime;
        address consumer;
        address broker;
    }

    struct Trade {
        uint id;
        uint index;
        address provider;
        address consumer;
        address broker;
        uint product;
        uint startTime;
        uint endTime;
        uint cost;
        uint frequency;
        address settlementContract;
    }

    mapping(uint => TradingRequest) private tradingRequests;
    uint[] private tradingRequestsIndex;

    mapping(uint => Trade) private trades;
    uint[] private tradesIndex;

    uint private tradeIdCounter;
    uint private requestIdCounter;

    uint private constant timeout = 86400;

    DeviceContract private deviceContract;
    ProductContract private productContract;
    BrokerContract private brokerContract;
    NegotiationContract private negotiationContract;
    RatingContract private ratingContract;

    constructor(
        address _deviceContractAddress,
        address _productContractAddress,
        address _brokerContractAddress,
        address _negotiationContractAddress,
        address _ratingContractAddress
    )
    public
    {
        deviceContract = DeviceContract(_deviceContractAddress);
        productContract = ProductContract(_productContractAddress);
        brokerContract = BrokerContract(_brokerContractAddress);
        negotiationContract = NegotiationContract(_negotiationContractAddress);
        ratingContract = RatingContract(_ratingContractAddress);
    }

    function requestTrading(uint _product, address _broker, uint _startTime, uint _endTime)
    external
    returns (uint id)
    {
        require(deviceContract.existsByAddressAndDeleted(msg.sender, false), 'consumer does not exist');
        require(brokerContract.existsByAddressAndDeleted(_broker, false), 'broker does not exist');
        require(productContract.existsByIdAndDeleted(_product, false), 'product does not exist');

        uint nextId = requestIdCounter++;
        TradingRequest storage tradingRequest = tradingRequests[nextId];

        tradingRequest.id = nextId;
        tradingRequest.product = _product;
        tradingRequest.cost = productContract.findCostById(_product);
        tradingRequest.frequency = productContract.findFrequencyById(_product);
        tradingRequest.consumer = msg.sender;
        tradingRequest.broker = _broker;
        tradingRequest.startTime = _startTime;
        tradingRequest.endTime = _endTime;
        tradingRequest.index = tradingRequestsIndex.push(tradingRequest.id) - 1;

        deviceContract.addTradingRequest(tradingRequest.consumer, tradingRequest.id);
        productContract.addTradingRequest(tradingRequest.product, tradingRequest.id);

        emit RequestedTrading(tradingRequest.id, tradingRequest.consumer, tradingRequest.product);
        return tradingRequest.id;
    }

    function acceptTradingRequest(uint _id)
    external
    returns (uint id)
    {
        require(containsTradingRequest(_id), 'request does not exist');

        TradingRequest memory tradingRequest = tradingRequests[_id];

        require(
            productContract.isOwnedByDevice(tradingRequest.product, msg.sender),
            'not authorized to accept request'
        );

        uint tradeId = create(
            tradingRequest.consumer,
            tradingRequest.broker,
            tradingRequest.product,
            tradingRequest.startTime,
            tradingRequest.endTime,
            tradingRequest.cost,
            tradingRequest.frequency
        );

        removeTradingRequest(_id);
        emit AcceptedTradingRequest(tradingRequest.id, tradeId);
        return tradingRequest.id;
    }

    function declineTradingRequest(uint _id)
    external
    returns (uint id)
    {
        require(containsTradingRequest(_id), 'request does not exist');

        TradingRequest memory tradingRequest = tradingRequests[_id];

        require(
            productContract.isOwnedByDevice(tradingRequest.product, msg.sender),
            'not authorized to decline request'
        );

        removeTradingRequest(_id);
        emit DeclinedTradingRequest(_id);
        return _id;
    }

    function create(uint _negotiation, address _broker)
    external
    returns (uint id)
    {
        require(negotiationContract.existsNegotiationById(_negotiation), 'negotiation does not exist');
        require(brokerContract.existsByAddressAndDeleted(_broker, false), 'broker does not exist');

        (, address consumer, uint product, address biddingContractAddr) =
        negotiationContract.findNegotiationById(_negotiation);

        require(
            productContract.isOwnedByDevice(product, msg.sender),
            'not authorized to create trade'
        );

        BiddingContract biddingContract = BiddingContract(biddingContractAddr);

        require(biddingContract.isAccepted(), 'no agreement reached in bidding contract');

        (uint cost, uint startTime, uint endTime) = biddingContract.lastBid();

        return create(consumer, _broker, product, startTime, endTime, cost, productContract.findFrequencyById(product));
    }

    function findTradingRequestByIndex(uint _index)
    external
    view
    returns (
        uint id,
        uint product,
        uint cost,
        uint startTime,
        uint endTime,
        address consumer,
        address broker
    )
    {
        require(_index >= 0 && _index < tradingRequestsIndex.length, 'request does not exist');

        TradingRequest memory tradingRequest = tradingRequests[tradingRequestsIndex[_index]];

        id = tradingRequest.id;
        product = tradingRequest.product;
        cost = tradingRequest.cost;
        startTime = tradingRequest.startTime;
        endTime = tradingRequest.endTime;
        consumer = tradingRequest.consumer;
        broker = tradingRequest.broker;
    }

    function findTradeByIndex(uint _index)
    external
    view
    returns (
        uint id,
        address provider,
        address consumer,
        address broker,
        uint product,
        uint startTime,
        uint endTime,
        uint cost,
        address settlementContract
    )
    {
        require(_index >= 0 && _index < tradesIndex.length, 'trade does not exist');

        Trade memory trade = trades[tradesIndex[_index]];

        id = trade.id;
        provider = trade.provider;
        consumer = trade.consumer;
        broker = trade.broker;
        product = trade.product;
        startTime = trade.startTime;
        endTime = trade.endTime;
        cost = trade.cost;
        settlementContract = trade.settlementContract;
    }

    function findTradingRequestById(uint _id)
    external
    view
    returns (
        uint id,
        uint product,
        uint cost,
        uint startTime,
        uint endTime,
        address consumer,
        address broker
    )
    {
        require(containsTradingRequest(_id), 'request does not exist');

        TradingRequest memory tradingRequest = tradingRequests[_id];

        id = tradingRequest.id;
        product = tradingRequest.product;
        cost = tradingRequest.cost;
        startTime = tradingRequest.startTime;
        endTime = tradingRequest.endTime;
        consumer = tradingRequest.consumer;
        broker = tradingRequest.broker;
    }

    function findTradeById(uint _id)
    external
    view
    returns (
        uint id,
        address provider,
        address consumer,
        address broker,
        uint product,
        uint startTime,
        uint endTime,
        uint cost,
        address settlementContract
    )
    {
        require(containsTrade(_id), 'trade does not exist');

        Trade memory trade = trades[_id];

        id = trade.id;
        provider = trade.provider;
        consumer = trade.consumer;
        broker = trade.broker;
        product = trade.product;
        startTime = trade.startTime;
        endTime = trade.endTime;
        cost = trade.cost;
        settlementContract = trade.settlementContract;
    }

    function countTradingRequests()
    external
    view
    returns (uint counter)
    {
        return tradingRequestsIndex.length;
    }

    function countTrades()
    external
    view
    returns (uint counter)
    {
        return tradesIndex.length;
    }

    function containsTradingRequest(uint _id)
    public
    view
    returns (bool isContained)
    {
        if (tradingRequestsIndex.length == 0) return false;
        return (tradingRequestsIndex[tradingRequests[_id].index] == _id);
    }

    function containsTrade(uint _id)
    public
    view
    returns (bool isContained)
    {
        if (tradesIndex.length == 0) return false;
        return (tradesIndex[trades[_id].index] == _id);
    }

    function create(
        address _consumer,
        address _broker,
        uint _product,
        uint _startTime,
        uint _endTime,
        uint _cost,
        uint _frequency
    )
    private
    returns (uint id)
    {
        require(deviceContract.existsByAddress(msg.sender), 'provider does not exist');
        require(deviceContract.existsByAddress(_consumer), 'consumer does not exist');
        require(brokerContract.existsByAddress(_broker), 'broker does not exist');
        require(productContract.isOwnedByDevice(_product, msg.sender), 'not authorized to create trade');

        uint nextId = tradeIdCounter++;

        Trade storage trade = trades[nextId];

        trade.id = nextId;
        trade.provider = msg.sender;
        trade.consumer = _consumer;
        trade.broker = _broker;
        trade.product = _product;
        trade.startTime = _startTime;
        trade.endTime = _endTime;
        trade.cost = _cost;
        trade.frequency = _frequency;
        trade.settlementContract = address(
            new SettlementContract(
                trade.provider,
                trade.consumer,
                trade.broker,
                trade.startTime,
                trade.endTime,
                trade.cost,
                trade.frequency,
                address(ratingContract)
            )
        );
        trade.index = tradesIndex.push(trade.id) - 1;

        deviceContract.addTrade(trade.consumer, trade.id);
        productContract.addTrade(trade.product, trade.id);
        brokerContract.addTrade(trade.broker, trade.id);
        ratingContract.addSettlementContract(trade.settlementContract);

        emit CreatedTrade(trade.broker, trade.id);

        return trade.id;
    }

    function removeTradingRequest(uint _id)
    private
    {
        require(containsTradingRequest(_id), 'request does not exist');

        TradingRequest memory deletedRequest = tradingRequests[_id];
        TradingRequest storage switchedRequest = tradingRequests[
        tradingRequestsIndex[tradingRequestsIndex.length - 1]
        ];

        tradingRequestsIndex[deletedRequest.index] = tradingRequestsIndex[switchedRequest.index];
        switchedRequest.index = deletedRequest.index;
        tradingRequestsIndex.length--;

        delete tradingRequests[_id];

        deviceContract.removeTradingRequest(deletedRequest.consumer, deletedRequest.id);
        productContract.removeTradingRequest(deletedRequest.product, deletedRequest.id);
    }

}
