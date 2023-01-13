import { mount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';
import TableItem from '../components/TableItem.vue';

describe('RegularButton.vue', () => {
  it('default props displayed correctly', () => {
    const wrapper = mount(TableItem);
    expect(wrapper.find('.place').html()).toContain(0);
    expect(wrapper.find('.points').html()).toContain(0);
    expect(wrapper.find('.name').html()).toContain('Name');
  });

  it('displays custom props', () => {
    const wrapper = mount(TableItem, {
      props: {
        points: 10,
        place: 3,
        name: 'Custom Name',
      },
    });
    expect(wrapper.find('.points').html()).toContain(10);
    expect(wrapper.find('.place').html()).toContain(3);
    expect(wrapper.find('.name').html()).toContain('Custom Name');
  });
});
