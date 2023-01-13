import { mount } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';
import MediumButton from '../components/MediumButton.vue';

describe('RegularButton.vue', () => {
  it('contains default slot text', () => {
    const wrapper = mount(MediumButton);
    console.log(wrapper.html());
    expect(wrapper.text()).toContain('Medium Button');
  });

  it('contains custom slot text', () => {
    const wrapper = mount(MediumButton, {
      slots: {
        default: 'Custom text',
      },
    });
    expect(wrapper.html()).toContain('Custom text');
  });

  it('has button--minimalistic class when minimalistic prop set', () => {
    const wrapper = mount(MediumButton, {
      props: { minimalistic: true },
    });
    expect(wrapper.classes()).toContain('button--minimalistic');
  });

  it('img has img--small when smallIcon prop set', () => {
    const wrapper = mount(MediumButton, {
      props: { smallIcon: true },
    });
    const img = wrapper.find('.img');
    expect(img.classes('img--small')).toBe(true);
  });
});
