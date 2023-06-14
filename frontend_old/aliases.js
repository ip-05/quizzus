import { resolve } from 'path';

const r = (p) => resolve(__dirname, p);

export const alias = {
  '~~': r('.'),
  '~~/': r('./'),
  '@@': r('.'),
  '@@/': r('./'),
  // ... other aliases
};
