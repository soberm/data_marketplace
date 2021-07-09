import {Location} from './location.enum';

export interface Broker {
  addr: string;
  user: string;
  name: string;
  hostAddr: string;
  location: Location;
  deleted: boolean;
}
