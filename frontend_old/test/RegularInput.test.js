import { mount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';
import RegularInput from '../components/RegularInput.vue';

describe('RegularInput.vue', () => {
  it('modelValue should be updated', async () => {
    const wrapper = mount(RegularInput, {
      props: {
        modelValue: null,
        'onUpdate:modelValue': (e) => wrapper.setProps({ modelValue: e }),
      },
    });

    await wrapper.find('input').setValue('test');
    expect(wrapper.props('modelValue')).toBe('test');
  });

  it('emits an update:modelValue when changed', () => {
    const wrapper = mount(RegularInput);

    wrapper.find('input').trigger('input');

    expect(wrapper.emitted()).toHaveProperty('update:modelValue');
  });

  it('default input placeholder to be "Input"', () => {
    const wrapper = mount(RegularInput);

    const placeholder = wrapper.find('input').attributes('placeholder');
    expect(placeholder).toContain('Input');
  });

  it('default input placeholder to be "Input"', () => {
    const wrapper = mount(RegularInput, {
      props: {
        placeholder: 'Custom Placeholder',
      },
    });

    const placeholder = wrapper.find('input').attributes('placeholder');
    expect(placeholder).toContain('Custom Placeholder');
  });

  it('has form--minimalistic class when minimalistic prop set', () => {
    const wrapper = mount(RegularInput, {
      props: {
        minimalistic: true,
      },
    });
    expect(wrapper.classes()).toContain('form--minimalistic');
  });
});
