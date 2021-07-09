pragma solidity >=0.4.21 <0.6.0;

contract BiddingContract {

    event MadeBid(
        address indexed bidder
    );

    event AcceptedBid(
        address indexed bidder
    );

    event CanceledBidding(
        address indexed bidder
    );

    address private provider;
    address private consumer;

    bool private accepted;
    bool private canceled;

    Bid[] private bids;

    struct Bid {
        uint price;
        uint startTime;
        uint endTime;
    }

    constructor(address _provider, address _consumer)
    public
    {
        provider = _provider;
        consumer = _consumer;
    }

    function makeBid(uint _price, uint _startTime, uint _endTime)
    external
    returns (uint index)
    {
        require(isActive(), 'not active');
        require(msg.sender == nextBidder(), 'not your turn');

        Bid memory bid = Bid(_price, _startTime, _endTime);
        index = bids.push(bid) - 1;
        emit MadeBid(msg.sender);
    }

    function accept()
    external
    {
        require(isActive(), 'not active');
        require(msg.sender == nextBidder(), 'not your turn');
        accepted = true;
        emit AcceptedBid(msg.sender);
    }

    function cancel()
    external
    {
        require(isActive(), 'not active');
        require(msg.sender == nextBidder(), 'not your turn');
        canceled = true;
        emit CanceledBidding(msg.sender);
    }

    function findByIndex(uint _index)
    external
    view
    returns (
        uint price,
        uint startTime,
        uint endTime
    )
    {
        require(_index >= 0 && _index < bids.length, 'bid does not exist');

        Bid memory bid = bids[_index];
        price = bid.price;
        startTime = bid.startTime;
        endTime = bid.endTime;
    }

    function lastBid()
    external
    view
    returns (
        uint price,
        uint startTime,
        uint endTime
    )
    {
        require(bids.length > 0, 'bid does not exist');

        Bid memory bid = bids[bids.length - 1];
        price = bid.price;
        startTime = bid.startTime;
        endTime = bid.endTime;
    }

    function count()
    external
    view
    returns (uint counter)
    {
        return bids.length;
    }

    function isAccepted()
    public
    view
    returns (bool _accepted)
    {
        _accepted = accepted;
    }

    function isCanceled()
    public
    view
    returns (bool _canceled)
    {
        _canceled = canceled;
    }

    function isActive()
    public
    view
    returns (bool _active)
    {
        return !accepted && !canceled;
    }

    function nextBidder()
    private
    view
    returns (address bidder)
    {
        if (bids.length % 2 == 0) return consumer;
        return provider;
    }

}
