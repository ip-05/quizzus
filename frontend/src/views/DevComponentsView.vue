<script setup lang="ts">
import { computed, reactive, ref } from 'vue';
import ButtonSecondary from '../components/buttons/ButtonSecondary.vue';
import InputText from '../components/inputs/InputText.vue';
import Textarea from '../components/inputs/Textarea.vue';
import Select from '../components/inputs/Select.vue';
import Switch from '../components/controls/Switch.vue';
import Radio from '../components/controls/Radio.vue';
import Checkbox from '../components/controls/Checkbox.vue';
import ButtonPrimary from '@/components/buttons/ButtonPrimary.vue';
import ButtonOption from '@/components/buttons/ButtonOption.vue';
import { Icons } from '@/types';
import ButtonLarge from '@/components/buttons/ButtonLarge.vue';
import GameOptionButton from '@/components/buttons/GameOptionButton.vue';
import GameRoundButton from '@/components/buttons/GameRoundButton.vue';

const inputValue = ref();
const testareaValue = ref();
const selectValue = ref();
const switchValue = ref();
const radioValue = ref();
const checkboxOneValue = ref(true);
const checkboxManyValue = ref(['red', 'green']);

const buttonIcon = ref<Icons>('star-black');
const changeToStarIcon = () => {
  buttonIcon.value = buttonIcon.value === 'star-black' ? 'star-yellow' : 'star-black';
};

const active = ref('blue');
</script>

