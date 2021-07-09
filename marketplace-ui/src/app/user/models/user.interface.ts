export interface User {
  addr: string;
  firstName: string;
  lastName: string;
  company: string;
  email: string;
  deleted: boolean;
  devices: string[];
  brokers: string[];
}
