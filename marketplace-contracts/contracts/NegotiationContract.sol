pragma solidity >=0.4.21 <0.6.0;

import "./ProductContract.sol";
import "./DeviceContract.sol";
import "./BiddingContract.sol";

contract NegotiationContract {

    event RequestedNegotiation(
        uint id,
        address indexed requester,
        uint indexed product
    );

    event AcceptedNegotiationRequest(
        uint indexed id,
        uint negotiation
    );

    event DeclinedNegotiationRequest(
        uint indexed id
    );

    struct Negotiation {
        uint id;
        uint index;
        uint product;
        address consumer;
        address biddingContract;
    }

    struct NegotiationRequest {
        uint id;
        uint index;
        uint product;
        address consumer;
    }

    uint private requestIdCounter;
    uint private negotiationIdCounter;

    mapping(uint => NegotiationRequest) private negotiationRequests;
    uint[] private negotiationRequestIndex;

    mapping(uint => Negotiation) private negotiations;
    uint[] private negotiationsIndex;

    DeviceContract private deviceContract;
    ProductContract private productContract;

    constructor(address _deviceContractAddress, address _productContractAddress)
    public
    {
        deviceContract = DeviceContract(_deviceContractAddress);
        productContract = ProductContract(_productContractAddress);
    }

    function requestNegotiation(uint _product)
    external
    returns (uint id)
    {
        require(deviceContract.existsByAddressAndDeleted(msg.sender, false), 'consumer does not exist');
        require(productContract.existsByIdAndDeleted(_product, false), 'product does not exist');

        uint nextId = requestIdCounter++;
        NegotiationRequest storage negotiationRequest = negotiationRequests[nextId];

        negotiationRequest.id = nextId;
        negotiationRequest.product = _product;
        negotiationRequest.consumer = msg.sender;
        negotiationRequest.index = negotiationRequestIndex.push(negotiationRequest.id) - 1;

        deviceContract.addNegotiationRequest(negotiationRequest.consumer, negotiationRequest.id);
        productContract.addNegotiationRequest(negotiationRequest.product, negotiationRequest.id);

        emit RequestedNegotiation(negotiationRequest.id, negotiationRequest.consumer, negotiationRequest.product);
        return negotiationRequest.id;
    }

    function acceptNegotiationRequest(uint _id)
    external
    returns (uint id)
    {
        require(existsNegotiationRequestById(_id), 'request does not exist');

        NegotiationRequest memory negotiationRequest = negotiationRequests[_id];

        require(
            productContract.isOwnedByDevice(negotiationRequest.product, msg.sender),
            'not authorized to accept request'
        );

        removeNegotiationRequest(_id);
        id = createNegotiation(negotiationRequest.consumer, negotiationRequest.product);

        emit AcceptedNegotiationRequest(_id, id);
        return id;
    }

    function declineNegotiationRequest(uint _id)
    external
    returns (uint id)
    {
        require(existsNegotiationRequestById(_id), 'request does not exist');

        NegotiationRequest memory negotiationRequest = negotiationRequests[_id];

        require(
            productContract.isOwnedByDevice(negotiationRequest.product, msg.sender),
            'not authorized to decline request'
        );

        removeNegotiationRequest(negotiationRequest.id);
        emit DeclinedNegotiationRequest(_id);
        return _id;
    }

    function findNegotiationRequestByIndex(uint _index)
    external
    view
    returns (
        uint id,
        uint product,
        address consumer
    )
    {
        require(_index >= 0 && _index < negotiationRequestIndex.length, 'request does not exist');

        NegotiationRequest memory negotiationRequest = negotiationRequests[negotiationRequestIndex[_index]];

        id = negotiationRequest.id;
        product = negotiationRequest.product;
        consumer = negotiationRequest.consumer;
    }

    function findNegotiationByIndex(uint _index)
    external
    view
    returns (
        uint id,
        address consumer,
        uint product,
        address biddingContract
    )
    {
        require(_index >= 0 && _index < negotiationsIndex.length, 'request does not exist');

        Negotiation memory negotiation = negotiations[negotiationsIndex[_index]];

        id = negotiation.id;
        product = negotiation.product;
        consumer = negotiation.consumer;
        biddingContract = negotiation.biddingContract;
    }

    function findNegotiationRequestById(uint _id)
    external
    view
    returns (
        uint id,
        uint product,
        address consumer
    )
    {
        require(existsNegotiationRequestById(_id), 'request does not exist');

        NegotiationRequest memory negotiationRequest = negotiationRequests[_id];

        id = negotiationRequest.id;
        product = negotiationRequest.product;
        consumer = negotiationRequest.consumer;
    }

    function findNegotiationById(uint _id)
    external
    view
    returns (
        uint id,
        address consumer,
        uint product,
        address biddingContract
    )
    {
        require(existsNegotiationById(_id), 'negotiation does not exist');

        Negotiation memory negotiation = negotiations[_id];

        id = negotiation.id;
        product = negotiation.product;
        consumer = negotiation.consumer;
        biddingContract = negotiation.biddingContract;
    }

    function countNegotiationRequests()
    external
    view
    returns (uint counter)
    {
        return negotiationRequestIndex.length;
    }

    function countNegotiations()
    external
    view
    returns (uint counter)
    {
        return negotiationsIndex.length;
    }

    function existsNegotiationRequestById(uint _id)
    public
    view
    returns (bool exists)
    {
        if (negotiationRequestIndex.length == 0) return false;
        return (negotiationRequestIndex[negotiationRequests[_id].index] == _id);
    }

    function existsNegotiationById(uint _id)
    public
    view
    returns (bool exists)
    {
        if (negotiationsIndex.length == 0) return false;
        return (negotiationsIndex[negotiations[_id].index] == _id);
    }

    function createNegotiation(address _consumer, uint _product)
    private
    returns (uint id)
    {
        uint nextId = negotiationIdCounter++;
        Negotiation storage negotiation = negotiations[nextId];

        negotiation.id = nextId;
        negotiation.consumer = _consumer;
        negotiation.product = _product;
        negotiation.biddingContract = address(new BiddingContract(msg.sender, negotiation.consumer));
        negotiation.index = negotiationsIndex.push(negotiation.id) - 1;

        deviceContract.addNegotiation(negotiation.consumer, negotiation.id);
        productContract.addNegotiation(negotiation.product, negotiation.id);

        return negotiation.id;
    }

    function removeNegotiationRequest(uint _id)
    private
    {
        NegotiationRequest memory deletedRequest = negotiationRequests[_id];
        NegotiationRequest storage switchedRequest = negotiationRequests[
        negotiationRequestIndex[negotiationRequestIndex.length - 1]
        ];

        negotiationRequestIndex[deletedRequest.index] = negotiationRequestIndex[switchedRequest.index];
        switchedRequest.index = deletedRequest.index;
        negotiationRequestIndex.length--;

        delete negotiationRequests[_id];
        deviceContract.removeNegotiationRequest(deletedRequest.consumer, deletedRequest.id);
        productContract.removeNegotiationRequest(deletedRequest.product, deletedRequest.id);
    }

}