<template>
  <main>
    <h2>Components</h2>
    <div class="section">
      <h4>Buttons</h4>

      <div>
        <ButtonPrimary icon="star-black"></ButtonPrimary>
        <ButtonPrimary disabled></ButtonPrimary>

        <ButtonPrimary icon="star-black" :style="'bordered'"></ButtonPrimary>
        <ButtonPrimary icon="star-grey" :style="'bordered'" disabled></ButtonPrimary>

        <ButtonPrimary icon="star-white" :style="'action'"></ButtonPrimary>
        <ButtonPrimary icon="star-white" :style="'action'" disabled></ButtonPrimary>
      </div>

      <div>
        <ButtonOption :icon="buttonIcon" @click="changeToStarIcon"></ButtonOption>
        <ButtonOption icon="star-grey" disabled></ButtonOption>
      </div>

      <div>
        <ButtonLarge icon="google" style="width: 320px">Sign in with Google</ButtonLarge>
        <ButtonLarge icon="google" icon-right="google">Google</ButtonLarge>
        <ButtonLarge style="width: 320px">Sign in with Telegram</ButtonLarge>
        <ButtonLarge icon="discord" style="width: 320px" disabled>Sign in with Discord</ButtonLarge>
      </div>

      <div>
        <GameOptionButton></GameOptionButton>
      </div>

      <div>
        <GameRoundButton icon="profile"></GameRoundButton>
      </div>
    </div>

    <div class="section">
      <h3>Inputs</h3>
      <h4>InputText: {{ inputValue }}</h4>
      <InputText v-model="inputValue"></InputText>
    </div>

    <div class="section">
      <h4>Textarea: {{ testareaValue }}</h4>
      <Textarea v-model="testareaValue"></Textarea>
    </div>

    <div class="section">
      <h4>Select {{ selectValue }}</h4>
      <Select v-model="selectValue" :options="[1, '2']" placeholder="Select options"></Select>
    </div>

    <div class="section">
      <h4>Switch {{ switchValue }}</h4>
      <Switch v-model:checked="switchValue" size="medium">Balls</Switch>
    </div>

    <div class="section">
      <h4>Radio {{ radioValue }}</h4>
      <div class="balls">
        <template v-for="name in ['red', 'blue', 'green', 'yellow']" :key="name">
          <Radio v-model="radioValue" :label="name" :value="name" :color="name" class="ball">
            <template #content>
              <GameOptionButton :active="radioValue === name" :color="name"
                >{{ name }}
              </GameOptionButton>
            </template>
          </Radio>
        </template>
      </div>
      <Radio v-model="radioValue" label="Foo" value="foo" color="red"> </Radio>
      <Radio v-model="radioValue" label="Bar" value="bar" color="blue"> </Radio>
      <Radio v-model="radioValue" label="Baz" value="baz" color="green"> </Radio>
      <Radio v-model="radioValue" label="Baz" value="baz4" color="yellow"> </Radio>
      <Radio v-model="radioValue" label="Baz" value="baz5"> </Radio>
      <!--
      <Radio v-model="radioValue" label="Foo" value="foo" color="red">
        <template #content>
          <GameOptionButton :active="radioValue" color="red">Option 1</GameOptionButton>
        </template>
      </Radio>
      <Radio v-model="radioValue" label="Bar" value="bar" color="blue">
        <template #content
          ><GameOptionButton :active="radioValue" color="blue">Option 1</GameOptionButton></template
        >
      </Radio>
      <Radio v-model="radioValue" label="Baz" value="baz" color="green">
        <template #content
          ><GameOptionButton :active="radioValue" color="green"
            >Option 1</GameOptionButton
          ></template
        >
      </Radio>
      <Radio v-model="radioValue" label="Baz" value="baz4" color="yellow">
        <template #content
          ><GameOptionButton :active="radioValue" color="yellow"
            >Option 1</GameOptionButton
          ></template
        >
      </Radio>
      <Radio v-model="radioValue" label="Baz" value="baz5">
        <template #content
          ><GameOptionButton :active="radioValue" color="red">Option 1</GameOptionButton></template
        >
      </Radio> -->
    </div>

    <div class="section">
      <!-- <div>
        <h4>checkboxOneValue {{ checkboxOneValue }}</h4>
        <Checkbox v-model="checkboxOneValue" label="Foo" value="foo" />

        <h4>checkboxManyValue {{ checkboxManyValue }}</h4>
        <Checkbox v-model="checkboxManyValue" value="1" color="blue" />
        <Checkbox v-model="checkboxManyValue" value="2" color="red" />
        <Checkbox v-model="checkboxManyValue" value="3" color="green" />
        <Checkbox v-model="checkboxManyValue" value="4" color="yellow" />
      </div> -->

      <div>
        <h4>checkboxManyValue {{ checkboxManyValue }}</h4>
        <template v-for="name in ['red', 'blue', 'green', 'yellow']" :key="name">
          <Checkbox v-model="checkboxManyValue" :label="name" :value="name" :color="name">
            <template #content>
              <GameOptionButton :active="checkboxManyValue.includes(name)" :color="name">{{
                name
              }}</GameOptionButton>
            </template>
          </Checkbox>
        </template>

        <!-- <Checkbox v-model="checkboxOneValue" label="Foo" value="foo">
          <template #content>
            <GameOptionButton :active="checkboxOneValue" color="red">Option 1</GameOptionButton>
          </template>
        </Checkbox>

        <h4>checkboxManyValue {{ checkboxManyValue }}</h4>
        <Checkbox v-model="checkboxManyValue" value="1" color="blue">
          <template #content>
            <GameOptionButton :active="checkboxManyValue" color="red">Option 1</GameOptionButton>
          </template>
        </Checkbox>
        <Checkbox v-model="checkboxManyValue" value="2" color="red">
          <template #content
            ><GameOptionButton :active="checkboxManyValue" color="blue"
              >Option 2</GameOptionButton
            ></template
          >
        </Checkbox>
        <Checkbox v-model="checkboxManyValue" value="3" color="green">
          <template #content
            ><GameOptionButton :active="checkboxManyValue" color="green"
              >Option 3</GameOptionButton
            ></template
          >
        </Checkbox>
        <Checkbox v-model="checkboxManyValue" value="4" color="yellow">
          <template #content
            ><GameOptionButton :active="checkboxManyValue" color="yellow"
              >Option 4</GameOptionButton
            ></template
          >
        </Checkbox> -->
      </div>
    </div>
  </main>
</template>

<style scoped>
.section {
  margin-bottom: 50px;
}

.w320 {
  width: 320px;
}

.balls {
  display: grid;
  gap: 20px;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
}

.ball {
  width: 100%;
}
</style>
