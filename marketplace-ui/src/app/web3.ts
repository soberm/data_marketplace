import {InjectionToken} from '@angular/core';
import Web3 from 'web3';

const Contract = require('@truffle/contract');

import UserContractAbi from './../../abi/UserContract.json';
import DeviceContractAbi from './../../abi//DeviceContract.json';
import ProductContractAbi from './../../abi/ProductContract.json';
import BrokerContractAbi from './../../abi/BrokerContract.json';
import TradingContractAbi from './../../abi/TradingContract.json';
import SettlementContractAbi from './../../abi/SettlementContract.json';

export const WEB3 = new InjectionToken<Web3>('web3', {
  providedIn: 'root',
  factory: () => {
    try {
      const key = 'ethereum';
      const provider = (key in window) ? window[key] : Web3.givenProvider;
      return new Web3(provider);
    } catch (err) {
      throw new Error('Non-Ethereum browser detected. You should consider trying Mist or MetaMask!');
    }
  }
});

export const UserContract = new InjectionToken<any>('userContract', {
  providedIn: 'root',
  factory: () => Contract(UserContractAbi),
});

export const DeviceContract = new InjectionToken<any>('deviceContract', {
  providedIn: 'root',
  factory: () => Contract(DeviceContractAbi),
});

export const ProductContract = new InjectionToken<any>('productContract', {
  providedIn: 'root',
  factory: () => Contract(ProductContractAbi),
});

export const BrokerContract = new InjectionToken<any>('brokerContract', {
  providedIn: 'root',
  factory: () => Contract(BrokerContractAbi),
});

export const TradingContract = new InjectionToken<any>('tradingContract', {
  providedIn: 'root',
  factory: () => Contract(TradingContractAbi),
});

export const SettlementContract = new InjectionToken<any>('settlementContract', {
  providedIn: 'root',
  factory: () => Contract(SettlementContractAbi),
});
