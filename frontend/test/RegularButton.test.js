import { mount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';
import RegularButton from '../components/RegularButton.vue';

describe('RegularButton.vue', () => {
  it('contains default slot text', () => {
    const wrapper = mount(RegularButton);
    expect(wrapper.text()).toContain('Regular Button');
  });

  it('contains custom slot text', () => {
    const wrapper = mount(RegularButton, {
      slots: {
        default: 'Custom text',
      },
    });
    expect(wrapper.html()).toContain('Custom text');
  });

  it('has button-active class when :active prop passed', () => {
    const wrapper = mount(RegularButton, {
      props: { active: true },
    });
    expect(wrapper.classes()).toContain('button--active');
  });

  it('has button-disabled class when :disabled prop passed', () => {
    const wrapper = mount(RegularButton, {
      props: { disabled: true },
    });
    expect(wrapper.classes()).toContain('button--disabled');
  });

  it('have both button-active button-disabled classes when :active and :disabled props passed', () => {
    const wrapper = mount(RegularButton, {
      props: { disabled: true, active: true },
    });
    expect(wrapper.classes('button--active')).toBe(true);
    expect(wrapper.classes('button--disabled')).toBe(true);
  });
});
