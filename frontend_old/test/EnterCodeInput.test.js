import { mount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';
import { createTestingPinia } from '@pinia/testing';
import EnterCodeInput from '../components/EnterCodeInput.vue';

import { useGameStore } from '../stores/game';

describe('EnterCodeInput.vue', () => {
  it('modelValue should be updated', async () => {
    const wrapper = mount(EnterCodeInput, {
      global: {
        plugins: [createTestingPinia()],
      },
    });

    const gameStore = useGameStore();

    const input = wrapper.find('input');
    await input.setValue('test');
    expect(input.text()).toBe('test');
  });

  // it('emits an update:modelValue when changed', () => {
  //   const wrapper = mount(EnterCodeInput);
  //   wrapper.find('input').trigger('input');

  //   expect(wrapper.emitted()).toHaveProperty('update:modelValue');
  // });

  // it('default input placeholder to be "Input"', () => {
  //   const wrapper = mount(EnterCodeInput);

  //   const placeholder = wrapper.find('input').attributes('placeholder');
  //   expect(placeholder).toContain('Input');
  // });

  // it('default input placeholder to be "Input"', () => {
  //   const wrapper = mount(EnterCodeInput, {
  //     props: {
  //       placeholder: 'Custom Placeholder',
  //     },
  //   });

  //   const placeholder = wrapper.find('input').attributes('placeholder');
  //   expect(placeholder).toContain('Custom Placeholder');
  // });

  // it('has form--minimalistic class when minimalistic prop set', () => {
  //   const wrapper = mount(EnterCodeInput, {
  //     props: {
  //       minimalistic: true,
  //     },
  //   });
  //   expect(wrapper.classes()).toContain('form--minimalistic');
  // });
});
