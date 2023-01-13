import { mount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';
import IslandButton from '../components/IslandButton.vue';

describe('RegularButton.vue', () => {
  it('default slot contains button', () => {
    const wrapper = mount(IslandButton, {
      props: {
        type: 'createroom',
      },
    });
    expect(wrapper.find('button').exists()).toBe(true);
  });

  it('contains custom slot element', () => {
    const wrapper = mount(IslandButton, {
      props: {
        type: 'createroom',
      },
      slots: {
        default: '<span>Some inserted element</span>',
      },
    });
    expect(wrapper.find('span').exists()).toBe(true);
  });

  it('has form--join class when prop is join', () => {
    const wrapper = mount(IslandButton, {
      props: {
        type: 'join',
      },
    });
    expect(wrapper.classes()).toContain('form--join');
  });
});
