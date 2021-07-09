import {Pipe, PipeTransform} from '@angular/core';
import {Buffer} from 'buffer';

@Pipe({name: 'hex'})
export class HexPipe implements PipeTransform {
  transform(value: Uint8Array | string): string {
    const buffer = Buffer.from(value as string, 'base64');
    return '0x' + buffer.toString('hex');
  }
}
